package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartVpecpRequest Request Object
type StartVpecpRequest struct {

	// 指定开启终端节点的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *StartVpecpReq `json:"body,omitempty"`
}

func (o StartVpecpRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartVpecpRequest struct{}"
	}

	return strings.Join([]string{"StartVpecpRequest", string(data)}, " ")
}
