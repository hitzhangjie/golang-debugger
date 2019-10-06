## History

### version 1 vs. version 2
The first version of DWARF proved to use excessive amounts of storage, and an incompatible successor, DWARF-2, superseded it and added various encoding schemes to reduce data size. DWARF did not immediately gain universal acceptance; for instance, when Sun Microsystems adopted ELF as part of their move to Solaris, they opted to continue using **stabs, in an embedding known as "stabs-in-elf"**. Linux followed suit, and DWARF-2 did not become the default until the late 1990s.

 
### version 3
The DWARF Workgroup of the Free Standards Group released DWARF version 3 in January 2006, adding (among other things) support for C++ namespaces, Fortran 90 allocatable data and additional compiler optimization techniques.

 
### version 4
The DWARF committee published version 4 of DWARF, which offers "improved data compression, better description of optimized code, and support for new language features in C++", in 2010.

 
### version 5
Version 5 of the DWARF format was published in February 2017. It "incorporates improvements in many areas: better data compression, separation of debugging data from executable files, improved description of macros and source files, faster searching for symbols, improved debugging of optimized code, as well as numerous improvements in functionality and performance."

Now golang build tools use Dwarf v4, while gcc has applied some features of Dwarf v5 for C++.

>If you're interested in golang build tools, you can watch this: 
>
>issue: https://github.com/golang/go/issues/26379.
