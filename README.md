# minijs

**minijs** is a JavaScript bytecode virtual machine implemented in Go. It compiles JavaScript code into bytecode and executes it in a virtual machine to optimize performance. With **high compatibility with Go**, it can be seamlessly embedded into a wide range of Go-based applications.

## Key Features

- **Bytecode Execution**: Compiles JavaScript code into bytecode for execution in the virtual machine, optimizing performance.
- **High Compatibility with Go**: Implemented in Go, it offers high compatibility with Go environments and can be easily embedded into various Go-based applications.

## Installation

**minijs** runs in a Go environment. Use the following commands to install and build the project locally:

```bash
git clone https://github.com/siyul-park/minijs.git
cd minijs
make build
```

## Usage

You can run JavaScript code interactively using the REPL (Read-Eval-Print Loop).

### Running REPL

```bash
minijs
```

### Example

```bash
> 'b'+'a'+ +'a'+'a'
"baNaNa"
```
