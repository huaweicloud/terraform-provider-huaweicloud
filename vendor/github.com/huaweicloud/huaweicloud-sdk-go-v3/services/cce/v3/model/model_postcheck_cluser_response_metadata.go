package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostcheckCluserResponseMetadata 升级后确认元数据
type PostcheckCluserResponseMetadata struct {

	// 任务ID
	Uid *string `json:"uid,omitempty"`
}

func (o PostcheckCluserResponseMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostcheckCluserResponseMetadata struct{}"
	}

	return strings.Join([]string{"PostcheckCluserResponseMetadata", string(data)}, " ")
}
