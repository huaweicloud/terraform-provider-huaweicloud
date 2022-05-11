package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneShowProtocolRequest struct {

	// 身份提供商ID。
	IdpId string `json:"idp_id"`

	// 待查询的协议ID。
	ProtocolId string `json:"protocol_id"`
}

func (o KeystoneShowProtocolRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowProtocolRequest struct{}"
	}

	return strings.Join([]string{"KeystoneShowProtocolRequest", string(data)}, " ")
}
