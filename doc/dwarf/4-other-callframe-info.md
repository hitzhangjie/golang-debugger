### 5.4.3 Call Frame Information 

#### 5.4.3.1 Introduction

Debuggers often need to be able to view and modify the state of any subroutine activation that is on the call stack. An activation consists of:

- A code location that is within the subroutine. This location is either the place where the program stopped when the debugger got control (e.g. a breakpoint), or is a place where a subroutine made a call or was interrupted by an asynchronous event (e.g. a signal).
- An area of memory that is allocated on a stack called a “call frame.” The call frame is identified by an address on the stack. We refer to this address as the Canonical Frame Address or CFA. Typically, the CFA is defined to be the value of the stack pointer at the call site in the previous frame (which may be different from its value on entry to the current frame).
- A set of registers that are in use by the subroutine at the code location.

Typically, a set of registers are designated to be preserved across a call. If a callee wishes to use such a register, it saves the value that the register had at entry time in its call frame and restores it on exit. The code that allocates space on the call frame stack and performs the save operation is called the subroutine’s prologue, and the code that performs the restore operation and deallocates the frame is called its epilogue. Typically, the prologue code is physically at the beginning of a subroutine and the epilogue code is at the end.



////////////////////////

Every processor has a certain way of determining how to pass parameters and return values, this is defined by the processor’s ABI (Application Binary Interface). In the simplest case, all functions have the same way to pass parameters and return values, and the debuggers know exactly how to get the parameters and return values. 

But actually, not every processor uses the same way to pass parameters and to return values. Besides, compilers may do some optimization to make generated instructions much smaller and faster. For example, maybe a simple function is created to use caller’s local variables as parameters instead of passing parameters to callee to avoid create frame, maybe to optimize the usage of registers, maybe others… The result is a small change takes place in the optimizations and the debugger may no longer be able to walk the stack to the calling functions. 

> for example, [tail recursion call optimization]( http://www.ruanyifeng.com/blog/2015/04/tail-call.html)

Dwarf call frame information (CFI) provides the debugger enough information about how a function is called, how to locate the parameters to the functions, how to locate the call frame for the calling function. This information is used by the debugger to unwind the stack, locating the previous function, the location where the function is called, and the values passed. 

Like the line number table, CFI is also encoded as a sequence of instructions that are interpreted to generate a table. There’s one row for each address that contains code. The first column contains the machine address, while others contain the registers’ values at when instruction at that address is executed. Like the line number table, the complete CFI is huge. Luckily, there’s very little change between two instructions, so the CFI encoding is quite compact. 

/////////////////





