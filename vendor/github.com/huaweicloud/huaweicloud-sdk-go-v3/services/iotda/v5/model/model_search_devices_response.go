package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SearchDevicesResponse Response Object
type SearchDevicesResponse struct {

	// 搜索设备结果列表。
	Devices *[]SearchDevice `json:"devices,omitempty"`

	// 满足查询条件的记录总数(只有条件为select count(*)/count(1)时单独返回)。
	Count          *int64 `json:"count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o SearchDevicesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SearchDevicesResponse struct{}"
	}

	return strings.Join([]string{"SearchDevicesResponse", string(data)}, " ")
}
