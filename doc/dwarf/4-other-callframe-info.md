### 5.4.3 Call Frame Information 

#### 5.4.3.1 Introduction

Debuggers often need to be able to **view and modify the state of any subroutine activation** that is on the call stack. An activation consists of:

- A code location that is within the subroutine. This location is either the place where the program stopped when the debugger got control (e.g. a breakpoint), or is a place where a subroutine made a call or was interrupted by an asynchronous event (e.g. a signal).
- An area of memory that is allocated on a stack called a “call frame.” The call frame is identified by an address on the stack. We refer to this address as the **Canonical Frame Address or CFA**. Typically, the CFA is defined to be the value of the stack pointer at the call site in the previous frame (which may be different from its value on entry to the current frame).
- A set of registers that are in use by the subroutine at the code location.

Typically, a set of registers are designated to be preserved across a call. If a callee wishes to use such a register, it saves the value that the register had at entry time in its call frame and restores it on exit. The code that allocates space on the call frame stack and performs the save operation is called the subroutine’s prologue, and the code that performs the restore operation and deallocates the frame is called its epilogue. Typically, the prologue code is physically at the beginning of a subroutine and the epilogue code is at the end.

To be able to view or modify an activation that is not on the top of the call frame stack, the debugger must **“virtually unwind” the stack of activations** until it finds the activation of interest. A debugger unwinds a stack in steps. Starting with the current activation it virtually restores any registers that were preserved by the current activation and computes the predecessor’s CFA and code location. This has the logical effect of returning from the current subroutine to its predecessor. We say that the debugger virtually unwinds the stack because the actual state of the target process is unchanged.

#### 5.4.3.2 Arch-Independent Way of Encoding

The unwinding operation needs to know where registers are saved and how to compute the predecessor’s CFA and code location. When considering an architecture-independent way of encoding this information one has to consider a number of special things.

- Prologue and epilogue code is not always in distinct blocks at the beginning and end of a subroutine. It is common to duplicate the epilogue code at the site of each return from the code. Sometimes a compiler breaks up the register save/unsave operations and moves them into the body of the subroutine to just where they are needed.
- Compilers use different ways to manage the call frame. Sometimes they use a frame pointer register, sometimes not.
- The algorithm to compute CFA changes as you progress through the prologue and epilogue code. (By definition, the CFA value does not change.)
- Some subroutines have no call frame.
- Sometimes a register is saved in another register that by convention does not need to be
  saved.
- Some architectures have special instructions that perform some or all of the register management in one instruction, leaving special information on the stack that indicates how registers are saved.
- Some architectures treat return address values specially. For example, in one architecture, the call instruction guarantees that the low order two bits will be zero and the return instruction ignores those bits. This leaves two bits of storage that are available to other uses that must be treated specially.

#### 5.4.3.3 Structure of Call Frame Information

DWARF supports virtual unwinding by defining an architecture independent basis for recording how procedures save and restore registers during their lifetimes. This basis must be augmented on some machines with specific information that is defined by an architecture specific ABI authoring committee, a hardware vendor, or a compiler producer. The body defining a specific augmentation is referred to below as the “augmenter.”

Abstractly, this mechanism describes a very large table that has the following structure:

<img src="assets/image-20191229130341692.png" alt="image-20191229130341692" style="zoom:5%;" />

The first column indicates an address for every location that contains code in a program. (In shared objects, this is an object-relative offset.) The remaining columns contain virtual unwinding rules that are associated with the indicated location.

The CFA column defines the rule which computes the Canonical Frame Address value; it may be either a register and a signed offset that are added together, or a DWARF expression that is evaluated.

The remaining columns are labeled by register number. This includes some registers that have special designation on some architectures such as the PC and the stack pointer register. (The actual mapping of registers for a particular architecture is defined by the augmenter.) The register columns contain rules that describe whether a given register has been saved and the rule to find the value for the register in the previous frame.

The register rules are:

- undefined, A register that has this rule has no recoverable value in the previous frame. (By convention, it is not preserved by a callee.)
- same value, This register has not been modified from the previous frame. (By convention, it is preserved by the callee, but the callee has not modified it.)
- offset(N), The previous value of this register is saved at the address CFA+N where CFA is the current CFA value and N is a signed offset.
- val_offset(N), The previous value of this register is the value CFA+N where CFA is the current CFA value and N is a signed offset.
- register(R), The previous value of this register is stored in another register numbered R.
- expression(E), The previous value of this register is located at the address produced by
  executing the DWARF expression E.
- val_expression(E), The previous value of this register is the value produced by executing the
  DWARF expression E.
- architectural, The rule is defined externally to this specification by the augmenter.

This table would be extremely large if actually constructed as described. Most of the entries at any point in the table are identical to the ones above them. The whole table can be represented quite compactly by recording just the differences starting at the beginning address of each subroutine in the program.

The virtual unwind information is encoded in a self-contained section called .debug_frame. Entries in a .debug_frame section are aligned on a multiple of the address size relative to the start of the section and come in two forms: a Common Information Entry (CIE) and a Frame Description Entry (FDE).

////////////////////////

Every processor has a certain way of determining how to pass parameters and return values, this is defined by the processor’s ABI (Application Binary Interface). In the simplest case, all functions have the same way to pass parameters and return values, and the debuggers know exactly how to get the parameters and return values. 

But actually, not every processor uses the same way to pass parameters and to return values. Besides, compilers may do some optimization to make generated instructions much smaller and faster. For example, maybe a simple function is created to use caller’s local variables as parameters instead of passing parameters to callee to avoid create frame, maybe to optimize the usage of registers, maybe others… The result is a small change takes place in the optimizations and the debugger may no longer be able to walk the stack to the calling functions. 

> for example, [tail recursion call optimization]( http://www.ruanyifeng.com/blog/2015/04/tail-call.html)

Dwarf call frame information (CFI) provides the debugger enough information about how a function is called, how to locate the parameters to the functions, how to locate the call frame for the calling function. This information is used by the debugger to unwind the stack, locating the previous function, the location where the function is called, and the values passed. 

Like the line number table, CFI is also encoded as a sequence of instructions that are interpreted to generate a table. There’s one row for each address that contains code. The first column contains the machine address, while others contain the registers’ values at when instruction at that address is executed. Like the line number table, the complete CFI is huge. Luckily, there’s very little change between two instructions, so the CFI encoding is quite compact. 

/////////////////





