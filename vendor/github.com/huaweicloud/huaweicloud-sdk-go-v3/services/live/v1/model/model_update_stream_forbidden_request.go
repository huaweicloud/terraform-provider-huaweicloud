package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateStreamForbiddenRequest struct {

	// op账号需要携带的特定project_id，当使用op账号时该值为所操作租户的project_id
	SpecifyProject *string `json:"specify_project,omitempty"`

	Body *StreamForbiddenSetting `json:"body,omitempty"`
}

func (o UpdateStreamForbiddenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateStreamForbiddenRequest struct{}"
	}

	return strings.Join([]string{"UpdateStreamForbiddenRequest", string(data)}, " ")
}
