package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListYmlsResponse Response Object
type ListYmlsResponse struct {

	// 集群参数配置列表。该对象中key值以具体获取为准，value拥有以下属性。 - id:  参数配置ID。 - key: 参数名称。 - vlaue:  参数值。 - defaultValue:  参数默认值。 - regex:  参数约束限制。 - desc:  参数中文描述。 - type:  参数类型描述。 - moduleDesc: 参数功能中文描述。 - modifyEnable: 参数是否可修改 true： 可以修改。 false：不可修改。 - enableValue: 参数支持修改的值。 - fileName: 参数存在的文件名称。默认为elasticsearch.yml。 - version:  版本信息。 - descENG: 参数英文描述。 - moduleDescENG:  参数功能英文描述。
	Configurations *interface{} `json:"configurations,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListYmlsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListYmlsResponse struct{}"
	}

	return strings.Join([]string{"ListYmlsResponse", string(data)}, " ")
}
