package password

import (
	"fmt"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword() (string, error) {
	pass, err := password.Generate(12, 3, 2, false, false)
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}

	return pass, nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashed), nil
}

func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}

	return nil
}
