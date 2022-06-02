package v2

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"
)

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
