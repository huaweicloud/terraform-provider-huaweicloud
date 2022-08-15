package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowDomainItemDetailsResponse struct {
	DomainItemDetails *DomainItemDetail `json:"domain_item_details,omitempty"`
	HttpStatusCode    int               `json:"-"`
}

func (o ShowDomainItemDetailsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainItemDetailsResponse struct{}"
	}

	return strings.Join([]string{"ShowDomainItemDetailsResponse", string(data)}, " ")
}
