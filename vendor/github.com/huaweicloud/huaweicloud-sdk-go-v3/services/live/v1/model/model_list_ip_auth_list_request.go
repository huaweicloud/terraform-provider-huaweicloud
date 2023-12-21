package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListIpAuthListRequest Request Object
type ListIpAuthListRequest struct {

	// 推流域名或播放域名
	Domain string `json:"domain"`
}

func (o ListIpAuthListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListIpAuthListRequest struct{}"
	}

	return strings.Join([]string{"ListIpAuthListRequest", string(data)}, " ")
}
