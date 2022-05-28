package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListPropertiesResponse struct {

	// 设备上报的属性执行结果。Json格式，具体格式需要应用和设备约定。
	Response       *interface{} `json:"response,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListPropertiesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPropertiesResponse struct{}"
	}

	return strings.Join([]string{"ListPropertiesResponse", string(data)}, " ")
}
