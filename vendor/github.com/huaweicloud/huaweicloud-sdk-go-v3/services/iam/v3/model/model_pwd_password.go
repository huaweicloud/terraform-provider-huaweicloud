package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type PwdPassword struct {
	User *PwdPasswordUser `json:"user"`
}

func (o PwdPassword) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PwdPassword struct{}"
	}

	return strings.Join([]string{"PwdPassword", string(data)}, " ")
}
