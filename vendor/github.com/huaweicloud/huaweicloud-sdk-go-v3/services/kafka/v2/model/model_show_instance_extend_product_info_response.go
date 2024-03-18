package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowInstanceExtendProductInfoResponse Response Object
type ShowInstanceExtendProductInfoResponse struct {

	// 表示[按需付费的](tag:hws,hws_hk,hws_ocb,ocb,ctc,sbc,hk_sbc,cmcc,g42,tm,hk_g42,hk_tm)产品列表。
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
