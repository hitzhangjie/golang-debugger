package frame_test

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	"syscall"
	"testing"

	helper "github.com/hitzhangjie/golang-debugger/_helper"
	dwarfhelper "github.com/hitzhangjie/golang-debugger/dwarf/_helper"
	"github.com/hitzhangjie/golang-debugger/dwarf/frame"
	"github.com/hitzhangjie/golang-debugger/proctl"
)

func TestFindReturnAddress(t *testing.T) {
	var (
		testfile, _ = filepath.Abs("../../_fixtures/testnextprog")
		dbframe     = dwarfhelper.GrabDebugFrameSection(testfile, t)
		fdes        = frame.Parse(dbframe)
		gsd         = dwarfhelper.GosymData(testfile, t)
	)

	helper.WithTestProcess(testfile, t, func(p *proctl.DebuggedProcess) {
		testsourcefile := testfile + ".go"
		start, _, err := gsd.LineToPC(testsourcefile, 26)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("start address: %X\n", start)

		_, err = p.Break(uintptr(start))
		if err != nil {
			t.Fatal(err)
		}

		err = p.Continue()
		if err != nil {
			t.Fatal(err)
		}

		regs, err := p.Registers()
		if err != nil {
			t.Fatal(err)
		}

		fde, err := fdes.FDEForPC(start)
		if err != nil {
			t.Fatal(err)
		}

		ret := fde.ReturnAddressOffset(start)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("ret offset: %X\n", ret)
		fmt.Printf("rsp value: %x\n", regs.Rsp)

		addr := uint64(int64(regs.Rsp) + ret)
		fmt.Printf("check addr: offset + rsp: %X\n", addr)
		data := make([]byte, 8)

		syscall.PtracePeekText(p.Pid, uintptr(addr), data)
		addr = binary.LittleEndian.Uint64(data)
		fmt.Printf("expected ret addr: %x\n", addr)

		if addr != 0x400dff {
			t.Fatalf("return address not found correctly, expected %#v got %#v", uintptr(0x400dff), addr)
		}
	})
}
