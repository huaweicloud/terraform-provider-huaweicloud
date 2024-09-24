package v5

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
)

type IoTDAClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewIoTDAClient(hcClient *httpclient.HcHttpClient) *IoTDAClient {
	return &IoTDAClient{HcClient: hcClient}
}

func IoTDAClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder().WithDerivedAuthServiceName("iotdm")
	return builder
}

// CreateAccessCode 生成接入凭证
//
// 接入凭证是用于客户端使用AMQP等协议与平台建链的一个认证凭据。只保留一条记录，如果重复调用只会重置接入凭证，使得之前的失效。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateAccessCode(request *model.CreateAccessCodeRequest) (*model.CreateAccessCodeResponse, error) {
	requestDef := GenReqDefForCreateAccessCode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAccessCodeResponse), nil
	}
}

// CreateAccessCodeInvoker 生成接入凭证
func (c *IoTDAClient) CreateAccessCodeInvoker(request *model.CreateAccessCodeRequest) *CreateAccessCodeInvoker {
	requestDef := GenReqDefForCreateAccessCode()
	return &CreateAccessCodeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddQueue 创建AMQP队列
//
// 应用服务器可调用此接口在物联网平台创建一个AMQP队列。每个租户只能创建100个队列，若超过规格，则创建失败，若队列名称与已有的队列名称相同，则创建失败。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) AddQueue(request *model.AddQueueRequest) (*model.AddQueueResponse, error) {
	requestDef := GenReqDefForAddQueue()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddQueueResponse), nil
	}
}

// AddQueueInvoker 创建AMQP队列
func (c *IoTDAClient) AddQueueInvoker(request *model.AddQueueRequest) *AddQueueInvoker {
	requestDef := GenReqDefForAddQueue()
	return &AddQueueInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchShowQueue 查询AMQP列表
//
// 应用服务器可调用此接口查询物联网平台中的AMQP队列信息列表。可通过队列名称作模糊查询，支持分页。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) BatchShowQueue(request *model.BatchShowQueueRequest) (*model.BatchShowQueueResponse, error) {
	requestDef := GenReqDefForBatchShowQueue()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchShowQueueResponse), nil
	}
}

// BatchShowQueueInvoker 查询AMQP列表
func (c *IoTDAClient) BatchShowQueueInvoker(request *model.BatchShowQueueRequest) *BatchShowQueueInvoker {
	requestDef := GenReqDefForBatchShowQueue()
	return &BatchShowQueueInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteQueue 删除AMQP队列
//
// 应用服务器可调用此接口在物联网平台上删除指定AMQP队列。若当前队列正在使用，则会删除失败。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteQueue(request *model.DeleteQueueRequest) (*model.DeleteQueueResponse, error) {
	requestDef := GenReqDefForDeleteQueue()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteQueueResponse), nil
	}
}

// DeleteQueueInvoker 删除AMQP队列
func (c *IoTDAClient) DeleteQueueInvoker(request *model.DeleteQueueRequest) *DeleteQueueInvoker {
	requestDef := GenReqDefForDeleteQueue()
	return &DeleteQueueInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowQueue 查询单个AMQP队列
//
// 应用服务器可调用此接口查询物联网平台中指定队列的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowQueue(request *model.ShowQueueRequest) (*model.ShowQueueResponse, error) {
	requestDef := GenReqDefForShowQueue()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowQueueResponse), nil
	}
}

// ShowQueueInvoker 查询单个AMQP队列
func (c *IoTDAClient) ShowQueueInvoker(request *model.ShowQueueRequest) *ShowQueueInvoker {
	requestDef := GenReqDefForShowQueue()
	return &ShowQueueInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddApplication 创建资源空间
//
// 资源空间对应的是物联网平台原有的应用，在物联网平台的含义与应用一致，只是变更了名称。应用服务器可以调用此接口创建资源空间。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) AddApplication(request *model.AddApplicationRequest) (*model.AddApplicationResponse, error) {
	requestDef := GenReqDefForAddApplication()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddApplicationResponse), nil
	}
}

// AddApplicationInvoker 创建资源空间
func (c *IoTDAClient) AddApplicationInvoker(request *model.AddApplicationRequest) *AddApplicationInvoker {
	requestDef := GenReqDefForAddApplication()
	return &AddApplicationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteApplication 删除资源空间
//
// 删除指定资源空间。删除资源空间属于高危操作，删除资源空间后，该空间下的产品、设备等资源将不可用，请谨慎操作！
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteApplication(request *model.DeleteApplicationRequest) (*model.DeleteApplicationResponse, error) {
	requestDef := GenReqDefForDeleteApplication()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteApplicationResponse), nil
	}
}

// DeleteApplicationInvoker 删除资源空间
func (c *IoTDAClient) DeleteApplicationInvoker(request *model.DeleteApplicationRequest) *DeleteApplicationInvoker {
	requestDef := GenReqDefForDeleteApplication()
	return &DeleteApplicationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowApplication 查询资源空间
//
// 资源空间对应的是物联网平台原有的应用，在物联网平台的含义与应用一致，只是变更了名称。应用服务器可以调用此接口查询指定资源空间详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowApplication(request *model.ShowApplicationRequest) (*model.ShowApplicationResponse, error) {
	requestDef := GenReqDefForShowApplication()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowApplicationResponse), nil
	}
}

// ShowApplicationInvoker 查询资源空间
func (c *IoTDAClient) ShowApplicationInvoker(request *model.ShowApplicationRequest) *ShowApplicationInvoker {
	requestDef := GenReqDefForShowApplication()
	return &ShowApplicationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowApplications 查询资源空间列表
//
// 资源空间对应的是物联网平台原有的应用，在物联网平台的含义与应用一致，只是变更了名称。应用服务器可以调用此接口查询资源空间列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowApplications(request *model.ShowApplicationsRequest) (*model.ShowApplicationsResponse, error) {
	requestDef := GenReqDefForShowApplications()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowApplicationsResponse), nil
	}
}

// ShowApplicationsInvoker 查询资源空间列表
func (c *IoTDAClient) ShowApplicationsInvoker(request *model.ShowApplicationsRequest) *ShowApplicationsInvoker {
	requestDef := GenReqDefForShowApplications()
	return &ShowApplicationsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateApplication 更新资源空间
//
// 应用服务器可以调用此接口更新资源空间的名称
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateApplication(request *model.UpdateApplicationRequest) (*model.UpdateApplicationResponse, error) {
	requestDef := GenReqDefForUpdateApplication()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateApplicationResponse), nil
	}
}

// UpdateApplicationInvoker 更新资源空间
func (c *IoTDAClient) UpdateApplicationInvoker(request *model.UpdateApplicationRequest) *UpdateApplicationInvoker {
	requestDef := GenReqDefForUpdateApplication()
	return &UpdateApplicationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAsyncCommand 下发异步设备命令
//
// 设备的产品模型中定义了物联网平台可向设备下发的命令，应用服务器可调用此接口向指定设备下发异步命令，以实现对设备的控制。平台负责将命令发送给设备，并将设备执行命令结果异步通知应用服务器。 命令执行结果支持灵活的数据流转，应用服务器通过调用物联网平台的创建规则触发条件（Resource:device.command.status，Event:update）、创建规则动作并激活规则后，当命令状态变更时，物联网平台会根据规则将结果发送到规则指定的服务器，如用户自定义的HTTP服务器，AMQP服务器，以及华为云的其他储存服务器等, 详情参考[[设备命令状态变更通知](https://support.huaweicloud.com/api-iothub/iot_06_v5_01212.html)](tag:hws)[[设备命令状态变更通知](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_01212.html)](tag:hws_hk)。
// 注意：
// - 此接口适用于NB设备异步命令下发，暂不支持其他协议类型设备命令下发。
// - 此接口仅支持单个设备异步命令下发，如需多个设备异步命令下发，请参见 [[创建批量任务](https://support.huaweicloud.com/api-iothub/iot_06_v5_0045.html)](tag:hws)[[创建批量任务](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0045.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateAsyncCommand(request *model.CreateAsyncCommandRequest) (*model.CreateAsyncCommandResponse, error) {
	requestDef := GenReqDefForCreateAsyncCommand()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAsyncCommandResponse), nil
	}
}

// CreateAsyncCommandInvoker 下发异步设备命令
func (c *IoTDAClient) CreateAsyncCommandInvoker(request *model.CreateAsyncCommandRequest) *CreateAsyncCommandInvoker {
	requestDef := GenReqDefForCreateAsyncCommand()
	return &CreateAsyncCommandInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAsyncDeviceCommand 查询指定id的命令
//
// 物联网平台可查询指定id的命令。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowAsyncDeviceCommand(request *model.ShowAsyncDeviceCommandRequest) (*model.ShowAsyncDeviceCommandResponse, error) {
	requestDef := GenReqDefForShowAsyncDeviceCommand()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAsyncDeviceCommandResponse), nil
	}
}

// ShowAsyncDeviceCommandInvoker 查询指定id的命令
func (c *IoTDAClient) ShowAsyncDeviceCommandInvoker(request *model.ShowAsyncDeviceCommandRequest) *ShowAsyncDeviceCommandInvoker {
	requestDef := GenReqDefForShowAsyncDeviceCommand()
	return &ShowAsyncDeviceCommandInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRoutingBacklogPolicy 新建数据流转积压策略
//
// 应用服务器可调用此接口在物联网平台创建数据流转积压策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateRoutingBacklogPolicy(request *model.CreateRoutingBacklogPolicyRequest) (*model.CreateRoutingBacklogPolicyResponse, error) {
	requestDef := GenReqDefForCreateRoutingBacklogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRoutingBacklogPolicyResponse), nil
	}
}

// CreateRoutingBacklogPolicyInvoker 新建数据流转积压策略
func (c *IoTDAClient) CreateRoutingBacklogPolicyInvoker(request *model.CreateRoutingBacklogPolicyRequest) *CreateRoutingBacklogPolicyInvoker {
	requestDef := GenReqDefForCreateRoutingBacklogPolicy()
	return &CreateRoutingBacklogPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRoutingBacklogPolicy 删除数据流转积压策略
//
// 应用服务器可调用此接口在物联网平台删除指定数据流转积压策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteRoutingBacklogPolicy(request *model.DeleteRoutingBacklogPolicyRequest) (*model.DeleteRoutingBacklogPolicyResponse, error) {
	requestDef := GenReqDefForDeleteRoutingBacklogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteRoutingBacklogPolicyResponse), nil
	}
}

// DeleteRoutingBacklogPolicyInvoker 删除数据流转积压策略
func (c *IoTDAClient) DeleteRoutingBacklogPolicyInvoker(request *model.DeleteRoutingBacklogPolicyRequest) *DeleteRoutingBacklogPolicyInvoker {
	requestDef := GenReqDefForDeleteRoutingBacklogPolicy()
	return &DeleteRoutingBacklogPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRoutingBacklogPolicy 查询数据流转积压策略列表
//
// 应用服务器可调用此接口查询在物联网平台设置的数据流转积压策略列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListRoutingBacklogPolicy(request *model.ListRoutingBacklogPolicyRequest) (*model.ListRoutingBacklogPolicyResponse, error) {
	requestDef := GenReqDefForListRoutingBacklogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRoutingBacklogPolicyResponse), nil
	}
}

// ListRoutingBacklogPolicyInvoker 查询数据流转积压策略列表
func (c *IoTDAClient) ListRoutingBacklogPolicyInvoker(request *model.ListRoutingBacklogPolicyRequest) *ListRoutingBacklogPolicyInvoker {
	requestDef := GenReqDefForListRoutingBacklogPolicy()
	return &ListRoutingBacklogPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRoutingBacklogPolicy 查询数据流转积压策略
//
// 应用服务器可调用此接口在物联网平台查询指定数据流转积压策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowRoutingBacklogPolicy(request *model.ShowRoutingBacklogPolicyRequest) (*model.ShowRoutingBacklogPolicyResponse, error) {
	requestDef := GenReqDefForShowRoutingBacklogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRoutingBacklogPolicyResponse), nil
	}
}

// ShowRoutingBacklogPolicyInvoker 查询数据流转积压策略
func (c *IoTDAClient) ShowRoutingBacklogPolicyInvoker(request *model.ShowRoutingBacklogPolicyRequest) *ShowRoutingBacklogPolicyInvoker {
	requestDef := GenReqDefForShowRoutingBacklogPolicy()
	return &ShowRoutingBacklogPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateRoutingBacklogPolicy 修改数据流转积压策略
//
// 应用服务器可调用此接口在物联网平台修改指定数据流转积压策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateRoutingBacklogPolicy(request *model.UpdateRoutingBacklogPolicyRequest) (*model.UpdateRoutingBacklogPolicyResponse, error) {
	requestDef := GenReqDefForUpdateRoutingBacklogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateRoutingBacklogPolicyResponse), nil
	}
}

