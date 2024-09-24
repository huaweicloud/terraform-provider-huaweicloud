package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ScaleGroupStatus 伸缩组当前详细状态信息
type ScaleGroupStatus struct {

	// 伸缩组名称
	Name *string `json:"name,omitempty"`

	// 伸缩组uuid
	Uid *string `json:"uid,omitempty"`

	// 伸缩组创建时间
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`

	// 伸缩组更新时间
	UpdateTimestamp *string `json:"updateTimestamp,omitempty"`

	// 伸缩组状态。 - 空值：可用（伸缩组当前节点数已达到预期，且无伸缩中的节点） - Synchronizing：伸缩中（伸缩组当前节点数未达到预期，且无伸缩中的节点） - Synchronized：伸缩等待中（伸缩组当前节点数未达到预期，或者存在伸缩中的节点） - SoldOut：伸缩组当前不可扩容（兼容字段，标记节点池资源售罄、资源配额不足等不可扩容状态） > 上述伸缩组状态已废弃，仅兼容保留，不建议使用，替代感知方式如下： > - 伸缩组扩缩状态：可通过desiredNodeCount/existingNodeCount/upcomingNodeCount节点状态统计信息，精确感知当前伸缩组扩缩状态。 > - 伸缩组可扩容状态：可通过conditions感知伸缩组详细状态，其中\"Scalable\"可替代SoldOut语义。 - Deleting：删除中 - Error：错误
	Phase *ScaleGroupStatusPhase `json:"phase,omitempty"`

	// 伸缩组期望节点数
	DesiredNodeCount *int32 `json:"desiredNodeCount,omitempty"`

	// 订单未支付节点个数
	UnpaidScaleNodeCount *int32 `json:"unpaidScaleNodeCount,omitempty"`

	ExistingNodeCount *ScaleGroupStatusExistingNodeCount `json:"existingNodeCount,omitempty"`

	UpcomingNodeCount *ScaleGroupStatusUpcomingNodeCount `json:"upcomingNodeCount,omitempty"`

	// 伸缩组禁止缩容的节点数
	ScaleDownDisabledNodeCount *int32 `json:"scaleDownDisabledNodeCount,omitempty"`

	// 伸缩组当前详细状态列表，详情参见Condition类型定义。
	Conditions *[]NodePoolCondition `json:"conditions,omitempty"`
}

func (o ScaleGroupStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleGroupStatus struct{}"
	}

	return strings.Join([]string{"ScaleGroupStatus", string(data)}, " ")
}

type ScaleGroupStatusPhase struct {
	value string
}

type ScaleGroupStatusPhaseEnum struct {
	SYNCHRONIZING ScaleGroupStatusPhase
	SYNCHRONIZED  ScaleGroupStatusPhase
	SOLD_OUT      ScaleGroupStatusPhase
	DELETING      ScaleGroupStatusPhase
	ERROR         ScaleGroupStatusPhase
}

func GetScaleGroupStatusPhaseEnum() ScaleGroupStatusPhaseEnum {
	return ScaleGroupStatusPhaseEnum{
		SYNCHRONIZING: ScaleGroupStatusPhase{
			value: "Synchronizing",
		},
		SYNCHRONIZED: ScaleGroupStatusPhase{
			value: "Synchronized",
		},
		SOLD_OUT: ScaleGroupStatusPhase{
			value: "SoldOut",
		},
		DELETING: ScaleGroupStatusPhase{
			value: "Deleting",
		},
		ERROR: ScaleGroupStatusPhase{
			value: "Error",
		},
	}
}

func (c ScaleGroupStatusPhase) Value() string {
	return c.value
}

func (c ScaleGroupStatusPhase) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ScaleGroupStatusPhase) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
