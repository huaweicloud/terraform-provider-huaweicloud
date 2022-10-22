package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreatePartitionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreatePartitionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePartitionResponse struct{}"
	}

	return strings.Join([]string{"CreatePartitionResponse", string(data)}, " ")
}
