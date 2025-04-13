package cache

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	c := NewLRUCache(2, 2*time.Second)

	c.Set("a", 1)
	val, ok := c.Get("a")
	if !ok {
		t.Errorf("Expected to get value for key 'a'")
	}
	if val.(int) != 1 {
		t.Errorf("Expected value 1, got %v", val)
	}
}

func TestOverwriteExistingKey(t *testing.T) {
	c := NewLRUCache(2, 2*time.Second)

	c.Set("a", 1)
	c.Set("a", 2)
	val, ok := c.Get("a")
	if !ok || val.(int) != 2 {
		t.Errorf("Expected updated value 2, got %v", val)
	}
}

func TestExpiration(t *testing.T) {
	c := NewLRUCache(2, 500*time.Millisecond)

	c.Set("a", 1)
	time.Sleep(600 * time.Millisecond)
	val, ok := c.Get("a")
	if ok || val != nil {
		t.Errorf("Expected key 'a' to be expired")
	}
}

func TestEviction(t *testing.T) {
	c := NewLRUCache(2, 5*time.Second)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3) // This should evict "a" (least recently used)

	if _, ok := c.Get("a"); ok {
		t.Errorf("Expected 'a' to be evicted")
	}

	if val, ok := c.Get("b"); !ok || val.(int) != 2 {
		t.Errorf("Expected 'b' to still be present")
	}

	if val, ok := c.Get("c"); !ok || val.(int) != 3 {
		t.Errorf("Expected 'c' to be present")
	}
}

func TestUpdateTTL(t *testing.T) {
	c := NewLRUCache(1, 1*time.Second)

	c.Set("a", 1)
	time.Sleep(500 * time.Millisecond)

	c.Set("a", 2) // Update value and TTL

	time.Sleep(700 * time.Millisecond) // total 1.2s, but TTL renewed

	val, ok := c.Get("a")
	if !ok || val.(int) != 2 {
		t.Errorf("Expected renewed TTL and value 2, got %v", val)
	}
}