// UpdateRoutingBacklogPolicyInvoker 修改数据流转积压策略
func (c *IoTDAClient) UpdateRoutingBacklogPolicyInvoker(request *model.UpdateRoutingBacklogPolicyRequest) *UpdateRoutingBacklogPolicyInvoker {
	requestDef := GenReqDefForUpdateRoutingBacklogPolicy()
	return &UpdateRoutingBacklogPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateBatchTask 创建批量任务
//
// 应用服务器可调用此接口为创建批量处理任务，对多个设备进行批量操作。当前支持批量软固件升级、批量创建设备、批量删除设备、批量冻结设备、批量解冻设备、批量创建命令、批量创建消息任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateBatchTask(request *model.CreateBatchTaskRequest) (*model.CreateBatchTaskResponse, error) {
	requestDef := GenReqDefForCreateBatchTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateBatchTaskResponse), nil
	}
}

// CreateBatchTaskInvoker 创建批量任务
func (c *IoTDAClient) CreateBatchTaskInvoker(request *model.CreateBatchTaskRequest) *CreateBatchTaskInvoker {
	requestDef := GenReqDefForCreateBatchTask()
	return &CreateBatchTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteBatchTask 删除批量任务
//
// 应用服务器可调用此接口删除物联网平台中已经完成（状态为成功，失败，部分成功，已停止）的批量任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteBatchTask(request *model.DeleteBatchTaskRequest) (*model.DeleteBatchTaskResponse, error) {
	requestDef := GenReqDefForDeleteBatchTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteBatchTaskResponse), nil
	}
}

// DeleteBatchTaskInvoker 删除批量任务
func (c *IoTDAClient) DeleteBatchTaskInvoker(request *model.DeleteBatchTaskRequest) *DeleteBatchTaskInvoker {
	requestDef := GenReqDefForDeleteBatchTask()
	return &DeleteBatchTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListBatchTasks 查询批量任务列表
//
// 应用服务器可调用此接口查询物联网平台中批量任务列表，每一个任务又包括具体的任务内容、任务状态、任务完成情况统计等。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListBatchTasks(request *model.ListBatchTasksRequest) (*model.ListBatchTasksResponse, error) {
	requestDef := GenReqDefForListBatchTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListBatchTasksResponse), nil
	}
}

// ListBatchTasksInvoker 查询批量任务列表
func (c *IoTDAClient) ListBatchTasksInvoker(request *model.ListBatchTasksRequest) *ListBatchTasksInvoker {
	requestDef := GenReqDefForListBatchTasks()
	return &ListBatchTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RetryBatchTask 重试批量任务
//
// 应用服务器可调用此接口重试批量任务，目前只支持task_type为firmwareUpgrade，softwareUpgrade。如果task_id对应任务已经成功、停止、正在停止、等待中或初始化中，则不可以调用该接口。如果请求Body为{}，则调用该接口后会重新执行所有状态为失败、失败待重试和已停止的子任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) RetryBatchTask(request *model.RetryBatchTaskRequest) (*model.RetryBatchTaskResponse, error) {
	requestDef := GenReqDefForRetryBatchTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RetryBatchTaskResponse), nil
	}
}

// RetryBatchTaskInvoker 重试批量任务
func (c *IoTDAClient) RetryBatchTaskInvoker(request *model.RetryBatchTaskRequest) *RetryBatchTaskInvoker {
	requestDef := GenReqDefForRetryBatchTask()
	return &RetryBatchTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowBatchTask 查询批量任务
//
// 应用服务器可调用此接口查询物联网平台中指定批量任务的信息，包括任务内容、任务状态、任务完成情况统计以及子任务列表等。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowBatchTask(request *model.ShowBatchTaskRequest) (*model.ShowBatchTaskResponse, error) {
	requestDef := GenReqDefForShowBatchTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBatchTaskResponse), nil
	}
}

// ShowBatchTaskInvoker 查询批量任务
func (c *IoTDAClient) ShowBatchTaskInvoker(request *model.ShowBatchTaskRequest) *ShowBatchTaskInvoker {
	requestDef := GenReqDefForShowBatchTask()
	return &ShowBatchTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopBatchTask 停止批量任务
//
// 应用服务器可调用此接口停止批量任务，目前只支持task_type为firmwareUpgrade，softwareUpgrade。如果task_id对应任务已经完成（成功、失败、部分成功，已经停止）或正在停止中，则不可以调用该接口。如果请求Body为{}，则调用该接口后会停止所有正在执行中、等待中和失败待重试状态的子任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) StopBatchTask(request *model.StopBatchTaskRequest) (*model.StopBatchTaskResponse, error) {
	requestDef := GenReqDefForStopBatchTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopBatchTaskResponse), nil
	}
}

// StopBatchTaskInvoker 停止批量任务
func (c *IoTDAClient) StopBatchTaskInvoker(request *model.StopBatchTaskRequest) *StopBatchTaskInvoker {
	requestDef := GenReqDefForStopBatchTask()
	return &StopBatchTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteBatchTaskFile 删除批量任务文件
//
// 应用服务器可调用此接口删除批量任务文件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteBatchTaskFile(request *model.DeleteBatchTaskFileRequest) (*model.DeleteBatchTaskFileResponse, error) {
	requestDef := GenReqDefForDeleteBatchTaskFile()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteBatchTaskFileResponse), nil
	}
}

// DeleteBatchTaskFileInvoker 删除批量任务文件
func (c *IoTDAClient) DeleteBatchTaskFileInvoker(request *model.DeleteBatchTaskFileRequest) *DeleteBatchTaskFileInvoker {
	requestDef := GenReqDefForDeleteBatchTaskFile()
	return &DeleteBatchTaskFileInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListBatchTaskFiles 查询批量任务文件列表
//
// 应用服务器可调用此接口查询批量任务文件列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListBatchTaskFiles(request *model.ListBatchTaskFilesRequest) (*model.ListBatchTaskFilesResponse, error) {
	requestDef := GenReqDefForListBatchTaskFiles()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListBatchTaskFilesResponse), nil
	}
}

// ListBatchTaskFilesInvoker 查询批量任务文件列表
func (c *IoTDAClient) ListBatchTaskFilesInvoker(request *model.ListBatchTaskFilesRequest) *ListBatchTaskFilesInvoker {
	requestDef := GenReqDefForListBatchTaskFiles()
	return &ListBatchTaskFilesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UploadBatchTaskFile 上传批量任务文件
//
// 应用服务器可调用此接口上传批量任务文件，用于创建批量任务。当前支持批量创建设备任务、批量删除设备任务、批量冻结设备任务、批量解冻设备任务的文件上传。
// - [批量注册设备模板](https://developer.obs.cn-north-4.myhuaweicloud.com/template/BatchCreateDevices_Template.xlsx)
//
// - [批量删除设备模板](https://developer.obs.cn-north-4.myhuaweicloud.com/template/BatchDeleteDevices_Template.xlsx)
//
// - [批量冻结设备模板](https://developer.obs.cn-north-4.myhuaweicloud.com/template/BatchFreezeDevices_Template.xlsx)
//
// - [批量解冻设备模板](https://developer.obs.cn-north-4.myhuaweicloud.com/template/BatchUnfreezeDevices_Template.xlsx)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UploadBatchTaskFile(request *model.UploadBatchTaskFileRequest) (*model.UploadBatchTaskFileResponse, error) {
	requestDef := GenReqDefForUploadBatchTaskFile()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UploadBatchTaskFileResponse), nil
	}
}

// UploadBatchTaskFileInvoker 上传批量任务文件
func (c *IoTDAClient) UploadBatchTaskFileInvoker(request *model.UploadBatchTaskFileRequest) *UploadBatchTaskFileInvoker {
	requestDef := GenReqDefForUploadBatchTaskFile()
	return &UploadBatchTaskFileInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddBridge 创建网桥
//
// 应用服务器可调用此接口在物联网平台创建一个网桥，仅在创建后的网桥才可以接入物联网平台。
// - 一个实例最多支持20个网桥。
// - 仅**标准版实例、企业版实例**支持该接口调用，基础版不支持。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) AddBridge(request *model.AddBridgeRequest) (*model.AddBridgeResponse, error) {
	requestDef := GenReqDefForAddBridge()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddBridgeResponse), nil
	}
}

// AddBridgeInvoker 创建网桥
func (c *IoTDAClient) AddBridgeInvoker(request *model.AddBridgeRequest) *AddBridgeInvoker {
	requestDef := GenReqDefForAddBridge()
	return &AddBridgeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteBridge 删除网桥
//
// 应用服务器可调用此接口在物联网平台上删除指定网桥。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteBridge(request *model.DeleteBridgeRequest) (*model.DeleteBridgeResponse, error) {
	requestDef := GenReqDefForDeleteBridge()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteBridgeResponse), nil
	}
}

// DeleteBridgeInvoker 删除网桥
func (c *IoTDAClient) DeleteBridgeInvoker(request *model.DeleteBridgeRequest) *DeleteBridgeInvoker {
	requestDef := GenReqDefForDeleteBridge()
	return &DeleteBridgeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListBridges 查询网桥列表
//
// 应用服务器可调用此接口在物联网平台查询网桥列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListBridges(request *model.ListBridgesRequest) (*model.ListBridgesResponse, error) {
	requestDef := GenReqDefForListBridges()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListBridgesResponse), nil
	}
}

// ListBridgesInvoker 查询网桥列表
func (c *IoTDAClient) ListBridgesInvoker(request *model.ListBridgesRequest) *ListBridgesInvoker {
	requestDef := GenReqDefForListBridges()
	return &ListBridgesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ResetBridgeSecret 重置网桥密钥
//
// 应用服务器可调用此接口在物联网平台上重置网桥密钥。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ResetBridgeSecret(request *model.ResetBridgeSecretRequest) (*model.ResetBridgeSecretResponse, error) {
	requestDef := GenReqDefForResetBridgeSecret()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetBridgeSecretResponse), nil
	}
}

// ResetBridgeSecretInvoker 重置网桥密钥
func (c *IoTDAClient) ResetBridgeSecretInvoker(request *model.ResetBridgeSecretRequest) *ResetBridgeSecretInvoker {
	requestDef := GenReqDefForResetBridgeSecret()
	return &ResetBridgeSecretInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BroadcastMessage 下发广播消息
//
// 应用服务器可调用此接口向订阅了指定Topic的所有在线设备发布广播消息。应用将广播消息下发给平台后，平台会先返回应用响应结果，再将消息广播给设备。
// 注意：
// - 此接口只适用于使用MQTT协议接入的设备。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) BroadcastMessage(request *model.BroadcastMessageRequest) (*model.BroadcastMessageResponse, error) {
	requestDef := GenReqDefForBroadcastMessage()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BroadcastMessageResponse), nil
	}
}

// BroadcastMessageInvoker 下发广播消息
func (c *IoTDAClient) BroadcastMessageInvoker(request *model.BroadcastMessageRequest) *BroadcastMessageInvoker {
	requestDef := GenReqDefForBroadcastMessage()
	return &BroadcastMessageInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddCertificate 上传设备CA证书
//
// 应用服务器可调用此接口在物联网平台上传设备CA证书
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) AddCertificate(request *model.AddCertificateRequest) (*model.AddCertificateResponse, error) {
	requestDef := GenReqDefForAddCertificate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddCertificateResponse), nil
	}
}

// AddCertificateInvoker 上传设备CA证书
func (c *IoTDAClient) AddCertificateInvoker(request *model.AddCertificateRequest) *AddCertificateInvoker {
	requestDef := GenReqDefForAddCertificate()
	return &AddCertificateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CheckCertificate 验证设备CA证书
//
// 应用服务器可调用此接口在物联网平台验证设备的CA证书，目的是为了验证用户持有设备CA证书的私钥
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CheckCertificate(request *model.CheckCertificateRequest) (*model.CheckCertificateResponse, error) {
	requestDef := GenReqDefForCheckCertificate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CheckCertificateResponse), nil
	}
}

// CheckCertificateInvoker 验证设备CA证书
func (c *IoTDAClient) CheckCertificateInvoker(request *model.CheckCertificateRequest) *CheckCertificateInvoker {
	requestDef := GenReqDefForCheckCertificate()
	return &CheckCertificateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteCertificate 删除设备CA证书
//
// 应用服务器可调用此接口在物联网平台删除设备CA证书
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteCertificate(request *model.DeleteCertificateRequest) (*model.DeleteCertificateResponse, error) {
	requestDef := GenReqDefForDeleteCertificate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteCertificateResponse), nil
	}
}

