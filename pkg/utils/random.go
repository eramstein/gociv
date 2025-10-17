package utils

import (
	"fmt"
	"math/rand/v2"
)

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func GetRandomFromArray[T comparable](a []T) T {
	return a[rand.IntN(len(a))]
}

// GetWeightedRandomFromMap picks a random key with probability proportional to its int weight
// weights can't be negative
func GetWeightedRandomFromMap[K comparable](m map[K]int) K {
	if len(m) == 0 {
		panic("GetWeightedRandomFromMap: empty map")
	}
	totalWeight := 0
	for _, w := range m {
		totalWeight += w
	}
	r := rand.IntN(totalWeight)
	acc := 0
	for k, w := range m {
		acc += w
		if r < acc {
			return k
		}
	}
	var zero K
	fmt.Println("GetWeightedRandomFromMap: no key found")
	return zero
}
