package constants

import (
	"fmt"

	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
)

var (
	_ = SystemRoleLists
	_ = SystemRoleMap
)

// Defines system role info. not be modified it.
var (
	SystemRoleSpaceAdmin = &pbmodel.SystemRole{
		Id:   fmt.Sprintf("%s%016x", IdPrefixRoleSystem, pbmodel.SystemRole_SpaceAdmin.Number()),
		Type: pbmodel.SystemRole_SpaceAdmin, Name: "空间管理员",
	}
	SystemRoleSpaceDeveloper = &pbmodel.SystemRole{
		Id:   fmt.Sprintf("%s%016x", IdPrefixRoleSystem, pbmodel.SystemRole_SpaceDeveloper.Number()),
		Type: pbmodel.SystemRole_SpaceDeveloper, Name: "开发",
	}
	SystemRoleSpaceOperator = &pbmodel.SystemRole{
		Id:   fmt.Sprintf("%s%016x", IdPrefixRoleSystem, pbmodel.SystemRole_SpaceOperator.Number()),
		Type: pbmodel.SystemRole_SpaceOperator, Name: "运维",
	}
	SystemRoleSpaceVisitor = &pbmodel.SystemRole{
		Id:   fmt.Sprintf("%s%016x", IdPrefixRoleSystem, pbmodel.SystemRole_SpaceVisitor.Number()),
		Type: pbmodel.SystemRole_SpaceVisitor, Name: "访客",
	}
)

// SystemRoleLists store all system role in a list.
var SystemRoleLists = []*pbmodel.SystemRole{
	SystemRoleSpaceAdmin,
	SystemRoleSpaceDeveloper,
	SystemRoleSpaceOperator,
	SystemRoleSpaceVisitor,
}

// SystemRoleMap store all system role in a map. key is <id>.
var SystemRoleMap = map[string]*pbmodel.SystemRole{
	SystemRoleSpaceAdmin.Id:     SystemRoleSpaceAdmin,
	SystemRoleSpaceDeveloper.Id: SystemRoleSpaceDeveloper,
	SystemRoleSpaceOperator.Id:  SystemRoleSpaceOperator,
	SystemRoleSpaceVisitor.Id:   SystemRoleSpaceVisitor,
}
