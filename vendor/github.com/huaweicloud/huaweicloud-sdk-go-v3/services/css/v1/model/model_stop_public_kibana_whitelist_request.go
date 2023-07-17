package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopPublicKibanaWhitelistRequest Request Object
type StopPublicKibanaWhitelistRequest struct {

	// 指定关闭Kibana公网访问控制的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o StopPublicKibanaWhitelistRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopPublicKibanaWhitelistRequest struct{}"
	}

	return strings.Join([]string{"StopPublicKibanaWhitelistRequest", string(data)}, " ")
}