// DeleteCertificateInvoker 删除设备CA证书
func (c *IoTDAClient) DeleteCertificateInvoker(request *model.DeleteCertificateRequest) *DeleteCertificateInvoker {
	requestDef := GenReqDefForDeleteCertificate()
	return &DeleteCertificateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListCertificates 获取设备CA证书列表
//
// 应用服务器可调用此接口在物联网平台获取设备CA证书列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListCertificates(request *model.ListCertificatesRequest) (*model.ListCertificatesResponse, error) {
	requestDef := GenReqDefForListCertificates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListCertificatesResponse), nil
	}
}

// ListCertificatesInvoker 获取设备CA证书列表
func (c *IoTDAClient) ListCertificatesInvoker(request *model.ListCertificatesRequest) *ListCertificatesInvoker {
	requestDef := GenReqDefForListCertificates()
	return &ListCertificatesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateCertificate 更新CA证书
//
// 应用服务器可调用此接口在物联网平台上更新CA证书。仅标准版实例、企业版实例支持该接口调用，基础版不支持。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateCertificate(request *model.UpdateCertificateRequest) (*model.UpdateCertificateResponse, error) {
	requestDef := GenReqDefForUpdateCertificate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCertificateResponse), nil
	}
}

// UpdateCertificateInvoker 更新CA证书
func (c *IoTDAClient) UpdateCertificateInvoker(request *model.UpdateCertificateRequest) *UpdateCertificateInvoker {
	requestDef := GenReqDefForUpdateCertificate()
	return &UpdateCertificateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateCommand 下发设备命令
//
// 设备的产品模型中定义了物联网平台可向设备下发的命令，应用服务器可调用此接口向指定设备下发命令，以实现对设备的同步控制。平台负责将命令以同步方式发送给设备，并将设备执行命令结果同步返回, 如果设备没有响应，平台会返回给应用服务器超时，平台超时时间是20秒。如果命令下发需要超过20秒，建议采用[[消息下发](https://support.huaweicloud.com/api-iothub/iot_06_v5_0059.html)](tag:hws)[[消息下发](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0059.html)](tag:hws_hk)。
// 注意：
// - 此接口适用于MQTT设备同步命令下发，暂不支持NB-IoT设备命令下发。
// - 此接口仅支持单个设备同步命令下发，如需多个设备同步命令下发，请参见 [[创建批量任务](https://support.huaweicloud.com/api-iothub/iot_06_v5_0045.html)](tag:hws)[[创建批量任务](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0045.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateCommand(request *model.CreateCommandRequest) (*model.CreateCommandResponse, error) {
	requestDef := GenReqDefForCreateCommand()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateCommandResponse), nil
	}
}

// CreateCommandInvoker 下发设备命令
func (c *IoTDAClient) CreateCommandInvoker(request *model.CreateCommandRequest) *CreateCommandInvoker {
	requestDef := GenReqDefForCreateCommand()
	return &CreateCommandInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDeviceAuthorizer 创建自定义鉴权
//
// 应用服务器可调用此接口在物联网平台创建一个自定义鉴权。自定义鉴权是指用户可以通过函数服务自定义实现鉴权逻辑，以对接入平台的设备进行身份认证。
// - 单个实例最大可配置10个自定义鉴权
// - 仅标准版实例、企业版实例支持该接口调用，基础版不支持。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateDeviceAuthorizer(request *model.CreateDeviceAuthorizerRequest) (*model.CreateDeviceAuthorizerResponse, error) {
	requestDef := GenReqDefForCreateDeviceAuthorizer()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDeviceAuthorizerResponse), nil
	}
}

// CreateDeviceAuthorizerInvoker 创建自定义鉴权
func (c *IoTDAClient) CreateDeviceAuthorizerInvoker(request *model.CreateDeviceAuthorizerRequest) *CreateDeviceAuthorizerInvoker {
	requestDef := GenReqDefForCreateDeviceAuthorizer()
	return &CreateDeviceAuthorizerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDeviceAuthorizer 删除自定义鉴权
//
// 应用服务器可调用此接口在物联网平台上删除指定自定义鉴权。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteDeviceAuthorizer(request *model.DeleteDeviceAuthorizerRequest) (*model.DeleteDeviceAuthorizerResponse, error) {
	requestDef := GenReqDefForDeleteDeviceAuthorizer()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDeviceAuthorizerResponse), nil
	}
}

// DeleteDeviceAuthorizerInvoker 删除自定义鉴权
func (c *IoTDAClient) DeleteDeviceAuthorizerInvoker(request *model.DeleteDeviceAuthorizerRequest) *DeleteDeviceAuthorizerInvoker {
	requestDef := GenReqDefForDeleteDeviceAuthorizer()
	return &DeleteDeviceAuthorizerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDeviceAuthorizers 查询自定义鉴权列表
//
// 应用服务器可调用此接口在物联网平台查询自定义鉴权列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListDeviceAuthorizers(request *model.ListDeviceAuthorizersRequest) (*model.ListDeviceAuthorizersResponse, error) {
	requestDef := GenReqDefForListDeviceAuthorizers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDeviceAuthorizersResponse), nil
	}
}

// ListDeviceAuthorizersInvoker 查询自定义鉴权列表
func (c *IoTDAClient) ListDeviceAuthorizersInvoker(request *model.ListDeviceAuthorizersRequest) *ListDeviceAuthorizersInvoker {
	requestDef := GenReqDefForListDeviceAuthorizers()
	return &ListDeviceAuthorizersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDeviceAuthorizer 查询自定义鉴权详情
//
// 应用服务器可调用此接口在物联网平台查询指定自定义鉴权ID的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDeviceAuthorizer(request *model.ShowDeviceAuthorizerRequest) (*model.ShowDeviceAuthorizerResponse, error) {
	requestDef := GenReqDefForShowDeviceAuthorizer()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDeviceAuthorizerResponse), nil
	}
}

// ShowDeviceAuthorizerInvoker 查询自定义鉴权详情
func (c *IoTDAClient) ShowDeviceAuthorizerInvoker(request *model.ShowDeviceAuthorizerRequest) *ShowDeviceAuthorizerInvoker {
	requestDef := GenReqDefForShowDeviceAuthorizer()
	return &ShowDeviceAuthorizerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDeviceAuthorizer 更新指定id的自定义鉴权
//
// 应用服务器可调用此接口在物联网平台更新指定id的自定义鉴权。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateDeviceAuthorizer(request *model.UpdateDeviceAuthorizerRequest) (*model.UpdateDeviceAuthorizerResponse, error) {
	requestDef := GenReqDefForUpdateDeviceAuthorizer()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDeviceAuthorizerResponse), nil
	}
}

// UpdateDeviceAuthorizerInvoker 更新指定id的自定义鉴权
func (c *IoTDAClient) UpdateDeviceAuthorizerInvoker(request *model.UpdateDeviceAuthorizerRequest) *UpdateDeviceAuthorizerInvoker {
	requestDef := GenReqDefForUpdateDeviceAuthorizer()
	return &UpdateDeviceAuthorizerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddDeviceGroup 添加设备组
//
// 应用服务器可调用此接口新建设备组，一个华为云账号下最多可有1,000个设备组，包括父设备组和子设备组。设备组的最大层级关系不超过5层，即群组形成的关系树最大深度不超过5。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) AddDeviceGroup(request *model.AddDeviceGroupRequest) (*model.AddDeviceGroupResponse, error) {
	requestDef := GenReqDefForAddDeviceGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddDeviceGroupResponse), nil
	}
}

// AddDeviceGroupInvoker 添加设备组
func (c *IoTDAClient) AddDeviceGroupInvoker(request *model.AddDeviceGroupRequest) *AddDeviceGroupInvoker {
	requestDef := GenReqDefForAddDeviceGroup()
	return &AddDeviceGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateOrDeleteDeviceInGroup 管理设备组中的设备
//
// 应用服务器可调用此接口管理设备组中的设备。单个设备组内最多添加20,000个设备，一个设备最多可以被添加到10个设备组中。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateOrDeleteDeviceInGroup(request *model.CreateOrDeleteDeviceInGroupRequest) (*model.CreateOrDeleteDeviceInGroupResponse, error) {
	requestDef := GenReqDefForCreateOrDeleteDeviceInGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateOrDeleteDeviceInGroupResponse), nil
	}
}

// CreateOrDeleteDeviceInGroupInvoker 管理设备组中的设备
func (c *IoTDAClient) CreateOrDeleteDeviceInGroupInvoker(request *model.CreateOrDeleteDeviceInGroupRequest) *CreateOrDeleteDeviceInGroupInvoker {
	requestDef := GenReqDefForCreateOrDeleteDeviceInGroup()
	return &CreateOrDeleteDeviceInGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDeviceGroup 删除设备组
//
// 应用服务器可调用此接口删除指定设备组，如果该设备组存在子设备组或者该设备组中存在设备，必须先删除子设备组并将设备从该设备组移除，才能删除该设备组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteDeviceGroup(request *model.DeleteDeviceGroupRequest) (*model.DeleteDeviceGroupResponse, error) {
	requestDef := GenReqDefForDeleteDeviceGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDeviceGroupResponse), nil
	}
}

// DeleteDeviceGroupInvoker 删除设备组
func (c *IoTDAClient) DeleteDeviceGroupInvoker(request *model.DeleteDeviceGroupRequest) *DeleteDeviceGroupInvoker {
	requestDef := GenReqDefForDeleteDeviceGroup()
	return &DeleteDeviceGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDeviceGroups 查询设备组列表
//
// 应用服务器可调用此接口查询物联网平台中的设备组信息列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListDeviceGroups(request *model.ListDeviceGroupsRequest) (*model.ListDeviceGroupsResponse, error) {
	requestDef := GenReqDefForListDeviceGroups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDeviceGroupsResponse), nil
	}
}

// ListDeviceGroupsInvoker 查询设备组列表
func (c *IoTDAClient) ListDeviceGroupsInvoker(request *model.ListDeviceGroupsRequest) *ListDeviceGroupsInvoker {
	requestDef := GenReqDefForListDeviceGroups()
	return &ListDeviceGroupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDeviceGroup 查询设备组
//
// 应用服务器可调用此接口查询指定设备组详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDeviceGroup(request *model.ShowDeviceGroupRequest) (*model.ShowDeviceGroupResponse, error) {
	requestDef := GenReqDefForShowDeviceGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDeviceGroupResponse), nil
	}
}

// ShowDeviceGroupInvoker 查询设备组
func (c *IoTDAClient) ShowDeviceGroupInvoker(request *model.ShowDeviceGroupRequest) *ShowDeviceGroupInvoker {
	requestDef := GenReqDefForShowDeviceGroup()
	return &ShowDeviceGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDevicesInGroup 查询设备组设备列表
//
// 应用服务器可调用此接口查询指定设备组下的设备列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDevicesInGroup(request *model.ShowDevicesInGroupRequest) (*model.ShowDevicesInGroupResponse, error) {
	requestDef := GenReqDefForShowDevicesInGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDevicesInGroupResponse), nil
	}
}

// ShowDevicesInGroupInvoker 查询设备组设备列表
func (c *IoTDAClient) ShowDevicesInGroupInvoker(request *model.ShowDevicesInGroupRequest) *ShowDevicesInGroupInvoker {
	requestDef := GenReqDefForShowDevicesInGroup()
	return &ShowDevicesInGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDeviceGroup 修改设备组
//
// 应用服务器可调用此接口修改物联网平台中指定设备组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateDeviceGroup(request *model.UpdateDeviceGroupRequest) (*model.UpdateDeviceGroupResponse, error) {
	requestDef := GenReqDefForUpdateDeviceGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDeviceGroupResponse), nil
	}
}

// UpdateDeviceGroupInvoker 修改设备组
func (c *IoTDAClient) UpdateDeviceGroupInvoker(request *model.UpdateDeviceGroupRequest) *UpdateDeviceGroupInvoker {
	requestDef := GenReqDefForUpdateDeviceGroup()
	return &UpdateDeviceGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddDevice 创建设备
//
// 应用服务器可调用此接口在物联网平台创建一个设备，仅在创建后设备才可以接入物联网平台。
//
// - 该接口支持使用gateway_id参数指定在父设备下创建一个子设备，并且支持多级子设备，当前最大支持二级子设备。
// - 该接口同时还支持对设备进行初始配置，接口会读取创建设备请求参数product_id对应的产品详情，如果产品的属性有定义默认值，则会将该属性默认值写入该设备的设备影子中。
// - 用户还可以使用创建设备请求参数shadow字段为设备指定初始配置，指定后将会根据service_id和desired设置的属性值与产品中对应属性的默认值比对，如果不同，则将以shadow字段中设置的属性值为准写入到设备影子中。
// - 该接口仅支持创建单个设备，如需批量注册设备，请参见 [[创建批量任务](https://support.huaweicloud.com/api-iothub/iot_06_v5_0045.html)](tag:hws)[[创建批量任务](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0045.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) AddDevice(request *model.AddDeviceRequest) (*model.AddDeviceResponse, error) {
	requestDef := GenReqDefForAddDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddDeviceResponse), nil
	}
}

