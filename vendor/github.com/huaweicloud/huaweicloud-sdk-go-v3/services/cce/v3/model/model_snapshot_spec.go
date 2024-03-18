package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SnapshotSpec struct {

	// 备份任务详情
	Items *[]SnapshotSpecItems `json:"items,omitempty"`
}

func (o SnapshotSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SnapshotSpec struct{}"
	}

	return strings.Join([]string{"SnapshotSpec", string(data)}, " ")
}
