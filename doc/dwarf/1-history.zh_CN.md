## 发展历史

DWARF标准，主要是面向开发者的，如调试信息的生产者、消费者，比如编译器、调试器以及其他希望用源码来描述二进制程序执行的工具。在我们开始动手写一个调试器之前，我们先尽力了解DWARF标准，包括DWARF的历史。

### 版本2 vs. 版本1

DWARF v1，其对应的调试信息占磁盘空间很大，并且与DWARF v2是不兼容的，DWARF v2比DWARF v1成功，也添加了各种各样的编码格式压缩数据尺寸。DWARF v2依然没有立即获得广泛的接纳。那时候Sun公司决定采用ELF作为Solaris平台上的文件格式，但是并没有选择将为ELF设计的DWARF作为首选的调试信息格式，而是继续使用了Stabs（stabs in elf）。Linux当时也是一样的选择，直到20世纪90年代才将DWARF作为了默认调试信息格式。

>DWARF调试信息的表示，版本2和版本1相比有些不同，DWARF v2和DWARF v1不是二进制兼容的。为了能够让DWARF信息的消费者依旧能够兼容版本DWARF v1、v2，DWARF v2建议将相应的调试信息存储在对象文件的不同section中，即.debug_info。


### 版本3 vs. 版本2
2006年1月份，Free Standards Group这个组织的DWARF工作组发布了DWARF v3，这个版本增加了对Java、C++ namespace、Fortran 90等的支持，也增加了一些针对编译器、连接器的优化技术。

>如，Common Information Entry （简称CIE）中字段 return_address_register 存储调用栈的返回地址，该字段使用无符号LEB编码算法进行编码，可以有效压缩小整数占用的存储空间。


### 版本4 vs. 版本3
2010年，DWARF委员会发布了DWARF v4，该版本的焦点围绕在改善数据压缩、更好地描述编译器优化后代码、增加对C++新特性的描述支持等。


### 版本5
2017年，DWARF v5发布，该版本在很多方面都做了改善、提升，包括更好的数据压缩、调试信息与可执行程序的分离、对macro和源文件的更好的描述、更快速的符号搜索、对编译器优化后代码的更好描述，以及其他功能、性能上的提升。

Now golang build tools use Dwarf v4, while gcc has applied some features of Dwarf v5 for C++.

> DWARF也是现在go语言工具链使用的调试信息格式，截止到go1.12.10，当前采用的版本是DWARF v4。在C++中，某些编译器如gcc已经开始应用了部分DWARF v5的特性，go语言也有这方面的讨论，如果对此感兴趣，可以关注go语言issue：: https://github.com/golang/go/issues/26379.