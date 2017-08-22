package chandict

import "testing"
import (
	"github.com/stretchr/testify/assert"
	"context"
)

func TestConcurrentDict(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d := NewChanDict(ctx)
	d.SetVal("foo", "bar")
	val, ok := d.ReadVal("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", val)
	assert.True(t, d.CasVal("foo", "bar", "baz", false))
	assert.False(t, d.CasVal("foo", "bar", "baz", false))

}
