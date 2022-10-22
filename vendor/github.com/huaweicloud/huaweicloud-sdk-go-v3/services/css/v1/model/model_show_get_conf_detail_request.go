package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowGetConfDetailRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`

	// 配置文件名称。
	Name string `json:"name"`
}

func (o ShowGetConfDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGetConfDetailRequest struct{}"
	}

	return strings.Join([]string{"ShowGetConfDetailRequest", string(data)}, " ")
}
