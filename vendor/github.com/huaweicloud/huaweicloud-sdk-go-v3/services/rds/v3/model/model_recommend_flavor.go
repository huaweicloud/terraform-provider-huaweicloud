package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// RecommendFlavor 规格信息
type RecommendFlavor struct {

	// 按照实例类型查询
	Type RecommendFlavorType `json:"type"`

	// 规格码
	FlavorRef string `json:"flavor_ref"`

	// CPU大小
	Cpu string `json:"cpu"`

	// 内存大小（单位：GB）
	Mem string `json:"mem"`

	// 规格类型
	GroupType string `json:"group_type"`

	// 磁盘规格信息
	VolumeFlavors []VolumeFlavor `json:"volume_flavors"`
}

func (o RecommendFlavor) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecommendFlavor struct{}"
	}

	return strings.Join([]string{"RecommendFlavor", string(data)}, " ")
}

type RecommendFlavorType struct {
	value string
}

type RecommendFlavorTypeEnum struct {
	HA     RecommendFlavorType
	SINGLE RecommendFlavorType
}

func GetRecommendFlavorTypeEnum() RecommendFlavorTypeEnum {
	return RecommendFlavorTypeEnum{
		HA: RecommendFlavorType{
			value: "Ha",
		},
		SINGLE: RecommendFlavorType{
			value: "Single",
		},
	}
}

func (c RecommendFlavorType) Value() string {
	return c.value
}

func (c RecommendFlavorType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecommendFlavorType) UnmarshalJSON(b []byte) error {
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
