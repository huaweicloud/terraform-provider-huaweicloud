package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowUrlTaskInfoResponse struct {

	// 查询结果总数
	Total *int32 `json:"total,omitempty"`

	// 查询当前页总数
	Count *int32 `json:"count,omitempty"`

	// url信息
	Results *[]Urls `json:"results,omitempty"`

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
