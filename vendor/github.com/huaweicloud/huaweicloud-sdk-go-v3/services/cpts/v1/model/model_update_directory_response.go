package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDirectoryResponse Response Object
type UpdateDirectoryResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateDirectoryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDirectoryResponse struct{}"
	}

	return strings.Join([]string{"UpdateDirectoryResponse", string(data)}, " ")
}
