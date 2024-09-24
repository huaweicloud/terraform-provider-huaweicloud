package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DeleteAutopilotClusterRequest Request Object
type DeleteAutopilotClusterRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 集群状态兼容Error参数，用于API平滑切换。 兼容场景下，errorStatus为空则屏蔽Error状态为Deleting状态。
	ErrorStatus *string `json:"errorStatus,omitempty"`

	// 是否删除SFS Turbo（极速文件存储卷）， 枚举取值： - true或block (执行删除流程，失败则阻塞后续流程) - try (执行删除流程，失败则忽略，并继续执行后续流程) - false或skip (跳过删除流程，默认选项)
	DeleteEfs *DeleteAutopilotClusterRequestDeleteEfs `json:"delete_efs,omitempty"`

	// 是否删除eni ports（原生弹性网卡）， 枚举取值： - true或block (执行删除流程，失败则阻塞后续流程，默认选项) - try (执行删除流程，失败则忽略，并继续执行后续流程) - false或skip (跳过删除流程)
	DeleteEni *DeleteAutopilotClusterRequestDeleteEni `json:"delete_eni,omitempty"`

	// 是否删除elb（弹性负载均衡）等集群Service/Ingress相关资源。 枚举取值： - true或block (执行删除流程，失败则阻塞后续流程，默认选项) - try (执行删除流程，失败则忽略，并继续执行后续流程) - false或skip (跳过删除流程)
	DeleteNet *DeleteAutopilotClusterRequestDeleteNet `json:"delete_net,omitempty"`

	// 是否删除obs（对象存储卷）， 枚举取值： - true或block (执行删除流程，失败则阻塞后续流程) - try (执行删除流程，失败则忽略，并继续执行后续流程) - false或skip (跳过删除流程，默认选项)
	DeleteObs *DeleteAutopilotClusterRequestDeleteObs `json:"delete_obs,omitempty"`

	// 是否删除sfs3.0（文件存储卷3.0）， 枚举取值： - true或block (执行删除流程，失败则阻塞后续流程) - try (执行删除流程，失败则忽略，并继续执行后续流程) - false或skip (跳过删除流程，默认选项)
	DeleteSfs30 *DeleteAutopilotClusterRequestDeleteSfs30 `json:"delete_sfs30,omitempty"`

	// 是否删除LTS LogStream（日志流）， 枚举取值： - true或block (执行删除流程，失败则阻塞后续流程) - try (执行删除流程，失败则忽略，并继续执行后续流程) - false或skip (跳过删除流程，默认选项)
	LtsReclaimPolicy *DeleteAutopilotClusterRequestLtsReclaimPolicy `json:"lts_reclaim_policy,omitempty"`
}

func (o DeleteAutopilotClusterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAutopilotClusterRequest struct{}"
	}

	return strings.Join([]string{"DeleteAutopilotClusterRequest", string(data)}, " ")
}

type DeleteAutopilotClusterRequestDeleteEfs struct {
	value string
}

type DeleteAutopilotClusterRequestDeleteEfsEnum struct {
	TRUE  DeleteAutopilotClusterRequestDeleteEfs
	BLOCK DeleteAutopilotClusterRequestDeleteEfs
	TRY   DeleteAutopilotClusterRequestDeleteEfs
	FALSE DeleteAutopilotClusterRequestDeleteEfs
	SKIP  DeleteAutopilotClusterRequestDeleteEfs
}

func GetDeleteAutopilotClusterRequestDeleteEfsEnum() DeleteAutopilotClusterRequestDeleteEfsEnum {
	return DeleteAutopilotClusterRequestDeleteEfsEnum{
		TRUE: DeleteAutopilotClusterRequestDeleteEfs{
			value: "true",
		},
		BLOCK: DeleteAutopilotClusterRequestDeleteEfs{
			value: "block",
		},
		TRY: DeleteAutopilotClusterRequestDeleteEfs{
			value: "try",
		},
		FALSE: DeleteAutopilotClusterRequestDeleteEfs{
			value: "false",
		},
		SKIP: DeleteAutopilotClusterRequestDeleteEfs{
			value: "skip",
		},
	}
}

