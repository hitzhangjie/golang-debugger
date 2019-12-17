package frame

import (
	dwarfhelper "github.com/hitzhangjie/golang-debugger/dwarf/_helper"
	"path/filepath"
	"testing"
)

// search FDE for target pc
func TestFDEForPC(t *testing.T) {
	fde1 := &FrameDescriptionEntry{AddressRange: &addrange{begin: 100, end: 200}}
	fde2 := &FrameDescriptionEntry{AddressRange: &addrange{begin: 50, end: 99}}
	fde3 := &FrameDescriptionEntry{AddressRange: &addrange{begin: 0, end: 49}}
	fde4 := &FrameDescriptionEntry{AddressRange: &addrange{begin: 201, end: 245}}

	tree := NewFrameIndex()
	tree.Put(fde1)
	tree.Put(fde2)
	tree.Put(fde3)
	tree.Put(fde4)

	node, ok := tree.Find(Addr(35))
	if !ok {
		t.Fatal("Could not find FDE")
	}

	if node != fde3 {
		t.Fatal("Got incorrect fde")
	}
}

func BenchmarkFDEForPC(b *testing.B) {
	var (
		testfile, _ = filepath.Abs("../../_fixtures/testnextprog")
		dbframe     = dwarfhelper.GrabDebugFrameSection(testfile, b)
		fdes        = Parse(dbframe)
		gsd         = dwarfhelper.GosymData(testfile, b)
	)

	pc, _, _ := gsd.LineToPC("/usr/local/go/src/pkg/runtime/memmove_amd64.s", 33)

	for i := 0; i < b.N; i++ {
		_, _ = fdes.FDEForPC(pc)
	}
}
