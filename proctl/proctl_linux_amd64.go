// Package proctl provides functions for attaching to and manipulating
// a process during the debug session.
package proctl

import (
	"bytes"
	"debug/dwarf"
	"debug/elf"
	"debug/gosym"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"../dwarf/frame"
	"../dwarf/line"
	"../dwarf/op"
	"unsafe"
	"strings"
)

// Struct representing a debugged process. Holds onto pid, register values,
// process struct and process state.
type DebuggedProcess struct {
	Pid             int
	Regs            *syscall.PtraceRegs
	Process         *os.Process
	ProcessState    *os.ProcessState
	Executable      *elf.File
	Symbols         []elf.Symbol
	GoSymTable      *gosym.Table
	FrameEntries    frame.FrameDescriptionEntries
	DebugLine       *line.DebugLineInfo
	BreakPoints     map[string]*BreakPoint
	TempBreakPoints map[uint64]*BreakPoint
}

// Represents a single breakpoint. Stores information on the break
// point including the byte of data that originally was stored at that
// address.
type BreakPoint struct {
	FunctionName string
	File         string
	Line         int
	Addr         uint64
	OriginalData []byte
}

type BreakPointExistsError struct {
	file string
	line int
	addr uintptr
}

type Variable struct {
	Name  string
	Value string
	Type  string
}

func (bpe BreakPointExistsError) Error() string {
	return fmt.Sprintf("Breakpoint exists at %s:%d at %x", bpe.file, bpe.line, bpe.addr)
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
		Pid:             pid,
		Regs:            &syscall.PtraceRegs{},
		Process:         proc,
		ProcessState:    ps,
		BreakPoints:     make(map[string]*BreakPoint),
		TempBreakPoints: make(map[uint64]*BreakPoint),
	}

	err = debuggedProc.LoadInformation()
	if err != nil {
		return nil, err
	}

	return &debuggedProc, nil
}

// Finds the executable from /proc/<pid>/exe and then
// uses that to parse the following information:
// * Dwarf .debug_frame section
// * Dwarf .debug_line section
// * Go symbol table.
func (dbp *DebuggedProcess) LoadInformation() error {
	err := dbp.findExecutable()
	if err != nil {
		return err
	}

	err = dbp.parseDebugFrame()
	if err != nil {
		return err
	}

	err = dbp.parseDebugLine()
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
		return nil, fmt.Errorf("Registers():%v", err)
	}

	return dbp.Regs, nil
}

