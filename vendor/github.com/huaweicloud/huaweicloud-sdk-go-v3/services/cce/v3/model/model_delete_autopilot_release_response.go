package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAutopilotReleaseResponse Response Object
type DeleteAutopilotReleaseResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteAutopilotReleaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAutopilotReleaseResponse struct{}"
	}

	return strings.Join([]string{"DeleteAutopilotReleaseResponse", string(data)}, " ")
}
