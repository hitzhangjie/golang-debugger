# 术语表

软件开发需要一整套的编译构建调试工具链的支持，其中还涉及到操作系统、处理器等多方面的知识，知识面覆盖非常广，有必要整理一个术语表规范术语的选择、使用，也方便读者快速查阅、理解。

| **Term**                     | **Description**                                              |
| ---------------------------- | ------------------------------------------------------------ |
| Source                       | Source code programming in go, etc                           |
| Compiler                     | Build object file based on source                            |
| Linker                       | Link object files, shared libraries, system startup code   to build executable file |
| Debugger                     | Attach a running process or load a core file, extract debugging   information from process or core file, inspect process running state via remapping   between memory address, instruction address and source |
| Dwarf                        | A standard to guide the compiler generating debugging information   in object files, guide linker to merge debugging information stored in several   object files, debugger will extract this debugging information to parse and   understand. Dwarf coordinates the work between compiler, linker and debugger |
| Debugger types               | Debugger can be classified into 2 types: instruction level   debugger and symbol level debugger. |
| Instruction level   debugger | It depends on ptrace syscall, no need of debugging symbols   table. It only works on instruction or assembly language level, it cannot set   a breakpoint on source statement |
| Symbol level   debugger      | It depends on on ptrace syscall, too. Besides, it can extract   and parse debugging symbols table, remap information between memory address,   instruction address and source |