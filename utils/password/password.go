package password

import (
	"github.com/DataWorkbench/common/constants"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func Encode(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Check(password string, encodePassword string) bool {
	ret := bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(password))
	return ret == nil
}

func RandomGenerateAccessKey() (string, string) {
	rand.Seed(time.Now().UnixNano())
	accessKeyId := make([]byte, constants.AccessKeyIdLength)
	secretKey := make([]byte, constants.SecretKeyLength)
	for i := 0; i < constants.AccessKeyIdLength; i++ {
		accessKeyId[i] = constants.AccessKeyIdLetters[rand.Intn(len(constants.AccessKeyIdLetters))]
	}
	for i := 0; i < constants.SecretKeyLength; i++ {
		secretKey[i] = constants.SecretKeyLetters[rand.Intn(len(constants.SecretKeyLetters))]
	}
	return string(accessKeyId), string(secretKey)
}

func GenerateSession() string {
	rand.Seed(time.Now().UnixNano())
	session := make([]byte, constants.SessionLength)
	for i := 0; i < constants.SessionLength; i++ {
		session[i] = constants.SessionLetters[rand.Intn(len(constants.SessionLetters))]
	}
	return string(session)
}
