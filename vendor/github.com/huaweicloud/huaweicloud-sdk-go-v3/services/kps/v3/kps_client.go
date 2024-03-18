package v3

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kps/v3/model"
)

type KpsClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewKpsClient(hcClient *httpclient.HcHttpClient) *KpsClient {
	return &KpsClient{HcClient: hcClient}
}

func KpsClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder()
	return builder
}

// AssociateKeypair 绑定SSH密钥对
//
// 给指定的虚拟机绑定（替换或重置，替换需提供虚拟机已配置的SSH密钥对私钥；重置不需要提供虚拟机的SSH密钥对私钥）新的SSH密钥对。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) AssociateKeypair(request *model.AssociateKeypairRequest) (*model.AssociateKeypairResponse, error) {
	requestDef := GenReqDefForAssociateKeypair()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociateKeypairResponse), nil
	}
}

// AssociateKeypairInvoker 绑定SSH密钥对
func (c *KpsClient) AssociateKeypairInvoker(request *model.AssociateKeypairRequest) *AssociateKeypairInvoker {
	requestDef := GenReqDefForAssociateKeypair()
	return &AssociateKeypairInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchAssociateKeypair 批量绑定SSH密钥对
//
// 给指定的虚拟机批量绑定新的SSH密钥对。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) BatchAssociateKeypair(request *model.BatchAssociateKeypairRequest) (*model.BatchAssociateKeypairResponse, error) {
	requestDef := GenReqDefForBatchAssociateKeypair()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchAssociateKeypairResponse), nil
	}
}

// BatchAssociateKeypairInvoker 批量绑定SSH密钥对
func (c *KpsClient) BatchAssociateKeypairInvoker(request *model.BatchAssociateKeypairRequest) *BatchAssociateKeypairInvoker {
	requestDef := GenReqDefForBatchAssociateKeypair()
	return &BatchAssociateKeypairInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ClearPrivateKey 清除私钥
//
// 清除SSH密钥对私钥。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) ClearPrivateKey(request *model.ClearPrivateKeyRequest) (*model.ClearPrivateKeyResponse, error) {
	requestDef := GenReqDefForClearPrivateKey()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ClearPrivateKeyResponse), nil
	}
}

// ClearPrivateKeyInvoker 清除私钥
func (c *KpsClient) ClearPrivateKeyInvoker(request *model.ClearPrivateKeyRequest) *ClearPrivateKeyInvoker {
	requestDef := GenReqDefForClearPrivateKey()
	return &ClearPrivateKeyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateKeypair 创建和导入SSH密钥对
//
// 创建和导入SSH密钥对
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) CreateKeypair(request *model.CreateKeypairRequest) (*model.CreateKeypairResponse, error) {
	requestDef := GenReqDefForCreateKeypair()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateKeypairResponse), nil
	}
}

// CreateKeypairInvoker 创建和导入SSH密钥对
func (c *KpsClient) CreateKeypairInvoker(request *model.CreateKeypairRequest) *CreateKeypairInvoker {
	requestDef := GenReqDefForCreateKeypair()
	return &CreateKeypairInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAllFailedTask 删除所有失败的任务
//
// 删除操作失败的任务信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) DeleteAllFailedTask(request *model.DeleteAllFailedTaskRequest) (*model.DeleteAllFailedTaskResponse, error) {
	requestDef := GenReqDefForDeleteAllFailedTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAllFailedTaskResponse), nil
	}
}

// DeleteAllFailedTaskInvoker 删除所有失败的任务
func (c *KpsClient) DeleteAllFailedTaskInvoker(request *model.DeleteAllFailedTaskRequest) *DeleteAllFailedTaskInvoker {
	requestDef := GenReqDefForDeleteAllFailedTask()
	return &DeleteAllFailedTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteFailedTask 删除失败的任务
//
// 删除失败的任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) DeleteFailedTask(request *model.DeleteFailedTaskRequest) (*model.DeleteFailedTaskResponse, error) {
	requestDef := GenReqDefForDeleteFailedTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteFailedTaskResponse), nil
	}
}

// DeleteFailedTaskInvoker 删除失败的任务
func (c *KpsClient) DeleteFailedTaskInvoker(request *model.DeleteFailedTaskRequest) *DeleteFailedTaskInvoker {
	requestDef := GenReqDefForDeleteFailedTask()
	return &DeleteFailedTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteKeypair 删除SSH密钥对
//
// 删除SSH密钥对。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) DeleteKeypair(request *model.DeleteKeypairRequest) (*model.DeleteKeypairResponse, error) {
	requestDef := GenReqDefForDeleteKeypair()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteKeypairResponse), nil
	}
}

// DeleteKeypairInvoker 删除SSH密钥对
func (c *KpsClient) DeleteKeypairInvoker(request *model.DeleteKeypairRequest) *DeleteKeypairInvoker {
	requestDef := GenReqDefForDeleteKeypair()
	return &DeleteKeypairInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DisassociateKeypair 解绑SSH密钥对
//
// 给指定的虚拟机解除绑定SSH密钥对并恢复SSH密码登录。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) DisassociateKeypair(request *model.DisassociateKeypairRequest) (*model.DisassociateKeypairResponse, error) {
	requestDef := GenReqDefForDisassociateKeypair()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DisassociateKeypairResponse), nil
	}
}

