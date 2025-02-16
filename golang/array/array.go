// Helper package to handle arrays.
package array

func Map[T1, T2 any](array []T1, mapper func(T1) T2) []T2 {
	result := make([]T2, len(array))
	for i, elem := range array {
		result[i] = mapper(elem)
	}
	return result
}

func Filter[T any](array []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(array))
	for _, elem := range array {
		if predicate(elem) {
			result = append(result, elem)
		}
	}
	// Adjust the capacity to the length of the result.
	return result[:len(result):len(result)]
}

func Contains[T comparable](array []T, elem T) bool {
	for _, e := range array {
		if e == elem {
			return true
		}
	}
	return false
}
