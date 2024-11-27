package goconv

func Map[E any, R any](items []E, f func(E) R) []R {
	result := make([]R, len(items))
	for i, v := range items {
		result[i] = f(v)
	}
	return result
}

func Delkey[M ~map[K]V, K comparable, V any](m M, keys ...K) M {
	dst := make(M)
	keymap := map[K]any{}

	for _, x := range keys {
		keymap[x] = struct{}{}
	}

	for k, v := range m {
		if _, ok := keymap[k]; !ok {
			dst[k] = v
			delete(keymap, k)
		}
	}

	return dst
}

func Assign[K comparable, V any, Map ~map[K]V](maps ...Map) Map {
	count := 0
	for i := range maps {
		count += len(maps[i])
	}

	out := make(Map, count)
	for i := range maps {
		for k := range maps[i] {
			out[k] = maps[i][k]
		}
	}

	return out
}