// AddDeviceInvoker 创建设备
func (c *IoTDAClient) AddDeviceInvoker(request *model.AddDeviceRequest) *AddDeviceInvoker {
	requestDef := GenReqDefForAddDevice()
	return &AddDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDevice 删除设备
//
// 应用服务器可调用此接口在物联网平台上删除指定设备。若设备下连接了非直连设备，则必须把设备下的非直连设备都删除后，才能删除该设备。该接口仅支持删除单个设备，如需批量删除设备，请参见 [[创建批量任务](https://support.huaweicloud.com/api-iothub/iot_06_v5_0045.html)](tag:hws)[[创建批量任务](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0045.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteDevice(request *model.DeleteDeviceRequest) (*model.DeleteDeviceResponse, error) {
	requestDef := GenReqDefForDeleteDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDeviceResponse), nil
	}
}

// DeleteDeviceInvoker 删除设备
func (c *IoTDAClient) DeleteDeviceInvoker(request *model.DeleteDeviceRequest) *DeleteDeviceInvoker {
	requestDef := GenReqDefForDeleteDevice()
	return &DeleteDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// FreezeDevice 冻结设备
//
// 应用服务器可调用此接口冻结设备，设备冻结后不能再连接上线，可以通过解冻设备接口解除设备冻结。注意，当前仅支持冻结与平台直连的设备。该接口仅支持冻结单个设备，如需批量冻结设备，请参见 [[创建批量任务](https://support.huaweicloud.com/api-iothub/iot_06_v5_0045.html)](tag:hws)[[创建批量任务](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0045.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) FreezeDevice(request *model.FreezeDeviceRequest) (*model.FreezeDeviceResponse, error) {
	requestDef := GenReqDefForFreezeDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.FreezeDeviceResponse), nil
	}
}

// FreezeDeviceInvoker 冻结设备
func (c *IoTDAClient) FreezeDeviceInvoker(request *model.FreezeDeviceRequest) *FreezeDeviceInvoker {
	requestDef := GenReqDefForFreezeDevice()
	return &FreezeDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDeviceGroupsByDevice 查询指定设备加入的设备组列表
//
// 应用服务器可调用此接口查询物联网平台中的某个设备加入的设备组信息列表。仅标准版实例、企业版实例支持该接口调用，基础版不支持。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListDeviceGroupsByDevice(request *model.ListDeviceGroupsByDeviceRequest) (*model.ListDeviceGroupsByDeviceResponse, error) {
	requestDef := GenReqDefForListDeviceGroupsByDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDeviceGroupsByDeviceResponse), nil
	}
}

// ListDeviceGroupsByDeviceInvoker 查询指定设备加入的设备组列表
func (c *IoTDAClient) ListDeviceGroupsByDeviceInvoker(request *model.ListDeviceGroupsByDeviceRequest) *ListDeviceGroupsByDeviceInvoker {
	requestDef := GenReqDefForListDeviceGroupsByDevice()
	return &ListDeviceGroupsByDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDevices 查询设备列表
//
// 应用服务器可调用此接口查询物联网平台中的设备信息列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListDevices(request *model.ListDevicesRequest) (*model.ListDevicesResponse, error) {
	requestDef := GenReqDefForListDevices()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDevicesResponse), nil
	}
}

// ListDevicesInvoker 查询设备列表
func (c *IoTDAClient) ListDevicesInvoker(request *model.ListDevicesRequest) *ListDevicesInvoker {
	requestDef := GenReqDefForListDevices()
	return &ListDevicesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ResetDeviceSecret 重置设备密钥
//
// 应用服务器可调用此接口重置设备密钥，携带指定密钥时平台将设备密钥重置为指定的密钥，不携带密钥时平台将自动生成一个新的随机密钥返回。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ResetDeviceSecret(request *model.ResetDeviceSecretRequest) (*model.ResetDeviceSecretResponse, error) {
	requestDef := GenReqDefForResetDeviceSecret()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetDeviceSecretResponse), nil
	}
}

// ResetDeviceSecretInvoker 重置设备密钥
func (c *IoTDAClient) ResetDeviceSecretInvoker(request *model.ResetDeviceSecretRequest) *ResetDeviceSecretInvoker {
	requestDef := GenReqDefForResetDeviceSecret()
	return &ResetDeviceSecretInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ResetFingerprint 重置设备指纹
//
// 应用服务器可调用此接口重置设备指纹。携带指定设备指纹时将之重置为指定值；不携带时将之置空，后续设备第一次接入时，该设备指纹的值将设置为第一次接入时的证书指纹。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ResetFingerprint(request *model.ResetFingerprintRequest) (*model.ResetFingerprintResponse, error) {
	requestDef := GenReqDefForResetFingerprint()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetFingerprintResponse), nil
	}
}

// ResetFingerprintInvoker 重置设备指纹
func (c *IoTDAClient) ResetFingerprintInvoker(request *model.ResetFingerprintRequest) *ResetFingerprintInvoker {
	requestDef := GenReqDefForResetFingerprint()
	return &ResetFingerprintInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SearchDevices 灵活搜索设备列表
//
// #### 接口说明
//
// 应用服务器使用SQL语句调用该接口，灵活的搜索所需要的设备资源列表
//
// #### 限制
//
// - 仅**标准版实例、企业版实例**支持该接口调用，基础版不支持。
// - 单账号调用该接口的 TPS 限制最大为1/S(每秒1次请求数)
//
// #### 类SQL语法使用说明
//
// 类SQL语句有select、from、where(可选)、order by(可选)、limit子句(可选)组成，长度限制为400个字符。子句里的内容大小写敏感，SQL语句的关键字大小写不敏感。
//
// 示例：
//
// &#x60;&#x60;&#x60;
// select * from device where device_id &#x3D; &#39;as********&#39; limit 0,5
// &#x60;&#x60;&#x60;
//
// ##### SELECT子句
//
// &#x60;&#x60;&#x60;
// select [field]/[count(*)/count(1)] from device
// &#x60;&#x60;&#x60;
//
// 其中field为需要获取的字段，请参考响应参数字段名称，也可填*，获取所有字段。
//
// 如果需要统计搜索的设备个数，请填count(*)或者count(1).
//
// ##### FROM子句
//
// &#x60;&#x60;&#x60;
// from device
// &#x60;&#x60;&#x60;
//
// from后为要查询的资源名，当前支持\&quot;device\&quot;
//
// ##### WHERE子句(可选)
//
// &#x60;&#x60;&#x60;
// WHERE [condition1] AND [condition2]
// &#x60;&#x60;&#x60;
//
// 最多支持5个condition，不支持嵌套；支持的检索字段请参见下面的**搜索条件字段说明**和**支持的运算符**章节
//
// 连接词支持AND、OR，优先级参考标准SQL语法，默认AND优先级高于OR。
//
// ##### LIMIT子句(可选)
//
// &#x60;&#x60;&#x60;
// limit [offset,] rows
// &#x60;&#x60;&#x60;
//
// offset标识搜索的偏移量，rows标识返回搜索结果的最大行数，例如：
//
// - limit n ;示例(select * from device limit 10)
//
//	最大返回n条结果数据
//
//   - limit m,n; 示例(select * from device limit 20,10)
//     搜索偏移量为m，最大返回n条结果数据
//
// ###### 限制
//
//	offset 最大 500， rows最大50，如果不填写limit子句，默认为limit 10
//
// ##### ORDER BY子句(可选)
//
// 用于实现自定义排序，当前支持自定义排序的字段为：\&quot;marker\&quot;。
//
// &#x60;&#x60;&#x60;
// order by marker [asc]/[desc]
// &#x60;&#x60;&#x60;
//
// 子句不填写时默认逻辑为随机排序
//
// #### 搜索条件字段说明
//
// | 字段名      | 类型   | 说明             | 取值范围                                                     |
// | :---------- | :----- | :--------------- | :----------------------------------------------------------- |
// | app_id      | string | 资源空间ID       | 长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。 |
// | device_id   | string | 设备ID           | 长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。 |
// | gateway_id  | string | 网关ID           | 长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。 |
// | product_id  | string | 设备关联的产品ID | 长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。 |
// | device_name | string | 设备名称         | 长度不超过256，只允许中文、字母、数字、以及_?&#39;#().,&amp;%@!-等字符的组合符。 |
// | node_id     | string | 设备标识码       | 长度不超过64，只允许字母、数字、下划线（_）、连接符（-）的组合 |
// | status      | string | 设备的状态       | ONLINE(在线)、OFFLINE(离线)、ABNORMAL(异常)、INACTIVE(未激活)、FROZEN(冻结) |
// | node_type   | string | 设备节点类型     | GATEWAY(直连设备或网关)、ENDPOINT(非直连设备)                |
// | tag_key     | string | 标签键           | 长度不超过64，只允许中文、字母、数字、以及_.-等字符的组合。  |
// | tag_value   | string | 标签值           | 长度不超过128，只允许中文、字母、数字、以及_.-等字符的组合。 |
// | sw_version  | string | 软件版本         | 长度不超过64，只允许字母、数字、下划线（_）、连接符（-）、英文点(.)的组合。 |
// | fw_version  | string | 固件版本         | 长度不超过64，只允许字母、数字、下划线（_）、连接符（-）、英文点(.)的组合。 |
// | group_id    | string | 群组Id           | 长度不超过36，十六进制字符串和连接符（-）的组合              |
// | create_time | string | 设备注册时间     | 格式：yyyy-MM-dd&#39;T&#39;HH:mm:ss.SSS&#39;Z&#39;，如：2015-06-06T12:10:10.000Z |
// | marker      | string | 结果记录ID       | 长度为24的十六进制字符串，如ffffffffffffffffffffffff         |
//
// #### 支持的运算符
//
// | 运算符  | 支持的字段                               |
// | ------- | ---------------------------------------- |
// | &#x3D;       | 所有                                     |
// | !&#x3D;      | 所有                                     |
// | &gt;       | create_time、marker                      |
// | &lt;       | create_time、marker                      |
// | like    | device_name、node_id、tag_key、tag_value |
// | in      | 除tag_key、tag_value以外字段             |
// | not  in | 除tag_key、tag_value以外字段             |
//
// #### SQL 限制
//
// - like: 只支持前缀匹配，不支持后缀匹配或者通配符匹配。前缀匹配不得少于4个字符，且不能包含任何特殊字符(只允许中文、字母、数字、下划线（_）、连接符（-）). 前缀后必须跟上\&quot;%\&quot;结尾。
// - 不支持除了count(*)/count(1)以外的其他任何函数。
// - 不支持其他SQL用法，如嵌套SQL、union、join、别名(Alias)等用法
// - SQL长度限制为400个字符，单个请求条件最大支持5个。
// - 不支持\&quot;null\&quot;和空字符串等条件值匹配
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) SearchDevices(request *model.SearchDevicesRequest) (*model.SearchDevicesResponse, error) {
	requestDef := GenReqDefForSearchDevices()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SearchDevicesResponse), nil
	}
}

// SearchDevicesInvoker 灵活搜索设备列表
func (c *IoTDAClient) SearchDevicesInvoker(request *model.SearchDevicesRequest) *SearchDevicesInvoker {
	requestDef := GenReqDefForSearchDevices()
	return &SearchDevicesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDevice 查询设备
//
// 应用服务器可调用此接口查询物联网平台中指定设备的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDevice(request *model.ShowDeviceRequest) (*model.ShowDeviceResponse, error) {
	requestDef := GenReqDefForShowDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDeviceResponse), nil
	}
}

// ShowDeviceInvoker 查询设备
func (c *IoTDAClient) ShowDeviceInvoker(request *model.ShowDeviceRequest) *ShowDeviceInvoker {
	requestDef := GenReqDefForShowDevice()
	return &ShowDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UnfreezeDevice 解冻设备
//
// 应用服务器可调用此接口解冻设备，解除冻结后，设备可以连接上线。该接口仅支持解冻单个设备，如需批量解冻设备，请参见 [[创建批量任务](https://support.huaweicloud.com/api-iothub/iot_06_v5_0045.html)](tag:hws)[[创建批量任务](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0045.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UnfreezeDevice(request *model.UnfreezeDeviceRequest) (*model.UnfreezeDeviceResponse, error) {
	requestDef := GenReqDefForUnfreezeDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UnfreezeDeviceResponse), nil
	}
}

