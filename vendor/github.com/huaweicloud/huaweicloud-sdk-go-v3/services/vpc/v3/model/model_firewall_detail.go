package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FirewallDetail
type FirewallDetail struct {

	// 功能说明：ACL唯一标识 取值范围：合法UUID的字符串
	Id string `json:"id"`

	// 功能说明：ACL名称 取值范围：0-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：ACL描述信息 取值范围：0-255个字符 约束：不能包含“<”和“>”。
	Description string `json:"description"`

	// 功能说明：资源所属项目ID
	ProjectId string `json:"project_id"`

	// 功能说明：ACL创建时间 取值范围：UTC时间格式：yyyy-MM-ddTHH:mm:ss；系统自动生成
	CreatedAt string `json:"created_at"`

	// 功能描述：ACL最近一次更新资源的时间 取值范围：UTC时间格式：yyyy-MM-ddTHH:mm:ss；系统自动生成
	UpdatedAt string `json:"updated_at"`

	// 功能说明：ACL是否开启 取值范围：true表示ACL开启；false表示ACL关闭
	AdminStateUp bool `json:"admin_state_up"`

	// 功能说明：网络ACL的状态
	Status string `json:"status"`

	// 功能说明：ACL企业项目ID。 取值范围：最大长度36字节，带“-”连字符的UUID格式，或者是字符串“0”。“0”表示默认企业项目。
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// 功能描述：ACL资源标签
	Tags []ResourceTag `json:"tags"`

	// 功能说明：ACL绑定的子网列表
	Associations []FirewallAssociation `json:"associations"`

	// 功能说明：ACL入方向规则列表
	IngressRules []FirewallRuleDetail `json:"ingress_rules"`

	// 功能说明：ACL出方向规则列表
	EgressRules []FirewallRuleDetail `json:"egress_rules"`
}

func (o FirewallDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FirewallDetail struct{}"
	}

	return strings.Join([]string{"FirewallDetail", string(data)}, " ")
}
