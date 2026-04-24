package ts_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sugarme/gotch"
	"github.com/sugarme/gotch/ts"
)

// printChunk prints a 2D tensor chunk followed by a "shape=[r c]" line,
// using Go stdout so that Example tests can capture the output. libtorch's
// native tensor.Print() writes through cgo directly to C++ std::cout,
// bypassing Go's os.Stdout capture — which is why the original Example
// tests silently mis-reported as failing.
func printChunk(t *ts.Tensor) {
	sz, err := t.Size()
	if err != nil {
		fmt.Printf("size err: %v\n", err)
		return
	}
	vals := t.Float64Values(false)
	if len(sz) != 2 {
		// Examples use 2D chunks only; bail out cleanly if this changes.
		fmt.Printf("expected 2D, got shape=%v values=%v\n", sz, vals)
		return
	}
	cols := int(sz[1])
	var b strings.Builder
	for i, v := range vals {
		if i > 0 && i%cols == 0 {
			b.WriteByte('\n')
		} else if i > 0 {
			b.WriteString("  ")
		}
		fmt.Fprintf(&b, "%g", v)
	}
	b.WriteByte('\n')
	fmt.Fprintf(&b, "shape=%v\n", sz)
	fmt.Print(b.String())
}

// ExampleTensor_Split demonstrates splitting a 5×2 tensor along dim 0 into
// fixed-size chunks of 2 rows. 5 rows ÷ 2 = three chunks sized [2, 2, 1].
func ExampleTensor_Split() {
	tensor := ts.MustArange(ts.FloatScalar(10), gotch.Float, gotch.CPU).MustView([]int64{5, 2}, true)
	splitTensors := tensor.MustSplit(2, 0, false)

	for _, st := range splitTensors {
		printChunk(st)
	}

	// Output:
	// 0  1
	// 2  3
	// shape=[2 2]
	// 4  5
	// 6  7
	// shape=[2 2]
	// 8  9
	// shape=[1 2]
}

// ExampleTensor_SplitWithSizes demonstrates splitting a 5×2 tensor along
// dim 0 into explicitly-sized chunks [1, 4].
func ExampleTensor_SplitWithSizes() {
	tensor := ts.MustArange(ts.FloatScalar(10), gotch.Float, gotch.CPU).MustView([]int64{5, 2}, true)
	splitTensors := tensor.MustSplitWithSizes([]int64{1, 4}, 0, false)

	for _, st := range splitTensors {
		printChunk(st)
	}

	// Output:
	// 0  1
	// shape=[1 2]
	// 2  3
	// 4  5
	// 6  7
	// 8  9
	// shape=[4 2]
}

// Test Unbind op specific for BFloat16/Half
func TestTensorUnbind(t *testing.T) {
	// device := gotch.CudaIfAvailable()
	device := gotch.CPU

	dtype := gotch.BFloat16
	// dtype := gotch.Half // <- NOTE. Libtorch API Error: "arange_cpu" not implemented for 'Half'

	x := ts.MustArange(ts.IntScalar(60), dtype, device).MustView([]int64{3, 4, 5}, true)

	out := x.MustUnbind(0, true)

	if len(out) != 3 {
		t.Errorf("Want 3, got %v\n", len(out))
	}
}
