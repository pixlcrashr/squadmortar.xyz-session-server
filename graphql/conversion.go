package graphql

import (
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/graphql/model"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/math"
	session2 "github.com/pixlcrashr/squadmortar.xyz-sessions-server/session"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/slice"
)

func UserToGraphQL(user session2.User) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		ClientGUID: user.ClientUuid().String(),
		Name:       user.Name(),
	}
}

func SessionToGraphQL(session session2.Session) *model.Session {
	if session == nil {
		return nil
	}

	return &model.Session{
		GUID:    session.Uuid().String(),
		Targets: slice.Map(session.Targets(), TargetToGraphQL),
		Users:   slice.Map(session.Users(), UserToGraphQL),
		Weapons: slice.Map(session.Weapons(), WeaponToGraphQL),
	}
}

func WeaponTypeToGraphQL(weaponType session2.WeaponType) model.WeaponType {
	switch weaponType {
	case session2.TechnicalMortarWeaponType:
		return model.WeaponTypeTechnicalMortar
	case session2.HellCannonWeaponType:
		return model.WeaponTypeHellCanon
	case session2.RocketsWeaponType:
		return model.WeaponTypeRockets
	case session2.StandardMortarWeaponType:
		return model.WeaponTypeStandardMortar
	default:
		return model.WeaponTypeStandardMortar
	}
}

func WeaponToGraphQL(weapon session2.Weapon) *model.Weapon {
	if weapon == nil {
		return nil
	}

	position := weapon.Position()

	return &model.Weapon{
		ID:       int(weapon.Id()),
		Active:   weapon.Active(),
		IsOwned:  weapon.IsOwned(),
		Owner:    UserToGraphQL(weapon.Owner()),
		Position: &position,
		Type:     WeaponTypeToGraphQL(weapon.Type()),
	}
}

func TargetToGraphQL(target session2.Target) *model.Target {
	if target == nil {
		return nil
	}

	position := target.Position()

	return &model.Target{
		ID:       int(target.Id()),
		Active:   target.Active(),
		IsOwned:  target.IsOwned(),
		Owner:    UserToGraphQL(target.Owner()),
		Position: &position,
	}
}

func SessionChangeToGraphQL(sessionChange *session2.SessionChange) *model.SessionUpdate {
	return &model.SessionUpdate{
		UserJoined:    UserToGraphQL(sessionChange.UserJoined),
		UserLeft:      UserToGraphQL(sessionChange.UserLeft),
		UserChanged:   UserToGraphQL(sessionChange.UserChanged),
		TargetAdded:   TargetToGraphQL(sessionChange.TargetAdded),
		TargetChanged: TargetToGraphQL(sessionChange.TargetChanged),
		TargetRemoved: TargetToGraphQL(sessionChange.TargetRemoved),
		WeaponAdded:   WeaponToGraphQL(sessionChange.WeaponAdded),
		WeaponChanged: WeaponToGraphQL(sessionChange.WeaponChanged),
		WeaponRemoved: WeaponToGraphQL(sessionChange.WeaponRemoved),
	}
}

func WeaponTypeFromGraphQL(weaponType model.WeaponType) session2.WeaponType {
	switch weaponType {
	case model.WeaponTypeTechnicalMortar:
		return session2.TechnicalMortarWeaponType
	case model.WeaponTypeHellCanon:
		return session2.HellCannonWeaponType
	case model.WeaponTypeRockets:
		return session2.RocketsWeaponType
	case model.WeaponTypeStandardMortar:
		return session2.StandardMortarWeaponType
	default:
		return session2.StandardMortarWeaponType
	}
}

func Vector3InputFromGraphQL(vector3Input model.Vector3Input) math.Vector3 {
	return math.Vector3{
		X: float32(vector3Input.X),
		Y: float32(vector3Input.Y),
		Z: float32(vector3Input.Z),
	}
}
