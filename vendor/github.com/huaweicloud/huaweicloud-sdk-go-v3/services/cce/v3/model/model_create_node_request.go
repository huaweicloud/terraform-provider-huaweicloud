package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateNodeRequest Request Object
type CreateNodeRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 标明是否为nodepool下发的请求。若不为“NodepoolScaleUp”将自动更新对应节点池的实例数
	NodepoolScaleUp *CreateNodeRequestNodepoolScaleUp `json:"nodepoolScaleUp,omitempty"`

	Body *NodeCreateRequest `json:"body,omitempty"`
}

func (o CreateNodeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateNodeRequest struct{}"
	}

	return strings.Join([]string{"CreateNodeRequest", string(data)}, " ")
}

type CreateNodeRequestNodepoolScaleUp struct {
	value string
}

type CreateNodeRequestNodepoolScaleUpEnum struct {
	NODEPOOL_SCALE_UP CreateNodeRequestNodepoolScaleUp
}

func GetCreateNodeRequestNodepoolScaleUpEnum() CreateNodeRequestNodepoolScaleUpEnum {
	return CreateNodeRequestNodepoolScaleUpEnum{
		NODEPOOL_SCALE_UP: CreateNodeRequestNodepoolScaleUp{
			value: "NodepoolScaleUp",
		},
	}
}

func (c CreateNodeRequestNodepoolScaleUp) Value() string {
	return c.value
}

func (c CreateNodeRequestNodepoolScaleUp) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateNodeRequestNodepoolScaleUp) UnmarshalJSON(b []byte) error {
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
