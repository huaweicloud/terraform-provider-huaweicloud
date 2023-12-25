package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePromInstanceResponse Response Object
type DeletePromInstanceResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeletePromInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePromInstanceResponse struct{}"
	}

	return strings.Join([]string{"DeletePromInstanceResponse", string(data)}, " ")
}
