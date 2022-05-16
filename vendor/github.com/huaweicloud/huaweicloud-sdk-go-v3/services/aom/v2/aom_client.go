package v2

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"
)

type AomClient struct {
	HcClient *http_client.HcHttpClient
}

func NewAomClient(hcClient *http_client.HcHttpClient) *AomClient {
	return &AomClient{HcClient: hcClient}
}

func AomClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// 添加阈值规则
//
// 该接口用于添加一条阈值规则。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) AddAlarmRule(request *model.AddAlarmRuleRequest) (*model.AddAlarmRuleResponse, error) {
	requestDef := GenReqDefForAddAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddAlarmRuleResponse), nil
	}
}

// 添加监控数据
//
// 该接口用于向服务端添加一条或多条监控数据。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) AddMetricData(request *model.AddMetricDataRequest) (*model.AddMetricDataResponse, error) {
	requestDef := GenReqDefForAddMetricData()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddMetricDataResponse), nil
	}
}

// 添加或修改服务发现规则
//
// 该接口用于添加或修改一条或多条服务发现规则。同一projectid下可添加的规则上限为100条。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) AddOrUpdateServiceDiscoveryRules(request *model.AddOrUpdateServiceDiscoveryRulesRequest) (*model.AddOrUpdateServiceDiscoveryRulesResponse, error) {
	requestDef := GenReqDefForAddOrUpdateServiceDiscoveryRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddOrUpdateServiceDiscoveryRulesResponse), nil
	}
}

// 统计事件告警信息
//
// 该接口用于分段统计指定条件下的事件、告警。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) CountEvents(request *model.CountEventsRequest) (*model.CountEventsResponse, error) {
	requestDef := GenReqDefForCountEvents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CountEventsResponse), nil
	}
}

// 删除阈值规则
//
// 该接口用于删除阈值规则。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) DeleteAlarmRule(request *model.DeleteAlarmRuleRequest) (*model.DeleteAlarmRuleResponse, error) {
	requestDef := GenReqDefForDeleteAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAlarmRuleResponse), nil
	}
}

// 批量删除阈值规则
//
// 批量删除阈值规则
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) DeleteAlarmRules(request *model.DeleteAlarmRulesRequest) (*model.DeleteAlarmRulesResponse, error) {
	requestDef := GenReqDefForDeleteAlarmRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAlarmRulesResponse), nil
	}
}

// 删除服务发现规则
//
// 该接口用于删除服务发现规则。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) DeleteserviceDiscoveryRules(request *model.DeleteserviceDiscoveryRulesRequest) (*model.DeleteserviceDiscoveryRulesResponse, error) {
	requestDef := GenReqDefForDeleteserviceDiscoveryRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteserviceDiscoveryRulesResponse), nil
	}
}

// 查询阈值规则列表
//
// 该接口用于查询阈值规则列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListAlarmRule(request *model.ListAlarmRuleRequest) (*model.ListAlarmRuleResponse, error) {
	requestDef := GenReqDefForListAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAlarmRuleResponse), nil
	}
}

// 查询事件告警信息
//
// 该接口用于查询对应用户的事件、告警。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListEvents(request *model.ListEventsRequest) (*model.ListEventsResponse, error) {
	requestDef := GenReqDefForListEvents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListEventsResponse), nil
	}
}

// 查询日志
//
// 该接口用于查询不同维度(例如集群、IP、应用等)下的日志内容，支持分页查询。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListLogItems(request *model.ListLogItemsRequest) (*model.ListLogItemsResponse, error) {
	requestDef := GenReqDefForListLogItems()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLogItemsResponse), nil
	}
}

// 查询指标
//
// 该接口用于查询系统当前可监控的指标列表，可以指定指标命名空间、指标名称、维度、所属资源的编号（格式为：resType_resId），分页查询的起始位置和返回的最大记录条数。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListMetricItems(request *model.ListMetricItemsRequest) (*model.ListMetricItemsResponse, error) {
	requestDef := GenReqDefForListMetricItems()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMetricItemsResponse), nil
	}
}

// 查询时序数据
//
// 该接口用于查询指定时间范围内的监控时序数据，可以通过参数指定需要查询的数据维度，数据周期等。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListSample(request *model.ListSampleRequest) (*model.ListSampleResponse, error) {
	requestDef := GenReqDefForListSample()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSampleResponse), nil
	}
}

// 查询时间序列
//
// 该接口用于查询系统当前可监控的时间序列列表，可以指定时间序列命名空间、名称、维度、所属资源的编号（格式为：resType_resId），分页查询的起始位置和返回的最大记录条数。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListSeries(request *model.ListSeriesRequest) (*model.ListSeriesResponse, error) {
	requestDef := GenReqDefForListSeries()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSeriesResponse), nil
	}
}

