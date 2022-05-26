package v1

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1/model"
)

type MpcClient struct {
	HcClient *http_client.HcHttpClient
}

func NewMpcClient(hcClient *http_client.HcHttpClient) *MpcClient {
	return &MpcClient{HcClient: hcClient}
}

func MpcClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// CreateAnimatedGraphicsTask 新建转动图任务
//
// 创建动图任务，用于将完整的视频文件或视频文件中的一部分转换为动态图文件，暂只支持输出GIF文件。
// 待转动图的视频文件需要存储在与媒体处理服务同区域的OBS桶中，且该OBS桶已授权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateAnimatedGraphicsTask(request *model.CreateAnimatedGraphicsTaskRequest) (*model.CreateAnimatedGraphicsTaskResponse, error) {
	requestDef := GenReqDefForCreateAnimatedGraphicsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAnimatedGraphicsTaskResponse), nil
	}
}

// CreateAnimatedGraphicsTaskInvoker 新建转动图任务
func (c *MpcClient) CreateAnimatedGraphicsTaskInvoker(request *model.CreateAnimatedGraphicsTaskRequest) *CreateAnimatedGraphicsTaskInvoker {
	requestDef := GenReqDefForCreateAnimatedGraphicsTask()
	return &CreateAnimatedGraphicsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAnimatedGraphicsTask 取消转动图任务
//
// 取消已下发的生成动图任务，仅支持取消正在排队中的任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteAnimatedGraphicsTask(request *model.DeleteAnimatedGraphicsTaskRequest) (*model.DeleteAnimatedGraphicsTaskResponse, error) {
	requestDef := GenReqDefForDeleteAnimatedGraphicsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAnimatedGraphicsTaskResponse), nil
	}
}

// DeleteAnimatedGraphicsTaskInvoker 取消转动图任务
func (c *MpcClient) DeleteAnimatedGraphicsTaskInvoker(request *model.DeleteAnimatedGraphicsTaskRequest) *DeleteAnimatedGraphicsTaskInvoker {
	requestDef := GenReqDefForDeleteAnimatedGraphicsTask()
	return &DeleteAnimatedGraphicsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAnimatedGraphicsTask 查询转动图任务
//
// 查询动图任务的状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListAnimatedGraphicsTask(request *model.ListAnimatedGraphicsTaskRequest) (*model.ListAnimatedGraphicsTaskResponse, error) {
	requestDef := GenReqDefForListAnimatedGraphicsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAnimatedGraphicsTaskResponse), nil
	}
}

// ListAnimatedGraphicsTaskInvoker 查询转动图任务
func (c *MpcClient) ListAnimatedGraphicsTaskInvoker(request *model.ListAnimatedGraphicsTaskRequest) *ListAnimatedGraphicsTaskInvoker {
	requestDef := GenReqDefForListAnimatedGraphicsTask()
	return &ListAnimatedGraphicsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateEditingJob 新建剪辑任务
//
// 创建剪辑任务，用于将多个视频文件进行裁剪成多个视频分段，并且可以把这些视频分段合并成一个视频，剪切和拼接功能可以单独使用。
// 待剪辑的视频文件需要存储在与媒体处理服务同区域的OBS桶中，且该OBS桶已授权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateEditingJob(request *model.CreateEditingJobRequest) (*model.CreateEditingJobResponse, error) {
	requestDef := GenReqDefForCreateEditingJob()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateEditingJobResponse), nil
	}
}

// CreateEditingJobInvoker 新建剪辑任务
func (c *MpcClient) CreateEditingJobInvoker(request *model.CreateEditingJobRequest) *CreateEditingJobInvoker {
	requestDef := GenReqDefForCreateEditingJob()
	return &CreateEditingJobInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteEditingJob 取消剪辑任务
//
// 取消已下发的生成剪辑任务，仅支持取消正在排队中的任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteEditingJob(request *model.DeleteEditingJobRequest) (*model.DeleteEditingJobResponse, error) {
	requestDef := GenReqDefForDeleteEditingJob()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteEditingJobResponse), nil
	}
}

// DeleteEditingJobInvoker 取消剪辑任务
func (c *MpcClient) DeleteEditingJobInvoker(request *model.DeleteEditingJobRequest) *DeleteEditingJobInvoker {
	requestDef := GenReqDefForDeleteEditingJob()
	return &DeleteEditingJobInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListEditingJob 查询剪辑任务
//
// 查询剪辑任务的状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListEditingJob(request *model.ListEditingJobRequest) (*model.ListEditingJobResponse, error) {
	requestDef := GenReqDefForListEditingJob()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListEditingJobResponse), nil
	}
}

