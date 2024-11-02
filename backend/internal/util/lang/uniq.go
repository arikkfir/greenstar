package lang

func Uniq[T comparable](s []T) []T {
	seen := make(map[T]bool)
	var result []T
	for _, item := range s {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
