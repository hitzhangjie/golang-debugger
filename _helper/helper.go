// Start a debugged process and run testfunc on it.
package helper

import (
	"os/exec"
	"runtime"
	"syscall"
	"testing"

	//"github.com/derekparker/dbg/proctl"
	"../proctl"
)

type testfunc func(p *proctl.DebuggedProcess)

func GetRegisters(p *proctl.DebuggedProcess, t *testing.T) *syscall.PtraceRegs {
	regs, err := p.Registers()
	if err != nil {
		t.Fatal("Registers():", err)
	}

	return regs
}

func WithTestProcess(name string, t *testing.T, fn testfunc) {
	runtime.LockOSThread()
	cmd, err := startTestProcess(name)
	if err != nil {
		t.Fatal("Starting test process:", err)
	}

	pid := cmd.Process.Pid
	p, err := proctl.NewDebugProcess(pid)
	if err != nil {
		t.Fatal("NewDebugProcess():", err)
	}
	defer cmd.Process.Kill()

	fn(p)
}

func startTestProcess(name string) (*exec.Cmd, error) {
	cmd := exec.Command(name)

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}