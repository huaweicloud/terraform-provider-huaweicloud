package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartPublicWhitelistRequest Request Object
type StartPublicWhitelistRequest struct {

	// 指定开启公网访问控制白名单集群ID。
	ClusterId string `json:"cluster_id"`

	Body *StartPublicWhitelistReq `json:"body,omitempty"`
}

func (o StartPublicWhitelistRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartPublicWhitelistRequest struct{}"
	}

	return strings.Join([]string{"StartPublicWhitelistRequest", string(data)}, " ")
}
