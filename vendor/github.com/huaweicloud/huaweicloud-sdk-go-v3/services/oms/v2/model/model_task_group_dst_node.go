package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 迁移任务组目的端节点信息
type TaskGroupDstNode struct {

	// 目的端桶的AK（最大长度100个字符）。
	Ak string `json:"ak"`

	// 目的端桶的SK（最大长度100个字符）。
	Sk string `json:"sk"`

	// 目的端桶所处的区域。  请与Endpoint对应的区域保持一致。
	Region string `json:"region"`

	// 目的端的桶名称
	Bucket string `json:"bucket"`

	// 华为云目的端信息，默认为HEC
	CloudType *string `json:"cloud_type,omitempty"`

	// 目的端桶内路径前缀（拼接在对象key前面,组成新的key,拼接后不能超过1024个字符）。
	SavePrefix *string `json:"save_prefix,omitempty"`
}

func (o TaskGroupDstNode) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskGroupDstNode struct{}"
	}

	return strings.Join([]string{"TaskGroupDstNode", string(data)}, " ")
}
