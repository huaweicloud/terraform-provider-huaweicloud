package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateRangeSwitchResponse struct {
	OriginRange    *OriginRangeBody `json:"origin_range,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o UpdateRangeSwitchResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateRangeSwitchResponse struct{}"
	}

	return strings.Join([]string{"UpdateRangeSwitchResponse", string(data)}, " ")
}
