package availablezones

import (
	"github.com/huaweicloud/golangsdk"
)

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
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*GetResponse, error) {
	var s GetResponse
	err := r.Result.ExtractInto(&s)
	return &s, err
}
