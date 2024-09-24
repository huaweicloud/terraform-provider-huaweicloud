package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotReleaseHistoryResponse Response Object
type ShowAutopilotReleaseHistoryResponse struct {
	Body           *[]ReleaseResp `json:"body,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ShowAutopilotReleaseHistoryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotReleaseHistoryResponse struct{}"
	}

	return strings.Join([]string{"ShowAutopilotReleaseHistoryResponse", string(data)}, " ")
}
