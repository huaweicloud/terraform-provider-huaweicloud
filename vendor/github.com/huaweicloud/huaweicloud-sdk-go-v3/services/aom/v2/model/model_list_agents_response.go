package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAgentsResponse Response Object
type ListAgentsResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListAgentsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAgentsResponse struct{}"
	}

	return strings.Join([]string{"ListAgentsResponse", string(data)}, " ")
}
