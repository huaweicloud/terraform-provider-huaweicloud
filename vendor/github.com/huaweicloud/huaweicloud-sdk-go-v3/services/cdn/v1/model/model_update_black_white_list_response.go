package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateBlackWhiteListResponse Response Object
type UpdateBlackWhiteListResponse struct {

	// 响应码，200：成功，400，失败。
	Code *string `json:"code,omitempty"`

	// 响应结果。
	Result *string `json:"result,omitempty"`

	// 响应体返回内容。
	Data *interface{} `json:"data,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateBlackWhiteListResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBlackWhiteListResponse struct{}"
	}

	return strings.Join([]string{"UpdateBlackWhiteListResponse", string(data)}, " ")
}
