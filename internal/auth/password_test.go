package auth_test

import (
	"adeeb_huma/internal/auth"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthPassword(t *testing.T) {
	t.Run("Testing HashPassword()", func(t *testing.T) {
		var password = "password"

		hash1, err := auth.HashPassword(password)
		if err != nil {
			t.Errorf("Failed to hash password")
		}
		err = bcrypt.CompareHashAndPassword([]byte(hash1), []byte(password))
		if err != nil {
			t.Errorf("Failed to compare the correct password with the hash")
		}

		hash2, err := auth.HashPassword(password)
		if err != nil {
			t.Errorf("Failed to hash password")
		}
		err = bcrypt.CompareHashAndPassword([]byte(hash2), []byte(password))
		if err != nil {
			t.Errorf("Failed to compare the correct password with the hash")
		}
	})

	t.Run("Testing VerifyPassword()", func(t *testing.T) {
		var password = "password"
		res1, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		hash1 := string(res1)

		err := auth.VerifyPassword(password, hash1)
		if err != nil {
			t.Errorf("Failed to verify correct password")
		}

		res2, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		hash2 := string(res2)
		err = auth.VerifyPassword(password, hash2)
		if err != nil {
			t.Errorf("Failed to verify correct password")
		}

		var wrong_pass = "wrong one"
		res3, _ := bcrypt.GenerateFromPassword([]byte(wrong_pass), bcrypt.DefaultCost)
		hash3 := string(res3)
		err = auth.VerifyPassword(wrong_pass, hash3)
		if err != nil {
			t.Errorf("Failed to verify wrong password")
		}

	})
}
