package v5

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
)

type IoTDAClient struct {
	HcClient *http_client.HcHttpClient
}

func NewIoTDAClient(hcClient *http_client.HcHttpClient) *IoTDAClient {
	return &IoTDAClient{HcClient: hcClient}
}

func IoTDAClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder().WithDerivedAuthServiceName("iotdm")
	return builder
}

// CreateAccessCode 生成接入凭证
//
// 接入凭证是用于客户端使用AMQP等协议与平台建链的一个认证凭据。只保留一条记录，如果重复调用只会重置接入凭证，使得之前的失效。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// CreateAsyncCommand 下发异步设备命令
//
// 设备的产品模型中定义了物联网平台可向设备下发的命令，应用服务器可调用此接口向指定设备下发异步命令，以实现对设备的控制。平台负责将命令发送给设备，并将设备执行命令结果异步通知应用服务器。 命令执行结果支持灵活的数据流转，应用服务器通过调用物联网平台的创建规则触发条件（Resource:device.command.status，Event:update）、创建规则动作并激活规则后，当命令状态变更时，物联网平台会根据规则将结果发送到规则指定的服务器，如用户自定义的HTTP服务器，AMQP服务器，以及华为云的其他储存服务器等, 详情参考[设备命令状态变更通知](https://support.huaweicloud.com/api-iothub/iot_06_v5_01212.html)。注意：此接口适用于NB设备异步命令下发，暂不支持其他协议类型设备命令下发。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// CreateBatchTask 创建批量任务
//
// 应用服务器可调用此接口为创建批量处理任务，对多个设备进行批量操作。当前支持批量软固件升级、批量创建设备、批量删除设备、批量冻结设备、批量解冻设备、批量创建命令、批量创建消息任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ListBatchTasks 查询批量任务列表
//
// 应用服务器可调用此接口查询物联网平台中批量任务列表，每一个任务又包括具体的任务内容、任务状态、任务完成情况统计等。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ShowBatchTask 查询批量任务
//
// 应用服务器可调用此接口查询物联网平台中指定批量任务的信息，包括任务内容、任务状态、任务完成情况统计以及子任务列表等。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// DeleteBatchTaskFile 删除批量任务文件
//
// 应用服务器可调用此接口删除批量任务文件。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// AddCertificate 上传设备CA证书
//
// 应用服务器可调用此接口在物联网平台上传设备的CA证书
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 应用服务器可调用此接口在物联网平台删除设备的CA证书
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 应用服务器可调用此接口在物联网平台获取设备的CA证书列表
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// CreateCommand 下发设备命令
//
// 设备的产品模型中定义了物联网平台可向设备下发的命令，应用服务器可调用此接口向指定设备下发命令，以实现对设备的同步控制。平台负责将命令以同步方式发送给设备，并将设备执行命令结果同步返回, 如果设备没有响应，平台会返回给应用服务器超时，平台超时间是20秒。注意：此接口适用于MQTT设备同步命令下发，暂不支持NB-IoT设备命令下发。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// AddDeviceGroup 添加设备组
//
// 应用服务器可调用此接口新建设备组，一个华为云账号下最多可有1,000个分组，包括父分组和子分组。设备组的最大层级关系不超过5层，即群组形成的关系树最大深度不超过5。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 应用服务器可调用此接口在物联网平台上删除指定设备。若设备下连接了非直连设备，则必须把设备下的非直连设备都删除后，才能删除该设备。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 应用服务器可调用此接口冻结设备，设备冻结后不能再连接上线，可以通过解冻设备接口解除设备冻结。注意，当前仅支持冻结与平台直连的设备。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ListDevices 查询设备列表
//
// 应用服务器可调用此接口查询物联网平台中的设备信息列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ShowDevice 查询设备
//
// 应用服务器可调用此接口查询物联网平台中指定设备的详细信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 应用服务器可调用此接口解冻设备，解除冻结后，设备可以连接上线。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ShowDeviceShadow 查询设备影子数据
//
// 应用服务器可调用此接口查询指定设备的设备影子信息，包括对设备的期望属性信息（desired区）和设备最新上报的属性信息（reported区）。
//
// 设备影子介绍：
// 设备影子是一个用于存储和检索设备当前状态信息的JSON文档。
// - 每个设备有且只有一个设备影子，由设备ID唯一标识
// - 设备影子仅保存最近一次设备的上报数据和预期数据
// - 无论该设备是否在线，都可以通过该影子获取和设置设备的属性
// - 设备上线或者设备上报属性时，如果desired区和reported区存在差异，则将差异部分下发给设备，配置的预期属性需在产品模型中定义且method具有可写属性“W”才可下发
//
// 限制：
// 设备影子JSON文档中的key不允许特殊字符：点(.)、dollar符号($)、空char(十六进制的ASCII码为00)。如果包含了以上特殊字符则无法正常刷新影子文档。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// - 设备影子仅保存最近一次设备的上报数据和预期数据
// - 无论该设备是否在线，都可以通过该影子获取和设置设备的属性
// - 设备上线或者设备上报属性时，如果desired区和reported区存在差异，则将差异部分下发给设备，配置的预期属性需在产品模型中定义且method具有可写属性“W”才可下发
//
// 限制：
// 设备影子JSON文档中的key不允许特殊字符：点(.)、dollar符号($)、空char(十六进制的ASCII码为00)。如果包含了以上特殊字符则无法正常刷新影子文档。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// CreateMessage 下发设备消息
//
// 物联网平台可向设备下发消息，应用服务器可调用此接口向指定设备下发消息，以实现对设备的控制。应用将消息下发给平台后，平台返回应用响应结果，平台再将消息发送给设备。注意：此接口适用于MQTT设备消息下发，暂不支持其他协议接入的设备消息下发。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// CreateProduct 创建产品
//
// 应用服务器可调用此接口创建产品。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 应用服务器可调用此接口删除已导入物联网平台的指定产品模型。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 应用服务器可调用此接口修改已导入物联网平台的指定产品模型，包括产品模型的服务、属性、命令等。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 设备的产品模型中定义了物联网平台可向设备下发的属性，应用服务器可调用此接口向设备发送指令用以查询设备的实时属性, 并由设备将属性查询的结果同步返回给应用服务器。注意：此接口适用于MQTT设备，暂不支持NB-IoT设备。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 设备的产品模型中定义了物联网平台可向设备下发的属性，应用服务器可调用此接口向指定设备下发属性。平台负责将属性以同步方式发送给设备，并将设备执行属性结果同步返回。注意：此接口适用于MQTT设备，暂不支持NB-IoT设备。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// CreateRoutingRule 创建规则触发条件
//
// 应用服务器可调用此接口在物联网平台创建一条规则触发条件。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
