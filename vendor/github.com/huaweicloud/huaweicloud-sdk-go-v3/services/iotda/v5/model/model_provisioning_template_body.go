package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ProvisioningTemplateBody 预调配模板详细内容。
type ProvisioningTemplateBody struct {

	// **参数说明**：预调配模板参数， ，配置格式为{\"parameter\":{\"type\":\"String\"}} ，其中parameter目前支持从预调配设备的证书的使用者字段提取内容，证书必须包含模板中定义的所有参数值，华为云IoT平台定义了可在预调配模板中声明和引用的如下参数: - iotda::certificate::country (国家/地区,C ) - iotda::certificate::organization (组织,O) - iotda::certificate::organizational_unit (组织单位,OU) - iotda::certificate::distinguished_name_qualifier (可辨别名称限定符,dnQualifier) - iotda::certificate::state_name (省市,ST) - iotda::certificate::common_name (公用名,CN) - iotda::certificate::serial_number (序列号,serialNumber)  type描述parameter的类型，目前仅支持string。  配置样例：  '{\"iotda::certificate::country\":{\"type\":\"String\"},  \"iotda::certificate::organization\":{\"type\":\"String\"},  \"iotda::certificate::organizational_unit\":{\"type\":\"String\"},  \"iotda::certificate::distinguished_name_qualifier\":{\"type\":\"String\"},  \"iotda::certificate::state_name\":{\"type\":\"String\"},  \"iotda::certificate::common_name\":{\"type\":\"String\"},  \"iotda::certificate::serial_number\":{\"type\":\"String\"}}'
	Parameters *interface{} `json:"parameters"`

	Resources *TemplateResource `json:"resources"`
}

func (o ProvisioningTemplateBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProvisioningTemplateBody struct{}"
	}

	return strings.Join([]string{"ProvisioningTemplateBody", string(data)}, " ")
}
