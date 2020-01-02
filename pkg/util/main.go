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

//HexaMask returns CIDR mask in format /16 from ffff0000 format
func HexaMask(hexa string) string {

	n := strings.Count(hexa, "0")

	switch {
	case n == 2:
		return string("24")
	case n == 4:
		return string("16")
	case n == 6:
		return string("8")
	case n == 0:
		return string("32")
	}
	return ""
}
