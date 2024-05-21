package qp

import "strings"

// CleanSliceString removes leading and trailing whitespace from each string in the given
// slice and returns a new slice with the cleaned strings.
func cleanSliceString(list []string) []string {
	var clean []string
	for _, v := range list {
		v = strings.Trim(v, " \t")
		if len(v) > 0 {
			clean = append(clean, v)
		}
	}
	return clean
}

// StringInSlice checks if a string is present in a slice of strings.
// It returns true if the string is found, otherwise false.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
