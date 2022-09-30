package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Setting struct {

	// 并行执行管道的Filters+Outputs阶段的工作线程数，默认值为CPU核数。
	Workers *int32 `json:"workers,omitempty"`

	// 单个工作线程在尝试执行其Filters和Outputs之前将从inputs收集的最大事件数，该值较大通常更有效，但会增加内存开销，默认为125。
	BatchSize *int32 `json:"batchSize,omitempty"`

	// 每个event被pipeline调度等待的最小时间。 单位毫秒。
	BatchDelayMs *int32 `json:"batchDelayMs,omitempty"`

	// 用于事件缓冲的内部队列模型。memory 为基于内存的传统队列，persisted为基于磁盘的ACKed持久化队列，默认值为memory。
	QueueType string `json:"queueType"`

	// 如果使用持久化队列，则表示强制执行检查点之前写入的最大事件数，默认值为1024。
	QueueCheckPointWrites *int32 `json:"queueCheckPointWrites,omitempty"`

	// 如果使用持久化队列，则表示持久化队列的总容量（以兆字节MB为单位），确保磁盘的容量大于该值，默认值为1024。
	QueueMaxBytesMb *int32 `json:"queueMaxBytesMb,omitempty"`
}

func (o Setting) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Setting struct{}"
	}

	return strings.Join([]string{"Setting", string(data)}, " ")
}
