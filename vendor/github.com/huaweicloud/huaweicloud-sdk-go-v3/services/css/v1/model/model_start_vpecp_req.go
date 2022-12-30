package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartVpecpReq struct {

	// 是否开启内网域名。 - ture：表示开启。 - false： 表示不开启。
	EndpointWithDnsName *bool `json:"endpointWithDnsName,omitempty"`
}

func (o StartVpecpReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartVpecpReq struct{}"
	}

	return strings.Join([]string{"StartVpecpReq", string(data)}, " ")
}
