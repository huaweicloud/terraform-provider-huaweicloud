package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateFlavorRequest struct {

	// 指定待更改规格的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *UpdateFlavorReq `json:"body,omitempty"`
}

func (o UpdateFlavorRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFlavorRequest struct{}"
	}

	return strings.Join([]string{"UpdateFlavorRequest", string(data)}, " ")
}
