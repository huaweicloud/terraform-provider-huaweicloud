package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteMuteRulesResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteMuteRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMuteRulesResponse struct{}"
	}

	return strings.Join([]string{"DeleteMuteRulesResponse", string(data)}, " ")
}
