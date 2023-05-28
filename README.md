# sctrl
Serial protocols controller

[![Go Reference](https://pkg.go.dev/badge/github.com/fdelbos/sctrl.svg)](https://pkg.go.dev/github.com/fdelbos/sctrl)
![CI Workflow](https://github.com/github/docs/actions/workflows/ci.yml/badge.svg)

A controller to use with all sorts of serial devices where commands are run one by one
in a non concurrent fashion and where sometimes notifications can appear.

## Test
- make sure you have mockery installed
```sh
go install github.com/vektra/mockery/v2@v2.20.0
```
- then run the tests:
```sh
make test
```
