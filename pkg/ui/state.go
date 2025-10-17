package ui

// UIState holds the current UI state
type UIState struct {
	SelectedRegionID int // -1 means no region selected
}

// Global UI state instance
var State = &UIState{
	SelectedRegionID: -1, // No region selected by default
}
