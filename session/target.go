package session

import (
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/eventhandler"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/math"
	"sync"
)

type TargetId int32

type Target interface {
	Id() TargetId
	Position() math.Vector3
	SetPosition(v math.Vector3)
	AddPosition(v math.Vector3)
	PositionChanged() eventhandler.Event[Target, PositionChangedEventArgs]
	Active() bool
	SetActive(v bool)
	ActiveChanged() eventhandler.Event[Target, ActiveChangedEventArgs]
	Owner() User
	SetOwner(User)
	OwnerChanged() eventhandler.Event[Target, OwnerChangedEventArgs]
	IsOwned() bool
}

type PositionChangedEventArgs struct {
	OldPosition math.Vector3
	NewPosition math.Vector3
}

type ActiveChangedEventArgs struct {
	OldActive bool
	NewActive bool
}

type OwnerChangedEventArgs struct {
	OldUser User
	NewUser User
}

type target struct {
	id       TargetId
	position math.Vector3
	active   bool
	owner    User

	positionEventHandler eventhandler.EventHandler[Target, PositionChangedEventArgs]
	activeEventHandler   eventhandler.EventHandler[Target, ActiveChangedEventArgs]
	ownerEventHandler    eventhandler.EventHandler[Target, OwnerChangedEventArgs]

	mtx sync.RWMutex
}

func (t *target) Id() TargetId {
	return t.id
}

func (t *target) Position() math.Vector3 {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.position
}

func (t *target) SetPosition(v math.Vector3) {
	t.mtx.Lock()

	old := t.position
	t.position = v

	t.mtx.Unlock()

	t.positionEventHandler.Invoke(t, PositionChangedEventArgs{
		OldPosition: old,
		NewPosition: t.position,
	})
}

func (t *target) AddPosition(v math.Vector3) {
	t.mtx.Lock()

	old := t.position
	t.position = t.position.Add(v)

	t.mtx.Unlock()

	t.positionEventHandler.Invoke(t, PositionChangedEventArgs{
		OldPosition: old,
		NewPosition: t.position,
	})
}

func (t *target) Active() bool {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.active
}

func (t *target) SetActive(v bool) {
	t.mtx.Lock()

	old := t.active
	t.active = v

	t.mtx.Unlock()

	t.activeEventHandler.Invoke(t, ActiveChangedEventArgs{
		OldActive: old,
		NewActive: v,
	})
}

func (t *target) Owner() User {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.owner
}

func (t *target) SetOwner(u User) {
	t.mtx.Lock()

	old := t.owner
	t.owner = u

	t.mtx.Unlock()

	t.ownerEventHandler.Invoke(t, OwnerChangedEventArgs{
		OldUser: old,
		NewUser: u,
	})
}

func (t *target) IsOwned() bool {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.owner != nil
}

func (t *target) PositionChanged() eventhandler.Event[Target, PositionChangedEventArgs] {
	return t.positionEventHandler
}

func (t *target) ActiveChanged() eventhandler.Event[Target, ActiveChangedEventArgs] {
	return t.activeEventHandler
}

func (t *target) OwnerChanged() eventhandler.Event[Target, OwnerChangedEventArgs] {
	return t.ownerEventHandler
}

func newTarget(id TargetId) Target {
	return &target{
		id,
		math.Vector3{},
		false,
		nil,
		eventhandler.New[Target, PositionChangedEventArgs](),
		eventhandler.New[Target, ActiveChangedEventArgs](),
		eventhandler.New[Target, OwnerChangedEventArgs](),
		sync.RWMutex{},
	}
}
