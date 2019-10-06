### 5.4.1 Line Number Table 

The Dwarf line number table contains the mapping between the memory address of executable code of a program and the source lines that corresponds to these address. 

In the simplest form, this can be looked as a matrix with one column contains the instruction address while another column contains the source line triplet (file, line, column). When setting a breakpoint of a source line, query this table to find the first instruction and set a breakpoing. When program has a fault during execucation, query current instruction address related source line to analyze it.

As we imagined, if this table were stored with one row one each instruction, this line number table would be too huge. How to compress it? Dwarf encodes it as a sequence of instructions called a line number table program. These instructions are interpreted by a simple finite state machine to recreate the complete line number table. Also, when recreating the complete line number table, only the first machine instruction of each source statement is stored into the table.

By means of this, line number table is compressed. 

