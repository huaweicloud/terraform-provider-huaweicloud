package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAiOpsResponse Response Object
type CreateAiOpsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateAiOpsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAiOpsResponse struct{}"
	}

	return strings.Join([]string{"CreateAiOpsResponse", string(data)}, " ")
}
