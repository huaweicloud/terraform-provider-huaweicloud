package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartConnectivityTestRequest Request Object
type StartConnectivityTestRequest struct {

	// 指定待测试的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *StartConnectivityTestReq `json:"body,omitempty"`
}

func (o StartConnectivityTestRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartConnectivityTestRequest struct{}"
	}

	return strings.Join([]string{"StartConnectivityTestRequest", string(data)}, " ")
}
