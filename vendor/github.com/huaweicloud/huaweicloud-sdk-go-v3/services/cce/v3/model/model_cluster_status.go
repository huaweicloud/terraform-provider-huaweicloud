package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClusterStatus
type ClusterStatus struct {

	// 集群状态，取值如下 - Available：可用，表示集群处于正常状态。 - Unavailable：不可用，表示集群异常，需手动删除。 - ScalingUp：扩容中，表示集群正处于扩容过程中。 - ScalingDown：缩容中，表示集群正处于缩容过程中。 - Creating：创建中，表示集群正处于创建过程中。 - Deleting：删除中，表示集群正处于删除过程中。 - Upgrading：升级中，表示集群正处于升级过程中。 - Resizing：规格变更中，表示集群正处于变更规格中。 - RollingBack：回滚中，表示集群正处于回滚过程中。 - RollbackFailed：回滚异常，表示集群回滚异常。 - Hibernating：休眠中，表示集群正处于休眠过程中。 - Hibernation：已休眠，表示集群正处于休眠状态。 - Awaking：唤醒中，表示集群正处于从休眠状态唤醒的过程中。 - Empty：集群无任何资源（已废弃） - Error：错误，表示集群资源异常，可尝试手动删除。
	Phase *string `json:"phase,omitempty"`

	// 任务ID,集群当前状态关联的任务ID。当前支持: - 创建集群时返回关联的任务ID，可通过任务ID查询创建集群的附属任务信息； - 删除集群或者删除集群失败时返回关联的任务ID，此字段非空时，可通过任务ID查询删除集群的附属任务信息。 > 任务信息具有一定时效性，仅用于短期跟踪任务进度，请勿用于集群状态判断等额外场景。
	JobID *string `json:"jobID,omitempty"`

	// 集群变为当前状态的原因，在集群在非“Available”状态下时，会返回此参数。
	Reason *string `json:"reason,omitempty"`

	// 集群变为当前状态的原因的详细信息，在集群在非“Available”状态下时，会返回此参数。
	Message *string `json:"message,omitempty"`

	// 集群中 kube-apiserver 的访问地址。
	Endpoints *[]ClusterEndpoints `json:"endpoints,omitempty"`

	// CBC资源锁定
	IsLocked *bool `json:"isLocked,omitempty"`

	// CBC资源锁定场景
	LockScene *string `json:"lockScene,omitempty"`

	// 锁定资源
	LockSource *string `json:"lockSource,omitempty"`

	// 锁定的资源ID
	LockSourceId *string `json:"lockSourceId,omitempty"`

	// 删除配置状态（仅删除请求响应包含）
	DeleteOption *interface{} `json:"deleteOption,omitempty"`

	// 删除状态信息（仅删除请求响应包含）
	DeleteStatus *interface{} `json:"deleteStatus,omitempty"`
}

func (o ClusterStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterStatus struct{}"
	}

	return strings.Join([]string{"ClusterStatus", string(data)}, " ")
}
