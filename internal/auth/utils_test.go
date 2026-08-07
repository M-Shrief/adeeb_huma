package auth_test

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"slices"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func TestAuthUtils(t *testing.T) {
	t.Run("Testing CreatePermission()", func(t *testing.T) {
		want := "DBA:read"
		got := auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read)
		if got != want {
			t.Errorf("failed to create the permission in the specified Convention")
		}

		want2 := "DBA:write"
		got2 := auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Write)
		if got2 != want2 {
			t.Errorf("failed to create the permission in the specified Convention")
		}
	})

	t.Run("Testing CreatePermissions()", func(t *testing.T) {
		roles1 := []database.RoleEnum{database.RoleEnum_Normal, database.RoleEnum_DBA}
		permissions1 := auth.CreatePermissions(roles1)

		perm1 := auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read)
		perm2 := auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write)
		perm3 := auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read)
		perm4 := auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Write)
		perm5 := auth.CreatePermission(database.RoleEnum_Banned, auth.OPEnum_Write)

		does_exist := slices.Contains(permissions1, perm1)
		if does_exist != true {
			t.Errorf("Permissions should have: %s", perm1)
		}
		does_exist = slices.Contains(permissions1, perm2)
		if does_exist != true {
			t.Errorf("Permissions should have: %s", perm2)
		}
		does_exist = slices.Contains(permissions1, perm3)
		if does_exist != true {
			t.Errorf("Permissions should have: %s", perm3)
		}
		does_exist = slices.Contains(permissions1, perm4)
		if does_exist != true {
			t.Errorf("Permissions should have: %s", perm4)
		}
		does_exist = slices.Contains(permissions1, perm5)
		if does_exist != false {
			t.Errorf("Permissions shouldn't have: %s", perm5)
		}
	})

	t.Run("Testing CheckPermissions()", func(t *testing.T) {
		t.Run("Approved normal users", func(t *testing.T) {
			normal_user_permissions := []string{
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
			}
			authorize_normal_users := []string{
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
			}
			is_authorized := auth.CheckPermissions(authorize_normal_users, normal_user_permissions, auth.OPEnum_Read)
			if is_authorized != true {
				t.Errorf("couldn't correctly approve normal user permissions")
			}
		})

		t.Run("Approved adminstrators", func(t *testing.T) {
			authorize_adminstrators := []string{
				auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Analytics, auth.OPEnum_Read),
			}

			admin_permissions1 := []string{
				auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Write),
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
			}
			is_authorized := auth.CheckPermissions(authorize_adminstrators, admin_permissions1, auth.OPEnum_Read)
			if is_authorized != true {
				t.Errorf("couldn't correctly approve adminstrator permissions")
			}

			admin_permissions2 := []string{
				auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Write),
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
			}
			is_authorized = auth.CheckPermissions(authorize_adminstrators, admin_permissions2, auth.OPEnum_Read)
			if is_authorized != true {
				t.Errorf("couldn't correctly approve adminstrator permissions")
			}
		})

		t.Run("Denied banned users", func(t *testing.T) {
			banned_user1 := []string{
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Banned, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Banned, auth.OPEnum_Write),
			}
			authorize_normal_users := []string{
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
			}
			is_authorized := auth.CheckPermissions(authorize_normal_users, banned_user1, auth.OPEnum_Read)
			if is_authorized != false {
				t.Errorf("couldn't correctly deny banned user")
			}

			banned_user2 := []string{
				auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Banned, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Banned, auth.OPEnum_Write),
			}
			authorize_adminstrators := []string{
				auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
				auth.CreatePermission(database.RoleEnum_Analytics, auth.OPEnum_Read),
			}
			is_authorized = auth.CheckPermissions(authorize_adminstrators, banned_user2, auth.OPEnum_Read)
			if is_authorized != false {
				t.Errorf("couldn't correctly deny banned user")
			}
		})
	})

	t.Run("Testing CheckAdminstration()", func(t *testing.T) {
		admin_permissions1 := []string{
			auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
			auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Write),
			auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
			auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
		}
		is_adminstrator := auth.CheckAdminstration(admin_permissions1, auth.OPEnum_Write)
		if is_adminstrator != true {
			t.Errorf("couldn't correctly approve adminstrator permissions")
		}

		admin_permissions2 := []string{
			auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Read),
			auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Write),
			auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
			auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
		}
		is_adminstrator = auth.CheckAdminstration(admin_permissions2, auth.OPEnum_Write)
		if is_adminstrator != true {
			t.Errorf("couldn't correctly approve adminstrator permissions")
		}

		admin_permissions3 := []string{
			auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
			auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
		}
		is_adminstrator = auth.CheckAdminstration(admin_permissions3, auth.OPEnum_Write)
		if is_adminstrator != false {
			t.Errorf("couldn't correctly deny non-adminstrator")
		}

		admin_permissions4 := []string{
			auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Read),
			auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Write),
			auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
			auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
			auth.CreatePermission(database.RoleEnum_Banned, auth.OPEnum_Read),
			auth.CreatePermission(database.RoleEnum_Banned, auth.OPEnum_Write),
		}
		is_adminstrator = auth.CheckAdminstration(admin_permissions4, auth.OPEnum_Write)
		if is_adminstrator != false {
			t.Errorf("couldn't correctly deny banned adminstrator")
		}

	})

	t.Run("Testing CheckOwnership()", func(t *testing.T) {
		right_uuid, _ := uuid.NewUUID()
		wrong_uuid, _ := uuid.NewUUID()
		var nil_uuid *uuid.UUID

		owner := make(map[string]interface{})
		owner["id"] = right_uuid.String()
		owner["username"] = "username1"
		owner["roles"] = []database.RoleEnum{database.RoleEnum_Normal}

		claims1 := make(jwt.MapClaims)
		claims1["user"] = owner

		is_owner := auth.CheckOwnership(&right_uuid, claims1)
		if is_owner != true {
			t.Errorf("Couldn't approve owner")
		}

		non_owner := make(map[string]interface{})
		non_owner["id"] = wrong_uuid.String()
		non_owner["username"] = "username1"
		non_owner["roles"] = []database.RoleEnum{database.RoleEnum_Normal}

		claims2 := make(jwt.MapClaims)
		claims2["user"] = non_owner
		is_owner = auth.CheckOwnership(&right_uuid, claims2)
		if is_owner != false {
			t.Errorf("Couldn't reject non-owner")
		}

		is_owner = auth.CheckOwnership(nil_uuid, claims1)
		if is_owner != false {
			t.Errorf("Couldn't reject user when there's no owner_id")
		}

	})

}
