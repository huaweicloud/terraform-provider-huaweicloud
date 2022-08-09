package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 服务配置信息
type ObsForwarding struct {

	// **参数说明**：OBS服务对应的region区域
	RegionName string `json:"region_name"`

	// **参数说明**：OBS服务对应的projectId信息
	ProjectId string `json:"project_id"`

	// **参数说明**：OBS服务对应的桶名称
	BucketName string `json:"bucket_name"`

	// **参数说明**：OBS服务对应桶的区域
	Location *string `json:"location,omitempty"`

	// **参数说明**：OBS服务中存储通道文件的自定义目录,多级目录可用(/)进行分隔，不可以斜杠(/)开头或结尾，不能包含两个以上相邻的斜杠(/) **取值范围**: 英文字母(a-zA-Z)、数字(0-9)、下划线(_)、中划线(-)、斜杠(/)和大括号({})，最大字符长度256个字符。其中大括号只能用于对应模板参数。 **模板参数**:    - \\{YYYY\\} 年   - \\{MM\\} 月   - \\{DD\\} 日   - \\{HH\\} 小时   - \\{appId\\} 应用ID   - \\{deviceId\\} 设备ID   例如:自定义目录结构为\\{YYYY\\}/\\{MM\\}/\\{DD\\}/\\{HH\\},则会在转发数据时，根据当前时间往对应的目录结构2021>08>11>09下生成对应的数据。
	FilePath *string `json:"file_path,omitempty"`
}

func (o ObsForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ObsForwarding struct{}"
	}

	return strings.Join([]string{"ObsForwarding", string(data)}, " ")
}
