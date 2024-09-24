package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListAutopilotClustersRequest Request Object
type ListAutopilotClustersRequest struct {

	// 集群状态兼容Error参数，用于API平滑切换。 兼容场景下，errorStatus为空则屏蔽Error状态为Deleting状态。
	ErrorStatus *string `json:"errorStatus,omitempty"`

	// 查询集群详细信息。  若设置为true，获取集群下节点总数(totalNodesNumber)、正常节点数(activeNodesNumber)、CPU总量(totalNodesCPU)、内存总量(totalNodesMemory)、已安装插件列表(installedAddonInstances)，已安装插件列表中包含名称(addonTemplateName)、版本号(version)、插件的状态信息(status)，放入到annotation中。
	Detail *string `json:"detail,omitempty"`

	// 集群状态，取值如下 - Available：可用，表示集群处于正常状态。 - Unavailable：不可用，表示集群异常，需手动删除。 - Creating：创建中，表示集群正处于创建过程中。 - Deleting：删除中，表示集群正处于删除过程中。 - Upgrading：升级中，表示集群正处于升级过程中。 - RollingBack：回滚中，表示集群正处于回滚过程中。 - RollbackFailed：回滚异常，表示集群回滚异常。 - Error：错误，表示集群资源异常，可尝试手动删除。
	Status *ListAutopilotClustersRequestStatus `json:"status,omitempty"`

	// 集群类型： - VirtualMachine：CCE集群
	Type *ListAutopilotClustersRequestType `json:"type,omitempty"`

	// 集群版本过滤
	Version *string `json:"version,omitempty"`
}

func (o ListAutopilotClustersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotClustersRequest struct{}"
	}

	return strings.Join([]string{"ListAutopilotClustersRequest", string(data)}, " ")
}

type ListAutopilotClustersRequestStatus struct {
	value string
}

type ListAutopilotClustersRequestStatusEnum struct {
	AVAILABLE       ListAutopilotClustersRequestStatus
	UNAVAILABLE     ListAutopilotClustersRequestStatus
	CREATING        ListAutopilotClustersRequestStatus
	DELETING        ListAutopilotClustersRequestStatus
	UPGRADING       ListAutopilotClustersRequestStatus
	ROLLING_BACK    ListAutopilotClustersRequestStatus
	ROLLBACK_FAILED ListAutopilotClustersRequestStatus
	ERROR           ListAutopilotClustersRequestStatus
}

func GetListAutopilotClustersRequestStatusEnum() ListAutopilotClustersRequestStatusEnum {
	return ListAutopilotClustersRequestStatusEnum{
		AVAILABLE: ListAutopilotClustersRequestStatus{
			value: "Available",
		},
		UNAVAILABLE: ListAutopilotClustersRequestStatus{
			value: "Unavailable",
		},
		CREATING: ListAutopilotClustersRequestStatus{
			value: "Creating",
		},
		DELETING: ListAutopilotClustersRequestStatus{
			value: "Deleting",
		},
		UPGRADING: ListAutopilotClustersRequestStatus{
			value: "Upgrading",
		},
		ROLLING_BACK: ListAutopilotClustersRequestStatus{
			value: "RollingBack",
		},
		ROLLBACK_FAILED: ListAutopilotClustersRequestStatus{
			value: "RollbackFailed",
		},
		ERROR: ListAutopilotClustersRequestStatus{
			value: "Error",
		},
	}
}

func (c ListAutopilotClustersRequestStatus) Value() string {
	return c.value
}

func (c ListAutopilotClustersRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListAutopilotClustersRequestStatus) UnmarshalJSON(b []byte) error {
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

type ListAutopilotClustersRequestType struct {
	value string
}

type ListAutopilotClustersRequestTypeEnum struct {
	VIRTUAL_MACHINE ListAutopilotClustersRequestType
}

func GetListAutopilotClustersRequestTypeEnum() ListAutopilotClustersRequestTypeEnum {
	return ListAutopilotClustersRequestTypeEnum{
		VIRTUAL_MACHINE: ListAutopilotClustersRequestType{
			value: "VirtualMachine",
		},
	}
}

func (c ListAutopilotClustersRequestType) Value() string {
	return c.value
}

func (c ListAutopilotClustersRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListAutopilotClustersRequestType) UnmarshalJSON(b []byte) error {
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
