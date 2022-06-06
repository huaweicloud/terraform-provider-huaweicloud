package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateCnfRequest struct {

	// 指定更新配置文件的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *CreateCnfReq `json:"body,omitempty"`
}

func (o UpdateCnfRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCnfRequest struct{}"
	}

	return strings.Join([]string{"UpdateCnfRequest", string(data)}, " ")
}
