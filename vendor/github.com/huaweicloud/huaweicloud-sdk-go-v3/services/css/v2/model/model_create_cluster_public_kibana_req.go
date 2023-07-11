package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterPublicKibanaReq Kibana公网访问信息。只有在authorityEnable设置为true时该参数配置生效。
type CreateClusterPublicKibanaReq struct {

	// 带宽大小。
	EipSize int32 `json:"eipSize"`

	ElbWhiteList *CreateClusterPublicKibanaElbWhiteList `json:"elbWhiteList"`
}

func (o CreateClusterPublicKibanaReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterPublicKibanaReq struct{}"
	}

	return strings.Join([]string{"CreateClusterPublicKibanaReq", string(data)}, " ")
}
