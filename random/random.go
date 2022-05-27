package random

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Generate(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(letters))
		result[i] = letters[idx]
	}

	return string(result)
}
