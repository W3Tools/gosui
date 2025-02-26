package transactions

import (
	"fmt"
	"sync"

	"github.com/W3Tools/gosui/types"
	"github.com/W3Tools/gosui/utils"
)

var (
	cache = &Cache{MoveFunction: make(map[string]*MoveFunctionCacheEntry)}
)

type Cache struct {
	mutex        sync.RWMutex
	MoveFunction map[string]*MoveFunctionCacheEntry
}

type MoveFunctionCacheEntry struct {
	Package    string
	Module     string
	Function   string
	Normalized *types.SuiMoveNormalizedFunction
}

// Cache: Move Function
func (c *Cache) GetMoveFunctionDefinition(pkg, mod, fn string) *MoveFunctionCacheEntry {
	name := fmt.Sprintf("%s::%s::%s", utils.NormalizeSuiAddress(pkg), mod, fn)

	return c.MoveFunction[name]
}

func (c *Cache) AddMoveFunctionDefinition(entry *MoveFunctionCacheEntry) {
	name := fmt.Sprintf("%s::%s::%s", utils.NormalizeSuiAddress(entry.Package), entry.Module, entry.Function)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.MoveFunction[name] = entry
}

func (c *Cache) DeleteMoveFunctionDefinition(pkg, mod, fn string) {
	name := fmt.Sprintf("%s::%s::%s", utils.NormalizeSuiAddress(pkg), mod, fn)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.MoveFunction, name)
}
