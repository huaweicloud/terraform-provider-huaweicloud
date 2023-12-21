package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type GeoBlockingConfigInfo struct {

	// 应用名
	App string `json:"app"`

	// 限制区域列表, 空列表表示不限制。 除中国以外，其他地区代码，2位字母大写。代码格式参阅[ISO 3166-1 alpha-2](https://www.iso.org/obp/ui/#search/code/) 包含如下部分取值： - CN-IN：中国大陆 - CN-HK：中国香港 - CN-MO：中国澳门 - CN-TW：中国台湾 - BR：巴西
	AreaWhitelist *[]string `json:"area_whitelist,omitempty"`
}

func (o GeoBlockingConfigInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GeoBlockingConfigInfo struct{}"
	}

	return strings.Join([]string{"GeoBlockingConfigInfo", string(data)}, " ")
}
