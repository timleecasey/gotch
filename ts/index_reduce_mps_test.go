package ts_test

import (
	"runtime"
	"testing"

	"github.com/sugarme/gotch"
	"github.com/sugarme/gotch/ts"
)

// TestIndexReduceAmaxMPS verifies that index_reduce(reduce="amax") executes on the
// MPS device under libtorch 2.12. The MPS kernel (index_reduce_mps_out) first ships
// in libtorch 2.12; on 2.10/2.11 this op has no MPS kernel and errors unless
// PYTORCH_ENABLE_MPS_FALLBACK routes it to CPU. The test runs with that env var
// UNSET, so a missing kernel surfaces as an error rather than a silent CPU fallback.
func TestIndexReduceAmaxMPS(t *testing.T) {
	if !gotch.MPS.IsAvailable() {
		t.Skip("MPS not available on this host")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	dev := gotch.Device(gotch.MPS)

	// dim=0, include_self=true, reduce=amax.
	//   self=[0,0,0], index=[0,1,1], source=[3,2,5]
	//   self[0] = max(0, 3)    = 3
	//   self[1] = max(0, 2, 5) = 5
	//   self[2] =     0        = 0
	self := ts.MustOfSlice([]float32{0, 0, 0}).MustTo(dev, true)
	index := ts.MustOfSlice([]int64{0, 1, 1}).MustTo(dev, true)
	source := ts.MustOfSlice([]float32{3, 2, 5}).MustTo(dev, true)
	defer self.MustDrop()
	defer index.MustDrop()
	defer source.MustDrop()

	out, err := self.IndexReduce(0, index, source, "amax", true, false)
	if err != nil {
		t.Fatalf("index_reduce(amax) errored on MPS — no MPS kernel (pre-2.12) or fallback disabled: %v", err)
	}
	defer out.MustDrop()

	if got := out.MustDevice().Value; got != gotch.MPS.Value {
		t.Fatalf("output device = %d, want MPS (%d) — op may have fallen back to CPU", got, gotch.MPS.Value)
	}

	cpu := out.MustTo(gotch.CPU, false)
	defer cpu.MustDrop()
	got := cpu.Float64Values()
	want := []float64{3, 5, 0}
	if len(got) != len(want) {
		t.Fatalf("result length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index_reduce(amax) result[%d] = %v, want %v (full=%v)", i, got[i], want[i], got)
		}
	}
}
