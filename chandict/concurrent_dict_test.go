package chandict

import "testing"
import (
	"context"
	"github.com/sargun/goconcurrency/types"
)

func TestConcurrentDict(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d := NewChanDict(ctx)
	types.TestImplementation(t, d)
}
