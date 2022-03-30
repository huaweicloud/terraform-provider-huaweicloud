package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ExpandGaussMySqlInstanceVolumeResponse struct {
	// 扩容后容量。

	Size *int32 `json:"size,omitempty"`
	// 订单号。

	OrderId        *string `json:"order_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ExpandGaussMySqlInstanceVolumeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExpandGaussMySqlInstanceVolumeResponse struct{}"
	}

	return strings.Join([]string{"ExpandGaussMySqlInstanceVolumeResponse", string(data)}, " ")
}