// ListEditingJobInvoker 查询剪辑任务
func (c *MpcClient) ListEditingJobInvoker(request *model.ListEditingJobRequest) *ListEditingJobInvoker {
	requestDef := GenReqDefForListEditingJob()
	return &ListEditingJobInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateEncryptTask 新建独立加密任务
//
// 支持独立加密，包括创建、查询、删除独立加密任务。
//
// 约束：
//   - 只支持转码后的文件进行加密。
//   - 加密的文件必须是m3u8或者mpd结尾的文件。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateEncryptTask(request *model.CreateEncryptTaskRequest) (*model.CreateEncryptTaskResponse, error) {
	requestDef := GenReqDefForCreateEncryptTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateEncryptTaskResponse), nil
	}
}

// CreateEncryptTaskInvoker 新建独立加密任务
func (c *MpcClient) CreateEncryptTaskInvoker(request *model.CreateEncryptTaskRequest) *CreateEncryptTaskInvoker {
	requestDef := GenReqDefForCreateEncryptTask()
	return &CreateEncryptTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteEncryptTask 取消独立加密任务
//
// 取消独立加密任务。
//
// 约束：
//
//   只能取消正在任务队列中排队的任务。已开始加密或已完成的加密任务不能取消。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteEncryptTask(request *model.DeleteEncryptTaskRequest) (*model.DeleteEncryptTaskResponse, error) {
	requestDef := GenReqDefForDeleteEncryptTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteEncryptTaskResponse), nil
	}
}

// DeleteEncryptTaskInvoker 取消独立加密任务
func (c *MpcClient) DeleteEncryptTaskInvoker(request *model.DeleteEncryptTaskRequest) *DeleteEncryptTaskInvoker {
	requestDef := GenReqDefForDeleteEncryptTask()
	return &DeleteEncryptTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListEncryptTask 查询独立加密任务
//
// 查询独立加密任务状态。返回任务执行结果或当前状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListEncryptTask(request *model.ListEncryptTaskRequest) (*model.ListEncryptTaskResponse, error) {
	requestDef := GenReqDefForListEncryptTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListEncryptTaskResponse), nil
	}
}

// ListEncryptTaskInvoker 查询独立加密任务
func (c *MpcClient) ListEncryptTaskInvoker(request *model.ListEncryptTaskRequest) *ListEncryptTaskInvoker {
	requestDef := GenReqDefForListEncryptTask()
	return &ListEncryptTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateExtractTask 新建视频解析任务
//
// 创建视频解析任务，解析视频元数据。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateExtractTask(request *model.CreateExtractTaskRequest) (*model.CreateExtractTaskResponse, error) {
	requestDef := GenReqDefForCreateExtractTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateExtractTaskResponse), nil
	}
}

// CreateExtractTaskInvoker 新建视频解析任务
func (c *MpcClient) CreateExtractTaskInvoker(request *model.CreateExtractTaskRequest) *CreateExtractTaskInvoker {
	requestDef := GenReqDefForCreateExtractTask()
	return &CreateExtractTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteExtractTask 取消视频解析任务
//
// 取消已下发的视频解析任务，仅支持取消正在排队中的任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteExtractTask(request *model.DeleteExtractTaskRequest) (*model.DeleteExtractTaskResponse, error) {
	requestDef := GenReqDefForDeleteExtractTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteExtractTaskResponse), nil
	}
}

// DeleteExtractTaskInvoker 取消视频解析任务
func (c *MpcClient) DeleteExtractTaskInvoker(request *model.DeleteExtractTaskRequest) *DeleteExtractTaskInvoker {
	requestDef := GenReqDefForDeleteExtractTask()
	return &DeleteExtractTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListExtractTask 查询视频解析任务
//
// 查询解析任务的状态和结果。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListExtractTask(request *model.ListExtractTaskRequest) (*model.ListExtractTaskResponse, error) {
	requestDef := GenReqDefForListExtractTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListExtractTaskResponse), nil
	}
}

// ListExtractTaskInvoker 查询视频解析任务
func (c *MpcClient) ListExtractTaskInvoker(request *model.ListExtractTaskRequest) *ListExtractTaskInvoker {
	requestDef := GenReqDefForListExtractTask()
	return &ListExtractTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateMbTasksReport 合并多声道任务、重置声轨任务上报接口
//
// ## 典型场景 ##
//   合并音频多声道文件任务、重置音频文件声轨任务上报结果接口。
// ## 接口功能 ##
//   合并音频多声道文件任务、重置音频文件声轨任务上报结果接口。
// ## 接口约束 ##
//   无。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateMbTasksReport(request *model.CreateMbTasksReportRequest) (*model.CreateMbTasksReportResponse, error) {
	requestDef := GenReqDefForCreateMbTasksReport()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateMbTasksReportResponse), nil
	}
}

