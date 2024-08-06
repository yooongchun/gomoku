package tests

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"testing"
)

func TestCache(t *testing.T) {
	c, _ := lru.New[uint64, any](2)

	c.Add(1, 1)
	val, ok := c.Get(1)
	if !ok || val.(int) != 1 {
		t.Errorf("Expected 1 but got %v", val)
	}

	c.Add(2, 2)
	val, ok = c.Get(2)
	if !ok || val.(int) != 2 {
		t.Errorf("Expected 2 but got %v", val)
	}

	c.Add(3, 3)
	val, ok = c.Get(3)

	if c.Contains(1) {
		t.Errorf("Expected false but got true")
	}

	if !c.Contains(2) {
		t.Errorf("Expected true but got false")
	}
}
