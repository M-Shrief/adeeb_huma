package auth

import (
	"adeeb_huma/database"
	"slices"
)

type OPEnum string

const (
	OPEnum_Write OPEnum = "write"
	OPEnum_Read  OPEnum = "read"
)

func CreatePermissions(roles []database.RoleEnum) []string {
	var permissions []string

	for _, role := range roles {
		permissions = append(permissions, CreatePermission(role, OPEnum_Read))
		permissions = append(permissions, CreatePermission(role, OPEnum_Write))
	}

	return permissions
}
func CreatePermission(role database.RoleEnum, op OPEnum) string {
	return string(role) + ":" + string(op)
}

func CheckPermissions(authorized_list, user_permissions []string, op OPEnum) bool {
	isAuthorized := false
	isBanned := false

	for _, perm := range user_permissions {
		if op == OPEnum_Write && perm == CreatePermission(database.RoleEnum_Banned, OPEnum_Write) {
			isBanned = true
			break
		} else if op == OPEnum_Read && perm == CreatePermission(database.RoleEnum_Banned, OPEnum_Read) {
			isBanned = true
			break
		}

		it_contains_perm := slices.Contains(authorized_list, perm)
		if it_contains_perm {
			isAuthorized = true
		}
	}

	if isBanned {
		isAuthorized = false
	}

	return isAuthorized
}
