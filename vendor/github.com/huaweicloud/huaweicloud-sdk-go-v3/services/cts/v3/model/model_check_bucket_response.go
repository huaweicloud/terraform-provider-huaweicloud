package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CheckBucketResponse struct {

	// 错误码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 错误信息。
	ErrorMessage *string `json:"error_message,omitempty"`

	// 返回的http状态码。
	ResponseCode *int32 `json:"response_code,omitempty"`

	// 是否成功转储。
	Success *bool `json:"success,omitempty"`
}

func (o CheckBucketResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckBucketResponse struct{}"
	}

	return strings.Join([]string{"CheckBucketResponse", string(data)}, " ")
}
