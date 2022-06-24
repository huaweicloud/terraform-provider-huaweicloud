package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

//
type Version struct {

	// 版本号，例如v1。
	Id *string `json:"id,omitempty"`

	// 链接地址信息。
	Links *[]Link `json:"links,omitempty"`

	// 版本状态。  取值“CURRENT”，表示该版本为主推版本。  取值\"SUPPORTED\"，表示支持该版本。  取值“DEPRECATED”，表示为废弃版本，存在后续删除的可能。
	Status *VersionStatus `json:"status,omitempty"`

	// 版本更新时间。  格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指UTC时间。
	Updated *string `json:"updated,omitempty"`
}

func (o Version) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Version struct{}"
	}

	return strings.Join([]string{"Version", string(data)}, " ")
}

type VersionStatus struct {
	value string
}

type VersionStatusEnum struct {
	CURRENT    VersionStatus
	DEPRECATED VersionStatus
	SUPPORTED  VersionStatus
}

func GetVersionStatusEnum() VersionStatusEnum {
	return VersionStatusEnum{
		CURRENT: VersionStatus{
			value: "CURRENT",
		},
		DEPRECATED: VersionStatus{
			value: "DEPRECATED",
		},
		SUPPORTED: VersionStatus{
			value: "SUPPORTED",
		},
	}
}

func (c VersionStatus) Value() string {
	return c.value
}

func (c VersionStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *VersionStatus) UnmarshalJSON(b []byte) error {
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
