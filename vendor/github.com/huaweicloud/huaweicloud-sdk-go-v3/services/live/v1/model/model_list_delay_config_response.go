package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDelayConfigResponse Response Object
type ListDelayConfigResponse struct {

	// 播放域名
	PlayDomain *string `json:"play_domain,omitempty"`

	// 直播延时配置
	DelayConfig    *[]DelayConfig `json:"delay_config,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListDelayConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDelayConfigResponse struct{}"
	}

	return strings.Join([]string{"ListDelayConfigResponse", string(data)}, " ")
}
