package ng

import (
	"fmt"
	"gociv/pkg/utils"
)

const (
	MapWidth    = 80
	MapHeight   = 50
	RegionCount = 200
)

func NewWorldMap() *WorldMap {
	mapData := WorldMap{
		Width:   MapWidth,
		Height:  MapHeight,
		Tiles:   [MapWidth * MapHeight]Tile{},
		Regions: [RegionCount]Region{},
	}

	initWorldMap(&mapData)

	initRegions(&mapData)

	assignVoronoiRegions(&mapData, RegionCount)

	return &mapData
}

// Start with all plains
func initWorldMap(world *WorldMap) {
	for i := 0; i < MapWidth*MapHeight; i++ {
		world.Tiles[i] = Tile{
			Col:     i % MapWidth,
			Row:     i / MapWidth,
			Terrain: Plains,
		}
	}
}

func initRegions(world *WorldMap) {
	avg := (MapWidth * MapHeight) / RegionCount
	capacity := avg + avg/4 // +25% pre-allocated headroom
	for i := 0; i < RegionCount; i++ {
		world.Regions[i] = Region{
			Id:      i,
			Name:    fmt.Sprintf("Region %d", i),
			TileIds: make([]int, 0, capacity),
			Biome:   utils.GetRandomFromArray([]TerrainType{Plains, Forest, Mountains, Water, Grassland, Ice, Savannah, Swamp}),
		}
	}
}

// partition the map into regions using Voronoi cells on the hex grid
func assignVoronoiRegions(world *WorldMap, numRegions int) {
	if numRegions <= 0 {
		return
	}

	seeds := utils.GetVoronoiSeeds(MapWidth, MapHeight, numRegions)

	for idx := range world.Tiles {
		t := &world.Tiles[idx]
		bestRegion := 0
		bestDist := int(^uint(0) >> 1)
		for i, s := range seeds {
			d := utils.GetHexDistance(t.Col, t.Row, s.C, s.R)
			if d < bestDist {
				bestDist = d
				bestRegion = i
			}
		}
		t.RegionId = bestRegion
		world.Regions[bestRegion].TileIds = append(world.Regions[bestRegion].TileIds, idx)
		t.Terrain = world.Regions[bestRegion].Biome
	}
}
