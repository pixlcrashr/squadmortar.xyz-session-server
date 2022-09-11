package storage

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/session"
	"sync"
)

type Storage interface {
	Create() (session.Session, error)
	Delete(uuid uuid.UUID) error
	Get(uuid uuid.UUID) (session.Session, error)
}

type storage struct {
	sessions map[string]session.Session

	mtx sync.Mutex
}

func (s *storage) Create() (session.Session, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	session := session.NewSession(
		uuid,
		30,
		200,
		200,
	)
	s.sessions[uuid.String()] = session

	return session, nil
}

func (s *storage) Get(uuid uuid.UUID) (session.Session, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	session, ok := s.sessions[uuid.String()]
	if !ok {
		return nil, errors.New("session does not exist")
	}

	return session, nil
}

func (s *storage) Delete(uuid uuid.UUID) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, ok := s.sessions[uuid.String()]
	if !ok {
		return errors.New("session was already deleted")
	}

	delete(s.sessions, uuid.String())

	return nil
}

func NewStorage() Storage {
	return &storage{
		make(map[string]session.Session, 0),
		sync.Mutex{},
	}
}
