package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SnapshotTaskStatus struct {

	// 最近一次备份的时间
	LatestBackupTime *string `json:"latestBackupTime,omitempty"`
}

func (o SnapshotTaskStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SnapshotTaskStatus struct{}"
	}

	return strings.Join([]string{"SnapshotTaskStatus", string(data)}, " ")
}
