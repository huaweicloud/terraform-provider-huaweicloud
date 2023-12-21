package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListGeoBlockingConfigRequest Request Object
type ListGeoBlockingConfigRequest struct {

	// 播放域名
	PlayDomain string `json:"play_domain"`
}

func (o ListGeoBlockingConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGeoBlockingConfigRequest struct{}"
	}

	return strings.Join([]string{"ListGeoBlockingConfigRequest", string(data)}, " ")
}
