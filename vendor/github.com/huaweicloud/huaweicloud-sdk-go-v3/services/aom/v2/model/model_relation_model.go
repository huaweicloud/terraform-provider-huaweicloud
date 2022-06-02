package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 查询条件。
type RelationModel struct {

	// 指定查询字段的key，对应metadata里面的key 。
	Key *string `json:"key,omitempty"`

	// 查询条件中指定key的值。
	Value *[]string `json:"value,omitempty"`

	// 该条件与其他条件的组合方式。 AND：必须满足所有条件； OR：可以满足其中一个条件； NOT：必须不满足所有条件。
	Relation *RelationModelRelation `json:"relation,omitempty"`
}

func (o RelationModel) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RelationModel struct{}"
	}

	return strings.Join([]string{"RelationModel", string(data)}, " ")
}

type RelationModelRelation struct {
	value string
}

type RelationModelRelationEnum struct {
	AND RelationModelRelation
	OR  RelationModelRelation
	NOT RelationModelRelation
}

func GetRelationModelRelationEnum() RelationModelRelationEnum {
	return RelationModelRelationEnum{
		AND: RelationModelRelation{
			value: "AND",
		},
		OR: RelationModelRelation{
			value: "OR",
		},
		NOT: RelationModelRelation{
			value: "NOT",
		},
	}
}

func (c RelationModelRelation) Value() string {
	return c.value
}

func (c RelationModelRelation) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RelationModelRelation) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
