package v1

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1/model"
)

type CdnClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewCdnClient(hcClient *httpclient.HcHttpClient) *CdnClient {
	return &CdnClient{HcClient: hcClient}
}

func CdnClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder().WithCredentialsType("global.Credentials")
	return builder
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// BatchDeleteTags 删除资源标签配置接口
//
// 用于删除资源标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) BatchDeleteTags(request *model.BatchDeleteTagsRequest) (*model.BatchDeleteTagsResponse, error) {
	requestDef := GenReqDefForBatchDeleteTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchDeleteTagsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// BatchDeleteTagsInvoker 删除资源标签配置接口
func (c *CdnClient) BatchDeleteTagsInvoker(request *model.BatchDeleteTagsRequest) *BatchDeleteTagsInvoker {
	requestDef := GenReqDefForBatchDeleteTags()
	return &BatchDeleteTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// CreateDomain 创建加速域名
//
// 创建加速域名。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) CreateDomain(request *model.CreateDomainRequest) (*model.CreateDomainResponse, error) {
	requestDef := GenReqDefForCreateDomain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDomainResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// CreateDomainInvoker 创建加速域名
func (c *CdnClient) CreateDomainInvoker(request *model.CreateDomainRequest) *CreateDomainInvoker {
	requestDef := GenReqDefForCreateDomain()
	return &CreateDomainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// CreatePreheatingTasks 创建预热缓存任务
//
// 创建预热任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) CreatePreheatingTasks(request *model.CreatePreheatingTasksRequest) (*model.CreatePreheatingTasksResponse, error) {
	requestDef := GenReqDefForCreatePreheatingTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePreheatingTasksResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// CreatePreheatingTasksInvoker 创建预热缓存任务
func (c *CdnClient) CreatePreheatingTasksInvoker(request *model.CreatePreheatingTasksRequest) *CreatePreheatingTasksInvoker {
	requestDef := GenReqDefForCreatePreheatingTasks()
	return &CreatePreheatingTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// CreateRefreshTasks 创建刷新缓存任务
//
// 创建刷新缓存任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) CreateRefreshTasks(request *model.CreateRefreshTasksRequest) (*model.CreateRefreshTasksResponse, error) {
	requestDef := GenReqDefForCreateRefreshTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRefreshTasksResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// CreateRefreshTasksInvoker 创建刷新缓存任务
func (c *CdnClient) CreateRefreshTasksInvoker(request *model.CreateRefreshTasksRequest) *CreateRefreshTasksInvoker {
	requestDef := GenReqDefForCreateRefreshTasks()
	return &CreateRefreshTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// CreateTags 创建资源标签配置接口
//
// 用于创建资源标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) CreateTags(request *model.CreateTagsRequest) (*model.CreateTagsResponse, error) {
	requestDef := GenReqDefForCreateTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTagsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// CreateTagsInvoker 创建资源标签配置接口
func (c *CdnClient) CreateTagsInvoker(request *model.CreateTagsRequest) *CreateTagsInvoker {
	requestDef := GenReqDefForCreateTags()
	return &CreateTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// DeleteDomain 删除加速域名
//
// 删除加速域名。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) DeleteDomain(request *model.DeleteDomainRequest) (*model.DeleteDomainResponse, error) {
	requestDef := GenReqDefForDeleteDomain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDomainResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// DeleteDomainInvoker 删除加速域名
func (c *CdnClient) DeleteDomainInvoker(request *model.DeleteDomainRequest) *DeleteDomainInvoker {
	requestDef := GenReqDefForDeleteDomain()
	return &DeleteDomainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// DisableDomain 停用加速域名
//
// 停用加速域名。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) DisableDomain(request *model.DisableDomainRequest) (*model.DisableDomainResponse, error) {
	requestDef := GenReqDefForDisableDomain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DisableDomainResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// DisableDomainInvoker 停用加速域名
func (c *CdnClient) DisableDomainInvoker(request *model.DisableDomainRequest) *DisableDomainInvoker {
	requestDef := GenReqDefForDisableDomain()
	return &DisableDomainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// EnableDomain 启用加速域名
//
// 启用加速域名。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) EnableDomain(request *model.EnableDomainRequest) (*model.EnableDomainResponse, error) {
	requestDef := GenReqDefForEnableDomain()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.EnableDomainResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// EnableDomainInvoker 启用加速域名
func (c *CdnClient) EnableDomainInvoker(request *model.EnableDomainRequest) *EnableDomainInvoker {
	requestDef := GenReqDefForEnableDomain()
	return &EnableDomainInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ListDomains 查询加速域名
//
// 查询加速域名信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ListDomains(request *model.ListDomainsRequest) (*model.ListDomainsResponse, error) {
	requestDef := GenReqDefForListDomains()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDomainsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ListDomainsInvoker 查询加速域名
func (c *CdnClient) ListDomainsInvoker(request *model.ListDomainsRequest) *ListDomainsInvoker {
	requestDef := GenReqDefForListDomains()
	return &ListDomainsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowBlackWhiteList 查询IP黑白名单
//
// 查询域名已经设置的IP黑白名单。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowBlackWhiteList(request *model.ShowBlackWhiteListRequest) (*model.ShowBlackWhiteListResponse, error) {
	requestDef := GenReqDefForShowBlackWhiteList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBlackWhiteListResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowBlackWhiteListInvoker 查询IP黑白名单
func (c *CdnClient) ShowBlackWhiteListInvoker(request *model.ShowBlackWhiteListRequest) *ShowBlackWhiteListInvoker {
	requestDef := GenReqDefForShowBlackWhiteList()
	return &ShowBlackWhiteListInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowCacheRules 查询缓存规则
//
// 查询缓存规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowCacheRules(request *model.ShowCacheRulesRequest) (*model.ShowCacheRulesResponse, error) {
	requestDef := GenReqDefForShowCacheRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCacheRulesResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowCacheRulesInvoker 查询缓存规则
func (c *CdnClient) ShowCacheRulesInvoker(request *model.ShowCacheRulesRequest) *ShowCacheRulesInvoker {
	requestDef := GenReqDefForShowCacheRules()
	return &ShowCacheRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowCertificatesHttpsInfo 查询所有绑定HTTPS证书的域名信息
//
// 查询所有绑定HTTPS证书的域名信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowCertificatesHttpsInfo(request *model.ShowCertificatesHttpsInfoRequest) (*model.ShowCertificatesHttpsInfoResponse, error) {
	requestDef := GenReqDefForShowCertificatesHttpsInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCertificatesHttpsInfoResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowCertificatesHttpsInfoInvoker 查询所有绑定HTTPS证书的域名信息
func (c *CdnClient) ShowCertificatesHttpsInfoInvoker(request *model.ShowCertificatesHttpsInfoRequest) *ShowCertificatesHttpsInfoInvoker {
	requestDef := GenReqDefForShowCertificatesHttpsInfo()
	return &ShowCertificatesHttpsInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainDetail 查询加速域名详情
//
// 查询加速域名详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowDomainDetail(request *model.ShowDomainDetailRequest) (*model.ShowDomainDetailResponse, error) {
	requestDef := GenReqDefForShowDomainDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainDetailResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainDetailInvoker 查询加速域名详情
func (c *CdnClient) ShowDomainDetailInvoker(request *model.ShowDomainDetailRequest) *ShowDomainDetailInvoker {
	requestDef := GenReqDefForShowDomainDetail()
	return &ShowDomainDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainFullConfig 查询域名配置接口
//
// 查询域名配置接口，
// 支持查询回源请求头、HTTP header配置、URL鉴权、证书、源站、回源协议、强制重定向、智能压缩、IPv6开关、状态码缓存时间、Range回源、User-Agent黑白名单、改写回源URL、自定义错误页面、缓存规则、IP黑白名单、防盗链、强制跳转。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowDomainFullConfig(request *model.ShowDomainFullConfigRequest) (*model.ShowDomainFullConfigResponse, error) {
	requestDef := GenReqDefForShowDomainFullConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainFullConfigResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainFullConfigInvoker 查询域名配置接口
func (c *CdnClient) ShowDomainFullConfigInvoker(request *model.ShowDomainFullConfigRequest) *ShowDomainFullConfigInvoker {
	requestDef := GenReqDefForShowDomainFullConfig()
	return &ShowDomainFullConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainItemDetails 批量查询域名的统计明细-按域名单独返回
//
// - 支持查询90天内的数据。
// - 查询跨度不能超过7天。
// - 最多同时指定100个域名。
// - 起始时间和结束时间，左闭右开，需要同时指定。
// - 开始时间、结束时间必须传毫秒级时间戳，且必须为5分钟整时刻点，如：0分、5分、10分、15分等，如果传的不是5分钟时刻点，返回数据可能与预期不一致。
// - 统一用开始时间表示一个时间段，如：2019-01-24 20:15:00 表示取 [20:15:00, 20:20:00)的统计数据，且左闭右开。
// - 流量类指标单位统一为Byte（字节）、带宽类指标单位统一为bit/s（比特/秒）、请求数类指标单位统一为次数。用于查询指定域名、指定统计指标的明细数据。
// - 如果传的是多个域名，则每个域名的数据分开返回。
// - 支持同时查询多个指标，不超过10个。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowDomainItemDetails(request *model.ShowDomainItemDetailsRequest) (*model.ShowDomainItemDetailsResponse, error) {
	requestDef := GenReqDefForShowDomainItemDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainItemDetailsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainItemDetailsInvoker 批量查询域名的统计明细-按域名单独返回
func (c *CdnClient) ShowDomainItemDetailsInvoker(request *model.ShowDomainItemDetailsRequest) *ShowDomainItemDetailsInvoker {
	requestDef := GenReqDefForShowDomainItemDetails()
	return &ShowDomainItemDetailsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainItemLocationDetails 批量查询域名的区域、运营商统计明细-按域名单独返回
//
// - 支持查询90天内的数据。
// - 查询跨度是7天。
// - 最多同时指定100个域名。
// - 起始时间和结束时间，左闭右开，需要同时指定。
// - 开始时间、结束时间必须传毫秒级时间戳，且必须为5分钟整时刻点，如：0分、5分、10分、15分等，如果传的不是5分钟时刻点，返回数据可能与预期不一致。
// - 统一用开始时间表示一个时间段，如：2019-01-24 20:15:00 表示取 [20:15:00, 20:20:00)的统计数据，且左闭右开。
// - 流量类指标单位统一为Byte（字节）、带宽类指标单位统一为bit/s（比特/秒）、请求数类指标单位统一为次数。
// - 用于查询指定域名、指定统计指标的明细数据。
// - 如果传的是多个域名，则每个域名的数据分开返回。
// - 支持按区域、运营商维度查询统计数据, 回源指标除外。
// - 支持同时查询多个指标，不超过10个。
// - 域名为海外加速场景不适用。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowDomainItemLocationDetails(request *model.ShowDomainItemLocationDetailsRequest) (*model.ShowDomainItemLocationDetailsResponse, error) {
	requestDef := GenReqDefForShowDomainItemLocationDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainItemLocationDetailsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainItemLocationDetailsInvoker 批量查询域名的区域、运营商统计明细-按域名单独返回
func (c *CdnClient) ShowDomainItemLocationDetailsInvoker(request *model.ShowDomainItemLocationDetailsRequest) *ShowDomainItemLocationDetailsInvoker {
	requestDef := GenReqDefForShowDomainItemLocationDetails()
	return &ShowDomainItemLocationDetailsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainLocationStats 按区域运营商查询域名统计数据
//
// - 支持查询90天内的数据。
//
// - 支持多指标同时查询，不超过5个。
//
// - 最多同时指定20个域名。
//
// - 起始时间和结束时间需要同时指定，左闭右开，毫秒级时间戳，且时间点必须为与查询时间间隔参数匹配的整时刻点。比如查询时间间隔为5分钟时，起始时间和结束时间必须为5分钟整时刻点，如：0分、5分、10分、15分等，如果时间点与时间间隔不匹配，返回数据可能与预期不一致。统一用开始时间表示一个时间段，如：2019-01-24 20:15:00 表示取 [20:15:00, 20:20:00)的统计数据，且左闭右开。
//
// - action取值：location_detail,location_summary
//
// - 流量类指标单位统一为Byte（字节）、带宽类指标单位统一为bit/s（比特/秒）、请求数类和状态码类指标单位统一为次数。用于查询指定域名、指定统计指标的区域运营商明细数据。
//
// - 单租户调用频率：15次/s。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowDomainLocationStats(request *model.ShowDomainLocationStatsRequest) (*model.ShowDomainLocationStatsResponse, error) {
	requestDef := GenReqDefForShowDomainLocationStats()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainLocationStatsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainLocationStatsInvoker 按区域运营商查询域名统计数据
func (c *CdnClient) ShowDomainLocationStatsInvoker(request *model.ShowDomainLocationStatsRequest) *ShowDomainLocationStatsInvoker {
	requestDef := GenReqDefForShowDomainLocationStats()
	return &ShowDomainLocationStatsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainStats 查询域名统计数据
//
// - 支持查询90天内的数据。
//
// - 支持多指标同时查询，不超过5个。
//
// - 最多同时指定20个域名。
//
// - 起始时间和结束时间需要同时指定，左闭右开，毫秒级时间戳，且时间点必须为与查询时间间隔参数匹配的整时刻点。比如查询时间间隔为5分钟时，起始时间和结束时间必须为5分钟整时刻点，如：0分、5分、10分、15分等，如果时间点与时间间隔不匹配，返回数据可能与预期不一致。统一用开始时间表示一个时间段，如：2019-01-24 20:15:00 表示取 [20:15:00, 20:20:00)的统计数据，且左闭右开。
//
// - action取值：detail,summary
//
// - 流量类指标单位统一为Byte（字节）、带宽类指标单位统一为bit/s（比特/秒）、请求数类和状态码类指标单位统一为次数。用于查询指定域名、指定统计指标的明细数据。
//
// - 单租户调用频率：15次/s。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowDomainStats(request *model.ShowDomainStatsRequest) (*model.ShowDomainStatsResponse, error) {
	requestDef := GenReqDefForShowDomainStats()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainStatsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowDomainStatsInvoker 查询域名统计数据
func (c *CdnClient) ShowDomainStatsInvoker(request *model.ShowDomainStatsRequest) *ShowDomainStatsInvoker {
	requestDef := GenReqDefForShowDomainStats()
	return &ShowDomainStatsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowHistoryTaskDetails 查询刷新预热任务详情
//
// 查询刷新预热任务详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowHistoryTaskDetails(request *model.ShowHistoryTaskDetailsRequest) (*model.ShowHistoryTaskDetailsResponse, error) {
	requestDef := GenReqDefForShowHistoryTaskDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowHistoryTaskDetailsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowHistoryTaskDetailsInvoker 查询刷新预热任务详情
func (c *CdnClient) ShowHistoryTaskDetailsInvoker(request *model.ShowHistoryTaskDetailsRequest) *ShowHistoryTaskDetailsInvoker {
	requestDef := GenReqDefForShowHistoryTaskDetails()
	return &ShowHistoryTaskDetailsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowHistoryTasks 查询刷新预热任务
//
// 查询刷新预热任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowHistoryTasks(request *model.ShowHistoryTasksRequest) (*model.ShowHistoryTasksResponse, error) {
	requestDef := GenReqDefForShowHistoryTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowHistoryTasksResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowHistoryTasksInvoker 查询刷新预热任务
func (c *CdnClient) ShowHistoryTasksInvoker(request *model.ShowHistoryTasksRequest) *ShowHistoryTasksInvoker {
	requestDef := GenReqDefForShowHistoryTasks()
	return &ShowHistoryTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowHttpInfo 查询HTTPS配置
//
// 获取加速域名证书。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowHttpInfo(request *model.ShowHttpInfoRequest) (*model.ShowHttpInfoResponse, error) {
	requestDef := GenReqDefForShowHttpInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowHttpInfoResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowHttpInfoInvoker 查询HTTPS配置
func (c *CdnClient) ShowHttpInfoInvoker(request *model.ShowHttpInfoRequest) *ShowHttpInfoInvoker {
	requestDef := GenReqDefForShowHttpInfo()
	return &ShowHttpInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowIpInfo 查询IP归属信息
//
// 查询IP归属信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowIpInfo(request *model.ShowIpInfoRequest) (*model.ShowIpInfoResponse, error) {
	requestDef := GenReqDefForShowIpInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowIpInfoResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowIpInfoInvoker 查询IP归属信息
func (c *CdnClient) ShowIpInfoInvoker(request *model.ShowIpInfoRequest) *ShowIpInfoInvoker {
	requestDef := GenReqDefForShowIpInfo()
	return &ShowIpInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowLogs 日志查询
//
// 查询日志下载链接，支持查询30天内的日志信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowLogs(request *model.ShowLogsRequest) (*model.ShowLogsResponse, error) {
	requestDef := GenReqDefForShowLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowLogsResponse), nil
	}
}

// ShowLogsInvoker 日志查询
func (c *CdnClient) ShowLogsInvoker(request *model.ShowLogsRequest) *ShowLogsInvoker {
	requestDef := GenReqDefForShowLogs()
	return &ShowLogsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowOriginHost 查询回源HOST
//
// 查询回源HOST。回源HOST是CDN节点在回源过程中，在源站访问的站点域名，即http请求头中的host信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowOriginHost(request *model.ShowOriginHostRequest) (*model.ShowOriginHostResponse, error) {
	requestDef := GenReqDefForShowOriginHost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowOriginHostResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowOriginHostInvoker 查询回源HOST
func (c *CdnClient) ShowOriginHostInvoker(request *model.ShowOriginHostRequest) *ShowOriginHostInvoker {
	requestDef := GenReqDefForShowOriginHost()
	return &ShowOriginHostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowQuota 查询用户配额
//
// 查询当前用户域名、刷新文件、刷新目录和预热的配额
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowQuota(request *model.ShowQuotaRequest) (*model.ShowQuotaResponse, error) {
	requestDef := GenReqDefForShowQuota()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowQuotaResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowQuotaInvoker 查询用户配额
func (c *CdnClient) ShowQuotaInvoker(request *model.ShowQuotaRequest) *ShowQuotaInvoker {
	requestDef := GenReqDefForShowQuota()
	return &ShowQuotaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowRefer 查询Referer过滤规则
//
// 查询Referer过滤规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowRefer(request *model.ShowReferRequest) (*model.ShowReferResponse, error) {
	requestDef := GenReqDefForShowRefer()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowReferResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowReferInvoker 查询Referer过滤规则
func (c *CdnClient) ShowReferInvoker(request *model.ShowReferRequest) *ShowReferInvoker {
	requestDef := GenReqDefForShowRefer()
	return &ShowReferInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowResponseHeader 查询响应头配置
//
// 列举header所有配置。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowResponseHeader(request *model.ShowResponseHeaderRequest) (*model.ShowResponseHeaderResponse, error) {
	requestDef := GenReqDefForShowResponseHeader()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowResponseHeaderResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowResponseHeaderInvoker 查询响应头配置
func (c *CdnClient) ShowResponseHeaderInvoker(request *model.ShowResponseHeaderRequest) *ShowResponseHeaderInvoker {
	requestDef := GenReqDefForShowResponseHeader()
	return &ShowResponseHeaderInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowTags 查询资源标签列表配置接口
//
// 用于查询资源标签列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowTags(request *model.ShowTagsRequest) (*model.ShowTagsResponse, error) {
	requestDef := GenReqDefForShowTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTagsResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowTagsInvoker 查询资源标签列表配置接口
func (c *CdnClient) ShowTagsInvoker(request *model.ShowTagsRequest) *ShowTagsInvoker {
	requestDef := GenReqDefForShowTags()
	return &ShowTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowTopUrl 查询TOP100 URL明细
//
// - 查询TOP100 URL明细。
//
// - 支持查询90天内的数据。
//
// - 查询跨度不能超过31天。
//
// - 起始时间和结束时间，左闭右开，需要同时指定。如查询2021-10-24 00:00:00 到 2021-10-25 00:00:00
// 的数据，表示取 [2021-10-24 00:00:00, 2021-10-25 00:00:00)的统计数据。
//
// - 开始时间、结束时间必须传毫秒级时间戳，且必须为凌晨0点整时刻点，如果传的不是凌晨0点整时刻点，返回数据可能与预期不一致。
//
// - 流量类指标单位统一为Byte（字节）、请求数类指标单位统一为次数。用于查询指定域名、指定统计指标的明细数据。
//
// - 单租户调用频率：5次/s。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowTopUrl(request *model.ShowTopUrlRequest) (*model.ShowTopUrlResponse, error) {
	requestDef := GenReqDefForShowTopUrl()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTopUrlResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowTopUrlInvoker 查询TOP100 URL明细
func (c *CdnClient) ShowTopUrlInvoker(request *model.ShowTopUrlRequest) *ShowTopUrlInvoker {
	requestDef := GenReqDefForShowTopUrl()
	return &ShowTopUrlInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowUrlTaskInfo 查询刷新预热URL记录
//
// 查询刷新预热URL记录。如需此接口，请提交工单开通。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) ShowUrlTaskInfo(request *model.ShowUrlTaskInfoRequest) (*model.ShowUrlTaskInfoResponse, error) {
	requestDef := GenReqDefForShowUrlTaskInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowUrlTaskInfoResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// ShowUrlTaskInfoInvoker 查询刷新预热URL记录
func (c *CdnClient) ShowUrlTaskInfoInvoker(request *model.ShowUrlTaskInfoRequest) *ShowUrlTaskInfoInvoker {
	requestDef := GenReqDefForShowUrlTaskInfo()
	return &ShowUrlTaskInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateBlackWhiteList 设置IP黑白名单
//
// 设置域名的IP黑白名单。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateBlackWhiteList(request *model.UpdateBlackWhiteListRequest) (*model.UpdateBlackWhiteListResponse, error) {
	requestDef := GenReqDefForUpdateBlackWhiteList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateBlackWhiteListResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateBlackWhiteListInvoker 设置IP黑白名单
func (c *CdnClient) UpdateBlackWhiteListInvoker(request *model.UpdateBlackWhiteListRequest) *UpdateBlackWhiteListInvoker {
	requestDef := GenReqDefForUpdateBlackWhiteList()
	return &UpdateBlackWhiteListInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateCacheRules 设置缓存规则
//
// 设置CDN节点上缓存资源的缓存策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateCacheRules(request *model.UpdateCacheRulesRequest) (*model.UpdateCacheRulesResponse, error) {
	requestDef := GenReqDefForUpdateCacheRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCacheRulesResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateCacheRulesInvoker 设置缓存规则
func (c *CdnClient) UpdateCacheRulesInvoker(request *model.UpdateCacheRulesRequest) *UpdateCacheRulesInvoker {
	requestDef := GenReqDefForUpdateCacheRules()
	return &UpdateCacheRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateDomainFullConfig 修改域名全量配置接口
//
// 修改域名配置接口，
// 支持修改回源请求头、HTTP header配置、URL鉴权、证书、源站、回源协议、强制重定向、智能压缩、IPv6开关、状态码缓存时间、Range回源、User-Agent黑白名单、改写回源URL、自定义错误页面、缓存规则、IP黑白名单、防盗链、强制跳转。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateDomainFullConfig(request *model.UpdateDomainFullConfigRequest) (*model.UpdateDomainFullConfigResponse, error) {
	requestDef := GenReqDefForUpdateDomainFullConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainFullConfigResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateDomainFullConfigInvoker 修改域名全量配置接口
func (c *CdnClient) UpdateDomainFullConfigInvoker(request *model.UpdateDomainFullConfigRequest) *UpdateDomainFullConfigInvoker {
	requestDef := GenReqDefForUpdateDomainFullConfig()
	return &UpdateDomainFullConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateDomainMultiCertificates 一个证书批量设置多个域名
//
// 一个证书配置多个域名，设置域名强制https回源参数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateDomainMultiCertificates(request *model.UpdateDomainMultiCertificatesRequest) (*model.UpdateDomainMultiCertificatesResponse, error) {
	requestDef := GenReqDefForUpdateDomainMultiCertificates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainMultiCertificatesResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateDomainMultiCertificatesInvoker 一个证书批量设置多个域名
func (c *CdnClient) UpdateDomainMultiCertificatesInvoker(request *model.UpdateDomainMultiCertificatesRequest) *UpdateDomainMultiCertificatesInvoker {
	requestDef := GenReqDefForUpdateDomainMultiCertificates()
	return &UpdateDomainMultiCertificatesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateDomainOrigin 修改源站信息。
//
// 修改加速域名的源站配置。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateDomainOrigin(request *model.UpdateDomainOriginRequest) (*model.UpdateDomainOriginResponse, error) {
	requestDef := GenReqDefForUpdateDomainOrigin()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainOriginResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateDomainOriginInvoker 修改源站信息。
func (c *CdnClient) UpdateDomainOriginInvoker(request *model.UpdateDomainOriginRequest) *UpdateDomainOriginInvoker {
	requestDef := GenReqDefForUpdateDomainOrigin()
	return &UpdateDomainOriginInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateFollow302Switch 开启/关闭回源跟随
//
// 开启此项配置后，当CDN节点回源请求源站返回301/302状态码时，CDN节点会先跳转到301/302对应地址获取资源并缓存后再返回给用户。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateFollow302Switch(request *model.UpdateFollow302SwitchRequest) (*model.UpdateFollow302SwitchResponse, error) {
	requestDef := GenReqDefForUpdateFollow302Switch()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateFollow302SwitchResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateFollow302SwitchInvoker 开启/关闭回源跟随
func (c *CdnClient) UpdateFollow302SwitchInvoker(request *model.UpdateFollow302SwitchRequest) *UpdateFollow302SwitchInvoker {
	requestDef := GenReqDefForUpdateFollow302Switch()
	return &UpdateFollow302SwitchInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateHttpsInfo 配置HTTPS
//
// 设置加速域名HTTPS。通过配置加速域名的HTTPS证书，并将其部署在全网CDN节点，实现HTTPS安全加速。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateHttpsInfo(request *model.UpdateHttpsInfoRequest) (*model.UpdateHttpsInfoResponse, error) {
	requestDef := GenReqDefForUpdateHttpsInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateHttpsInfoResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateHttpsInfoInvoker 配置HTTPS
func (c *CdnClient) UpdateHttpsInfoInvoker(request *model.UpdateHttpsInfoRequest) *UpdateHttpsInfoInvoker {
	requestDef := GenReqDefForUpdateHttpsInfo()
	return &UpdateHttpsInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateOriginHost 修改回源HOST。
//
// 修改回源HOST。回源HOST是CDN节点在回源过程中，在源站访问的站点域名，即http请求头中的host信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateOriginHost(request *model.UpdateOriginHostRequest) (*model.UpdateOriginHostResponse, error) {
	requestDef := GenReqDefForUpdateOriginHost()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateOriginHostResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateOriginHostInvoker 修改回源HOST。
func (c *CdnClient) UpdateOriginHostInvoker(request *model.UpdateOriginHostRequest) *UpdateOriginHostInvoker {
	requestDef := GenReqDefForUpdateOriginHost()
	return &UpdateOriginHostInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdatePrivateBucketAccess 修改私有桶开启关闭状态
//
// 修改私有桶开启关闭状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdatePrivateBucketAccess(request *model.UpdatePrivateBucketAccessRequest) (*model.UpdatePrivateBucketAccessResponse, error) {
	requestDef := GenReqDefForUpdatePrivateBucketAccess()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePrivateBucketAccessResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdatePrivateBucketAccessInvoker 修改私有桶开启关闭状态
func (c *CdnClient) UpdatePrivateBucketAccessInvoker(request *model.UpdatePrivateBucketAccessRequest) *UpdatePrivateBucketAccessInvoker {
	requestDef := GenReqDefForUpdatePrivateBucketAccess()
	return &UpdatePrivateBucketAccessInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateRangeSwitch 开启/关闭Range回源
//
// Range回源是指源站在收到CDN节点回源请求时，根据http请求头中的Range信息返回指定范围的数据给CDN节点。
//
// 开启Range回源前需要确认源站是否支持Range请求，若源站不支持Range请求，开启Range回源将导致资源无法缓存。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateRangeSwitch(request *model.UpdateRangeSwitchRequest) (*model.UpdateRangeSwitchResponse, error) {
	requestDef := GenReqDefForUpdateRangeSwitch()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateRangeSwitchResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateRangeSwitchInvoker 开启/关闭Range回源
func (c *CdnClient) UpdateRangeSwitchInvoker(request *model.UpdateRangeSwitchRequest) *UpdateRangeSwitchInvoker {
	requestDef := GenReqDefForUpdateRangeSwitch()
	return &UpdateRangeSwitchInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateRefer 设置Referer过滤规则
//
// 设置Referer过滤规则。通过设置过滤策略，对访问者身份进行识别和过滤，实现限制访问来源的目的。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateRefer(request *model.UpdateReferRequest) (*model.UpdateReferResponse, error) {
	requestDef := GenReqDefForUpdateRefer()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateReferResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateReferInvoker 设置Referer过滤规则
func (c *CdnClient) UpdateReferInvoker(request *model.UpdateReferRequest) *UpdateReferInvoker {
	requestDef := GenReqDefForUpdateRefer()
	return &UpdateReferInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateResponseHeader 新增/修改响应头配置
//
// 新增/修改域名响应头配置。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CdnClient) UpdateResponseHeader(request *model.UpdateResponseHeaderRequest) (*model.UpdateResponseHeaderResponse, error) {
	requestDef := GenReqDefForUpdateResponseHeader()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateResponseHeaderResponse), nil
	}
}

// Deprecated: This function is deprecated and will be removed in the future versions.
// UpdateResponseHeaderInvoker 新增/修改响应头配置
func (c *CdnClient) UpdateResponseHeaderInvoker(request *model.UpdateResponseHeaderRequest) *UpdateResponseHeaderInvoker {
	requestDef := GenReqDefForUpdateResponseHeader()
	return &UpdateResponseHeaderInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
