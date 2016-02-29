# Gomutate [![Build Status](https://travis-ci.org/zabawaba99/gomutate.svg?branch=master)](https://travis-ci.org/zabawaba99/gomutate) [![Coverage Status](https://coveralls.io/repos/github/zabawaba99/gomutate/badge.svg?branch=master)](https://coveralls.io/github/zabawaba99/gomutate?branch=master)

A mutation testing tool for Go programs.

Gomutate was inspired by the [Pitest](http://pitest.org/) tool
used for mutation testing on the JVM.

## Usage

Install the binary:

```bash
go get github.com/zabawaba99/gomutate/cmd/gomutate
```

Navigate to the directory of your project and run:

```bash
gomutate ./...
```

The above command will run the default mutations on all packages in your
project.

### Flags

```
Usage:
  gomutate [OPTIONS]

Application Options:
  -d, --debug    Show debug information
  -m, --mutator= The mutators to apply (conditionals-boundary)

Help Options:
  -h, --help     Show this help message
```

## Available mutators

* [conditionals-boundary](http://pitest.org/quickstart/mutators/#CONDITIONALS_BOUNDARY)
* [negate-conditionals](http://pitest.org/quickstart/mutators/#NEGATE_CONDITIONALS)

## Mutation Tests

A chunk of the documentation below was taken from [Pitest](http://pitest.org/)
since they did a great job of explaining what mutation testing is and why
you should do it.

### What is mutation testing?

Mutation testing is conceptually quite simple.

Faults (or **mutations**) are automatically seeded into your code, then your tests are run.
If your tests fail then the mutation is **killed**, if your tests pass then the
mutation **lived**.

The quality of your tests can be gauged from the percentage of mutations killed.

### What?

To put it another way - Gomutate runs your unit tests against automatically modified
versions of your application code. When the application code changes, it should produce
different results and cause the unit tests to fail. If a unit test does not fail in this
situation, it may indicate an issue with the test suite.

### Why?

Traditional test coverage (i.e line, statement, branch etc) measures only which code is
**executed** by your tests. It does **not** check that your tests are actually able to
**detect faults** in the executed code. It is therefore only able to identify code
the is definitely **not tested**.

The most extreme example of the problem are tests with no assertions. Fortunately these
are uncommon in most code bases. Much more common is code that is only partially tested
by its suite. A suite that only **partially tests** code can still execute all its
branches.

As it is actually able to detect whether each statement is meaningfully tested, mutation
testing is the **gold standard** against which all other types of coverage are measured.
