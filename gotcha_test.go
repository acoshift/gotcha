package gotcha

import (
	"reflect"
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	cases := []struct {
		index interface{}
		value interface{}
	}{
		{1, 1},
		{2, "test"},
		{"test", 123},
		{"test2", "test"},
		{nil, 345},
		{1.3, nil},
		{1.6, 1.3},
		{1, 5},
		{1, "test"},
		{1, nil},
		{1, 1.4},
		{true, false},
	}

	g := New()
	for _, c := range cases {
		g.Set(c.index, c.value)
		v := g.Get(c.index)
		if v != c.value {
			t.Errorf("expected value to be %v; got %v", c.value, v)
		}
	}

	// try get not seted value
	if g.Get(12345) != nil {
		t.Errorf("expected not set value to be nil")
	}

	if g.MustGet(9999) != nil {
		t.Errorf("expected not set value to be nil")
	}
}

func TestUnset(t *testing.T) {
	g := New()
	g.Set(1, 5)
	g.Set(2, 6)
	g.Unset(1)
	if g.Get(1) != nil {
		t.Errorf("expected unseted value to be nil")
	}
}

func TestPurge(t *testing.T) {
	g := New()
	g.Set(1, 2)
	g.Set(2, 3)
	g.Purge()
	if g.Get(1) != nil || g.Get(2) != nil || len(g.d) > 0 {
		t.Errorf("expected purged data to be nil")
	}
}

func TestMulti(t *testing.T) {
	indexes := []interface{}{1, 2, 3, "a", "aaaaaaaaaaaaaa", 1.3}
	values := []interface{}{"aaaaaa", "b", 3, 5, 2.9, nil}

	g := New()
	g.SetMulti(indexes, values)
	res := g.GetMulti(indexes)
	if !reflect.DeepEqual(res, values) {
		t.Errorf("expected values to be %v; got %v", values, res)
	}
}

func TestExtend(t *testing.T) {
	g := New()
	g.SetTTL(1, 2, time.Millisecond*10)
	g.Extend(1, time.Millisecond*20)

	if g.d[1].ttl != time.Millisecond*20 {
		t.Errorf("expected ttl to be %v; got %v", time.Millisecond*20, g.d[1].ttl)
	}

	time.Sleep(time.Millisecond * 30)
	if g.Get(1) != nil || g.MustGet(1) != 2 {
		t.Errorf("expected must get can get expired index")
	}
	g.MustExtend(1, time.Millisecond*10)
	if g.Expired(1) {
		t.Errorf("expected must extend to extend expired index")
	}
}

func TestExists(t *testing.T) {
	g := New()
	g.Set(1, true)
	g.Set(2, nil)
	if !g.Exists(1) {
		t.Errorf("expected 1 exists")
	}
	if !g.Exists(2) {
		t.Errorf("expected 2 exists")
	}
	if g.Exists(3) {
		t.Errorf("expected 3 not exists")
	}
}

func TestExpired(t *testing.T) {
	g := New()
	g.SetTTL(1, 2, time.Millisecond*10)
	if g.Expired(1) {
		t.Errorf("expected 1 not expired")
	}
	time.Sleep(time.Millisecond * 15)
	if !g.Expired(1) {
		t.Errorf("expected 1 expired")
	}
}

func TestCleanup(t *testing.T) {
	g := New()
	g.SetTTL(1, 2, time.Millisecond*10)
	g.SetTTL(2, 3, time.Millisecond*20)
	time.Sleep(time.Millisecond * 15)
	g.Cleanup()
	if g.Exists(1) {
		t.Errorf("expected 1 not exists")
	}
	if !g.Exists(2) {
		t.Errorf("expected 2 exists")
	}
}

func TestTimestamp(t *testing.T) {
	g := New()
	g.SetTTL(1, 2, time.Millisecond*10)
	ts := g.Timestamp(1)
	if ts == nil {
		t.Errorf("expected timestamp not nil")
	}
	ts = g.Timestamp(2)
	if ts != nil {
		t.Errorf("expected timestamp to be nil")
	}
}