// CreateMbTasksReportInvoker 合并多声道任务、重置声轨任务上报接口
func (c *MpcClient) CreateMbTasksReportInvoker(request *model.CreateMbTasksReportRequest) *CreateMbTasksReportInvoker {
	requestDef := GenReqDefForCreateMbTasksReport()
	return &CreateMbTasksReportInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateMergeChannelsTask 创建声道合并任务
//
// 创建声道合并任务，合并声道任务支持将每个声道各放一个文件中的片源，合并为单个音频文件。
// 执行合并声道的源音频文件需要存储在与媒体处理服务同区域的OBS桶中，且该OBS桶已授权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateMergeChannelsTask(request *model.CreateMergeChannelsTaskRequest) (*model.CreateMergeChannelsTaskResponse, error) {
	requestDef := GenReqDefForCreateMergeChannelsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateMergeChannelsTaskResponse), nil
	}
}

// CreateMergeChannelsTaskInvoker 创建声道合并任务
func (c *MpcClient) CreateMergeChannelsTaskInvoker(request *model.CreateMergeChannelsTaskRequest) *CreateMergeChannelsTaskInvoker {
	requestDef := GenReqDefForCreateMergeChannelsTask()
	return &CreateMergeChannelsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateResetTracksTask 创建音轨重置任务
//
// 创建音轨重置任务，重置音轨任务支持按人工指定关系声道layout，语言标签，转封装片源，使其满足转码输入。
// 执行音轨重置的源音频文件需要存储在与媒体处理服务同区域的OBS桶中，且该OBS桶已授权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateResetTracksTask(request *model.CreateResetTracksTaskRequest) (*model.CreateResetTracksTaskResponse, error) {
	requestDef := GenReqDefForCreateResetTracksTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateResetTracksTaskResponse), nil
	}
}

// CreateResetTracksTaskInvoker 创建音轨重置任务
func (c *MpcClient) CreateResetTracksTaskInvoker(request *model.CreateResetTracksTaskRequest) *CreateResetTracksTaskInvoker {
	requestDef := GenReqDefForCreateResetTracksTask()
	return &CreateResetTracksTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteMergeChannelsTask 取消声道合并任务
//
// 取消合并音频多声道文件。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteMergeChannelsTask(request *model.DeleteMergeChannelsTaskRequest) (*model.DeleteMergeChannelsTaskResponse, error) {
	requestDef := GenReqDefForDeleteMergeChannelsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteMergeChannelsTaskResponse), nil
	}
}

// DeleteMergeChannelsTaskInvoker 取消声道合并任务
func (c *MpcClient) DeleteMergeChannelsTaskInvoker(request *model.DeleteMergeChannelsTaskRequest) *DeleteMergeChannelsTaskInvoker {
	requestDef := GenReqDefForDeleteMergeChannelsTask()
	return &DeleteMergeChannelsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteResetTracksTask 取消音轨重置任务
//
// 取消重置音频文件声轨任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteResetTracksTask(request *model.DeleteResetTracksTaskRequest) (*model.DeleteResetTracksTaskResponse, error) {
	requestDef := GenReqDefForDeleteResetTracksTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteResetTracksTaskResponse), nil
	}
}

// DeleteResetTracksTaskInvoker 取消音轨重置任务
func (c *MpcClient) DeleteResetTracksTaskInvoker(request *model.DeleteResetTracksTaskRequest) *DeleteResetTracksTaskInvoker {
	requestDef := GenReqDefForDeleteResetTracksTask()
	return &DeleteResetTracksTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMergeChannelsTask 查询声道合并任务
//
// 查询声道合并任务的状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListMergeChannelsTask(request *model.ListMergeChannelsTaskRequest) (*model.ListMergeChannelsTaskResponse, error) {
	requestDef := GenReqDefForListMergeChannelsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMergeChannelsTaskResponse), nil
	}
}

// ListMergeChannelsTaskInvoker 查询声道合并任务
func (c *MpcClient) ListMergeChannelsTaskInvoker(request *model.ListMergeChannelsTaskRequest) *ListMergeChannelsTaskInvoker {
	requestDef := GenReqDefForListMergeChannelsTask()
	return &ListMergeChannelsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListResetTracksTask 查询音轨重置任务
//
// 查询音轨重置任务的状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListResetTracksTask(request *model.ListResetTracksTaskRequest) (*model.ListResetTracksTaskResponse, error) {
	requestDef := GenReqDefForListResetTracksTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListResetTracksTaskResponse), nil
	}
}

