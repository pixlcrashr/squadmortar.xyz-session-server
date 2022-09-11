package session

import (
	"github.com/google/uuid"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/eventhandler"
	"sync"
)

type User interface {
	ClientUuid() uuid.UUID
	Name() string
	SetName(name string)

	NameChanged() eventhandler.Event[User, NameChangedEventArgs]
}

type NameChangedEventArgs struct {
	OldName string
	NewName string
}

type user struct {
	clientUuid uuid.UUID
	name       string

	nameChangedEventHandler eventhandler.EventHandler[User, NameChangedEventArgs]

	mtx sync.RWMutex
}

func (u *user) ClientUuid() uuid.UUID {
	return u.clientUuid
}

func (u *user) Name() string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()

	return u.name
}

func (u *user) SetName(name string) {
	u.mtx.Lock()

	old := u.name
	u.name = name

	u.mtx.Unlock()

	u.nameChangedEventHandler.Invoke(u, NameChangedEventArgs{
		OldName: old,
		NewName: name,
	})
}

func (u *user) NameChanged() eventhandler.Event[User, NameChangedEventArgs] {
	return u.nameChangedEventHandler
}

func newUser(clientUuid uuid.UUID, name string) User {
	return &user{
		clientUuid,
		name,
		eventhandler.New[User, NameChangedEventArgs](),
		sync.RWMutex{},
	}
}
