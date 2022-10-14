package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TestCaseStage struct {

	// 压力阶段模式，0：时长模式；1：次数模式
	OperateMode *int32 `json:"operate_mode,omitempty"`

	// 阶段名称
	Name *string `json:"name,omitempty"`

	// 压测时长（单位：秒）
	Time *int32 `json:"time,omitempty"`

	// 开始时间
	StartTime *int32 `json:"start_time,omitempty"`

	// 结束时间
	EndTime *int32 `json:"end_time,omitempty"`

	// 最大并发数
	IssueNum *int32 `json:"issue_num,omitempty"`

	// 次数模式发送总次数
	Count *int32 `json:"count,omitempty"`

	// 压力模式，0：并发模式；1：TPS模式；2：摸高模式；3：浪涌并发模式；4：浪涌TPS模式；5：震荡并发模式；6：震荡TPS模式；7：智能摸高模式
	PressureMode *int32 `json:"pressure_mode,omitempty"`

	// TPS模式下TPS值
	TpsValue *int32 `json:"tps_value,omitempty"`

	// 起始并发数
	CurrentUserNum *int32 `json:"current_user_num,omitempty"`

	// 起始tps值
	CurrentTps *int32 `json:"current_tps,omitempty"`

	// 调压模式，0：自动调压模式；1：手动调压模式
	VoltageRegulatingMode *int32 `json:"voltage_regulating_mode,omitempty"`

	// 浪涌/浪涌模式下最大并发数
	Maximum *int32 `json:"maximum,omitempty"`

	// 浪涌/浪涌模式下最小并发数
	Minimum *int32 `json:"minimum,omitempty"`

	// 震荡/浪涌次数
	LoopCount *int32 `json:"loop_count,omitempty"`

	// 浪涌模式下峰值持续时间
	MaxDuration *int32 `json:"max_duration,omitempty"`

	// 摸高模式下爬坡时长（单位：秒）
	RampUp *int32 `json:"ramp_up,omitempty"`

	PeakLoadKpis *StageKpiItems `json:"peak_load_kpis,omitempty"`

	// 智能摸高模式下单步执行时长
	StepDuration *int32 `json:"step_duration,omitempty"`

	// 智能摸高模式下递增并发数
	StepSize *int32 `json:"step_size,omitempty"`
}

func (o TestCaseStage) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TestCaseStage struct{}"
	}

	return strings.Join([]string{"TestCaseStage", string(data)}, " ")
}
