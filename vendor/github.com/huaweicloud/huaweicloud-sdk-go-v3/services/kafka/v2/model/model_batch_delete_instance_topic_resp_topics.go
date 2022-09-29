package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchDeleteInstanceTopicRespTopics struct {

	// Topic名称。
	Id *string `json:"id,omitempty"`

	// 是否删除成功。
	Success *bool `json:"success,omitempty"`
}

func (o BatchDeleteInstanceTopicRespTopics) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteInstanceTopicRespTopics struct{}"
	}

	return strings.Join([]string{"BatchDeleteInstanceTopicRespTopics", string(data)}, " ")
}