// UnfreezeDeviceInvoker 解冻设备
func (c *IoTDAClient) UnfreezeDeviceInvoker(request *model.UnfreezeDeviceRequest) *UnfreezeDeviceInvoker {
	requestDef := GenReqDefForUnfreezeDevice()
	return &UnfreezeDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDevice 修改设备
//
// 应用服务器可调用此接口修改物联网平台中指定设备的基本信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateDevice(request *model.UpdateDeviceRequest) (*model.UpdateDeviceResponse, error) {
	requestDef := GenReqDefForUpdateDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDeviceResponse), nil
	}
}

// UpdateDeviceInvoker 修改设备
func (c *IoTDAClient) UpdateDeviceInvoker(request *model.UpdateDeviceRequest) *UpdateDeviceInvoker {
	requestDef := GenReqDefForUpdateDevice()
	return &UpdateDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDeviceProxy 创建设备代理
//
// 应用服务器可调用此接口在物联网平台创建一个动态设备代理规则，用于子设备自主选择网关设备上线和上报消息，即代理组下的任意网关下的子设备均可以通过代理组里其他设备上线([[网关更新子设备状态](https://support.huaweicloud.com/api-iothub/iot_06_v5_3022.html)](tag:hws) [[网关更新子设备状态](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_3022.html)](tag:hws_hk))然后进行数据上报([[网关批量设备属性上报](https://support.huaweicloud.com/api-iothub/iot_06_v5_3006.html)](tag:hws) [[网关更新子设备状态](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_3006.html)](tag:hws_hk))。
// - 单实例最多可以配置10个设备代理
// - 单账号调用该接口的 TPS 限制最大为1/S(每秒1次请求数)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateDeviceProxy(request *model.CreateDeviceProxyRequest) (*model.CreateDeviceProxyResponse, error) {
	requestDef := GenReqDefForCreateDeviceProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDeviceProxyResponse), nil
	}
}

// CreateDeviceProxyInvoker 创建设备代理
func (c *IoTDAClient) CreateDeviceProxyInvoker(request *model.CreateDeviceProxyRequest) *CreateDeviceProxyInvoker {
	requestDef := GenReqDefForCreateDeviceProxy()
	return &CreateDeviceProxyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDeviceProxy 删除设备代理
//
// 应用服务器可调用此接口在物联网平台上删除指定设备代理。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteDeviceProxy(request *model.DeleteDeviceProxyRequest) (*model.DeleteDeviceProxyResponse, error) {
	requestDef := GenReqDefForDeleteDeviceProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDeviceProxyResponse), nil
	}
}

// DeleteDeviceProxyInvoker 删除设备代理
func (c *IoTDAClient) DeleteDeviceProxyInvoker(request *model.DeleteDeviceProxyRequest) *DeleteDeviceProxyInvoker {
	requestDef := GenReqDefForDeleteDeviceProxy()
	return &DeleteDeviceProxyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDeviceProxies 查询设备代理列表
//
// 应用服务器可调用此接口查询物联网平台中的设备代理列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListDeviceProxies(request *model.ListDeviceProxiesRequest) (*model.ListDeviceProxiesResponse, error) {
	requestDef := GenReqDefForListDeviceProxies()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDeviceProxiesResponse), nil
	}
}

// ListDeviceProxiesInvoker 查询设备代理列表
func (c *IoTDAClient) ListDeviceProxiesInvoker(request *model.ListDeviceProxiesRequest) *ListDeviceProxiesInvoker {
	requestDef := GenReqDefForListDeviceProxies()
	return &ListDeviceProxiesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDeviceProxy 查询设备代理详情
//
// 应用服务器可调用此接口查询物联网平台中指定设备代理的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDeviceProxy(request *model.ShowDeviceProxyRequest) (*model.ShowDeviceProxyResponse, error) {
	requestDef := GenReqDefForShowDeviceProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDeviceProxyResponse), nil
	}
}

// ShowDeviceProxyInvoker 查询设备代理详情
func (c *IoTDAClient) ShowDeviceProxyInvoker(request *model.ShowDeviceProxyRequest) *ShowDeviceProxyInvoker {
	requestDef := GenReqDefForShowDeviceProxy()
	return &ShowDeviceProxyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDeviceProxy 修改设备代理
//
// 应用服务器可调用此接口修改物联网平台中指定设备代理的基本信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateDeviceProxy(request *model.UpdateDeviceProxyRequest) (*model.UpdateDeviceProxyResponse, error) {
	requestDef := GenReqDefForUpdateDeviceProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDeviceProxyResponse), nil
	}
}

// UpdateDeviceProxyInvoker 修改设备代理
func (c *IoTDAClient) UpdateDeviceProxyInvoker(request *model.UpdateDeviceProxyRequest) *UpdateDeviceProxyInvoker {
	requestDef := GenReqDefForUpdateDeviceProxy()
	return &UpdateDeviceProxyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDeviceShadow 查询设备影子数据
//
// 应用服务器可调用此接口查询指定设备的设备影子信息，包括对设备的期望属性信息（desired区）和设备最新上报的属性信息（reported区）。
//
// 设备影子介绍：
// 设备影子是一个用于存储和检索设备当前状态信息的JSON文档。
// - 每个设备有且只有一个设备影子，由设备ID唯一标识
// - 设备影子用于存储设备上报的(状态)属性和应用程序期望的设备(状态)属性
// - 无论该设备是否在线，都可以通过该影子获取和设置设备的属性
// - 设备上线或者设备上报属性时，如果desired区和reported区存在差异，则将差异部分下发给设备，配置的预期属性需在产品模型中定义且method具有可写属性“W”才可下发
//
// 限制：
// 设备影子JSON文档中的key不允许特殊字符：点(.)、dollar符号($)、空char(十六进制的ASCII码为00)。如果包含了以上特殊字符则无法正常刷新影子文档。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDeviceShadow(request *model.ShowDeviceShadowRequest) (*model.ShowDeviceShadowResponse, error) {
	requestDef := GenReqDefForShowDeviceShadow()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDeviceShadowResponse), nil
	}
}

// ShowDeviceShadowInvoker 查询设备影子数据
func (c *IoTDAClient) ShowDeviceShadowInvoker(request *model.ShowDeviceShadowRequest) *ShowDeviceShadowInvoker {
	requestDef := GenReqDefForShowDeviceShadow()
	return &ShowDeviceShadowInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDeviceShadowDesiredData 配置设备影子预期数据
//
// 应用服务器可调用此接口配置设备影子的预期属性（desired区），当设备上线或者设备上报属性时把属性下发给设备。
//
// 设备影子介绍：
// 设备影子是一个用于存储和检索设备当前状态信息的JSON文档。
// - 每个设备有且只有一个设备影子，由设备ID唯一标识
// - 设备影子用于存储设备上报的(状态)属性和应用程序期望的设备(状态)属性
// - 无论该设备是否在线，都可以通过该影子获取和设置设备的属性
// - 设备上线或者设备上报属性时，如果desired区和reported区存在差异，则将差异部分下发给设备，配置的预期属性需在产品模型中定义且method具有可写属性“W”才可下发
// - 该接口仅支持配置单个设备的设备影子的预期数据，如需多个设备的设备影子配置，请参见 [[创建批量任务](https://support.huaweicloud.com/api-iothub/iot_06_v5_0045.html)](tag:hws)[[创建批量任务](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0045.html)](tag:hws_hk)。
//
// 限制：
// 设备影子JSON文档中的key不允许特殊字符：点(.)、dollar符号($)、空char(十六进制的ASCII码为00)。如果包含了以上特殊字符则无法正常刷新影子文档。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateDeviceShadowDesiredData(request *model.UpdateDeviceShadowDesiredDataRequest) (*model.UpdateDeviceShadowDesiredDataResponse, error) {
	requestDef := GenReqDefForUpdateDeviceShadowDesiredData()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDeviceShadowDesiredDataResponse), nil
	}
}

// UpdateDeviceShadowDesiredDataInvoker 配置设备影子预期数据
func (c *IoTDAClient) UpdateDeviceShadowDesiredDataInvoker(request *model.UpdateDeviceShadowDesiredDataRequest) *UpdateDeviceShadowDesiredDataInvoker {
	requestDef := GenReqDefForUpdateDeviceShadowDesiredData()
	return &UpdateDeviceShadowDesiredDataInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRoutingFlowControlPolicy 新建数据流转流控策略
//
// 应用服务器可调用此接口在物联网平台创建数据流转流控策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateRoutingFlowControlPolicy(request *model.CreateRoutingFlowControlPolicyRequest) (*model.CreateRoutingFlowControlPolicyResponse, error) {
	requestDef := GenReqDefForCreateRoutingFlowControlPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRoutingFlowControlPolicyResponse), nil
	}
}

// CreateRoutingFlowControlPolicyInvoker 新建数据流转流控策略
func (c *IoTDAClient) CreateRoutingFlowControlPolicyInvoker(request *model.CreateRoutingFlowControlPolicyRequest) *CreateRoutingFlowControlPolicyInvoker {
	requestDef := GenReqDefForCreateRoutingFlowControlPolicy()
	return &CreateRoutingFlowControlPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRoutingFlowControlPolicy 删除数据流转流控策略
//
// 应用服务器可调用此接口在物联网平台删除指定数据流转流控策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteRoutingFlowControlPolicy(request *model.DeleteRoutingFlowControlPolicyRequest) (*model.DeleteRoutingFlowControlPolicyResponse, error) {
	requestDef := GenReqDefForDeleteRoutingFlowControlPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteRoutingFlowControlPolicyResponse), nil
	}
}

// DeleteRoutingFlowControlPolicyInvoker 删除数据流转流控策略
func (c *IoTDAClient) DeleteRoutingFlowControlPolicyInvoker(request *model.DeleteRoutingFlowControlPolicyRequest) *DeleteRoutingFlowControlPolicyInvoker {
	requestDef := GenReqDefForDeleteRoutingFlowControlPolicy()
	return &DeleteRoutingFlowControlPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRoutingFlowControlPolicy 查询数据流转流控策略列表
//
// 应用服务器可调用此接口查询在物联网平台设置的数据流转流控策略列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListRoutingFlowControlPolicy(request *model.ListRoutingFlowControlPolicyRequest) (*model.ListRoutingFlowControlPolicyResponse, error) {
	requestDef := GenReqDefForListRoutingFlowControlPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRoutingFlowControlPolicyResponse), nil
	}
}

// ListRoutingFlowControlPolicyInvoker 查询数据流转流控策略列表
func (c *IoTDAClient) ListRoutingFlowControlPolicyInvoker(request *model.ListRoutingFlowControlPolicyRequest) *ListRoutingFlowControlPolicyInvoker {
	requestDef := GenReqDefForListRoutingFlowControlPolicy()
	return &ListRoutingFlowControlPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRoutingFlowControlPolicy 查询数据流转流控策略
//
// 应用服务器可调用此接口在物联网平台查询指定数据流转流控策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowRoutingFlowControlPolicy(request *model.ShowRoutingFlowControlPolicyRequest) (*model.ShowRoutingFlowControlPolicyResponse, error) {
	requestDef := GenReqDefForShowRoutingFlowControlPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRoutingFlowControlPolicyResponse), nil
	}
}

// ShowRoutingFlowControlPolicyInvoker 查询数据流转流控策略
func (c *IoTDAClient) ShowRoutingFlowControlPolicyInvoker(request *model.ShowRoutingFlowControlPolicyRequest) *ShowRoutingFlowControlPolicyInvoker {
	requestDef := GenReqDefForShowRoutingFlowControlPolicy()
	return &ShowRoutingFlowControlPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateRoutingFlowControlPolicy 修改数据流转流控策略
//
// 应用服务器可调用此接口在物联网平台修改指定数据流转流控策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateRoutingFlowControlPolicy(request *model.UpdateRoutingFlowControlPolicyRequest) (*model.UpdateRoutingFlowControlPolicyResponse, error) {
	requestDef := GenReqDefForUpdateRoutingFlowControlPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateRoutingFlowControlPolicyResponse), nil
	}
}

