package ng

// odd-r neighbor deltas ([col,row] offsets)
var EvenRowDeltas = [6][2]int{{+1, 0}, {0, -1}, {-1, -1}, {-1, 0}, {-1, +1}, {0, +1}}
var OddRowDeltas = [6][2]int{{+1, 0}, {+1, -1}, {0, -1}, {-1, 0}, {0, +1}, {+1, +1}}

func (worldMap *WorldMap) GetTileAt(col, row int) *Tile {
	return &worldMap.Tiles[row*worldMap.Width+col]
}

// GetAdjacentTiles returns neighboring tiles around (col,row) using odd-r offset coordinates.
func (worldMap *WorldMap) GetAdjacentTiles(col, row int) []*Tile {

	var deltas [6][2]int
	if row%2 == 0 { // even rows
		deltas = EvenRowDeltas
	} else { // odd rows
		deltas = OddRowDeltas
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
