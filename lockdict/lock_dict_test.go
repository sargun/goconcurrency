package lockdict

import "testing"
import "github.com/stretchr/testify/assert"

func TestConcurrentDict(t *testing.T) {
	d := NewLockDict()
	d.SetVal("foo", "bar")
	val, ok := d.ReadVal("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", val)
	assert.True(t, d.CasVal("foo", "bar", "baz", false))
	assert.False(t, d.CasVal("foo", "bar", "baz", false))

}
