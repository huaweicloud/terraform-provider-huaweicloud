package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTopUrlRequest Request Object
type ShowTopUrlRequest struct {

	// 查询起始时间戳，只能传0点毫秒时间戳
	StartTime int64 `json:"start_time"`

	// 查询结束时间戳，只能传0点毫秒时间戳
	EndTime int64 `json:"end_time"`

	// 域名列表，多个域名以逗号（半角）分隔，如：www.test1.com,www.test2.com all表示查询名下全部域名。如果域名在查询时间段内无数据，结果将不返回该域名的信息。
	DomainName string `json:"domain_name"`

	// - 参数类型支持：flux(流量),req_num(请求数)
	StatType string `json:"stat_type"`

	// 服务区域：mainland_china(大陆)，outside_mainland_china(海外)，默认为global(全球)
	ServiceArea *string `json:"service_area,omitempty"`

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowTopUrlRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTopUrlRequest struct{}"
	}

	return strings.Join([]string{"ShowTopUrlRequest", string(data)}, " ")
}
