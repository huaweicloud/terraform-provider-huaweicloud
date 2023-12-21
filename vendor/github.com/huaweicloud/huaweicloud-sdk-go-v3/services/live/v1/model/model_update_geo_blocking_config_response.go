package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateGeoBlockingConfigResponse Response Object
type UpdateGeoBlockingConfigResponse struct {
	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateGeoBlockingConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateGeoBlockingConfigResponse struct{}"
	}

	return strings.Join([]string{"UpdateGeoBlockingConfigResponse", string(data)}, " ")
}
