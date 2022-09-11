package session

import (
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/eventhandler"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/math"
	"sync"
)

type WeaponId int32

type WeaponType int32

const (
	StandardMortarWeaponType WeaponType = iota
	TechnicalMortarWeaponType
	RocketsWeaponType
	HellCannonWeaponType
)

type Weapon interface {
	Id() WeaponId
	Type() WeaponType
	Position() math.Vector3
	SetPosition(v math.Vector3)
	AddPosition(v math.Vector3)
	PositionChanged() eventhandler.Event[Weapon, PositionChangedEventArgs]
	Active() bool
	SetActive(v bool)
	ActiveChanged() eventhandler.Event[Weapon, ActiveChangedEventArgs]
	Owner() User
	SetOwner(User)
	OwnerChanged() eventhandler.Event[Weapon, OwnerChangedEventArgs]
	IsOwned() bool
}

type weapon struct {
	id       WeaponId
	typ      WeaponType
	position math.Vector3
	active   bool
	owner    User

	positionEventHandler eventhandler.EventHandler[Weapon, PositionChangedEventArgs]
	activeEventHandler   eventhandler.EventHandler[Weapon, ActiveChangedEventArgs]
	ownerEventHandler    eventhandler.EventHandler[Weapon, OwnerChangedEventArgs]

	mtx sync.RWMutex
}

func (w *weapon) Id() WeaponId {
	return w.id
}

func (w *weapon) Type() WeaponType {
	return w.typ
}

func (w *weapon) Position() math.Vector3 {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	return w.position
}

func (w *weapon) SetPosition(v math.Vector3) {
	w.mtx.Lock()

	old := w.position
	w.position = v

	w.mtx.Unlock()

	w.positionEventHandler.Invoke(w, PositionChangedEventArgs{
		OldPosition: old,
		NewPosition: w.position,
	})
}

func (w *weapon) AddPosition(v math.Vector3) {
	w.mtx.Lock()

	old := w.position
	w.position = w.position.Add(v)

	w.mtx.Unlock()

	w.positionEventHandler.Invoke(w, PositionChangedEventArgs{
		OldPosition: old,
		NewPosition: w.position,
	})
}

func (w *weapon) Active() bool {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	return w.active
}

func (w *weapon) SetActive(v bool) {
	w.mtx.Lock()

	old := w.active
	w.active = v

	w.mtx.Unlock()

	w.activeEventHandler.Invoke(w, ActiveChangedEventArgs{
		OldActive: old,
		NewActive: v,
	})
}

func (w *weapon) Owner() User {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	return w.owner
}

func (w *weapon) SetOwner(u User) {
	w.mtx.Lock()

	old := w.owner
	w.owner = u

	w.mtx.Unlock()

	w.ownerEventHandler.Invoke(w, OwnerChangedEventArgs{
		OldUser: old,
		NewUser: u,
	})
}

func (w *weapon) IsOwned() bool {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	return w.owner != nil
}

func (w *weapon) PositionChanged() eventhandler.Event[Weapon, PositionChangedEventArgs] {
	return w.positionEventHandler
}

func (w *weapon) ActiveChanged() eventhandler.Event[Weapon, ActiveChangedEventArgs] {
	return w.activeEventHandler
}

func (w *weapon) OwnerChanged() eventhandler.Event[Weapon, OwnerChangedEventArgs] {
	return w.ownerEventHandler
}

func newWeapon(id WeaponId, typ WeaponType) Weapon {
	return &weapon{
		id,
		typ,
		math.Vector3{},
		false,
		nil,
		eventhandler.New[Weapon, PositionChangedEventArgs](),
		eventhandler.New[Weapon, ActiveChangedEventArgs](),
		eventhandler.New[Weapon, OwnerChangedEventArgs](),
		sync.RWMutex{},
	}
}
