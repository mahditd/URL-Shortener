package usecase

import (
	"crypto/rand"
	"math/big"
)

func generateCode() (string, error) {

	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, 6)

	for i := range result {

		number, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))

		if err != nil {
			return "", err
		}

		result[i] = chars[number.Int64()]
	}

	return string(result), nil
}
