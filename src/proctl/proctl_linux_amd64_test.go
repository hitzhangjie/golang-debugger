package proctl

import (
	"_helper"
	"bytes"
	"path/filepath"
	"syscall"
	"testing"
)

func dataAtAddr(pid int, addr uint64) ([]byte, error) {
	data := make([]byte, 1)
	_, err := syscall.PtracePeekData(pid, uintptr(addr), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func assertNoError(err error, t *testing.T, s string) {
	if err != nil {
		t.Fatal(s, ":", err)
	}
}

func currentPC(p *DebuggedProcess, t *testing.T) uint64 {
	pc, err := p.CurrentPC()
	if err != nil {
		t.Fatal(err)
	}

	return pc
}

func currentLineNumber(p *DebuggedProcess, t *testing.T) int {
	pc := currentPC(p, t)
	_, l, _ := p.GoSymTable.PCToLine(pc)

	return l
}

func TestAttachProcess(t *testing.T) {
	helper.WithTestProcess("../_fixtures/testprog", t, func(p *DebuggedProcess) {
		if !p.ProcessState.Stopped() {
			t.Errorf("Process was not stopped correctly")
		}
	})
}

func TestStep(t *testing.T) {
	helper.WithTestProcess("../_fixtures/testprog", t, func(p *DebuggedProcess) {
		if p.ProcessState.Exited() {
			t.Fatal("Process already exited")
		}

		regs := helper.GetRegisters(p, t)
		rip := regs.PC()

		err := p.Step()
		assertNoError(err, t, "Step()")

		regs = helper.GetRegisters(p, t)
		if rip >= regs.PC() {
			t.Errorf("Expected %#v to be greater than %#v", regs.PC(), rip)
		}
	})
}

func TestContinue(t *testing.T) {
	helper.WithTestProcess("../_fixtures/continuetestprog", t, func(p *DebuggedProcess) {
		if p.ProcessState.Exited() {
			t.Fatal("Process already exited")
		}

		err := p.Continue()
		assertNoError(err, t, "Continue()")

		if p.ProcessState.ExitStatus() != 0 {
			t.Fatal("Process did not exit successfully")
		}
	})
}

func TestBreakPoint(t *testing.T) {
	helper.WithTestProcess("../_fixtures/testprog", t, func(p *DebuggedProcess) {
		sleepytimefunc := p.GoSymTable.LookupFunc("main.sleepytime")
		sleepyaddr := sleepytimefunc.Entry

		bp, err := p.Break(uintptr(sleepyaddr))
		assertNoError(err, t, "Break()")

		breakpc := bp.Addr + 1
		err = p.Continue()
		assertNoError(err, t, "Continue()")

		regs := helper.GetRegisters(p, t)

		pc := regs.PC()
		if pc != breakpc {
			t.Fatalf("Break not respected:\nPC:%d\nFN:%d\n", pc, breakpc)
		}

		err = p.Step()
		assertNoError(err, t, "Step()")

		regs = helper.GetRegisters(p, t)

		pc = regs.PC()
		if pc == breakpc {
			t.Fatalf("Step not respected:\nPC:%d\nFN:%d\n", pc, breakpc)
		}
	})
}

func TestBreakPointWithNonExistantFunction(t *testing.T) {
	helper.WithTestProcess("../_fixtures/testprog", t, func(p *DebuggedProcess) {
		_, err := p.Break(uintptr(0))
		if err == nil {
			t.Fatal("Should not be able to break at non existant function")
		}
	})
}

func TestClearBreakPoint(t *testing.T) {
	helper.WithTestProcess("../_fixtures/testprog", t, func(p *DebuggedProcess) {
		fn := p.GoSymTable.LookupFunc("main.sleepytime")
		bp, err := p.Break(uintptr(fn.Entry))
		assertNoError(err, t, "Break()")

		int3, err := dataAtAddr(p.Pid, bp.Addr)
		if err != nil {
			t.Fatal(err)
		}

		bp, err = p.Clear(fn.Entry)
		assertNoError(err, t, "Clear()")

		data, err := dataAtAddr(p.Pid, bp.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if bytes.Equal(data, int3) {
			t.Fatalf("Breakpoint was not cleared data: %#v, int3: %#v", data, int3)
		}

		if len(p.BreakPoints) != 0 {
			t.Fatal("Breakpoint not removed internally")
		}
	})
}

func TestNext(t *testing.T) {
	var (
		ln             int
		err            error
		executablePath = "../_fixtures/testnextprog"
	)

	testcases := []struct {
		begin, end int
	}{
		{17, 19},
		{19, 20},
		{20, 22},
		{22, 19},
		{19, 20},
		{20, 22},
		{22, 19},
		{19, 25},
		{25, 26},
		{26, 30},
		{30, 31},
	}

	fp, err := filepath.Abs("../_fixtures/testnextprog.go")
	if err != nil {
		t.Fatal(err)
	}

	helper.WithTestProcess(executablePath, t, func(p *DebuggedProcess) {
		pc, _, _ := p.GoSymTable.LineToPC(fp, testcases[0].begin)
		_, err := p.Break(uintptr(pc))
		assertNoError(err, t, "Break()")
		assertNoError(p.Continue(), t, "Continue()")

		for _, tc := range testcases {
			ln = currentLineNumber(p, t)
			if ln != tc.begin {
				t.Fatalf("Program not stopped at correct spot expected %d was %d", tc.begin, ln)
			}

			assertNoError(p.Next(), t, "Next() returned an error")

			ln = currentLineNumber(p, t)
			if ln != tc.end {
				t.Fatalf("Program did not continue to correct next location expected %d was %d", tc.end, ln)
			}
		}
	})
}

func TestVariableEvaluation(t *testing.T) {
	executablePath := "../_fixtures/testvariables"

	fp, err := filepath.Abs(executablePath + ".go")
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		name    string
		value   string
		varType string
	}{
		{"a1", "foo", "struct string"},
		{"a2", "6", "int"},
		{"a3", "7.23", "float64"},
		{"a4", "[2]int [1 2]", "[2]int"},
		{"a5", "len: 5 cap: 5 [1 2 3 4 5]", "struct []int"},
	}

	helper.WithTestProcess(executablePath, t, func(p *DebuggedProcess) {
		pc, _, _ := p.GoSymTable.LineToPC(fp, 21)

		_, err := p.Break(uintptr(pc))
		assertNoError(err, t, "Break() returned an error")

		err = p.Continue()
		assertNoError(err, t, "Continue() returned an error")

		for _, tc := range testcases {
			variable, err := p.EvalSymbol(tc.name)
			assertNoError(err, t, "Variable() returned an error")

			if variable.Name != tc.name {
				t.Fatalf("Expected %s got %s\n", tc.name, variable.Name)
			}

			if variable.Type != tc.varType {
				t.Fatalf("Expected %s got %s\n", tc.varType, variable.Type)
			}

			if variable.Value != tc.value {
				t.Fatalf("Expected %#v got %#v\n", tc.value, variable.Value)
			}
		}
	})
}
