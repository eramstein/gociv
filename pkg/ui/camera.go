package ui

import (
	"gociv/pkg/ng"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// UpdateCamera updates the renderer camera with keyboard pan, right-drag, and mouse wheel zoom.
// Clamps camera target within the map's world extents.
func UpdateCamera(renderer *Renderer, mapData *ng.WorldMap) {
	const panSpeed = float32(600)
	dt := rl.GetFrameTime()

	// Keyboard pan
	move := rl.Vector2{X: 0, Y: 0}
	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		move.X += panSpeed * dt
	}
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		move.X -= panSpeed * dt
	}
	if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		move.Y += panSpeed * dt
	}
	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		move.Y -= panSpeed * dt
	}
	renderer.Camera.Target = rl.Vector2{X: renderer.Camera.Target.X + move.X, Y: renderer.Camera.Target.Y + move.Y}

	// Right mouse drag pan
	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		delta := rl.GetMouseDelta()
		renderer.Camera.Target = rl.Vector2{
			X: renderer.Camera.Target.X - delta.X/renderer.Camera.Zoom,
			Y: renderer.Camera.Target.Y - delta.Y/renderer.Camera.Zoom,
		}
	}

	// Mouse wheel zoom
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		zoom := renderer.Camera.Zoom + wheel*0.1
		if zoom < 0.2 {
			zoom = 0.2
		}
		if zoom > 2.0 {
			zoom = 2.0
		}
		renderer.Camera.Zoom = zoom
	}

	// Clamp target to world bounds considering viewport and zoom
	minX, minY, maxX, maxY := computeWorldBounds(mapData)
	halfW := (float32(rl.GetScreenWidth()) / 2) / renderer.Camera.Zoom
	halfH := (float32(rl.GetScreenHeight()) / 2) / renderer.Camera.Zoom

	worldW := maxX - minX
	worldH := maxY - minY

	// Add a small inner padding so edges don't touch the window borders
	padding := TileSize * 0.5

	if halfW*2 >= worldW {
		renderer.Camera.Target.X = minX + worldW/2
	} else {
		if renderer.Camera.Target.X < minX+halfW+padding {
			renderer.Camera.Target.X = minX + halfW + padding
		}
		if renderer.Camera.Target.X > maxX-halfW-padding {
			renderer.Camera.Target.X = maxX - halfW - padding
		}
	}

	if halfH*2 >= worldH {
		renderer.Camera.Target.Y = minY + worldH/2
	} else {
		if renderer.Camera.Target.Y < minY+halfH+padding {
			renderer.Camera.Target.Y = minY + halfH + padding
		}
		if renderer.Camera.Target.Y > maxY-halfH-padding {
			renderer.Camera.Target.Y = maxY - halfH - padding
		}
	}
}

// computeWorldBounds returns the min/max world coordinates that enclose the map.
func computeWorldBounds(mapData *ng.WorldMap) (float32, float32, float32, float32) {
	// Matches HexagonToPixel for odd-r pointy-top layout.
	// Expand bounds by polygon extents so tiles are fully visible at edges.
	hexWidth := float32(math.Sqrt(3)) * TileSize
	horizStep := hexWidth
	vertStep := 1.5 * TileSize

	// Expand centers by polygon reach: horizontal half-width, vertical radius
	minX := -horizStep / 2
	maxX := float32(mapData.Width+1) * (horizStep)
	minY := -vertStep / 2
	maxY := float32(mapData.Height+1) * (vertStep)

	return minX, minY, maxX, maxY
}
