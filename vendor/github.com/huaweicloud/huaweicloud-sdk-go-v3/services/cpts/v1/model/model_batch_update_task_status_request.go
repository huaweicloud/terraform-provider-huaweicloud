package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchUpdateTaskStatusRequest Request Object
type BatchUpdateTaskStatusRequest struct {

	// 工程id
	TestSuitId int32 `json:"test_suit_id"`

	Body *BatchUpdateTaskStatusRequestBody `json:"body,omitempty"`
}

func (o BatchUpdateTaskStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchUpdateTaskStatusRequest struct{}"
	}

	return strings.Join([]string{"BatchUpdateTaskStatusRequest", string(data)}, " ")
}
