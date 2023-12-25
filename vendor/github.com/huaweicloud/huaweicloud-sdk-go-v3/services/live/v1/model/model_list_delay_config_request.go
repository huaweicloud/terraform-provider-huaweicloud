package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDelayConfigRequest Request Object
type ListDelayConfigRequest struct {

	// 播放域名
	PlayDomain string `json:"play_domain"`
}

func (o ListDelayConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDelayConfigRequest struct{}"
	}

	return strings.Join([]string{"ListDelayConfigRequest", string(data)}, " ")
}
