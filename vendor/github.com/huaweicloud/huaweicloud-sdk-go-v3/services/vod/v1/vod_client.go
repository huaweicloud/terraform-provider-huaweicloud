package v1

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"
)

type VodClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewVodClient(hcClient *httpclient.HcHttpClient) *VodClient {
	return &VodClient{HcClient: hcClient}
}

func VodClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder()
	return builder
}

// CancelAssetTranscodeTask 取消媒资转码任务
//
// 取消媒资转码任务，只能取消排队中的转码任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CancelAssetTranscodeTask(request *model.CancelAssetTranscodeTaskRequest) (*model.CancelAssetTranscodeTaskResponse, error) {
	requestDef := GenReqDefForCancelAssetTranscodeTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CancelAssetTranscodeTaskResponse), nil
	}
}

// CancelAssetTranscodeTaskInvoker 取消媒资转码任务
func (c *VodClient) CancelAssetTranscodeTaskInvoker(request *model.CancelAssetTranscodeTaskRequest) *CancelAssetTranscodeTaskInvoker {
	requestDef := GenReqDefForCancelAssetTranscodeTask()
	return &CancelAssetTranscodeTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CancelExtractAudioTask 取消提取音频任务
//
// 取消提取音频任务，只有排队中的提取音频任务才可以取消。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CancelExtractAudioTask(request *model.CancelExtractAudioTaskRequest) (*model.CancelExtractAudioTaskResponse, error) {
	requestDef := GenReqDefForCancelExtractAudioTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CancelExtractAudioTaskResponse), nil
	}
}

// CancelExtractAudioTaskInvoker 取消提取音频任务
func (c *VodClient) CancelExtractAudioTaskInvoker(request *model.CancelExtractAudioTaskRequest) *CancelExtractAudioTaskInvoker {
	requestDef := GenReqDefForCancelExtractAudioTask()
	return &CancelExtractAudioTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CheckMd5Duplication 上传检验
//
// 校验媒资文件是否已存储于视频点播服务中。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CheckMd5Duplication(request *model.CheckMd5DuplicationRequest) (*model.CheckMd5DuplicationResponse, error) {
	requestDef := GenReqDefForCheckMd5Duplication()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CheckMd5DuplicationResponse), nil
	}
}

// CheckMd5DuplicationInvoker 上传检验
func (c *VodClient) CheckMd5DuplicationInvoker(request *model.CheckMd5DuplicationRequest) *CheckMd5DuplicationInvoker {
	requestDef := GenReqDefForCheckMd5Duplication()
	return &CheckMd5DuplicationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ConfirmAssetUpload 确认媒资上传
//
// 媒资分段上传完成后，需要调用此接口通知点播服务媒资上传的状态，表示媒资上传创建完成。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ConfirmAssetUpload(request *model.ConfirmAssetUploadRequest) (*model.ConfirmAssetUploadResponse, error) {
	requestDef := GenReqDefForConfirmAssetUpload()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ConfirmAssetUploadResponse), nil
	}
}

// ConfirmAssetUploadInvoker 确认媒资上传
func (c *VodClient) ConfirmAssetUploadInvoker(request *model.ConfirmAssetUploadRequest) *ConfirmAssetUploadInvoker {
	requestDef := GenReqDefForConfirmAssetUpload()
	return &ConfirmAssetUploadInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ConfirmImageUpload 确认水印图片上传
//
// 确认水印图片上传状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ConfirmImageUpload(request *model.ConfirmImageUploadRequest) (*model.ConfirmImageUploadResponse, error) {
	requestDef := GenReqDefForConfirmImageUpload()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ConfirmImageUploadResponse), nil
	}
}

// ConfirmImageUploadInvoker 确认水印图片上传
func (c *VodClient) ConfirmImageUploadInvoker(request *model.ConfirmImageUploadRequest) *ConfirmImageUploadInvoker {
	requestDef := GenReqDefForConfirmImageUpload()
	return &ConfirmImageUploadInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAssetByFileUpload 创建媒资：上传方式
//
// 调用该接口创建媒资时，需要将对应的媒资文件上传到点播服务的OBS桶中。
//
// 若上传的单媒资文件大小小于20M，则可以直接用PUT方法对该接口返回的地址进行上传。具体使用方法请参考[示例1：媒资上传（20M以下）](https://support.huaweicloud.com/api-vod/vod_04_0195.html)。
//
// 若上传的单个媒资大小大于20M，则需要进行二进制流分割后上传，该接口的具体使用方法请参考[示例2：媒资分段上传（20M以上）](https://support.huaweicloud.com/api-vod/vod_04_0216.html)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateAssetByFileUpload(request *model.CreateAssetByFileUploadRequest) (*model.CreateAssetByFileUploadResponse, error) {
	requestDef := GenReqDefForCreateAssetByFileUpload()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAssetByFileUploadResponse), nil
	}
}

// CreateAssetByFileUploadInvoker 创建媒资：上传方式
func (c *VodClient) CreateAssetByFileUploadInvoker(request *model.CreateAssetByFileUploadRequest) *CreateAssetByFileUploadInvoker {
	requestDef := GenReqDefForCreateAssetByFileUpload()
	return &CreateAssetByFileUploadInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAssetCategory 创建媒资分类
//
// 创建媒资分类。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateAssetCategory(request *model.CreateAssetCategoryRequest) (*model.CreateAssetCategoryResponse, error) {
	requestDef := GenReqDefForCreateAssetCategory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAssetCategoryResponse), nil
	}
}