// ListResetTracksTaskInvoker 查询音轨重置任务
func (c *MpcClient) ListResetTracksTaskInvoker(request *model.ListResetTracksTaskRequest) *ListResetTracksTaskInvoker {
	requestDef := GenReqDefForListResetTracksTask()
	return &ListResetTracksTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateMediaProcessTask 创建视频增强任务
//
// ## 典型场景 ##
//   创建视频增强任务。
//
// ## 接口功能 ##
//   创建视频增强任务。
//
// ## 接口约束 ##
//   无。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateMediaProcessTask(request *model.CreateMediaProcessTaskRequest) (*model.CreateMediaProcessTaskResponse, error) {
	requestDef := GenReqDefForCreateMediaProcessTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateMediaProcessTaskResponse), nil
	}
}

// CreateMediaProcessTaskInvoker 创建视频增强任务
func (c *MpcClient) CreateMediaProcessTaskInvoker(request *model.CreateMediaProcessTaskRequest) *CreateMediaProcessTaskInvoker {
	requestDef := GenReqDefForCreateMediaProcessTask()
	return &CreateMediaProcessTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteMediaProcessTask 取消视频增强任务
//
// ## 典型场景 ##
//   取消视频增强任务。
//
// ## 接口功能 ##
//   取消视频增强任务。
//
// ## 接口约束 ##
//   仅可删除正在排队中的任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteMediaProcessTask(request *model.DeleteMediaProcessTaskRequest) (*model.DeleteMediaProcessTaskResponse, error) {
	requestDef := GenReqDefForDeleteMediaProcessTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteMediaProcessTaskResponse), nil
	}
}

// DeleteMediaProcessTaskInvoker 取消视频增强任务
func (c *MpcClient) DeleteMediaProcessTaskInvoker(request *model.DeleteMediaProcessTaskRequest) *DeleteMediaProcessTaskInvoker {
	requestDef := GenReqDefForDeleteMediaProcessTask()
	return &DeleteMediaProcessTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMediaProcessTask 查询视频增强任务
//
// ## 典型场景 ##
//   查询视频增强任务。
//
// ## 接口功能 ##
//   查询视频增强任务。
//
// ## 接口约束 ##
//   无。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListMediaProcessTask(request *model.ListMediaProcessTaskRequest) (*model.ListMediaProcessTaskResponse, error) {
	requestDef := GenReqDefForListMediaProcessTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMediaProcessTaskResponse), nil
	}
}

// ListMediaProcessTaskInvoker 查询视频增强任务
func (c *MpcClient) ListMediaProcessTaskInvoker(request *model.ListMediaProcessTaskRequest) *ListMediaProcessTaskInvoker {
	requestDef := GenReqDefForListMediaProcessTask()
	return &ListMediaProcessTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateMpeCallBack mpe通知mpc
//
// ## 典型场景 ##
//   mpe通知mpc。
// ## 接口功能 ##
//   mpe调用此接口通知mpc转封装等结果。
// ## 接口约束 ##
//   无。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateMpeCallBack(request *model.CreateMpeCallBackRequest) (*model.CreateMpeCallBackResponse, error) {
	requestDef := GenReqDefForCreateMpeCallBack()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateMpeCallBackResponse), nil
	}
}

// CreateMpeCallBackInvoker mpe通知mpc
func (c *MpcClient) CreateMpeCallBackInvoker(request *model.CreateMpeCallBackRequest) *CreateMpeCallBackInvoker {
	requestDef := GenReqDefForCreateMpeCallBack()
	return &CreateMpeCallBackInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateQualityEnhanceTemplate 创建视频增强模板
//
// 创建视频增强模板
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateQualityEnhanceTemplate(request *model.CreateQualityEnhanceTemplateRequest) (*model.CreateQualityEnhanceTemplateResponse, error) {
	requestDef := GenReqDefForCreateQualityEnhanceTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateQualityEnhanceTemplateResponse), nil
	}
}

