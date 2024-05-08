package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowEngineInstanceExtendProductInfoResponse Response Object
type ShowEngineInstanceExtendProductInfoResponse struct {

	// 消息引擎类型:kafka。
	Engine *string `json:"engine,omitempty"`

	// 消息引擎支持的版本。
	Versions *[]string `json:"versions,omitempty"`

	// 规格变更的产品信息。
	Products       *[]ExtendProductInfoEntity `json:"products,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ShowEngineInstanceExtendProductInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowEngineInstanceExtendProductInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowEngineInstanceExtendProductInfoResponse", string(data)}, " ")
}
