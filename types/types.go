package types

type ConcurrentDict interface {
	// SetKey sets the key to val, val no matter if it exists or not
	SetKey(key, val string)
	// CasVal performs an atomic-compare-and-swap operation for a given key. setOnNotExists overwrites the "nil" value
	CasVal(key, oldVal, newVal string, setOnNotExists bool) bool
	// ReadKey returns the given value associated with a key, and if it exists in the map or not.
	ReadKey(key string) (string, bool)
	// DeleteKey removes a key from the map. If the key does not exist, we perform the operation without error
	DeleteKey(key string)
}
