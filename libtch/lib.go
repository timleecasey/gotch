package libtch

// #cgo CFLAGS: -I${SRCDIR} -O3 -Wall -Wno-unused-variable -Wno-deprecated-declarations -Wno-c++11-narrowing -g -Wno-sign-compare -Wno-unused-function
// #cgo CFLAGS: -I/Users/tlc/src/gotch/libtorch-2.10.0-macos/include/torch/csrc/api/include
// #cgo CFLAGS: -I/Users/tlc/src/gotch/libtorch-2.10.0-macos/include
// #cgo LDFLAGS: -L/Users/tlc/src/gotch/libtorch-2.10.0-macos/lib -ltorch -ltorch_cpu -lc10
// #cgo LDFLAGS: -Wl,-rpath,/Users/tlc/src/gotch/libtorch-2.10.0-macos/lib
// #cgo CXXFLAGS: -std=c++17
import "C"
