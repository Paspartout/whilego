# Todo

## v0.1 - Beta

- [ ] Finish parser
	- [ ] More SeqExpr Tests
	- [ ] Parse WhileExpr
		- [ ] Tests
		- [ ] Implement
	- [ ] Test some error cases
- [ ] Interpreter
	- [ ] Write tests
	- [ ] Interpret program to output x0
	- [ ] Parse program to output x0
- [ ] cmd: whilego [filename]
	- [ ] interpret program from stdin and write x0 to stdout
	- [ ] interpret program from [filename] if provided
- [ ] TravisCi
- [ ] README, LICENSE

## v1.0 - Release

- [ ] Beta Test, Feedback

- [ ] Ergonomics
	- [ ] Macros using [templates](https://golang.org/pkg/text/template/)?
- [ ] Online Interpreter using [GopherJS](https://github.com/gopherjs/gopherjs)
- [ ] Debugging
	- [ ] Step through every expression, printing every variable

## Probably future versions

- [ ] Transpiling
	- [ ] Format Code/Pretty Printing
	- [ ] Transpile to C
	- [ ] Transpile to Go
	- [ ] Transpile to JavaScript?
- [ ] Compiling
	- [ ] Compile to ASM/IR and use the [Go Assembler](https://golang.org/doc/asm)
	- [ ] llvm IR
- [ ] Debugging
	- [ ] Breakpoints?
	- [ ] Watchpoints?

## Random thoughts

- Programs are functions P(x1,...,xk) -> x0

