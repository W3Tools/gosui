package transactions

import (
	"fmt"
	"sync"

	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/gosui/types"
	"github.com/W3Tools/gosui/utils"
)

var cache = &Cache{
	MoveFunction:            make(map[string]*MoveFunctionCacheEntry),
	SharedOrImmutableObject: make(map[string]*SharedObjectCacheEntry),
}

type Cache struct {
	mutex                   sync.RWMutex
	MoveFunction            map[string]*MoveFunctionCacheEntry
	SharedOrImmutableObject map[string]*SharedObjectCacheEntry
}

type MoveFunctionCacheEntry struct {
	Package    string
	Module     string
	Function   string
	Normalized *types.SuiMoveNormalizedFunction
}

type SharedObjectCacheEntry struct {
	ObjectId             *sui_types.ObjectID
	InitialSharedVersion *uint64
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

// Cache: Shared Object
func (c *Cache) GetSharedObject(id string) *SharedObjectCacheEntry {
	return c.SharedOrImmutableObject[id]
}

func (c *Cache) GetSharedObjects(ids []string) []*SharedObjectCacheEntry {
	entries := make([]*SharedObjectCacheEntry, 0)

	for _, id := range ids {
		entries = append(entries, c.SharedOrImmutableObject[id])
	}

	return entries
}

func (c *Cache) AddSharedObject(entry *SharedObjectCacheEntry) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.SharedOrImmutableObject[entry.ObjectId.String()] = entry
}

func (c *Cache) AddSharedObjects(entries []*SharedObjectCacheEntry) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, entry := range entries {
		c.SharedOrImmutableObject[entry.ObjectId.String()] = entry
	}
}

func (c *Cache) DeleteSharedObject(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.SharedOrImmutableObject, id)
}

func (c *Cache) DeleteSharedObjects(ids []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, id := range ids {
		delete(c.SharedOrImmutableObject, id)
	}
}

// Convert SharedObjectCacheEntry to sui_types.ObjectArg
func (entry *SharedObjectCacheEntry) ToObjectArg(mutable bool) *sui_types.ObjectArg {
	objectArg := new(sui_types.ObjectArg)

	objectArg.SharedObject = &struct {
		Id                   move_types.AccountAddress
		InitialSharedVersion uint64
		Mutable              bool
	}{
		Id:                   *entry.ObjectId,
		InitialSharedVersion: *entry.InitialSharedVersion,
		Mutable:              mutable,
	}

	return objectArg
}
