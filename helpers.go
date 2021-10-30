package main

import "strconv"

// Default string to defaultString if the string is empty
func defaultEmptyString(s *string, defaultString string) {
	if *s == "" {
		*s = defaultString
	}
}

// Convert string to int, if unsuccessful return 0
func intOrZero(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
