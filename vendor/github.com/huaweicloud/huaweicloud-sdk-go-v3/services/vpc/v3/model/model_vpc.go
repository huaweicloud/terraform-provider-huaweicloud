package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type Vpc struct {

	// 功能描述：VPC对应的唯一标识 取值范围：带“-”的UUID格式
	Id string `json:"id"`

	// 功能说明：VPC对应的名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：VPC的描述信息 取值范围：0-255个字符，不能包含“<”和“>”
	Description string `json:"description"`

	// 功能说明：VPC下可用子网的范围 取值范围： - 10.0.0.0/8~10.255.255.240/28 - 172.16.0.0/12 ~ 172.31.255.240/28 - 192.168.0.0/16 ~ 192.168.255.240/28 不指定cidr时，默认值为“” 约束：必须是ipv4 cidr格式，例如:192.168.0.0/16
	Cidr string `json:"cidr"`

	// 功能描述：VPC的扩展网段 取值范围： 约束：目前只支持ipv4
	ExtendCidrs []string `json:"extend_cidrs"`

	// 功能说明：VPC对应的状态 取值范围：PENDING：创建中；ACTIVE：创建成功
	Status string `json:"status"`

	// 功能说明：VPC所属的项目ID
	ProjectId string `json:"project_id"`

	// 功能说明：VPC所属的企业项目ID。 取值范围：最大长度36字节，带“-”连字符的UUID格式，或者是字符串“0”。“0”表示默认企业项目。
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// 功能说明：VPC创建时间 取值范围：UTC时间格式：yyyy-MM-ddTHH:mm:ss
	CreatedAt *sdktime.SdkTime `json:"created_at"`

	// 功能说明：VPC更新时间 取值范围：UTC时间格式：yyyy-MM-ddTHH:mm:ss
	UpdatedAt *sdktime.SdkTime `json:"updated_at"`

	// 功能说明：VPC关联资源类型和数量 取值范围：目前只返回VPC关联的routetable和virsubnet
	CloudResources []CloudResource `json:"cloud_resources"`

	// 功能说明：VPC的标签信息，详情参见Tag对象 取值范围：0-10个标签键值对
	Tags []Tag `json:"tags"`
}

func (o Vpc) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Vpc struct{}"
	}

	return strings.Join([]string{"Vpc", string(data)}, " ")
}
