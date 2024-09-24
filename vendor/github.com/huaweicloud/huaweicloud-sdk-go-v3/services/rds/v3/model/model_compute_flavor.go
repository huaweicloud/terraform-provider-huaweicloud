package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ComputeFlavor 查询数据库可变更规格接口，响应体中的计算规格详情
type ComputeFlavor struct {

	// 规格ID，该字段唯一。
	Id string `json:"id"`

	// 资源规格编码。例如：rds.pg.m1.xlarge.rr。  更多规格说明请参考数据库实例规格。  “rds”代表RDS产品。 “pg”代表数据库引擎。 “m1.xlarge”代表性能规格，为高内存类型。 “rr”表示只读实例（“.ha”表示主备实例）。
	Code string `json:"code"`

	// CPU大小。例如：1表示1U。
	Vcpus string `json:"vcpus"`

	// 内存大小，单位为GB。
	Ram string `json:"ram"`

	// 规格所在az的状态，包含以下状态：  normal：在售。 unsupported：暂不支持该规格。 sellout：售罄。
	AzStatus map[string]string `json:"az_status"`
}

func (o ComputeFlavor) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ComputeFlavor struct{}"
	}

	return strings.Join([]string{"ComputeFlavor", string(data)}, " ")
}
