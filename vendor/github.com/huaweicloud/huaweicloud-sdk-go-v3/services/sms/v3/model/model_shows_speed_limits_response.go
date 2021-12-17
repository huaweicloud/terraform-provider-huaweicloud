package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowsSpeedLimitsResponse struct {
	// 按时间段限速信息

	SpeedLimit     *[]SpeedLimitlJson `json:"speed_limit,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ShowsSpeedLimitsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowsSpeedLimitsResponse struct{}"
	}

	return strings.Join([]string{"ShowsSpeedLimitsResponse", string(data)}, " ")
}
