package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type UpdateInstanceReq struct {

	// 实例名称。  由英文字符开头，只能由英文字母、数字、中划线、下划线组成，长度为4~64的字符。
	Name *string `json:"name,omitempty"`

	// 实例的描述信息。  长度不超过1024的字符串。[且字符串不能包含\">\"与\"<\"，字符串首字符不能为\"=\",\"+\",\"-\",\"@\"的全角和半角字符。](tag:hcs,fcs)  > \\与\"在json报文中属于特殊字符，如果参数值中需要显示\\或者\"字符，请在字符前增加转义字符\\，比如\\\\或者\\\"。
	Description *string `json:"description,omitempty"`

	// 维护时间窗开始时间，格式为HH:mm:ss。
	MaintainBegin *string `json:"maintain_begin,omitempty"`

	// 维护时间窗结束时间，格式为HH:mm:ss。
	MaintainEnd *string `json:"maintain_end,omitempty"`

	// 安全组ID。  获取方法如下：登录虚拟私有云服务的控制台界面，在安全组的详情页面查找安全组ID。
	SecurityGroupId *string `json:"security_group_id,omitempty"`

	// 容量阈值策略。  支持两种策略模式： - produce_reject: 生产受限 - time_base: 自动删除
	RetentionPolicy *UpdateInstanceReqRetentionPolicy `json:"retention_policy,omitempty"`

	// 企业项目。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o UpdateInstanceReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceReq struct{}"
	}

	return strings.Join([]string{"UpdateInstanceReq", string(data)}, " ")
}

type UpdateInstanceReqRetentionPolicy struct {
	value string
}

type UpdateInstanceReqRetentionPolicyEnum struct {
	PRODUCE_REJECT UpdateInstanceReqRetentionPolicy
	TIME_BASE      UpdateInstanceReqRetentionPolicy
}

func GetUpdateInstanceReqRetentionPolicyEnum() UpdateInstanceReqRetentionPolicyEnum {
	return UpdateInstanceReqRetentionPolicyEnum{
		PRODUCE_REJECT: UpdateInstanceReqRetentionPolicy{
			value: "produce_reject",
		},
		TIME_BASE: UpdateInstanceReqRetentionPolicy{
			value: "time_base",
		},
	}
}

func (c UpdateInstanceReqRetentionPolicy) Value() string {
	return c.value
}

func (c UpdateInstanceReqRetentionPolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateInstanceReqRetentionPolicy) UnmarshalJSON(b []byte) error {
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
