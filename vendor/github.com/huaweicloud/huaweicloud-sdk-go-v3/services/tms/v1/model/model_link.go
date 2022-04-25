package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// API的URL地址。
type Link struct {

	// API的URL地址。
	Href string `json:"href"`

	// self
	Rel string `json:"rel"`
}

func (o Link) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Link struct{}"
	}

	return strings.Join([]string{"Link", string(data)}, " ")
}
