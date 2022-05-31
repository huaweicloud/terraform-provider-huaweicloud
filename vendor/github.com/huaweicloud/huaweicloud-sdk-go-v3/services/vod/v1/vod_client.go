package v1

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"
)

type VodClient struct {
	HcClient *http_client.HcHttpClient
}

func NewVodClient(hcClient *http_client.HcHttpClient) *VodClient {
	return &VodClient{HcClient: hcClient}
}

func VodClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// 取消媒资转码任务
//
// 取消媒资转码任务，只能取消排队中的转码任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CancelAssetTranscodeTask(request *model.CancelAssetTranscodeTaskRequest) (*model.CancelAssetTranscodeTaskResponse, error) {
	requestDef := GenReqDefForCancelAssetTranscodeTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CancelAssetTranscodeTaskResponse), nil
	}
}

// 取消提取音频任务
//
// 取消提取音频任务，只有排队中的提取音频任务才可以取消。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CancelExtractAudioTask(request *model.CancelExtractAudioTaskRequest) (*model.CancelExtractAudioTaskResponse, error) {
	requestDef := GenReqDefForCancelExtractAudioTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CancelExtractAudioTaskResponse), nil
	}
}

// 上传检验
//
// 校验媒资文件是否已存储于视频点播服务中。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CheckMd5Duplication(request *model.CheckMd5DuplicationRequest) (*model.CheckMd5DuplicationResponse, error) {
	requestDef := GenReqDefForCheckMd5Duplication()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CheckMd5DuplicationResponse), nil
	}
}

// 确认媒资上传
//
// 媒资分段上传完成后，需要调用此接口通知点播服务媒资上传的状态，表示媒资上传创建完成。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ConfirmAssetUpload(request *model.ConfirmAssetUploadRequest) (*model.ConfirmAssetUploadResponse, error) {
	requestDef := GenReqDefForConfirmAssetUpload()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ConfirmAssetUploadResponse), nil
	}
}

// 确认水印图片上传
//
// 确认水印图片上传状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ConfirmImageUpload(request *model.ConfirmImageUploadRequest) (*model.ConfirmImageUploadResponse, error) {
	requestDef := GenReqDefForConfirmImageUpload()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ConfirmImageUploadResponse), nil
	}
}

// 创建媒资：上传方式
//
// 调用该接口创建媒资时，需要将对应的媒资文件上传到点播服务的OBS桶中。
//
// 若上传的单媒资文件大小小于20M，则可以直接用PUT方法对该接口返回的地址进行上传。具体使用方法请参考[示例1：媒资上传（20M以下）](https://support.huaweicloud.com/api-vod/vod_04_0195.html)。
//
// 若上传的单个媒资大小大于20M，则需要进行二进制流分割后上传，该接口的具体使用方法请参考[示例2：媒资分段上传（20M以上）](https://support.huaweicloud.com/api-vod/vod_04_0216.html)。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreateAssetByFileUpload(request *model.CreateAssetByFileUploadRequest) (*model.CreateAssetByFileUploadResponse, error) {
	requestDef := GenReqDefForCreateAssetByFileUpload()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAssetByFileUploadResponse), nil
	}
}

// 创建媒资分类
//
// 创建媒资分类。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreateAssetCategory(request *model.CreateAssetCategoryRequest) (*model.CreateAssetCategoryResponse, error) {
	requestDef := GenReqDefForCreateAssetCategory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAssetCategoryResponse), nil
	}
}

// 媒资处理
//
// 实现视频转码、截图、加密等处理。既可以同时启动多种操作，也可以只启动一种操作。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreateAssetProcessTask(request *model.CreateAssetProcessTaskRequest) (*model.CreateAssetProcessTaskResponse, error) {
	requestDef := GenReqDefForCreateAssetProcessTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAssetProcessTaskResponse), nil
	}
}

// 创建审核媒资任务
//
// 对上传的媒资进行审核。审核后，可以调用[查询媒资详细信息](https://support.huaweicloud.com/api-vod/vod_04_0202.html)接口查看审核结果。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreateAssetReviewTask(request *model.CreateAssetReviewTaskRequest) (*model.CreateAssetReviewTaskResponse, error) {
	requestDef := GenReqDefForCreateAssetReviewTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAssetReviewTaskResponse), nil
	}
}

// 音频提取
//
// 本接口为异步接口，创建音频提取任务下发成功后会返回asset_id和提取的audio_asset_id，但此时音频提取任务并没有立即完成，可通过消息订阅界面配置的音频提取完成事件来获取音频提取任务完成与否。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreateExtractAudioTask(request *model.CreateExtractAudioTaskRequest) (*model.CreateExtractAudioTaskResponse, error) {
	requestDef := GenReqDefForCreateExtractAudioTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateExtractAudioTaskResponse), nil
	}
}

