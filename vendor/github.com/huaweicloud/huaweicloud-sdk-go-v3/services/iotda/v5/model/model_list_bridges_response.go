package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListBridgesResponse Response Object
type ListBridgesResponse struct {

	// 网桥列表。
	Bridges *[]BridgeResponse `json:"bridges,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ListBridgesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListBridgesResponse struct{}"
	}

	return strings.Join([]string{"ListBridgesResponse", string(data)}, " ")
}
