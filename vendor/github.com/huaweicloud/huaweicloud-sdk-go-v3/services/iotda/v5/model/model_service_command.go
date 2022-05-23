package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 命令服务对象。
type ServiceCommand struct {

	// **参数说明**：设备命令名称。注：设备服务内不允许重复。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	CommandName string `json:"command_name"`

	// **参数说明**：设备命令的参数列表。
	Paras *[]ServiceCommandPara `json:"paras,omitempty"`

	// **参数说明**：设备命令的响应列表。
	Responses *[]ServiceCommandResponse `json:"responses,omitempty"`
}

func (o ServiceCommand) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServiceCommand struct{}"
	}

	return strings.Join([]string{"ServiceCommand", string(data)}, " ")
}
