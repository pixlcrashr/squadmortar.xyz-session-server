package session

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/pubsub"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/slice"
	"sync"
)

type SessionChange struct {
	UserLeft      User
	UserJoined    User
	UserChanged   User
	TargetAdded   Target
	TargetChanged Target
	TargetRemoved Target
	WeaponAdded   Weapon
	WeaponChanged Weapon
	WeaponRemoved Weapon
}

type Session interface {
	pubsub.Subscriber[SessionChange]

	Uuid() uuid.UUID

	MaxUsers() int
	MaxWeapons() int
	MaxTargets() int

	SetMaxUsers(v int)
	SetMaxWeapons(v int)
	SetMaxTargets(v int)

	Join(clientUuid uuid.UUID) (User, error)
	Quit(clientUuid uuid.UUID) (User, error)

	Users() []User
	Weapons() []Weapon
	Targets() []Target

	User(clientUuid uuid.UUID) (User, error)
	Weapon(id WeaponId) (Weapon, error)
	Target(id TargetId) (Target, error)

	AddWeapon(weaponType WeaponType) (Weapon, error)
	AddTarget() (Target, error)

	RemoveWeapon(id WeaponId) (Weapon, error)
	RemoveTarget(id TargetId) (Target, error)
}

type session struct {
	uuid uuid.UUID

	maxUsers   int
	maxWeapons int
	maxTargets int

	weaponIdCounter WeaponId
	targetIdCounter TargetId

	users   map[string]User
	weapons map[WeaponId]Weapon
	targets map[TargetId]Target

	mtx sync.RWMutex

	updateSubject pubsub.Subject[SessionChange]
}

func (s *session) Uuid() uuid.UUID {
	return s.uuid
}

func (s *session) nextWeaponId() WeaponId {
	s.weaponIdCounter++
	return s.weaponIdCounter
}

func (s *session) nextTargetId() TargetId {
	s.targetIdCounter++
	return s.targetIdCounter
}

func (s *session) AddWeapon(weaponType WeaponType) (Weapon, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(s.weapons) >= s.maxWeapons {
		return nil, errors.New("maximum weapons per sessions reached")
	}

	id := s.nextWeaponId()
	weapon := newWeapon(id, weaponType)

	weapon.PositionChanged().Add(s.weaponPositionChanged)
	weapon.ActiveChanged().Add(s.weaponActiveChanged)
	weapon.OwnerChanged().Add(s.weaponOwnerChanged)

	s.weapons[id] = weapon

	s.updateSubject.Publish(SessionChange{
		WeaponAdded: weapon,
	})

	return weapon, nil
}

func (s *session) weaponPositionChanged(sender Weapon, args PositionChangedEventArgs) {
	s.updateSubject.Publish(SessionChange{
		WeaponChanged: sender,
	})
}

func (s *session) weaponActiveChanged(sender Weapon, args ActiveChangedEventArgs) {
	s.updateSubject.Publish(SessionChange{
		WeaponChanged: sender,
	})
}

func (s *session) weaponOwnerChanged(sender Weapon, args OwnerChangedEventArgs) {
	s.updateSubject.Publish(SessionChange{
		WeaponChanged: sender,
	})
}

func (s *session) AddTarget() (Target, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(s.targets) >= s.maxTargets {
		return nil, errors.New("maximum targets per sessions reached")
	}

	id := s.nextTargetId()
	target := newTarget(id)

	target.PositionChanged().Add(s.targetPositionChanged)
	target.ActiveChanged().Add(s.targetActiveChanged)
	target.OwnerChanged().Add(s.targetOwnerChanged)

	s.targets[id] = target

	s.updateSubject.Publish(SessionChange{
		TargetAdded: target,
	})

	return target, nil
}

func (s *session) targetPositionChanged(sender Target, args PositionChangedEventArgs) {
	s.updateSubject.Publish(SessionChange{
		TargetChanged: sender,
	})
}

func (s *session) targetActiveChanged(sender Target, args ActiveChangedEventArgs) {
	s.updateSubject.Publish(SessionChange{
		TargetChanged: sender,
	})
}

func (s *session) targetOwnerChanged(sender Target, args OwnerChangedEventArgs) {
	s.updateSubject.Publish(SessionChange{
		TargetChanged: sender,
	})
}

