# sctrl
Serial protocols controller

[![Go Reference](https://pkg.go.dev/badge/github.com/fdelbos/sctrl.svg)](https://pkg.go.dev/github.com/fdelbos/sctrl)
[![CI Tests](https://github.com/fdelbos/sctrl/actions/workflows/ci.yaml/badge.svg?branch=master)](https://github.com/fdelbos/sctrl/actions?branch=master)

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
