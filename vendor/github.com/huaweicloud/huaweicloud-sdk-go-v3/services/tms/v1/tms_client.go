package v1

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1/model"
)

type TmsClient struct {
	HcClient *http_client.HcHttpClient
}

func NewTmsClient(hcClient *http_client.HcHttpClient) *TmsClient {
	return &TmsClient{HcClient: hcClient}
}

func TmsClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder().WithCredentialsType("global.Credentials")
	return builder
}

// CreatePredefineTags 创建预定义标签
//
// 用于创建预定标签。用户创建预定义标签后，可以使用预定义标签来给资源创建标签。该接口支持幂等特性和处理批量数据。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *TmsClient) CreatePredefineTags(request *model.CreatePredefineTagsRequest) (*model.CreatePredefineTagsResponse, error) {
	requestDef := GenReqDefForCreatePredefineTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePredefineTagsResponse), nil
	}
}

// CreatePredefineTagsInvoker 创建预定义标签
func (c *TmsClient) CreatePredefineTagsInvoker(request *model.CreatePredefineTagsRequest) *CreatePredefineTagsInvoker {
	requestDef := GenReqDefForCreatePredefineTags()
	return &CreatePredefineTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeletePredefineTags 删除预定义标签
//
// 用于删除预定标签。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *TmsClient) DeletePredefineTags(request *model.DeletePredefineTagsRequest) (*model.DeletePredefineTagsResponse, error) {
	requestDef := GenReqDefForDeletePredefineTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeletePredefineTagsResponse), nil
	}
}

// DeletePredefineTagsInvoker 删除预定义标签
func (c *TmsClient) DeletePredefineTagsInvoker(request *model.DeletePredefineTagsRequest) *DeletePredefineTagsInvoker {
	requestDef := GenReqDefForDeletePredefineTags()
	return &DeletePredefineTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListApiVersions 查询API版本列表
//
// 查询标签管理服务的API版本列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *TmsClient) ListApiVersions(request *model.ListApiVersionsRequest) (*model.ListApiVersionsResponse, error) {
	requestDef := GenReqDefForListApiVersions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListApiVersionsResponse), nil
	}
}

// ListApiVersionsInvoker 查询API版本列表
func (c *TmsClient) ListApiVersionsInvoker(request *model.ListApiVersionsRequest) *ListApiVersionsInvoker {
	requestDef := GenReqDefForListApiVersions()
	return &ListApiVersionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPredefineTags 查询预定义标签列表
//
// 用于查询预定义标签列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *TmsClient) ListPredefineTags(request *model.ListPredefineTagsRequest) (*model.ListPredefineTagsResponse, error) {
	requestDef := GenReqDefForListPredefineTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPredefineTagsResponse), nil
	}
}

// ListPredefineTagsInvoker 查询预定义标签列表
func (c *TmsClient) ListPredefineTagsInvoker(request *model.ListPredefineTagsRequest) *ListPredefineTagsInvoker {
	requestDef := GenReqDefForListPredefineTags()
	return &ListPredefineTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowApiVersion 查询API版本号详情
//
// 查询指定的标签管理服务API版本号详情。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *TmsClient) ShowApiVersion(request *model.ShowApiVersionRequest) (*model.ShowApiVersionResponse, error) {
	requestDef := GenReqDefForShowApiVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowApiVersionResponse), nil
	}
}

// ShowApiVersionInvoker 查询API版本号详情
func (c *TmsClient) ShowApiVersionInvoker(request *model.ShowApiVersionRequest) *ShowApiVersionInvoker {
	requestDef := GenReqDefForShowApiVersion()
	return &ShowApiVersionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTagQuota 查询标签配额
//
// 查询标签的配额信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *TmsClient) ShowTagQuota(request *model.ShowTagQuotaRequest) (*model.ShowTagQuotaResponse, error) {
	requestDef := GenReqDefForShowTagQuota()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTagQuotaResponse), nil
	}
}

// ShowTagQuotaInvoker 查询标签配额
func (c *TmsClient) ShowTagQuotaInvoker(request *model.ShowTagQuotaRequest) *ShowTagQuotaInvoker {
	requestDef := GenReqDefForShowTagQuota()
	return &ShowTagQuotaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePredefineTags 修改预定义标签
//
// 修改预定义标签。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *TmsClient) UpdatePredefineTags(request *model.UpdatePredefineTagsRequest) (*model.UpdatePredefineTagsResponse, error) {
	requestDef := GenReqDefForUpdatePredefineTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePredefineTagsResponse), nil
	}
}

// UpdatePredefineTagsInvoker 修改预定义标签
func (c *TmsClient) UpdatePredefineTagsInvoker(request *model.UpdatePredefineTagsRequest) *UpdatePredefineTagsInvoker {
	requestDef := GenReqDefForUpdatePredefineTags()
	return &UpdatePredefineTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
