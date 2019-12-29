### 5.4.3 调用栈信息（Call Frame Information）

#### 5.4.3.1 介绍

调试器通常需要能够查看和修改调用堆栈上任何子例程 “**活动记录（栈帧）**” 的状态。 一个活动记录包括：

- 子例程中的代码位置，该位置要么是调试器获得控制权时程序停止的位置（例如断点），要么是子例程进行调用或被异步事件（例如信号）中断的位置；
- 在堆栈上分配的内存区域称为“调用帧”。调用帧由堆栈上的地址标识。 我们将此地址称为“**Canonical Frame Address （规范帧地址）**”或“**CFA**”。 通常，CFA被定义为前一栈帧在调用当前子例程时的堆栈指针的值（可能与进入当前帧时的值不同）;
- 子例程在代码位置使用的一组寄存器；

通常，在调用子例程（函数）会指定一组寄存器将其状态进行保存。 如果被调用的子例程要使用一个寄存器，它就要在函数入口处保存该寄存器的原始值到栈帧中，并在退出时将其恢复。 

- 在调用栈上分配栈空间并执行保存寄存器状态任务的这部分代码，称为函数序言（prologue）；
- 执行寄存器状态恢复并销毁栈帧任务的这部分代码称为函数后记（epilogue）。

通常，序言代码实际上在子例程的开头，而后记代码在结尾。

为了能够查看或修改不在调用栈顶部的活动记录，调试器必须“虚拟地展开（virtually unwind）”活动记录（栈帧）堆栈，直到找到感兴趣的活动记录为止。调试器分步展开堆栈。 从当前活动记录（栈帧）开始，它实际上恢复了当前活动记录（栈帧）在函数入口处时保存的所有寄存器，并计算了调用方的CFA和代码位置。 这在逻辑上，效果等同于当前子例程通过return返回其调用方。 我们说调试器是“虚拟地展开”堆栈的，因为目标进程的实际状态是不会变的。

> 联想下gdb调试的过程，通过bt可以看到调用栈，然后通过frame N来选择指定的栈帧，这个时候就是一个虚拟地展开调用栈的过程，目标栈帧中的寄存器状态被恢复，为什么说是虚拟地展开堆栈？因为当我们调试器控制程序恢复执行时，还是会按照CFI表执行一遍指令回到栈顶，而寄存器的状态又回到了frame N选择栈帧之前的状态，目标进程的实际状态并没有发生改变。

#### 5.4.3.2 架构无关编码方式

展开堆栈操作，需要知道寄存器的保存位置以及如何计算调用方的CFA和代码位置。在考虑体系结构无关的信息编码方式时，有些特殊事项需要考虑：

- 子例程（函数）序言和后记代码，并不总是位于子例程的开头和结尾这两个不同的块中。通常子例程后记部分代码会被复制到每次return返回操作的地方。有时，编译器也会将寄存器保存、取消保存操作分割开，并将它们移到子例程代码需要用到它们的位置；
- 编译器会使用不同的方式来管理调用栈，有时是通过一个栈指针，有时可能不是；
- 随着子例程序言和后记部分代码的执行，计算CFA的算法也会发生变化（根据定义，CFA值不变）；
- 一些子例程调用是没有调用栈帧的（如可能通过“尾递归”优化掉了）；
- 有时将一个寄存器的值保存在另一个寄存器中，但是后者可能按照惯例是不需要在子例程序言中存储的；
- 某些体系结构具有特殊的指令，这些指令可以在一条指令中执行部分或全部的寄存器管理，而在堆栈上留下一些特殊信息来指示寄存器该如何保存；
- 一些体系结构处理返回地址值比较特殊。 例如，在有的体系结构中，调用指令可确保调用地址低两位为零，而返回指令将忽略这些位。 这留下了两个存储位，可供其他用途使用，必须对其进行特殊处理。

#### 5.4.3.3 调用栈帧信息（CFI）结构

DWARF定义了独立于体系结构的基本要素来支持“虚拟展开（virtually unwind）”调用栈，这些基础要素能够记录子例程调用期间如何保存和恢复寄存器的状态。对于某些特定机器，其可能拥有些诸如体系结构特定的ABI委员会、硬件供应商或编译器生产商定义的信息，需要借助这些信息对DWARF基本要素进行补充。

