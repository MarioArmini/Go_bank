package utils

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

func GetRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(uint64(time.Now().UnixNano()))

	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func GetRandomInt() int64 {
	return rand.Int63n(99999) + 1000
}

func GetRandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}

func GetRandomInterestRate(min, max float64) string {
	interestRate := min + rand.Float64()*(max-min)

	return fmt.Sprintf("%.2f", interestRate)
}
