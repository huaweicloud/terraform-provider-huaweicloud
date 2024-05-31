package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowUrlTaskInfoResponse Response Object
type ShowUrlTaskInfoResponse struct {

	// 查询结果总数。
	Total *int32 `json:"total,omitempty"`

	// 当前页查询到的总数。
	Count *int32 `json:"count,omitempty"`

	// url信息。
	Result *[]Urls `json:"result,omitempty"`

	XRequestId     *string `json:"X-request-id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowUrlTaskInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUrlTaskInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowUrlTaskInfoResponse", string(data)}, " ")
}
