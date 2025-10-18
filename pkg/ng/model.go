package ng

type GameData struct {
	WorldMap WorldMap
}

type WorldMap struct {
	Width   int
	Height  int
	Tiles   [MapWidth * MapHeight]Tile
	Regions [RegionCount]Region
}

type Tile struct {
	RegionId int
	Col      int
	Row      int
	Terrain  TerrainType
	Feature  FeatureType
	Border   int8 // bit 0-5 represent directions 0-5 (E, NE, NW, W, SW, SE)
}

type Region struct {
	Id       int
	Name     string
	TileIds  []int
	Centroid [2]int
}
