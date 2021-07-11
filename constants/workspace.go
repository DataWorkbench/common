package constants

// Workspace status.
const (
	SpaceStatusEnabled  int32 = iota + 1 // => "enabled"
	SpaceStatusDisabled                  // => "disabled"
)

// Workspace System roles.
const (
	RoleTypeSystem int32 = iota + 1
	RoleTypeCustom
)

// Workspace system role id.
const (
	RoleIdSpaceOwner     = IdPrefixRoleSystem + "1000000000000001"
	RoleIdSpaceAdmin     = IdPrefixRoleSystem + "1000000000000002"
	RoleIdSpaceDeveloper = IdPrefixRoleSystem + "1000000000000003"
	RoleIdSpaceOperator  = IdPrefixRoleSystem + "1000000000000004"
	RoleIdSpaceVisitor   = IdPrefixRoleSystem + "1000000000000005"
)

// Operation type.
const (
	OpTypeUnknown int32 = iota + 1
	OpTypeRead
	OpTypeWrite
	OpTypeDelete
)
