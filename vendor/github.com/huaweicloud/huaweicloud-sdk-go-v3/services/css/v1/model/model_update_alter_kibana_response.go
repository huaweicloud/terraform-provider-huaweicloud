package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAlterKibanaResponse Response Object
type UpdateAlterKibanaResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateAlterKibanaResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAlterKibanaResponse struct{}"
	}

	return strings.Join([]string{"UpdateAlterKibanaResponse", string(data)}, " ")
}
