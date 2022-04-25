package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 创建VPC的请求体
type CreateVpcOption struct {

	// 功能描述：VPC的名称信息 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：VPC的描述信息 取值范围：0-255个字符，不能包含“<”和“>”。
	Description *string `json:"description,omitempty"`

	// 功能说明：vpc下可用子网的范围 取值范围： −10.0.0.0/8~10.255.255.240/28 −172.16.0.0/12 ~ 172.31.255.240/28 −192.168.0.0/16 ~ 192.168.255.240/28 约束：必须是cidr格式，例如:192.168.0.0/16
	Cidr string `json:"cidr"`

	// 功能说明：企业项目ID。创建vpc时，给vpc绑定企业项目ID。 取值范围：最大长度36字节，带“-”连字符的UUID格式，或者是字符串“0”。“0”表示默认企业项目。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 功能说明：VPC的标签信息，详情参见Tag对象 取值范围：0-10个标签键值对
	Tags *[]Tag `json:"tags,omitempty"`
}

func (o CreateVpcOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateVpcOption struct{}"
	}

	return strings.Join([]string{"CreateVpcOption", string(data)}, " ")
}
