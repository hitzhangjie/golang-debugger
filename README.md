# How to develop a (golang) debugger

## Introduction

This project aims to introduce how to develop a (golang) debugger, including Operating System's support, how to coordinate work between compiler, linker and debugger, debugging information standard, mapping between machine instruction and source code, etc. 

Thanks to [delve](github.com/go-delve/delve) and the author [derek parker](https://twitter.com/derkthedaring?lang=en) and other contributors. I learned a lot from them. I want to share the knowledge to develop a (golang) debugger. I hope this project can be useful for developers interested in debugging topic.

To develop a symbolic debugger need to combine the knowledge of CPU instruction (like instruction patching), Operating System (like linux ptrace and OS scheduler), compilers, linkers, loaders, debuggers (how to coordinate the work between them), executable file format (how to store debugging information), debugging information format (how to describe source code, how to map between instruction and source, vice versa), and features of different programming languages (like goroutine concept), so I think it's also a good chance to improve the understanding of computer technology.

I think it's very helpful, So I am really excited to write this documents.

## Contact

Please email me **hit.zhangjie@gmail.com**, I will respond as soon as possible.

