package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowLogsRequest struct {

	// 只支持单个域名，如：www.test1.com。
	DomainName string `json:"domain_name"`

	// 查询开始时间，查询开始时间到开始时间+1天内的日志数据，取值范围是距离当前30天内。
	QueryDate int64 `json:"query_date"`

	// 单页最大数量，取值范围为1-10000。
	PageSize *int32 `json:"page_size,omitempty"`

	// 当前查询第几页，取值范围为1-65535。
	PageNumber *int32 `json:"page_number,omitempty"`

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowLogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowLogsRequest struct{}"
	}

	return strings.Join([]string{"ShowLogsRequest", string(data)}, " ")
}
