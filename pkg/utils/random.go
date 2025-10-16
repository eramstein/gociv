package utils

import "math/rand/v2"

func GetRandomFromArray[T comparable](a []T) T {
	return a[rand.IntN(len(a))]
}
