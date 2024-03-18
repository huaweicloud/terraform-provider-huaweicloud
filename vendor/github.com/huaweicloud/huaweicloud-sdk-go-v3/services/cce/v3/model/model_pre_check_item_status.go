package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PreCheckItemStatus 检查项状态信息
type PreCheckItemStatus struct {

	// 检查项名称
	Name *string `json:"name,omitempty"`

	// 检查项类型，取值如下 - Exception: 异常类，需要用户解决 - Risk：风险类，用户确认后可选择跳过
	Kind *string `json:"kind,omitempty"`

	// 检查项分组，取值如下 - LimitCheck: 集群限制检查 - MasterCheck：控制节点检查 - NodeCheck：用户节点检查 - AddonCheck：插件检查 - ExecuteException：检查流程错误
	Group *string `json:"group,omitempty"`

	// 检查项风险级别，取值如下 - Info: 提示级别 - Warning：风险级别 - Fatal：严重级别
	Level *string `json:"level,omitempty"`

	// 状态，取值如下 - Init: 初始化 - Running 运行中 - Success 成功 - Failed 失败
	Phase *string `json:"phase,omitempty"`

	// 提示信息
	Message *string `json:"message,omitempty"`

	RiskSource *RiskSource `json:"riskSource,omitempty"`

	// 错误码集合
	ErrorCodes *[]string `json:"errorCodes,omitempty"`
}

func (o PreCheckItemStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PreCheckItemStatus struct{}"
	}

	return strings.Join([]string{"PreCheckItemStatus", string(data)}, " ")
}