// CreateAssetCategoryInvoker 创建媒资分类
func (c *VodClient) CreateAssetCategoryInvoker(request *model.CreateAssetCategoryRequest) *CreateAssetCategoryInvoker {
	requestDef := GenReqDefForCreateAssetCategory()
	return &CreateAssetCategoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAssetProcessTask 媒资处理
//
// 实现视频转码、截图、加密等处理。既可以同时启动多种操作，也可以只启动一种操作。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateAssetProcessTask(request *model.CreateAssetProcessTaskRequest) (*model.CreateAssetProcessTaskResponse, error) {
	requestDef := GenReqDefForCreateAssetProcessTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAssetProcessTaskResponse), nil
	}
}

// CreateAssetProcessTaskInvoker 媒资处理
func (c *VodClient) CreateAssetProcessTaskInvoker(request *model.CreateAssetProcessTaskRequest) *CreateAssetProcessTaskInvoker {
	requestDef := GenReqDefForCreateAssetProcessTask()
	return &CreateAssetProcessTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAssetReviewTask 创建审核媒资任务
//
// 对上传的媒资进行审核。审核后，可以调用[查询媒资详细信息](https://support.huaweicloud.com/api-vod/vod_04_0202.html)接口查看审核结果。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateAssetReviewTask(request *model.CreateAssetReviewTaskRequest) (*model.CreateAssetReviewTaskResponse, error) {
	requestDef := GenReqDefForCreateAssetReviewTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAssetReviewTaskResponse), nil
	}
}

// CreateAssetReviewTaskInvoker 创建审核媒资任务
func (c *VodClient) CreateAssetReviewTaskInvoker(request *model.CreateAssetReviewTaskRequest) *CreateAssetReviewTaskInvoker {
	requestDef := GenReqDefForCreateAssetReviewTask()
	return &CreateAssetReviewTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateExtractAudioTask 音频提取
//
// 本接口为异步接口，创建音频提取任务下发成功后会返回asset_id和提取的audio_asset_id，但此时音频提取任务并没有立即完成，可通过消息订阅界面配置的音频提取完成事件来获取音频提取任务完成与否。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateExtractAudioTask(request *model.CreateExtractAudioTaskRequest) (*model.CreateExtractAudioTaskResponse, error) {
	requestDef := GenReqDefForCreateExtractAudioTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateExtractAudioTaskResponse), nil
	}
}

// CreateExtractAudioTaskInvoker 音频提取
func (c *VodClient) CreateExtractAudioTaskInvoker(request *model.CreateExtractAudioTaskRequest) *CreateExtractAudioTaskInvoker {
	requestDef := GenReqDefForCreateExtractAudioTask()
	return &CreateExtractAudioTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePreheatingAsset CDN预热
//
// 媒资发布后，可通过指定媒资ID或URL向CDN预热。用户初次请求时，将由CDN节点提供请求媒资，加快用户下载缓存时间，提高用户体验。单租户每天最多预热1000个。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreatePreheatingAsset(request *model.CreatePreheatingAssetRequest) (*model.CreatePreheatingAssetResponse, error) {
	requestDef := GenReqDefForCreatePreheatingAsset()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePreheatingAssetResponse), nil
	}
}

// CreatePreheatingAssetInvoker CDN预热
func (c *VodClient) CreatePreheatingAssetInvoker(request *model.CreatePreheatingAssetRequest) *CreatePreheatingAssetInvoker {
	requestDef := GenReqDefForCreatePreheatingAsset()
	return &CreatePreheatingAssetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTakeOverTask 创建媒资：OBS托管方式
//
// 通过存量托管的方式，将已存储在OBS桶中的音视频文件同步到点播服务。
//
// OBS托管方式分为增量托管和存量托管，增量托管暂只支持通过视频点播控制台配置，配置后，若OBS有新增音视频文件，则会自动同步到点播服务中，具体请参见[增量托管](https://support.huaweicloud.com/usermanual-vod/vod010032.html)。两个托管方式都需要先将对应的OBS桶授权给点播服务，具体请参见[桶授权](https://support.huaweicloud.com/usermanual-vod/vod010031.html)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateTakeOverTask(request *model.CreateTakeOverTaskRequest) (*model.CreateTakeOverTaskResponse, error) {
	requestDef := GenReqDefForCreateTakeOverTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTakeOverTaskResponse), nil
	}
}

// CreateTakeOverTaskInvoker 创建媒资：OBS托管方式
func (c *VodClient) CreateTakeOverTaskInvoker(request *model.CreateTakeOverTaskRequest) *CreateTakeOverTaskInvoker {
	requestDef := GenReqDefForCreateTakeOverTask()
	return &CreateTakeOverTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTemplateGroup 创建自定义转码模板组
//
// 创建自定义转码模板组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateTemplateGroup(request *model.CreateTemplateGroupRequest) (*model.CreateTemplateGroupResponse, error) {
	requestDef := GenReqDefForCreateTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTemplateGroupResponse), nil
	}
}

