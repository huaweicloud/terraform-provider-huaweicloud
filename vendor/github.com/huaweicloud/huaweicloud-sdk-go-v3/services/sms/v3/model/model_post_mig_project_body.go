package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SMS迁移项目信息
type PostMigProjectBody struct {
	// 迁移项目名称

	Name string `json:"name"`
	// 迁移项目描述

	Description *string `json:"description,omitempty"`
	// 是否为默认模板

	Isdefault *bool `json:"isdefault,omitempty"`
	// 区域名称

	Region string `json:"region"`
	// 迁移后是否启动目的端虚拟机

	StartTargetServer *bool `json:"start_target_server,omitempty"`
	// 限制迁移速率，单位：Mbps

	SpeedLimit *int32 `json:"speed_limit,omitempty"`
	// 是否使用公网IP迁移

	UsePublicIp bool `json:"use_public_ip"`
	// 是否是已经存在的服务器

	ExistServer bool `json:"exist_server"`
	// 迁移项目类型

	Type PostMigProjectBodyType `json:"type"`
	// 企业项目名称

	EnterpriseProject *string `json:"enterprise_project,omitempty"`
	// 首次复制或者同步后 是否继续持续同步

	Syncing bool `json:"syncing"`
}

func (o PostMigProjectBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostMigProjectBody struct{}"
	}

	return strings.Join([]string{"PostMigProjectBody", string(data)}, " ")
}

type PostMigProjectBodyType struct {
	value string
}

type PostMigProjectBodyTypeEnum struct {
	MIGRATE_BLOCK PostMigProjectBodyType
	MIGRATE_FILE  PostMigProjectBodyType
}

func GetPostMigProjectBodyTypeEnum() PostMigProjectBodyTypeEnum {
	return PostMigProjectBodyTypeEnum{
		MIGRATE_BLOCK: PostMigProjectBodyType{
			value: "MIGRATE_BLOCK",
		},
		MIGRATE_FILE: PostMigProjectBodyType{
			value: "MIGRATE_FILE",
		},
	}
}

func (c PostMigProjectBodyType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PostMigProjectBodyType) UnmarshalJSON(b []byte) error {
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
