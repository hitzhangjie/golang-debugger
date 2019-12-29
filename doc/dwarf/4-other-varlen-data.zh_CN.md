### 5.4.4 Variable Length Data

在整个DWARF中有大量信息使用整数值来表示，从数据段偏移量，到数组或结构体的大小，等等。由于大多数整数值是小整数，用几位就可以表示，因此这意味着数据主要由零组成，对应的bits相当于被浪费了。

DWARF定义了一个可变长度的整数，称为Little Endian Base 128（LEB128为有符号整数，ULEB128为无符号整数），它能够压缩实际占用的字节数，减小编码后的数据量。

Wiki: https://en.wikipedia.org/wiki/LEB128 
