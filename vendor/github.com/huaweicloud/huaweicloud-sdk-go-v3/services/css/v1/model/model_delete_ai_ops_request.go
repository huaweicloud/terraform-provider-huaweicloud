package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAiOpsRequest Request Object
type DeleteAiOpsRequest struct {

	// 指定待删除的集群ID。
	ClusterId string `json:"cluster_id"`

	// 指定检测任务ID。
	AiopsId string `json:"aiops_id"`
}

func (o DeleteAiOpsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAiOpsRequest struct{}"
	}

	return strings.Join([]string{"DeleteAiOpsRequest", string(data)}, " ")
}