// CreateQualityEnhanceTemplateInvoker 创建视频增强模板
func (c *MpcClient) CreateQualityEnhanceTemplateInvoker(request *model.CreateQualityEnhanceTemplateRequest) *CreateQualityEnhanceTemplateInvoker {
	requestDef := GenReqDefForCreateQualityEnhanceTemplate()
	return &CreateQualityEnhanceTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteQualityEnhanceTemplate 删除用户视频增强模板
//
// 删除用户视频增强模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteQualityEnhanceTemplate(request *model.DeleteQualityEnhanceTemplateRequest) (*model.DeleteQualityEnhanceTemplateResponse, error) {
	requestDef := GenReqDefForDeleteQualityEnhanceTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteQualityEnhanceTemplateResponse), nil
	}
}

// DeleteQualityEnhanceTemplateInvoker 删除用户视频增强模板
func (c *MpcClient) DeleteQualityEnhanceTemplateInvoker(request *model.DeleteQualityEnhanceTemplateRequest) *DeleteQualityEnhanceTemplateInvoker {
	requestDef := GenReqDefForDeleteQualityEnhanceTemplate()
	return &DeleteQualityEnhanceTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListQualityEnhanceDefaultTemplate 查询视频增强预置模板
//
// 查询视频增强预置模板，返回所有结果。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListQualityEnhanceDefaultTemplate(request *model.ListQualityEnhanceDefaultTemplateRequest) (*model.ListQualityEnhanceDefaultTemplateResponse, error) {
	requestDef := GenReqDefForListQualityEnhanceDefaultTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListQualityEnhanceDefaultTemplateResponse), nil
	}
}

// ListQualityEnhanceDefaultTemplateInvoker 查询视频增强预置模板
func (c *MpcClient) ListQualityEnhanceDefaultTemplateInvoker(request *model.ListQualityEnhanceDefaultTemplateRequest) *ListQualityEnhanceDefaultTemplateInvoker {
	requestDef := GenReqDefForListQualityEnhanceDefaultTemplate()
	return &ListQualityEnhanceDefaultTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateQualityEnhanceTemplate 更新视频增强模板
//
// 更新视频增强模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) UpdateQualityEnhanceTemplate(request *model.UpdateQualityEnhanceTemplateRequest) (*model.UpdateQualityEnhanceTemplateResponse, error) {
	requestDef := GenReqDefForUpdateQualityEnhanceTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateQualityEnhanceTemplateResponse), nil
	}
}

// UpdateQualityEnhanceTemplateInvoker 更新视频增强模板
func (c *MpcClient) UpdateQualityEnhanceTemplateInvoker(request *model.UpdateQualityEnhanceTemplateRequest) *UpdateQualityEnhanceTemplateInvoker {
	requestDef := GenReqDefForUpdateQualityEnhanceTemplate()
	return &UpdateQualityEnhanceTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTranscodeDetail 查询媒资转码详情
//
// 查询媒资转码详情
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListTranscodeDetail(request *model.ListTranscodeDetailRequest) (*model.ListTranscodeDetailResponse, error) {
	requestDef := GenReqDefForListTranscodeDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTranscodeDetailResponse), nil
	}
}

// ListTranscodeDetailInvoker 查询媒资转码详情
func (c *MpcClient) ListTranscodeDetailInvoker(request *model.ListTranscodeDetailRequest) *ListTranscodeDetailInvoker {
	requestDef := GenReqDefForListTranscodeDetail()
	return &ListTranscodeDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CancelRemuxTask 取消转封装任务
//
// 取消已下发的转封装任务，仅支持取消正在排队中的任务。。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CancelRemuxTask(request *model.CancelRemuxTaskRequest) (*model.CancelRemuxTaskResponse, error) {
	requestDef := GenReqDefForCancelRemuxTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CancelRemuxTaskResponse), nil
	}
}

// CancelRemuxTaskInvoker 取消转封装任务
func (c *MpcClient) CancelRemuxTaskInvoker(request *model.CancelRemuxTaskRequest) *CancelRemuxTaskInvoker {
	requestDef := GenReqDefForCancelRemuxTask()
	return &CancelRemuxTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRemuxTask 新建转封装任务
//
// 创建转封装任务，转换音视频文件的格式，但不改变其分辨率和码率。
// 待转封装的媒资文件需要存储在与媒体处理服务同区域的OBS桶中，且该OBS桶已授权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateRemuxTask(request *model.CreateRemuxTaskRequest) (*model.CreateRemuxTaskResponse, error) {
	requestDef := GenReqDefForCreateRemuxTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRemuxTaskResponse), nil
	}
}

// CreateRemuxTaskInvoker 新建转封装任务
func (c *MpcClient) CreateRemuxTaskInvoker(request *model.CreateRemuxTaskRequest) *CreateRemuxTaskInvoker {
	requestDef := GenReqDefForCreateRemuxTask()
	return &CreateRemuxTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRetryRemuxTask 重试转封装任务
//
// 对失败的转封装任务进行重试。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateRetryRemuxTask(request *model.CreateRetryRemuxTaskRequest) (*model.CreateRetryRemuxTaskResponse, error) {
	requestDef := GenReqDefForCreateRetryRemuxTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRetryRemuxTaskResponse), nil
	}
}

