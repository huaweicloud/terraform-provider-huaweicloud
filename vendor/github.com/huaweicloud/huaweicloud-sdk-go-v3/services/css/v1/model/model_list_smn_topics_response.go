package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSmnTopicsResponse Response Object
type ListSmnTopicsResponse struct {

	// SMN主题名称列表。
	TopicsName     *[]string `json:"topicsName,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ListSmnTopicsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSmnTopicsResponse struct{}"
	}

	return strings.Join([]string{"ListSmnTopicsResponse", string(data)}, " ")
}
