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

			// draw textured hex (assumes texture exists for each terrain)
			tex := renderer.TerrainTex[tile.Terrain]
			src := rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height))
			dst := rl.NewRectangle(center.X, center.Y, hexWidth, hexHeight)
			origin := rl.NewVector2(hexWidth/2, hexHeight/2)

			//compute a stable tint per region using HSV; semi-transparent so base texture shows
			// hueIndex := tile.RegionId % 50
			// hue := float32((hueIndex*37)%50) * 360.0 / 50.0
			// tint := rl.ColorFromHSV(hue, 0.35, 1.0)
			// tint.A = 200
			// rl.DrawTexturePro(tex, src, dst, origin, 0, tint)

			rl.DrawTexturePro(tex, src, dst, origin, 0, rl.White)
			// rl.DrawText(fmt.Sprintf("%d,%d", col, r32), int32(center.X)-int32(TileSize/2), int32(center.Y)-3, 6, rl.RayWhite)
			// rl.DrawText(fmt.Sprintf("%d,%d", mapData.Regions[tile.RegionId].Centroid[0], mapData.Regions[tile.RegionId].Centroid[1]), int32(center.X), int32(center.Y), 16, rl.RayWhite)
		}
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
