package auth

import (
	"adeeb_huma/config"
	"adeeb_huma/database"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTUserClaim struct {
	ID       string              `json:"id"`
	Username string              `json:"username"`
	Roles    []database.RoleEnum `json:"roles"`
}

func CreateJWT(ttl time.Duration, user JWTUserClaim, permissions []string) (string, error) {
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["user"] = user               // Our custom data.
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["permissions"] = permissions
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(config.JWT_PRIVATE_KEY)
	if err != nil {
		return "", fmt.Errorf("create token error: %w", err)
	}

	return token, nil
}

func VerifyJWT(authHeader string) (jwt.MapClaims, error) {
	if authHeader == "" {
		return nil, fmt.Errorf("Empty header")
	}

	token := authHeader[7:]

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return config.JWT_PUBLIC_KEY, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	cmp := time.Now().Unix()
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	} else if claims.VerifyExpiresAt(cmp, true) == false {
		return nil, fmt.Errorf("error: token is already expired")
	} else if claims.VerifyIssuedAt(cmp, true) == false {
		return nil, fmt.Errorf("error: token is issued in invalid time")
	}

	return claims, nil
}
