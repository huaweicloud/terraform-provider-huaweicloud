package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateYmlsResponse struct {

	// 修改是否成功。 - true: 修改成功。 - false: 修改失败。
	Acknowledged *bool `json:"acknowledged,omitempty"`

	// 错误信息描述。当acknowledged为true时，该字段返回null。
	ExternalMessage *string `json:"externalMessage,omitempty"`

	// HTTP错误信息。默认为null。
	HttpErrorResponse *string `json:"httpErrorResponse,omitempty"`
	HttpStatusCode    int     `json:"-"`
}

func (o UpdateYmlsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateYmlsResponse struct{}"
	}

	return strings.Join([]string{"UpdateYmlsResponse", string(data)}, " ")
}
