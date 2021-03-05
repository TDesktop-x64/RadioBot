package utils

import (
	"fmt"
	"log"
)

func SecondsToMinutes(sec int64) string {
	seconds := sec % 60
	minutes := sec / 60
	str := fmt.Sprintf("%02d:%02d", minutes, seconds)
	return str
}

func IsEmpty(text string) string {
	if text == "" {
		return "Unknown"
	}
	return text
}

func ContainsInt(s []int, v int) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func ContainsInt32(s []int32, v int32) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func FilterInt32(s []int32, cb func(s int32) bool) []int32 {
	results := []int32{}

	for _, i := range s {
		result := cb(i)

		if result {
			results = append(results, i)
		}
	}

	return results
}

func CheckPortIsValid(method string, port int) {
	if port < 1024 || port > 65535 {
		log.Fatal(method+" port range: 1024-65535, but current port is ", port)
	}
}
