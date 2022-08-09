package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowDomainItemLocationDetailsRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，不传表示查询默认项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 查询开始时间戳，必须设为5分钟整时刻点
	StartTime int64 `json:"start_time"`

	// 查询结束时间戳，必须设为5分钟整时刻点
	EndTime int64 `json:"end_time"`

	// 域名列表，多个域名以逗号（半角）分隔，如：www.test1.com,www.test2.com，all表示查询名下全部域名
	DomainName string `json:"domain_name"`

	// 指标类型列表 网络资源消耗：bw（带宽），flux（流量），ipv6_bw(ipv6带宽)，ipv6_flux(ipv6流量), https_bw(https带宽)，https_flux(https流量) 访问情况：req_num（请求总数），hit_num（请求命中次数），req_time(请求时长) HTTP状态码（组合指标）：status_code_2xx(状态码2xx)，status_code_3xx(状态码3xx)，status_code_4xx(状态码4xx)，status_code_5xx(状态码5xx)
	StatType string `json:"stat_type"`

	// 区域列表，以逗号分隔，all表示查询全部区域
	Region string `json:"region"`

	// 运营商列表，以逗号分隔，all表示查询全部运营商
	Isp string `json:"isp"`
}

func (o ShowDomainItemLocationDetailsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainItemLocationDetailsRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainItemLocationDetailsRequest", string(data)}, " ")
}
