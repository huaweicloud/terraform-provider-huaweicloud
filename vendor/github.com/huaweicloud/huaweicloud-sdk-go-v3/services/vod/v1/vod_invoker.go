package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"
)

type CancelAssetTranscodeTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CancelAssetTranscodeTaskInvoker) Invoke() (*model.CancelAssetTranscodeTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CancelAssetTranscodeTaskResponse), nil
	}
}

type CancelExtractAudioTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CancelExtractAudioTaskInvoker) Invoke() (*model.CancelExtractAudioTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CancelExtractAudioTaskResponse), nil
	}
}

type CheckMd5DuplicationInvoker struct {
	*invoker.BaseInvoker
}

func (i *CheckMd5DuplicationInvoker) Invoke() (*model.CheckMd5DuplicationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CheckMd5DuplicationResponse), nil
	}
}

type ConfirmAssetUploadInvoker struct {
	*invoker.BaseInvoker
}

func (i *ConfirmAssetUploadInvoker) Invoke() (*model.ConfirmAssetUploadResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ConfirmAssetUploadResponse), nil
	}
}

type ConfirmImageUploadInvoker struct {
	*invoker.BaseInvoker
}

func (i *ConfirmImageUploadInvoker) Invoke() (*model.ConfirmImageUploadResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ConfirmImageUploadResponse), nil
	}
}

type CreateAssetByFileUploadInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAssetByFileUploadInvoker) Invoke() (*model.CreateAssetByFileUploadResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAssetByFileUploadResponse), nil
	}
}

type CreateAssetCategoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAssetCategoryInvoker) Invoke() (*model.CreateAssetCategoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAssetCategoryResponse), nil
	}
}

type CreateAssetProcessTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAssetProcessTaskInvoker) Invoke() (*model.CreateAssetProcessTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAssetProcessTaskResponse), nil
	}
}

type CreateAssetReviewTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAssetReviewTaskInvoker) Invoke() (*model.CreateAssetReviewTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAssetReviewTaskResponse), nil
	}
}

type CreateExtractAudioTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateExtractAudioTaskInvoker) Invoke() (*model.CreateExtractAudioTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateExtractAudioTaskResponse), nil
	}
}

type CreatePreheatingAssetInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePreheatingAssetInvoker) Invoke() (*model.CreatePreheatingAssetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePreheatingAssetResponse), nil
	}
}

type CreateTakeOverTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTakeOverTaskInvoker) Invoke() (*model.CreateTakeOverTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTakeOverTaskResponse), nil
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

type CreateTemplateGroupCollectionInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTemplateGroupCollectionInvoker) Invoke() (*model.CreateTemplateGroupCollectionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTemplateGroupCollectionResponse), nil
	}
}

type CreateTranscodeTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTranscodeTemplateInvoker) Invoke() (*model.CreateTranscodeTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTranscodeTemplateResponse), nil
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

type DeleteAssetCategoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAssetCategoryInvoker) Invoke() (*model.DeleteAssetCategoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAssetCategoryResponse), nil
	}
}

type DeleteAssetsInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAssetsInvoker) Invoke() (*model.DeleteAssetsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAssetsResponse), nil
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

type DeleteTemplateGroupCollectionInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTemplateGroupCollectionInvoker) Invoke() (*model.DeleteTemplateGroupCollectionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTemplateGroupCollectionResponse), nil
	}
}

type DeleteTranscodeProductInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTranscodeProductInvoker) Invoke() (*model.DeleteTranscodeProductResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTranscodeProductResponse), nil
	}
}

type DeleteTranscodeTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTranscodeTemplateInvoker) Invoke() (*model.DeleteTranscodeTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTranscodeTemplateResponse), nil
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

type ListAssetCategoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAssetCategoryInvoker) Invoke() (*model.ListAssetCategoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAssetCategoryResponse), nil
	}
}

type ListAssetDailySummaryLogInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAssetDailySummaryLogInvoker) Invoke() (*model.ListAssetDailySummaryLogResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAssetDailySummaryLogResponse), nil
	}
}

type ListAssetListInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAssetListInvoker) Invoke() (*model.ListAssetListResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAssetListResponse), nil
	}
}

type ListDomainLogsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDomainLogsInvoker) Invoke() (*model.ListDomainLogsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDomainLogsResponse), nil
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

type ListTemplateGroupCollectionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTemplateGroupCollectionInvoker) Invoke() (*model.ListTemplateGroupCollectionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTemplateGroupCollectionResponse), nil
	}
}

type ListTopStatisticsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTopStatisticsInvoker) Invoke() (*model.ListTopStatisticsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTopStatisticsResponse), nil
	}
}

type ListTranscodeTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTranscodeTemplateInvoker) Invoke() (*model.ListTranscodeTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTranscodeTemplateResponse), nil
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

type ModifySubtitleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ModifySubtitleInvoker) Invoke() (*model.ModifySubtitleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ModifySubtitleResponse), nil
	}
}

type PublishAssetFromObsInvoker struct {
	*invoker.BaseInvoker
}

func (i *PublishAssetFromObsInvoker) Invoke() (*model.PublishAssetFromObsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.PublishAssetFromObsResponse), nil
	}
}

type PublishAssetsInvoker struct {
	*invoker.BaseInvoker
}

func (i *PublishAssetsInvoker) Invoke() (*model.PublishAssetsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.PublishAssetsResponse), nil
	}
}

type ShowAssetCipherInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAssetCipherInvoker) Invoke() (*model.ShowAssetCipherResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAssetCipherResponse), nil
	}
}

type ShowAssetDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAssetDetailInvoker) Invoke() (*model.ShowAssetDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAssetDetailResponse), nil
	}
}

type ShowAssetMetaInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAssetMetaInvoker) Invoke() (*model.ShowAssetMetaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAssetMetaResponse), nil
	}
}

type ShowAssetTempAuthorityInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAssetTempAuthorityInvoker) Invoke() (*model.ShowAssetTempAuthorityResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAssetTempAuthorityResponse), nil
	}
}

type ShowCdnStatisticsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowCdnStatisticsInvoker) Invoke() (*model.ShowCdnStatisticsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCdnStatisticsResponse), nil
	}
}

type ShowPreheatingAssetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPreheatingAssetInvoker) Invoke() (*model.ShowPreheatingAssetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPreheatingAssetResponse), nil
	}
}

type ShowVodRetrievalInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowVodRetrievalInvoker) Invoke() (*model.ShowVodRetrievalResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowVodRetrievalResponse), nil
	}
}

type ShowVodStatisticsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowVodStatisticsInvoker) Invoke() (*model.ShowVodStatisticsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowVodStatisticsResponse), nil
	}
}

type UnpublishAssetsInvoker struct {
	*invoker.BaseInvoker
}

func (i *UnpublishAssetsInvoker) Invoke() (*model.UnpublishAssetsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UnpublishAssetsResponse), nil
	}
}

type UpdateAssetInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAssetInvoker) Invoke() (*model.UpdateAssetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAssetResponse), nil
	}
}

type UpdateAssetCategoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAssetCategoryInvoker) Invoke() (*model.UpdateAssetCategoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAssetCategoryResponse), nil
	}
}

type UpdateAssetMetaInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAssetMetaInvoker) Invoke() (*model.UpdateAssetMetaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAssetMetaResponse), nil
	}
}

type UpdateBucketAuthorizedInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateBucketAuthorizedInvoker) Invoke() (*model.UpdateBucketAuthorizedResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateBucketAuthorizedResponse), nil
	}
}

type UpdateCoverByThumbnailInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateCoverByThumbnailInvoker) Invoke() (*model.UpdateCoverByThumbnailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateCoverByThumbnailResponse), nil
	}
}

type UpdateStorageModeInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateStorageModeInvoker) Invoke() (*model.UpdateStorageModeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateStorageModeResponse), nil
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

type UpdateTemplateGroupCollectionInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTemplateGroupCollectionInvoker) Invoke() (*model.UpdateTemplateGroupCollectionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTemplateGroupCollectionResponse), nil
	}
}

type UpdateTranscodeTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTranscodeTemplateInvoker) Invoke() (*model.UpdateTranscodeTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTranscodeTemplateResponse), nil
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

type UploadMetaDataByUrlInvoker struct {
	*invoker.BaseInvoker
}

func (i *UploadMetaDataByUrlInvoker) Invoke() (*model.UploadMetaDataByUrlResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UploadMetaDataByUrlResponse), nil
	}
}

type ListTakeOverTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTakeOverTaskInvoker) Invoke() (*model.ListTakeOverTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTakeOverTaskResponse), nil
	}
}

type ShowTakeOverAssetDetailsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTakeOverAssetDetailsInvoker) Invoke() (*model.ShowTakeOverAssetDetailsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTakeOverAssetDetailsResponse), nil
	}
}

type ShowTakeOverTaskDetailsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTakeOverTaskDetailsInvoker) Invoke() (*model.ShowTakeOverTaskDetailsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTakeOverTaskDetailsResponse), nil
	}
}
