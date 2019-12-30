package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	// LowerLetters is the list of lowercase letters.
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters is the list of uppercase letters.
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Digits is the list of permitted digits.
	Digits = "0123456789"

	// KeyLength is the length of the Key
	KeyLength = 10
)

//RandomKey returns a ramdom alphanumeric string
func RandomKey() string {

	rand.Seed(time.Now().UnixNano())

	chars := []rune(LowerLetters +
		UpperLetters +
		Digits)

	var b strings.Builder
	for i := 0; i < KeyLength; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()

	return str
}
