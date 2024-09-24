package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAgenciesRequest Request Object
type ListAgenciesRequest struct {

	// 委托方账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`

	// 被委托方账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	TrustDomainId *string `json:"trust_domain_id,omitempty"`

	// 委托名，获取方式请参见：[获取委托名、委托ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	Name *string `json:"name,omitempty"`

	// 分页查询时数据的页数，查询值最小为1。需要与per_page同时存在。
	Page *int32 `json:"page,omitempty"`

	// 分页查询时每页的数据个数，取值范围为[1,500]。需要与page同时存在。
	PerPage *int32 `json:"per_page,omitempty"`
}

func (o ListAgenciesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAgenciesRequest struct{}"
	}

	return strings.Join([]string{"ListAgenciesRequest", string(data)}, " ")
}
