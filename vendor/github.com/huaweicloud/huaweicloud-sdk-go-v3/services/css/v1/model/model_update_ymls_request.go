package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateYmlsRequest Request Object
type UpdateYmlsRequest struct {

	// 指定修改参数配置的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *UpdateYmlsReq `json:"body,omitempty"`
}

func (o UpdateYmlsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateYmlsRequest struct{}"
	}

	return strings.Join([]string{"UpdateYmlsRequest", string(data)}, " ")
}
