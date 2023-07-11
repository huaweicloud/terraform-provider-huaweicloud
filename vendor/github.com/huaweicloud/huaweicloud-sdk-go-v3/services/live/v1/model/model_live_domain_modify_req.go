package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type LiveDomainModifyReq struct {

	// 直播域名，不允许修改
	Domain string `json:"domain"`

	// 直播域名状态，通过修改此字段，实现域名的启用和停用。注意：域名处于“配置中”状态时，不允对该域名执行启停操作。
	Status *LiveDomainModifyReqStatus `json:"status,omitempty"`

	// 企业项目ID
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o LiveDomainModifyReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LiveDomainModifyReq struct{}"
	}

	return strings.Join([]string{"LiveDomainModifyReq", string(data)}, " ")
}

type LiveDomainModifyReqStatus struct {
	value string
}

type LiveDomainModifyReqStatusEnum struct {
	ON  LiveDomainModifyReqStatus
	OFF LiveDomainModifyReqStatus
}

func GetLiveDomainModifyReqStatusEnum() LiveDomainModifyReqStatusEnum {
	return LiveDomainModifyReqStatusEnum{
		ON: LiveDomainModifyReqStatus{
			value: "on",
		},
		OFF: LiveDomainModifyReqStatus{
			value: "off",
		},
	}
}

func (c LiveDomainModifyReqStatus) Value() string {
	return c.value
}

func (c LiveDomainModifyReqStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *LiveDomainModifyReqStatus) UnmarshalJSON(b []byte) error {
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
