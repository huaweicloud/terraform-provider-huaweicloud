package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowInstanceExtendProductInfoResponse struct {

	// 表示按需付费的产品列表。
	Hourly *[]ShowInstanceExtendProductInfoRespHourly `json:"hourly,omitempty"`

	// 表示包年包月的产品列表。当前暂不支持通过API创建包年包月的Kafka实例。
	Monthly        *[]ShowInstanceExtendProductInfoRespHourly `json:"monthly,omitempty"`
	HttpStatusCode int                                        `json:"-"`
}

func (o ShowInstanceExtendProductInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceExtendProductInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowInstanceExtendProductInfoResponse", string(data)}, " ")
}
