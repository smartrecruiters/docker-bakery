package commons

import "sort"

// Index returns the first index of the target string t, or -1 if no match is found.
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// RemoveLast removes the last element from the slice and returns the altered slice.
func RemoveLast(s []string) []string {
	return RemoveIndex(s, len(s)-1)
}

// RemoveIndex removes the element from the slice by its index.
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// Include returns true if the target string t is in the slice.
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

// Contains returns true if the target string t is in the slice.
func Contains(vs []string, t string) bool {
	return Include(vs, t)
}

// Any returns true if one of the strings in the slice satisfies the predicate f.
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// All returns true if all of the strings in the slice satisfy the predicate f.
func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// Filter returns a new slice containing all strings in the slice that satisfy the predicate f.
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Map returns a new slice containing the results of applying the function f to each string in the original slice.
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// SortKeys sorts keys of the map of interfaces and returns them in the slice.
func SortKeys(m map[string]interface{}) []string {
	i, sorted := 0, make([]string, len(m))
	for k := range m {
		sorted[i] = k
		i++
	}
	sort.Strings(sorted)
	return sorted
}

// SortMapKeys sorts keys of the map of strings and returns them in the slice.
func SortMapKeys(m map[string]string) []string {
	i, sorted := 0, make([]string, len(m))
	for k := range m {
		sorted[i] = k
		i++
	}
	sort.Strings(sorted)
	return sorted
}
