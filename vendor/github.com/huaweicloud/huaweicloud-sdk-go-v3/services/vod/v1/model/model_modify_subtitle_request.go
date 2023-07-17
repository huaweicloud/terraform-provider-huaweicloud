package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifySubtitleRequest Request Object
type ModifySubtitleRequest struct {
	Body *SubtitleModifyReq `json:"body,omitempty"`
}

func (o ModifySubtitleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifySubtitleRequest struct{}"
	}

	return strings.Join([]string{"ModifySubtitleRequest", string(data)}, " ")
}
