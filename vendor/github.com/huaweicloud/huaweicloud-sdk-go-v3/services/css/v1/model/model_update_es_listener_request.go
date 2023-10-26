package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateEsListenerRequest Request Object
type UpdateEsListenerRequest struct {

	// 指定待操作的集群ID。
	ClusterId string `json:"cluster_id"`

	// 指定待操作得监听器ID。
	ListenerId string `json:"listener_id"`

	Body *UpdateEsListenerRequestBody `json:"body,omitempty"`
}

func (o UpdateEsListenerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateEsListenerRequest struct{}"
	}

	return strings.Join([]string{"UpdateEsListenerRequest", string(data)}, " ")
}
