package auth

import "golang.org/x/crypto/bcrypt"

// create hash password
func HashPassword(password string) (string, error) {
	hass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hass), nil
}

// compare hashed password
func ComparePassword(hashed string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), password)
	return err == nil
}
