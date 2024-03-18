package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SnapshotTaskMetadata struct {

	// 任务的ID。
	Uid *string `json:"uid,omitempty"`

	// 任务的创建时间。
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`

	// 任务的更新时间。
	UpdateTimestamp *string `json:"updateTimestamp,omitempty"`
}

func (o SnapshotTaskMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SnapshotTaskMetadata struct{}"
	}

	return strings.Join([]string{"SnapshotTaskMetadata", string(data)}, " ")
}
