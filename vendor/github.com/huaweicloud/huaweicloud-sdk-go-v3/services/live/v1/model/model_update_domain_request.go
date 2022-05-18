package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateDomainRequest struct {
	Body *LiveDomainModifyReq `json:"body,omitempty"`
}

func (o UpdateDomainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainRequest struct{}"
	}

	return strings.Join([]string{"UpdateDomainRequest", string(data)}, " ")
}
