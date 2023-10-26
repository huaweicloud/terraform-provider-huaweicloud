package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteBatchTaskResponse Response Object
type DeleteBatchTaskResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteBatchTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteBatchTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteBatchTaskResponse", string(data)}, " ")
}