// CreateTemplateGroupInvoker 创建自定义转码模板组
func (c *VodClient) CreateTemplateGroupInvoker(request *model.CreateTemplateGroupRequest) *CreateTemplateGroupInvoker {
	requestDef := GenReqDefForCreateTemplateGroup()
	return &CreateTemplateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTemplateGroupCollection 创建转码模板组集合
//
// 创建转码模板组集合
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateTemplateGroupCollection(request *model.CreateTemplateGroupCollectionRequest) (*model.CreateTemplateGroupCollectionResponse, error) {
	requestDef := GenReqDefForCreateTemplateGroupCollection()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTemplateGroupCollectionResponse), nil
	}
}

// CreateTemplateGroupCollectionInvoker 创建转码模板组集合
func (c *VodClient) CreateTemplateGroupCollectionInvoker(request *model.CreateTemplateGroupCollectionRequest) *CreateTemplateGroupCollectionInvoker {
	requestDef := GenReqDefForCreateTemplateGroupCollection()
	return &CreateTemplateGroupCollectionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTranscodeTemplate 创建自定义转码模板
//
// 创建自定义转码模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateTranscodeTemplate(request *model.CreateTranscodeTemplateRequest) (*model.CreateTranscodeTemplateResponse, error) {
	requestDef := GenReqDefForCreateTranscodeTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTranscodeTemplateResponse), nil
	}
}

// CreateTranscodeTemplateInvoker 创建自定义转码模板
func (c *VodClient) CreateTranscodeTemplateInvoker(request *model.CreateTranscodeTemplateRequest) *CreateTranscodeTemplateInvoker {
	requestDef := GenReqDefForCreateTranscodeTemplate()
	return &CreateTranscodeTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateWatermarkTemplate 创建水印模板
//
// 创建水印模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) CreateWatermarkTemplate(request *model.CreateWatermarkTemplateRequest) (*model.CreateWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForCreateWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateWatermarkTemplateResponse), nil
	}
}

// CreateWatermarkTemplateInvoker 创建水印模板
func (c *VodClient) CreateWatermarkTemplateInvoker(request *model.CreateWatermarkTemplateRequest) *CreateWatermarkTemplateInvoker {
	requestDef := GenReqDefForCreateWatermarkTemplate()
	return &CreateWatermarkTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAssetCategory 删除媒资分类
//
// 删除媒资分类。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) DeleteAssetCategory(request *model.DeleteAssetCategoryRequest) (*model.DeleteAssetCategoryResponse, error) {
	requestDef := GenReqDefForDeleteAssetCategory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAssetCategoryResponse), nil
	}
}

// DeleteAssetCategoryInvoker 删除媒资分类
func (c *VodClient) DeleteAssetCategoryInvoker(request *model.DeleteAssetCategoryRequest) *DeleteAssetCategoryInvoker {
	requestDef := GenReqDefForDeleteAssetCategory()
	return &DeleteAssetCategoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAssets 删除媒资
//
// 删除媒资。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) DeleteAssets(request *model.DeleteAssetsRequest) (*model.DeleteAssetsResponse, error) {
	requestDef := GenReqDefForDeleteAssets()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAssetsResponse), nil
	}
}

// DeleteAssetsInvoker 删除媒资
func (c *VodClient) DeleteAssetsInvoker(request *model.DeleteAssetsRequest) *DeleteAssetsInvoker {
	requestDef := GenReqDefForDeleteAssets()
	return &DeleteAssetsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTemplateGroup 删除自定义转码模板组
//
// 删除自定义转码模板组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) DeleteTemplateGroup(request *model.DeleteTemplateGroupRequest) (*model.DeleteTemplateGroupResponse, error) {
	requestDef := GenReqDefForDeleteTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTemplateGroupResponse), nil
	}
}

// DeleteTemplateGroupInvoker 删除自定义转码模板组
func (c *VodClient) DeleteTemplateGroupInvoker(request *model.DeleteTemplateGroupRequest) *DeleteTemplateGroupInvoker {
	requestDef := GenReqDefForDeleteTemplateGroup()
	return &DeleteTemplateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTemplateGroupCollection 删除转码模板组集合
//
// 删除转码模板组集合
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) DeleteTemplateGroupCollection(request *model.DeleteTemplateGroupCollectionRequest) (*model.DeleteTemplateGroupCollectionResponse, error) {
	requestDef := GenReqDefForDeleteTemplateGroupCollection()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTemplateGroupCollectionResponse), nil
	}
}

// DeleteTemplateGroupCollectionInvoker 删除转码模板组集合
func (c *VodClient) DeleteTemplateGroupCollectionInvoker(request *model.DeleteTemplateGroupCollectionRequest) *DeleteTemplateGroupCollectionInvoker {
	requestDef := GenReqDefForDeleteTemplateGroupCollection()
	return &DeleteTemplateGroupCollectionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTranscodeProduct 删除转码产物
//
// 删除转码产物。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) DeleteTranscodeProduct(request *model.DeleteTranscodeProductRequest) (*model.DeleteTranscodeProductResponse, error) {
	requestDef := GenReqDefForDeleteTranscodeProduct()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTranscodeProductResponse), nil
	}
}

// DeleteTranscodeProductInvoker 删除转码产物
func (c *VodClient) DeleteTranscodeProductInvoker(request *model.DeleteTranscodeProductRequest) *DeleteTranscodeProductInvoker {
	requestDef := GenReqDefForDeleteTranscodeProduct()
	return &DeleteTranscodeProductInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTranscodeTemplate 删除自定义模板
//
// 删除自定义模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) DeleteTranscodeTemplate(request *model.DeleteTranscodeTemplateRequest) (*model.DeleteTranscodeTemplateResponse, error) {
	requestDef := GenReqDefForDeleteTranscodeTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTranscodeTemplateResponse), nil
	}
}