抽象地，此机制描述了具有以下结构的非常大的表（CFI信息表）：

<img src="assets/image-20191229130341692.png" alt="image-20191229130341692" style="zoom:5%;" />

- 第一列，指令地址。表示程序中包含代码的每个位置的地址（在共享对象中，这是一个相对于对象的偏移量）；
- 第二列，CFA（Canonical Frame Address），规范栈帧地址，调用方调用当前子例程时栈指针值；
- 其他列，各寄存器对应的虚拟展开规则（virtual unwinding rules）；

CFA列，定义了计算规范栈帧地址值的规则，它可以是寄存器、带符号偏移量组合在一起进行计算得到，也可以由求值的DWARF表达式计算得到。

其余列由寄存器编号标记。 其中包括一些在某些架构上具有特殊名称的寄存器，例如PC和堆栈指针寄存器（针对特定体系结构的寄存器的实际映射由扩展器augmenter定义）。寄存器列的描述，包含是否已保存给定寄存器，以及如何在前一栈帧中查找寄存器值的规则。

这里的寄存器规则，包括：

- undefined，该规则表示对应寄存器在前一个栈帧中没有可恢复的值。通常是，在调用callee的时候没有对相关寄存器的状态进行保存；
- same value，该规则表示对应寄存器的值与前一个栈帧中寄存器的值相同。通常是，在调用callee的时候对相关寄存器的状态进行了保存，但是并没有进行修改；
- offset(N)，该规则表示对应寄存器的值被保存在CFA+N对应的地址处，CFA就是当前的CFA值，N是有符号偏移量；
- val_offset(N)，该规则表示对应寄存器的值就是CFA+N的值，CFA就是当前的CFA值，N是有符号偏移量；
- register(R)，该规则表示对应寄存器的值，被保存在另一个寄存器R中；
- expression(E)，该规则表示对应寄存器的值，保存在DWARF表达式E对应的内存地址中；
- val_expression(E)，该规则表示对应寄存器的值，就是DWARF表达式E的值；
- architectural，该规则不是当前规范内的定义，它由其他增强器（augmenter）定义；

如果按照上述CFI信息表对表结构进行定义、数据编码，则该表将非常大、数据量会非常多。实际情况是，该表中相邻各行，他们在多数列上的值是相同的，因此我们可以只记录它们的差异。通过仅记录程序中各个子例程的起始地址开始的差异，可以非常紧凑地表示整个表。

上述CFI信息表（virtual unwind rules）被编码在 “**.debug_frame**” section 中。 .debug_frame节中的条目相对于该节的开头按地址大小的倍数对齐，并以两种形式出现：

- 公共信息条目（Common Information Entry, CIE）；
- 帧描述条目（Frame Descriptor Entry, FDE）；

> 如果函数的代码段地址范围不是连续的，可能存在多个CIEs和FDEs。

##### 5.4.3.3.1 公共信息条目（Common Information Entry）

A Common Information Entry holds information that is shared among many Frame Description Entries. There is at least one CIE in every non-empty .debug_frame section. A CIE contains the following fields, in order:

每个公共信息条目CIE的信息，可能会被很多帧描述条目FDE所共享。每个非空的.debug_frame section中至少包含一个CIE，每个CIE都包含如下字段，按照字段存储顺序依次是：

1. length (初始长度)，常量，指明了该CIE结构的大小（字节数量），不包含该字段本身。length字段所占字节数，加上length的值，必须是按照address size对齐；

2. CIE_id (4字节或8字节)，常量，用语CIEs、FDEs；

3. version(ubyte)，版本号，该值与CFI信息有关，与DWARF版本无关；

4. augmentation (UTF-8字符串)

  null结尾的UTF-8字符串，用于标志当前CIE和使用它的FDEs的扩展信息，如果一个reader遇到一个未知的augmentation字符串，只可以读取如下字段；

  - CIE: length, CIE_id, version, augmentation
  - FDE: length, CIE_pointer, initial_location, address_range

  如果没有augmentation，该字段值就是0，一个字节。augmentation字符串，允许用户向CIE、FDE添加一些目标机器相关的信息，来指导如何解开一个堆栈。例如，动态分配的数据可在函数退出时进行释放，可以将这些信息作为augmentation信息。.debug_frame只使用UTF-8编码。

