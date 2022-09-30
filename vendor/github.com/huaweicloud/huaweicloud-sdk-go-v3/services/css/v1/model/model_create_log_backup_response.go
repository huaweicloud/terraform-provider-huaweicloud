package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateLogBackupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateLogBackupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateLogBackupResponse struct{}"
	}

	return strings.Join([]string{"CreateLogBackupResponse", string(data)}, " ")
}
