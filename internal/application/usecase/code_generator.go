package usecase

import (
	"crypto/rand"
	"math/big"
)

func generateCode() string {

	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, 6)

	for i := range result {

		number, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))

		result[i] = chars[number.Int64()]
	}

	return string(result)
}
