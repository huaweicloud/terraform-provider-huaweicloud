package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateClusterNameReq struct {

	// 修改后集群名称。
	DisplayName string `json:"displayName"`
}

func (o UpdateClusterNameReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateClusterNameReq struct{}"
	}

	return strings.Join([]string{"UpdateClusterNameReq", string(data)}, " ")
}
