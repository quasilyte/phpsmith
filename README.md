# phpsmith

[![Tests](https://github.com/quasilyte/phpsmith/workflows/Tests/badge.svg)](https://github.com/quasilyte/phpsmith/blob/master/.github/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/quasilyte/phpsmith)](https://goreportcard.com/report/github.com/quasilyte/phpsmith)
[![Go Reference](https://pkg.go.dev/badge/github.com/quasilyte/phpsmith.svg)](https://pkg.go.dev/github.com/quasilyte/phpsmith)

`phpsmith` creates random PHP and KPHP programs to test their compilers and runtimes.

It can create both valid (but very contrived) and invalid programs.

### Features

* Random valid php code generation
* Random invalid php code generation
* Reproducible program generation by seed
* Big generation variability

### How it works

Phpsmith can be executed in two modes: `fuzz`, `generate`:

- `fuzz`:
    - infinitely generate php programs
    - run it on php and kphp
    - catch exceptions, segmentation faults, fatal errors
    - compare results between php and kphp
    - save diff in logs

- `generate`:
    - generate php program by provided seed

### Installation

```bash
go install github.com/quasilyte/phpsmith
```

### Pre-installation requirements

- php ≥7.0
- kphp ≥9.3.0

### Usage

Run `phpsmith` without arguments to get help output.

```
Possible commands are:

  version     print phpsmith version info to stdout and exit
  fuzz        run fuzzing using the provided configuration
  generate    generate a program using the provided configuration
```

`fuzz` command examples:

```bash
phpsmith fuzz -o ~/phpsmith_out
```

`generate` command examples:

```bash
phpsmith generate -seed 1651182107
```
