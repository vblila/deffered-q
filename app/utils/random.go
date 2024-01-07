package utils

import (
	"crypto/rand"
	"math"
)

const randomIdLength = 10
const randomIdAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func RandomId() string {
	buf := make([]byte, randomIdLength)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	k := 255.0 / float64(len(randomIdAlphabet)-1)

	randomId := make([]byte, randomIdLength)

	for i := 0; i < randomIdLength; i++ {
		key := int(math.Round(float64(buf[i]) / k))
		randomId[i] = randomIdAlphabet[key]
	}

	return string(randomId)
}