func (s *session) RemoveWeapon(id WeaponId) (Weapon, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	weapon, ok := s.weapons[id]
	if !ok {
		return nil, errors.New("weapon is already removed")
	}

	delete(s.weapons, id)

	weapon.PositionChanged().Remove(s.weaponPositionChanged)
	weapon.ActiveChanged().Remove(s.weaponActiveChanged)
	weapon.OwnerChanged().Remove(s.weaponOwnerChanged)

	s.updateSubject.Publish(SessionChange{
		WeaponRemoved: weapon,
	})

	return weapon, nil
}

func (s *session) RemoveTarget(id TargetId) (Target, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	target, ok := s.targets[id]
	if !ok {
		return nil, errors.New("target is already removed")
	}

	delete(s.targets, id)

	target.PositionChanged().Remove(s.targetPositionChanged)
	target.ActiveChanged().Remove(s.targetActiveChanged)
	target.OwnerChanged().Remove(s.targetOwnerChanged)

	s.updateSubject.Publish(SessionChange{
		TargetRemoved: target,
	})

	return target, nil
}

func (s *session) Subscribe() pubsub.Subscription[SessionChange] {
	return s.updateSubject.Subscribe()
}

func (s *session) MaxUsers() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.maxUsers
}

func (s *session) MaxWeapons() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.maxWeapons
}

func (s *session) MaxTargets() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.maxTargets
}

func (s *session) SetMaxUsers(v int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.maxUsers = v
}

func (s *session) SetMaxWeapons(v int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.maxWeapons = v
}

func (s *session) SetMaxTargets(v int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.maxTargets = v
}

func (s *session) Users() []User {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return slice.MapValuesToSlice(s.users)
}

func (s *session) Weapons() []Weapon {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return slice.MapValuesToSlice(s.weapons)
}

func (s *session) Targets() []Target {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return slice.MapValuesToSlice(s.targets)
}

func (s *session) User(clientUuid uuid.UUID) (User, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	user, ok := s.users[clientUuid.String()]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *session) Weapon(id WeaponId) (Weapon, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	weapon, ok := s.weapons[id]
	if !ok {
		return nil, errors.New("weapon not found")
	}

	return weapon, nil
}

func (s *session) Target(id TargetId) (Target, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	target, ok := s.targets[id]
	if !ok {
		return nil, errors.New("target not found")
	}

	return target, nil
}

func (s *session) Join(clientUuid uuid.UUID) (User, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	user, err := s.join(clientUuid)
	if err != nil {
		return nil, err
	}

	s.updateSubject.Publish(SessionChange{
		UserJoined: user,
	})

	return user, nil
}

func (s *session) join(clientUuid uuid.UUID) (User, error) {
	if _, ok := s.users[clientUuid.String()]; ok {
		return nil, errors.New("client already joined")
	}

	user := newUser(clientUuid, "")
	user.NameChanged().Add(s.userNameChanged)

	s.users[clientUuid.String()] = user

	return user, nil
}

func (s *session) userNameChanged(sender User, args NameChangedEventArgs) {
	s.updateSubject.Publish(SessionChange{
		UserChanged: sender,
	})
}

func (s *session) Quit(clientUuid uuid.UUID) (User, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	user, err := s.quit(clientUuid)
	if err != nil {
		return nil, err
	}

	s.updateSubject.Publish(SessionChange{
		UserLeft: user,
	})

	return s.quit(clientUuid)
}

func (s *session) quit(clientUuid uuid.UUID) (User, error) {
	user, ok := s.users[clientUuid.String()]
	if !ok {
		return nil, errors.New("user already quit")
	}

	delete(s.users, clientUuid.String())

	user.NameChanged().Remove(s.userNameChanged)

	return user, nil
}

func NewSession(uuid uuid.UUID, maxUsers int, maxWeapons int, maxTargets int) Session {
	return &session{
		uuid,
		maxUsers,
		maxWeapons,
		maxTargets,

		0,
		0,

		make(map[string]User, 0),
		make(map[WeaponId]Weapon, 0),
		make(map[TargetId]Target, 0),

		sync.RWMutex{},

		pubsub.NewSubject[SessionChange](),
	}
}
