package helpers

// Deduplicate 去重
func Deduplicate[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	var result []T

	for _, v := range slice {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// GetUnique 获取第二个切片中不存在于第一个切片的元素
//
// s2 相对于 s1 的差集
//
// s2 - s1 或 s2 \ s1
func GetUnique[T comparable](s1, s2 []T) []T {
	m := make(map[T]struct{})
	for _, v := range s1 {
		m[v] = struct{}{}
	}

	var result []T
	for _, v := range s2 {
		if _, found := m[v]; !found {
			result = append(result, v)
		}
	}
	return result
}
