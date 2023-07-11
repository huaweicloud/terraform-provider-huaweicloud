package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SupportVersion 支持的版本，包含如下:   - hss.version.basic ：基础版策略组   - hss.version.advanced : 专业版策略组   - hss.version.enterprise : 企业版策略组   - hss.version.premium : 旗舰版策略组   - hss.version.wtp : 网页防篡改版策略组   - hss.version.container.enterprise : 容器版策略组
type SupportVersion struct {
}

func (o SupportVersion) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SupportVersion struct{}"
	}

	return strings.Join([]string{"SupportVersion", string(data)}, " ")
}
