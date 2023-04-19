package util

import (
	"archive/zip"
	"io"
	"log"
	"math/rand"
	"os"
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

// Create Zip file from a file path
func CreateZipFile(filePath string) error {

	zipfile, err := os.Create("stack.zip")
	if err != nil {
		return err
	}

	defer zipfile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	// Open the file to be zipped
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add the file to the zip archive
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}
