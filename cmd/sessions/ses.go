package sessions

import (
	"fmt"
	"sync"
	"time"

	"github.com/ivm97/schedule_izh/models"
)

type Session struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[int64]Item
}

type Item struct {
	Value      models.SearchData
	Created    time.Time
	Expiration int64
}

func New(defaultExpiration, cleanupInterval time.Duration) *Session {
	items := make(map[int64]Item)
	session := Session{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	if cleanupInterval > 0 {
		session.StartGC()
	}

	return &session
}

func (s *Session) Set(key int64, value models.SearchData, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = s.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	s.Lock()
	defer s.Unlock()

	s.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (s *Session) Get(key int64) (models.SearchData, bool) {
	s.RLock()
	defer s.RUnlock()

	item, found := s.items[key]

	if !found {
		return models.SearchData{}, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return models.SearchData{}, false
		}
	}

	return item.Value, true
}

func (s *Session) Delete(key int64) error {
	s.Lock()
	defer s.Unlock()
	if _, found := s.items[key]; !found {
		return fmt.Errorf("key %d not found", key)
	}

	delete(s.items, key)

	return nil
}

func (s *Session) StartGC() {
	go s.GC()
}

func (s *Session) GC() {
	for {
		<-time.After(s.cleanupInterval)

		if s.items == nil {
			return
		}

		if keys := s.expiredKeys(); keys != nil {
			s.clearItems(keys)
		}
	}
}

func (s *Session) expiredKeys() (keys []int64) {
	s.RLock()

	defer s.RUnlock()

	for k, i := range s.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (s *Session) clearItems(keys []int64) {
	s.Lock()

	defer s.Unlock()

	for _, k := range keys {
		delete(s.items, k)
	}
}
