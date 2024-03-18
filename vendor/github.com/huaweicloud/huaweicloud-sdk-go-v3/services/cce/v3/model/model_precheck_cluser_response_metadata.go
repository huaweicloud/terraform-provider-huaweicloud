package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PrecheckCluserResponseMetadata 升级前检查元数据
type PrecheckCluserResponseMetadata struct {

	// 检查任务ID
	Uid *string `json:"uid,omitempty"`
}

func (o PrecheckCluserResponseMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrecheckCluserResponseMetadata struct{}"
	}

	return strings.Join([]string{"PrecheckCluserResponseMetadata", string(data)}, " ")
}
