package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowLogsRequest Request Object
type ShowLogsRequest struct {

	// 只支持单个域名，如：www.test1.com。
	DomainName string `json:"domain_name"`

	// 查询开始时间，时间格式为整点毫秒时间戳，此参数传空值时默认为当天0点。
	StartTime *int64 `json:"start_time,omitempty"`

	// 查询结束时间（不包含结束时间），时间格式为整点毫秒时间戳，与开始时间的最大跨度为30天，此参数传空值时默认为开始时间加1天。
	EndTime *int64 `json:"end_time,omitempty"`

	// 单页最大数量，取值范围为1-10000，默认值：10。
	PageSize *int32 `json:"page_size,omitempty"`

	// 当前查询第几页，取值范围为1-65535，默认值：1。
	PageNumber *int32 `json:"page_number,omitempty"`

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子帐号调用接口时，该参数必传。  您可以通过调用企业项目管理服务（EPS）的查询企业项目列表接口（ListEnterpriseProject）查询企业项目id。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowLogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowLogsRequest struct{}"
	}

	return strings.Join([]string{"ShowLogsRequest", string(data)}, " ")
}
