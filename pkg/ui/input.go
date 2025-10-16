package ui

import (
	"gociv/pkg/ng"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ClickEvent struct {
	Col int32
	Row int32
}

// UpdateInput listens for a left mouse click and converts the pixel position
// to odd-r pointy-top hex coordinates using the current tile size.
// Returns a ClickEvent when a tile is clicked, or nil if no click.
func UpdateInput(gameData *ng.GameData, renderer *Renderer) {
	// Quick save/load
	if rl.IsKeyPressed(rl.KeyF5) {
		if err := ng.SaveGameToFile("quicksave", *gameData); err != nil {
			log.Printf("quicksave failed: %v", err)
		}
	}
	if rl.IsKeyPressed(rl.KeyF4) {
		if err := ng.LoadGameFromFile("quicksave", gameData); err != nil {
			log.Printf("quickload failed: %v", err)
		}
	}

	// Click to interact with tiles using world-space position
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		screen := rl.GetMousePosition()
		world := rl.GetScreenToWorld2D(screen, renderer.Camera)
		col, row := PixelToHexagon(world.X, world.Y)

		// bounds check
		if col >= 0 && row >= 0 && int(col) < gameData.WorldMap.Width && int(row) < gameData.WorldMap.Height {
			tile := gameData.WorldMap.GetTileAt(int(col), int(row))
			tile.UpdateTerrain(ng.Forest)
			adjacentTiles := gameData.WorldMap.GetAdjacentTiles(int(col), int(row))
			for _, adjacentTile := range adjacentTiles {
				adjacentTile.UpdateTerrain(ng.Mountains)
			}
		}
	}
}
