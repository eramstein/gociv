package utils

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Pre-computed hexagon vertex offsets (relative to center)
// These are the unit vectors for a regular hexagon, scaled by radius
var HexVertexOffsets = [6]rl.Vector2{
	{X: 0.8660254, Y: 0.5},   // Vertex 0
	{X: 0.8660254, Y: -0.5},  // Vertex 1
	{X: 0.0, Y: -1.0},        // Vertex 2
	{X: -0.8660254, Y: -0.5}, // Vertex 3
	{X: -0.8660254, Y: 0.5},  // Vertex 4
	{X: 0.0, Y: 1.0},         // Vertex 5
}

// Map border directions to hexagon edges
// Each direction corresponds to the edge between two vertices
var EdgeMapping = [6][2]int{
	{0, 1}, // Direction 0 (East) -> edge between vertex 0 and 1
	{1, 2}, // Direction 1 (NE) -> edge between vertex 1 and 2
	{2, 3}, // Direction 2 (NW) -> edge between vertex 2 and 3
	{3, 4}, // Direction 3 (West) -> edge between vertex 3 and 4
	{4, 5}, // Direction 4 (SW) -> edge between vertex 4 and 5
	{5, 0}, // Direction 5 (SE) -> edge between vertex 5 and 0
}

// oddRToCube converts odd-r offset (col,row) to cube coordinates (x,y,z) for pointy-top hexes.
func OddRToCube(col, row int) (int, int, int) {
	q := col - (row-(row&1))/2 // axial q
	r := row                   // axial r
	x := q
	z := r
	y := -x - z
	return x, y, z
}

// cubeToOddR converts cube coords back to odd-r offset (col,row).
func CubeToOddR(x, y, z int) (int, int) {
	// axial from cube
	q := x
	r := z
	// odd-r from axial
	row := r
	col := q + (row-(row&1))/2
	return col, row
}

// cubeRound rounds floating cube coordinates to the nearest valid cube coordinate.
func CubeRound(xf, yf, zf float64) (int, int, int) {
	rx := math.Round(xf)
	ry := math.Round(yf)
	rz := math.Round(zf)

	dx := math.Abs(rx - xf)
	dy := math.Abs(ry - yf)
	dz := math.Abs(rz - zf)

	if dx > dy && dx > dz {
		rx = -ry - rz
	} else if dy > dz {
		ry = -rx - rz
	} else {
		rz = -rx - ry
	}

	return int(rx), int(ry), int(rz)
}

// GetHexDistance computes distance between two hexes given odd-r offset (pointy-top) coords.
func GetHexDistance(col1, row1, col2, row2 int) int {
	// Convert odd-r to axial (q,r)
	q1 := col1 - (row1-(row1&1))/2
	q2 := col2 - (row2-(row2&1))/2
	r1a := row1
	r2a := row2

	// Axial to cube (x=q, z=r, y=-x-z)
	x1, z1 := q1, r1a
	y1 := -x1 - z1
	x2, z2 := q2, r2a
	y2 := -x2 - z2

	dx := absInt(x1 - x2)
	dy := absInt(y1 - y2)
	dz := absInt(z1 - z2)
	return (dx + dy + dz) / 2
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
