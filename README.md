# sctrl
Serial protocols controller

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
