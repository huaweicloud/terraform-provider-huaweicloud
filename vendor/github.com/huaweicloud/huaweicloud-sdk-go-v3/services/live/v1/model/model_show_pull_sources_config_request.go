package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPullSourcesConfigRequest Request Object
type ShowPullSourcesConfigRequest struct {

	// 播放域名
	PlayDomain string `json:"play_domain"`
}

func (o ShowPullSourcesConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPullSourcesConfigRequest struct{}"
	}

	return strings.Join([]string{"ShowPullSourcesConfigRequest", string(data)}, " ")
}
