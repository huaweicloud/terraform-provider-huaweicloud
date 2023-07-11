package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopPublicWhitelistRequest Request Object
type StopPublicWhitelistRequest struct {

	// 指定关闭公网访问控制白名单集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o StopPublicWhitelistRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopPublicWhitelistRequest struct{}"
	}

	return strings.Join([]string{"StopPublicWhitelistRequest", string(data)}, " ")
}
