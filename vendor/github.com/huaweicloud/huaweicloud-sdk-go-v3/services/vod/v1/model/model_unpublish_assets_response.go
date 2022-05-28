package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UnpublishAssetsResponse struct {
	AssetInfoArray *[]AssetInfo `json:"asset_info_array,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o UnpublishAssetsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnpublishAssetsResponse struct{}"
	}

	return strings.Join([]string{"UnpublishAssetsResponse", string(data)}, " ")
}