// CDN预热
//
// 媒资发布后，可通过指定媒资ID或URL向CDN预热。用户初次请求时，将由CDN节点提供请求媒资，加快用户下载缓存时间，提高用户体验。单租户每天最多预热1000个。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreatePreheatingAsset(request *model.CreatePreheatingAssetRequest) (*model.CreatePreheatingAssetResponse, error) {
	requestDef := GenReqDefForCreatePreheatingAsset()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePreheatingAssetResponse), nil
	}
}

// 创建媒资：OBS托管方式
//
// 通过存量托管的方式，将已存储在OBS桶中的音视频文件同步到点播服务。
//
// OBS托管方式分为增量托管和存量托管，增量托管暂只支持通过视频点播控制台配置，配置后，若OBS有新增音视频文件，则会自动同步到点播服务中，具体请参见[增量托管](https://support.huaweicloud.com/usermanual-vod/vod010032.html)。两个托管方式都需要先将对应的OBS桶授权给点播服务，具体请参见[桶授权](https://support.huaweicloud.com/usermanual-vod/vod010031.html)。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreateTakeOverTask(request *model.CreateTakeOverTaskRequest) (*model.CreateTakeOverTaskResponse, error) {
	requestDef := GenReqDefForCreateTakeOverTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTakeOverTaskResponse), nil
	}
}

// 创建自定义转码模板组
//
// 创建自定义转码模板组。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreateTemplateGroup(request *model.CreateTemplateGroupRequest) (*model.CreateTemplateGroupResponse, error) {
	requestDef := GenReqDefForCreateTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTemplateGroupResponse), nil
	}
}

// 创建水印模板
//
// 创建水印模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) CreateWatermarkTemplate(request *model.CreateWatermarkTemplateRequest) (*model.CreateWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForCreateWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateWatermarkTemplateResponse), nil
	}
}

// 删除媒资分类
//
// 删除媒资分类。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) DeleteAssetCategory(request *model.DeleteAssetCategoryRequest) (*model.DeleteAssetCategoryResponse, error) {
	requestDef := GenReqDefForDeleteAssetCategory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAssetCategoryResponse), nil
	}
}

// 删除媒资
//
// 删除媒资。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) DeleteAssets(request *model.DeleteAssetsRequest) (*model.DeleteAssetsResponse, error) {
	requestDef := GenReqDefForDeleteAssets()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAssetsResponse), nil
	}
}

// 删除自定义转码模板组
//
// 删除自定义转码模板组。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) DeleteTemplateGroup(request *model.DeleteTemplateGroupRequest) (*model.DeleteTemplateGroupResponse, error) {
	requestDef := GenReqDefForDeleteTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTemplateGroupResponse), nil
	}
}

// 删除水印模板
//
// 删除水印模板
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) DeleteWatermarkTemplate(request *model.DeleteWatermarkTemplateRequest) (*model.DeleteWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForDeleteWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteWatermarkTemplateResponse), nil
	}
}

// 查询指定分类信息
//
// 查询指定分类信息，及其子分类（即下一级分类）的列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ListAssetCategory(request *model.ListAssetCategoryRequest) (*model.ListAssetCategoryResponse, error) {
	requestDef := GenReqDefForListAssetCategory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAssetCategoryResponse), nil
	}
}

// 查询媒资列表
//
// 查询媒资列表，列表中的每一条记录包含媒资的概要信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ListAssetList(request *model.ListAssetListRequest) (*model.ListAssetListResponse, error) {
	requestDef := GenReqDefForListAssetList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAssetListResponse), nil
	}
}

// 查询域名播放日志
//
// 查询指定点播域名某段时间内在CDN的相关日志。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ListDomainLogs(request *model.ListDomainLogsRequest) (*model.ListDomainLogsResponse, error) {
	requestDef := GenReqDefForListDomainLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDomainLogsResponse), nil
	}
}

// 查询转码模板组列表
//
// 查询转码模板组列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ListTemplateGroup(request *model.ListTemplateGroupRequest) (*model.ListTemplateGroupResponse, error) {
	requestDef := GenReqDefForListTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTemplateGroupResponse), nil
	}
}

// 查询TopN媒资信息
//
// 查询指定域名在指定日期播放次数排名Top 100的媒资统计数据。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ListTopStatistics(request *model.ListTopStatisticsRequest) (*model.ListTopStatisticsResponse, error) {
	requestDef := GenReqDefForListTopStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTopStatisticsResponse), nil
	}
}

