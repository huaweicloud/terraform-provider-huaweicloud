package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowTopUrlResponse struct {

	// 服务区域。
	ServiceArea *string `json:"service_area,omitempty"`

	// 详情数据对象。
	TopUrlSummary  *[]TopUrlSummary `json:"top_url_summary,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ShowTopUrlResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTopUrlResponse struct{}"
	}

	return strings.Join([]string{"ShowTopUrlResponse", string(data)}, " ")
}
