package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAgentsRequest Request Object
type ListAgentsRequest struct {

	// - 查询集群主机时，填写集群id。 - 查询用户自定义主机时，填写“APM”。
	ClusterId string `json:"cluster_id"`

	// - 查询集群主机时，填写命名空间。 - 查询用户自定义主机时，填写“APM”。
	Namespace string `json:"namespace"`
}

func (o ListAgentsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAgentsRequest struct{}"
	}

	return strings.Join([]string{"ListAgentsRequest", string(data)}, " ")
}
