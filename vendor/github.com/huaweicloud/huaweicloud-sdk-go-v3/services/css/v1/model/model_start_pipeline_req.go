package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartPipelineReq struct {

	// 是否保持常驻。   - true: 开启保持常驻。  - false: 关闭保持常驻。  开启“保持常驻”适用于需要长期运行的业务。开启“保持常驻”以后，将会在每个节点上面配置一个守护进程，当logstash服务出现故障的时候，会主动拉起并修复。“保持常驻”不适用于短期运行的业务，因为多次主动拉起logstash服务会导致数据迁移重复。
	KeepAlive *bool `json:"keepAlive,omitempty"`

	// 配置文件名称。
	Names []string `json:"names"`
}

func (o StartPipelineReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartPipelineReq struct{}"
	}

	return strings.Join([]string{"StartPipelineReq", string(data)}, " ")
}
