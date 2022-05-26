package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListPredefineTagsRequest struct {

	// 键，支持模糊查询，不区分大小写，如果包含“non-URL-safe”的字符，需要进行“urlencoded”。
	Key *string `json:"key,omitempty"`

	// 值，支持模糊查询，不区分大小写，如果包含“non-URL-safe”的字符，需要进行“urlencoded”。
	Value *string `json:"value,omitempty"`

	// 查询记录数。 最小为1，最大为1000，未输入时默认为10，为0时不限制查询数据条数。
	Limit *int32 `json:"limit,omitempty"`

	// 分页位置标识（索引）。 从marker指定索引的下一条数据开始查询。 说明： 查询第一页数据时，不需要传入此参数，查询后续页码数据时，将查询前一页数据响应体中marker值配入此参数，当返回的tags为空列表时表示查询到最后一页。
	Marker *string `json:"marker,omitempty"`

	// 排序字段： 可输入的值包含（区分大小写）：update_time（更新时间）、key（键）、value（值）。 只能选择以上排序字段中的一个，并按照排序方法字段order_method进行排序，如果不传则默认值为：update_time。 如以下： 若该字段为update_time，则剩余两个默认字段排序为key升序，value升序。 若该字段如为key，则剩余两个默认字段排序为update_time降序，value升序。 若该字段如为value，则剩余两个默认字段排序为update_time降序，key升序。 若该字段不传，默认字段为update_time，则剩余两个默认字段排序为key升序，value升序。
	OrderField *string `json:"order_field,omitempty"`

	// order_field字段的排序方法。 可输入的值包含（区分大小写）： asc（升序） desc（降序） 只能选择以上值的其中之一。 不传则默认值为：desc
	OrderMethod *ListPredefineTagsRequestOrderMethod `json:"order_method,omitempty"`
}

func (o ListPredefineTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPredefineTagsRequest struct{}"
	}

	return strings.Join([]string{"ListPredefineTagsRequest", string(data)}, " ")
}

type ListPredefineTagsRequestOrderMethod struct {
	value string
}

type ListPredefineTagsRequestOrderMethodEnum struct {
	ASC  ListPredefineTagsRequestOrderMethod
	DESC ListPredefineTagsRequestOrderMethod
}

func GetListPredefineTagsRequestOrderMethodEnum() ListPredefineTagsRequestOrderMethodEnum {
	return ListPredefineTagsRequestOrderMethodEnum{
		ASC: ListPredefineTagsRequestOrderMethod{
			value: "asc",
		},
		DESC: ListPredefineTagsRequestOrderMethod{
			value: "desc",
		},
	}
}

func (c ListPredefineTagsRequestOrderMethod) Value() string {
	return c.value
}

func (c ListPredefineTagsRequestOrderMethod) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListPredefineTagsRequestOrderMethod) UnmarshalJSON(b []byte) error {
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