// DeleteTranscodeTemplateInvoker 删除自定义模板
func (c *VodClient) DeleteTranscodeTemplateInvoker(request *model.DeleteTranscodeTemplateRequest) *DeleteTranscodeTemplateInvoker {
	requestDef := GenReqDefForDeleteTranscodeTemplate()
	return &DeleteTranscodeTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteWatermarkTemplate 删除水印模板
//
// 删除水印模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) DeleteWatermarkTemplate(request *model.DeleteWatermarkTemplateRequest) (*model.DeleteWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForDeleteWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteWatermarkTemplateResponse), nil
	}
}

// DeleteWatermarkTemplateInvoker 删除水印模板
func (c *VodClient) DeleteWatermarkTemplateInvoker(request *model.DeleteWatermarkTemplateRequest) *DeleteWatermarkTemplateInvoker {
	requestDef := GenReqDefForDeleteWatermarkTemplate()
	return &DeleteWatermarkTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAssetCategory 查询指定分类信息
//
// 查询指定分类信息，及其子分类（即下一级分类）的列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListAssetCategory(request *model.ListAssetCategoryRequest) (*model.ListAssetCategoryResponse, error) {
	requestDef := GenReqDefForListAssetCategory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAssetCategoryResponse), nil
	}
}

// ListAssetCategoryInvoker 查询指定分类信息
func (c *VodClient) ListAssetCategoryInvoker(request *model.ListAssetCategoryRequest) *ListAssetCategoryInvoker {
	requestDef := GenReqDefForListAssetCategory()
	return &ListAssetCategoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAssetDailySummaryLog 查询媒资日播放统计数据
//
// 查询媒资日播放统计数据。
//
// 使用媒资日播放统计查询API前，需要先提交工单开通统计功能，才能触发统计任务。
//
// 支持查询最近一年的播放统计数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListAssetDailySummaryLog(request *model.ListAssetDailySummaryLogRequest) (*model.ListAssetDailySummaryLogResponse, error) {
	requestDef := GenReqDefForListAssetDailySummaryLog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAssetDailySummaryLogResponse), nil
	}
}

// ListAssetDailySummaryLogInvoker 查询媒资日播放统计数据
func (c *VodClient) ListAssetDailySummaryLogInvoker(request *model.ListAssetDailySummaryLogRequest) *ListAssetDailySummaryLogInvoker {
	requestDef := GenReqDefForListAssetDailySummaryLog()
	return &ListAssetDailySummaryLogInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAssetList 查询媒资列表
//
// 查询媒资列表，列表中的每一条记录包含媒资的概要信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListAssetList(request *model.ListAssetListRequest) (*model.ListAssetListResponse, error) {
	requestDef := GenReqDefForListAssetList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAssetListResponse), nil
	}
}

// ListAssetListInvoker 查询媒资列表
func (c *VodClient) ListAssetListInvoker(request *model.ListAssetListRequest) *ListAssetListInvoker {
	requestDef := GenReqDefForListAssetList()
	return &ListAssetListInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDomainLogs 查询域名播放日志
//
// 查询指定点播域名某段时间内在CDN的相关日志。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListDomainLogs(request *model.ListDomainLogsRequest) (*model.ListDomainLogsResponse, error) {
	requestDef := GenReqDefForListDomainLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDomainLogsResponse), nil
	}
}

// ListDomainLogsInvoker 查询域名播放日志
func (c *VodClient) ListDomainLogsInvoker(request *model.ListDomainLogsRequest) *ListDomainLogsInvoker {
	requestDef := GenReqDefForListDomainLogs()
	return &ListDomainLogsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTemplateGroup 查询转码模板组列表
//
// 查询转码模板组列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListTemplateGroup(request *model.ListTemplateGroupRequest) (*model.ListTemplateGroupResponse, error) {
	requestDef := GenReqDefForListTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTemplateGroupResponse), nil
	}
}

// ListTemplateGroupInvoker 查询转码模板组列表
func (c *VodClient) ListTemplateGroupInvoker(request *model.ListTemplateGroupRequest) *ListTemplateGroupInvoker {
	requestDef := GenReqDefForListTemplateGroup()
	return &ListTemplateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTemplateGroupCollection 查询自定义模板组集合
//
// 查询转码模板组集合
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListTemplateGroupCollection(request *model.ListTemplateGroupCollectionRequest) (*model.ListTemplateGroupCollectionResponse, error) {
	requestDef := GenReqDefForListTemplateGroupCollection()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTemplateGroupCollectionResponse), nil
	}
}

// ListTemplateGroupCollectionInvoker 查询自定义模板组集合
func (c *VodClient) ListTemplateGroupCollectionInvoker(request *model.ListTemplateGroupCollectionRequest) *ListTemplateGroupCollectionInvoker {
	requestDef := GenReqDefForListTemplateGroupCollection()
	return &ListTemplateGroupCollectionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTopStatistics 查询TopN媒资信息
//
// 查询指定域名在指定日期播放次数排名Top 100的媒资统计数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListTopStatistics(request *model.ListTopStatisticsRequest) (*model.ListTopStatisticsResponse, error) {
	requestDef := GenReqDefForListTopStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTopStatisticsResponse), nil
	}
}

