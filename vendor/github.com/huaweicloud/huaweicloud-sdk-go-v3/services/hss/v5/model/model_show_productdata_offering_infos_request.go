package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowProductdataOfferingInfosRequest Request Object
type ShowProductdataOfferingInfosRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 站点信息：   - HWC_CN ：中国站   - HWC_HK ：国际站   - HWC_EU : 欧洲站
	SiteCode *string `json:"site_code,omitempty"`
}

func (o ShowProductdataOfferingInfosRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowProductdataOfferingInfosRequest struct{}"
	}

	return strings.Join([]string{"ShowProductdataOfferingInfosRequest", string(data)}, " ")
}
