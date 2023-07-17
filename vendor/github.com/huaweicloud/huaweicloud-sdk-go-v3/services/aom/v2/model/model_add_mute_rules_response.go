package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddMuteRulesResponse Response Object
type AddMuteRulesResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o AddMuteRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddMuteRulesResponse struct{}"
	}

	return strings.Join([]string{"AddMuteRulesResponse", string(data)}, " ")
}
