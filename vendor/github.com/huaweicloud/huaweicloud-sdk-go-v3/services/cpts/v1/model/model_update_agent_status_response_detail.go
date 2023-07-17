package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateAgentStatusResponseDetail struct {

	// 全链路应用id
	Id int32 `json:"id"`

	// 全链路应用状态，枚举值：CREATING，FAILED，NORMAL，DELETE
	Status string `json:"status"`

	// 全链路应用更新时间
	UpdateTime int64 `json:"update_time"`

	Result *AgentConfig `json:"result,omitempty"`
}

func (o UpdateAgentStatusResponseDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAgentStatusResponseDetail struct{}"
	}

	return strings.Join([]string{"UpdateAgentStatusResponseDetail", string(data)}, " ")
}
