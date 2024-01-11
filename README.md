# The Whirl Programming Language

[![Test Status](https://github.com/whirl-lang/whirl/workflows/test/badge.svg)](https://github.com/whirl-lang/whirl/actions)

This is the main source code repository for Whirl. It contains the compiler, standard library, and documentation.

## Quick links

- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Dependencies](#dependencies)
- [Syntax](#syntax)
  - [Comments](#comments)
  - [Variables](#variables)
  - [Control flow](#control-flow)
  - [Functions](#functions)
- [License](#license)

## Installation

Pre-built binaries can be found in the Releases section on GitHub. Otherwise, Whirl can be built from source.

```bash
git clone https://github.com/whirl-lang/whirl
cd whirl
go build cmd/whirl/main.go
```

## Usage

```bash
A statically typed, compiled programming language.

Usage: whirl <PATH>

Arguments:
  <PATH>  Path to the script to execute

Options:
  -h, --help     Print help
  -V, --version  Print version
```

## Examples

Examples can be found in the [examples](examples) directory.

## Dependencies

Whirl requires the [Tiny C](https://bellard.org/tcc/) compiler to be downloaded on the system and added to the PATH.

## Syntax

### Comments

```r
 $ This is a comment
```

### Variables

```rust
let x = 1;
let y = 2;
let z = x + y;
```

### Control flow

```c
if x == 1 {
  printf("x is 1");
} else if x == 2 {
  printf("x is 2")
} else {
  printf("x is neither 1 or 2")
}
```

### Functions

```rust
proc main() :: int {
  escape 0;
}
```

## License

Whirl is distributed under the MIT license. See [LICENSE](LICENSE) for more information.
