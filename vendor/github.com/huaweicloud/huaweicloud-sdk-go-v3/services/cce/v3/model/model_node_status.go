package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// NodeStatus
type NodeStatus struct {

	// 节点状态：节点资源生命周期管理（如安装卸载等）状态和集群内k8s node状态的综合体现，取值如下 - Build：创建中，表示节点正处于创建过程中。 - Installing：纳管中，表示节点正处于纳管过程中。 - Upgrading：升级中，表示节点正处于升级过程中。 - Active：正常，表示节点处于正常状态。 - Abnormal：异常，表示节点处于异常状态。 - Deleting： 删除中，表示节点正处于删除过程中。 - Error：故障，表示节点处于故障状态。
	Phase *NodeStatusPhase `json:"phase,omitempty"`

	// 节点最近一次状态检查时间。集群处于异常、冻结或者中间态（例如创建中）时，节点的状态检查动作可能受影响。检查时间超过5分的节点状态不具有参考意义。
	LastProbeTime *string `json:"lastProbeTime,omitempty"`

	// 创建或删除时的任务ID。
	JobID *string `json:"jobID,omitempty"`

	// 底层云服务器或裸金属节点ID。
	ServerId *string `json:"serverId,omitempty"`

	// 节点主网卡私有网段IP地址。
	PrivateIP *string `json:"privateIP,omitempty"`

	// 节点主网卡私有网段IPv6地址。
	PrivateIPv6IP *string `json:"privateIPv6IP,omitempty"`

	// 节点弹性公网IP地址。如果ECS的数据没有实时同步，可在界面上通过“同步节点信息”手动进行更新。
	PublicIP *string `json:"publicIP,omitempty"`

	DeleteStatus *DeleteStatus `json:"deleteStatus,omitempty"`
}

func (o NodeStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeStatus struct{}"
	}

	return strings.Join([]string{"NodeStatus", string(data)}, " ")
}

type NodeStatusPhase struct {
	value string
}

type NodeStatusPhaseEnum struct {
	BUILD      NodeStatusPhase
	INSTALLING NodeStatusPhase
	UPGRADING  NodeStatusPhase
	ACTIVE     NodeStatusPhase
	ABNORMAL   NodeStatusPhase
	DELETING   NodeStatusPhase
	ERROR      NodeStatusPhase
}

func GetNodeStatusPhaseEnum() NodeStatusPhaseEnum {
	return NodeStatusPhaseEnum{
		BUILD: NodeStatusPhase{
			value: "Build",
		},
		INSTALLING: NodeStatusPhase{
			value: "Installing",
		},
		UPGRADING: NodeStatusPhase{
			value: "Upgrading",
		},
		ACTIVE: NodeStatusPhase{
			value: "Active",
		},
		ABNORMAL: NodeStatusPhase{
			value: "Abnormal",
		},
		DELETING: NodeStatusPhase{
			value: "Deleting",
		},
		ERROR: NodeStatusPhase{
			value: "Error",
		},
	}
}

func (c NodeStatusPhase) Value() string {
	return c.value
}

func (c NodeStatusPhase) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NodeStatusPhase) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
