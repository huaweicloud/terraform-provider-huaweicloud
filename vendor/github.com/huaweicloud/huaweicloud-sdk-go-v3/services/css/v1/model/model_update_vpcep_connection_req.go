package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateVpcepConnectionReq struct {

	// 期望的操作行为。 - receive: 允许连接 - reject: 拒绝连接
	Action string `json:"action"`

	// 终端节点ID列表。
	EndpointIdList []string `json:"endpointIdList"`
}

func (o UpdateVpcepConnectionReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateVpcepConnectionReq struct{}"
	}

	return strings.Join([]string{"UpdateVpcepConnectionReq", string(data)}, " ")
}
