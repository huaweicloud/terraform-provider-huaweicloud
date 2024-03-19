package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type GroupInfoSimple struct {

	// 创建时间。
	CreatedAt *int64 `json:"createdAt,omitempty"`

	// 消费组ID。
	GroupId *string `json:"group_id,omitempty"`

	// 消费组状态。包含以下状态： - Dead：消费组内没有任何成员，且没有任何元数据。 - Empty：消费组内没有任何成员，存在元数据。 - PreparingRebalance：准备开启rebalance。 - CompletingRebalance：所有成员加入group。 - Stable：消费组内成员可正常消费。
	State *string `json:"state,omitempty"`

	// 协调器编号。
	CoordinatorId *int32 `json:"coordinator_id,omitempty"`

	// 消费组描述。
	GroupDesc *string `json:"group_desc,omitempty"`

	// 堆积数。
	Lag *int64 `json:"lag,omitempty"`
}

func (o GroupInfoSimple) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GroupInfoSimple struct{}"
	}

	return strings.Join([]string{"GroupInfoSimple", string(data)}, " ")
}
