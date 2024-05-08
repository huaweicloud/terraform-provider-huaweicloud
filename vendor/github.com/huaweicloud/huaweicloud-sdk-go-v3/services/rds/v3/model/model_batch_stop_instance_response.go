package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchStopInstanceResponse Response Object
type BatchStopInstanceResponse struct {

	// 停止实例结果列表
	Records *[]ShutdownInstanceRecord `json:"records,omitempty"`

	XRequestId     *string `json:"X-request-id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o BatchStopInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchStopInstanceResponse struct{}"
	}

	return strings.Join([]string{"BatchStopInstanceResponse", string(data)}, " ")
}
