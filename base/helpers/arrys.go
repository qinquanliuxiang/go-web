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

type Named interface {
	GetName() string
}

// FindMissingByName 返回 target 中不在 existing.Name 集合中的 name
func FindMissingByName[T Named](existing []T, target []string) []string {
	nameSet := make(map[string]struct{}, len(existing))
	for _, item := range existing {
		nameSet[item.GetName()] = struct{}{}
	}

	var missing []string
	for _, name := range target {
		if _, ok := nameSet[name]; !ok {
			missing = append(missing, name)
		}
	}
	return missing
}

type IDed interface {
	GetID() int
}

// FindMissingByID 返回 target 中不在 existing.ID 集合中的 id
func FindMissingByID[T IDed](existing []T, target []int) []int {
	idSet := make(map[int]struct{}, len(existing))
	for _, item := range existing {
		idSet[item.GetID()] = struct{}{}
	}

	var missing []int
	for _, id := range target {
		if _, ok := idSet[id]; !ok {
			missing = append(missing, id)
		}
	}
	return missing
}
