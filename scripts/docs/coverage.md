# How to generate code coverage chart(s) for a specified service or all services

Using this script can help you quickly test the test cases of a specific service or all services and generate their code coverage charts.

**Note:** The execution directory must be located in the terraform-provider-huaweicloud directory or its subdirectories.

## Usage Guide

```bash
$ ./coverage.sh [options] [args] ([options] [args]...)
```

## Operators

**-o**: The name of the service (package) where the test case to be tested is located, such as 'aad'.

**-f**: The name or prefix of the test case to be tested (the prefix can specify multiple tests to run simultaneously).

**-n**: The number of concurrent tests, defaults to `10`.

**-s**: The name of the service (package) where the resource code of the test case to be tested is located.  
    By default, the name is consistent with the service (package) name of the test case.  
    It only needs to be specified when it is inconsistent (in rare cases).

**-h | --help**: Get operator usage tips.

## Example

### Executes all tests for all services and generates code coverage charts (one for each service)

```bash
# concurrency is 10.
$ ./coverage.sh
```

### Executes all tests starting with TestAccForward under the aad package and generate coverage charts

```bash
# concurrency is 2.
$ ./coverage.sh -o aad -f TestAccForward -n 2
```