5. address_size (ubyte)，该CIE中以及使用该CIE的其他FDEs中，目标机器地址占用几个字节，如果该frame存在一个编译单元，其中的address size必须与这里的address size相同；

6. segment_size (ubyte)，该CIE中以及使用该CIE的其他FDEs中，段选择器占用几个字节；

7. code_alignment_factor (unsigned LEB128)，常量，指令地址偏移量 = operand * code_alignment_factor；
8. data_alignment_factor (signed LEB128)，常量，偏移量 = operand * data_alignment_factor；
9. return_address_register (unsigned LEB128)，常量，指示返回地址存储在哪里，可能是物理寄存器或内存
10. initial_instructions (array of ubyte)，一系列rules，用于指示如何创建CFI信息表的初始设置；
  在执行initial instructions之前，所有列的默认生成规则都是undefined，不过, ABI authoring body 或者 compilation system authoring body 也可以为某列或者所有列指定其他的默认规则；
11. padding (array of ubyte)，字节填充，通过DW_CFA_nop指令填充结构体，使CIE结构体大小满足length要求，length值加字段字节数必须按照address size对齐；
  

##### 5.4.3.3.2 帧描述条目（Frame Descriptor Entry）

An FDE contains the following fields, in order:

一个FDE包含如下字段，按照字段顺序依次如下：

