package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateFlavorByTypeRequest Request Object
type UpdateFlavorByTypeRequest struct {

	// 指定待更改规格的集群ID。
	ClusterId string `json:"cluster_id"`

	// 指定待更改的集群节点类型。 - ess：数据节点。 - ess-cold：冷数据节点。 - ess-client：Client节点。 - ess-master：Master节点。
	Types string `json:"types"`

	Body *UpdateFlavorByTypeReq `json:"body,omitempty"`
}

func (o UpdateFlavorByTypeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFlavorByTypeRequest struct{}"
	}

	return strings.Join([]string{"UpdateFlavorByTypeRequest", string(data)}, " ")
}