// UpdateRoutingFlowControlPolicyInvoker 修改数据流转流控策略
func (c *IoTDAClient) UpdateRoutingFlowControlPolicyInvoker(request *model.UpdateRoutingFlowControlPolicyRequest) *UpdateRoutingFlowControlPolicyInvoker {
	requestDef := GenReqDefForUpdateRoutingFlowControlPolicy()
	return &UpdateRoutingFlowControlPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateMessage 下发设备消息
//
// 物联网平台可向设备下发消息，应用服务器可调用此接口向指定设备下发消息，以实现对设备的控制。应用将消息下发给平台后，平台返回应用响应结果，平台再将消息发送给设备。平台返回应用响应结果不一定是设备接收结果，建议用户应用通过订阅[[设备消息状态变更通知](https://support.huaweicloud.com/api-iothub/iot_06_v5_01203.html)](tag:hws)[[设备消息状态变更通知](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_01203.html)](tag:hws_hk)，订阅后平台会将设备接收结果推送给订阅的应用。
// 注意：
// - 此接口适用于MQTT设备消息下发，暂不支持其他协议接入的设备消息下发。
// - 此接口仅支持单个设备消息下发，如需多个设备消息下发，请参见 [[创建批量任务](https://support.huaweicloud.com/api-iothub/iot_06_v5_0045.html)](tag:hws)[[创建批量任务](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0045.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateMessage(request *model.CreateMessageRequest) (*model.CreateMessageResponse, error) {
	requestDef := GenReqDefForCreateMessage()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateMessageResponse), nil
	}
}

// CreateMessageInvoker 下发设备消息
func (c *IoTDAClient) CreateMessageInvoker(request *model.CreateMessageRequest) *CreateMessageInvoker {
	requestDef := GenReqDefForCreateMessage()
	return &CreateMessageInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDeviceMessages 查询设备消息
//
// 应用服务器可调用此接口查询平台下发给设备的消息，平台为每个设备默认最多保存20条消息，超过20条后， 后续的消息会替换下发最早的消息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListDeviceMessages(request *model.ListDeviceMessagesRequest) (*model.ListDeviceMessagesResponse, error) {
	requestDef := GenReqDefForListDeviceMessages()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDeviceMessagesResponse), nil
	}
}

// ListDeviceMessagesInvoker 查询设备消息
func (c *IoTDAClient) ListDeviceMessagesInvoker(request *model.ListDeviceMessagesRequest) *ListDeviceMessagesInvoker {
	requestDef := GenReqDefForListDeviceMessages()
	return &ListDeviceMessagesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDeviceMessage 查询指定消息id的消息
//
// 应用服务器可调用此接口查询平台下发给设备的指定消息id的消息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDeviceMessage(request *model.ShowDeviceMessageRequest) (*model.ShowDeviceMessageResponse, error) {
	requestDef := GenReqDefForShowDeviceMessage()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDeviceMessageResponse), nil
	}
}

// ShowDeviceMessageInvoker 查询指定消息id的消息
func (c *IoTDAClient) ShowDeviceMessageInvoker(request *model.ShowDeviceMessageRequest) *ShowDeviceMessageInvoker {
	requestDef := GenReqDefForShowDeviceMessage()
	return &ShowDeviceMessageInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateOtaPackage 创建OTA升级包
//
// 用户可调用此接口创建升级包关联OBS对象
// 使用前提：使用该API需要您授权设备接入服务(IoTDA)的实例访问对象存储服务(OBS)以及 密钥管理服务(KMS Administrator)的权限。在“[[统一身份认证服务（IAM）](https://console.huaweicloud.com/iam)](tag:hws)[[统一身份认证服务（IAM）](https://console-intl.huaweicloud.com/iam)](tag:hws_hk) - 委托”中将委托名称为iotda_admin_trust的委托授权KMS Administrator和OBS OperateAccess
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateOtaPackage(request *model.CreateOtaPackageRequest) (*model.CreateOtaPackageResponse, error) {
	requestDef := GenReqDefForCreateOtaPackage()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateOtaPackageResponse), nil
	}
}

// CreateOtaPackageInvoker 创建OTA升级包
func (c *IoTDAClient) CreateOtaPackageInvoker(request *model.CreateOtaPackageRequest) *CreateOtaPackageInvoker {
	requestDef := GenReqDefForCreateOtaPackage()
	return &CreateOtaPackageInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteOtaPackage 删除OTA升级包
//
// 用户可调用此接口删除关联OBS对象的升级包信息，不会删除OBS上对象
// 使用前提：使用该API需要您授权设备接入服务(IoTDA)的实例访问对象存储服务(OBS)以及 密钥管理服务(KMS Administrator)的权限。在“[[统一身份认证服务（IAM）](https://console.huaweicloud.com/iam)](tag:hws)[[统一身份认证服务（IAM）](https://console-intl.huaweicloud.com/iam)](tag:hws_hk) - 委托”中将委托名称为iotda_admin_trust的委托授权KMS Administrator和OBS OperateAccess
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteOtaPackage(request *model.DeleteOtaPackageRequest) (*model.DeleteOtaPackageResponse, error) {
	requestDef := GenReqDefForDeleteOtaPackage()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteOtaPackageResponse), nil
	}
}

// DeleteOtaPackageInvoker 删除OTA升级包
func (c *IoTDAClient) DeleteOtaPackageInvoker(request *model.DeleteOtaPackageRequest) *DeleteOtaPackageInvoker {
	requestDef := GenReqDefForDeleteOtaPackage()
	return &DeleteOtaPackageInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListOtaPackageInfo 查询OTA升级包列表
//
// 用户可调用此接口查询关联OBS对象的升级包列表
// 使用前提：使用该API需要您授权设备接入服务(IoTDA)的实例访问对象存储服务(OBS)以及 密钥管理服务(KMS Administrator)的权限。在“[[统一身份认证服务（IAM）](https://console.huaweicloud.com/iam)](tag:hws)[[统一身份认证服务（IAM）](https://console-intl.huaweicloud.com/iam)](tag:hws_hk) - 委托”中将委托名称为iotda_admin_trust的委托授权KMS Administrator和OBS OperateAccess
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListOtaPackageInfo(request *model.ListOtaPackageInfoRequest) (*model.ListOtaPackageInfoResponse, error) {
	requestDef := GenReqDefForListOtaPackageInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOtaPackageInfoResponse), nil
	}
}

// ListOtaPackageInfoInvoker 查询OTA升级包列表
func (c *IoTDAClient) ListOtaPackageInfoInvoker(request *model.ListOtaPackageInfoRequest) *ListOtaPackageInfoInvoker {
	requestDef := GenReqDefForListOtaPackageInfo()
	return &ListOtaPackageInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowOtaPackage 获取OTA升级包详情
//
// 用户可调用此接口查询关联OBS对象的升级包详情
// 使用前提：使用该API需要您授权设备接入服务(IoTDA)的实例访问对象存储服务(OBS)以及 密钥管理服务(KMS Administrator)的权限。在“[[统一身份认证服务（IAM）](https://console.huaweicloud.com/iam)](tag:hws)[[统一身份认证服务（IAM）](https://console-intl.huaweicloud.com/iam)](tag:hws_hk) - 委托”中将委托名称为iotda_admin_trust的委托授权KMS Administrator和OBS OperateAccess
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowOtaPackage(request *model.ShowOtaPackageRequest) (*model.ShowOtaPackageResponse, error) {
	requestDef := GenReqDefForShowOtaPackage()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowOtaPackageResponse), nil
	}
}

// ShowOtaPackageInvoker 获取OTA升级包详情
func (c *IoTDAClient) ShowOtaPackageInvoker(request *model.ShowOtaPackageRequest) *ShowOtaPackageInvoker {
	requestDef := GenReqDefForShowOtaPackage()
	return &ShowOtaPackageInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BindDevicePolicy 绑定设备策略
//
// 应用服务器可调用此接口在物联网平台上为批量设备绑定目标策略，目前支持绑定目标类型为：设备、产品，当目标类型为产品时，该产品下所有设备都会生效。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) BindDevicePolicy(request *model.BindDevicePolicyRequest) (*model.BindDevicePolicyResponse, error) {
	requestDef := GenReqDefForBindDevicePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BindDevicePolicyResponse), nil
	}
}

// BindDevicePolicyInvoker 绑定设备策略
func (c *IoTDAClient) BindDevicePolicyInvoker(request *model.BindDevicePolicyRequest) *BindDevicePolicyInvoker {
	requestDef := GenReqDefForBindDevicePolicy()
	return &BindDevicePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDevicePolicy 创建设备策略
//
// 应用服务器可调用此接口在物联网平台创建一个策略，该策略需要绑定到设备和产品下才能生效。
// - 一个实例最多能创建50个设备策略。
// - 仅**标准版实例、企业版实例**支持该接口调用，基础版不支持。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateDevicePolicy(request *model.CreateDevicePolicyRequest) (*model.CreateDevicePolicyResponse, error) {
	requestDef := GenReqDefForCreateDevicePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDevicePolicyResponse), nil
	}
}

// CreateDevicePolicyInvoker 创建设备策略
func (c *IoTDAClient) CreateDevicePolicyInvoker(request *model.CreateDevicePolicyRequest) *CreateDevicePolicyInvoker {
	requestDef := GenReqDefForCreateDevicePolicy()
	return &CreateDevicePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDevicePolicy 删除设备策略
//
// 应用服务器可调用此接口在物联网平台上删除指定策略，注意：删除策略同时会解绑该策略下所有绑定对象。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteDevicePolicy(request *model.DeleteDevicePolicyRequest) (*model.DeleteDevicePolicyResponse, error) {
	requestDef := GenReqDefForDeleteDevicePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDevicePolicyResponse), nil
	}
}

// DeleteDevicePolicyInvoker 删除设备策略
func (c *IoTDAClient) DeleteDevicePolicyInvoker(request *model.DeleteDevicePolicyRequest) *DeleteDevicePolicyInvoker {
	requestDef := GenReqDefForDeleteDevicePolicy()
	return &DeleteDevicePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDevicePolicies 查询设备策略列表
//
// 应用服务器可调用此接口在物联网平台查询策略列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListDevicePolicies(request *model.ListDevicePoliciesRequest) (*model.ListDevicePoliciesResponse, error) {
	requestDef := GenReqDefForListDevicePolicies()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDevicePoliciesResponse), nil
	}
}

// ListDevicePoliciesInvoker 查询设备策略列表
func (c *IoTDAClient) ListDevicePoliciesInvoker(request *model.ListDevicePoliciesRequest) *ListDevicePoliciesInvoker {
	requestDef := GenReqDefForListDevicePolicies()
	return &ListDevicePoliciesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDevicePolicy 查询设备策略详情
//
// 应用服务器可调用此接口在物联网平台查询指定策略ID的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDevicePolicy(request *model.ShowDevicePolicyRequest) (*model.ShowDevicePolicyResponse, error) {
	requestDef := GenReqDefForShowDevicePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDevicePolicyResponse), nil
	}
}

// ShowDevicePolicyInvoker 查询设备策略详情
func (c *IoTDAClient) ShowDevicePolicyInvoker(request *model.ShowDevicePolicyRequest) *ShowDevicePolicyInvoker {
	requestDef := GenReqDefForShowDevicePolicy()
	return &ShowDevicePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTargetsInDevicePolicy 查询设备策略绑定的目标列表
//
// 应用服务器可调用此接口在物联网平台上查询指定策略ID下绑定的目标列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowTargetsInDevicePolicy(request *model.ShowTargetsInDevicePolicyRequest) (*model.ShowTargetsInDevicePolicyResponse, error) {
	requestDef := GenReqDefForShowTargetsInDevicePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTargetsInDevicePolicyResponse), nil
	}
}

// ShowTargetsInDevicePolicyInvoker 查询设备策略绑定的目标列表
func (c *IoTDAClient) ShowTargetsInDevicePolicyInvoker(request *model.ShowTargetsInDevicePolicyRequest) *ShowTargetsInDevicePolicyInvoker {
	requestDef := GenReqDefForShowTargetsInDevicePolicy()
	return &ShowTargetsInDevicePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UnbindDevicePolicy 解绑设备策略
//
// 应用服务器可调用此接口在物联网平台上解除指定策略下绑定的目标对象。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UnbindDevicePolicy(request *model.UnbindDevicePolicyRequest) (*model.UnbindDevicePolicyResponse, error) {
	requestDef := GenReqDefForUnbindDevicePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UnbindDevicePolicyResponse), nil
	}
}

// UnbindDevicePolicyInvoker 解绑设备策略
func (c *IoTDAClient) UnbindDevicePolicyInvoker(request *model.UnbindDevicePolicyRequest) *UnbindDevicePolicyInvoker {
	requestDef := GenReqDefForUnbindDevicePolicy()
	return &UnbindDevicePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDevicePolicy 更新设备策略信息
//
// 应用服务器可调用此接口在物联网平台更新策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateDevicePolicy(request *model.UpdateDevicePolicyRequest) (*model.UpdateDevicePolicyResponse, error) {
	requestDef := GenReqDefForUpdateDevicePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDevicePolicyResponse), nil
	}
}

