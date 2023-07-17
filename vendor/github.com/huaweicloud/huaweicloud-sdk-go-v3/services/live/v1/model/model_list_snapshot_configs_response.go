package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSnapshotConfigsResponse Response Object
type ListSnapshotConfigsResponse struct {

	// 总条目数
	Total *int32 `json:"total,omitempty"`

	SnapshotConfigList *LiveSnapshotConfig `json:"snapshot_config_list,omitempty"`

	// 每页记录数
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量
	Offset         *int32 `json:"offset,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListSnapshotConfigsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSnapshotConfigsResponse struct{}"
	}

	return strings.Join([]string{"ListSnapshotConfigsResponse", string(data)}, " ")
}
