package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartVpecpReq struct {

	// 开启终端节点。
	EndpointWithDnsName bool `json:"endpointWithDnsName"`
}

func (o StartVpecpReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartVpecpReq struct{}"
	}

	return strings.Join([]string{"StartVpecpReq", string(data)}, " ")
}
