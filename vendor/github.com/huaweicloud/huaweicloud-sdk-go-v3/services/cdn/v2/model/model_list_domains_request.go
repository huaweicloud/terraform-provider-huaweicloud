package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDomainsRequest Request Object
type ListDomainsRequest struct {

	// 加速域名，采用模糊匹配的方式。（长度限制为1-255字符）。
	DomainName *string `json:"domain_name,omitempty"`

	// 加速域名的业务类型。取值： - web（网站加速） - download（文件下载加速） - video（点播加速） - wholeSite（全站加速）
	BusinessType *string `json:"business_type,omitempty"`

	// 加速域名状态。取值意义： - online表示“已开启” - offline表示“已停用” - configuring表示“配置中” - configure_failed表示“配置失败” - checking表示“审核中” - check_failed表示“审核未通过” - deleting表示“删除中”。
	DomainStatus *string `json:"domain_status,omitempty"`

	// 华为云CDN提供的加速服务范围，包含： - mainland_china 中国大陆 - outside_mainland_china 中国大陆境外 - global 全球。
	ServiceArea *string `json:"service_area,omitempty"`

	// 每页加速域名的数量，取值范围1-10000，默认值为30。
	PageSize *int32 `json:"page_size,omitempty"`

	// 查询的页码，即：从哪一页开始查询，取值范围1-65535，默认值为1。
	PageNumber *int32 `json:"page_number,omitempty"`

	// 展示标签标识 true：展示 false：不展示。
	ShowTags *bool `json:"show_tags,omitempty"`

	// 精准匹配 true：开启 false：关闭。
	ExactMatch *bool `json:"exact_match,omitempty"`

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子帐号调用接口时，该参数必传。  您可以通过调用企业项目管理服务（EPS）的查询企业项目列表接口（ListEnterpriseProject）查询企业项目id。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListDomainsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDomainsRequest struct{}"
	}

	return strings.Join([]string{"ListDomainsRequest", string(data)}, " ")
}
