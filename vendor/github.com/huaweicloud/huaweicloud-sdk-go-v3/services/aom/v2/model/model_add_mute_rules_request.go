package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddMuteRulesRequest Request Object
type AddMuteRulesRequest struct {
	Body *MuteRule `json:"body,omitempty"`
}

func (o AddMuteRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddMuteRulesRequest struct{}"
	}

	return strings.Join([]string{"AddMuteRulesRequest", string(data)}, " ")
}
