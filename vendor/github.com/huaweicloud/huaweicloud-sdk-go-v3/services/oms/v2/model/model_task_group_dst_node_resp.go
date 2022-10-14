package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 迁移任务组目的端节点信息
type TaskGroupDstNodeResp struct {

	// 目的端桶的名称。
	Bucket *string `json:"bucket,omitempty"`

	// 目的端桶所处的区域。
	Region *string `json:"region,omitempty"`

	// 目的端桶内路径前缀（拼接在对象key前面,组成新的key,拼接后不能超过1024个字符）。
	SavePrefix *string `json:"save_prefix,omitempty"`
}

func (o TaskGroupDstNodeResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskGroupDstNodeResp struct{}"
	}

	return strings.Join([]string{"TaskGroupDstNodeResp", string(data)}, " ")
}
