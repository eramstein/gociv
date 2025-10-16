package utils

import (
	"math/rand/v2"
)

type Tile struct {
	RegionId int
	Col      int
	Row      int
}

type VoronoiSeed struct{ C, R int }

func GetVoronoiSeeds(width, height, numRegions int) []VoronoiSeed {
	tiles := make([]Tile, width*height)
	for i := 0; i < width*height; i++ {
		tiles[i] = Tile{
			Col: i % width,
			Row: i / width,
		}
	}

	seeds := make([]VoronoiSeed, 0, numRegions)
	used := make(map[int]struct{})
	for len(seeds) < numRegions {
		c := int(rand.IntN(width))
		r := int(rand.IntN(height))
		key := r*width + c
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		seeds = append(seeds, VoronoiSeed{C: c, R: r})
	}

	// Lloyd relaxation: alternate assignment and centroid update
	const iterations = 2
	for it := 0; it < iterations; it++ {
		// Assign tiles to nearest seed
		for idx := range tiles {
			t := &tiles[idx]
			bestRegion := 0
			bestDist := int(^uint(0) >> 1)
			for i, s := range seeds {
				d := GetHexDistance(t.Col, t.Row, s.C, s.R)
				if d < bestDist {
					bestDist = d
					bestRegion = i
				}
			}
			t.RegionId = bestRegion
		}

		// Recompute seeds at region centroids using cube coords with proper rounding
		sumX := make([]float64, numRegions)
		sumY := make([]float64, numRegions)
		sumZ := make([]float64, numRegions)
		count := make([]int, numRegions)
		for i := range tiles {
			t := &tiles[i]
			rx, ry, rz := OddRToCube(t.Col, t.Row)
			rid := t.RegionId
			sumX[rid] += float64(rx)
			sumY[rid] += float64(ry)
			sumZ[rid] += float64(rz)
			count[rid]++
		}

		for i := 0; i < numRegions; i++ {
			if count[i] == 0 {
				// If empty, keep previous seed
				continue
			}
			avgX := sumX[i] / float64(count[i])
			avgY := sumY[i] / float64(count[i])
			avgZ := sumZ[i] / float64(count[i])
			cx, cy, cz := CubeRound(avgX, avgY, avgZ)
			col, row := CubeToOddR(cx, cy, cz)
			// clamp to map bounds
			if col < 0 {
				col = 0
			}
			if row < 0 {
				row = 0
			}
			if col >= width {
				col = width - 1
			}
			if row >= height {
				row = height - 1
			}
			seeds[i] = VoronoiSeed{C: col, R: row}
		}
	}

	return seeds
}
