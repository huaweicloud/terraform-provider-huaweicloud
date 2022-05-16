package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateDomainMappingRequest struct {

	// op账号需要携带的特定project_id，当使用op账号时该值为所操作租户的project_id
	SpecifyProject *string `json:"specify_project,omitempty"`

	Body *DomainMapping `json:"body,omitempty"`
}

func (o CreateDomainMappingRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDomainMappingRequest struct{}"
	}

	return strings.Join([]string{"CreateDomainMappingRequest", string(data)}, " ")
}
