package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeCoreResponse Response Object
type UpgradeCoreResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpgradeCoreResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeCoreResponse struct{}"
	}

	return strings.Join([]string{"UpgradeCoreResponse", string(data)}, " ")
}