// 查询系统中已有服务发现规则
//
// 该接口用于查询系统当前已存在的服务发现规则。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListServiceDiscoveryRules(request *model.ListServiceDiscoveryRulesRequest) (*model.ListServiceDiscoveryRulesResponse, error) {
	requestDef := GenReqDefForListServiceDiscoveryRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListServiceDiscoveryRulesResponse), nil
	}
}

// 上报事件告警信息
//
// 该接口用于上报对应用户的事件、告警。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) PushEvents(request *model.PushEventsRequest) (*model.PushEventsResponse, error) {
	requestDef := GenReqDefForPushEvents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.PushEventsResponse), nil
	}
}

// 查询单条阈值规则
//
// 该接口用于查询单条阈值规则。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ShowAlarmRule(request *model.ShowAlarmRuleRequest) (*model.ShowAlarmRuleResponse, error) {
	requestDef := GenReqDefForShowAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAlarmRuleResponse), nil
	}
}

// 查询监控数据
//
// 该接口用于查询指定时间范围内指标的监控数据，可以通过参数指定需要查询的数据维度，数据周期等。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ShowMetricsData(request *model.ShowMetricsDataRequest) (*model.ShowMetricsDataResponse, error) {
	requestDef := GenReqDefForShowMetricsData()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowMetricsDataResponse), nil
	}
}

// 修改阈值规则
//
// 该接口用于修改一条阈值规则。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) UpdateAlarmRule(request *model.UpdateAlarmRuleRequest) (*model.UpdateAlarmRuleResponse, error) {
	requestDef := GenReqDefForUpdateAlarmRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAlarmRuleResponse), nil
	}
}

// 瞬时数据查询
//
// 该接口用于查询PromQL(Prometheus Query Language)。 在特定时间点下的计算结果。（注：接口目前开放的region为：北京四、上海一和广州）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListInstantQueryAomPromGet(request *model.ListInstantQueryAomPromGetRequest) (*model.ListInstantQueryAomPromGetResponse, error) {
	requestDef := GenReqDefForListInstantQueryAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstantQueryAomPromGetResponse), nil
	}
}

// 瞬时数据查询
//
// 该接口用于查询PromQL(Prometheus Query Language) 在特定时间点下的计算结果。（注：接口目前开放的region为：北京四、上海一和广州）
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListInstantQueryAomPromPost(request *model.ListInstantQueryAomPromPostRequest) (*model.ListInstantQueryAomPromPostResponse, error) {
	requestDef := GenReqDefForListInstantQueryAomPromPost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstantQueryAomPromPostResponse), nil
	}
}

// 查询标签值
//
// 该接口用于查询带有指定标签的时间序列列表。（注：接口目前开放的region为：北京四、上海一和广州）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListLabelValuesAomPromGet(request *model.ListLabelValuesAomPromGetRequest) (*model.ListLabelValuesAomPromGetResponse, error) {
	requestDef := GenReqDefForListLabelValuesAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLabelValuesAomPromGetResponse), nil
	}
}

// 获取标签名列表
//
// 该接口用于获取标签名列表。（注：接口目前开放的region为：北京四、上海一和广州）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListLabelsAomPromGet(request *model.ListLabelsAomPromGetRequest) (*model.ListLabelsAomPromGetResponse, error) {
	requestDef := GenReqDefForListLabelsAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLabelsAomPromGetResponse), nil
	}
}

// 获取标签名列表
//
// 该接口用于获取标签名列表。（注：接口目前开放的region为：北京四、上海一和广州）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListLabelsAomPromPost(request *model.ListLabelsAomPromPostRequest) (*model.ListLabelsAomPromPostResponse, error) {
	requestDef := GenReqDefForListLabelsAomPromPost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListLabelsAomPromPostResponse), nil
	}
}

// 元数据查询
//
// 该接口用于查询序列及序列标签的元数据。（注：接口目前开放的region为：北京四、上海一和广州）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListMetadataAomPromGet(request *model.ListMetadataAomPromGetRequest) (*model.ListMetadataAomPromGetResponse, error) {
	requestDef := GenReqDefForListMetadataAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMetadataAomPromGetResponse), nil
	}
}

// 区间数据查询
//
// 该接口用于查询PromQL(Prometheus Query Language)在一段时间返回内的计算结果。（注：接口目前开放的region为：北京四、上海一和广州）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListRangeQueryAomPromGet(request *model.ListRangeQueryAomPromGetRequest) (*model.ListRangeQueryAomPromGetResponse, error) {
	requestDef := GenReqDefForListRangeQueryAomPromGet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRangeQueryAomPromGetResponse), nil
	}
}

// 区间数据查询
//
// 该接口用于查询PromQL(Prometheus Query Language)在一段时间返回内的计算结果。（注：接口目前开放的region为：北京四、上海一和广州）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *AomClient) ListRangeQueryAomPromPost(request *model.ListRangeQueryAomPromPostRequest) (*model.ListRangeQueryAomPromPostResponse, error) {
	requestDef := GenReqDefForListRangeQueryAomPromPost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRangeQueryAomPromPostResponse), nil
	}
}
