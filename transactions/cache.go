package transactions

import (
	"fmt"
	"sync"

	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/gosui/types"
	"github.com/W3Tools/gosui/utils"
)

// cache defines a global instance of Cache for Move functions and shared or immutable objects.
var cache = &Cache{
	MoveFunction:            make(map[string]*MoveFunctionCacheEntry),
	SharedOrImmutableObject: make(map[string]*SharedObjectCacheEntry),
}

// Cache defines a cache for Move functions and shared or immutable objects.
type Cache struct {
	mutex                   sync.RWMutex
	MoveFunction            map[string]*MoveFunctionCacheEntry
	SharedOrImmutableObject map[string]*SharedObjectCacheEntry
}

// MoveFunctionCacheEntry defines a cache entry for a Move function.
type MoveFunctionCacheEntry struct {
	Package    string
	Module     string
	Function   string
	Normalized *types.SuiMoveNormalizedFunction
}

// SharedObjectCacheEntry defines a cache entry for a shared or immutable object.
type SharedObjectCacheEntry struct {
	ObjectID             *sui_types.ObjectID
	InitialSharedVersion *uint64
}

// GetMoveFunctionDefinition retrieves a cached Move function definition by package, module, and function name.
func (c *Cache) GetMoveFunctionDefinition(pkg, mod, fn string) *MoveFunctionCacheEntry {
	name := fmt.Sprintf("%s::%s::%s", utils.NormalizeSuiAddress(pkg), mod, fn)

	return c.MoveFunction[name]
}

// AddMoveFunctionDefinition adds a Move function definition to the cache.
func (c *Cache) AddMoveFunctionDefinition(entry *MoveFunctionCacheEntry) {
	name := fmt.Sprintf("%s::%s::%s", utils.NormalizeSuiAddress(entry.Package), entry.Module, entry.Function)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.MoveFunction[name] = entry
}

// DeleteMoveFunctionDefinition removes a Move function definition from the cache by package, module, and function name.
func (c *Cache) DeleteMoveFunctionDefinition(pkg, mod, fn string) {
	name := fmt.Sprintf("%s::%s::%s", utils.NormalizeSuiAddress(pkg), mod, fn)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.MoveFunction, name)
}

// GetSharedObject retrieves a shared or immutable object from the cache by its ID.
func (c *Cache) GetSharedObject(id string) *SharedObjectCacheEntry {
	return c.SharedOrImmutableObject[id]
}

// GetSharedObjects retrieves multiple shared or immutable objects from the cache by their IDs.
func (c *Cache) GetSharedObjects(ids []string) []*SharedObjectCacheEntry {
	entries := make([]*SharedObjectCacheEntry, 0)

	for _, id := range ids {
		entries = append(entries, c.SharedOrImmutableObject[id])
	}

	return entries
}

// AddSharedObject adds a shared or immutable object to the cache.
func (c *Cache) AddSharedObject(entry *SharedObjectCacheEntry) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.SharedOrImmutableObject[entry.ObjectID.String()] = entry
}

// AddSharedObjects adds multiple shared or immutable objects to the cache.
func (c *Cache) AddSharedObjects(entries []*SharedObjectCacheEntry) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, entry := range entries {
		c.SharedOrImmutableObject[entry.ObjectID.String()] = entry
	}
}

// DeleteSharedObject removes a shared or immutable object from the cache by its ID.
func (c *Cache) DeleteSharedObject(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.SharedOrImmutableObject, id)
}

// DeleteSharedObjects removes multiple shared or immutable objects from the cache by their IDs.
func (c *Cache) DeleteSharedObjects(ids []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, id := range ids {
		delete(c.SharedOrImmutableObject, id)
	}
}

// ToObjectArg encodes a SharedObjectCacheEntry as a Sui ObjectArg.
func (entry *SharedObjectCacheEntry) ToObjectArg(mutable bool) *sui_types.ObjectArg {
	objectArg := new(sui_types.ObjectArg)

	objectArg.SharedObject = &struct {
		Id                   move_types.AccountAddress
		InitialSharedVersion uint64
		Mutable              bool
	}{
		Id:                   *entry.ObjectID,
		InitialSharedVersion: *entry.InitialSharedVersion,
		Mutable:              mutable,
	}

	return objectArg
}
