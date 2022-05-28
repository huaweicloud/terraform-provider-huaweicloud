package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteBatchTaskFileResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteBatchTaskFileResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteBatchTaskFileResponse struct{}"
	}

	return strings.Join([]string{"DeleteBatchTaskFileResponse", string(data)}, " ")
}
