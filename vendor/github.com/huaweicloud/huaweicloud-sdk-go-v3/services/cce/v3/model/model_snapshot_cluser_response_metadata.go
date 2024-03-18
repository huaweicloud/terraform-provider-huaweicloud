package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SnapshotCluserResponseMetadata 备份任务数据
type SnapshotCluserResponseMetadata struct {

	// API版本，默认为v3.1
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 任务类型
	Kind *string `json:"kind,omitempty"`
}

func (o SnapshotCluserResponseMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SnapshotCluserResponseMetadata struct{}"
	}

	return strings.Join([]string{"SnapshotCluserResponseMetadata", string(data)}, " ")
}
