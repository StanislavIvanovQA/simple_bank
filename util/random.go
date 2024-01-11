package util

import (
	"github.com/go-faker/faker/v4"
	"log"
	"math/rand"
)

func RandomUserName() string {
	return faker.Username()
}

func RandomName() string {
	return faker.Name()
}

func RandomInt64(min, max int) int64 {
	randomInt, err := faker.RandomInt(min, max)
	if err != nil {
		log.Fatal(err)
	}

	return int64(randomInt[0])
}

func RandomCurrency() string {
	currencies := []string{EUR, CAD, USD}
	return currencies[rand.Intn(len(currencies))]
}

// RandomEmail Randomly created email
func RandomEmail() string {
	return faker.Email()
}

func RandomPassword() string {
	return faker.Password()
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
