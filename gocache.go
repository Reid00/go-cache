package gocache

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

// Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}
