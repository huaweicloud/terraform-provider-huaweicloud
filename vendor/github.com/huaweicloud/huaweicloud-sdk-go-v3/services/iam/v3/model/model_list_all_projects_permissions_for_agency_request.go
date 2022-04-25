package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListAllProjectsPermissionsForAgencyRequest struct {

	// 委托ID，获取方式请参见：[获取委托名、委托ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	AgencyId string `json:"agency_id"`

	// 账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`
}

func (o ListAllProjectsPermissionsForAgencyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAllProjectsPermissionsForAgencyRequest struct{}"
	}

	return strings.Join([]string{"ListAllProjectsPermissionsForAgencyRequest", string(data)}, " ")
}
