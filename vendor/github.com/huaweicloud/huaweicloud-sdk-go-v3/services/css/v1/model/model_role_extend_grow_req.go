package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RoleExtendGrowReq 集群扩容请求详细描述。
type RoleExtendGrowReq struct {

	// 扩容实例类型。取值为ess、ess-master或ess-client，可以选择其中一个或多个之间的组合但不可以重复选择。  ess-master、ess-client节点只支持增加实例个数。
	Type string `json:"type"`

	// 扩容实例个数。集群已有实例个数和增加实例个数总和不能超过32。若无需扩容该参数将该参数设置为0即可。
	Nodesize int32 `json:"nodesize"`

	// 扩容实例存储容量。集群原实例存储容量和扩容实例存储容量之和不能超过创建集群时对应默认实例存储容量上限。若无需扩容该参数将该参数设置为0即可。[当集群为包周期集群时不支持同时修改节点个数和磁盘容量。](tag: hc,tag: hws)  单位：GB。  - ess节点、ess-cold节点扩容步长为20。  - ess-master节点、ess-client节点不允许扩容存储。
	Disksize int32 `json:"disksize"`
}

func (o RoleExtendGrowReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleExtendGrowReq struct{}"
	}

	return strings.Join([]string{"RoleExtendGrowReq", string(data)}, " ")
}
