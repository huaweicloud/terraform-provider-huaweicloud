package availablezones

// GetResponse response
type GetResponse struct {
	RegionID       string          `json:"regionId"`
	AvailableZones []AvailableZone `json:"available_zones"`
}

// AvailableZone for dms
type AvailableZone struct {
	ID                   string `json:"id"`
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	Port                 string `json:"port"`
	ResourceAvailability string `json:"resource_availability"`
	SoldOut              bool   `json:"soldOut"`
	DefaultAz            bool   `json:"default_az"`
	RemainTime           uint64 `json:"remain_time"`
	Ipv6Enable           bool   `json:"ipv6_enable"`
}
