package lockdict

import "testing"
import "github.com/sargun/goconcurrency/types"

func TestConcurrentDict(t *testing.T) {
	d := NewLockDict()
	types.TestImplementation(t, d)

}
