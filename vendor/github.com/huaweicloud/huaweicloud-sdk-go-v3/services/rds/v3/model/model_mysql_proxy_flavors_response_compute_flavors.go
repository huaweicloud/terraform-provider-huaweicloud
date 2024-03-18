package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlProxyFlavorsResponseComputeFlavors struct {

	// 数据库代理规格ID。
	Id *string `json:"id,omitempty"`

	// 数据库代理规格码。
	Code *string `json:"code,omitempty"`

	// CPU大小。例如：1表示1U。
	Cpu *string `json:"cpu,omitempty"`

	// 内存大小，单位为GB。
	Mem *string `json:"mem,omitempty"`

	// 数据库类型。
	DbType *string `json:"db_type,omitempty"`

	// 可用区信息，其中key是该规格绑定的可用区，value是该规格在对应可用区中的状态。 取值范围：     normal：正常     abandon：禁用      - 仅展示数据库主实例所在可用区规格状态。
	AzStatus *interface{} `json:"az_status,omitempty"`
}

func (o MysqlProxyFlavorsResponseComputeFlavors) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlProxyFlavorsResponseComputeFlavors struct{}"
	}

	return strings.Join([]string{"MysqlProxyFlavorsResponseComputeFlavors", string(data)}, " ")
}
