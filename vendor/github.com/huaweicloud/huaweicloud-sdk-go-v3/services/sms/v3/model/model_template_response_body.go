package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 自动创建虚拟机模板
type TemplateResponseBody struct {
	// 模板ID

	Id *string `json:"id,omitempty"`
	// 模板名称

	Name string `json:"name"`
	// 是否是通用模板，如果模板关联一个任务，则不算通用模板

	IsTemplate *string `json:"is_template,omitempty"`
	// Region信息

	Region string `json:"region"`
	// 项目ID

	Projectid string `json:"projectid"`
	// 目标端服务器名称

	TargetServerName string `json:"target_server_name"`
	// 可用区

	AvailabilityZone string `json:"availability_zone"`
	// 磁盘类型

	Volumetype TemplateResponseBodyVolumetype `json:"volumetype"`
	// 虚拟机规格

	Flavor string `json:"flavor"`

	Vpc *VpcObject `json:"vpc"`
	// 网卡信息，支持多个网卡，如果是自动创建，只填一个，id使用“autoCreate”

	Nics []Nics `json:"nics"`
	// 安全组，支持多个安全组，如果是自动创建，只填一个，id使用“autoCreate”

	SecurityGroups []SgObject `json:"security_groups"`

	Publicip *PublicIp `json:"publicip"`
	// 磁盘信息

	Disk []TemplateDisk `json:"disk"`
	// 数据盘磁盘类型

	DataVolumeType TemplateResponseBodyDataVolumeType `json:"data_volume_type"`
	// 目的端密码

	TargetPassword string `json:"target_password"`
}

func (o TemplateResponseBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TemplateResponseBody struct{}"
	}

	return strings.Join([]string{"TemplateResponseBody", string(data)}, " ")
}

type TemplateResponseBodyVolumetype struct {
	value string
}

type TemplateResponseBodyVolumetypeEnum struct {
	SAS  TemplateResponseBodyVolumetype
	SSD  TemplateResponseBodyVolumetype
	SATA TemplateResponseBodyVolumetype
}

func GetTemplateResponseBodyVolumetypeEnum() TemplateResponseBodyVolumetypeEnum {
	return TemplateResponseBodyVolumetypeEnum{
		SAS: TemplateResponseBodyVolumetype{
			value: "SAS",
		},
		SSD: TemplateResponseBodyVolumetype{
			value: "SSD",
		},
		SATA: TemplateResponseBodyVolumetype{
			value: "SATA",
		},
	}
}

func (c TemplateResponseBodyVolumetype) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TemplateResponseBodyVolumetype) UnmarshalJSON(b []byte) error {
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

type TemplateResponseBodyDataVolumeType struct {
	value string
}

type TemplateResponseBodyDataVolumeTypeEnum struct {
	SAS  TemplateResponseBodyDataVolumeType
	SSD  TemplateResponseBodyDataVolumeType
	SATA TemplateResponseBodyDataVolumeType
}

func GetTemplateResponseBodyDataVolumeTypeEnum() TemplateResponseBodyDataVolumeTypeEnum {
	return TemplateResponseBodyDataVolumeTypeEnum{
		SAS: TemplateResponseBodyDataVolumeType{
			value: "SAS",
		},
		SSD: TemplateResponseBodyDataVolumeType{
			value: "SSD",
		},
		SATA: TemplateResponseBodyDataVolumeType{
			value: "SATA",
		},
	}
}

func (c TemplateResponseBodyDataVolumeType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TemplateResponseBodyDataVolumeType) UnmarshalJSON(b []byte) error {
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
