package types

type ConcurrentDict interface {
	SetVal(key, val string)
	CasVal(key, oldVal, newVal string, setOnNotExists bool) bool
	ReadVal(key string) (string, bool)
}
