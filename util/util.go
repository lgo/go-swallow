package util

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateJobID makes a 24-char random string
func GenerateJobID() string {
	b := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
