package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 创建或修改规则条件时，指定资源及其范围
type RoutingRuleSubject struct {

	// **参数说明**：资源名称。 **取值范围**： - device：设备。 - device.property：设备属性。 - device.message：设备消息。 - device.message.status：设备消息状态。 - device.status：设备状态。 - batchtask：批量任务。 - product：产品。 - device.command.status：设备异步命令状态。
	Resource string `json:"resource"`

	// **参数说明**：资源事件。 **取值范围**：与资源有关，不同的资源，事件不同。event需要与resource关联使用，具体的“resource：event”映射关系如下： - device：create（设备添加） - device：delete（设备删除） - device：update（设备更新） - device.status：update （设备状态变更） - device.property：report（设备属性上报） - device.message：report（设备消息上报） - device.message.status：update（设备消息状态变更） - batchtask：update （批量任务状态变更） - product：create（产品添加） - product：delete（产品删除） - product：update（产品更新） - device.command.status：update（设备异步命令状态更新）
	Event string `json:"event"`
}

func (o RoutingRuleSubject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoutingRuleSubject struct{}"
	}

	return strings.Join([]string{"RoutingRuleSubject", string(data)}, " ")
}
