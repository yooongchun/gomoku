package ai

import (
	"testing"
)

func TestCache(t *testing.T) {
	c := NewCache(2)

	key1 := uint64(123)
	c.Put(key1, 1)
	if c.Get(key1).(int) != 1 {
		t.Errorf("Expected 1 but got %v", c.Get(key1))
	}
	key2 := uint64(234)
	c.Put(key2, 2)
	if c.Get(key2).(int) != 2 {
		t.Errorf("Expected 2 but got %v", c.Get(key2))
	}
	key3 := uint64(345)
	c.Put(key3, 3)
	if c.Get(key3).(int) != 3 {
		t.Errorf("Expected 3 but got %v", c.Get(key3))
	}

	if c.Get(key1) != nil {
		t.Errorf("Expected 0 but got %v", c.Get(key1))
	}

	if c.Has(key1) {
		t.Errorf("Expected false but got true")
	}

	if !c.Has(key2) {
		t.Errorf("Expected true but got false")
	}
}
