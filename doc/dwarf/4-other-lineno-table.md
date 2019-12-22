### 5.4.1 Line Number Table 

A source-level debugger will need to know how to associate locations in the source files with the corresponding machine instruction addresses in the executable object or the shared objects used by that executable object. Such an association would make it possible for the debugger user to specify machine instruction addresses in terms of source locations. This would be done by specifying the line number and the source file containing the statement. The debugger can also use this information to display locations in terms of the source files and to single step from line to line, or statement to statement.

Line number information generated for a compilation unit is represented in the **.debug_line** section of an object file and is referenced by a corresponding compilation unit debugging information entry (see Section 3.1.1 in DWARF v4) in the .debug_info section.

The Dwarf line number table contains the mapping between the memory address of executable code of a program and the source lines that corresponds to these address. 

In the simplest form, this can be looked as a matrix with one column contains the instruction address while another column contains the source line triplet (file, line, column). When setting a breakpoint of a source line, query this table to find the first instruction and set a breakpoint. When program has a fault during execution, query current instruction address related source line to analyze it.

As we imagined, if this table were stored with one row one each instruction, this line number table would be too huge. How to compress it? Dwarf encodes it as a sequence of instructions called a line number table program. These instructions are interpreted by a simple finite state machine to recreate the complete line number table. Also, when recreating the complete line number table, only the first machine instruction of each source statement is stored into the table.

By means of this, line number table is compressed. 