1. length (初始长度)，常量，指明该函数对应header以及instruction流的字节数量，不包含该字段本身。length字段大小（字节数），加上length值，必须是address size（FDE引用的CIE中有定义）的整数倍，即按address size对齐；
2. CIE_pointer (4或8字节），常量，该FDE引用的CIE在.debug_frame的偏移量；
3. initial_location (段选择器，以及目标地址），该table entry对应第一个指令地址，如果segment_size（引用的CIE中定义）非0, initial_location前还需要加一个段选择器；
4. address_range (target address)，该FDE描述的程序指令占用的字节数量；
5. instructions (array of ubyte)，FDE中包含的指令序列，在后面进行描述；
6. padding (array of ubyte)，字节填充，通过DW_CFA_nop指令填充结构体，使FDE结构体大小满足length字段要求；

#### 5.4.3.4 调用帧指令（Call Frame Instructions）

调用帧指令部分，每个指令采用0个或多个操作数，一些操作数可能被编码为操作码的一部分（请参见DWARF v4 section 7.23）。这些说明在以下各节中定义。
有的调用栈帧指令，其操作数通过DWARF表达式编码（请参见DWARF v4 section 2.5.1）。以下DWARF运算符不能在此类操作数（DWARF表达式）中使用：

- DW_OP_call2，DW_OP_call4和DW_OP_call_ref运算符在这些指令的操作数中没有意义，因为没有从调用帧信息到任何相应的调试编译单元信息的映射，因此无法解释调用偏移。
- DW_OP_push_object_address在这些指令的操作数中没有意义，因为没有对象上下文可提供的要push的值.
- DW_OP_call_frame_cfa在这些指令的操作数中没有意义，因为它的使用是循环的。

上述限制适用的调用帧指令包括DW_CFA_def_cfa_expression，DW_CFA_expression和DW_CFA_val_expression。

##### 5.4.3.4.1 CFI表行创建指令（Row Creation Instructions）

1. DW_CFA_set_loc
   
   DW_CFA_set_loc指令采用代表目标地址的单个操作数。 所需的操作是使用指定的地址作为新位置来创建新的表行。新行中的所有其他值最初都与当前行相同。 新位置值始终大于当前位置值。 如果此FDE的CIE的segment_size字段不为零，还需要在在初始位置之前加上段选择器。
   
2. DW_CFA_advance_loc

   DW_CFA_advance指令采用单个操作数（在操作码中编码），该操作数表示常数增量。 所需的操作是使用位置值创建一个新表行，该位置值是通过获取当前条目的位置值并加上delta * code_alignment_factor的值来计算的。 新行中的所有其他值最初都与当前行相同。

3. DW_CFA_advance_loc1
   DW_CFA_advance_loc1指令采用一个表示常量增量的单个ubyte操作数。 除了增量操作数的编码和大小外，该指令与DW_CFA_advance_loc相同。

4. DW_CFA_advance_loc2
   DW_CFA_advance_loc2指令采用单个uhalf操作数表示常数增量。 除了增量操作数的编码和大小外，该指令与DW_CFA_advance_loc相同。

5. DW_CFA_advance_loc4
   DW_CFA_advance_loc4指令采用单个uword操作数来表示恒定增量。 除了增量操作数的编码和大小外，该指令与DW_CFA_advance_loc相同。

##### 5.4.3.4.2 CFI表CFA定义指令（CFA Definition Instructions）

1. DW_CFA_def_cfa
   The DW_CFA_def_cfa instruction takes two unsigned LEB128 operands representing a register number and a (non-factored) offset. The required action is to define the current CFA rule to use the provided register and offset.
2. DW_CFA_def_cfa_sf
   The DW_CFA_def_cfa_sf instruction takes two operands: an unsigned LEB128 value representing a register number and a signed LEB128 factored offset. This instruction is identical to DW_CFA_def_cfa except that the second operand is signed and factored. The resulting offset is factored_offset * data_alignment_factor.
3. DW_CFA_def_cfa_register
   The DW_CFA_def_cfa_register instruction takes a single unsigned LEB128 operand representing a register number. The required action is to define the current CFA rule to use the provided register (but to keep the old offset). This operation is valid only if the current CFA rule is defined to use a register and offset.
4. DW_CFA_def_cfa_offset
   The DW_CFA_def_cfa_offset instruction takes a single unsigned LEB128 operand representing a (non-factored) offset. The required action is to define the current CFA rule to use the provided offset (but to keep the old register). This operation is valid only if the current CFA rule is defined to use a register and offset.
5. DW_CFA_def_cfa_offset_sf
   The DW_CFA_def_cfa_offset_sf instruction takes a signed LEB128 operand representing a factored offset. This instruction is identical to DW_CFA_def_cfa_offset except that the operand is signed and factored. The resulting offset is factored_offset * data_alignment_factor. This operation is valid only if the current CFA rule is defined to use a register and offset.
6. DW_CFA_def_cfa_expression
   The DW_CFA_def_cfa_expression instruction takes a single operand encoded as a DW_FORM_exprloc value representing a DWARF expression. The required action is to establish that expression as the means by which the current CFA is computed.
   See Section 6.4.2 regarding restrictions on the DWARF expression operators that can be used.

##### 5.4.3.4.3 CFI表寄存器规则指令（Register Rule Instructions）

1. DW_CFA_undefined
   The DW_CFA_undefined instruction takes a single unsigned LEB128 operand that represents a register number. The required action is to set the rule for the specified register to “undefined.”
2. DW_CFA_same_value
   The DW_CFA_same_value instruction takes a single unsigned LEB128 operand that represents a register number. The required action is to set the rule for the specified register to “same value.”
3. DW_CFA_offset
   The DW_CFA_offset instruction takes two operands: a register number (encoded with the opcode) and an unsigned LEB128 constant representing a factored offset. The required action is to change the rule for the register indicated by the register number to be an offset(N) rule where the value of N is factored offset * data_alignment_factor.
4. DW_CFA_offset_extended
   The DW_CFA_offset_extended instruction takes two unsigned LEB128 operands representing a register number and a factored offset. This instruction is identical to DW_CFA_offset except for the encoding and size of the register operand.
5. DW_CFA_offset_extended_sf
   The DW_CFA_offset_extended_sf instruction takes two operands: an unsigned LEB128 value representing a register number and a signed LEB128 factored offset. This instruction is identical to DW_CFA_offset_extended except that the second operand is signed and factored. The resulting offset is factored_offset * data_alignment_factor.
6. DW_CFA_val_offset
   The DW_CFA_val_offset instruction takes two unsigned LEB128 operands representing a register number and a factored offset. The required action is to change the rule for the register indicated by the register number to be a val_offset(N) rule where the value of N is factored_offset * data_alignment_factor.
7. DW_CFA_val_offset_sf
   The DW_CFA_val_offset_sf instruction takes two operands: an unsigned LEB128 value representing a register number and a signed LEB128 factored offset. This instruction is identical to DW_CFA_val_offset except that the second operand is signed and factored. The resulting offset is factored_offset * data_alignment_factor.
8. DW_CFA_register
   The DW_CFA_register instruction takes two unsigned LEB128 operands representing register numbers. The required action is to set the rule for the first register to be register(R) where R is the second register.
9. DW_CFA_expression
   The DW_CFA_expression instruction takes two operands: an unsigned LEB128 value representing a register number, and a DW_FORM_block value representing a DWARF expression. The required action is to change the rule for the register indicated by the register number to be an expression(E) rule where E is the DWARF expression. That is, the DWARF expression computes the address. The value of the CFA is pushed on the DWARF evaluation stack prior to execution of the DWARF expression.
   See Section 6.4.2 regarding restrictions on the DWARF expression operators that can be used.
10. DW_CFA_val_expression
    The DW_CFA_val_expression instruction takes two operands: an unsigned LEB128 value representing a register number, and a DW_FORM_block value representing a DWARF expression. The required action is to change the rule for the register indicated by the register number to be a val_expression(E) rule where E is the DWARF expression. That is, the DWARF expression computes the value of the given register. The value of the CFA is pushed on the DWARF evaluation stack prior to execution of the DWARF expression.
    See Section 6.4.2 regarding restrictions on the DWARF expression operators that can be used.
11. DW_CFA_restore
    The DW_CFA_restore instruction takes a single operand (encoded with the opcode) that represents a register number. The required action is to change the rule for the indicated register to the rule assigned it by the initial_instructions in the CIE.
12. DW_CFA_restore_extended
    The DW_CFA_restore_extended instruction takes a single unsigned LEB128 operand that represents a register number. This instruction is identical to DW_CFA_restore except for the encoding and size of the register operand.

##### 5.4.3.4.4 CFI表行状态指令（Row State Instructions）

The next two instructions provide the ability to stack and retrieve complete register states. They may be useful, for example, for a compiler that moves epilogue code into the body of a function.

1. DW_CFA_remember_state
   The DW_CFA_remember_state instruction takes no operands. The required action is to push the set of rules for every register onto an implicit stack.
2. DW_CFA_restore_state
   The DW_CFA_restore_state instruction takes no operands. The required action is to pop the set of rules off the implicit stack and place them in the current row.

##### 5.4.3.4.5 CFI表字节填充指令（Padding Instruction）

1. DW_CFA_nop
   The DW_CFA_nop instruction has no operands and no required actions. It is used as padding to make a CIE or FDE an appropriate size.

#### 5.4.3.5 Call Frame Instruction Usage

#### 5.4.3.6 Example



/////////////////////////////////////

每个处理器都有自己的方式来决定“**如何传递参数和返回值**”，这由处理器的“**ABI（应用程序二进制接口）**”定义。

**在最简单的情况下，所有函数都采用相同的方式来传递参数和返回值**，并且调试器确切地知道如何获取参数和返回值。

**但是实际上，并非每个处理器都使用相同的方式来传递参数和返回值**，不同型号的处理器ABI不同，这个好理解。

**此外，编译器可能会进行一些优化，以使生成的指令更精炼、更快**。 例如，我们创建一个简单的函数，它使用调用者的局部变量作为参数，编译器可能会将其优化（如尾递归优化），而不是将参数传递给被调用函数，这样就可以避免创建新的栈帧，当然也可能会优化寄存器的使用，也许还有其他的优化……这里一点小优化，可能会导致调试器不能准确定位调用函数的栈帧。

> 关于编译器对尾递归的优化，可以参考博文： [tail recursion call optimization]( http://www.ruanyifeng.com/blog/2015/04/tail-call.html)，当前go编译器还不支持尾递归优化，gcc是支持的。

DWARF中的调用栈信息（Call Frame Information，简称CFI）为调试器提供了如下信息，函数是如何被调用的，如何找到函数参数，如何找到调用函数（caller）的调用帧信息。 调试器借助CFI可以展开调用栈、查找上一个函数、确定当前函数的被调用位置以及传递的参数值。

与行号表一样，CFI也被编码为一系列字节码指令，这些指令由CFI状态机解释、执行，以创建完整的CFI表。 每个包含代码的地址都有一行。 第一列包含机器地址，而其他列则包含执行该地址处的指令时（fixme 执行前还是执行后）的寄存器值。 像行号表一样，完整的CFI也非常庞大。 幸运的是，两条指令之间的变化很小，因此CFI编码非常紧凑。

////////////////////

调试器通常需要能够查看和修改调用堆栈上任何子例程激活的状态。 激活包括：