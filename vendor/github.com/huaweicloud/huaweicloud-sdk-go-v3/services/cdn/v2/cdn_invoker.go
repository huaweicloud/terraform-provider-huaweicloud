package v2

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
)

type BatchCopyDomainInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCopyDomainInvoker) Invoke() (*model.BatchCopyDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCopyDomainResponse), nil
	}
}

type BatchDeleteTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteTagsInvoker) Invoke() (*model.BatchDeleteTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteTagsResponse), nil
	}
}

type CreateDomainInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDomainInvoker) Invoke() (*model.CreateDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDomainResponse), nil
	}
}

type CreatePreheatingTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePreheatingTasksInvoker) Invoke() (*model.CreatePreheatingTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePreheatingTasksResponse), nil
	}
}

type CreateRefreshTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRefreshTasksInvoker) Invoke() (*model.CreateRefreshTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRefreshTasksResponse), nil
	}
}

type CreateTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTagsInvoker) Invoke() (*model.CreateTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTagsResponse), nil
	}
}

type DeleteDomainInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDomainInvoker) Invoke() (*model.DeleteDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDomainResponse), nil
	}
}

type DisableDomainInvoker struct {
	*invoker.BaseInvoker
}

func (i *DisableDomainInvoker) Invoke() (*model.DisableDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DisableDomainResponse), nil
	}
}

type DownloadRegionCarrierExcelInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *DownloadRegionCarrierExcelInvoker) Invoke() (*model.DownloadRegionCarrierExcelResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DownloadRegionCarrierExcelResponse), nil
	}
}

type DownloadStatisticsExcelInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *DownloadStatisticsExcelInvoker) Invoke() (*model.DownloadStatisticsExcelResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DownloadStatisticsExcelResponse), nil
	}
}

type EnableDomainInvoker struct {
	*invoker.BaseInvoker
}

func (i *EnableDomainInvoker) Invoke() (*model.EnableDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.EnableDomainResponse), nil
	}
}

type ListCdnDomainTopRefersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListCdnDomainTopRefersInvoker) Invoke() (*model.ListCdnDomainTopRefersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListCdnDomainTopRefersResponse), nil
	}
}

type ListDomainsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDomainsInvoker) Invoke() (*model.ListDomainsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDomainsResponse), nil
	}
}

type SetChargeModesInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetChargeModesInvoker) Invoke() (*model.SetChargeModesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetChargeModesResponse), nil
	}
}

type ShowBandwidthCalcInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowBandwidthCalcInvoker) Invoke() (*model.ShowBandwidthCalcResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBandwidthCalcResponse), nil
	}
}

type ShowCertificatesHttpsInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowCertificatesHttpsInfoInvoker) Invoke() (*model.ShowCertificatesHttpsInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCertificatesHttpsInfoResponse), nil
	}
}

type ShowChargeModesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowChargeModesInvoker) Invoke() (*model.ShowChargeModesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowChargeModesResponse), nil
	}
}

type ShowDomainDetailByNameInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainDetailByNameInvoker) Invoke() (*model.ShowDomainDetailByNameResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainDetailByNameResponse), nil
	}
}

type ShowDomainFullConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainFullConfigInvoker) Invoke() (*model.ShowDomainFullConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainFullConfigResponse), nil
	}
}

type ShowDomainLocationStatsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainLocationStatsInvoker) Invoke() (*model.ShowDomainLocationStatsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainLocationStatsResponse), nil
	}
}

type ShowDomainStatsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainStatsInvoker) Invoke() (*model.ShowDomainStatsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainStatsResponse), nil
	}
}

type ShowHistoryTaskDetailsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowHistoryTaskDetailsInvoker) Invoke() (*model.ShowHistoryTaskDetailsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowHistoryTaskDetailsResponse), nil
	}
}

type ShowHistoryTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowHistoryTasksInvoker) Invoke() (*model.ShowHistoryTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowHistoryTasksResponse), nil
	}
}

type ShowIpInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowIpInfoInvoker) Invoke() (*model.ShowIpInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowIpInfoResponse), nil
	}
}

type ShowLogsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowLogsInvoker) Invoke() (*model.ShowLogsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowLogsResponse), nil
	}
}

type ShowQuotaInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowQuotaInvoker) Invoke() (*model.ShowQuotaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowQuotaResponse), nil
	}
}

type ShowTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTagsInvoker) Invoke() (*model.ShowTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTagsResponse), nil
	}
}

type ShowTopDomainNamesInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowTopDomainNamesInvoker) Invoke() (*model.ShowTopDomainNamesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTopDomainNamesResponse), nil
	}
}

type ShowTopUrlInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTopUrlInvoker) Invoke() (*model.ShowTopUrlResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTopUrlResponse), nil
	}
}

type ShowUrlTaskInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowUrlTaskInfoInvoker) Invoke() (*model.ShowUrlTaskInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUrlTaskInfoResponse), nil
	}
}

type ShowVerifyDomainOwnerInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowVerifyDomainOwnerInfoInvoker) Invoke() (*model.ShowVerifyDomainOwnerInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowVerifyDomainOwnerInfoResponse), nil
	}
}

type UpdateDomainFullConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainFullConfigInvoker) Invoke() (*model.UpdateDomainFullConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainFullConfigResponse), nil
	}
}

type UpdateDomainMultiCertificatesInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainMultiCertificatesInvoker) Invoke() (*model.UpdateDomainMultiCertificatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainMultiCertificatesResponse), nil
	}
}

type UpdatePrivateBucketAccessInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePrivateBucketAccessInvoker) Invoke() (*model.UpdatePrivateBucketAccessResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePrivateBucketAccessResponse), nil
	}
}

type VerifyDomainOwnerInvoker struct {
	*invoker.BaseInvoker
}

func (i *VerifyDomainOwnerInvoker) Invoke() (*model.VerifyDomainOwnerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.VerifyDomainOwnerResponse), nil
	}
}