// CreateRetryRemuxTaskInvoker 重试转封装任务
func (c *MpcClient) CreateRetryRemuxTaskInvoker(request *model.CreateRetryRemuxTaskRequest) *CreateRetryRemuxTaskInvoker {
	requestDef := GenReqDefForCreateRetryRemuxTask()
	return &CreateRetryRemuxTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteRemuxTask 删除转封装任务(仅供Console调用)
//
// 删除转封装任务
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteRemuxTask(request *model.DeleteRemuxTaskRequest) (*model.DeleteRemuxTaskResponse, error) {
	requestDef := GenReqDefForDeleteRemuxTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteRemuxTaskResponse), nil
	}
}

// DeleteRemuxTaskInvoker 删除转封装任务(仅供Console调用)
func (c *MpcClient) DeleteRemuxTaskInvoker(request *model.DeleteRemuxTaskRequest) *DeleteRemuxTaskInvoker {
	requestDef := GenReqDefForDeleteRemuxTask()
	return &DeleteRemuxTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRemuxTask 查询转封装任务
//
// 查询转封装任务状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListRemuxTask(request *model.ListRemuxTaskRequest) (*model.ListRemuxTaskResponse, error) {
	requestDef := GenReqDefForListRemuxTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRemuxTaskResponse), nil
	}
}

// ListRemuxTaskInvoker 查询转封装任务
func (c *MpcClient) ListRemuxTaskInvoker(request *model.ListRemuxTaskRequest) *ListRemuxTaskInvoker {
	requestDef := GenReqDefForListRemuxTask()
	return &ListRemuxTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTemplateGroup 新建转码模板组
//
// 新建转码模板组，最多支持一进六出。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateTemplateGroup(request *model.CreateTemplateGroupRequest) (*model.CreateTemplateGroupResponse, error) {
	requestDef := GenReqDefForCreateTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTemplateGroupResponse), nil
	}
}

// CreateTemplateGroupInvoker 新建转码模板组
func (c *MpcClient) CreateTemplateGroupInvoker(request *model.CreateTemplateGroupRequest) *CreateTemplateGroupInvoker {
	requestDef := GenReqDefForCreateTemplateGroup()
	return &CreateTemplateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTemplateGroup 删除转码模板组
//
// 删除转码模板组。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteTemplateGroup(request *model.DeleteTemplateGroupRequest) (*model.DeleteTemplateGroupResponse, error) {
	requestDef := GenReqDefForDeleteTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTemplateGroupResponse), nil
	}
}

// DeleteTemplateGroupInvoker 删除转码模板组
func (c *MpcClient) DeleteTemplateGroupInvoker(request *model.DeleteTemplateGroupRequest) *DeleteTemplateGroupInvoker {
	requestDef := GenReqDefForDeleteTemplateGroup()
	return &DeleteTemplateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTemplateGroup 查询转码模板组
//
// 查询转码模板组列表。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListTemplateGroup(request *model.ListTemplateGroupRequest) (*model.ListTemplateGroupResponse, error) {
	requestDef := GenReqDefForListTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTemplateGroupResponse), nil
	}
}

// ListTemplateGroupInvoker 查询转码模板组
func (c *MpcClient) ListTemplateGroupInvoker(request *model.ListTemplateGroupRequest) *ListTemplateGroupInvoker {
	requestDef := GenReqDefForListTemplateGroup()
	return &ListTemplateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTemplateGroup 更新转码模板组
//
// 修改模板组接口。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) UpdateTemplateGroup(request *model.UpdateTemplateGroupRequest) (*model.UpdateTemplateGroupResponse, error) {
	requestDef := GenReqDefForUpdateTemplateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTemplateGroupResponse), nil
	}
}

// UpdateTemplateGroupInvoker 更新转码模板组
func (c *MpcClient) UpdateTemplateGroupInvoker(request *model.UpdateTemplateGroupRequest) *UpdateTemplateGroupInvoker {
	requestDef := GenReqDefForUpdateTemplateGroup()
	return &UpdateTemplateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateThumbnailsTask 新建截图任务
//
// 新建截图任务，视频截图将从首帧开始，按设置的时间间隔截图，最后截取末帧。
// 待截图的视频文件需要存储在与媒体处理服务同区域的OBS桶中，且该OBS桶已授权。
//
// 约束：
//   暂只支持生成JPG格式的图片文件。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateThumbnailsTask(request *model.CreateThumbnailsTaskRequest) (*model.CreateThumbnailsTaskResponse, error) {
	requestDef := GenReqDefForCreateThumbnailsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateThumbnailsTaskResponse), nil
	}
}

