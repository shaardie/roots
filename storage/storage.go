package storage

import (
	"fmt"
	"sync"
)

type WorldMap struct {
	Countries map[string]int
	m         sync.Mutex
}

type Storage struct {
	Storage map[string]*WorldMap
	m       *sync.Mutex
}

func New() *Storage {
	return &Storage{
		Storage: make(map[string]*WorldMap),
		m:       &sync.Mutex{},
	}
}

func (s *Storage) New(name string) (*WorldMap, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.Storage[name]; ok {
		return nil, fmt.Errorf("key already exists")
	}

	m := NewWorldMap()
	s.Storage[name] = m
	return m, nil
}

func (s *Storage) Get(name string) (*WorldMap, error) {
	s.m.Lock()
	defer s.m.Unlock()

	m, ok := s.Storage[name]
	if !ok {
		return nil, fmt.Errorf("key does not exist")
	}
	return m, nil
}

func NewWorldMap() *WorldMap {
	return &WorldMap{
		Countries: make(map[string]int),
		m:         sync.Mutex{},
	}
}

func (w *WorldMap) Add(countryCode string) error {
	w.m.Lock()
	defer w.m.Unlock()

	if _, ok := w.Countries[countryCode]; ok {
		w.Countries[countryCode]++
	} else {
		w.Countries[countryCode] = 1
	}

	return nil
}
