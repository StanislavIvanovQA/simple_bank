package testUtils

import (
	"github.com/go-faker/faker/v4"
	"log"
	"math/rand"
)

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
	currencies := []string{"EUR", "USD", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
