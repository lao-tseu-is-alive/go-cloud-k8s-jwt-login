package crypto

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	emptyStringSha256            = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	emptyStringMd5               = "d41d8cd98f00b204e9800998ecf8427e"
	testStringSha256             = "4acf0b39d9c4766709a3689f553ac01ab550545ffa4544dfc0b2cea82fba02a3"
	testStringHashSaltFromSha256 = "$2a$04$Clft9Wh3JQqdPyjKDOk3Bu.N4vHEoeSZ8mec7P2sHpivKNx78CyZm"
)

// HashAndSalt calculates the salted hash for the given password
func HashAndSalt(plainPwd string) (string, error) {
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower than the MinCost (4)
	bytePassword := []byte(plainPwd)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("error doing bcrypt.GenerateFromPassword : %v", err)
	}
	return string(hash), nil
}

// ComparePasswords verify if a given password (plainPwd) corresponds to the hashedPwd from User Database
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	// golog.Un(golog.Trace("ComparePasswords hash [%v] pwd [%v]", hashedPwd, plainPwd))
	// we need to convert strings to byte slice
	byteHash := []byte(hashedPwd)
	bytePassword := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)
	if err != nil {
		return false
	}
	return true
}

// Sha256Hash calculates the sha256 hash of a string
func Sha256Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ValidatePasswordHash verifies that the given string is a valid password hash
func ValidatePasswordHash(s string) bool {
	if len(s) < 32 { // 32 is for md5 so at least this
		return false
	}
	if s == emptyStringMd5 {
		return false
	}
	if s == emptyStringSha256 {
		return false
	}
	return true
}
