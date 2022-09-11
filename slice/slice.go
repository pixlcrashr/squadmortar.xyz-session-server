package slice

func MapValuesToSlice[K comparable, T any](m map[K]T) []T {
	vs := make([]T, len(m))
	i := 0
	for _, v := range m {
		vs[i] = v
		i++
	}

	return vs
}

type MapFunc[V any, W any] func(V) W

func Map[V any, W any](vs []V, fn MapFunc[V, W]) []W {
	ws := make([]W, len(vs))

	for i, v := range vs {
		ws[i] = fn(v)
	}

	return ws
}
