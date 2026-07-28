package utils

func EnsureSliceItemsAreUnique[T comparable](input []T) []T {
	seen := make(map[T]struct{}, len(input))
	result := make([]T, 0, len(input))

	for _, v := range input {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
