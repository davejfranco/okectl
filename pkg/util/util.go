package util

import (
	"log"
	"math/rand"
	"net/netip"
	"regexp"
	"strings"
)

func FatalIfError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// RandomString returns a random string of length n
func RandomString(length int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// HexaMask returns CIDR mask in format /16 from ffff0000 format
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

func OCIDvalidator(ocid string) bool {
	//regex := "ocid1.(tenancy|vcn|intance|privateip).oc1(..|phx|iad|)[a-zA-Z0-9]*"
	regex, _ := regexp.Compile("ocid1.(tenancy|compartment|vcn|intance|vnic|oke).oc1.(..|phx|iad|).[a-zA-Z0-9]*$")
	return regex.Match([]byte(ocid))
}

func IsvalidCIDR(cidr string) bool {
	//validate CIDR
	_, err := netip.ParsePrefix(cidr)
	if err != nil {
		return false
	}
	return true
}
