package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateStreamForbiddenRequest struct {
	Body *StreamForbiddenSetting `json:"body,omitempty"`
}

func (o CreateStreamForbiddenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateStreamForbiddenRequest struct{}"
	}

	return strings.Join([]string{"CreateStreamForbiddenRequest", string(data)}, " ")
}
