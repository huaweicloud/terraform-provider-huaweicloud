package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListLabelsAomPromPostResponse struct {

	// 响应状态。
	Status *string `json:"status,omitempty"`

	// 标签值信息。
	Data           *[]string `json:"data,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ListLabelsAomPromPostResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLabelsAomPromPostResponse struct{}"
	}

	return strings.Join([]string{"ListLabelsAomPromPostResponse", string(data)}, " ")
}
