package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 自动创建虚拟机模板
type TemplateRequest struct {
	// 模板名称

	Name string `json:"name"`
	// 是否是通用模板，如果模板关联一个任务，则不算通用模板

	IsTemplate bool `json:"is_template"`
	// Region信息

	Region string `json:"region"`
	// 项目ID

	Projectid string `json:"projectid"`
	// 目标端服务器名称

	TargetServerName *string `json:"target_server_name,omitempty"`
	// 可用区

	AvailabilityZone *string `json:"availability_zone,omitempty"`
	// 磁盘类型

	Volumetype *TemplateRequestVolumetype `json:"volumetype,omitempty"`
	// 虚拟机规格

	Flavor *string `json:"flavor,omitempty"`

	Vpc *VpcObject `json:"vpc,omitempty"`
	// 网卡信息，支持多个网卡，如果是自动创建，只填一个，id使用“autoCreate”

	Nics *[]Nics `json:"nics,omitempty"`
	// 安全组，支持多个安全组，如果是自动创建，只填一个，id使用“autoCreate”

	SecurityGroups *[]SgObject `json:"security_groups,omitempty"`

	Publicip *PublicIp `json:"publicip,omitempty"`
	// 磁盘信息

	Disk *[]TemplateDisk `json:"disk,omitempty"`
	// 数据盘磁盘类型

	DataVolumeType *TemplateRequestDataVolumeType `json:"data_volume_type,omitempty"`
	// 目的端密码

	TargetPassword *string `json:"target_password,omitempty"`
}

func (o TemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TemplateRequest struct{}"
	}

	return strings.Join([]string{"TemplateRequest", string(data)}, " ")
}

type TemplateRequestVolumetype struct {
	value string
}

type TemplateRequestVolumetypeEnum struct {
	SAS  TemplateRequestVolumetype
	SSD  TemplateRequestVolumetype
	SATA TemplateRequestVolumetype
}

func GetTemplateRequestVolumetypeEnum() TemplateRequestVolumetypeEnum {
	return TemplateRequestVolumetypeEnum{
		SAS: TemplateRequestVolumetype{
			value: "SAS",
		},
		SSD: TemplateRequestVolumetype{
			value: "SSD",
		},
		SATA: TemplateRequestVolumetype{
			value: "SATA",
		},
	}
}

func (c TemplateRequestVolumetype) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TemplateRequestVolumetype) UnmarshalJSON(b []byte) error {
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

type TemplateRequestDataVolumeType struct {
	value string
}

type TemplateRequestDataVolumeTypeEnum struct {
	SAS  TemplateRequestDataVolumeType
	SSD  TemplateRequestDataVolumeType
	SATA TemplateRequestDataVolumeType
}

func GetTemplateRequestDataVolumeTypeEnum() TemplateRequestDataVolumeTypeEnum {
	return TemplateRequestDataVolumeTypeEnum{
		SAS: TemplateRequestDataVolumeType{
			value: "SAS",
		},
		SSD: TemplateRequestDataVolumeType{
			value: "SSD",
		},
		SATA: TemplateRequestDataVolumeType{
			value: "SATA",
		},
	}
}

func (c TemplateRequestDataVolumeType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TemplateRequestDataVolumeType) UnmarshalJSON(b []byte) error {
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
