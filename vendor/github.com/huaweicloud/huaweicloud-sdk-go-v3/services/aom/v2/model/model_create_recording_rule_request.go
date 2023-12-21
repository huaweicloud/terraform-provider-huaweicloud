package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateRecordingRuleRequest Request Object
type CreateRecordingRuleRequest struct {

	// prometheus实例id。
	PrometheusInstance string `json:"prometheus_instance"`

	Body *RecordingRuleRequest `json:"body,omitempty"`
}

func (o CreateRecordingRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRecordingRuleRequest struct{}"
	}

	return strings.Join([]string{"CreateRecordingRuleRequest", string(data)}, " ")
}
