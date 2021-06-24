package main

type memoryDB struct {
	maps map[string]savedToken
}

func newMemoryDB() *memoryDB {
	return &memoryDB{
		maps: map[string]savedToken{},
	}
}
