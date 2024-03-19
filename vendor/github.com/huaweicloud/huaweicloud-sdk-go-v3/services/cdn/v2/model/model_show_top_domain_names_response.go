package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTopDomainNamesResponse Response Object
type ShowTopDomainNamesResponse struct {

	// top域名信息
	TopDomainNames *[]map[string]interface{} `json:"top_domain_names,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o ShowTopDomainNamesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTopDomainNamesResponse struct{}"
	}

	return strings.Join([]string{"ShowTopDomainNamesResponse", string(data)}, " ")
}
