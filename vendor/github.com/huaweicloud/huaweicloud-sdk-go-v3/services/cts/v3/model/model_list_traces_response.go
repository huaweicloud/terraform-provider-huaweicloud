package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTracesResponse struct {

	// 本次查询事件列表返回事件数组。
	Traces *[]Traces `json:"traces,omitempty"`

	MetaData       *MetaData `json:"meta_data,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ListTracesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTracesResponse struct{}"
	}

	return strings.Join([]string{"ListTracesResponse", string(data)}, " ")
}
