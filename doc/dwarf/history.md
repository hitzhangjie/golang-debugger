## History

Dwarf standards, its intended audience is the developers of both producers and consumers of debugging information, typically language compilers, debuggers and other tools that need to interpret a binary program in terms of its original source. So before we set out to writing code, we’d better thoroughly learn Dwarf standards, even its history.

### version 2 vs. version 1

The first version of DWARF proved to use excessive amounts of storage, and an incompatible successor, DWARF-2, superseded it and added various encoding schemes to reduce data size. DWARF did not immediately gain universal acceptance; for instance, when Sun Microsystems adopted ELF as part of their move to Solaris, they opted to continue using **stabs, in an embedding known as "stabs-in-elf"**. Linux followed suit, and DWARF-2 did not become the default until the late 1990s.

>The representation of information changed from Version 1 to Version 2, so that Version 2 DWARF information is not binary compatible with Version 1 information. To make it easier for consumers to support both Version 1 and Version 2 DWARF information, the Version 2 information has been moved to a different object file section, .debug_info. 


### version 3 vs. version 2
The DWARF Workgroup of the Free Standards Group released DWARF version 3 in January 2006, adding (among other things) support for Java, C++ namespaces, Fortran 90 allocatable data and additional optimization techniques for compilers and linkers.

>The return_address_register field in a Common Information Entry record for call frame information is changed to unsigned LEB representation.


### version 4 vs. version 3
The DWARF committee published version 4 of DWARF, which offers "improved data compression, better description of optimized code, and support for new language features in C++", in 2010.


### version 5
Version 5 of the DWARF format was published in February 2017. It "incorporates improvements in many areas: better data compression, separation of debugging data from executable files, improved description of macros and source files, faster searching for symbols, improved debugging of optimized code, as well as numerous improvements in functionality and performance."

Now golang build tools use Dwarf v4, while gcc has applied some features of Dwarf v5 for C++.

>If you're interested in golang build tools, please watch issue: https://github.com/golang/go/issues/26379.
>
>I also test `gcc` on macOS 10.15, it generates dwarf information separately with the executable file.
