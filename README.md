# Golang Debugger

## 关于本书

Hi，我是一名开发者，对未知的东西感到好奇，喜欢求根问底，不喜欢模棱两可、得过且过。自己平时对计算机硬件、操作系统、开发工具链的东西也多少有点涉猎，有一次在使用delve的过程中，突发奇想，希望从一个调试器的角度入口来一窥计算机世界的秘密。

为什么要从调试器角度入手？
- 调试过程，并不只是调试器的工作，也涉及到到了源码、编译器、连接器、调试信息标准，因此从调试器视角来看，它看到的是一连串的协作过程，可以给开发者更宏观的视角来审视软件开发的位置，也为开发者更全面地认识我国的IT产业提供了一个窗口；
- 调试标准，调试信息格式有多种标准，在了解调试信息标准的过程中，可以更好地理解处理器、操作系统、编程语言等的设计思想，如果能结合开源调试器学习还可以了解、验证某些语言特性的设计实现；
- 调试需要与操作系统交互来实现，调试给了一个更加直接、快速的途径让我们一窥操作系统的工作原理，如任务调度、信号处理、虚拟内存管理等。操作系统离我们那么近但是在认识上离我们又那么远，加强操作系统知识的普及程度对于我们构建更加全面、立体化的IT产业链也有帮助；
- 此外，调试器是每个开发者都接触过的常用工具，我也希望借此机会剖析下调试器的常用功能的设计实现、调试的一些技巧，也领略下调试信息标准制定者的高屋建瓴的设计思想，站在巨人的肩膀上体验标准的美的一面。

## 项目介绍

该项目“**golang debugger**”，是一款面向go语言的调试器，现在业界已经有针对go语言的调试器了，如gdb、dlv等等，那么为什么还要从头再开发一款调试器呢？项目初衷并不是为了开发一款新的调试器，现在上也不是。

我的初衷希望从调试器为切入点，将作者多年以来掌握的知识进行融会贯通，这里的内容涉及go语言本身（类型系统、协程调度）、编译器与调试器的协作（DWARF）、操作系统内核（虚拟内存、任务调度、系统调用、指令patch）以及处理器相关指令等诸多内容。

简言之，就是希望能从开发一个go语言调试器作为入口切入，帮助初学者快速上手go语言开发，也在循序渐进、拔高过程中慢慢体会操作系统、编译器、调试器、处理器之间的协作过程、加深对计算机系统全局的认识。由于本人水平有限，不可能完全从0开始自研一款调试器，特别是针对go这样一门快速演进中的语言，所以选择了参考开源社区中某些已有的调试器实现gdb、delve作为参考，结合相关规范、标准慢慢钻研的方式。

希望该项目及相关书籍，能顺利完成，也算是我磨练心性、自我救赎的一种方式，最后，如果能对大家确实起到帮助的作用那是再好不过了。

## 开发计划

- ~ - 2018.11.30 完成Linux平台调试器开发
- ~ - 2018.12.31 完成调试器开发文档撰写
- ~ - 2019.01.31 完成《从0开发go调试器》相关章节
- ~ - 2019.02.28 完成《从调试器看go类型系统》相关章节
- ~ - 2019.03.31 完成《从调试器看go调度系统》相关章节
- ~ - 2019.04.30 通读全文，理顺各个章节内容
- ~ - 2019.05.31 与出版社进行沟通，决定是否出版
    - 能出版就出版
    - 出版不了就做成免费的电子书分发

>备注：项目已经成功延期了一年，so sad...  重新制定下项目开发计划

- ~ - 2019.10.06~2019.10.13 回顾下调试标准Dwarf v4的内容
- ~ - 2019.10.14~2019.10.20 基于go v1.12.6+linux，开发指令级调试器
- ~ - 2019.10.21~2019.10.27 熟练掌握go标准库中debug、elf相关操作
- ~ - 2019.10.28~2019.11.03 基于go v1.12.6+linux，开发符号级调试器
    - 实现ELF的解析
    - 实现.debug_info的解析
    - 实现.debug_line的解析
    - ...

# 意见反馈

请邮件联系 `hit.zhangjie@gmail.com`，标题中请注明来意`golang debugger交流`。

