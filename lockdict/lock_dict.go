package lockdict

import (
	"github.com/sargun/goconcurrency/types"
	"sync"
)

var _ types.ConcurrentDict = (*LockDict)(nil)

func NewLockDict() *LockDict {
	return &LockDict{
		dict: make(map[string]string),
		lock: sync.RWMutex{},
	}
}

type LockDict struct {
	dict map[string]string
	lock sync.RWMutex
}

func (dict *LockDict) SetKey(key, val string) {
	dict.lock.Lock()
	defer dict.lock.Unlock()
	dict.dict[key] = val
}
func (dict *LockDict) CasVal(key, oldVal, newVal string, setOnNotExists bool) bool {
	dict.lock.Lock()
	defer dict.lock.Unlock()
	if val, exists := dict.dict[key]; exists && val == oldVal {
		dict.dict[key] = newVal
		return true
	} else if !exists && setOnNotExists {
		dict.dict[key] = newVal
		return true
	}

	return false
}
func (dict *LockDict) ReadKey(key string) (string, bool) {
	dict.lock.RLock()
	defer dict.lock.RUnlock()
	val, exists := dict.dict[key]
	return val, exists
}

func (dict *LockDict) DeleteKey(key string) {
	dict.lock.Lock()
	defer dict.lock.Unlock()
	delete(dict.dict, key)
}
