### 5.2.4 Describing Executable Code

#### 5.2.4.1 Functions and SubPrograms

1. subprogram

   Functions (also called subprogram) may has return value or not, Dwarf use DIE DW_AT_subprogram to represent both these two cases. This DIE has a name, a source location triplet, and an attribute which indicates whether the subprogram is external, that is, visible outside the current compilation unit.

2. subprogram address range

   A subprogram DIE has attributes DW_AT_low_pc, DW_AT_high_pc to give the low and high memory addresses that the subprogram occupies. In some cases, maybe the subprogram memory address is continuous or not, if not, there’ll be a list of memory ranges. The low pc is assumed to be the entry point of subprogram unless another one is specified explicitly.

3. subprogram return type

   A subprogram’s return value’s type is described by the attribute DW_AT_type within DIE DW_TAG_subprogram. If no value returned, this attribute doesn’t exist. If return type is defined within the same scope with this subprogram, the return type DIE will also be an sibling of this DIE subprogram.

4. subprogram formal parameters

   A subprogram may have zero or several formal parameters, which are described by DIE DW_TAG_formal_parameter, will be listed after the DIE subprogram as the same order as declared in parameter list, though DIE of parameter type may be interspersed. Mostly, these formal parameters are stored in registers.

5. subprogram variables

   A subprogram body may contains local variables, these variables are described by DIE DW_TAG_variables listing after formal parameters’ DIEs. Mostly, these local variables are allocated in stack. 
6. lexical block

   Most programming language support lexical block, there’s may be some lexcical blocks in subprogram, which can be described DIE DW_TAG_lexcical_block. Lexical block may contain variable and lexical block DIEs, too. 

   Following is an example showing how to describe a C function.

![img](assets/clip_image009.png)

​	Generated Dwarf information is as following: 

![img](assets/clip_image010.png)

Referring to 1)~5) content, this example easy to be understood.

#### 5.2.4.2 Compilation Unit

Most Program contain more than one source file. When building program, each source file are treated as a independent compilation unit, which will be independently compiled to *.o (such as C), then these object files will be linked with system specific startup code and system libraries to build the executable program. 

Dwarf adopts the terminology Compilation Unit from C as the DIE name, DW_TAG_compilation_unit. The DIE contains general information about the compilation, including the directory and name of the file name, the used programming language, producer that generated the Dwarf information and the offsets to help locate the line number and macro information. 

If the compilation unit takes up continuous memory (i.e., it’s loaded into memory in one piece), then there’re values for the low and high memory addresses for the unit, which are low pc and high pc attributes. This helps debugger determine which compilation generate the code (instruction) at particular memory address much more easiliy.  

If the compilation is not continuous, then a list of the memory address that the code takes up is provided by the compiler and linker.

#### 5.2.4.3 Data Encoding 

Dwarf data conceptually is a tree of DIE, DIE may has children or siblings, each DIE may has several attributes. Dwarf data is unwieldly, so it must be compressed to reduce the size, then the compressed data is stored into the object file. 

Dwarf provides serveral methods to compress the data. 

1. Use prefix traversal to flatten the tree. 

   Prefix traversal the Dwarf tree, the accessed tree node DIE is stored. By this way, the links between DIE and its children DIEs and sibling DIEs are eliminated. When reading  Dwarf data, maybe jumping to the next sibling DIE is needed, the sibling DIE can be stored as an attribute in current DIE. 

2. Use abbreviation to avoid store duplicated values.

   Instead of storing the value of the TAG and attribute-value pairs, only an index into a table of abbreviations is stored, followed by attributes codes. Each abbreviation gives the TAG value, a flag indicating whether the DIE has children, and a list of attributes with the type of value it expects. 

Figure 9 is an example of using abbreviation:

![img](assets/clip_image011.png)

Less commonly used are features of Dwarf 3 and 4 which allow references from one compilation unit to the Dwarf data stored in another compilation unit.

