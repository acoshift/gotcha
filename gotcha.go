package gotcha

import (
	"sync"
	"time"
)

// Gotcha type
type Gotcha struct {
	m *sync.RWMutex
	d map[interface{}]*item
}

type item struct {
	v   interface{}
	ts  time.Time
	ttl time.Duration
}

// New creates a gotcha
func New() *Gotcha {
	return &Gotcha{
		m: &sync.RWMutex{},
		d: map[interface{}]*item{},
	}
}

// Get retrieves data for an index
// return nil if not found
func (g *Gotcha) Get(index interface{}) interface{} {
	g.m.RLock()
	defer g.m.RUnlock()
	it := g.d[index]
	if it == nil || expired(it) {
		return nil
	}
	return it.v
}

// MustGet gets data for an index even if it already expired
func (g *Gotcha) MustGet(index interface{}) interface{} {
	g.m.RLock()
	defer g.m.RUnlock()
	it := g.d[index]
	if it == nil {
		return nil
	}
	return it.v
}

// GetMulti retrieves multiple data
// faster than Get
func (g *Gotcha) GetMulti(indexes []interface{}) []interface{} {
	its := make([]interface{}, len(indexes))
	g.m.RLock()
	defer g.m.RUnlock()
	for i := range its {
		it := g.d[indexes[i]]
		if it != nil && !expired(it) {
			its[i] = it.v
		}
	}
	return its
}

// Set sets data for an index
func (g *Gotcha) Set(index interface{}, value interface{}) {
	g.SetTTL(index, value, 0)
}

// SetTTL sets data for an index with ttl
func (g *Gotcha) SetTTL(index interface{}, value interface{}, ttl time.Duration) {
	g.m.Lock()
	defer g.m.Unlock()
	it := g.d[index]
	if it == nil {
		it = &item{}
		g.d[index] = it
	}
	it.ts = time.Now()
	it.ttl = ttl
	it.v = value
}

// SetMulti sets data for indexes
func (g *Gotcha) SetMulti(indexes []interface{}, values []interface{}) {
	g.SetMultiTTL(indexes, values, 0)
}

// SetMultiTTL sets data for indexes with a ttl
func (g *Gotcha) SetMultiTTL(indexes []interface{}, values []interface{}, ttl time.Duration) {
	g.m.Lock()
	defer g.m.Unlock()
	now := time.Now()
	for i := range indexes {
		index := indexes[i]
		it := g.d[index]
		if it == nil {
			it = &item{}
			g.d[index] = it
		}
		it.ts = now
		it.ttl = ttl
		it.v = values[i]
	}
}

// Unset removes data for an index
func (g *Gotcha) Unset(index interface{}) {
	g.m.Lock()
	defer g.m.Unlock()
	delete(g.d, index)
}

// Purge removes all data
func (g *Gotcha) Purge() {
	g.m.Lock()
	defer g.m.Unlock()
	g.d = map[interface{}]*item{}
}

// Extend replaces old ttl with the new one if value not expired
func (g *Gotcha) Extend(index interface{}, ttl time.Duration) {
	g.m.Lock()
	defer g.m.Unlock()
	it := g.d[index]
	if it != nil && !expired(it) {
		it.ts = time.Now()
		it.ttl = ttl
	}
}

// MustExtend force replace old ttl with the new one even if it already expired
func (g *Gotcha) MustExtend(index interface{}, ttl time.Duration) {
	g.m.Lock()
	defer g.m.Unlock()
	it := g.d[index]
	if it != nil {
		it.ts = time.Now()
		it.ttl = ttl
	}
}

// Exists checks is index exists but not check is expired
func (g *Gotcha) Exists(index interface{}) bool {
	g.m.RLock()
	defer g.m.RUnlock()
	return g.d[index] != nil
}

// Expired checks is index expired
// return false if expired or not exists
func (g *Gotcha) Expired(index interface{}) bool {
	g.m.RLock()
	defer g.m.RUnlock()
	it := g.d[index]
	return it != nil && expired(it)
}

// Cleanup removes all expired data
func (g *Gotcha) Cleanup() {
	g.m.Lock()
	defer g.m.Unlock()
	for index, it := range g.d {
		if expired(it) {
			delete(g.d, index)
		}
	}
}

// Timestamp returns timestamp of an index
func (g *Gotcha) Timestamp(index interface{}) *time.Time {
	g.m.RLock()
	defer g.m.RUnlock()
	it := g.d[index]
	if it == nil {
		return nil
	}
	return &it.ts
}

// helper funcitons
func expired(it *item) bool {
	if it.ttl <= 0 {
		return false
	}
	return time.Now().After(it.ts.Add(it.ttl))
}
