package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateAssetMetaResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateAssetMetaResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAssetMetaResponse struct{}"
	}

	return strings.Join([]string{"UpdateAssetMetaResponse", string(data)}, " ")
}
