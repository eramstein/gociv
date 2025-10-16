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
}

type Region struct {
	Id      int
	Name    string
	TileIds []int
	Biome   TerrainType
}
