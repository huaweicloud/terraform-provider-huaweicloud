package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAiOpsRequest Request Object
type CreateAiOpsRequest struct {

	// 指定待操作的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *CreateAiOpsRequestBody `json:"body,omitempty"`
}

func (o CreateAiOpsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAiOpsRequest struct{}"
	}

	return strings.Join([]string{"CreateAiOpsRequest", string(data)}, " ")
}
