# DBG

### What is DBG?

DBG is a Go debugger, written primarily in Go. It really needs a new name.

### Building

Currently, DBG requires the following [patch](https://codereview.appspot.com/117280043/) to be applied to your Go source to build.

### Features

* Attach to (trace) a running process
* Set breakpoints
* Single step through a process
* Next through a process (step over / out of subroutines)
* Never retype commands, empty line defaults to previous command

### Usage

* `break` - Set break point at the entry point of a function, or at a specific file/line. Example: `break foo.go:13`.

* `step` - Single step through program.

* `next` - Step over to next source line.

### Upcoming features

* Handle Gos multithreaded nature better (follow goroutine accross thread contexts)
* In-scope variable evaluation
* In-scope variable setting
* Readline integration
* Ability to launch debugging session from debugged program, with breakpoint set correctly
* Support for OS X

### License

MIT
