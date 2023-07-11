package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddIndependentNodeRequest Request Object
type AddIndependentNodeRequest struct {

	// 指定需要独立master或client的集群ID。
	ClusterId string `json:"cluster_id"`

	// 指定待新增独立节点类型。 - ess-master：Master节点。 - ess-client：Client节点。
	Type string `json:"type"`

	Body *IndependentReq `json:"body,omitempty"`
}

func (o AddIndependentNodeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddIndependentNodeRequest struct{}"
	}

	return strings.Join([]string{"AddIndependentNodeRequest", string(data)}, " ")
}
