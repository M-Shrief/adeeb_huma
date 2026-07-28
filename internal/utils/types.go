package utils

import "fmt"

// It's used to add typing for a value with type interface{} but it's actually a []string
//
// Mainly used to get permissions []string from JWT claim
func InterfaceToStringSlice(input interface{}) ([]string, error) {
	list := input.([]interface{})
	result := make([]string, len(list))
	for i, v := range list {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a string: %T", i, v)
		}
		result[i] = s
	}
	return result, nil
}
