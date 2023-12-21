package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DeleteNodeRequest Request Object
type DeleteNodeRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 节点ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	NodeId string `json:"node_id"`

	// 集群状态兼容Error参数，用于API平滑切换。 兼容场景下，errorStatus为空则屏蔽Error状态为Deleting状态。
	ErrorStatus *string `json:"errorStatus,omitempty"`

	// 标明是否为nodepool下发的请求。若不为“NoScaleDown”将自动更新对应节点池的实例数
	NodepoolScaleDown *DeleteNodeRequestNodepoolScaleDown `json:"nodepoolScaleDown,omitempty"`
}

func (o DeleteNodeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteNodeRequest struct{}"
	}

	return strings.Join([]string{"DeleteNodeRequest", string(data)}, " ")
}

type DeleteNodeRequestNodepoolScaleDown struct {
	value string
}

type DeleteNodeRequestNodepoolScaleDownEnum struct {
	NO_SCALE_DOWN DeleteNodeRequestNodepoolScaleDown
}

func GetDeleteNodeRequestNodepoolScaleDownEnum() DeleteNodeRequestNodepoolScaleDownEnum {
	return DeleteNodeRequestNodepoolScaleDownEnum{
		NO_SCALE_DOWN: DeleteNodeRequestNodepoolScaleDown{
			value: "NoScaleDown",
		},
	}
}

func (c DeleteNodeRequestNodepoolScaleDown) Value() string {
	return c.value
}

func (c DeleteNodeRequestNodepoolScaleDown) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteNodeRequestNodepoolScaleDown) UnmarshalJSON(b []byte) error {
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
