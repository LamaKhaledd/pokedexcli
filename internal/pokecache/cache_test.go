package pokecache

import (
	"testing"
	"time"
)

func TestCache_AddAndGet(t *testing.T) {
	cache := NewCache(2 * time.Second)

	key := "https://pokeapi.co/api/v2/location-area"
	val := []byte("cached data")

	cache.Add(key, val)

	got, found := cache.Get(key)
	if !found {
		t.Fatal("Expected to find key in cache")
	}

	if string(got) != string(val) {
		t.Fatalf("Expected '%s', got '%s'", val, got)
	}

	time.Sleep(3 * time.Second)

	_, found = cache.Get(key)
	if found {
		t.Fatal("Expected key to be removed after reaping")
	}
}

