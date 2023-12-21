package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeIsolatedFileResponse Response Object
type ChangeIsolatedFileResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeIsolatedFileResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeIsolatedFileResponse struct{}"
	}

	return strings.Join([]string{"ChangeIsolatedFileResponse", string(data)}, " ")
}
