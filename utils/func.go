package utils

import (
	"fmt"
	"log"
)

// SecondsToMinutes convert seconds to minutes
func SecondsToMinutes(sec int64) string {
	seconds := sec % 60
	minutes := sec / 60
	str := fmt.Sprintf("%02d:%02d", minutes, seconds)
	return str
}

// IsEmpty Check string is empty, return "Unknown" if empty
func IsEmpty(text string) string {
	if text == "" {
		return "Unknown"
	}
	return text
}

// ContainsInt check value is contains []int
func ContainsInt(s []int, v int) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsInt64 check value is contains []int32
func ContainsInt64(s []int64, v int64) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsString check value is contains []string
func ContainsString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// FilterString make new []int32 without specific value
func FilterString(s []string, cb func(s string) bool) []string {
	results := []string{}

	for _, i := range s {
		result := cb(i)

		if result {
			results = append(results, i)
		}
	}

	return results
}

// CheckPortIsValid check port is small than 1024 or bigger than 65535
func CheckPortIsValid(method string, port int) {
	if port < 1024 || port > 65535 {
		log.Fatal(method+" port range: 1024-65535, but current port is ", port)
	}
}
