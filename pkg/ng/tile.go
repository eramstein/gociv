package ng

func (tile *Tile) UpdateTerrain(terrain TerrainType) {
	tile.Terrain = terrain
}

// GetAdjacentTiles returns neighboring tiles around this tile using odd-r offset coordinates.
// Caller must provide the world to resolve neighbor tiles.
func (tile *Tile) GetAdjacentTiles(world *WorldMap) []*Tile {
	return world.GetAdjacentTiles(tile.Col, tile.Row)
}
