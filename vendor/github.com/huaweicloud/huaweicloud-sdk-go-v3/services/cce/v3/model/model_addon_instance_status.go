package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// AddonInstanceStatus 插件状态信息
type AddonInstanceStatus struct {

	// 插件实例状态, 取值如下 - running：运行中，表示插件全部实例状态都在运行中，插件正常使用。 - abnormal：不可用，表示插件状态异常，插件不可使用。可单击插件名称查看实例异常事件。 - installing：安装中，表示插件正在安装中。 - installFailed：安装失败，表示插件安装失败，需要卸载后重新安装。 - upgrading：升级中，表示插件正在更新中。 - upgradeFailed：升级失败，表示插件升级失败，可重试升级或卸载后重新安装。 - deleting：删除中，表示插件正在删除中。 - deleteFailed：删除失败，表示插件删除失败，可重试卸载。 - deleteSuccess：删除成功，表示插件删除成功。 - available：部分就绪，表示插件下只有部分实例状态为运行中，插件部分功能可用。 - rollbacking：回滚中，表示插件正在回滚中。 - rollbackFailed：回滚失败，表示插件回滚失败，可重试回滚或卸载后重新安装。 - unknown：未知状态，表示插件模板实例不存在。
	Status AddonInstanceStatusStatus `json:"status"`

	// 插件安装失败原因
	Reason string `json:"Reason"`

	// 安装错误详情
	Message string `json:"message"`

	// 此插件版本，支持升级的集群版本
	TargetVersions *[]string `json:"targetVersions,omitempty"`

	CurrentVersion *Versions `json:"currentVersion"`

	// 是否支持回滚到插件升级前的插件版本
	IsRollbackable *bool `json:"isRollbackable,omitempty"`

	// 插件升级或回滚前的版本
	PreviousVersion *string `json:"previousVersion,omitempty"`
}

func (o AddonInstanceStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddonInstanceStatus struct{}"
	}

	return strings.Join([]string{"AddonInstanceStatus", string(data)}, " ")
}

type AddonInstanceStatusStatus struct {
	value string
}

type AddonInstanceStatusStatusEnum struct {
	RUNNING         AddonInstanceStatusStatus
	ABNORMAL        AddonInstanceStatusStatus
	INSTALLING      AddonInstanceStatusStatus
	INSTALL_FAILED  AddonInstanceStatusStatus
	UPGRADING       AddonInstanceStatusStatus
	UPGRADE_FAILED  AddonInstanceStatusStatus
	DELETING        AddonInstanceStatusStatus
	DELETE_SUCCESS  AddonInstanceStatusStatus
	DELETE_FAILED   AddonInstanceStatusStatus
	AVAILABLE       AddonInstanceStatusStatus
	ROLLBACKING     AddonInstanceStatusStatus
	ROLLBACK_FAILED AddonInstanceStatusStatus
}

func GetAddonInstanceStatusStatusEnum() AddonInstanceStatusStatusEnum {
	return AddonInstanceStatusStatusEnum{
		RUNNING: AddonInstanceStatusStatus{
			value: "running",
		},
		ABNORMAL: AddonInstanceStatusStatus{
			value: "abnormal",
		},
		INSTALLING: AddonInstanceStatusStatus{
			value: "installing",
		},
		INSTALL_FAILED: AddonInstanceStatusStatus{
			value: "installFailed",
		},
		UPGRADING: AddonInstanceStatusStatus{
			value: "upgrading",
		},
		UPGRADE_FAILED: AddonInstanceStatusStatus{
			value: "upgradeFailed",
		},
		DELETING: AddonInstanceStatusStatus{
			value: "deleting",
		},
		DELETE_SUCCESS: AddonInstanceStatusStatus{
			value: "deleteSuccess",
		},
		DELETE_FAILED: AddonInstanceStatusStatus{
			value: "deleteFailed",
		},
		AVAILABLE: AddonInstanceStatusStatus{
			value: "available",
		},
		ROLLBACKING: AddonInstanceStatusStatus{
			value: "rollbacking",
		},
		ROLLBACK_FAILED: AddonInstanceStatusStatus{
			value: "rollbackFailed",
		},
	}
}

func (c AddonInstanceStatusStatus) Value() string {
	return c.value
}

func (c AddonInstanceStatusStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AddonInstanceStatusStatus) UnmarshalJSON(b []byte) error {
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
