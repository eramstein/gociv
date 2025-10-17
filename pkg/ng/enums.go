package ng

type TerrainType int

const (
	NoTerrain TerrainType = iota
	Plains
	Forest
	Mountains
	Water
	Grassland
	Ice
	Savannah
	Swamp
	Wasteland
)

type FeatureType int

const (
	NoFeature FeatureType = iota
	City
)
