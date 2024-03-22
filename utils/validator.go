package utils

import "regexp"

func IsEmail(str string) bool {
	// Simple regex for checking email; you might want to use a more robust version
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(str)
}
