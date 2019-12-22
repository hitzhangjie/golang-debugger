前面描述的调试信息是存储在.debug_info或者.debug_types中的：

- .debug_types中主要存储类型相关的DIEs；
- .debug_info中主要存储变量、常量、可执行代码等相关的DIEs；

还有些其他的调试信息，不是存储在.debug_info、.debug_types中的。



调试的时候，指定一个符号、类型名、指令地址，如何快速找到对应的源码行：

- 比较笨的办法，就是遍历所有的DIEs，去匹配符号名、类型名、指令地址范围是否匹配，但是这样效率太低了；

- 必须想办法加速访问，算法中经常要平衡“时间复杂度”和“空间复杂度”，这里就是个典型的空间换时间的问题，DWARF中设计了几张表用来加速查询。

  

To make lookups of program entities (data objects, functions and types) by name or by address faster, a producer of DWARF information may provide three different types of tables containing information about the debugging information entries owned by a particular compilation unit entry in a more condensed format.

- Lookup by Name, two tables are maintained in separate object file sections named `.debug_pubnames` for objects and functions, and `.debug_pubtypes` for types. Each table consists of sets of variable length entries. Each set describes the names of global objects and functions, or global tpes, respectively, whose definitions are repersented by debugging information entries owned by a single compilation unit.
- Lookup by Address, a table is maintained in a separate object file section named `.debug_aranges`. The table consists of sets of variable length entries, each set describing the portion of the program’s address space that is covered by a single compilation unit.

每个编译单元，为了支持更快速地查询，DWARF信息生成的时候，编译器可能会生成3个不同的表用语加速查询，这几个表的数据更加紧凑，比全局遍历所有DIEs会快很多。根据符号名查询用.debug_pubnames，根据类型名查询用.debug_pubtypes，根据地址进行查询用.debug_aranges。







