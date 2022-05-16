package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListDomainPermissionsForAgencyRequest struct {

	// 委托方账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`

	// 委托ID，获取方式请参见：[获取委托名、委托ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	AgencyId string `json:"agency_id"`
}

func (o ListDomainPermissionsForAgencyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDomainPermissionsForAgencyRequest struct{}"
	}

	return strings.Join([]string{"ListDomainPermissionsForAgencyRequest", string(data)}, " ")
}
