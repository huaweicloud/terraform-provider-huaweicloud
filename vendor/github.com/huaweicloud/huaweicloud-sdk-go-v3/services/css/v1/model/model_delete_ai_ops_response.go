package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAiOpsResponse Response Object
type DeleteAiOpsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteAiOpsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAiOpsResponse struct{}"
	}

	return strings.Join([]string{"DeleteAiOpsResponse", string(data)}, " ")
}
