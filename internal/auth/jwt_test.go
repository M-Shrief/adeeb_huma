package auth_test

import (
	"adeeb_huma/config"
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"adeeb_huma/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Custom function to make paths start from the project root
// so we can read .env file, and load config
func Init() {
	_, filename, _, _ := runtime.Caller(0)
	// Adjust the number of "../" to reach your project root
	dir := filepath.Join(filepath.Dir(filename), "..", "..")
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
}

func TestAuthJWT(t *testing.T) {
	Init()
	config.LoadENV()

	t.Run("Testing CreateJWT()", func(t *testing.T) {
		id, _ := uuid.NewUUID()
		roles1 := []database.RoleEnum{database.RoleEnum_DBA, database.RoleEnum_Normal}

		token, err := auth.CreateJWT(
			time.Hour*2,
			auth.JWTUserClaim{
				ID:       id.String(),
				Username: "username1",
				Roles:    roles1,
			},
			auth.CreatePermissions(roles1),
		)

		if err != nil {
			t.Errorf("Failed to create JWT token")
		}

		parsed_token, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
			}
			return config.JWT_PUBLIC_KEY, nil
		})
		if err != nil {
			t.Errorf("Failed to parse JWT token")
		}

		claims, ok := parsed_token.Claims.(jwt.MapClaims)
		cmp := time.Now().Unix()
		if !ok || !parsed_token.Valid {
			t.Errorf("failed to create valid JWT token")
		} else if claims.VerifyExpiresAt(cmp, true) == false {
			t.Errorf("failed to create non-expired JWT token")
		} else if claims.VerifyIssuedAt(cmp, true) == false {
			t.Errorf("failed to issue JWT token in a valid time")
		}

		claims_user, ok := claims["user"].(map[string]interface{})
		if !ok {
			t.Errorf("Failed to create JWT token")
		}

		claims_id, ok := claims_user["id"].(string)
		if !ok {
			t.Errorf("Failed to create JWT token")
		}
		if claims_id != id.String() {
			t.Errorf("failed to get the right ID from JWT's claims")
		}

		claims_username, ok := claims_user["username"].(string)
		if !ok {
			t.Errorf("Failed to create JWT token")
		}
		if claims_username != "username1" {
			t.Errorf("failed to the right username from JWT's claims")
		}

		claims_roles, err := utils.InterfaceToStringSlice(claims_user["roles"])
		if len(claims_roles) != len(roles1) {
			t.Errorf("couldn't extract user.roles from JWT's claim")
		}
		if claims_roles[0] != string(roles1[0]) {
			t.Errorf("couldn't get the right role from user.roles in JWT's claims")
		}

		// Checking if claims' permission is a valid []string
		claim_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
		if err != nil {
			t.Errorf("couldn't extract permissions from JWT's claim as []string")
		}

		if slices.Contains(claim_permissions, auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read)) != true {
			t.Errorf("claims' permission doesn't contain DBA:read")
		}
		if slices.Contains(claim_permissions, auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Write)) != true {
			t.Errorf("claims' permission doesn't contain DBA:Write")
		}
		if slices.Contains(claim_permissions, auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read)) != true {
			t.Errorf("claims' permission doesn't contain Normal:read")
		}
		if slices.Contains(claim_permissions, auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write)) != true {
			t.Errorf("claims' permission doesn't contain Normal:Write")
		}
	})

	t.Run("Testing VerifyJWT()", func(t *testing.T) {
		hour_ago := time.Now().Add(time.Hour * -1).UTC()
		valid_exp := hour_ago.Add(time.Hour * 2).Unix()
		invalid_exp := hour_ago.Add(time.Minute * 30).Unix()
		valid_iat := hour_ago.Unix()
		invalid_iat := time.Now().Add(time.Hour * 1).Unix()

		claims1 := make(jwt.MapClaims)
		claims1["exp"] = valid_exp
		claims1["iat"] = valid_iat
		token1, _ := jwt.NewWithClaims(jwt.SigningMethodRS256, claims1).SignedString(config.JWT_PRIVATE_KEY)

		_, err := auth.VerifyJWT("Bearer " + token1)
		if err != nil {
			t.Errorf("Couldn't verify valid JWT")
		}

		claims2 := make(jwt.MapClaims)
		claims2["exp"] = invalid_exp
		claims2["iat"] = valid_iat
		token2, _ := jwt.NewWithClaims(jwt.SigningMethodRS256, claims2).SignedString(config.JWT_PRIVATE_KEY)
		_, err = auth.VerifyJWT("Bearer " + token2)
		if err == nil {
			t.Errorf("Couldn't verify invalid JWT, it shou'd have error because claims.exp is invalid: %v", err)
		}

		claims3 := make(jwt.MapClaims)
		claims3["exp"] = valid_exp
		claims3["iat"] = invalid_iat
		token3, _ := jwt.NewWithClaims(jwt.SigningMethodRS256, claims3).SignedString(config.JWT_PRIVATE_KEY)
		_, err = auth.VerifyJWT("Bearer " + token3)
		if err == nil {
			t.Errorf("Couldn't verify invalid JWT, it shou'd have error because claims.iat is invalid: %v", err)
		}

	})
}
