package types

import "testing"
import "github.com/stretchr/testify/assert"


func TestImplementation(t *testing.T, d ConcurrentDict) {
	d.SetKey("foo", "bar")
	val, ok := d.ReadKey("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", val)
	assert.True(t, d.CasVal("foo", "bar", "baz", false))
	assert.False(t, d.CasVal("foo", "bar", "baz", false))
	d.DeleteKey("foo")
	_, ok = d.ReadKey("foo")
	assert.False(t, ok)
}