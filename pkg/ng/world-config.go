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
	LowElevation:    {Mountains: -30},
	MediumElevation: {Mountains: 30},
	HighElevation:   {Mountains: 150, Swamp: -10},
}

var RegionHumidityToTerrain = map[RegionHumidity]map[TerrainType]int{
	LowHumidity:    {Grassland: -20, Swamp: -20, Savannah: 30, Plains: 40},
	MediumHumidity: {Forest: 40},
	HighHumidity:   {Grassland: 40, Swamp: 40, Water: 20},
}

var RegionTemperatureToTerrain = map[RegionTemperature]map[TerrainType]int{
	LowTemperature:    {Ice: 150, Grassland: -20, Savannah: -50},
	MediumTemperature: {},
	HighTemperature:   {Ice: -80, Savannah: 50, Plains: 20, Swamp: 20},
}
