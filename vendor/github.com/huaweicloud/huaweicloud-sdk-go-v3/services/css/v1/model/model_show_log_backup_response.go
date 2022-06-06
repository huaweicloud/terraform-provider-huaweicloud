package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowLogBackupResponse struct {
	LogList        *[]LogList `json:"logList,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o ShowLogBackupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowLogBackupResponse struct{}"
	}

	return strings.Join([]string{"ShowLogBackupResponse", string(data)}, " ")
}
