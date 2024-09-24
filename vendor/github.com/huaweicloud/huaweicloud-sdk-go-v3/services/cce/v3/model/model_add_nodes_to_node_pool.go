package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddNodesToNodePool 自定义节点池纳管节点参数。
type AddNodesToNodePool struct {

	// 服务器ID，获取方式请参见ECS/BMS相关资料。
	ServerID string `json:"serverID"`
}

func (o AddNodesToNodePool) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddNodesToNodePool struct{}"
	}

	return strings.Join([]string{"AddNodesToNodePool", string(data)}, " ")
}
