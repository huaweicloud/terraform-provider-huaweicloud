package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Storage 磁盘初始化配置管理参数。  该参数配置逻辑较为复杂，详细说明请参见[节点磁盘挂载](node_storage_example.xml)。  该参数缺省时，按照extendParam中的DockerLVMConfigOverride（已废弃）参数进行磁盘管理。此参数对1.15.11及以上集群版本支持。  > 如存在节点规格涉及本地盘并同时使用云硬盘场景时，请勿缺省此参数，避免出现将用户未期望的磁盘分区。  > 如希望数据盘取值范围调整至20-32768，请勿缺省此参数。  > 如希望使用共享磁盘空间(取消runtime和kubernetes分区)，请勿缺省此参数，共享磁盘空间请参考[数据盘空间分配说明](cce_01_0341.xml)。
type Storage struct {

	// 磁盘选择，根据matchLabels和storageType对匹配的磁盘进行管理。磁盘匹配存在先后顺序，靠前的匹配规则优先匹配。
	StorageSelectors []StorageSelectors `json:"storageSelectors"`

	// 由多个存储设备组成的存储组，用于各个存储空间的划分。
	StorageGroups []StorageGroups `json:"storageGroups"`
}

func (o Storage) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Storage struct{}"
	}

	return strings.Join([]string{"Storage", string(data)}, " ")
}
