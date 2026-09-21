package service

import (
	"api-students/app/model"
	"api-students/helper"
)

func CanAccessStudent(
	current model.AuthUser,
	ownerID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {
	if current.ID == ownerID {
		return true
	}

	return perms.Has(current.Role, anyPermission)
}