// ListTopStatisticsInvoker 查询TopN媒资信息
func (c *VodClient) ListTopStatisticsInvoker(request *model.ListTopStatisticsRequest) *ListTopStatisticsInvoker {
	requestDef := GenReqDefForListTopStatistics()
	return &ListTopStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTranscodeTemplate 查询转码模板列表
//
// 查询转码模板列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListTranscodeTemplate(request *model.ListTranscodeTemplateRequest) (*model.ListTranscodeTemplateResponse, error) {
	requestDef := GenReqDefForListTranscodeTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTranscodeTemplateResponse), nil
	}
}

// ListTranscodeTemplateInvoker 查询转码模板列表
func (c *VodClient) ListTranscodeTemplateInvoker(request *model.ListTranscodeTemplateRequest) *ListTranscodeTemplateInvoker {
	requestDef := GenReqDefForListTranscodeTemplate()
	return &ListTranscodeTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListWatermarkTemplate 查询水印列表
//
// 查询水印模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListWatermarkTemplate(request *model.ListWatermarkTemplateRequest) (*model.ListWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForListWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListWatermarkTemplateResponse), nil
	}
}

// ListWatermarkTemplateInvoker 查询水印列表
func (c *VodClient) ListWatermarkTemplateInvoker(request *model.ListWatermarkTemplateRequest) *ListWatermarkTemplateInvoker {
	requestDef := GenReqDefForListWatermarkTemplate()
	return &ListWatermarkTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ModifySubtitle 多字幕封装
//
// 多字幕封装，仅支持 HLS VTT格式
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ModifySubtitle(request *model.ModifySubtitleRequest) (*model.ModifySubtitleResponse, error) {
	requestDef := GenReqDefForModifySubtitle()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ModifySubtitleResponse), nil
	}
}

// ModifySubtitleInvoker 多字幕封装
func (c *VodClient) ModifySubtitleInvoker(request *model.ModifySubtitleRequest) *ModifySubtitleInvoker {
	requestDef := GenReqDefForModifySubtitle()
	return &ModifySubtitleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// PublishAssetFromObs 创建媒资：OBS转存方式
//
// 若您在使用点播服务前，已经在OBS桶中存储了音视频文件，您可以使用该接口将存储在OBS桶中的音视频文件转存到点播服务中，使用点播服务的音视频管理功能。调用该接口前，您需要调用[桶授权](https://support.huaweicloud.com/api-vod/vod_04_0199.html)接口，将存储音视频文件的OBS桶授权给点播服务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) PublishAssetFromObs(request *model.PublishAssetFromObsRequest) (*model.PublishAssetFromObsResponse, error) {
	requestDef := GenReqDefForPublishAssetFromObs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.PublishAssetFromObsResponse), nil
	}
}

// PublishAssetFromObsInvoker 创建媒资：OBS转存方式
func (c *VodClient) PublishAssetFromObsInvoker(request *model.PublishAssetFromObsRequest) *PublishAssetFromObsInvoker {
	requestDef := GenReqDefForPublishAssetFromObs()
	return &PublishAssetFromObsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// PublishAssets 媒资发布
//
// 将媒资设置为发布状态。支持批量发布。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) PublishAssets(request *model.PublishAssetsRequest) (*model.PublishAssetsResponse, error) {
	requestDef := GenReqDefForPublishAssets()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.PublishAssetsResponse), nil
	}
}

// PublishAssetsInvoker 媒资发布
func (c *VodClient) PublishAssetsInvoker(request *model.PublishAssetsRequest) *PublishAssetsInvoker {
	requestDef := GenReqDefForPublishAssets()
	return &PublishAssetsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAssetCipher 密钥查询
//
// 终端播放HLS加密视频时，向租户管理系统请求密钥，租户管理系统先查询其本地有没有已缓存的密钥，没有时则调用此接口向VOD查询。该接口的具体使用场景请参见[通过HLS加密防止视频泄露](https://support.huaweicloud.com/bestpractice-vod/vod_10_0004.html)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowAssetCipher(request *model.ShowAssetCipherRequest) (*model.ShowAssetCipherResponse, error) {
	requestDef := GenReqDefForShowAssetCipher()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetCipherResponse), nil
	}
}

// ShowAssetCipherInvoker 密钥查询
func (c *VodClient) ShowAssetCipherInvoker(request *model.ShowAssetCipherRequest) *ShowAssetCipherInvoker {
	requestDef := GenReqDefForShowAssetCipher()
	return &ShowAssetCipherInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAssetDetail 查询指定媒资的详细信息
//
// 查询指定媒资的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowAssetDetail(request *model.ShowAssetDetailRequest) (*model.ShowAssetDetailResponse, error) {
	requestDef := GenReqDefForShowAssetDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetDetailResponse), nil
	}
}

// ShowAssetDetailInvoker 查询指定媒资的详细信息
func (c *VodClient) ShowAssetDetailInvoker(request *model.ShowAssetDetailRequest) *ShowAssetDetailInvoker {
	requestDef := GenReqDefForShowAssetDetail()
	return &ShowAssetDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAssetMeta 查询媒资信息
//
// 查询媒资信息，支持指定媒资ID、分类、状态、起止时间查询。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowAssetMeta(request *model.ShowAssetMetaRequest) (*model.ShowAssetMetaResponse, error) {
	requestDef := GenReqDefForShowAssetMeta()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetMetaResponse), nil
	}
}

