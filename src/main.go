package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"command"
	"proctl"
)

type term struct {
	stdin *bufio.Reader
}

func main() {
	// We must ensure here that we are running on the same thread during
	// the execution of dbg. This is due to the fact that ptrace(2) expects
	// all commands after PTRACE_ATTACH to come from the same thread.
	runtime.LockOSThread()

	var (
		pid  int
		proc string
		t    = newTerm()
		cmds = command.DebugCommands()
	)

	flag.IntVar(&pid, "pid", 0, "Pid of running process to attach to.")
	flag.StringVar(&proc, "proc", "", "Path to process to run and debug.")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	dbgproc := beginTrace(pid, proc)

	for {
		cmdstr, err := t.promptForInput()
		if err != nil {
			die(1, "Prompt for input failed.\n")
		}

		cmdstr, args := parseCommand(cmdstr)

		if cmdstr == "exit" {
			handleExit(t, dbgproc, 0)
		}

		cmd := cmds.Find(cmdstr)
		err = cmd(dbgproc, args...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: %s\n", err)
		}
	}
}

func beginTrace(pid int, proc string) *proctl.DebuggedProcess {
	var (
		err     error
		dbgproc *proctl.DebuggedProcess
	)

	if pid != 0 {
		dbgproc, err = proctl.NewDebugProcess(pid)
		if err != nil {
			die(1, "Could not start debugging process:", err)
		}
	}

	if proc != "" {
		proc := exec.Command(proc)
		proc.Stdout = os.Stdout

		err = proc.Start()
		if err != nil {
			die(1, "Could not start process:", err)
		}

		dbgproc, err = proctl.NewDebugProcess(proc.Process.Pid)
		if err != nil {
			die(1, "Could not start debugging process:", err)
		}
	}

	return dbgproc
}

func handleExit(t *term, dbp *proctl.DebuggedProcess, status int) {
	fmt.Println("Would you like to kill the process? [y/n]")
	answer, err := t.stdin.ReadString('\n')
	if err != nil {
		die(2, err.Error())
	}

	fmt.Println("Detaching from process...")
	err = syscall.PtraceDetach(dbp.Process.Pid)
	if err != nil {
		die(2, "Could not detach", err)
	}

	if answer == "y\n" {
		fmt.Println("Killing process", dbp.Process.Pid)

		err := dbp.Process.Kill()
		if err != nil {
			fmt.Println("Could not kill process", err)
		}
	}

	die(status, "Hope I was of service hunting your bug!")
}

func die(status int, args ...interface{}) {
	fmt.Fprint(os.Stderr, args)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(status)
}

func newTerm() *term {
	return &term{
		stdin: bufio.NewReader(os.Stdin),
	}
}

func parseCommand(cmdstr string) (string, []string) {
	vals := strings.Split(cmdstr, " ")
	return vals[0], vals[1:]
}

func (t *term) promptForInput() (string, error) {
	fmt.Print("dbg> ")

	line, err := t.stdin.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(line, "\n"), nil
}