// CreateThumbnailsTaskInvoker 新建截图任务
func (c *MpcClient) CreateThumbnailsTaskInvoker(request *model.CreateThumbnailsTaskRequest) *CreateThumbnailsTaskInvoker {
	requestDef := GenReqDefForCreateThumbnailsTask()
	return &CreateThumbnailsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteThumbnailsTask 取消截图任务
//
// 取消已下发截图任务。
// 只能取消已接受尚在队列中等待处理的任务，已完成或正在执行阶段的任务不能取消。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteThumbnailsTask(request *model.DeleteThumbnailsTaskRequest) (*model.DeleteThumbnailsTaskResponse, error) {
	requestDef := GenReqDefForDeleteThumbnailsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteThumbnailsTaskResponse), nil
	}
}

// DeleteThumbnailsTaskInvoker 取消截图任务
func (c *MpcClient) DeleteThumbnailsTaskInvoker(request *model.DeleteThumbnailsTaskRequest) *DeleteThumbnailsTaskInvoker {
	requestDef := GenReqDefForDeleteThumbnailsTask()
	return &DeleteThumbnailsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListThumbnailsTask 查询截图任务
//
// 查询截图任务状态。返回任务执行结果，包括状态、输入、输出等信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListThumbnailsTask(request *model.ListThumbnailsTaskRequest) (*model.ListThumbnailsTaskResponse, error) {
	requestDef := GenReqDefForListThumbnailsTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListThumbnailsTaskResponse), nil
	}
}

// ListThumbnailsTaskInvoker 查询截图任务
func (c *MpcClient) ListThumbnailsTaskInvoker(request *model.ListThumbnailsTaskRequest) *ListThumbnailsTaskInvoker {
	requestDef := GenReqDefForListThumbnailsTask()
	return &ListThumbnailsTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTranscodingTask 新建转码任务
//
// 新建转码任务可以将视频进行转码，并在转码过程中压制水印、视频截图等。视频转码前需要配置转码模板。
// 待转码的音视频需要存储在与媒体处理服务同区域的OBS桶中，且该OBS桶已授权。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateTranscodingTask(request *model.CreateTranscodingTaskRequest) (*model.CreateTranscodingTaskResponse, error) {
	requestDef := GenReqDefForCreateTranscodingTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTranscodingTaskResponse), nil
	}
}

// CreateTranscodingTaskInvoker 新建转码任务
func (c *MpcClient) CreateTranscodingTaskInvoker(request *model.CreateTranscodingTaskRequest) *CreateTranscodingTaskInvoker {
	requestDef := GenReqDefForCreateTranscodingTask()
	return &CreateTranscodingTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTranscodingTask 取消转码任务
//
// 取消已下发转码任务。
// 只能取消正在转码任务队列中排队的转码任务。已开始转码或已完成的转码任务不能取消。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteTranscodingTask(request *model.DeleteTranscodingTaskRequest) (*model.DeleteTranscodingTaskResponse, error) {
	requestDef := GenReqDefForDeleteTranscodingTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTranscodingTaskResponse), nil
	}
}

// DeleteTranscodingTaskInvoker 取消转码任务
func (c *MpcClient) DeleteTranscodingTaskInvoker(request *model.DeleteTranscodingTaskRequest) *DeleteTranscodingTaskInvoker {
	requestDef := GenReqDefForDeleteTranscodingTask()
	return &DeleteTranscodingTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTranscodingTask 查询转码任务
//
// 查询转码任务状态。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListTranscodingTask(request *model.ListTranscodingTaskRequest) (*model.ListTranscodingTaskResponse, error) {
	requestDef := GenReqDefForListTranscodingTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTranscodingTaskResponse), nil
	}
}

// ListTranscodingTaskInvoker 查询转码任务
func (c *MpcClient) ListTranscodingTaskInvoker(request *model.ListTranscodingTaskRequest) *ListTranscodingTaskInvoker {
	requestDef := GenReqDefForListTranscodingTask()
	return &ListTranscodingTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTransTemplate 新建转码模板
//
// 新建转码模板，采用自定义的模板转码。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateTransTemplate(request *model.CreateTransTemplateRequest) (*model.CreateTransTemplateResponse, error) {
	requestDef := GenReqDefForCreateTransTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTransTemplateResponse), nil
	}
}

