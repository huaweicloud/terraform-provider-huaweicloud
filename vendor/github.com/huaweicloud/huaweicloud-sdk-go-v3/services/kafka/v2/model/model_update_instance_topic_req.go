package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceTopicReq 修改的topic列表。
type UpdateInstanceTopicReq struct {

	// 修改的topic列表。
	Topics *[]UpdateInstanceTopicReqTopics `json:"topics,omitempty"`
}

func (o UpdateInstanceTopicReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceTopicReq struct{}"
	}

	return strings.Join([]string{"UpdateInstanceTopicReq", string(data)}, " ")
}
