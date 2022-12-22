package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListWtpProtectHostResponse struct {

	// data list
	DataList *[]WtpProtectHostResponseInfo `json:"data_list,omitempty"`

	// total number
	TotalNum       *int32 `json:"total_num,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListWtpProtectHostResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListWtpProtectHostResponse struct{}"
	}

	return strings.Join([]string{"ListWtpProtectHostResponse", string(data)}, " ")
}
