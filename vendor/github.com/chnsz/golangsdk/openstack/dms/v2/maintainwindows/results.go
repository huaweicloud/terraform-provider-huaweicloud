package maintainwindows

// GetResponse response
type GetResponse struct {
	MaintainWindows []MaintainWindow `json:"maintain_windows"`
}

// MaintainWindow for dms
type MaintainWindow struct {
	ID      int    `json:"seq"`
	Begin   string `json:"begin"`
	End     string `json:"end"`
	Default bool   `json:"default"`
}
