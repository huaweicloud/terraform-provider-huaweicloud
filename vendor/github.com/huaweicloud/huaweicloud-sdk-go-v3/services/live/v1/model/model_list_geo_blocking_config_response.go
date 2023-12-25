package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListGeoBlockingConfigResponse Response Object
type ListGeoBlockingConfigResponse struct {

	// 直播播放域名
	PlayDomain *string `json:"play_domain,omitempty"`

	// 应用列表
	Apps *[]GeoBlockingConfigInfo `json:"apps,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListGeoBlockingConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGeoBlockingConfigResponse struct{}"
	}

	return strings.Join([]string{"ListGeoBlockingConfigResponse", string(data)}, " ")
}
