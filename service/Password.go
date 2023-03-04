package service

import "golang.org/x/crypto/bcrypt"

func GetHashedPasswordOrEmpty(password string) ([]byte, error) {
	if len(password) == 0 {
		return []byte{}, nil
	}
	return HashPassword(password)
}

func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return bytes, err
}

func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