// ShowAssetMetaInvoker 查询媒资信息
func (c *VodClient) ShowAssetMetaInvoker(request *model.ShowAssetMetaRequest) *ShowAssetMetaInvoker {
	requestDef := GenReqDefForShowAssetMeta()
	return &ShowAssetMetaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAssetTempAuthority 获取分段上传授权
//
// 客户端请求创建媒资时，如果媒资文件超过20MB，需采用分段的方式向OBS上传，在每次与OBS交互前，客户端需通过此接口获取到授权方可与OBS交互。
//
// 该接口可以获取[初始化多段上传任务](https://support.huaweicloud.com/api-obs/obs_04_0098.html)、[上传段](https://support.huaweicloud.com/api-obs/obs_04_0099.html)、[合并段](https://support.huaweicloud.com/api-obs/obs_04_0102.html)、[列举已上传段](https://support.huaweicloud.com/api-obs/obs_04_0101.html)、[取消段合并](https://support.huaweicloud.com/api-obs/obs_04_0103.html)的带有临时授权的URL，用户需要根据OBS的接口文档配置相应请求的HTTP请求方法、请求头、请求体，然后请求对应的带有临时授权的URL。
//
// 视频分段上传方式和OBS的接口文档保持一致，包括HTTP请求方法、请求头、请求体等各种入参，此接口的作用是为用户生成带有鉴权信息的URL（鉴权信息即query_str），用来替换OBS接口中对应的URL，临时给用户开通向点播服务的桶上传文件的权限。
//
// 调用获取授权接口时需要传入bucket、object_key、http_verb，其中bucket和object_key是由[创建媒资：上传方式](https://support.huaweicloud.com/api-vod/vod_04_0196.html)接口中返回的响应体中的target字段获得的bucket和object，http_verb需要根据指定的操作选择。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowAssetTempAuthority(request *model.ShowAssetTempAuthorityRequest) (*model.ShowAssetTempAuthorityResponse, error) {
	requestDef := GenReqDefForShowAssetTempAuthority()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetTempAuthorityResponse), nil
	}
}

// ShowAssetTempAuthorityInvoker 获取分段上传授权
func (c *VodClient) ShowAssetTempAuthorityInvoker(request *model.ShowAssetTempAuthorityRequest) *ShowAssetTempAuthorityInvoker {
	requestDef := GenReqDefForShowAssetTempAuthority()
	return &ShowAssetTempAuthorityInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowCdnStatistics 查询CDN统计信息
//
// 查询CDN的统计数据，包括流量、峰值带宽、请求总数、请求命中率、流量命中率。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowCdnStatistics(request *model.ShowCdnStatisticsRequest) (*model.ShowCdnStatisticsResponse, error) {
	requestDef := GenReqDefForShowCdnStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCdnStatisticsResponse), nil
	}
}

// ShowCdnStatisticsInvoker 查询CDN统计信息
func (c *VodClient) ShowCdnStatisticsInvoker(request *model.ShowCdnStatisticsRequest) *ShowCdnStatisticsInvoker {
	requestDef := GenReqDefForShowCdnStatistics()
	return &ShowCdnStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowPreheatingAsset 查询CDN预热
//
// 查询预热结果。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowPreheatingAsset(request *model.ShowPreheatingAssetRequest) (*model.ShowPreheatingAssetResponse, error) {
	requestDef := GenReqDefForShowPreheatingAsset()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowPreheatingAssetResponse), nil
	}
}

// ShowPreheatingAssetInvoker 查询CDN预热
func (c *VodClient) ShowPreheatingAssetInvoker(request *model.ShowPreheatingAssetRequest) *ShowPreheatingAssetInvoker {
	requestDef := GenReqDefForShowPreheatingAsset()
	return &ShowPreheatingAssetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowVodRetrieval 查询取回数据信息
//
// ## 典型场景 ##
//
//	用于查询点播低频和归档取回量统计数据。&lt;br/&gt;
//
// ## 接口功能 ##
//
//	用于查询点播低频和归档取回量统计数据。&lt;br/&gt;
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowVodRetrieval(request *model.ShowVodRetrievalRequest) (*model.ShowVodRetrievalResponse, error) {
	requestDef := GenReqDefForShowVodRetrieval()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowVodRetrievalResponse), nil
	}
}

// ShowVodRetrievalInvoker 查询取回数据信息
func (c *VodClient) ShowVodRetrievalInvoker(request *model.ShowVodRetrievalRequest) *ShowVodRetrievalInvoker {
	requestDef := GenReqDefForShowVodRetrieval()
	return &ShowVodRetrievalInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowVodStatistics 查询源站统计信息
//
// 查询点播源站的统计数据，包括流量、存储空间、转码时长。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowVodStatistics(request *model.ShowVodStatisticsRequest) (*model.ShowVodStatisticsResponse, error) {
	requestDef := GenReqDefForShowVodStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowVodStatisticsResponse), nil
	}
}

