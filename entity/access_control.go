package entity

//AccessControl only keeps allowed permissions
// if we want to use denied permissions : add IsDenied bool as field of AccessControl
type AccessControl struct {
	ID uint
	ActorID uint
	ActorType ActorType
	PermissionID uint
}

type ActorType string

const(
	RoleActorType= "role"
	UserActorType= "user"
)