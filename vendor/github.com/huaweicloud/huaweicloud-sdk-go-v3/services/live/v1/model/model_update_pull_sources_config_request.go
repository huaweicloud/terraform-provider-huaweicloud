package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePullSourcesConfigRequest Request Object
type UpdatePullSourcesConfigRequest struct {
	Body *ModifyPullSourcesConfig `json:"body,omitempty"`
}

func (o UpdatePullSourcesConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePullSourcesConfigRequest struct{}"
	}

	return strings.Join([]string{"UpdatePullSourcesConfigRequest", string(data)}, " ")
}
