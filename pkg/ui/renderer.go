package ui

import (
	"gociv/pkg/ng"
	"log"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Renderer owns rendering-related state (camera, fonts, sprites, etc.).
// It starts minimal and will grow as features are added.
type Renderer struct {
	Camera     rl.Camera2D
	TerrainTex map[ng.TerrainType]rl.Texture2D
	FeatureTex map[ng.FeatureType]rl.Texture2D
}

func NewRenderer() *Renderer {
	r := &Renderer{
		Camera: rl.Camera2D{
			Target:   rl.Vector2{X: 0, Y: 0},
			Offset:   rl.Vector2{X: float32(rl.GetScreenWidth()) / 2, Y: float32(rl.GetScreenHeight()) / 2},
			Rotation: 0,
			Zoom:     1,
		},
		TerrainTex: make(map[ng.TerrainType]rl.Texture2D),
		FeatureTex: make(map[ng.FeatureType]rl.Texture2D),
	}

	start := time.Now()
	// Load terrain textures (PNGs). Files should exist at these paths.
	r.TerrainTex[ng.Plains] = rl.LoadTexture("assets/images/tiles/plains.png")
	r.TerrainTex[ng.Forest] = rl.LoadTexture("assets/images/tiles/forest.png")
	r.TerrainTex[ng.Mountains] = rl.LoadTexture("assets/images/tiles/mountains.png")
	r.TerrainTex[ng.Water] = rl.LoadTexture("assets/images/tiles/water.png")
	r.TerrainTex[ng.Grassland] = rl.LoadTexture("assets/images/tiles/grasslands.png")
	r.TerrainTex[ng.Ice] = rl.LoadTexture("assets/images/tiles/ice.png")
	r.TerrainTex[ng.Savannah] = rl.LoadTexture("assets/images/tiles/savannah.png")
	r.TerrainTex[ng.Swamp] = rl.LoadTexture("assets/images/tiles/swamp.png")
	r.TerrainTex[ng.Wasteland] = rl.LoadTexture("assets/images/tiles/wasteland.png")

	// Load feature textures (PNGs). Files should exist at these paths.
	r.FeatureTex[ng.City] = rl.LoadTexture("assets/images/features/city.png")

	log.Printf("Loaded terrain textures in %v", time.Since(start))

	return r
}

// Close releases renderer-owned GPU resources.
func (r *Renderer) Close() {
	for _, tex := range r.TerrainTex {
		if tex.ID != 0 {
			rl.UnloadTexture(tex)
		}
	}
}

// Render draws the world using the renderer's camera.
func (r *Renderer) Render(mapData *ng.WorldMap) {
	rl.BeginMode2D(r.Camera)
	DrawMap(mapData, r.Camera, r)
	rl.EndMode2D()
}
