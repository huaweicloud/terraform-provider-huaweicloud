package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotReleasesResponse Response Object
type ListAutopilotReleasesResponse struct {
	Body           *[]ReleaseResp `json:"body,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListAutopilotReleasesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotReleasesResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotReleasesResponse", string(data)}, " ")
}
