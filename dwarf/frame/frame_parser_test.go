package frame_test

import (
	//"testing"

	//"../../_helper"
	//"../frame"
)

// func TestParse(t *testing.T) {
// 	var (
// 		data = dwarfhelper.GrabDebugFrameSection("../../_fixtures/testprog", t)
// 		fe   = frame.Parse(data)[0]
// 		ce   = fe.CIE
// 	)

// 	if ce.Length != 16 {
// 		t.Error("Length was not parsed correctly, got ", ce.Length)
// 	}

// 	if ce.Version != 0x3 {
// 		t.Fatalf("Version was not parsed correctly expected %#v got %#v", 0x3, ce.Version)
// 	}

// 	if ce.Augmentation != "" {
// 		t.Fatal("Augmentation was not parsed correctly")
// 	}

// 	if ce.CodeAlignmentFactor != 0x1 {
// 		t.Fatal("Code Alignment Factor was not parsed correctly")
// 	}

// 	if ce.DataAlignmentFactor != -4 {
// 		t.Fatalf("Data Alignment Factor was not parsed correctly got %#v", ce.DataAlignmentFactor)
// 	}

// 	if fe.Length != 32 {
// 		t.Fatal("Length was not parsed correctly, got ", fe.Length)
// 	}

// }