// ShowVodStatisticsInvoker 查询源站统计信息
func (c *VodClient) ShowVodStatisticsInvoker(request *model.ShowVodStatisticsRequest) *ShowVodStatisticsInvoker {
	requestDef := GenReqDefForShowVodStatistics()
	return &ShowVodStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UnpublishAssets 媒资发布取消
//
// 将媒资设置为未发布状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UnpublishAssets(request *model.UnpublishAssetsRequest) (*model.UnpublishAssetsResponse, error) {
	requestDef := GenReqDefForUnpublishAssets()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UnpublishAssetsResponse), nil
	}
}

// UnpublishAssetsInvoker 媒资发布取消
func (c *VodClient) UnpublishAssetsInvoker(request *model.UnpublishAssetsRequest) *UnpublishAssetsInvoker {
	requestDef := GenReqDefForUnpublishAssets()
	return &UnpublishAssetsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAsset 视频更新
//
// 媒资创建后，单独上传封面、更新视频文件或更新已有封面。
//
// 如果是更新视频文件，更新完后要通过[确认媒资上传](https://support.huaweicloud.com/api-vod/vod_04_0198.html)接口通知点播服务。
//
// 如果是更新封面或单独上传封面，则不需通知。
//
// 更新视频可以使用分段上传，具体方式可以参考[示例2：媒资分段上传（20M以上）](https://support.huaweicloud.com/api-vod/vod_04_0216.html)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateAsset(request *model.UpdateAssetRequest) (*model.UpdateAssetResponse, error) {
	requestDef := GenReqDefForUpdateAsset()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAssetResponse), nil
	}
}

// UpdateAssetInvoker 视频更新
func (c *VodClient) UpdateAssetInvoker(request *model.UpdateAssetRequest) *UpdateAssetInvoker {
	requestDef := GenReqDefForUpdateAsset()
	return &UpdateAssetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAssetCategory 修改媒资分类
//
// 修改媒资分类。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateAssetCategory(request *model.UpdateAssetCategoryRequest) (*model.UpdateAssetCategoryResponse, error) {
	requestDef := GenReqDefForUpdateAssetCategory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAssetCategoryResponse), nil
	}
}

// UpdateAssetCategoryInvoker 修改媒资分类
func (c *VodClient) UpdateAssetCategoryInvoker(request *model.UpdateAssetCategoryRequest) *UpdateAssetCategoryInvoker {
	requestDef := GenReqDefForUpdateAssetCategory()
	return &UpdateAssetCategoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAssetMeta 修改媒资属性
//
// 修改媒资属性。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateAssetMeta(request *model.UpdateAssetMetaRequest) (*model.UpdateAssetMetaResponse, error) {
	requestDef := GenReqDefForUpdateAssetMeta()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAssetMetaResponse), nil
	}
}

// UpdateAssetMetaInvoker 修改媒资属性
func (c *VodClient) UpdateAssetMetaInvoker(request *model.UpdateAssetMetaRequest) *UpdateAssetMetaInvoker {
	requestDef := GenReqDefForUpdateAssetMeta()
	return &UpdateAssetMetaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateBucketAuthorized 桶授权
//
// 用户可以通过该接口将OBS桶授权给点播服务或取消点播服务的授权。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateBucketAuthorized(request *model.UpdateBucketAuthorizedRequest) (*model.UpdateBucketAuthorizedResponse, error) {
	requestDef := GenReqDefForUpdateBucketAuthorized()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateBucketAuthorizedResponse), nil
	}
}

// UpdateBucketAuthorizedInvoker 桶授权
func (c *VodClient) UpdateBucketAuthorizedInvoker(request *model.UpdateBucketAuthorizedRequest) *UpdateBucketAuthorizedInvoker {
	requestDef := GenReqDefForUpdateBucketAuthorized()
	return &UpdateBucketAuthorizedInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateCoverByThumbnail 设置封面
//
// 将视频截图生成的某张图片设置成封面。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateCoverByThumbnail(request *model.UpdateCoverByThumbnailRequest) (*model.UpdateCoverByThumbnailResponse, error) {
	requestDef := GenReqDefForUpdateCoverByThumbnail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCoverByThumbnailResponse), nil
	}
}

// UpdateCoverByThumbnailInvoker 设置封面
func (c *VodClient) UpdateCoverByThumbnailInvoker(request *model.UpdateCoverByThumbnailRequest) *UpdateCoverByThumbnailInvoker {
	requestDef := GenReqDefForUpdateCoverByThumbnail()
	return &UpdateCoverByThumbnailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateStorageMode 修改媒资文件在obs的存储模式
//
// ## 接口功能 ##
//
//	修改媒资文件在obs的存储模式&lt;br/&gt;
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateStorageMode(request *model.UpdateStorageModeRequest) (*model.UpdateStorageModeResponse, error) {
	requestDef := GenReqDefForUpdateStorageMode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateStorageModeResponse), nil
	}
}

// UpdateStorageModeInvoker 修改媒资文件在obs的存储模式
func (c *VodClient) UpdateStorageModeInvoker(request *model.UpdateStorageModeRequest) *UpdateStorageModeInvoker {
	requestDef := GenReqDefForUpdateStorageMode()
	return &UpdateStorageModeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTemplateGroup 修改自定义转码模板组
//
// 修改自定义转码模板组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateTemplateGroup(request *model.UpdateTemplateGroupRequest) (*model.UpdateTemplateGroupResponse, error) {
	requestDef := GenReqDefForUpdateTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTemplateGroupResponse), nil
	}
}