// 查询水印列表
//
// 查询水印模板
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ListWatermarkTemplate(request *model.ListWatermarkTemplateRequest) (*model.ListWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForListWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListWatermarkTemplateResponse), nil
	}
}

// 创建媒资：OBS转存方式
//
// 若您在使用点播服务前，已经在OBS桶中存储了音视频文件，您可以使用该接口将存储在OBS桶中的音视频文件转存到点播服务中，使用点播服务的音视频管理功能。调用该接口前，您需要调用[桶授权](https://support.huaweicloud.com/api-vod/vod_04_0199.html)接口，将存储音视频文件的OBS桶授权给点播服务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) PublishAssetFromObs(request *model.PublishAssetFromObsRequest) (*model.PublishAssetFromObsResponse, error) {
	requestDef := GenReqDefForPublishAssetFromObs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.PublishAssetFromObsResponse), nil
	}
}

// 媒资发布
//
// 将媒资设置为发布状态。支持批量发布。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) PublishAssets(request *model.PublishAssetsRequest) (*model.PublishAssetsResponse, error) {
	requestDef := GenReqDefForPublishAssets()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.PublishAssetsResponse), nil
	}
}

// 密钥查询
//
// 终端播放HLS加密视频时，向租户管理系统请求密钥，租户管理系统先查询其本地有没有已缓存的密钥，没有时则调用此接口向VOD查询。该接口的具体使用场景请参见[通过HLS加密防止视频泄露](https://support.huaweicloud.com/bestpractice-vod/vod_10_0004.html)。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowAssetCipher(request *model.ShowAssetCipherRequest) (*model.ShowAssetCipherResponse, error) {
	requestDef := GenReqDefForShowAssetCipher()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetCipherResponse), nil
	}
}

// 查询指定媒资的详细信息
//
// 查询指定媒资的详细信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowAssetDetail(request *model.ShowAssetDetailRequest) (*model.ShowAssetDetailResponse, error) {
	requestDef := GenReqDefForShowAssetDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetDetailResponse), nil
	}
}

// 查询媒资信息
//
// 查询媒资信息，支持指定媒资ID、分类、状态、起止时间查询。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowAssetMeta(request *model.ShowAssetMetaRequest) (*model.ShowAssetMetaResponse, error) {
	requestDef := GenReqDefForShowAssetMeta()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetMetaResponse), nil
	}
}

// 获取分段上传授权
//
// 客户端请求创建媒资时，如果媒资文件超过20MB，需采用分段的方式向OBS上传，在每次与OBS交互前，客户端需通过此接口获取到授权方可与OBS交互。
//
// 该接口可以获取[初始化多段上传任务](https://support.huaweicloud.com/api-obs/obs_04_0098.html)、[上传段](https://support.huaweicloud.com/api-obs/obs_04_0099.html)、[合并段](https://support.huaweicloud.com/api-obs/obs_04_0102.html)、[列举已上传段](https://support.huaweicloud.com/api-obs/obs_04_0101.html)、[取消段合并](https://support.huaweicloud.com/api-obs/obs_04_0103.html)的带有临时授权的URL，用户需要根据OBS的接口文档配置相应请求的HTTP请求方法、请求头、请求体，然后请求对应的带有临时授权的URL。
//
// 视频分段上传方式和OBS的接口文档保持一致，包括HTTP请求方法、请求头、请求体等各种入参，此接口的作用是为用户生成带有鉴权信息的URL（鉴权信息即query_str），用来替换OBS接口中对应的URL，临时给用户开通向点播服务的桶上传文件的权限。
//
// 调用获取授权接口时需要传入bucket、object_key、http_verb，其中bucket和object_key是由[创建媒资：上传方式](https://support.huaweicloud.com/api-vod/vod_04_0196.html)接口中返回的响应体中的target字段获得的bucket和object，http_verb需要根据指定的操作选择。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowAssetTempAuthority(request *model.ShowAssetTempAuthorityRequest) (*model.ShowAssetTempAuthorityResponse, error) {
	requestDef := GenReqDefForShowAssetTempAuthority()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAssetTempAuthorityResponse), nil
	}
}

// 查询CDN统计信息
//
// 查询CDN的统计数据，包括流量、峰值带宽、请求总数、请求命中率、流量命中率。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowCdnStatistics(request *model.ShowCdnStatisticsRequest) (*model.ShowCdnStatisticsResponse, error) {
	requestDef := GenReqDefForShowCdnStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCdnStatisticsResponse), nil
	}
}

// 查询CDN预热
//
// 查询预热结果。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowPreheatingAsset(request *model.ShowPreheatingAssetRequest) (*model.ShowPreheatingAssetResponse, error) {
	requestDef := GenReqDefForShowPreheatingAsset()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowPreheatingAssetResponse), nil
	}
}

