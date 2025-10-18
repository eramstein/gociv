package ng

import (
	"fmt"
	"gociv/pkg/utils"
	"math"
	"math/rand"
	"time"
)

func NewWorldMap() *WorldMap {
	start := time.Now()
	defer func() { fmt.Printf("NewWorldMap took %s\n", time.Since(start)) }()
	mapData := WorldMap{
		Width:   MapWidth,
		Height:  MapHeight,
		Tiles:   [MapWidth * MapHeight]Tile{},
		Regions: [RegionCount]Region{},
	}

	initWorldMap(&mapData)
	initRegions(&mapData)
	assignTilesToRegions(&mapData, RegionCount)
	setRegionCentroids(&mapData)
	regionBiomes := makeRegionBiomes(&mapData)
	assignTerrainToTiles(&mapData, regionBiomes)
	assignFeaturesToRegions(&mapData, regionBiomes)
	computeRegionBorders(&mapData)
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
		}
	}
}

// partition the map into regions using Voronoi cells on the hex grid
func assignTilesToRegions(world *WorldMap, numRegions int) {
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
	}
}

func setRegionCentroids(world *WorldMap) {
	for i := 0; i < RegionCount; i++ {
		region := &world.Regions[i]
		if len(region.TileIds) == 0 {
			continue
		}
		avgX := 0.0
		avgY := 0.0
		for _, tileId := range region.TileIds {
			t := &world.Tiles[tileId]
			avgX += float64(t.Col)
			avgY += float64(t.Row)
		}
		avgX /= float64(len(region.TileIds))
		avgY /= float64(len(region.TileIds))
		region.Centroid = [2]int{int(math.Round(avgX)), int(math.Round(avgY))}
	}
}

func makeRegionBiomes(world *WorldMap) []RegionBuildData {
	regions := make([]RegionBuildData, 0, RegionCount)
	chaosRegions := make(map[int]bool, ChaosRegionCount)
	sylvanRegions := make(map[int]bool, SylvanRegionCount)
	for i := 0; i < ChaosRegionCount; i++ {
		chaosRegions[rand.Intn(RegionCount)] = true
	}
	for i := 0; i < SylvanRegionCount; i++ {
		sylvanRegions[rand.Intn(RegionCount)] = true
	}
	for i := 0; i < RegionCount; i++ {
		region := world.Regions[i]

		// type
		isOcean := region.Centroid[1] < OceanBoundary || region.Centroid[0] < OceanBoundary || region.Centroid[0] > MapWidth-OceanBoundary || region.Centroid[1] > MapHeight-OceanBoundary
		regionType := NoRegionType
		isChaos := chaosRegions[i]
		isSylvan := sylvanRegions[i]
		isInnerSea := rand.Intn(100) < InnerSeaRatio
		if isOcean || isInnerSea {
			regionType = Ocean
		}
		if isChaos {
			regionType = Chaos
		}
		if isSylvan {
			regionType = Sylvan
		}

		// temperature
		temperature := MediumTemperature
		if region.Centroid[1] < MapHeight/4 {
			temperature = LowTemperature
		}
		if region.Centroid[1] > MapHeight/3 && region.Centroid[1] < 2*MapHeight/3 {
			temperature = HighTemperature
		}

		regions = append(regions, RegionBuildData{
			Type:        regionType,
			Elevation:   utils.GetRandomFromArray([]RegionElevation{LowElevation, LowElevation, LowElevation, MediumElevation, MediumElevation, HighElevation}),
			Humidity:    utils.GetRandomFromArray([]RegionHumidity{LowHumidity, MediumHumidity, HighHumidity}),
			Temperature: temperature,
		})
	}
	return regions
}

func assignTerrainToTiles(world *WorldMap, regionBiomes []RegionBuildData) {
	for idx := range world.Tiles {
		t := &world.Tiles[idx]
		biome := regionBiomes[t.RegionId]
		switch biome.Type {
		case Ocean:
			t.Terrain = Water
		case Chaos:
			t.Terrain = Wasteland
		case Sylvan:
			t.Terrain = Forest
		default:
			t.Terrain = getTileTerrainFromBiome(biome)
		}
	}
}

func getTileTerrainFromBiome(biome RegionBuildData) TerrainType {
	// first check for mountains
	elevationAdjustments := RegionElevationToTerrain[biome.Elevation]
	if rand.Intn(100) < elevationAdjustments[Mountains] {
		return Mountains
	}
	// else get random other based on temperature and humidity
	prevalences := make(map[TerrainType]int, len(BaseTerrainPrevalence))
	for k, v := range BaseTerrainPrevalence {
		prevalences[k] = v
	}

	biomeAdjustments := BiomeToTerrain[biome.Humidity][biome.Temperature]
	for terrain := range prevalences {
		prevalences[terrain] += biomeAdjustments[terrain]
		if prevalences[terrain] < 0 {
			prevalences[terrain] = 0
		}
	}
	//fmt.Println(prevalences)
	return utils.GetWeightedRandomFromMap(prevalences)
}

func assignFeaturesToRegions(world *WorldMap, regionBiomes []RegionBuildData) {
	for i := 0; i < RegionCount; i++ {
		region := &world.Regions[i]

		if regionBiomes[i].Type == Ocean || regionBiomes[i].Type == Chaos || regionBiomes[i].Type == Sylvan {
			continue
		}

		// city towards center of region
		centroidTile := world.GetTileAt(region.Centroid[0], region.Centroid[1])
		if centroidTile.Terrain != Water {
			centroidTile.Feature = City
		} else {
			for _, tileId := range region.TileIds {
				tile := &world.Tiles[tileId]
				if tile.Terrain != Water {
					tile.Feature = City
					break
				}
			}
		}
	}
}

func computeRegionBorders(world *WorldMap) {
	for idx := range world.Tiles {
		tile := &world.Tiles[idx]
		borderDirections := int8(0)
		for direction := 0; direction < 6; direction++ {
			if !isFromSameRegion(world, tile, direction) {
				borderDirections |= 1 << direction
			}
		}
		tile.Border = borderDirections
	}
}

// 0-5 represent directions 0-5 (E, NE, NW, W, SW, SE)
func isFromSameRegion(world *WorldMap, tile *Tile, direction int) bool {
	var deltas [6][2]int
	if tile.Row%2 == 0 { // even rows
		deltas = EvenRowDeltas
	} else { // odd rows
		deltas = OddRowDeltas
	}
	if tile.Col+deltas[direction][0] < 0 ||
		tile.Col+deltas[direction][0] >= world.Width ||
		tile.Row+deltas[direction][1] < 0 ||
		tile.Row+deltas[direction][1] >= world.Height {
		return false
	}
	otherTile := world.GetTileAt(tile.Col+deltas[direction][0], tile.Row+deltas[direction][1])
	if otherTile == nil {
		return false
	}
	return tile.RegionId == otherTile.RegionId
}
