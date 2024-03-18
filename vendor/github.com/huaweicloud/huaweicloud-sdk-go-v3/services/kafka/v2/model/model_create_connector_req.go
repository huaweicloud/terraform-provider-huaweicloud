package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type CreateConnectorReq struct {

	// 部署Smart Connect的规格，基准带宽，表示单位时间内传送的最大数据量。请保持和当前实例规格一致。仅老规格实例需要填写。 取值范围：   - 100MB   - 300MB   - 600MB   - 1200MB
	Specification *CreateConnectorReqSpecification `json:"specification,omitempty"`

	// Smart Connect节点数量。不能小于2个。 如果不填，默认是2个。
	NodeCnt *string `json:"node_cnt,omitempty"`

	// 转储节点规格编码。仅老规格实例需要填写。
	SpecCode *string `json:"spec_code,omitempty"`
}

func (o CreateConnectorReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateConnectorReq struct{}"
	}

	return strings.Join([]string{"CreateConnectorReq", string(data)}, " ")
}

type CreateConnectorReqSpecification struct {
	value string
}

type CreateConnectorReqSpecificationEnum struct {
	E_100_MB  CreateConnectorReqSpecification
	E_300_MB  CreateConnectorReqSpecification
	E_600_MB  CreateConnectorReqSpecification
	E_1200_MB CreateConnectorReqSpecification
}

func GetCreateConnectorReqSpecificationEnum() CreateConnectorReqSpecificationEnum {
	return CreateConnectorReqSpecificationEnum{
		E_100_MB: CreateConnectorReqSpecification{
			value: "100MB",
		},
		E_300_MB: CreateConnectorReqSpecification{
			value: "300MB",
		},
		E_600_MB: CreateConnectorReqSpecification{
			value: "600MB",
		},
		E_1200_MB: CreateConnectorReqSpecification{
			value: "1200MB",
		},
	}
}

func (c CreateConnectorReqSpecification) Value() string {
	return c.value
}

func (c CreateConnectorReqSpecification) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateConnectorReqSpecification) UnmarshalJSON(b []byte) error {
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
