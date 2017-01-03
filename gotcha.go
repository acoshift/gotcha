package gotcha

import (
	"sync"
	"time"
)

// Gotcha type
type Gotcha struct {
	m *sync.RWMutex
	d map[interface{}]*item
	f map[interface{}]func()
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
		f: map[interface{}]func(){},
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
	g.f = map[interface{}]func(){}
}

// Filler registers a filler for an index
func (g *Gotcha) Filler(index interface{}, f func()) {
	g.m.RLock()
	defer g.m.RUnlock()
	g.f[index] = f
}
