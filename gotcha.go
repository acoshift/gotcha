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
	ts  int64
	ttl int64
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
		if it != nil {
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
func (g *Gotcha) SetTTL(index interface{}, value interface{}, ttl int64) {
	g.m.Lock()
	defer g.m.Unlock()
	it := g.d[index]
	if it == nil {
		it = &item{}
		g.d[index] = it
	}
	it.ts = time.Now().UnixNano()
	it.ttl = ttl
	it.v = value
}

// SetMulti sets data for indexes
func (g *Gotcha) SetMulti(indexes []interface{}, values []interface{}) {
	g.SetMultiTTL(indexes, values, 0)
}

// SetMultiTTL sets data for indexes with a ttl
func (g *Gotcha) SetMultiTTL(indexes []interface{}, values []interface{}, ttl int64) {
	g.m.Lock()
	defer g.m.Unlock()
	now := time.Now().UnixNano()
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

// Extend replace old ttl with the new one
func (g *Gotcha) Extend(index interface{}, ttl int64) {
	g.m.Lock()
	defer g.m.Unlock()
	it := g.d[index]
	if it != nil {
		it.ts = time.Now().UnixNano()
		it.ttl = ttl
	}
}
