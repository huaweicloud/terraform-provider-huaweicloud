package v2

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
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

// AddAlarmRule 添加阈值规则
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

// AddAlarmRuleInvoker 添加阈值规则
func (c *AomClient) AddAlarmRuleInvoker(request *model.AddAlarmRuleRequest) *AddAlarmRuleInvoker {
	requestDef := GenReqDefForAddAlarmRule()
	return &AddAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddMetricData 添加监控数据
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

// AddMetricDataInvoker 添加监控数据
func (c *AomClient) AddMetricDataInvoker(request *model.AddMetricDataRequest) *AddMetricDataInvoker {
	requestDef := GenReqDefForAddMetricData()
	return &AddMetricDataInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddOrUpdateServiceDiscoveryRules 添加或修改服务发现规则
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

// AddOrUpdateServiceDiscoveryRulesInvoker 添加或修改服务发现规则
func (c *AomClient) AddOrUpdateServiceDiscoveryRulesInvoker(request *model.AddOrUpdateServiceDiscoveryRulesRequest) *AddOrUpdateServiceDiscoveryRulesInvoker {
	requestDef := GenReqDefForAddOrUpdateServiceDiscoveryRules()
	return &AddOrUpdateServiceDiscoveryRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CountEvents 统计事件告警信息
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

// CountEventsInvoker 统计事件告警信息
func (c *AomClient) CountEventsInvoker(request *model.CountEventsRequest) *CountEventsInvoker {
	requestDef := GenReqDefForCountEvents()
	return &CountEventsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAlarmRule 删除阈值规则
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

// DeleteAlarmRuleInvoker 删除阈值规则
func (c *AomClient) DeleteAlarmRuleInvoker(request *model.DeleteAlarmRuleRequest) *DeleteAlarmRuleInvoker {
	requestDef := GenReqDefForDeleteAlarmRule()
	return &DeleteAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAlarmRules 批量删除阈值规则
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

// DeleteAlarmRulesInvoker 批量删除阈值规则
func (c *AomClient) DeleteAlarmRulesInvoker(request *model.DeleteAlarmRulesRequest) *DeleteAlarmRulesInvoker {
	requestDef := GenReqDefForDeleteAlarmRules()
	return &DeleteAlarmRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteserviceDiscoveryRules 删除服务发现规则
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

// DeleteserviceDiscoveryRulesInvoker 删除服务发现规则
func (c *AomClient) DeleteserviceDiscoveryRulesInvoker(request *model.DeleteserviceDiscoveryRulesRequest) *DeleteserviceDiscoveryRulesInvoker {
	requestDef := GenReqDefForDeleteserviceDiscoveryRules()
	return &DeleteserviceDiscoveryRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAlarmRule 查询阈值规则列表
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

// ListAlarmRuleInvoker 查询阈值规则列表
func (c *AomClient) ListAlarmRuleInvoker(request *model.ListAlarmRuleRequest) *ListAlarmRuleInvoker {
	requestDef := GenReqDefForListAlarmRule()
	return &ListAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListEvents 查询事件告警信息
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

// ListEventsInvoker 查询事件告警信息
func (c *AomClient) ListEventsInvoker(request *model.ListEventsRequest) *ListEventsInvoker {
	requestDef := GenReqDefForListEvents()
	return &ListEventsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLogItems 查询日志
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

// ListLogItemsInvoker 查询日志
func (c *AomClient) ListLogItemsInvoker(request *model.ListLogItemsRequest) *ListLogItemsInvoker {
	requestDef := GenReqDefForListLogItems()
	return &ListLogItemsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMetricItems 查询指标
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

// ListMetricItemsInvoker 查询指标
func (c *AomClient) ListMetricItemsInvoker(request *model.ListMetricItemsRequest) *ListMetricItemsInvoker {
	requestDef := GenReqDefForListMetricItems()
	return &ListMetricItemsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSample 查询时序数据
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

// ListSampleInvoker 查询时序数据
func (c *AomClient) ListSampleInvoker(request *model.ListSampleRequest) *ListSampleInvoker {
	requestDef := GenReqDefForListSample()
	return &ListSampleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSeries 查询时间序列
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

// ListSeriesInvoker 查询时间序列
func (c *AomClient) ListSeriesInvoker(request *model.ListSeriesRequest) *ListSeriesInvoker {
	requestDef := GenReqDefForListSeries()
	return &ListSeriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListServiceDiscoveryRules 查询系统中已有服务发现规则
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

// ListServiceDiscoveryRulesInvoker 查询系统中已有服务发现规则
func (c *AomClient) ListServiceDiscoveryRulesInvoker(request *model.ListServiceDiscoveryRulesRequest) *ListServiceDiscoveryRulesInvoker {
	requestDef := GenReqDefForListServiceDiscoveryRules()
	return &ListServiceDiscoveryRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// PushEvents 上报事件告警信息
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

// PushEventsInvoker 上报事件告警信息
func (c *AomClient) PushEventsInvoker(request *model.PushEventsRequest) *PushEventsInvoker {
	requestDef := GenReqDefForPushEvents()
	return &PushEventsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAlarmRule 查询单条阈值规则
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

// ShowAlarmRuleInvoker 查询单条阈值规则
func (c *AomClient) ShowAlarmRuleInvoker(request *model.ShowAlarmRuleRequest) *ShowAlarmRuleInvoker {
	requestDef := GenReqDefForShowAlarmRule()
	return &ShowAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowMetricsData 查询监控数据
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

// ShowMetricsDataInvoker 查询监控数据
func (c *AomClient) ShowMetricsDataInvoker(request *model.ShowMetricsDataRequest) *ShowMetricsDataInvoker {
	requestDef := GenReqDefForShowMetricsData()
	return &ShowMetricsDataInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAlarmRule 修改阈值规则
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

// UpdateAlarmRuleInvoker 修改阈值规则
func (c *AomClient) UpdateAlarmRuleInvoker(request *model.UpdateAlarmRuleRequest) *UpdateAlarmRuleInvoker {
	requestDef := GenReqDefForUpdateAlarmRule()
	return &UpdateAlarmRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstantQueryAomPromGet 瞬时数据查询
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

// ListInstantQueryAomPromGetInvoker 瞬时数据查询
func (c *AomClient) ListInstantQueryAomPromGetInvoker(request *model.ListInstantQueryAomPromGetRequest) *ListInstantQueryAomPromGetInvoker {
	requestDef := GenReqDefForListInstantQueryAomPromGet()
	return &ListInstantQueryAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstantQueryAomPromPost 瞬时数据查询
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

// ListInstantQueryAomPromPostInvoker 瞬时数据查询
func (c *AomClient) ListInstantQueryAomPromPostInvoker(request *model.ListInstantQueryAomPromPostRequest) *ListInstantQueryAomPromPostInvoker {
	requestDef := GenReqDefForListInstantQueryAomPromPost()
	return &ListInstantQueryAomPromPostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLabelValuesAomPromGet 查询标签值
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

// ListLabelValuesAomPromGetInvoker 查询标签值
func (c *AomClient) ListLabelValuesAomPromGetInvoker(request *model.ListLabelValuesAomPromGetRequest) *ListLabelValuesAomPromGetInvoker {
	requestDef := GenReqDefForListLabelValuesAomPromGet()
	return &ListLabelValuesAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLabelsAomPromGet 获取标签名列表
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

// ListLabelsAomPromGetInvoker 获取标签名列表
func (c *AomClient) ListLabelsAomPromGetInvoker(request *model.ListLabelsAomPromGetRequest) *ListLabelsAomPromGetInvoker {
	requestDef := GenReqDefForListLabelsAomPromGet()
	return &ListLabelsAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListLabelsAomPromPost 获取标签名列表
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

// ListLabelsAomPromPostInvoker 获取标签名列表
func (c *AomClient) ListLabelsAomPromPostInvoker(request *model.ListLabelsAomPromPostRequest) *ListLabelsAomPromPostInvoker {
	requestDef := GenReqDefForListLabelsAomPromPost()
	return &ListLabelsAomPromPostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMetadataAomPromGet 元数据查询
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

// ListMetadataAomPromGetInvoker 元数据查询
func (c *AomClient) ListMetadataAomPromGetInvoker(request *model.ListMetadataAomPromGetRequest) *ListMetadataAomPromGetInvoker {
	requestDef := GenReqDefForListMetadataAomPromGet()
	return &ListMetadataAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRangeQueryAomPromGet 区间数据查询
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

// ListRangeQueryAomPromGetInvoker 区间数据查询
func (c *AomClient) ListRangeQueryAomPromGetInvoker(request *model.ListRangeQueryAomPromGetRequest) *ListRangeQueryAomPromGetInvoker {
	requestDef := GenReqDefForListRangeQueryAomPromGet()
	return &ListRangeQueryAomPromGetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRangeQueryAomPromPost 区间数据查询
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

// ListRangeQueryAomPromPostInvoker 区间数据查询
func (c *AomClient) ListRangeQueryAomPromPostInvoker(request *model.ListRangeQueryAomPromPostRequest) *ListRangeQueryAomPromPostInvoker {
	requestDef := GenReqDefForListRangeQueryAomPromPost()
	return &ListRangeQueryAomPromPostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
