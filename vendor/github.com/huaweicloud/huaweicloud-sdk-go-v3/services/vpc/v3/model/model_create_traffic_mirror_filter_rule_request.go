package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorFilterRuleRequest Request Object
type CreateTrafficMirrorFilterRuleRequest struct {
	Body *CreateTrafficMirrorFilterRuleRequestBody `json:"body,omitempty"`
}

func (o CreateTrafficMirrorFilterRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorFilterRuleRequest struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorFilterRuleRequest", string(data)}, " ")
}
