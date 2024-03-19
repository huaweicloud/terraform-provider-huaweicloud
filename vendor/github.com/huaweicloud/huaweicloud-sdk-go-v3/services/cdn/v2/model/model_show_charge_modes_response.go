package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowChargeModesResponse Response Object
type ShowChargeModesResponse struct {

	// 计费模式查询结果
	Result         *[]map[string]interface{} `json:"result,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o ShowChargeModesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowChargeModesResponse struct{}"
	}

	return strings.Join([]string{"ShowChargeModesResponse", string(data)}, " ")
}