// UpdateDevicePolicyInvoker 更新设备策略信息
func (c *IoTDAClient) UpdateDevicePolicyInvoker(request *model.UpdateDevicePolicyRequest) *UpdateDevicePolicyInvoker {
	requestDef := GenReqDefForUpdateDevicePolicy()
	return &UpdateDevicePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateProduct 创建产品
//
// 应用服务器可调用此接口创建产品。此接口仅创建了产品，没有创建和安装插件，如果需要对数据进行编解码，还需要在平台开发和安装插件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateProduct(request *model.CreateProductRequest) (*model.CreateProductResponse, error) {
	requestDef := GenReqDefForCreateProduct()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateProductResponse), nil
	}
}

// CreateProductInvoker 创建产品
func (c *IoTDAClient) CreateProductInvoker(request *model.CreateProductRequest) *CreateProductInvoker {
	requestDef := GenReqDefForCreateProduct()
	return &CreateProductInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteProduct 删除产品
//
// 应用服务器可调用此接口删除已导入物联网平台的指定产品模型。此接口仅删除了产品，未删除关联的插件，在产品下存在设备时，该产品不允许删除。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteProduct(request *model.DeleteProductRequest) (*model.DeleteProductResponse, error) {
	requestDef := GenReqDefForDeleteProduct()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteProductResponse), nil
	}
}

// DeleteProductInvoker 删除产品
func (c *IoTDAClient) DeleteProductInvoker(request *model.DeleteProductRequest) *DeleteProductInvoker {
	requestDef := GenReqDefForDeleteProduct()
	return &DeleteProductInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProducts 查询产品列表
//
// 应用服务器可调用此接口查询已导入物联网平台的产品模型信息列表，了解产品模型的概要信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListProducts(request *model.ListProductsRequest) (*model.ListProductsResponse, error) {
	requestDef := GenReqDefForListProducts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProductsResponse), nil
	}
}

// ListProductsInvoker 查询产品列表
func (c *IoTDAClient) ListProductsInvoker(request *model.ListProductsRequest) *ListProductsInvoker {
	requestDef := GenReqDefForListProducts()
	return &ListProductsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowProduct 查询产品
//
// 应用服务器可调用此接口查询已导入物联网平台的指定产品模型详细信息，包括产品模型的服务、属性、命令等。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowProduct(request *model.ShowProductRequest) (*model.ShowProductResponse, error) {
	requestDef := GenReqDefForShowProduct()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowProductResponse), nil
	}
}

// ShowProductInvoker 查询产品
func (c *IoTDAClient) ShowProductInvoker(request *model.ShowProductRequest) *ShowProductInvoker {
	requestDef := GenReqDefForShowProduct()
	return &ShowProductInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateProduct 修改产品
//
// 应用服务器可调用此接口修改已导入物联网平台的指定产品模型，包括产品模型的服务、属性、命令等。此接口仅修改了产品，未修改和安装插件，如果修改了产品中的service定义，且在平台中有对应的插件，请修改并重新安装插件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateProduct(request *model.UpdateProductRequest) (*model.UpdateProductResponse, error) {
	requestDef := GenReqDefForUpdateProduct()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateProductResponse), nil
	}
}

// UpdateProductInvoker 修改产品
func (c *IoTDAClient) UpdateProductInvoker(request *model.UpdateProductRequest) *UpdateProductInvoker {
	requestDef := GenReqDefForUpdateProduct()
	return &UpdateProductInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProperties 查询设备属性
//
// 设备的产品模型中定义了物联网平台可向设备下发的属性，应用服务器可调用此接口向设备发送指令用以查询设备的实时属性, 并由设备将属性查询的结果同步返回给应用服务器。
// 注意：此接口适用于MQTT设备，暂不支持NB-IoT设备。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListProperties(request *model.ListPropertiesRequest) (*model.ListPropertiesResponse, error) {
	requestDef := GenReqDefForListProperties()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPropertiesResponse), nil
	}
}

// ListPropertiesInvoker 查询设备属性
func (c *IoTDAClient) ListPropertiesInvoker(request *model.ListPropertiesRequest) *ListPropertiesInvoker {
	requestDef := GenReqDefForListProperties()
	return &ListPropertiesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateProperties 修改设备属性
//
// 设备的产品模型中定义了物联网平台可向设备下发的属性，应用服务器可调用此接口向指定设备下发属性。平台负责将属性以同步方式发送给设备，并将设备执行属性结果同步返回。
// 注意：此接口适用于MQTT设备，暂不支持NB-IoT设备。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateProperties(request *model.UpdatePropertiesRequest) (*model.UpdatePropertiesResponse, error) {
	requestDef := GenReqDefForUpdateProperties()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePropertiesResponse), nil
	}
}

// UpdatePropertiesInvoker 修改设备属性
func (c *IoTDAClient) UpdatePropertiesInvoker(request *model.UpdatePropertiesRequest) *UpdatePropertiesInvoker {
	requestDef := GenReqDefForUpdateProperties()
	return &UpdatePropertiesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateProvisioningTemplate 创建预调配模板
//
// 应用服务器可调用此接口在物联网平台创建一个预调配模板。用户的设备未在平台注册时，可以通过预调配模板在设备首次接入物联网平台时将设备信息自动注册到物联网平台。
// - 该预调配模板至少需要绑定到一个设备CA证书下才能生效。
// - 一个实例最多可有10个预调配模板。
// - 仅标准版实例、企业版实例支持该接口调用，基础版不支持。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateProvisioningTemplate(request *model.CreateProvisioningTemplateRequest) (*model.CreateProvisioningTemplateResponse, error) {
	requestDef := GenReqDefForCreateProvisioningTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateProvisioningTemplateResponse), nil
	}
}

// CreateProvisioningTemplateInvoker 创建预调配模板
func (c *IoTDAClient) CreateProvisioningTemplateInvoker(request *model.CreateProvisioningTemplateRequest) *CreateProvisioningTemplateInvoker {
	requestDef := GenReqDefForCreateProvisioningTemplate()
	return &CreateProvisioningTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteProvisioningTemplate 删除预调配模板
//
// 应用服务器可调用此接口在物联网平台上删除指定预调配模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteProvisioningTemplate(request *model.DeleteProvisioningTemplateRequest) (*model.DeleteProvisioningTemplateResponse, error) {
	requestDef := GenReqDefForDeleteProvisioningTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteProvisioningTemplateResponse), nil
	}
}

// DeleteProvisioningTemplateInvoker 删除预调配模板
func (c *IoTDAClient) DeleteProvisioningTemplateInvoker(request *model.DeleteProvisioningTemplateRequest) *DeleteProvisioningTemplateInvoker {
	requestDef := GenReqDefForDeleteProvisioningTemplate()
	return &DeleteProvisioningTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProvisioningTemplates 查询预调配模板列表
//
// 应用服务器可调用此接口在物联网平台查询预调配模板列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListProvisioningTemplates(request *model.ListProvisioningTemplatesRequest) (*model.ListProvisioningTemplatesResponse, error) {
	requestDef := GenReqDefForListProvisioningTemplates()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProvisioningTemplatesResponse), nil
	}
}

// ListProvisioningTemplatesInvoker 查询预调配模板列表
func (c *IoTDAClient) ListProvisioningTemplatesInvoker(request *model.ListProvisioningTemplatesRequest) *ListProvisioningTemplatesInvoker {
	requestDef := GenReqDefForListProvisioningTemplates()
	return &ListProvisioningTemplatesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowProvisioningTemplate 查询预调配模板详情
//
// 应用服务器可调用此接口在物联网平台查询指定预调配模板ID的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowProvisioningTemplate(request *model.ShowProvisioningTemplateRequest) (*model.ShowProvisioningTemplateResponse, error) {
	requestDef := GenReqDefForShowProvisioningTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowProvisioningTemplateResponse), nil
	}
}

// ShowProvisioningTemplateInvoker 查询预调配模板详情
func (c *IoTDAClient) ShowProvisioningTemplateInvoker(request *model.ShowProvisioningTemplateRequest) *ShowProvisioningTemplateInvoker {
	requestDef := GenReqDefForShowProvisioningTemplate()
	return &ShowProvisioningTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateProvisioningTemplate 更新指定id的预调配模板信息
//
// 应用服务器可调用此接口在物联网平台更新指定id的预调配模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateProvisioningTemplate(request *model.UpdateProvisioningTemplateRequest) (*model.UpdateProvisioningTemplateResponse, error) {
	requestDef := GenReqDefForUpdateProvisioningTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateProvisioningTemplateResponse), nil
	}
}

// UpdateProvisioningTemplateInvoker 更新指定id的预调配模板信息
func (c *IoTDAClient) UpdateProvisioningTemplateInvoker(request *model.UpdateProvisioningTemplateRequest) *UpdateProvisioningTemplateInvoker {
	requestDef := GenReqDefForUpdateProvisioningTemplate()
	return &UpdateProvisioningTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRoutingRule 创建规则触发条件
//
// 应用服务器可调用此接口在物联网平台创建一条规则触发条件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateRoutingRule(request *model.CreateRoutingRuleRequest) (*model.CreateRoutingRuleResponse, error) {
	requestDef := GenReqDefForCreateRoutingRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRoutingRuleResponse), nil
	}
}

// CreateRoutingRuleInvoker 创建规则触发条件
func (c *IoTDAClient) CreateRoutingRuleInvoker(request *model.CreateRoutingRuleRequest) *CreateRoutingRuleInvoker {
	requestDef := GenReqDefForCreateRoutingRule()
	return &CreateRoutingRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRuleAction 创建规则动作
//
// 应用服务器可调用此接口在物联网平台创建一条规则动作。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateRuleAction(request *model.CreateRuleActionRequest) (*model.CreateRuleActionResponse, error) {
	requestDef := GenReqDefForCreateRuleAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRuleActionResponse), nil
	}
}

// CreateRuleActionInvoker 创建规则动作
func (c *IoTDAClient) CreateRuleActionInvoker(request *model.CreateRuleActionRequest) *CreateRuleActionInvoker {
	requestDef := GenReqDefForCreateRuleAction()
	return &CreateRuleActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRoutingRule 删除规则触发条件
//
// 应用服务器可调用此接口删除物联网平台中的指定规则条件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteRoutingRule(request *model.DeleteRoutingRuleRequest) (*model.DeleteRoutingRuleResponse, error) {
	requestDef := GenReqDefForDeleteRoutingRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteRoutingRuleResponse), nil
	}
}

// DeleteRoutingRuleInvoker 删除规则触发条件
func (c *IoTDAClient) DeleteRoutingRuleInvoker(request *model.DeleteRoutingRuleRequest) *DeleteRoutingRuleInvoker {
	requestDef := GenReqDefForDeleteRoutingRule()
	return &DeleteRoutingRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRuleAction 删除规则动作
//
// 应用服务器可调用此接口删除物联网平台中的指定规则动作。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteRuleAction(request *model.DeleteRuleActionRequest) (*model.DeleteRuleActionResponse, error) {
	requestDef := GenReqDefForDeleteRuleAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteRuleActionResponse), nil
	}
}

// DeleteRuleActionInvoker 删除规则动作
func (c *IoTDAClient) DeleteRuleActionInvoker(request *model.DeleteRuleActionRequest) *DeleteRuleActionInvoker {
	requestDef := GenReqDefForDeleteRuleAction()
	return &DeleteRuleActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRoutingRules 查询规则条件列表
//
// 应用服务器可调用此接口查询物联网平台中设置的规则条件列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListRoutingRules(request *model.ListRoutingRulesRequest) (*model.ListRoutingRulesResponse, error) {
	requestDef := GenReqDefForListRoutingRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRoutingRulesResponse), nil
	}
}

// ListRoutingRulesInvoker 查询规则条件列表
func (c *IoTDAClient) ListRoutingRulesInvoker(request *model.ListRoutingRulesRequest) *ListRoutingRulesInvoker {
	requestDef := GenReqDefForListRoutingRules()
	return &ListRoutingRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRuleActions 查询规则动作列表
//
// 应用服务器可调用此接口查询物联网平台中设置的规则动作列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListRuleActions(request *model.ListRuleActionsRequest) (*model.ListRuleActionsResponse, error) {
	requestDef := GenReqDefForListRuleActions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRuleActionsResponse), nil
	}
}

// ListRuleActionsInvoker 查询规则动作列表
func (c *IoTDAClient) ListRuleActionsInvoker(request *model.ListRuleActionsRequest) *ListRuleActionsInvoker {
	requestDef := GenReqDefForListRuleActions()
	return &ListRuleActionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRoutingRule 查询规则条件
//
// 应用服务器可调用此接口查询物联网平台中指定规则条件的配置信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowRoutingRule(request *model.ShowRoutingRuleRequest) (*model.ShowRoutingRuleResponse, error) {
	requestDef := GenReqDefForShowRoutingRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRoutingRuleResponse), nil
	}
}