func (c DeleteAutopilotClusterRequestDeleteEfs) Value() string {
	return c.value
}

func (c DeleteAutopilotClusterRequestDeleteEfs) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteAutopilotClusterRequestDeleteEfs) UnmarshalJSON(b []byte) error {
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

type DeleteAutopilotClusterRequestDeleteEni struct {
	value string
}

type DeleteAutopilotClusterRequestDeleteEniEnum struct {
	TRUE  DeleteAutopilotClusterRequestDeleteEni
	BLOCK DeleteAutopilotClusterRequestDeleteEni
	TRY   DeleteAutopilotClusterRequestDeleteEni
	FALSE DeleteAutopilotClusterRequestDeleteEni
	SKIP  DeleteAutopilotClusterRequestDeleteEni
}

func GetDeleteAutopilotClusterRequestDeleteEniEnum() DeleteAutopilotClusterRequestDeleteEniEnum {
	return DeleteAutopilotClusterRequestDeleteEniEnum{
		TRUE: DeleteAutopilotClusterRequestDeleteEni{
			value: "true",
		},
		BLOCK: DeleteAutopilotClusterRequestDeleteEni{
			value: "block",
		},
		TRY: DeleteAutopilotClusterRequestDeleteEni{
			value: "try",
		},
		FALSE: DeleteAutopilotClusterRequestDeleteEni{
			value: "false",
		},
		SKIP: DeleteAutopilotClusterRequestDeleteEni{
			value: "skip",
		},
	}
}

func (c DeleteAutopilotClusterRequestDeleteEni) Value() string {
	return c.value
}

func (c DeleteAutopilotClusterRequestDeleteEni) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteAutopilotClusterRequestDeleteEni) UnmarshalJSON(b []byte) error {
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

type DeleteAutopilotClusterRequestDeleteNet struct {
	value string
}

type DeleteAutopilotClusterRequestDeleteNetEnum struct {
	TRUE  DeleteAutopilotClusterRequestDeleteNet
	BLOCK DeleteAutopilotClusterRequestDeleteNet
	TRY   DeleteAutopilotClusterRequestDeleteNet
	FALSE DeleteAutopilotClusterRequestDeleteNet
	SKIP  DeleteAutopilotClusterRequestDeleteNet
}

func GetDeleteAutopilotClusterRequestDeleteNetEnum() DeleteAutopilotClusterRequestDeleteNetEnum {
	return DeleteAutopilotClusterRequestDeleteNetEnum{
		TRUE: DeleteAutopilotClusterRequestDeleteNet{
			value: "true",
		},
		BLOCK: DeleteAutopilotClusterRequestDeleteNet{
			value: "block",
		},
		TRY: DeleteAutopilotClusterRequestDeleteNet{
			value: "try",
		},
		FALSE: DeleteAutopilotClusterRequestDeleteNet{
			value: "false",
		},
		SKIP: DeleteAutopilotClusterRequestDeleteNet{
			value: "skip",
		},
	}
}

func (c DeleteAutopilotClusterRequestDeleteNet) Value() string {
	return c.value
}

func (c DeleteAutopilotClusterRequestDeleteNet) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteAutopilotClusterRequestDeleteNet) UnmarshalJSON(b []byte) error {
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

type DeleteAutopilotClusterRequestDeleteObs struct {
	value string
}

type DeleteAutopilotClusterRequestDeleteObsEnum struct {
	TRUE  DeleteAutopilotClusterRequestDeleteObs
	BLOCK DeleteAutopilotClusterRequestDeleteObs
	TRY   DeleteAutopilotClusterRequestDeleteObs
	FALSE DeleteAutopilotClusterRequestDeleteObs
	SKIP  DeleteAutopilotClusterRequestDeleteObs
}

func GetDeleteAutopilotClusterRequestDeleteObsEnum() DeleteAutopilotClusterRequestDeleteObsEnum {
	return DeleteAutopilotClusterRequestDeleteObsEnum{
		TRUE: DeleteAutopilotClusterRequestDeleteObs{
			value: "true",
		},
		BLOCK: DeleteAutopilotClusterRequestDeleteObs{
			value: "block",
		},
		TRY: DeleteAutopilotClusterRequestDeleteObs{
			value: "try",
		},
		FALSE: DeleteAutopilotClusterRequestDeleteObs{
			value: "false",
		},
		SKIP: DeleteAutopilotClusterRequestDeleteObs{
			value: "skip",
		},
	}
}

func (c DeleteAutopilotClusterRequestDeleteObs) Value() string {
	return c.value
}

func (c DeleteAutopilotClusterRequestDeleteObs) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteAutopilotClusterRequestDeleteObs) UnmarshalJSON(b []byte) error {
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

type DeleteAutopilotClusterRequestDeleteSfs30 struct {
	value string
}

type DeleteAutopilotClusterRequestDeleteSfs30Enum struct {
	TRUE  DeleteAutopilotClusterRequestDeleteSfs30
	BLOCK DeleteAutopilotClusterRequestDeleteSfs30
	TRY   DeleteAutopilotClusterRequestDeleteSfs30
	FALSE DeleteAutopilotClusterRequestDeleteSfs30
	SKIP  DeleteAutopilotClusterRequestDeleteSfs30
}

func GetDeleteAutopilotClusterRequestDeleteSfs30Enum() DeleteAutopilotClusterRequestDeleteSfs30Enum {
	return DeleteAutopilotClusterRequestDeleteSfs30Enum{
		TRUE: DeleteAutopilotClusterRequestDeleteSfs30{
			value: "true",
		},
		BLOCK: DeleteAutopilotClusterRequestDeleteSfs30{
			value: "block",
		},
		TRY: DeleteAutopilotClusterRequestDeleteSfs30{
			value: "try",
		},
		FALSE: DeleteAutopilotClusterRequestDeleteSfs30{
			value: "false",
		},
		SKIP: DeleteAutopilotClusterRequestDeleteSfs30{
			value: "skip",
		},
	}
}

func (c DeleteAutopilotClusterRequestDeleteSfs30) Value() string {
	return c.value
}

func (c DeleteAutopilotClusterRequestDeleteSfs30) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteAutopilotClusterRequestDeleteSfs30) UnmarshalJSON(b []byte) error {
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

type DeleteAutopilotClusterRequestLtsReclaimPolicy struct {
	value string
}

type DeleteAutopilotClusterRequestLtsReclaimPolicyEnum struct {
	TRUE  DeleteAutopilotClusterRequestLtsReclaimPolicy
	BLOCK DeleteAutopilotClusterRequestLtsReclaimPolicy
	TRY   DeleteAutopilotClusterRequestLtsReclaimPolicy
	FALSE DeleteAutopilotClusterRequestLtsReclaimPolicy
	SKIP  DeleteAutopilotClusterRequestLtsReclaimPolicy
}

func GetDeleteAutopilotClusterRequestLtsReclaimPolicyEnum() DeleteAutopilotClusterRequestLtsReclaimPolicyEnum {
	return DeleteAutopilotClusterRequestLtsReclaimPolicyEnum{
		TRUE: DeleteAutopilotClusterRequestLtsReclaimPolicy{
			value: "true",
		},
		BLOCK: DeleteAutopilotClusterRequestLtsReclaimPolicy{
			value: "block",
		},
		TRY: DeleteAutopilotClusterRequestLtsReclaimPolicy{
			value: "try",
		},
		FALSE: DeleteAutopilotClusterRequestLtsReclaimPolicy{
			value: "false",
		},
		SKIP: DeleteAutopilotClusterRequestLtsReclaimPolicy{
			value: "skip",
		},
	}
}

func (c DeleteAutopilotClusterRequestLtsReclaimPolicy) Value() string {
	return c.value
}

func (c DeleteAutopilotClusterRequestLtsReclaimPolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteAutopilotClusterRequestLtsReclaimPolicy) UnmarshalJSON(b []byte) error {
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
