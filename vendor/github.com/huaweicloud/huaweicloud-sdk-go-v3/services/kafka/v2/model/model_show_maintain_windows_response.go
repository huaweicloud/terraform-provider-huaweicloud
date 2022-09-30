package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowMaintainWindowsResponse struct {

	// 支持的维护时间窗列表。
	MaintainWindows *[]MaintainWindowsEntity `json:"maintain_windows,omitempty"`
	HttpStatusCode  int                      `json:"-"`
}

func (o ShowMaintainWindowsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMaintainWindowsResponse struct{}"
	}

	return strings.Join([]string{"ShowMaintainWindowsResponse", string(data)}, " ")
}
