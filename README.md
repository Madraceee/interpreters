# Interpreters

Built 2 interpreter for the Language [Lox](https://craftinginterpreters.com/the-lox-language.html) based on the book "Crafting Interpreter".

(Glox)Golang Interpreter is a AST interpreter. Given code, an AST(Abstract Syntax Tree) is built and executed.
(Clox)C Interpreter is a byte code interpreter with its own VM. Given code, it is converted into bytecode and then executed.

Both are under construction.
This project is used for learning purpose.

If u want to learn more about interpreters, check out my [blog](https://blog.iamnitheesh.com). I have my learnings stored there.

#### Projects based on my learnings
- [Markdown to Blog(HTML)](https://github.com/Madraceee/md-to-html)
  - Built a transpiler which converts Markdown to AST and then to HTML


## How to Run

To build (Requires Golang, make)
```
cd glox
make build
```
To run 
```
./bin/glox FILENAME
```
