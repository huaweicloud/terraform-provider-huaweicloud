package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Response Object
type ShowMigprojectResponse struct {
	// 迁移项目名称

	Name *string `json:"name,omitempty"`
	// 迁移项目描述

	Description *string `json:"description,omitempty"`
	// 是否为默认模板

	Isdefault *bool `json:"isdefault,omitempty"`
	// 区域名称

	Region *string `json:"region,omitempty"`
	// 迁移后是否启动目的端虚拟机

	StartTargetServer *bool `json:"start_target_server,omitempty"`
	// 限制迁移速率，单位：Mbps

	SpeedLimit *int32 `json:"speed_limit,omitempty"`
	// 是否使用公网IP迁移

	UsePublicIp *bool `json:"use_public_ip,omitempty"`
	// 是否是已经存在的服务器

	ExistServer *bool `json:"exist_server,omitempty"`
	// 迁移项目类型

	Type *ShowMigprojectResponseType `json:"type,omitempty"`
	// 企业项目名称

	EnterpriseProject *string `json:"enterprise_project,omitempty"`
	// 首次复制或者同步后 是否继续持续同步

	Syncing        *bool `json:"syncing,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ShowMigprojectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMigprojectResponse struct{}"
	}

	return strings.Join([]string{"ShowMigprojectResponse", string(data)}, " ")
}

type ShowMigprojectResponseType struct {
	value string
}

type ShowMigprojectResponseTypeEnum struct {
	MIGRATE_BLOCK ShowMigprojectResponseType
	MIGRATE_FILE  ShowMigprojectResponseType
}

func GetShowMigprojectResponseTypeEnum() ShowMigprojectResponseTypeEnum {
	return ShowMigprojectResponseTypeEnum{
		MIGRATE_BLOCK: ShowMigprojectResponseType{
			value: "MIGRATE_BLOCK",
		},
		MIGRATE_FILE: ShowMigprojectResponseType{
			value: "MIGRATE_FILE",
		},
	}
}

func (c ShowMigprojectResponseType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowMigprojectResponseType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
