package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1/model"
)

type BatchDeleteTagsInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *DisableDomainInvoker) Invoke() (*model.DisableDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DisableDomainResponse), nil
	}
}

type EnableDomainInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *EnableDomainInvoker) Invoke() (*model.EnableDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.EnableDomainResponse), nil
	}
}

type ListDomainsInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ListDomainsInvoker) Invoke() (*model.ListDomainsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDomainsResponse), nil
	}
}

type ShowBlackWhiteListInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowBlackWhiteListInvoker) Invoke() (*model.ShowBlackWhiteListResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBlackWhiteListResponse), nil
	}
}

type ShowCacheRulesInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowCacheRulesInvoker) Invoke() (*model.ShowCacheRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCacheRulesResponse), nil
	}
}

type ShowCertificatesHttpsInfoInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowCertificatesHttpsInfoInvoker) Invoke() (*model.ShowCertificatesHttpsInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCertificatesHttpsInfoResponse), nil
	}
}

type ShowDomainDetailInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowDomainDetailInvoker) Invoke() (*model.ShowDomainDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainDetailResponse), nil
	}
}

type ShowDomainFullConfigInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowDomainFullConfigInvoker) Invoke() (*model.ShowDomainFullConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainFullConfigResponse), nil
	}
}

type ShowDomainItemDetailsInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowDomainItemDetailsInvoker) Invoke() (*model.ShowDomainItemDetailsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainItemDetailsResponse), nil
	}
}

type ShowDomainItemLocationDetailsInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowDomainItemLocationDetailsInvoker) Invoke() (*model.ShowDomainItemLocationDetailsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainItemLocationDetailsResponse), nil
	}
}

type ShowDomainLocationStatsInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowHistoryTasksInvoker) Invoke() (*model.ShowHistoryTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowHistoryTasksResponse), nil
	}
}

type ShowHttpInfoInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowHttpInfoInvoker) Invoke() (*model.ShowHttpInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowHttpInfoResponse), nil
	}
}

type ShowIpInfoInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
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

type ShowOriginHostInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowOriginHostInvoker) Invoke() (*model.ShowOriginHostResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowOriginHostResponse), nil
	}
}

type ShowQuotaInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowQuotaInvoker) Invoke() (*model.ShowQuotaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowQuotaResponse), nil
	}
}

type ShowReferInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowReferInvoker) Invoke() (*model.ShowReferResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowReferResponse), nil
	}
}

type ShowResponseHeaderInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowResponseHeaderInvoker) Invoke() (*model.ShowResponseHeaderResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowResponseHeaderResponse), nil
	}
}

type ShowTagsInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowTagsInvoker) Invoke() (*model.ShowTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTagsResponse), nil
	}
}

type ShowTopUrlInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *ShowUrlTaskInfoInvoker) Invoke() (*model.ShowUrlTaskInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUrlTaskInfoResponse), nil
	}
}

type UpdateBlackWhiteListInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateBlackWhiteListInvoker) Invoke() (*model.UpdateBlackWhiteListResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateBlackWhiteListResponse), nil
	}
}

type UpdateCacheRulesInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateCacheRulesInvoker) Invoke() (*model.UpdateCacheRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateCacheRulesResponse), nil
	}
}

type UpdateDomainFullConfigInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
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

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateDomainMultiCertificatesInvoker) Invoke() (*model.UpdateDomainMultiCertificatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainMultiCertificatesResponse), nil
	}
}

type UpdateDomainOriginInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateDomainOriginInvoker) Invoke() (*model.UpdateDomainOriginResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainOriginResponse), nil
	}
}

type UpdateFollow302SwitchInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateFollow302SwitchInvoker) Invoke() (*model.UpdateFollow302SwitchResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateFollow302SwitchResponse), nil
	}
}

type UpdateHttpsInfoInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateHttpsInfoInvoker) Invoke() (*model.UpdateHttpsInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateHttpsInfoResponse), nil
	}
}

type UpdateOriginHostInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateOriginHostInvoker) Invoke() (*model.UpdateOriginHostResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateOriginHostResponse), nil
	}
}

type UpdatePrivateBucketAccessInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdatePrivateBucketAccessInvoker) Invoke() (*model.UpdatePrivateBucketAccessResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePrivateBucketAccessResponse), nil
	}
}

type UpdateRangeSwitchInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateRangeSwitchInvoker) Invoke() (*model.UpdateRangeSwitchResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateRangeSwitchResponse), nil
	}
}

type UpdateReferInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateReferInvoker) Invoke() (*model.UpdateReferResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateReferResponse), nil
	}
}

type UpdateResponseHeaderInvoker struct {
	*invoker.BaseInvoker
}

// Deprecated: This function is deprecated and will be removed in the future versions.
func (i *UpdateResponseHeaderInvoker) Invoke() (*model.UpdateResponseHeaderResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateResponseHeaderResponse), nil
	}
}
