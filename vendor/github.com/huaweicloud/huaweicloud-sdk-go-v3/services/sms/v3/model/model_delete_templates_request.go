package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteTemplatesRequest struct {
	Body *DeletetemplatesReq `json:"body,omitempty"`
}

func (o DeleteTemplatesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTemplatesRequest struct{}"
	}

	return strings.Join([]string{"DeleteTemplatesRequest", string(data)}, " ")
}
