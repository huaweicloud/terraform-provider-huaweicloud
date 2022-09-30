package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartConnectivityTestReq struct {

	// 地址和端口列表。
	AddressAndPorts []AddressAndPorts `json:"addressAndPorts"`
}

func (o StartConnectivityTestReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartConnectivityTestReq struct{}"
	}

	return strings.Join([]string{"StartConnectivityTestReq", string(data)}, " ")
}
