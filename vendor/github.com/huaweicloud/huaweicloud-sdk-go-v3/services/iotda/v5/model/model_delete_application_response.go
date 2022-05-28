package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteApplicationResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteApplicationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteApplicationResponse struct{}"
	}

	return strings.Join([]string{"DeleteApplicationResponse", string(data)}, " ")
}
