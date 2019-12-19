## 发展历史

DWARF标准，主要是面向开发者的，如调试信息的生产者、消费者，比如编译器、调试器以及其他希望用源码来描述二进制程序执行的工具。在我们开始动手写一个调试器之前，我们先尽力了解DWARF标准，包括DWARF的历史。

### 版本2 vs. 版本1

The first version of DWARF proved to use excessive amounts of storage, and an incompatible successor, DWARF-2, superseded it and added various encoding schemes to reduce data size. DWARF did not immediately gain universal acceptance; for instance, when Sun Microsystems adopted ELF as part of their move to Solaris, they opted to continue using **stabs, in an embedding known as "stabs-in-elf"**. Linux followed suit, and DWARF-2 did not become the default until the late 1990s.

DWARF的第一版，其对应的调试信息占磁盘空间很大，并且与DWARF的第二版是不兼容的，DWARF2比DWARF1成功，也添加了各种各样的编码格式降低数据尺寸。DWARF没有立即获到广泛认可、同意，

>The representation of information changed from Version 1 to Version 2, so that Version 2 DWARF information is not binary compatible with Version 1 information. To make it easier for consumers to support both Version 1 and Version 2 DWARF information, the Version 2 information has been moved to a different object file section, .debug_info. 


### 版本3 vs. 版本2
The DWARF Workgroup of the Free Standards Group released DWARF version 3 in January 2006, adding (among other things) support for Java, C++ namespaces, Fortran 90 allocatable data and additional optimization techniques for compilers and linkers.

>The return_address_register field in a Common Information Entry record for call frame information is changed to unsigned LEB representation.


### 版本4 vs. 版本3
The DWARF committee published version 4 of DWARF, which offers "improved data compression, better description of optimized code, and support for new language features in C++", in 2010.


### 版本5
Version 5 of the DWARF format was published in February 2017. It "incorporates improvements in many areas: better data compression, separation of debugging data from executable files, improved description of macros and source files, faster searching for symbols, improved debugging of optimized code, as well as numerous improvements in functionality and performance."

Now golang build tools use Dwarf v4, while gcc has applied some features of Dwarf v5 for C++.

> If you're interested in golang build tools, please watch issue: https://github.com/golang/go/issues/26379.