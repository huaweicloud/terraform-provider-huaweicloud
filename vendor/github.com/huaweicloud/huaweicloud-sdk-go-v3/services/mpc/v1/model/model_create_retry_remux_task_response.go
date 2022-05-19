package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateRetryRemuxTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateRetryRemuxTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRetryRemuxTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateRetryRemuxTaskResponse", string(data)}, " ")
}
