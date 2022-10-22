package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListEngineProductsResponse struct {

	// 分布式消息服务的产品类型。
	Engine *string `json:"engine,omitempty"`

	// 支持的产品版本类型。
	Versions *[]string `json:"versions,omitempty"`

	// 产品规格的详细信息。
	Products       *[]ListEngineProductsEntity `json:"products,omitempty"`
	HttpStatusCode int                         `json:"-"`
}

func (o ListEngineProductsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEngineProductsResponse struct{}"
	}

	return strings.Join([]string{"ListEngineProductsResponse", string(data)}, " ")
}
