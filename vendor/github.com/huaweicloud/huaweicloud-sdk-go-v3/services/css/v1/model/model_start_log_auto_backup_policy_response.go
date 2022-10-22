package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type StartLogAutoBackupPolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StartLogAutoBackupPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartLogAutoBackupPolicyResponse struct{}"
	}

	return strings.Join([]string{"StartLogAutoBackupPolicyResponse", string(data)}, " ")
}
