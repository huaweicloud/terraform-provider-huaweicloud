package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type ProtocolResult struct {

	// 协议ID。
	Id string `json:"id"`

	// 映射ID。
	MappingId string `json:"mapping_id"`

	Links *ProtocolLinks `json:"links"`
}

func (o ProtocolResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtocolResult struct{}"
	}

	return strings.Join([]string{"ProtocolResult", string(data)}, " ")
}