// 查询源站统计信息
//
// 查询点播源站的统计数据，包括流量、存储空间、转码时长。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowVodStatistics(request *model.ShowVodStatisticsRequest) (*model.ShowVodStatisticsResponse, error) {
	requestDef := GenReqDefForShowVodStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowVodStatisticsResponse), nil
	}
}

// 媒资发布取消
//
// 将媒资设置为未发布状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UnpublishAssets(request *model.UnpublishAssetsRequest) (*model.UnpublishAssetsResponse, error) {
	requestDef := GenReqDefForUnpublishAssets()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UnpublishAssetsResponse), nil
	}
}

// 视频更新
//
// 媒资创建后，单独上传封面、更新视频文件或更新已有封面。
//
// 如果是更新视频文件，更新完后要通过[确认媒资上传](https://support.huaweicloud.com/api-vod/vod_04_0198.html)接口通知点播服务。
//
// 如果是更新封面或单独上传封面，则不需通知。
//
// 更新视频可以使用分段上传，具体方式可以参考[示例2：媒资分段上传（20M以上）](https://support.huaweicloud.com/api-vod/vod_04_0216.html)。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UpdateAsset(request *model.UpdateAssetRequest) (*model.UpdateAssetResponse, error) {
	requestDef := GenReqDefForUpdateAsset()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAssetResponse), nil
	}
}

// 修改媒资分类
//
// 修改媒资分类。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UpdateAssetCategory(request *model.UpdateAssetCategoryRequest) (*model.UpdateAssetCategoryResponse, error) {
	requestDef := GenReqDefForUpdateAssetCategory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAssetCategoryResponse), nil
	}
}

// 修改媒资属性
//
// 修改媒资属性。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UpdateAssetMeta(request *model.UpdateAssetMetaRequest) (*model.UpdateAssetMetaResponse, error) {
	requestDef := GenReqDefForUpdateAssetMeta()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAssetMetaResponse), nil
	}
}

// 桶授权
//
// 用户可以通过该接口将OBS桶授权给点播服务或取消点播服务的授权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UpdateBucketAuthorized(request *model.UpdateBucketAuthorizedRequest) (*model.UpdateBucketAuthorizedResponse, error) {
	requestDef := GenReqDefForUpdateBucketAuthorized()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateBucketAuthorizedResponse), nil
	}
}

// 设置封面
//
// 将视频截图生成的某张图片设置成封面。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UpdateCoverByThumbnail(request *model.UpdateCoverByThumbnailRequest) (*model.UpdateCoverByThumbnailResponse, error) {
	requestDef := GenReqDefForUpdateCoverByThumbnail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCoverByThumbnailResponse), nil
	}
}

// 修改自定义转码模板组
//
// 修改自定义转码模板组。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UpdateTemplateGroup(request *model.UpdateTemplateGroupRequest) (*model.UpdateTemplateGroupResponse, error) {
	requestDef := GenReqDefForUpdateTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTemplateGroupResponse), nil
	}
}

// 修改水印模板
//
// 修改水印模板
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UpdateWatermarkTemplate(request *model.UpdateWatermarkTemplateRequest) (*model.UpdateWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForUpdateWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateWatermarkTemplateResponse), nil
	}
}

// 创建媒资：URL拉取注入
//
// 基于音视频源文件URL，将音视频文件离线拉取上传到点播服务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) UploadMetaDataByUrl(request *model.UploadMetaDataByUrlRequest) (*model.UploadMetaDataByUrlResponse, error) {
	requestDef := GenReqDefForUploadMetaDataByUrl()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UploadMetaDataByUrlResponse), nil
	}
}

// 查询托管任务
//
// 查询OBS存量托管任务列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ListTakeOverTask(request *model.ListTakeOverTaskRequest) (*model.ListTakeOverTaskResponse, error) {
	requestDef := GenReqDefForListTakeOverTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTakeOverTaskResponse), nil
	}
}

// 查询托管媒资详情
//
// 查询OBS托管媒资的详细信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowTakeOverAssetDetails(request *model.ShowTakeOverAssetDetailsRequest) (*model.ShowTakeOverAssetDetailsResponse, error) {
	requestDef := GenReqDefForShowTakeOverAssetDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTakeOverAssetDetailsResponse), nil
	}
}

// 查询托管任务详情
//
// 查询OBS存量托管任务详情。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VodClient) ShowTakeOverTaskDetails(request *model.ShowTakeOverTaskDetailsRequest) (*model.ShowTakeOverTaskDetailsResponse, error) {
	requestDef := GenReqDefForShowTakeOverTaskDetails()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTakeOverTaskDetailsResponse), nil
	}
}
