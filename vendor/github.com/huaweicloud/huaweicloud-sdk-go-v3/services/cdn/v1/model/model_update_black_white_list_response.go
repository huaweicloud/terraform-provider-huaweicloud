package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateBlackWhiteListResponse struct {

	// 编码
	Code *string `json:"code,omitempty"`

	// 结果
	Result *string `json:"result,omitempty"`

	// 数据
	Data           *interface{} `json:"data,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o UpdateBlackWhiteListResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBlackWhiteListResponse struct{}"
	}

	return strings.Join([]string{"UpdateBlackWhiteListResponse", string(data)}, " ")
}
