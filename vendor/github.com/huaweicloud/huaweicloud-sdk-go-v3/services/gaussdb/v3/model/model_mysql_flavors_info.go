package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlFlavorsInfo struct {
	// CPU大小。例如：1表示1U。

	Vcpus string `json:"vcpus"`
	// 内存大小，单位为GB。

	Ram string `json:"ram"`
	// 规格类型，取值为arm和x86。

	Type string `json:"type"`
	// 规格ID，该字段唯一

	Id string `json:"id"`
	// 资源规格编码，同创建指定的flavor_ref。例如：gaussdb.mysql.xlarge.x86.4。

	SpecCode string `json:"spec_code"`
	// 数据库版本号。

	VersionName string `json:"version_name"`
	// 实例类型。目前仅支持Cluster。

	InstanceMode string `json:"instance_mode"`
	// 规格所在az的状态，包含以下状态： - normal，在售 - unsupported，暂不支持该规格 - sellout，售罄。

	AzStatus map[string]string `json:"az_status"`
}

func (o MysqlFlavorsInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlFlavorsInfo struct{}"
	}

	return strings.Join([]string{"MysqlFlavorsInfo", string(data)}, " ")
}
