package gotcha

import (
	"testing"
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
	if g.Get(1) != nil || g.Get(2) != nil || len(g.d) > 0 || len(g.f) > 0 {
		t.Errorf("expected purged data to be nil")
	}
}