// Returns the value of the named symbol.
func (dbp *DebuggedProcess) EvalSymbol(name string) (*Variable, error) {
	data, err := dbp.Executable.DWARF()
	if err != nil {
		return nil, err
	}

	reader := data.Reader()

	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return nil, err
		}

		if entry.Tag == dwarf.TagVariable {
			n, ok := entry.Val(dwarf.AttrName).(string)
			if !ok {
				continue
			}

			if n == name {
				offset, ok := entry.Val(dwarf.AttrType).(dwarf.Offset)
				if !ok {
					continue
				}

				t, err := data.Type(offset)
				if err != nil {
					return nil, err
				}

				instructions, ok := entry.Val(dwarf.AttrLocation).([]byte)
				if !ok {
					continue
				}

				val, err := dbp.extractValue(instructions, t)
				if err != nil {
					return nil, err
				}

				return &Variable{Name: n, Type: t.String(), Value: val}, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find symbol value for %s", name)
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

	if bytes.Equal(originalData, int3) {
		return nil, BreakPointExistsError{f, l, addr}
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
		return nil, fmt.Errorf("No breakpoint currently set for %s", bp.FunctionName)
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

	bp, ok := dbp.PCtoBP(regs.PC() - 1)
	if ok {
		// Clear the breakpoint so that we can continue execution.
		_, err = dbp.Clear(bp.Addr)
		if err != nil {
			return err
		}

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
		return fmt.Errorf("step failed:%v", err.Error())
	}

	return nil
}

// Step over function calls.
func (dbp *DebuggedProcess) Next() error {
	addrs := make([]uint64, 0, 3)
	pc, err := dbp.CurrentPC()
	if err != nil {
		return err
	}

	pc-- // account for breakpoint instruction

	fde, err := dbp.FrameEntries.FDEForPC(pc)
	if err != nil {
		return err
	}

	loc := dbp.DebugLine.NextLocAfterPC(pc)
	addrs = append(addrs, loc.Address)
	if !fde.AddressRange.Cover(loc.Address) {
		// Next line is outside current frame, use return addr.
		addr := dbp.ReturnAddressFromOffset(fde.ReturnAddressOffset(pc))
		loc = dbp.DebugLine.LocationInfoForPC(addr)
		addrs = append(addrs, loc.Address)
	}

	if loc.Delta < 0 {
		// We are likely in a loop, set breakpoints at entry and exit.
		entry := dbp.DebugLine.LoopEntryLocation(loc.Line)
		exit := dbp.DebugLine.LoopExitLocation(loc.Address)
		addrs = append(addrs, entry.Address, exit.Address)
	}

	for _, addr := range addrs {
		bp, err := dbp.Break(uintptr(addr))
		if err != nil {
			if _, ok := err.(BreakPointExistsError); !ok {
				return err
			}

			continue
		}
		dbp.TempBreakPoints[addr] = bp
	}

	err = dbp.Continue()
	if err != nil {
		return err
	}

	if bp, ok := dbp.TempBreakPoints[pc]; ok {
		_, err := dbp.Clear(bp.Addr)
		if err != nil {
			return err
		}
	}

	return nil
}

// Continue process until next breakpoint.
func (dbp *DebuggedProcess) Continue() error {
	//bugfix: step() re-inserted the breakpoint at exit phase,
	// continue() won't go to next breakpoint causing this step().

	// Stepping first will ensure we are able to continue
	// past a breakpoint if that's currently where we are stopped.
	//err := dbp.Step()

	regs, err := dbp.Registers()
	if err != nil {
		return err
	}

	// check whether previous single byte instruction is int3,
	// if it indeed is, restore the instruction which was overwritten by int3.
	bp, ok := dbp.PCtoBP(regs.PC() - 1)
	if ok {
		// Clear the breakpoint so that we can continue execution.
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

	return dbp.handleResult(syscall.PtraceCont(dbp.Pid, 0))
}

func (dbp *DebuggedProcess) CurrentPC() (uint64, error) {
	regs, err := dbp.Registers()
	if err != nil {
		return 0, err
	}

	return regs.Rip, nil
}

// Extracts the value from the instructions given in the DW_AT_location entry.
// We execute the stack program described in the DW_OP_* instruction stream, and
// then grab the value from the other processes memory.
func (dbp *DebuggedProcess) extractValue(instructions []byte, typ interface{}) (string, error) {
	regs, err := dbp.Registers()
	if err != nil {
		return "", err
	}

	fde, err := dbp.FrameEntries.FDEForPC(regs.PC())
	if err != nil {
		return "", err
	}

	fctx := fde.EstablishFrame(regs.PC())
	cfaOffset := fctx.CFAOffset()

	off, err := op.ExecuteStackProgram(cfaOffset, instructions)
	if err != nil {
		return "", err
	}

	offset := uintptr(int64(regs.Rsp) + off)

	switch t := typ.(type) {
	case *dwarf.StructType:
		ty := strings.Split(t.String(), " ")
		switch ty[1] {
		case "string":
			return dbp.readString(offset)
		}
	case *dwarf.ArrayType:
		return dbp.readIntArray(offset, t)
	case *dwarf.IntType:
		return dbp.readInt(offset)
	case *dwarf.FloatType:
		return dbp.readFloat64(offset)
	}

	return "", fmt.Errorf("could not find value for type %s", typ)
}

func (dbp *DebuggedProcess) readString(addr uintptr) (string, error) {
	val, err := dbp.readMemory(addr, 8)
	if err != nil {
		return "", err
	}

	// deref the pointer to the string
	addr = uintptr(binary.LittleEndian.Uint64(val))
	val, err = dbp.readMemory(addr, 16)
	if err != nil {
		return "", err
	}

	i := bytes.IndexByte(val, 0x0)
	val = val[:i]
	str := *(*string)(unsafe.Pointer(&val))
	return str, nil
}

func (dbp *DebuggedProcess) readIntArray(addr uintptr, t *dwarf.ArrayType) (string, error) {
	var (
		number  uint64
		members = make([]uint64, 0, t.Size()/8)
	)

	val, err := dbp.readMemory(addr, uintptr(t.Size()))
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(val)
	for {
		err := binary.Read(buf, binary.LittleEndian, &number)
		if err != nil {
			break
		}

		members = append(members, number)
	}

	str := fmt.Sprintf("%s %d", t.String(), members)

	return str, err
}

func (dbp *DebuggedProcess) readInt(addr uintptr) (string, error) {
	val, err := dbp.readMemory(addr, 8)
	if err != nil {
		return "", err
	}

	n := binary.LittleEndian.Uint64(val)

	return strconv.Itoa(int(n)), nil
}
func (dbp *DebuggedProcess) readFloat64(addr uintptr) (string, error) {
	var n float64
	val, err := dbp.readMemory(addr, 8)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(val)
	binary.Read(buf, binary.LittleEndian, &n)

	return strconv.FormatFloat(n, 'f', -1, 64), nil
}

func (dbp *DebuggedProcess) readMemory(addr uintptr, size uintptr) ([]byte, error) {
	buf := make([]byte, size)

	_, err := syscall.PtracePeekData(dbp.Pid, addr, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
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

func (dbp *DebuggedProcess) parseDebugLine() error {
	debugLine, err := dbp.Executable.Section(".debug_line").Data()
	if err != nil {
		return err
	}

	dbp.DebugLine = line.Parse(debugLine)
	return nil
}

func (dbp *DebuggedProcess) parseDebugFrame() error {
	debugFrame, err := dbp.Executable.Section(".debug_frame").Data()
	if err != nil {
		return err
	}

	dbp.FrameEntries = frame.Parse(debugFrame)
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

// Takes an offset from RSP and returns the address of the
// instruction the currect function is going to return to.
func (dbp *DebuggedProcess) ReturnAddressFromOffset(offset int64) uint64 {
	regs, err := dbp.Registers()
	if err != nil {
		panic("Could not obtain register values")
	}

	retaddr := int64(regs.Rsp) + offset
	data := make([]byte, 8)
	syscall.PtracePeekText(dbp.Pid, uintptr(retaddr), data)
	return binary.LittleEndian.Uint64(data)
}
