package ng

const (
	MapWidth          = 80
	MapHeight         = 50
	RegionCount       = 200
	OceanBoundary     = 5  // if the region centroid is this close to map border, it is ocean
	ChaosRegionCount  = 2  // number of regions that will be chaos
	SylvanRegionCount = 1  // number of regions that will be forest
	InnerSeaRatio     = 10 // percentage chance of a region being an inner sea
)

// these fields are only useful during world building
type RegionBuildData struct {
	Type        RegionType
	Elevation   RegionElevation
	Humidity    RegionHumidity
	Temperature RegionTemperature
}

// special regions types
type RegionType int

const (
	NoRegionType RegionType = iota
	Ocean
	Chaos
	Sylvan
)

type RegionElevation int

const (
	NoRegionElevation RegionElevation = iota
	LowElevation
	MediumElevation
	HighElevation
)

type RegionHumidity int

const (
	NoRegionHumidity RegionHumidity = iota
	LowHumidity
	MediumHumidity
	HighHumidity
)

type RegionTemperature int

const (
	NoRegionTemperature RegionTemperature = iota
	LowTemperature
	MediumTemperature
	HighTemperature
)

var BaseTerrainPrevalence = map[TerrainType]int{
	Plains:    50,
	Forest:    60,
	Mountains: 30,
	Water:     10,
	Grassland: 50,
	Ice:       0,
	Savannah:  20,
	Swamp:     10,
	Wasteland: 0,
}

var RegionElevationToTerrain = map[RegionElevation]map[TerrainType]int{
	LowElevation:    {Mountains: 5},
	MediumElevation: {Mountains: 20},
	HighElevation:   {Mountains: 50},
}

var BiomeToTerrain = map[RegionHumidity]map[RegionTemperature]map[TerrainType]int{
	LowHumidity: {
		LowTemperature:    {Ice: 80, Plains: 60},
		MediumTemperature: {Forest: 50, Plains: 60},
		HighTemperature:   {Savannah: 100, Plains: 60},
	},
	MediumHumidity: {
		LowTemperature:    {Ice: 40, Forest: 80},
		MediumTemperature: {Forest: 100},
		HighTemperature:   {Forest: 80},
	},
	HighHumidity: {
		LowTemperature:    {Ice: 40, Swamp: 20, Grassland: 60},
		MediumTemperature: {Swamp: 20, Forest: 60, Grassland: 60},
		HighTemperature:   {Swamp: 20, Grassland: 60},
	},
}
