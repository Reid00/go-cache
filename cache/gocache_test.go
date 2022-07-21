package cache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "ABC",
	"Jack": "DEF",
	"Reid": "GHI",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expected := []byte("key")

	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expected) {
		t.Errorf("calback failed")
	}
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))

	goCache := NewGroup("scores", 1<<12, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowBD] search key: ", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := goCache.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", view)
		}

		if _, err := goCache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := goCache.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown key should be empty, but got %s", view)
	}
}
