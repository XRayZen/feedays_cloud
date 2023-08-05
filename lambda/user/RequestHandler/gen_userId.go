package RequestHandler

import (
	"user/DBRepo"

	// "fmt"
	"math/rand"
)

func GenRandomUserID(repo DBRepo.DBRepo,ip string) (string, error) {
	return RandomString(15), nil
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
