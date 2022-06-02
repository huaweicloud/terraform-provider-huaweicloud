package v3

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kps/v3/model"
)

type KpsClient struct {
	HcClient *http_client.HcHttpClient
}

func NewKpsClient(hcClient *http_client.HcHttpClient) *KpsClient {
	return &KpsClient{HcClient: hcClient}
}

func KpsClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// AssociateKeypair 绑定SSH密钥对
//
// 给指定的虚拟机绑定（替换或重置，替换需提供虚拟机已配置的SSH密钥对私钥；重置不需要提供虚拟机的SSH密钥对私钥）新的SSH密钥对。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// CreateKeypair 创建和导入SSH密钥对
//
// 创建和导入SSH密钥对
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ListFailedTask 查询失败的任务信息
//
// 查询绑定、解绑等操作失败的任务信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
