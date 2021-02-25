package utils

import "os"

// MustGetHomeDir gets the user home directory
// Panic if an error occurs
func MustGetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return homeDir
}
