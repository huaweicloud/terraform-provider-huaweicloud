package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1/model"
)

type CreateAnimatedGraphicsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAnimatedGraphicsTaskInvoker) Invoke() (*model.CreateAnimatedGraphicsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAnimatedGraphicsTaskResponse), nil
	}
}

type DeleteAnimatedGraphicsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAnimatedGraphicsTaskInvoker) Invoke() (*model.DeleteAnimatedGraphicsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAnimatedGraphicsTaskResponse), nil
	}
}

type ListAnimatedGraphicsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAnimatedGraphicsTaskInvoker) Invoke() (*model.ListAnimatedGraphicsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAnimatedGraphicsTaskResponse), nil
	}
}

type CreateEditingJobInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateEditingJobInvoker) Invoke() (*model.CreateEditingJobResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateEditingJobResponse), nil
	}
}

type DeleteEditingJobInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteEditingJobInvoker) Invoke() (*model.DeleteEditingJobResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteEditingJobResponse), nil
	}
}

type ListEditingJobInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListEditingJobInvoker) Invoke() (*model.ListEditingJobResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListEditingJobResponse), nil
	}
}

type CreateEncryptTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateEncryptTaskInvoker) Invoke() (*model.CreateEncryptTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateEncryptTaskResponse), nil
	}
}

type DeleteEncryptTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteEncryptTaskInvoker) Invoke() (*model.DeleteEncryptTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteEncryptTaskResponse), nil
	}
}

type ListEncryptTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListEncryptTaskInvoker) Invoke() (*model.ListEncryptTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListEncryptTaskResponse), nil
	}
}

type CreateExtractTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateExtractTaskInvoker) Invoke() (*model.CreateExtractTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateExtractTaskResponse), nil
	}
}

type DeleteExtractTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteExtractTaskInvoker) Invoke() (*model.DeleteExtractTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteExtractTaskResponse), nil
	}
}

type ListExtractTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListExtractTaskInvoker) Invoke() (*model.ListExtractTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListExtractTaskResponse), nil
	}
}

type CreateMbTasksReportInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateMbTasksReportInvoker) Invoke() (*model.CreateMbTasksReportResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateMbTasksReportResponse), nil
	}
}

type CreateMergeChannelsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateMergeChannelsTaskInvoker) Invoke() (*model.CreateMergeChannelsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateMergeChannelsTaskResponse), nil
	}
}

type CreateResetTracksTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateResetTracksTaskInvoker) Invoke() (*model.CreateResetTracksTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateResetTracksTaskResponse), nil
	}
}

type DeleteMergeChannelsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteMergeChannelsTaskInvoker) Invoke() (*model.DeleteMergeChannelsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteMergeChannelsTaskResponse), nil
	}
}

type DeleteResetTracksTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteResetTracksTaskInvoker) Invoke() (*model.DeleteResetTracksTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteResetTracksTaskResponse), nil
	}
}

type ListMergeChannelsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListMergeChannelsTaskInvoker) Invoke() (*model.ListMergeChannelsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListMergeChannelsTaskResponse), nil
	}
}

type ListResetTracksTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListResetTracksTaskInvoker) Invoke() (*model.ListResetTracksTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListResetTracksTaskResponse), nil
	}
}

type CreateMediaProcessTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateMediaProcessTaskInvoker) Invoke() (*model.CreateMediaProcessTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateMediaProcessTaskResponse), nil
	}
}

type DeleteMediaProcessTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteMediaProcessTaskInvoker) Invoke() (*model.DeleteMediaProcessTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteMediaProcessTaskResponse), nil
	}
}

type ListMediaProcessTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListMediaProcessTaskInvoker) Invoke() (*model.ListMediaProcessTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListMediaProcessTaskResponse), nil
	}
}

type CreateMpeCallBackInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateMpeCallBackInvoker) Invoke() (*model.CreateMpeCallBackResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateMpeCallBackResponse), nil
	}
}

type CreateQualityEnhanceTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateQualityEnhanceTemplateInvoker) Invoke() (*model.CreateQualityEnhanceTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateQualityEnhanceTemplateResponse), nil
	}
}

type DeleteQualityEnhanceTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteQualityEnhanceTemplateInvoker) Invoke() (*model.DeleteQualityEnhanceTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteQualityEnhanceTemplateResponse), nil
	}
}

type ListQualityEnhanceDefaultTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListQualityEnhanceDefaultTemplateInvoker) Invoke() (*model.ListQualityEnhanceDefaultTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListQualityEnhanceDefaultTemplateResponse), nil
	}
}

type UpdateQualityEnhanceTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateQualityEnhanceTemplateInvoker) Invoke() (*model.UpdateQualityEnhanceTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateQualityEnhanceTemplateResponse), nil
	}
}

type ListTranscodeDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTranscodeDetailInvoker) Invoke() (*model.ListTranscodeDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTranscodeDetailResponse), nil
	}
}

type CancelRemuxTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CancelRemuxTaskInvoker) Invoke() (*model.CancelRemuxTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CancelRemuxTaskResponse), nil
	}
}

type CreateRemuxTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRemuxTaskInvoker) Invoke() (*model.CreateRemuxTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRemuxTaskResponse), nil
	}
}

type CreateRetryRemuxTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRetryRemuxTaskInvoker) Invoke() (*model.CreateRetryRemuxTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRetryRemuxTaskResponse), nil
	}
}

type DeleteRemuxTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteRemuxTaskInvoker) Invoke() (*model.DeleteRemuxTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteRemuxTaskResponse), nil
	}
}

type ListRemuxTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRemuxTaskInvoker) Invoke() (*model.ListRemuxTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRemuxTaskResponse), nil
	}
}

type CreateTemplateGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTemplateGroupInvoker) Invoke() (*model.CreateTemplateGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTemplateGroupResponse), nil
	}
}

type DeleteTemplateGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTemplateGroupInvoker) Invoke() (*model.DeleteTemplateGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTemplateGroupResponse), nil
	}
}

type ListTemplateGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTemplateGroupInvoker) Invoke() (*model.ListTemplateGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTemplateGroupResponse), nil
	}
}

type UpdateTemplateGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTemplateGroupInvoker) Invoke() (*model.UpdateTemplateGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTemplateGroupResponse), nil
	}
}

type CreateThumbnailsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateThumbnailsTaskInvoker) Invoke() (*model.CreateThumbnailsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateThumbnailsTaskResponse), nil
	}
}

type DeleteThumbnailsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteThumbnailsTaskInvoker) Invoke() (*model.DeleteThumbnailsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteThumbnailsTaskResponse), nil
	}
}

type ListThumbnailsTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListThumbnailsTaskInvoker) Invoke() (*model.ListThumbnailsTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListThumbnailsTaskResponse), nil
	}
}

type CreateTranscodingTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTranscodingTaskInvoker) Invoke() (*model.CreateTranscodingTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTranscodingTaskResponse), nil
	}
}

type DeleteTranscodingTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTranscodingTaskInvoker) Invoke() (*model.DeleteTranscodingTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTranscodingTaskResponse), nil
	}
}

type ListTranscodingTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTranscodingTaskInvoker) Invoke() (*model.ListTranscodingTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTranscodingTaskResponse), nil
	}
}

type CreateTransTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTransTemplateInvoker) Invoke() (*model.CreateTransTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTransTemplateResponse), nil
	}
}

type DeleteTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTemplateInvoker) Invoke() (*model.DeleteTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTemplateResponse), nil
	}
}

type ListTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTemplateInvoker) Invoke() (*model.ListTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTemplateResponse), nil
	}
}

type UpdateTransTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTransTemplateInvoker) Invoke() (*model.UpdateTransTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTransTemplateResponse), nil
	}
}

type CreateWatermarkTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateWatermarkTemplateInvoker) Invoke() (*model.CreateWatermarkTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateWatermarkTemplateResponse), nil
	}
}

type DeleteWatermarkTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteWatermarkTemplateInvoker) Invoke() (*model.DeleteWatermarkTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteWatermarkTemplateResponse), nil
	}
}

type ListWatermarkTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListWatermarkTemplateInvoker) Invoke() (*model.ListWatermarkTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListWatermarkTemplateResponse), nil
	}
}

type UpdateWatermarkTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateWatermarkTemplateInvoker) Invoke() (*model.UpdateWatermarkTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateWatermarkTemplateResponse), nil
	}
}
