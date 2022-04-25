package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneShowEndpointRequest struct {

	// 待查询的终端节点ID。
	EndpointId string `json:"endpoint_id"`
}

func (o KeystoneShowEndpointRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowEndpointRequest struct{}"
	}

	return strings.Join([]string{"KeystoneShowEndpointRequest", string(data)}, " ")
}
