package chandict

import "testing"
import (
	"github.com/sargun/goconcurrency/types"
	"runtime"
)

func TestConcurrentDict(t *testing.T) {
	d := NewChanDict()
	types.TestImplementation(t, d)
}

func TestBreak(t *testing.T) {
	innerBitBreak()
	/* Force GC, to require finalizer to run */
	runtime.GC()
}
func innerBitBreak() {
	d := NewChanDict()
	d.ReadKey("foo")
}
