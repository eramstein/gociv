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
}

type Region struct {
	Id       int
	Name     string
	TileIds  []int
	Centroid [2]int
}
