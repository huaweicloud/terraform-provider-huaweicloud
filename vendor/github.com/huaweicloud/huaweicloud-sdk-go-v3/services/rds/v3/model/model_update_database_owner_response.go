package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDatabaseOwnerResponse Response Object
type UpdateDatabaseOwnerResponse struct {

	// 调用正常时，返回“successful”。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateDatabaseOwnerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDatabaseOwnerResponse struct{}"
	}

	return strings.Join([]string{"UpdateDatabaseOwnerResponse", string(data)}, " ")
}
