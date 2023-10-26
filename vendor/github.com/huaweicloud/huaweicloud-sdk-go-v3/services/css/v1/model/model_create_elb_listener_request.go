package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateElbListenerRequest Request Object
type CreateElbListenerRequest struct {

	// 指定待更改集群名称的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *CreateEsListenerRequestBody `json:"body,omitempty"`
}

func (o CreateElbListenerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateElbListenerRequest struct{}"
	}

	return strings.Join([]string{"CreateElbListenerRequest", string(data)}, " ")
}