// ShowRoutingRuleInvoker 查询规则条件
func (c *IoTDAClient) ShowRoutingRuleInvoker(request *model.ShowRoutingRuleRequest) *ShowRoutingRuleInvoker {
	requestDef := GenReqDefForShowRoutingRule()
	return &ShowRoutingRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRuleAction 查询规则动作
//
// 应用服务器可调用此接口查询物联网平台中指定规则动作的配置信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowRuleAction(request *model.ShowRuleActionRequest) (*model.ShowRuleActionResponse, error) {
	requestDef := GenReqDefForShowRuleAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRuleActionResponse), nil
	}
}

// ShowRuleActionInvoker 查询规则动作
func (c *IoTDAClient) ShowRuleActionInvoker(request *model.ShowRuleActionRequest) *ShowRuleActionInvoker {
	requestDef := GenReqDefForShowRuleAction()
	return &ShowRuleActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateRoutingRule 修改规则触发条件
//
// 应用服务器可调用此接口修改物联网平台中指定规则条件的配置参数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateRoutingRule(request *model.UpdateRoutingRuleRequest) (*model.UpdateRoutingRuleResponse, error) {
	requestDef := GenReqDefForUpdateRoutingRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateRoutingRuleResponse), nil
	}
}

// UpdateRoutingRuleInvoker 修改规则触发条件
func (c *IoTDAClient) UpdateRoutingRuleInvoker(request *model.UpdateRoutingRuleRequest) *UpdateRoutingRuleInvoker {
	requestDef := GenReqDefForUpdateRoutingRule()
	return &UpdateRoutingRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateRuleAction 修改规则动作
//
// 应用服务器可调用此接口修改物联网平台中指定规则动作的配置。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateRuleAction(request *model.UpdateRuleActionRequest) (*model.UpdateRuleActionResponse, error) {
	requestDef := GenReqDefForUpdateRuleAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateRuleActionResponse), nil
	}
}

// UpdateRuleActionInvoker 修改规则动作
func (c *IoTDAClient) UpdateRuleActionInvoker(request *model.UpdateRuleActionRequest) *UpdateRuleActionInvoker {
	requestDef := GenReqDefForUpdateRuleAction()
	return &UpdateRuleActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeRuleStatus 修改规则状态
//
// 应用服务器可调用此接口修改物联网平台中指定规则的状态，激活或者去激活规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ChangeRuleStatus(request *model.ChangeRuleStatusRequest) (*model.ChangeRuleStatusResponse, error) {
	requestDef := GenReqDefForChangeRuleStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeRuleStatusResponse), nil
	}
}

// ChangeRuleStatusInvoker 修改规则状态
func (c *IoTDAClient) ChangeRuleStatusInvoker(request *model.ChangeRuleStatusRequest) *ChangeRuleStatusInvoker {
	requestDef := GenReqDefForChangeRuleStatus()
	return &ChangeRuleStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRule 创建规则
//
// 应用服务器可调用此接口在物联网平台创建一条规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CreateRule(request *model.CreateRuleRequest) (*model.CreateRuleResponse, error) {
	requestDef := GenReqDefForCreateRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRuleResponse), nil
	}
}

// CreateRuleInvoker 创建规则
func (c *IoTDAClient) CreateRuleInvoker(request *model.CreateRuleRequest) *CreateRuleInvoker {
	requestDef := GenReqDefForCreateRule()
	return &CreateRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRule 删除规则
//
// 应用服务器可调用此接口删除物联网平台中的指定规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteRule(request *model.DeleteRuleRequest) (*model.DeleteRuleResponse, error) {
	requestDef := GenReqDefForDeleteRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteRuleResponse), nil
	}
}

// DeleteRuleInvoker 删除规则
func (c *IoTDAClient) DeleteRuleInvoker(request *model.DeleteRuleRequest) *DeleteRuleInvoker {
	requestDef := GenReqDefForDeleteRule()
	return &DeleteRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRules 查询规则列表
//
// 应用服务器可调用此接口查询物联网平台中设置的规则列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListRules(request *model.ListRulesRequest) (*model.ListRulesResponse, error) {
	requestDef := GenReqDefForListRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRulesResponse), nil
	}
}

// ListRulesInvoker 查询规则列表
func (c *IoTDAClient) ListRulesInvoker(request *model.ListRulesRequest) *ListRulesInvoker {
	requestDef := GenReqDefForListRules()
	return &ListRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRule 查询规则
//
// 应用服务器可调用此接口查询物联网平台中指定规则的配置信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowRule(request *model.ShowRuleRequest) (*model.ShowRuleResponse, error) {
	requestDef := GenReqDefForShowRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRuleResponse), nil
	}
}

// ShowRuleInvoker 查询规则
func (c *IoTDAClient) ShowRuleInvoker(request *model.ShowRuleRequest) *ShowRuleInvoker {
	requestDef := GenReqDefForShowRule()
	return &ShowRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateRule 修改规则
//
// 应用服务器可调用此接口修改物联网平台中指定规则的配置。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UpdateRule(request *model.UpdateRuleRequest) (*model.UpdateRuleResponse, error) {
	requestDef := GenReqDefForUpdateRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateRuleResponse), nil
	}
}

// UpdateRuleInvoker 修改规则
func (c *IoTDAClient) UpdateRuleInvoker(request *model.UpdateRuleRequest) *UpdateRuleInvoker {
	requestDef := GenReqDefForUpdateRule()
	return &UpdateRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListResourcesByTags 按标签查询资源
//
// 应用服务器可调用此接口查询绑定了指定标签的资源。当前支持标签的资源有Device(设备)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListResourcesByTags(request *model.ListResourcesByTagsRequest) (*model.ListResourcesByTagsResponse, error) {
	requestDef := GenReqDefForListResourcesByTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListResourcesByTagsResponse), nil
	}
}

// ListResourcesByTagsInvoker 按标签查询资源
func (c *IoTDAClient) ListResourcesByTagsInvoker(request *model.ListResourcesByTagsRequest) *ListResourcesByTagsInvoker {
	requestDef := GenReqDefForListResourcesByTags()
	return &ListResourcesByTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// TagDevice 绑定标签
//
// 应用服务器可调用此接口为指定资源绑定标签。当前支持标签的资源有Device(设备)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) TagDevice(request *model.TagDeviceRequest) (*model.TagDeviceResponse, error) {
	requestDef := GenReqDefForTagDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.TagDeviceResponse), nil
	}
}

// TagDeviceInvoker 绑定标签
func (c *IoTDAClient) TagDeviceInvoker(request *model.TagDeviceRequest) *TagDeviceInvoker {
	requestDef := GenReqDefForTagDevice()
	return &TagDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UntagDevice 解绑标签
//
// 应用服务器可调用此接口为指定资源解绑标签。当前支持标签的资源有Device(设备)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) UntagDevice(request *model.UntagDeviceRequest) (*model.UntagDeviceResponse, error) {
	requestDef := GenReqDefForUntagDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UntagDeviceResponse), nil
	}
}

// UntagDeviceInvoker 解绑标签
func (c *IoTDAClient) UntagDeviceInvoker(request *model.UntagDeviceRequest) *UntagDeviceInvoker {
	requestDef := GenReqDefForUntagDevice()
	return &UntagDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddTunnel 创建设备隧道
//
// 用户可以通过该接口创建隧道（WebSocket协议），应用服务器和设备可以通过该隧道进行数据传输。
//
// - 该API接口在基础版不支持。
// - 该API调用后平台会向对应的MQTT/MQTTS设备下发隧道地址及密钥，同时给应用服务器也返回隧道地址及密钥，设备可以通过该地址及密钥创建WebSocket协议连接。
// - 一个设备无法创建多个隧道。
// - 具体应用可见“设备远程登录”功能，请参见[[设备远程登录](https://support.huaweicloud.com/usermanual-iothub/iot_01_00301.html)](tag:hws)[[设备远程登录](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_00301.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) AddTunnel(request *model.AddTunnelRequest) (*model.AddTunnelResponse, error) {
	requestDef := GenReqDefForAddTunnel()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddTunnelResponse), nil
	}
}

// AddTunnelInvoker 创建设备隧道
func (c *IoTDAClient) AddTunnelInvoker(request *model.AddTunnelRequest) *AddTunnelInvoker {
	requestDef := GenReqDefForAddTunnel()
	return &AddTunnelInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CloseDeviceTunnel 关闭设备隧道
//
// 应用服务器可通过该接口关闭某个设备隧道。关闭后可以再次连接。
// - 该API接口在基础版不支持。
// - 具体应用可见“设备远程登录”功能，请参见[[设备远程登录](https://support.huaweicloud.com/usermanual-iothub/iot_01_00301.html)](tag:hws)[[设备远程登录](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_00301.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) CloseDeviceTunnel(request *model.CloseDeviceTunnelRequest) (*model.CloseDeviceTunnelResponse, error) {
	requestDef := GenReqDefForCloseDeviceTunnel()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CloseDeviceTunnelResponse), nil
	}
}

// CloseDeviceTunnelInvoker 关闭设备隧道
func (c *IoTDAClient) CloseDeviceTunnelInvoker(request *model.CloseDeviceTunnelRequest) *CloseDeviceTunnelInvoker {
	requestDef := GenReqDefForCloseDeviceTunnel()
	return &CloseDeviceTunnelInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDeviceTunnel 删除设备隧道
//
// 用户可通过该接口删除某个设备隧道。删除后该通道不存在，无法再次连接。
// - 该API接口在基础版不支持。
// - 具体应用可见“设备远程登录”功能，请参见[[设备远程登录](https://support.huaweicloud.com/usermanual-iothub/iot_01_00301.html)](tag:hws)[[设备远程登录](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_00301.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) DeleteDeviceTunnel(request *model.DeleteDeviceTunnelRequest) (*model.DeleteDeviceTunnelResponse, error) {
	requestDef := GenReqDefForDeleteDeviceTunnel()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDeviceTunnelResponse), nil
	}
}

// DeleteDeviceTunnelInvoker 删除设备隧道
func (c *IoTDAClient) DeleteDeviceTunnelInvoker(request *model.DeleteDeviceTunnelRequest) *DeleteDeviceTunnelInvoker {
	requestDef := GenReqDefForDeleteDeviceTunnel()
	return &DeleteDeviceTunnelInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDeviceTunnels 查询设备所有隧道
//
// 用户可通过该接口查询某项目下的所有设备隧道，以实现对设备管理。应用服务器可通过此接口向平台查询设备隧道建立的情况。
// - 该API接口在基础版不支持。
// - 具体应用可见“设备远程登录”功能，请参见[[设备远程登录](https://support.huaweicloud.com/usermanual-iothub/iot_01_00301.html)](tag:hws)[[设备远程登录](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_00301.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ListDeviceTunnels(request *model.ListDeviceTunnelsRequest) (*model.ListDeviceTunnelsResponse, error) {
	requestDef := GenReqDefForListDeviceTunnels()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDeviceTunnelsResponse), nil
	}
}

// ListDeviceTunnelsInvoker 查询设备所有隧道
func (c *IoTDAClient) ListDeviceTunnelsInvoker(request *model.ListDeviceTunnelsRequest) *ListDeviceTunnelsInvoker {
	requestDef := GenReqDefForListDeviceTunnels()
	return &ListDeviceTunnelsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDeviceTunnel 查询设备隧道
//
// 用户可通过该接口查询某项目中的某个设备隧道，查看该设备隧道的信息与连接情况。应用服务器可调用此接口向平台查询设备隧道建立情况。
// - 该API接口在基础版不支持。
// - 具体应用可见“设备远程登录”功能，请参见[[设备远程登录](https://support.huaweicloud.com/usermanual-iothub/iot_01_00301.html)](tag:hws)[[设备远程登录](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_00301.html)](tag:hws_hk)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IoTDAClient) ShowDeviceTunnel(request *model.ShowDeviceTunnelRequest) (*model.ShowDeviceTunnelResponse, error) {
	requestDef := GenReqDefForShowDeviceTunnel()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDeviceTunnelResponse), nil
	}
}

// ShowDeviceTunnelInvoker 查询设备隧道
func (c *IoTDAClient) ShowDeviceTunnelInvoker(request *model.ShowDeviceTunnelRequest) *ShowDeviceTunnelInvoker {
	requestDef := GenReqDefForShowDeviceTunnel()
	return &ShowDeviceTunnelInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
