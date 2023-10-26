package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListFirewallRequest Request Object
type ListFirewallRequest struct {

	// 功能说明：每页返回的个数 取值范围：0~2000
	Limit *int32 `json:"limit,omitempty"`

	// 分页查询起始的资源ID，为空时查询第一页
	Marker *string `json:"marker,omitempty"`

	// ACL唯一标识，填写后接口按照id进行过滤，支持多id同时过滤
	Id *[]string `json:"id,omitempty"`

	// ACL名称，填写后按照名称进行过滤，支持多名称同时过滤
	Name *[]string `json:"name,omitempty"`

	// ACL的状态
	Status *ListFirewallRequestStatus `json:"status,omitempty"`

	// ACL是否启用
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// 功能说明：企业项目ID。可以使用该字段过滤某个企业项目下的ACL。  取值范围：最大长度36字节，带“-”连字符的UUID格式，或者是字符串“0”。“0”表示默认企业项目。若需要查询当前用户所有企业项目绑定的ACL，请传参all_granted_eps。
	EnterpriseProjectId *[]string `json:"enterprise_project_id,omitempty"`
}

func (o ListFirewallRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFirewallRequest struct{}"
	}

	return strings.Join([]string{"ListFirewallRequest", string(data)}, " ")
}

type ListFirewallRequestStatus struct {
	value string
}

type ListFirewallRequestStatusEnum struct {
	ACTIVE   ListFirewallRequestStatus
	INACTIVE ListFirewallRequestStatus
}

func GetListFirewallRequestStatusEnum() ListFirewallRequestStatusEnum {
	return ListFirewallRequestStatusEnum{
		ACTIVE: ListFirewallRequestStatus{
			value: "ACTIVE",
		},
		INACTIVE: ListFirewallRequestStatus{
			value: "INACTIVE",
		},
	}
}

func (c ListFirewallRequestStatus) Value() string {
	return c.value
}

func (c ListFirewallRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListFirewallRequestStatus) UnmarshalJSON(b []byte) error {
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
