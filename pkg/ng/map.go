package ng

func (worldMap *WorldMap) GetTileAt(col, row int) *Tile {
	return &worldMap.Tiles[row*worldMap.Width+col]
}

// GetAdjacentTiles returns neighboring tiles around (col,row) using odd-r offset coordinates.
func (worldMap *WorldMap) GetAdjacentTiles(col, row int) []*Tile {
	// odd-r neighbor deltas
	evenRow := [6][2]int{{+1, 0}, {0, -1}, {-1, -1}, {-1, 0}, {-1, +1}, {0, +1}}
	oddRow := [6][2]int{{+1, 0}, {+1, -1}, {0, -1}, {-1, 0}, {0, +1}, {+1, +1}}

	var deltas [6][2]int
	if row%2 == 0 { // even rows
		deltas = evenRow
	} else { // odd rows
		deltas = oddRow
	}

	neighbors := make([]*Tile, 0, 6)
	for _, d := range deltas {
		nc := col + d[0]
		nr := row + d[1]
		if nc >= 0 && nr >= 0 && nc < worldMap.Width && nr < worldMap.Height {
			neighbors = append(neighbors, worldMap.GetTileAt(nc, nr))
		}
	}
	return neighbors
}
