package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResetReplicaReqPartitions struct {

	// 分区ID。
	Partition *int32 `json:"partition,omitempty"`

	// 副本期望所在的broker ID。其中Array首位为leader副本，所有分区需要有同样数量的副本，副本数不能大于总broker的数量。
	Replicas *[]int32 `json:"replicas,omitempty"`
}

func (o ResetReplicaReqPartitions) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetReplicaReqPartitions struct{}"
	}

	return strings.Join([]string{"ResetReplicaReqPartitions", string(data)}, " ")
}
