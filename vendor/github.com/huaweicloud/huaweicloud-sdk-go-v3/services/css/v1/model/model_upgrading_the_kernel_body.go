package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type UpgradingTheKernelBody struct {

	// 目标镜像版本ID。
	TargetImageId string `json:"target_image_id"`

	// 升级类型。 - same：同版本。 - cross：跨版本。
	UpgradeType UpgradingTheKernelBodyUpgradeType `json:"upgrade_type"`

	// 是否进行备份校验。 - true：进行校验。 - false：不进行校验。
	IndicesBackupCheck bool `json:"indices_backup_check"`

	// 委托名称，委托给CSS，允许CSS调用您的其他云服务。
	Agency string `json:"agency"`

	// 是否校验负载。默认为true。 - true：进行校验。 - false：不进行校验。
	CheckLoad *bool `json:"check_load,omitempty"`
}

func (o UpgradingTheKernelBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradingTheKernelBody struct{}"
	}

	return strings.Join([]string{"UpgradingTheKernelBody", string(data)}, " ")
}

type UpgradingTheKernelBodyUpgradeType struct {
	value string
}

type UpgradingTheKernelBodyUpgradeTypeEnum struct {
	SAME  UpgradingTheKernelBodyUpgradeType
	CROSS UpgradingTheKernelBodyUpgradeType
}

func GetUpgradingTheKernelBodyUpgradeTypeEnum() UpgradingTheKernelBodyUpgradeTypeEnum {
	return UpgradingTheKernelBodyUpgradeTypeEnum{
		SAME: UpgradingTheKernelBodyUpgradeType{
			value: "same",
		},
		CROSS: UpgradingTheKernelBodyUpgradeType{
			value: "cross",
		},
	}
}

func (c UpgradingTheKernelBodyUpgradeType) Value() string {
	return c.value
}

func (c UpgradingTheKernelBodyUpgradeType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpgradingTheKernelBodyUpgradeType) UnmarshalJSON(b []byte) error {
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
