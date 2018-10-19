package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	f, e := elf.Open("../_fixtures/testprog")
	checkError(e)
	fmt.Println("open ELF file success")
	fmt.Println()

	// read section data by section type
	for i := 0; i <= 18; i++ {
		s := f.SectionByType(elf.SectionType(i))
		fmt.Printf("SectionType: %s, data: %v\n", elf.SectionType(i).String(), s)
	}
	fmt.Println()

	// read section data by section name
	s := f.Section(".debug_info")
	fmt.Printf(".debug_info: %v\n", s)

	s = f.Section(".text")
	fmt.Printf(".text: %v\n\n", s)

	/// read dwarf and traverse each DIE entry
	dbg, e := f.DWARF()
	checkError(e)
	fmt.Println("Read Dwarf data success")
	r := dbg.Reader()
	for {
		entry, _ := r.Next()
		if entry == nil {
			break
		}

		fmt.Printf("DIE: %#v\n", entry)

		time.Sleep(time.Millisecond * 500)

		lr, _ := dbg.LineReader(entry)
		if lr != nil {
			le := dwarf.LineEntry{}
			for {
				e := lr.Next(&le)
				if e == io.EOF {
					break
				}
				fmt.Printf("\t\tline: %#v\n", le)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
	fmt.Println()

	// read symbols
	/*
		syms, e := f.Symbols()
		checkError(e)
		fmt.Printf("Symbols: %v\n\n", syms)

		dsyms, e := f.DynamicSymbols()
		fmt.Printf("Dyn Symbols: %v\n\n", dsyms)

		isyms, e := f.ImportedSymbols()
		fmt.Printf("Import Symbols: %v\n\n", isyms)

		libs, e := f.ImportedLibraries()
		fmt.Printf("Import Libraries: %v\n\n", libs)
	*/

	e = f.Close()
	checkError(e)
	fmt.Println("close ELF file success")

}

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
