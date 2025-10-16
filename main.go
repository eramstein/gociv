package main

import (
	"gociv/pkg/ng"
	"gociv/pkg/ui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1600, 900, "Gociv")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	gameData := ng.NewGame()
	renderer := ui.NewRenderer()
	defer renderer.Close()

	for !rl.WindowShouldClose() {
		// Update
		ui.UpdateCamera(renderer, &gameData.WorldMap)
		ui.UpdateInput(gameData, renderer)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		renderer.Render(&gameData.WorldMap)
		// rl.DrawFPS(700, 10)
		rl.EndDrawing()
	}
}
