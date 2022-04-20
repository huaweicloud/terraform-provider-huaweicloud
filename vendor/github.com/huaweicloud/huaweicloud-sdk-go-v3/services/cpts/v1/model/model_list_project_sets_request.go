package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListProjectSetsRequest struct {
	// 查询偏移

	Offset *int32 `json:"offset,omitempty"`
	// 查询数量

	Limit *int32 `json:"limit,omitempty"`
}

func (o ListProjectSetsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProjectSetsRequest struct{}"
	}

	return strings.Join([]string{"ListProjectSetsRequest", string(data)}, " ")
}
