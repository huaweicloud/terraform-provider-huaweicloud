package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateNewCaseResponse Response Object
type UpdateNewCaseResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateNewCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateNewCaseResponse struct{}"
	}

	return strings.Join([]string{"UpdateNewCaseResponse", string(data)}, " ")
}
