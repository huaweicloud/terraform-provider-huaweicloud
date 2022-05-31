package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowAssetCipherResponse struct {

	// 媒资ID。
	AssetId *string `json:"asset_id,omitempty"`

	// 密钥密文。
	Edk *string `json:"edk,omitempty"`

	// 密钥明文。
	Dk             *string `json:"dk,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowAssetCipherResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetCipherResponse struct{}"
	}

	return strings.Join([]string{"ShowAssetCipherResponse", string(data)}, " ")
}