// UpdateTemplateGroupInvoker 修改自定义转码模板组
func (c *VodClient) UpdateTemplateGroupInvoker(request *model.UpdateTemplateGroupRequest) *UpdateTemplateGroupInvoker {
	requestDef := GenReqDefForUpdateTemplateGroup()
	return &UpdateTemplateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTemplateGroupCollection 修改转码模板组集合
//
// 修改转码模板组结合
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateTemplateGroupCollection(request *model.UpdateTemplateGroupCollectionRequest) (*model.UpdateTemplateGroupCollectionResponse, error) {
	requestDef := GenReqDefForUpdateTemplateGroupCollection()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTemplateGroupCollectionResponse), nil
	}
}

// UpdateTemplateGroupCollectionInvoker 修改转码模板组集合
func (c *VodClient) UpdateTemplateGroupCollectionInvoker(request *model.UpdateTemplateGroupCollectionRequest) *UpdateTemplateGroupCollectionInvoker {
	requestDef := GenReqDefForUpdateTemplateGroupCollection()
	return &UpdateTemplateGroupCollectionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTranscodeTemplate 修改转码模板
//
// 修改转码模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateTranscodeTemplate(request *model.UpdateTranscodeTemplateRequest) (*model.UpdateTranscodeTemplateResponse, error) {
	requestDef := GenReqDefForUpdateTranscodeTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTranscodeTemplateResponse), nil
	}
}

// UpdateTranscodeTemplateInvoker 修改转码模板
func (c *VodClient) UpdateTranscodeTemplateInvoker(request *model.UpdateTranscodeTemplateRequest) *UpdateTranscodeTemplateInvoker {
	requestDef := GenReqDefForUpdateTranscodeTemplate()
	return &UpdateTranscodeTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateWatermarkTemplate 修改水印模板
//
// 修改水印模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UpdateWatermarkTemplate(request *model.UpdateWatermarkTemplateRequest) (*model.UpdateWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForUpdateWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateWatermarkTemplateResponse), nil
	}
}

// UpdateWatermarkTemplateInvoker 修改水印模板
func (c *VodClient) UpdateWatermarkTemplateInvoker(request *model.UpdateWatermarkTemplateRequest) *UpdateWatermarkTemplateInvoker {
	requestDef := GenReqDefForUpdateWatermarkTemplate()
	return &UpdateWatermarkTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UploadMetaDataByUrl 创建媒资：URL拉取注入
//
// 基于音视频源文件URL，将音视频文件离线拉取上传到点播服务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) UploadMetaDataByUrl(request *model.UploadMetaDataByUrlRequest) (*model.UploadMetaDataByUrlResponse, error) {
	requestDef := GenReqDefForUploadMetaDataByUrl()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UploadMetaDataByUrlResponse), nil
	}
}

// UploadMetaDataByUrlInvoker 创建媒资：URL拉取注入
func (c *VodClient) UploadMetaDataByUrlInvoker(request *model.UploadMetaDataByUrlRequest) *UploadMetaDataByUrlInvoker {
	requestDef := GenReqDefForUploadMetaDataByUrl()
	return &UploadMetaDataByUrlInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTakeOverTask 查询托管任务
//
// 查询OBS存量托管任务列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ListTakeOverTask(request *model.ListTakeOverTaskRequest) (*model.ListTakeOverTaskResponse, error) {
	requestDef := GenReqDefForListTakeOverTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTakeOverTaskResponse), nil
	}
}

// ListTakeOverTaskInvoker 查询托管任务
func (c *VodClient) ListTakeOverTaskInvoker(request *model.ListTakeOverTaskRequest) *ListTakeOverTaskInvoker {
	requestDef := GenReqDefForListTakeOverTask()
	return &ListTakeOverTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTakeOverAssetDetails 查询托管媒资详情
//
// 查询OBS托管媒资的详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowTakeOverAssetDetails(request *model.ShowTakeOverAssetDetailsRequest) (*model.ShowTakeOverAssetDetailsResponse, error) {
	requestDef := GenReqDefForShowTakeOverAssetDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTakeOverAssetDetailsResponse), nil
	}
}

// ShowTakeOverAssetDetailsInvoker 查询托管媒资详情
func (c *VodClient) ShowTakeOverAssetDetailsInvoker(request *model.ShowTakeOverAssetDetailsRequest) *ShowTakeOverAssetDetailsInvoker {
	requestDef := GenReqDefForShowTakeOverAssetDetails()
	return &ShowTakeOverAssetDetailsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTakeOverTaskDetails 查询托管任务详情
//
// 查询OBS存量托管任务详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VodClient) ShowTakeOverTaskDetails(request *model.ShowTakeOverTaskDetailsRequest) (*model.ShowTakeOverTaskDetailsResponse, error) {
	requestDef := GenReqDefForShowTakeOverTaskDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTakeOverTaskDetailsResponse), nil
	}
}

// ShowTakeOverTaskDetailsInvoker 查询托管任务详情
func (c *VodClient) ShowTakeOverTaskDetailsInvoker(request *model.ShowTakeOverTaskDetailsRequest) *ShowTakeOverTaskDetailsInvoker {
	requestDef := GenReqDefForShowTakeOverTaskDetails()
	return &ShowTakeOverTaskDetailsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
