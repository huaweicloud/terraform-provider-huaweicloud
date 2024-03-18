package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AppRules 服务参数。
type AppRules struct {

	// 规则创建时间(创建时不传，修改时传查询返回的createTime)。
	CreateTime *string `json:"createTime,omitempty"`

	// true、false 规则是否启用。
	Enable bool `json:"enable"`

	// aom_inventory_rules_event规则事件名称，对于服务发现固定为\"aom_inventory_rules_event\"。
	EventName string `json:"eventName"`

	// 主机ID（暂不使用，传空即可）。
	Hostid *[]string `json:"hostid,omitempty"`

	// 创建时填空，修改时填规则ID。
	Id string `json:"id"`

	// 规则名称。 字符长度为4到63位，以小写字母a-z开头，只能包含0-9/a-z/-，不能以-结尾。
	Name string `json:"name"`

	// 租户从IAM申请到的projectid，一般为32位字符串。
	Projectid string `json:"projectid"`

	Spec *AppRulesSpec `json:"spec"`

	// 自定义描述信息
	Desc *string `json:"desc,omitempty"`
}

func (o AppRules) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AppRules struct{}"
	}

	return strings.Join([]string{"AppRules", string(data)}, " ")
}
