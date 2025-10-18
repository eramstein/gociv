package ui

import (
	"gociv/pkg/ng"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	TileSize = float32(60)
)

var TerrainTypeColors = map[ng.TerrainType]rl.Color{
	ng.Plains:    rl.Green,
	ng.Forest:    rl.DarkGreen,
	ng.Mountains: rl.DarkGray,
}

/*
The map hexagons are pointy top, with an odd-r offset to make the map rectangular.
*/
func DrawMap(mapData *ng.WorldMap, camera rl.Camera2D, renderer *Renderer) {
	hexWidth := float32(math.Sqrt(3)) * TileSize
	halfHexW := hexWidth / 2
	hexHeight := 2 * TileSize

	// World-space visible rect with a small margin so edges draw fully
	halfW := (float32(rl.GetScreenWidth()) / 2) / camera.Zoom
	halfH := (float32(rl.GetScreenHeight()) / 2) / camera.Zoom
	minX := camera.Target.X - halfW - halfHexW
	maxX := camera.Target.X + halfW + halfHexW
	minY := camera.Target.Y - halfH - TileSize
	maxY := camera.Target.Y + halfH + TileSize

	// Limit iteration to rows intersecting the vertical view range
	vertStep := float32(1.5) * TileSize
	firstRow := int(math.Floor(float64((minY - TileSize) / vertStep)))
	lastRow := int(math.Ceil(float64((maxY - TileSize) / vertStep)))
	if firstRow < 0 {
		firstRow = 0
	}
	if lastRow >= mapData.Height {
		lastRow = mapData.Height - 1
	}

	for row := firstRow; row <= lastRow; row++ {
		rowStart := row * mapData.Width
		rowEnd := rowStart + mapData.Width
		for i := rowStart; i < rowEnd; i++ {
			tile := mapData.Tiles[i]
			col := int32(i % mapData.Width)
			r32 := int32(row)

			center := HexagonToPixel(col, r32)

			// Culling: skip hex if fully outside view
			if center.X+halfHexW < minX || center.X-halfHexW > maxX || center.Y+TileSize < minY || center.Y-TileSize > maxY {
				continue
			}

			drawTile(&tile, renderer, center, hexWidth, hexHeight)
		}
	}
}

func drawTile(tile *ng.Tile, renderer *Renderer, center rl.Vector2, hexWidth float32, hexHeight float32) {
	// determine color modulation based on selection
	var color rl.Color
	if State.SelectedRegionID == -1 || State.SelectedRegionID == tile.RegionId {
		color = rl.White
	} else {
		color = rl.Fade(rl.White, 0.75)
	}

	// draw textured hex (assumes texture exists for each terrain)
	tex := renderer.TerrainTex[tile.Terrain]
	src := rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height))
	dst := rl.NewRectangle(center.X, center.Y, hexWidth, hexHeight)
	origin := rl.NewVector2(hexWidth/2, hexHeight/2)
	rl.DrawTexturePro(tex, src, dst, origin, 0, color)

	// draw the feature if it exists
	if tile.Feature != ng.NoFeature {
		tex := renderer.FeatureTex[tile.Feature]
		src := rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height))
		dst := rl.NewRectangle(center.X, center.Y, hexWidth/2, hexHeight/2)
		origin := rl.NewVector2(hexWidth/4, hexHeight/4)
		rl.DrawTexturePro(tex, src, dst, origin, 0, color)
	}

	//rl.DrawText(fmt.Sprintf("%d,%d [%s]", tile.Col, tile.Row, borderStr), int32(center.X)-int32(TileSize/2), int32(center.Y)-3, 6, rl.RayWhite)

	// Draw region borders if any exist
	if State.SelectedRegionID == tile.RegionId && tile.Border != 0 {
		drawRegionBorders(tile, center, hexWidth, hexHeight)
	}
}

// HexagonToPixel converts odd-r offset hex coordinates (col, row) to pixel center for pointy-top hexes.
// A margin of one radius is applied so tiles are fully visible.
func HexagonToPixel(col, row int32) rl.Vector2 {
	hexWidth := float32(math.Sqrt(3)) * TileSize
	horizStep := hexWidth
	vertStep := 1.5 * TileSize

	cx := TileSize + float32(col)*horizStep
	cy := TileSize + float32(row)*vertStep
	if row%2 != 0 {
		cx += hexWidth / 2.0
	}

	return rl.Vector2{X: cx, Y: cy}
}

// PixelToHexagon converts a pixel position to the nearest odd-r offset hex (col, row)
// for pointy-top hexes. This is an approximation using simple rounding suitable for picking.
func PixelToHexagon(x, y float32) (int32, int32) {
	hexWidth := float32(math.Sqrt(3)) * TileSize
	horizStep := hexWidth
	vertStep := 1.5 * TileSize

	// remove offset
	py := y - TileSize
	row := int32(math.Round(float64(py / vertStep)))

	px := x - TileSize
	if row%2 != 0 {
		px -= hexWidth / 2.0
	}

	col := int32(math.Round(float64(px / horizStep)))
	return col, row
}

// drawRegionBorders draws thick borders on hexagon edges that are region boundaries
func drawRegionBorders(tile *ng.Tile, center rl.Vector2, hexWidth float32, hexHeight float32) {

	borderColor := rl.Black
	borderThickness := float32(5)

	vertices := make([]rl.Vector2, 6)
	radius := hexHeight / 2 // pointy hexagon points are on a circle of radius size/2

	// compute coordinates of the hexagon points
	for i := 0; i < 6; i++ {
		angle := (float64(i) - 0.5) * math.Pi / 3.0 // 60 degrees in radians
		vertices[i] = rl.Vector2{
			X: center.X + radius*float32(math.Cos(angle)),
			Y: center.Y - radius*float32(math.Sin(angle)),
		}
	}

	// Map border directions to hexagon edges
	// Each direction corresponds to the edge between two vertices
	edgeMapping := [6][2]int{
		{0, 1}, // Direction 0 (East) -> edge between vertex 0 and 1
		{1, 2}, // Direction 1 (NE) -> edge between vertex 1 and 2
		{2, 3}, // Direction 2 (NW) -> edge between vertex 2 and 3
		{3, 4}, // Direction 3 (West) -> edge between vertex 3 and 4
		{4, 5}, // Direction 4 (SW) -> edge between vertex 4 and 5
		{5, 0}, // Direction 5 (SE) -> edge between vertex 5 and 0
	}

	// Draw border lines for each direction that has a border
	for direction := 0; direction < 6; direction++ {
		if tile.Border&(1<<direction) != 0 {
			startIdx := edgeMapping[direction][0]
			endIdx := edgeMapping[direction][1]
			start := vertices[startIdx]
			end := vertices[endIdx]
			rl.DrawLineEx(start, end, borderThickness, borderColor)
		}
	}
}
