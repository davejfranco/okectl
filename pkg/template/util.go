package template

import (
	"fmt"
	"os"
)

func okectlDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	directory := fmt.Sprintf("%s/%s", currentDir, ".okectl")

	// Check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// Directory doesn't exist, so create it
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return directory
}
