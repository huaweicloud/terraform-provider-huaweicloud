package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePortTagResponse Response Object
type DeletePortTagResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeletePortTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePortTagResponse struct{}"
	}

	return strings.Join([]string{"DeletePortTagResponse", string(data)}, " ")
}
