package utils

import "strings"

func ErrsFlatten(errs map[string]map[string]string) map[string][]string {
	errsFlatten := make(map[string][]string)
	for field, rules := range errs {
		for _, msg := range rules {
			errsFlatten[strings.ToLower(field)] = append(errsFlatten[strings.ToLower(field)], msg)
		}
	}
	return errsFlatten
}
