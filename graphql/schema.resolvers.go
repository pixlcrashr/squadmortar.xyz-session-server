package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/google/uuid"
	auth2 "github.com/pixlcrashr/squadmortar.xyz-sessions-server/auth"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/graphql/generated"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/graphql/model"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/math"
	session3 "github.com/pixlcrashr/squadmortar.xyz-sessions-server/session"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/slice"
)

var ErrUserNotInSession = errors.New("user is not in session")

// Authenticate is the resolver for the authenticate field.
func (r *mutationResolver) Authenticate(ctx context.Context) (string, error) {
	clientUuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return auth2.GenerateToken(clientUuid, r.EcdsaKey)
}

// ChangeUserName is the resolver for the changeUserName field.
func (r *mutationResolver) ChangeUserName(ctx context.Context, sessionGUID string, name string) (*model.User, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	uuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(uuid)
	if err != nil {
		return nil, err
	}

	user, err := session.User(clientUuid)
	if err != nil {
		return nil, ErrUserNotInSession
	}

	user.SetName(name)

	return UserToGraphQL(user), nil
}

// CreateSession is the resolver for the createSession field.
func (r *mutationResolver) CreateSession(ctx context.Context) (*model.Session, error) {
	_, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	session, err := r.SessionStorage.Create()
	if err != nil {
		return nil, err
	}

	return SessionToGraphQL(session), nil
}

// JoinSession is the resolver for the joinSession field.
func (r *mutationResolver) JoinSession(ctx context.Context, sessionGUID string) (*model.User, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	user, err := session.Join(clientUuid)
	if err != nil {
		return nil, err
	}

	return UserToGraphQL(user), nil
}

// QuitSession is the resolver for the quitSession field.
func (r *mutationResolver) QuitSession(ctx context.Context, sessionGUID string) (string, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return "", auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return "", err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return "", err
	}

	_, err = session.Quit(clientUuid)
	if err != nil {
		return "", err
	}

	return sessionGUID, nil
}

// AddWeapon is the resolver for the addWeapon field.
func (r *mutationResolver) AddWeapon(ctx context.Context, sessionGUID string, weaponType model.WeaponType) (*model.Weapon, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	weapon, err := session.AddWeapon(WeaponTypeFromGraphQL(weaponType))
	if err != nil {
		return nil, err
	}

	return WeaponToGraphQL(weapon), nil
}

// AddTarget is the resolver for the addTarget field.
func (r *mutationResolver) AddTarget(ctx context.Context, sessionGUID string) (*model.Target, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	target, err := session.AddTarget()
	if err != nil {
		return nil, err
	}

	return TargetToGraphQL(target), nil
}

// Target is the resolver for the target field.
func (r *mutationResolver) Target(ctx context.Context, sessionGUID string, input model.TargetInput) (*model.Target, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	target, err := session.Target(session3.TargetId(input.ID))
	if err != nil {
		return nil, err
	}

	if input.Active != nil {
		target.SetActive(*input.Active)
	}

	if input.Position != nil {
		target.SetPosition(Vector3InputFromGraphQL(*input.Position))
	}

	return TargetToGraphQL(target), nil
}

// Weapon is the resolver for the weapon field.
func (r *mutationResolver) Weapon(ctx context.Context, sessionGUID string, input model.WeaponInput) (*model.Weapon, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	weapon, err := session.Weapon(session3.WeaponId(input.ID))
	if err != nil {
		return nil, err
	}

	if input.Active != nil {
		weapon.SetActive(*input.Active)
	}

	if input.Position != nil {
		weapon.SetPosition(Vector3InputFromGraphQL(*input.Position))
	}

	return WeaponToGraphQL(weapon), nil
}

// AcquireTarget is the resolver for the acquireTarget field.
func (r *mutationResolver) AcquireTarget(ctx context.Context, sessionGUID string, id int) (*model.Target, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	user, err := session.User(clientUuid)
	if err != nil {
		return nil, ErrUserNotInSession
	}

	target, err := session.Target(session3.TargetId(id))
	if err != nil {
		return nil, err
	}

	target.SetOwner(user)

	return TargetToGraphQL(target), nil
}

// ReleaseTarget is the resolver for the releaseTarget field.
func (r *mutationResolver) ReleaseTarget(ctx context.Context, sessionGUID string, id int) (*model.Target, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	target, err := session.Target(session3.TargetId(id))
	if err != nil {
		return nil, err
	}

	target.SetOwner(nil)

	return TargetToGraphQL(target), nil
}

// AcquireWeapon is the resolver for the acquireWeapon field.
func (r *mutationResolver) AcquireWeapon(ctx context.Context, sessionGUID string, id int) (*model.Weapon, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	user, err := session.User(clientUuid)
	if err != nil {
		return nil, ErrUserNotInSession
	}

	weapon, err := session.Weapon(session3.WeaponId(id))
	if err != nil {
		return nil, err
	}

	weapon.SetOwner(user)

	return WeaponToGraphQL(weapon), nil
}

// ReleaseWeapon is the resolver for the releaseWeapon field.
func (r *mutationResolver) ReleaseWeapon(ctx context.Context, sessionGUID string, id int) (*model.Weapon, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	weapon, err := session.Weapon(session3.WeaponId(id))
	if err != nil {
		return nil, err
	}

	weapon.SetOwner(nil)

	return WeaponToGraphQL(weapon), nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, sessionGUID string) ([]*model.User, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	return slice.Map(session.Users(), UserToGraphQL), nil
}

// Targets is the resolver for the targets field.
func (r *queryResolver) Targets(ctx context.Context, sessionGUID string) ([]*model.Target, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	return slice.Map(session.Targets(), TargetToGraphQL), nil
}

// Weapons is the resolver for the weapons field.
func (r *queryResolver) Weapons(ctx context.Context, sessionGUID string) ([]*model.Weapon, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	session, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := session.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	return slice.Map(session.Weapons(), WeaponToGraphQL), nil
}

// SessionUpdates is the resolver for the sessionUpdates field.
func (r *subscriptionResolver) SessionUpdates(ctx context.Context, sessionGUID string) (<-chan *model.SessionUpdate, error) {
	clientUuid, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, auth2.ErrNotAuthenticated
	}

	sessionUuid, err := uuid.Parse(sessionGUID)
	if err != nil {
		return nil, err
	}

	s, err := r.SessionStorage.Get(sessionUuid)
	if err != nil {
		return nil, err
	}

	if _, err := s.User(clientUuid); err != nil {
		return nil, ErrUserNotInSession
	}

	ch := make(chan *model.SessionUpdate, 4)

	go func() {
		sub := s.Subscribe()

		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				close(ch)
				return
			case change := <-sub.Chan():
				ch <- SessionChangeToGraphQL(&change)
			}
		}
	}()

	return ch, nil
}

// X is the resolver for the x field.
func (r *vector3Resolver) X(ctx context.Context, obj *math.Vector3) (float64, error) {
	if obj == nil {
		return 0, nil
	}

	return float64(obj.X), nil
}

// Y is the resolver for the y field.
func (r *vector3Resolver) Y(ctx context.Context, obj *math.Vector3) (float64, error) {
	if obj == nil {
		return 0, nil
	}

	return float64(obj.Y), nil
}

// Z is the resolver for the z field.
func (r *vector3Resolver) Z(ctx context.Context, obj *math.Vector3) (float64, error) {
	if obj == nil {
		return 0, nil
	}

	return float64(obj.Z), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

// Vector3 returns generated.Vector3Resolver implementation.
func (r *Resolver) Vector3() generated.Vector3Resolver { return &vector3Resolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type vector3Resolver struct{ *Resolver }
