// Package proctl provides functions for attaching to and manipulating
// a process during the debug session.
package proctl

import (
	"debug/elf"
	"debug/gosym"
	"fmt"
	"os"
	"syscall"
)

// Struct representing a debugged process. Holds onto pid, register values,
// process struct and process state.
type DebuggedProcess struct {
	Pid          int
	Regs         *syscall.PtraceRegs
	Process      *os.Process
	ProcessState *os.ProcessState
	Executable   *elf.File
	Symbols      []elf.Symbol
	GoSymTable   *gosym.Table
	BreakPoints  map[string]*BreakPoint
}

// Represents a single breakpoint. Stores information on the break
// point include the byte of data that originally was stored at that
// address.
type BreakPoint struct {
	FunctionName string
	File         string
	Line         int
	Addr         uint64
	OriginalData []byte
}

// Returns a new DebuggedProcess struct with sensible defaults.
func NewDebugProcess(pid int) (*DebuggedProcess, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	err = syscall.PtraceAttach(pid)
	if err != nil {
		return nil, err
	}

	ps, err := proc.Wait()
	if err != nil {
		return nil, err
	}

	debuggedProc := DebuggedProcess{
		Pid:          pid,
		Regs:         &syscall.PtraceRegs{},
		Process:      proc,
		ProcessState: ps,
		BreakPoints:  make(map[string]*BreakPoint),
	}

	err = debuggedProc.LoadInformation()
	if err != nil {
		return nil, err
	}

	return &debuggedProc, nil
}

// Finds the executable from /proc/<pid>/exe and then
// uses that to parse the Go symbol table.
func (dbp *DebuggedProcess) LoadInformation() error {
	err := dbp.findExecutable()
	if err != nil {
		return err
	}

	err = dbp.obtainGoSymbols()
	if err != nil {
		return err
	}

	return nil
}

// Obtains register values from the debugged process.
func (dbp *DebuggedProcess) Registers() (*syscall.PtraceRegs, error) {
	err := syscall.PtraceGetRegs(dbp.Pid, dbp.Regs)
	if err != nil {
		return nil, fmt.Errorf("Registers: %v", err)
	}

	return dbp.Regs, nil
}

// Sets a breakpoint in the running process.
func (dbp *DebuggedProcess) Break(addr uintptr) (*BreakPoint, error) {
	var (
		int3         = []byte{0xCC}
		f, l, fn     = dbp.GoSymTable.PCToLine(uint64(addr))
		originalData = make([]byte, 1)
	)

	_, err := syscall.PtracePeekData(dbp.Pid, addr, originalData)
	if err != nil {
		return nil, err
	}

	_, err = syscall.PtracePokeData(dbp.Pid, addr, int3)
	if err != nil {
		return nil, err
	}

	breakpoint := &BreakPoint{
		FunctionName: fn.Name,
		File:         f,
		Line:         l,
		Addr:         uint64(addr),
		OriginalData: originalData,
	}

	fname := fmt.Sprintf("%s:%d", f, l)
	dbp.BreakPoints[fname] = breakpoint

	return breakpoint, nil
}

// Clears a breakpoint.
func (dbp *DebuggedProcess) Clear(pc uint64) (*BreakPoint, error) {
	bp, ok := dbp.PCtoBP(pc)
	if !ok {
		if bp != nil {
			return nil, fmt.Errorf("No breakpoint currently set for %s", bp.FunctionName)
		} else {
			return nil, fmt.Errorf("No breakpoint currently set for addr %#x", pc)
		}
	}
	_, err := syscall.PtracePokeData(dbp.Pid, uintptr(bp.Addr), bp.OriginalData)
	if err != nil {
		return nil, err
	}

	delete(dbp.BreakPoints, fmt.Sprintf("%s:%d", bp.File, bp.Line))

	return bp, nil
}

// Steps through process.
func (dbp *DebuggedProcess) Step() (err error) {
	regs, err := dbp.Registers()
	if err != nil {
		return err
	}

	// check whether previous single byte instruction is int3,
	// if it indeed is, restore the instruction which was overwritten by int3.
	bp, ok := dbp.PCtoBP(regs.PC() - 1)
	if ok {
		_, err = dbp.Clear(bp.Addr)
		if err != nil {
			return err
		}

		// Reset instruction pointer to our restored instruction.
		//regs.Rip -= 1
		//syscall.PtraceSetRegs(dbp.Pid, regs)

		// Reset program counter to our restored instruction.
		regs.SetPC(bp.Addr)
		err = syscall.PtraceSetRegs(dbp.Pid, regs)
		if err != nil {
			return err
		}

		// Restore breakpoint now that we have passed it.
		defer func() {
			_, err = dbp.Break(uintptr(bp.Addr))
		}()
	}

	err = dbp.handleResult(syscall.PtraceSingleStep(dbp.Pid))
	if err != nil {
		return fmt.Errorf("step failed: %v", err.Error())
	}

	return nil
}

// Continue process until next breakpoint.
func (dbp *DebuggedProcess) Continue() error {
	// Stepping first will ensure we are able to continue
	// past a breakpoint if that's currently where we are stopped.
	err := dbp.Step()
	if err != nil {
		return err
	}

	return dbp.handleResult(syscall.PtraceCont(dbp.Pid, 0))
}

func (dbp *DebuggedProcess) handleResult(err error) error {
	if err != nil {
		return err
	}

	ps, err := dbp.Process.Wait()
	if err != nil {
		return err
	}

	dbp.ProcessState = ps

	return nil
}

func (dbp *DebuggedProcess) findExecutable() error {
	procpath := fmt.Sprintf("/proc/%d/exe", dbp.Pid)

	f, err := os.Open(procpath)
	if err != nil {
		return err
	}

	elffile, err := elf.NewFile(f)
	if err != nil {
		return err
	}

	dbp.Executable = elffile

	return nil
}

func (dbp *DebuggedProcess) obtainGoSymbols() error {
	symdat, err := dbp.Executable.Section(".gosymtab").Data()
	if err != nil {
		return err
	}

	pclndat, err := dbp.Executable.Section(".gopclntab").Data()
	if err != nil {
		return err
	}

	pcln := gosym.NewLineTable(pclndat, dbp.Executable.Section(".text").Addr)
	tab, err := gosym.NewTable(symdat, pcln)
	if err != nil {
		return err
	}

	dbp.GoSymTable = tab

	return nil
}

// Converts a program counter value into a breakpoint, if one was set
// for the function containing pc.
func (dbp *DebuggedProcess) PCtoBP(pc uint64) (*BreakPoint, bool) {
	f, l, _ := dbp.GoSymTable.PCToLine(pc)
	bp, ok := dbp.BreakPoints[fmt.Sprintf("%s:%d", f, l)]
	return bp, ok
}

