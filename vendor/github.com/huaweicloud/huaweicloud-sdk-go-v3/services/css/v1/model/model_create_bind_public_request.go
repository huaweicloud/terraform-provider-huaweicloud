package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateBindPublicRequest struct {

	// 指定开启公网访问的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *BindPublicReq `json:"body,omitempty"`
}

func (o CreateBindPublicRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateBindPublicRequest struct{}"
	}

	return strings.Join([]string{"CreateBindPublicRequest", string(data)}, " ")
}
