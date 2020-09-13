package pg

import (
	"math/rand"

	"github.com/itimofeev/shaolink/internal/util"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	util.CheckOK(length > 0, "length should be positive")
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = alpha[rand.Intn(len(alpha))] // nolint:gosec // ok for demo example
	}

	return string(b)
}
