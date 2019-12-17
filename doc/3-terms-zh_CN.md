# Terms

在本书中，我们将介绍编译器，链接器，操作系统，调试器和调试信息标准，软件开发等方面的知识，使用到的术语会非常多。我们在此处列出常见和重要的术语，以便读者可以方便地查找。

| **Term**                   | **Description**                                              |
| :------------------------- | :----------------------------------------------------------- |
| Source                     | 源代码，如go语言编写的源代码                                 |
| Compiler                   | 编译器，编译源代码为目标文件                                 |
| Linker                     | 链接器，将目标文件、共享库、系统启动代码链接到一起构建可执行程序 |
| Debugger                   | 调试器，连到一个正在运行的进程或者装载一个core文件，加载程序或core文件调试符号信息，探查、修改进程运行时状态，如查看内存、寄存器情况 |
| Dwarf                      | Dwarf，是一种调试信息标准，指导编译器将调试信息生成到目标文件中，指导链接器合并存储在多个目标文件中的调试信息，调试器将加载此调试信息。简言之，Dwarf用来协调编译器，链接器和调试器之间的工作 |
| Debugger types             | 通常，调试器可以分为两种类型：指令级调试器和符号级调试器     |
| Instruction level debugger | 指令级调试器，其操作的对象是机器指令。通过处理器指令patch技术就可以实现指令级调试，不需要调试符号表。它仅适用于指令或汇编语言级别的操作，不支持源代码级别的操作。 |
| Symbol level debugger      | 符号级调试器，其操作的对象不仅是机器指令，更重要的是支持源代码级的操作。它可以提取和解析调试符号表，建立内存地址、指令地址和源代码之间的映射关系，支持在源代码语句上设置断点的时候，将其转换为精确的机器指令断点，也支持其他方便的操作 |