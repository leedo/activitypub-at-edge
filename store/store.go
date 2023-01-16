package store

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fastly/compute-sdk-go/objectstore"
)

const lockTTL = time.Second * 60

type Store struct {
	s     *objectstore.Store
	locks map[string]*Lock
	id    string
}

type Lock struct {
	id   string
	time time.Time
}

func Open(name string) (*Store, error) {
	s, err := objectstore.Open(name)
	if err != nil {
		return nil, err
	}
	return &Store{
		s:     s,
		id:    os.Getenv("FASTLY_TRACE_ID"),
		locks: make(map[string]*Lock),
	}, nil
}

func (l Lock) Age() time.Duration {
	return time.Now().Sub(l.time)
}

func (l Lock) Expired() bool {
	return l.Age() > lockTTL
}

func (l Lock) String() string {
	return fmt.Sprintf("%s:%d", l.id, time.Now().Unix())
}

var KeyNotFound = errors.New("key not found")

func (s *Store) Lookup(key string) (string, error) {
	e, err := s.s.Lookup(key)
	if err != nil {
		return "", err
	}

	v := e.String()
	if v == "" {
		return "", KeyNotFound
	}

	return v, nil
}

func (s *Store) Insert(key string, value string) error {
	return s.s.Insert(key, strings.NewReader(value))
}

func (s *Store) readLock(key string) *Lock {
	e, err := s.s.Lookup(key)

	if err != nil {
		return nil
	}

	v := e.String()
	if v == "" {
		return nil
	}

	parts := strings.SplitN(v, ":", 2)
	t, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil
	}

	return &Lock{
		id:   parts[0],
		time: time.Unix(t, 0),
	}
}

func (s *Store) Lock(name string) error {
	key := name + ".lock"

	// we already have this lock
	if _, ok := s.locks[key]; ok {
		return nil
	}

	lock := s.readLock(key)

	// lock already exists
	if lock != nil {
		if lock.id != s.id && !lock.Expired() {
			return fmt.Errorf("lock already exists with id: %s", lock.id)
		}
	} else {
		lock = &Lock{s.id, time.Now()}
	}

	if err := s.s.Insert(key, strings.NewReader(lock.String())); err != nil {
		return err
	}

	s.locks[name] = lock
	return nil
}

func (s *Store) Unlock(name string) error {
	key := name + ".lock"

	// we don't have this lock
	if _, ok := s.locks[key]; !ok {
		return nil
	}

	lock := s.readLock(key)

	// lock already exists?
	if lock != nil {
		if lock.id != s.id && !lock.Expired() {
			return fmt.Errorf("cannot remove lock from another instance: %s", lock.id)
		}
	}

	if err := s.s.Insert(key, strings.NewReader("")); err != nil {
		return err
	}

	delete(s.locks, key)
	return nil
}
