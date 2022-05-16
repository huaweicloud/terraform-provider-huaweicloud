package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListCustomPoliciesRequest struct {

	// 分页查询时数据的页数，查询值最小为1。需要与per_page同时存在。
	Page *int32 `json:"page,omitempty"`

	// 分页查询时每页的数据个数，取值范围为[1,300]。需要与page同时存在。
	PerPage *int32 `json:"per_page,omitempty"`
}

func (o ListCustomPoliciesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListCustomPoliciesRequest struct{}"
	}

	return strings.Join([]string{"ListCustomPoliciesRequest", string(data)}, " ")
}