// CreateTransTemplateInvoker 新建转码模板
func (c *MpcClient) CreateTransTemplateInvoker(request *model.CreateTransTemplateRequest) *CreateTransTemplateInvoker {
	requestDef := GenReqDefForCreateTransTemplate()
	return &CreateTransTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTemplate 删除转码模板
//
// 删除转码模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteTemplate(request *model.DeleteTemplateRequest) (*model.DeleteTemplateResponse, error) {
	requestDef := GenReqDefForDeleteTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTemplateResponse), nil
	}
}

// DeleteTemplateInvoker 删除转码模板
func (c *MpcClient) DeleteTemplateInvoker(request *model.DeleteTemplateRequest) *DeleteTemplateInvoker {
	requestDef := GenReqDefForDeleteTemplate()
	return &DeleteTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTemplate 查询转码模板
//
// 查询用户自定义转码配置模板。
// 支持指定模板ID查询，或分页全量查询。转码配置模板ID，最多10个。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListTemplate(request *model.ListTemplateRequest) (*model.ListTemplateResponse, error) {
	requestDef := GenReqDefForListTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTemplateResponse), nil
	}
}

// ListTemplateInvoker 查询转码模板
func (c *MpcClient) ListTemplateInvoker(request *model.ListTemplateRequest) *ListTemplateInvoker {
	requestDef := GenReqDefForListTemplate()
	return &ListTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTransTemplate 更新转码模板
//
// 更新转码模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) UpdateTransTemplate(request *model.UpdateTransTemplateRequest) (*model.UpdateTransTemplateResponse, error) {
	requestDef := GenReqDefForUpdateTransTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTransTemplateResponse), nil
	}
}

// UpdateTransTemplateInvoker 更新转码模板
func (c *MpcClient) UpdateTransTemplateInvoker(request *model.UpdateTransTemplateRequest) *UpdateTransTemplateInvoker {
	requestDef := GenReqDefForUpdateTransTemplate()
	return &UpdateTransTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateWatermarkTemplate 新建水印模板
//
// 自定义水印模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) CreateWatermarkTemplate(request *model.CreateWatermarkTemplateRequest) (*model.CreateWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForCreateWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateWatermarkTemplateResponse), nil
	}
}

// CreateWatermarkTemplateInvoker 新建水印模板
func (c *MpcClient) CreateWatermarkTemplateInvoker(request *model.CreateWatermarkTemplateRequest) *CreateWatermarkTemplateInvoker {
	requestDef := GenReqDefForCreateWatermarkTemplate()
	return &CreateWatermarkTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteWatermarkTemplate 删除水印模板
//
// 删除自定义水印模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) DeleteWatermarkTemplate(request *model.DeleteWatermarkTemplateRequest) (*model.DeleteWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForDeleteWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteWatermarkTemplateResponse), nil
	}
}

// DeleteWatermarkTemplateInvoker 删除水印模板
func (c *MpcClient) DeleteWatermarkTemplateInvoker(request *model.DeleteWatermarkTemplateRequest) *DeleteWatermarkTemplateInvoker {
	requestDef := GenReqDefForDeleteWatermarkTemplate()
	return &DeleteWatermarkTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListWatermarkTemplate 查询水印模板
//
// 查询自定义水印模板。支持指定模板ID查询，或分页全量查询。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) ListWatermarkTemplate(request *model.ListWatermarkTemplateRequest) (*model.ListWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForListWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListWatermarkTemplateResponse), nil
	}
}

// ListWatermarkTemplateInvoker 查询水印模板
func (c *MpcClient) ListWatermarkTemplateInvoker(request *model.ListWatermarkTemplateRequest) *ListWatermarkTemplateInvoker {
	requestDef := GenReqDefForListWatermarkTemplate()
	return &ListWatermarkTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateWatermarkTemplate 更新水印模板
//
// 更新自定义水印模板。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *MpcClient) UpdateWatermarkTemplate(request *model.UpdateWatermarkTemplateRequest) (*model.UpdateWatermarkTemplateResponse, error) {
	requestDef := GenReqDefForUpdateWatermarkTemplate()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateWatermarkTemplateResponse), nil
	}
}

// UpdateWatermarkTemplateInvoker 更新水印模板
func (c *MpcClient) UpdateWatermarkTemplateInvoker(request *model.UpdateWatermarkTemplateRequest) *UpdateWatermarkTemplateInvoker {
	requestDef := GenReqDefForUpdateWatermarkTemplate()
	return &UpdateWatermarkTemplateInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
