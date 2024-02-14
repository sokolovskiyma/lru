package lru

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	// Create a new LRU cache with a capacity of 3
	lru := newLRUCache[string, int](3)

	// Test 1: Adding and getting values
	lru.Update("key1", 1)
	lru.Update("key2", 2)
	lru.Update("key3", 3)

	// Check the length after adding 3 entries
	if lru.length != 3 {
		t.Errorf("Expected length: 3, Got: %d", lru.length)
	}

	// Check the values retrieved from the cache
	checkValue(t, &lru, 1, "key1")
	checkValue(t, &lru, 2, "key2")
	checkValue(t, &lru, 3, "key3")

	// Test 2: Updating an existing value
	lru.Update("key2", 22)

	// Check the length after updating an entry
	if lru.length != 3 {
		t.Errorf("Expected length: 3, Got: %d", lru.length)
	}

	// Check the updated value
	checkValue(t, &lru, 22, "key2")

	// Test 3: Adding more values to trigger eviction
	lru.Update("key4", 4)
	lru.Update("key5", 5)

	// Check the length after adding 2 more entries
	if lru.length != 3 {
		t.Errorf("Expected length: 3, Got: %d", lru.length)
	}

	// Check that the least recently used entry was evicted
	checkValue(t, &lru, 0, "key1")

	// Check the values of remaining entries
	checkValue(t, &lru, 22, "key2")
	checkValue(t, &lru, 4, "key4")
	checkValue(t, &lru, 5, "key5")
}

func checkValue[K comparable, V comparable](t *testing.T, lru *LRUCache[K, V], expected V, key K) {
	actual, exists := lru.Get(key)
	if !exists || actual != expected {
		t.Errorf("Expected value for %v: %v, Got: %v", key, expected, actual)
	}
}

/* BENCHMARKS */

func BenchmarkSetGet(b *testing.B) {
	// preparation
	lru := newLRUCache[int, int](10_000)

	// test
	for i := 0; i < b.N; i++ {
		lru.Update(i, i)
		value, exists := lru.Get(i)
		if !exists || value != i {
			b.Fatalf("value does not exist or %v != %v", value, i)
		}
	}
}

func BenchmarkSetGetT(b *testing.B) {
	// preparation
	lru := newLRUCacheT[int, int](10_000)

	// test
	for i := 0; i < b.N; i++ {
		lru.Update(i, i)
		value, exists := lru.Get(i)
		if !exists || value != i {
			b.Fatalf("value does not exist or %v != %v", value, i)
		}
	}
}
