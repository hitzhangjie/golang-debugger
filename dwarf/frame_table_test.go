package frame

import (
	"debug/elf"
	"debug/gosym"
	"encoding/binary"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	//"github.com/derekparker/dbg/_helper"
	//"github.com/derekparker/dbg/proctl"
	"../helper"
	"../proctl"
	"fmt"
	"unsafe"
)

var testfile string

func init() {
	testfile, _ = filepath.Abs("../_fixtures/testprog")
}

func parseGoSym(t *testing.T, exe *elf.File) *gosym.Table {
	symdat, err := exe.Section(".gosymtab").Data()
	if err != nil {
		t.Fatal(err)
	}

	pclndat, err := exe.Section(".gopclntab").Data()
	if err != nil {
		t.Fatal(err)
	}

	pcln := gosym.NewLineTable(pclndat, exe.Section(".text").Addr)
	tab, err := gosym.NewTable(symdat, pcln)
	if err != nil {
		t.Fatal(err)
	}

	return tab
}

func gosymData(t *testing.T) *gosym.Table {
	f, err := os.Open(testfile)
	if err != nil {
		t.Fatal(err)
	}

	e, err := elf.NewFile(f)
	if err != nil {
		t.Fatal(err)
	}

	return parseGoSym(t, e)
}

func TestFindReturnAddress(t *testing.T) {
	var (
		dbframe = grabDebugFrameSection(testfile, t)
		fdes    = Parse(dbframe)
		gsd     = gosymData(t)
	)

	helper.WithTestProcess("testprog", t, func(p *proctl.DebuggedProcess) {
		//testsourcefile := testfile + ".go"
		testsourcefile := "/home/zhangjie/gg/debug/_fixtures/testprog.go" // use abs path
		start, _, err := gsd.LineToPC(testsourcefile, 9)
		if err != nil {
			t.Fatal(err)
		}

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

		end, _, err := gsd.LineToPC(testsourcefile, 19)
		if err != nil {
			t.Fatal(err)
		}

		ret := fde.ReturnAddressOffset(start)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(fde.CIE.ReturnAddressRegister)

		// why do we need minus 8 from rsp + ret?
		// Because golang CIE uses R16 as the return_address_register, R16 is defined as cfa-8 in x86_64.

		//addr := uint64(int64(regs.Rsp) + ret)
		addr := uint64(int64(regs.Rsp) + ret - int64(unsafe.Sizeof(uintptr(0))))

		data := make([]byte, 8)

		syscall.PtracePeekText(p.Pid, uintptr(addr), data)
		addr = binary.LittleEndian.Uint64(data)

		if addr != end {
			t.Fatalf("return address not found correctly, expected %#v got %#v", end, addr)
		}
	})
}

