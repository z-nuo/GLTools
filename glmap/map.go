package glmap

// Keys 返回 map 中所有 key。
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回 map 中所有 value。
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// HasKey 判断 map 中是否存在指定 key。
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// Merge 合并两个 map，right 中相同 key 的值会覆盖 left。
func Merge[K comparable, V any](left map[K]V, right map[K]V) map[K]V {
	out := make(map[K]V, len(left)+len(right))
	for k, v := range left {
		out[k] = v
	}
	for k, v := range right {
		out[k] = v
	}
	return out
}