// DisassociateKeypairInvoker 解绑SSH密钥对
func (c *KpsClient) DisassociateKeypairInvoker(request *model.DisassociateKeypairRequest) *DisassociateKeypairInvoker {
	requestDef := GenReqDefForDisassociateKeypair()
	return &DisassociateKeypairInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ExportPrivateKey 导出私钥
//
// 导出指定密钥对的私钥。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) ExportPrivateKey(request *model.ExportPrivateKeyRequest) (*model.ExportPrivateKeyResponse, error) {
	requestDef := GenReqDefForExportPrivateKey()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ExportPrivateKeyResponse), nil
	}
}

// ExportPrivateKeyInvoker 导出私钥
func (c *KpsClient) ExportPrivateKeyInvoker(request *model.ExportPrivateKeyRequest) *ExportPrivateKeyInvoker {
	requestDef := GenReqDefForExportPrivateKey()
	return &ExportPrivateKeyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ImportPrivateKey 导入私钥
//
// 导入私钥到指定密钥对。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) ImportPrivateKey(request *model.ImportPrivateKeyRequest) (*model.ImportPrivateKeyResponse, error) {
	requestDef := GenReqDefForImportPrivateKey()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ImportPrivateKeyResponse), nil
	}
}

// ImportPrivateKeyInvoker 导入私钥
func (c *KpsClient) ImportPrivateKeyInvoker(request *model.ImportPrivateKeyRequest) *ImportPrivateKeyInvoker {
	requestDef := GenReqDefForImportPrivateKey()
	return &ImportPrivateKeyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListFailedTask 查询失败的任务信息
//
// 查询绑定、解绑等操作失败的任务信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) ListFailedTask(request *model.ListFailedTaskRequest) (*model.ListFailedTaskResponse, error) {
	requestDef := GenReqDefForListFailedTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListFailedTaskResponse), nil
	}
}

// ListFailedTaskInvoker 查询失败的任务信息
func (c *KpsClient) ListFailedTaskInvoker(request *model.ListFailedTaskRequest) *ListFailedTaskInvoker {
	requestDef := GenReqDefForListFailedTask()
	return &ListFailedTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListKeypairDetail 查询SSH密钥对详细信息
//
// 查询SSH密钥对详细信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) ListKeypairDetail(request *model.ListKeypairDetailRequest) (*model.ListKeypairDetailResponse, error) {
	requestDef := GenReqDefForListKeypairDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListKeypairDetailResponse), nil
	}
}

// ListKeypairDetailInvoker 查询SSH密钥对详细信息
func (c *KpsClient) ListKeypairDetailInvoker(request *model.ListKeypairDetailRequest) *ListKeypairDetailInvoker {
	requestDef := GenReqDefForListKeypairDetail()
	return &ListKeypairDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListKeypairTask 查询任务信息
//
// 根据SSH密钥对接口返回的task_id，查询SSH密钥对当前任务的执行状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) ListKeypairTask(request *model.ListKeypairTaskRequest) (*model.ListKeypairTaskResponse, error) {
	requestDef := GenReqDefForListKeypairTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListKeypairTaskResponse), nil
	}
}

// ListKeypairTaskInvoker 查询任务信息
func (c *KpsClient) ListKeypairTaskInvoker(request *model.ListKeypairTaskRequest) *ListKeypairTaskInvoker {
	requestDef := GenReqDefForListKeypairTask()
	return &ListKeypairTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListKeypairs 查询SSH密钥对列表
//
// 查询SSH密钥对列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) ListKeypairs(request *model.ListKeypairsRequest) (*model.ListKeypairsResponse, error) {
	requestDef := GenReqDefForListKeypairs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListKeypairsResponse), nil
	}
}

// ListKeypairsInvoker 查询SSH密钥对列表
func (c *KpsClient) ListKeypairsInvoker(request *model.ListKeypairsRequest) *ListKeypairsInvoker {
	requestDef := GenReqDefForListKeypairs()
	return &ListKeypairsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRunningTask 查询正在处理的任务信息
//
// 查询正在处理的任务信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) ListRunningTask(request *model.ListRunningTaskRequest) (*model.ListRunningTaskResponse, error) {
	requestDef := GenReqDefForListRunningTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRunningTaskResponse), nil
	}
}

// ListRunningTaskInvoker 查询正在处理的任务信息
func (c *KpsClient) ListRunningTaskInvoker(request *model.ListRunningTaskRequest) *ListRunningTaskInvoker {
	requestDef := GenReqDefForListRunningTask()
	return &ListRunningTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateKeypairDescription 更新SSH密钥对描述
//
// 更新SSH密钥对描述。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *KpsClient) UpdateKeypairDescription(request *model.UpdateKeypairDescriptionRequest) (*model.UpdateKeypairDescriptionResponse, error) {
	requestDef := GenReqDefForUpdateKeypairDescription()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateKeypairDescriptionResponse), nil
	}
}

// UpdateKeypairDescriptionInvoker 更新SSH密钥对描述
func (c *KpsClient) UpdateKeypairDescriptionInvoker(request *model.UpdateKeypairDescriptionRequest) *UpdateKeypairDescriptionInvoker {
	requestDef := GenReqDefForUpdateKeypairDescription()
	return &UpdateKeypairDescriptionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
