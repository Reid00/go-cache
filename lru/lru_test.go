package lru

import (
	"reflect"
	"testing"
)

type String string

// Value interface
func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)

	lru.Add("key1", String("1234"))

	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "v1", "v2", "v3"

	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)

	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	// key1 should be removed
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(10), callback)

	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expected := []string{"key1", "k2"}

	if !reflect.DeepEqual(expected, keys) {
		t.Fatalf("call OnEvicted failed, expect keys=%v, actual=%v", expected, keys)
	}

}

func TestAdd(t *testing.T) {
	lru := New(int64(0), nil)

	lru.Add("key", String("1"))
	lru.Add("key", String("22"))

	if lru.nbyte != int64(len("key")+len("22")) {
		t.Fatalf("expected=%v, but got = %v", len("key")+len("22"), lru.nbyte)
	}
}
