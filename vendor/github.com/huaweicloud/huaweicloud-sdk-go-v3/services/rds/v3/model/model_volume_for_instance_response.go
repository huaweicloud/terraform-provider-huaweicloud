package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type VolumeForInstanceResponse struct {

	// 磁盘类型。 取值范围如下，区分大小写： - COMMON，表示SATA。 - HIGH，表示SAS。 - ULTRAHIGH，表示SSD。 - ULTRAHIGHPRO，表示SSD尊享版，仅支持超高性能型尊享版。 - CLOUDSSD，表示SSD云盘，仅支持通用型和独享型规格实例。 - LOCALSSD，表示本地SSD。 - ESSD，表示极速型SSD，仅支持独享型规格实例。
	Type VolumeForInstanceResponseType `json:"type"`

	// 磁盘大小，单位为GB。
	Size int32 `json:"size"`
}

func (o VolumeForInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VolumeForInstanceResponse struct{}"
	}

	return strings.Join([]string{"VolumeForInstanceResponse", string(data)}, " ")
}

type VolumeForInstanceResponseType struct {
	value string
}

type VolumeForInstanceResponseTypeEnum struct {
	ULTRAHIGH    VolumeForInstanceResponseType
	HIGH         VolumeForInstanceResponseType
	COMMON       VolumeForInstanceResponseType
	NVMESSD      VolumeForInstanceResponseType
	ULTRAHIGHPRO VolumeForInstanceResponseType
	CLOUDSSD     VolumeForInstanceResponseType
	LOCALSSD     VolumeForInstanceResponseType
	ESSD         VolumeForInstanceResponseType
}

func GetVolumeForInstanceResponseTypeEnum() VolumeForInstanceResponseTypeEnum {
	return VolumeForInstanceResponseTypeEnum{
		ULTRAHIGH: VolumeForInstanceResponseType{
			value: "ULTRAHIGH",
		},
		HIGH: VolumeForInstanceResponseType{
			value: "HIGH",
		},
		COMMON: VolumeForInstanceResponseType{
			value: "COMMON",
		},
		NVMESSD: VolumeForInstanceResponseType{
			value: "NVMESSD",
		},
		ULTRAHIGHPRO: VolumeForInstanceResponseType{
			value: "ULTRAHIGHPRO",
		},
		CLOUDSSD: VolumeForInstanceResponseType{
			value: "CLOUDSSD",
		},
		LOCALSSD: VolumeForInstanceResponseType{
			value: "LOCALSSD",
		},
		ESSD: VolumeForInstanceResponseType{
			value: "ESSD",
		},
	}
}

func (c VolumeForInstanceResponseType) Value() string {
	return c.value
}

func (c VolumeForInstanceResponseType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *VolumeForInstanceResponseType) UnmarshalJSON(b []byte) error {
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
