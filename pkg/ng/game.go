package ng

func NewGame() *GameData {
	return &GameData{
		WorldMap: *NewWorldMap(),
	}
}
