package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateQuotasOrderRequest Request Object
type CreateQuotasOrderRequest struct {

	// 缺省值:application/json; charset=utf-8
	ContentType *string `json:"Content-Type,omitempty"`

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *CreateQuotasOrderRequestInfo `json:"body,omitempty"`
}

func (o CreateQuotasOrderRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateQuotasOrderRequest struct{}"
	}

	return strings.Join([]string{"CreateQuotasOrderRequest", string(data)}, " ")
}
