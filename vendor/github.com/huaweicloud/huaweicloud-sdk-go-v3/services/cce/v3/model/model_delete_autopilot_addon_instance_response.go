package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAutopilotAddonInstanceResponse Response Object
type DeleteAutopilotAddonInstanceResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteAutopilotAddonInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAutopilotAddonInstanceResponse struct{}"
	}

	return strings.Join([]string{"DeleteAutopilotAddonInstanceResponse", string(data)}, " ")
}
