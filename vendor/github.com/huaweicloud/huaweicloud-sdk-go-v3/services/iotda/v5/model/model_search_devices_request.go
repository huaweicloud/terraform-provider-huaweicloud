package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SearchDevicesRequest Request Object
type SearchDevicesRequest struct {

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。您可以在IoTDA管理控制台界面，选择左侧导航栏“总览”页签查看当前实例的ID
	InstanceId *string `json:"Instance-Id,omitempty"`

	Body *SearchSql `json:"body,omitempty"`
}

func (o SearchDevicesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SearchDevicesRequest struct{}"
	}

	return strings.Join([]string{"SearchDevicesRequest", string(data)}, " ")
}
