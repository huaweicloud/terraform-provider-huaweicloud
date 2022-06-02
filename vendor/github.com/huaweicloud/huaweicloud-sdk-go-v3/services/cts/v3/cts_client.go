package v3

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"
)

type CtsClient struct {
	HcClient *http_client.HcHttpClient
}

func NewCtsClient(hcClient *http_client.HcHttpClient) *CtsClient {
	return &CtsClient{HcClient: hcClient}
}

func CtsClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// CreateNotification 创建关键操作通知
//
// 配置关键操作通知，可在发生特定操作时，使用预先创建好的SMN主题，向用户手机、邮箱发送消息，也可直接发送http/https消息。常用于实时感知高危操作、触发特定操作或对接用户自有审计分析系统。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) CreateNotification(request *model.CreateNotificationRequest) (*model.CreateNotificationResponse, error) {
	requestDef := GenReqDefForCreateNotification()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateNotificationResponse), nil
	}
}

// CreateNotificationInvoker 创建关键操作通知
func (c *CtsClient) CreateNotificationInvoker(request *model.CreateNotificationRequest) *CreateNotificationInvoker {
	requestDef := GenReqDefForCreateNotification()
	return &CreateNotificationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTracker 创建追踪器
//
// 云审计服务开通后系统会自动创建一个追踪器，用来关联系统记录的所有操作。目前，一个云账户在一个Region下支持创建一个管理类追踪器和多个数据类追踪器。
// 云审计服务支持在管理控制台查询近7天内的操作记录。如需保存更长时间的操作记录，您可以在创建追踪器之后通过对象存储服务（Object Storage Service，以下简称OBS）将操作记录实时保存至OBS桶中。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) CreateTracker(request *model.CreateTrackerRequest) (*model.CreateTrackerResponse, error) {
	requestDef := GenReqDefForCreateTracker()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTrackerResponse), nil
	}
}

// CreateTrackerInvoker 创建追踪器
func (c *CtsClient) CreateTrackerInvoker(request *model.CreateTrackerRequest) *CreateTrackerInvoker {
	requestDef := GenReqDefForCreateTracker()
	return &CreateTrackerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteNotification 删除关键操作通知
//
// 云审计服务支持删除已创建的关键操作通知。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) DeleteNotification(request *model.DeleteNotificationRequest) (*model.DeleteNotificationResponse, error) {
	requestDef := GenReqDefForDeleteNotification()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteNotificationResponse), nil
	}
}

// DeleteNotificationInvoker 删除关键操作通知
func (c *CtsClient) DeleteNotificationInvoker(request *model.DeleteNotificationRequest) *DeleteNotificationInvoker {
	requestDef := GenReqDefForDeleteNotification()
	return &DeleteNotificationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTracker 删除追踪器
//
// 云审计服务目前仅支持删除已创建的数据类追踪器。删除追踪器对已有的操作记录没有影响，当您重新开通云审计服务后，依旧可以查看已有的操作记录。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) DeleteTracker(request *model.DeleteTrackerRequest) (*model.DeleteTrackerResponse, error) {
	requestDef := GenReqDefForDeleteTracker()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTrackerResponse), nil
	}
}

// DeleteTrackerInvoker 删除追踪器
func (c *CtsClient) DeleteTrackerInvoker(request *model.DeleteTrackerRequest) *DeleteTrackerInvoker {
	requestDef := GenReqDefForDeleteTracker()
	return &DeleteTrackerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListNotifications 查询关键操作通知
//
// 查询创建的关键操作通知规则。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) ListNotifications(request *model.ListNotificationsRequest) (*model.ListNotificationsResponse, error) {
	requestDef := GenReqDefForListNotifications()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListNotificationsResponse), nil
	}
}

// ListNotificationsInvoker 查询关键操作通知
func (c *CtsClient) ListNotificationsInvoker(request *model.ListNotificationsRequest) *ListNotificationsInvoker {
	requestDef := GenReqDefForListNotifications()
	return &ListNotificationsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListQuotas 查询租户追踪器配额信息
//
// 查询租户追踪器配额信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) ListQuotas(request *model.ListQuotasRequest) (*model.ListQuotasResponse, error) {
	requestDef := GenReqDefForListQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListQuotasResponse), nil
	}
}

// ListQuotasInvoker 查询租户追踪器配额信息
func (c *CtsClient) ListQuotasInvoker(request *model.ListQuotasRequest) *ListQuotasInvoker {
	requestDef := GenReqDefForListQuotas()
	return &ListQuotasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTraces 查询事件列表
//
// 通过事件列表查询接口，可以查出系统记录的7天内资源操作记录。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) ListTraces(request *model.ListTracesRequest) (*model.ListTracesResponse, error) {
	requestDef := GenReqDefForListTraces()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTracesResponse), nil
	}
}

// ListTracesInvoker 查询事件列表
func (c *CtsClient) ListTracesInvoker(request *model.ListTracesRequest) *ListTracesInvoker {
	requestDef := GenReqDefForListTraces()
	return &ListTracesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTrackers 查询追踪器
//
// 开通云审计服务成功后，您可以在追踪器信息页面查看追踪器的详细信息。详细信息主要包括追踪器名称，用于存储操作事件的OBS桶名称和OBS桶中的事件文件前缀。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) ListTrackers(request *model.ListTrackersRequest) (*model.ListTrackersResponse, error) {
	requestDef := GenReqDefForListTrackers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTrackersResponse), nil
	}
}

// ListTrackersInvoker 查询追踪器
func (c *CtsClient) ListTrackersInvoker(request *model.ListTrackersRequest) *ListTrackersInvoker {
	requestDef := GenReqDefForListTrackers()
	return &ListTrackersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateNotification 修改关键操作通知
//
// 云审计服务支持修改已创建关键操作通知配置项，通过notification_id的字段匹配修改对象，notification_id必须已经存在。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) UpdateNotification(request *model.UpdateNotificationRequest) (*model.UpdateNotificationResponse, error) {
	requestDef := GenReqDefForUpdateNotification()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateNotificationResponse), nil
	}
}

// UpdateNotificationInvoker 修改关键操作通知
func (c *CtsClient) UpdateNotificationInvoker(request *model.UpdateNotificationRequest) *UpdateNotificationInvoker {
	requestDef := GenReqDefForUpdateNotification()
	return &UpdateNotificationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTracker 修改追踪器
//
// 云审计服务支持修改已创建追踪器的配置项，包括OBS桶转储、关键事件通知、事件转储加密、通过LTS对管理类事件进行检索、事件文件完整性校验以及追踪器启停状态等相关参数，修改追踪器对已有的操作记录没有影响。修改追踪器完成后，系统立即以新的规则开始记录操作。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CtsClient) UpdateTracker(request *model.UpdateTrackerRequest) (*model.UpdateTrackerResponse, error) {
	requestDef := GenReqDefForUpdateTracker()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTrackerResponse), nil
	}
}

// UpdateTrackerInvoker 修改追踪器
func (c *CtsClient) UpdateTrackerInvoker(request *model.UpdateTrackerRequest) *UpdateTrackerInvoker {
	requestDef := GenReqDefForUpdateTracker()
	return &UpdateTrackerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
