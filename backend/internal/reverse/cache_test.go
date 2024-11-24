package reverse

import (
	"testing"
)

func TestCacheSetAndGet(t *testing.T) {
	cache := NewCache(3)

	// Set entries in the cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Retrieve entries and verify values
	if value, found := cache.Get("key1"); !found || value != "value1" {
		t.Errorf("Expected to find key1 with value 'value1', got found=%v, value=%v", found, value)
	}

	if value, found := cache.Get("key2"); !found || value != "value2" {
		t.Errorf("Expected to find key2 with value 'value2', got found=%v, value=%v", found, value)
	}

	if value, found := cache.Get("key3"); !found || value != "value3" {
		t.Errorf("Expected to find key3 with value 'value3', got found=%v, value=%v", found, value)
	}
}

func TestCacheEviction(t *testing.T) {
	cache := NewCache(2)

	// Add entries beyond capacity
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	evicted := cache.Set("key3", "value3")

	// Verify eviction occurred
	if !evicted {
		t.Error("Expected an eviction when setting key3")
	}

	// Verify that the oldest entry (key1) was evicted
	if _, found := cache.Get("key1"); found {
		t.Error("Expected key1 to be evicted, but it was found")
	}

	// Verify that other entries are still present
	if value, found := cache.Get("key2"); !found || value != "value2" {
		t.Errorf("Expected to find key2 with value 'value2', got found=%v, value=%v", found, value)
	}
	if value, found := cache.Get("key3"); !found || value != "value3" {
		t.Errorf("Expected to find key3 with value 'value3', got found=%v, value=%v", found, value)
	}
}

func TestCacheUpdate(t *testing.T) {
	cache := NewCache(2)

	// Add an entry and then update it
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key1", "updated_value1")

	// Verify updated value
	if value, found := cache.Get("key1"); !found || value != "updated_value1" {
		t.Errorf("Expected key1 to have updated value 'updated_value1', got found=%v, value=%v", found, value)
	}
}

func TestCacheLRUBehavior(t *testing.T) {
	cache := NewCache(2)

	// Add entries and access them to update LRU order
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Get("key1") // Access key1 to make it most recently used
	evicted := cache.Set("key3", "value3")

	// Verify eviction occurred
	if !evicted {
		t.Error("Expected an eviction when setting key3")
	}

	// Verify that the least recently used entry (key2) was evicted
	if _, found := cache.Get("key2"); found {
		t.Error("Expected key2 to be evicted, but it was found")
	}

	// Verify that other entries are still present
	if value, found := cache.Get("key1"); !found || value != "value1" {
		t.Errorf("Expected to find key1 with value 'value1', got found=%v, value=%v", found, value)
	}
	if value, found := cache.Get("key3"); !found || value != "value3" {
		t.Errorf("Expected to find key3 with value 'value3', got found=%v, value=%v", found, value)
	}
}

func TestCacheCapacityZero(t *testing.T) {
	cache := NewCache(0)

	// Attempt to add an entry to a cache with capacity 0
	evicted := cache.Set("key1", "value1")

	// Verify eviction occurred
	if evicted {
		t.Error("Nothing should be evicted, as nothing was set")
	}

	// Verify no entries can be retrieved
	if _, found := cache.Get("key1"); found {
		t.Error("Expected no entries to be stored in a cache with capacity 0")
	}
}
