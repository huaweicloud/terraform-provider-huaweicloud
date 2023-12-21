package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateGeoBlockingConfigRequest Request Object
type UpdateGeoBlockingConfigRequest struct {

	// 播放域名
	PlayDomain string `json:"play_domain"`

	Body *GeoBlockingConfigInfo `json:"body,omitempty"`
}

func (o UpdateGeoBlockingConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateGeoBlockingConfigRequest struct{}"
	}

	return strings.Join([]string{"UpdateGeoBlockingConfigRequest", string(data)}, " ")
}
