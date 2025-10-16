package utils

import "math"

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
