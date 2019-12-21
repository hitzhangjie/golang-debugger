DWARF使用一系列的调试信息入口（DIEs）来对源程序进行描述，一个调试信息入口（DIE）或者一组调试信息入口（DIEs）共同对源程序中的实体进行描述。

每个调试信息入口都包含一个标签（tag）以及一系列的属性（attributes）：

- tag指明了当前调试信息入口描述的实体属于哪一种类型，如类型、变量、函数等；
- attribute定义了调试信息入口的一些特征，如函数的返回值类型是int类型;

调试信息入口存储在.debug_info和.debug_types中，后者多是描述一些类型定义，前者描述变量、代码等。

> 如果编译器对调试信息进行了压缩，压缩后的调试信息将存储在目标文件中的”z”前缀的section中，如未压缩的调试信息入口信息对应section是.debug_info，那么压缩后将存储在.zdebug_info中。
