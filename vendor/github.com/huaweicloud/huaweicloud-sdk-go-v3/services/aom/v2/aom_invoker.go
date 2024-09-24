package v2

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"
)

type AddActionRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddActionRuleInvoker) Invoke() (*model.AddActionRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddActionRuleResponse), nil
	}
}

type AddAlarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddAlarmRuleInvoker) Invoke() (*model.AddAlarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddAlarmRuleResponse), nil
	}
}

type AddEvent2alarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddEvent2alarmRuleInvoker) Invoke() (*model.AddEvent2alarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddEvent2alarmRuleResponse), nil
	}
}

type AddMetricDataInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddMetricDataInvoker) Invoke() (*model.AddMetricDataResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddMetricDataResponse), nil
	}
}

type AddMuteRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddMuteRulesInvoker) Invoke() (*model.AddMuteRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddMuteRulesResponse), nil
	}
}

type AddOrUpdateMetricOrEventAlarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddOrUpdateMetricOrEventAlarmRuleInvoker) Invoke() (*model.AddOrUpdateMetricOrEventAlarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddOrUpdateMetricOrEventAlarmRuleResponse), nil
	}
}

type AddOrUpdateServiceDiscoveryRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddOrUpdateServiceDiscoveryRulesInvoker) Invoke() (*model.AddOrUpdateServiceDiscoveryRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddOrUpdateServiceDiscoveryRulesResponse), nil
	}
}

type CountEventsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CountEventsInvoker) Invoke() (*model.CountEventsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CountEventsResponse), nil
	}
}

type DeleteActionRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteActionRuleInvoker) Invoke() (*model.DeleteActionRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteActionRuleResponse), nil
	}
}

type DeleteAlarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAlarmRuleInvoker) Invoke() (*model.DeleteAlarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAlarmRuleResponse), nil
	}
}

type DeleteAlarmRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAlarmRulesInvoker) Invoke() (*model.DeleteAlarmRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAlarmRulesResponse), nil
	}
}

type DeleteEvent2alarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteEvent2alarmRuleInvoker) Invoke() (*model.DeleteEvent2alarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteEvent2alarmRuleResponse), nil
	}
}

type DeleteMetricOrEventAlarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteMetricOrEventAlarmRuleInvoker) Invoke() (*model.DeleteMetricOrEventAlarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteMetricOrEventAlarmRuleResponse), nil
	}
}

type DeleteMuteRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteMuteRulesInvoker) Invoke() (*model.DeleteMuteRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteMuteRulesResponse), nil
	}
}

type DeleteserviceDiscoveryRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteserviceDiscoveryRulesInvoker) Invoke() (*model.DeleteserviceDiscoveryRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteserviceDiscoveryRulesResponse), nil
	}
}

type ListActionRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListActionRuleInvoker) Invoke() (*model.ListActionRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListActionRuleResponse), nil
	}
}

type ListAgentsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAgentsInvoker) Invoke() (*model.ListAgentsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAgentsResponse), nil
	}
}

type ListAlarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAlarmRuleInvoker) Invoke() (*model.ListAlarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAlarmRuleResponse), nil
	}
}

type ListEvent2alarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListEvent2alarmRuleInvoker) Invoke() (*model.ListEvent2alarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListEvent2alarmRuleResponse), nil
	}
}

type ListEventsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListEventsInvoker) Invoke() (*model.ListEventsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListEventsResponse), nil
	}
}

type ListLogItemsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListLogItemsInvoker) Invoke() (*model.ListLogItemsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListLogItemsResponse), nil
	}
}

type ListMetricItemsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListMetricItemsInvoker) Invoke() (*model.ListMetricItemsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListMetricItemsResponse), nil
	}
}

type ListMetricOrEventAlarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListMetricOrEventAlarmRuleInvoker) Invoke() (*model.ListMetricOrEventAlarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListMetricOrEventAlarmRuleResponse), nil
	}
}

type ListMuteRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListMuteRuleInvoker) Invoke() (*model.ListMuteRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListMuteRuleResponse), nil
	}
}

type ListNotifiedHistoriesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListNotifiedHistoriesInvoker) Invoke() (*model.ListNotifiedHistoriesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListNotifiedHistoriesResponse), nil
	}
}

type ListPermissionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPermissionsInvoker) Invoke() (*model.ListPermissionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPermissionsResponse), nil
	}
}

type ListSampleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSampleInvoker) Invoke() (*model.ListSampleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSampleResponse), nil
	}
}

type ListSeriesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSeriesInvoker) Invoke() (*model.ListSeriesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSeriesResponse), nil
	}
}

type ListServiceDiscoveryRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListServiceDiscoveryRulesInvoker) Invoke() (*model.ListServiceDiscoveryRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListServiceDiscoveryRulesResponse), nil
	}
}

type PushEventsInvoker struct {
	*invoker.BaseInvoker
}

func (i *PushEventsInvoker) Invoke() (*model.PushEventsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.PushEventsResponse), nil
	}
}

type ShowActionRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowActionRuleInvoker) Invoke() (*model.ShowActionRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowActionRuleResponse), nil
	}
}

type ShowAlarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAlarmRuleInvoker) Invoke() (*model.ShowAlarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAlarmRuleResponse), nil
	}
}

type ShowMetricsDataInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowMetricsDataInvoker) Invoke() (*model.ShowMetricsDataResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowMetricsDataResponse), nil
	}
}

type UpdateActionRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateActionRuleInvoker) Invoke() (*model.UpdateActionRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateActionRuleResponse), nil
	}
}

type UpdateAlarmRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAlarmRuleInvoker) Invoke() (*model.UpdateAlarmRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAlarmRuleResponse), nil
	}
}

type UpdateEventRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateEventRuleInvoker) Invoke() (*model.UpdateEventRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateEventRuleResponse), nil
	}
}

type UpdateMuteRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateMuteRuleInvoker) Invoke() (*model.UpdateMuteRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateMuteRuleResponse), nil
	}
}

type CreatePromInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePromInstanceInvoker) Invoke() (*model.CreatePromInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePromInstanceResponse), nil
	}
}

type CreateRecordingRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRecordingRuleInvoker) Invoke() (*model.CreateRecordingRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRecordingRuleResponse), nil
	}
}

type DeletePromInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeletePromInstanceInvoker) Invoke() (*model.DeletePromInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeletePromInstanceResponse), nil
	}
}

type ListAccessCodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAccessCodeInvoker) Invoke() (*model.ListAccessCodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAccessCodeResponse), nil
	}
}

type ListInstantQueryAomPromGetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstantQueryAomPromGetInvoker) Invoke() (*model.ListInstantQueryAomPromGetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstantQueryAomPromGetResponse), nil
	}
}

type ListInstantQueryAomPromPostInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstantQueryAomPromPostInvoker) Invoke() (*model.ListInstantQueryAomPromPostResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstantQueryAomPromPostResponse), nil
	}
}

type ListLabelValuesAomPromGetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListLabelValuesAomPromGetInvoker) Invoke() (*model.ListLabelValuesAomPromGetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListLabelValuesAomPromGetResponse), nil
	}
}

type ListLabelsAomPromGetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListLabelsAomPromGetInvoker) Invoke() (*model.ListLabelsAomPromGetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListLabelsAomPromGetResponse), nil
	}
}

type ListLabelsAomPromPostInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListLabelsAomPromPostInvoker) Invoke() (*model.ListLabelsAomPromPostResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListLabelsAomPromPostResponse), nil
	}
}

type ListMetadataAomPromGetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListMetadataAomPromGetInvoker) Invoke() (*model.ListMetadataAomPromGetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListMetadataAomPromGetResponse), nil
	}
}

type ListPromInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPromInstanceInvoker) Invoke() (*model.ListPromInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPromInstanceResponse), nil
	}
}

type ListRangeQueryAomPromGetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRangeQueryAomPromGetInvoker) Invoke() (*model.ListRangeQueryAomPromGetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRangeQueryAomPromGetResponse), nil
	}
}

type ListRangeQueryAomPromPostInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRangeQueryAomPromPostInvoker) Invoke() (*model.ListRangeQueryAomPromPostResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRangeQueryAomPromPostResponse), nil
	}
}
