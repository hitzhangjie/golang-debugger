DWARF is a widely used, standardized debugging data format. DWARF was originally designed along with Executable and Linkable Format (ELF), although it is independent of object file formats. The name is a medieval fantasy complement to "ELF" that has no official meaning, although the backronym '**Debugging With Attributed Record Formats**' was later proposed.

Dwarf uses **DIE (Debugging Information Entry) with TAG and Attributes** to describe nearly everything, including variables, base datatypes and compound datatypes, functions, compilation unit, etc.

Dwarf uses some **Encoding Methods** to shrink the size of debugging information which maybe large.

Dwarf defines some important data, including **Line Number Table**, **Call Frame Information**, etc. Thanks to this, developers can add breakpoints at source statement level, or use `frame N`, `bt` to traverse the callstack, or do something others like that.

Ah, there're too many smart and important thinkings in Dwarf standard. I cannot list them one by one, if you're interested please read Dwarf standard. Try to read it.

