package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"
)

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

type CreateDomainMappingInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDomainMappingInvoker) Invoke() (*model.CreateDomainMappingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDomainMappingResponse), nil
	}
}

type CreateRecordCallbackConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRecordCallbackConfigInvoker) Invoke() (*model.CreateRecordCallbackConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRecordCallbackConfigResponse), nil
	}
}

type CreateRecordIndexInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRecordIndexInvoker) Invoke() (*model.CreateRecordIndexResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRecordIndexResponse), nil
	}
}

type CreateRecordRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRecordRuleInvoker) Invoke() (*model.CreateRecordRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRecordRuleResponse), nil
	}
}

type CreateStreamForbiddenInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateStreamForbiddenInvoker) Invoke() (*model.CreateStreamForbiddenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateStreamForbiddenResponse), nil
	}
}

type CreateTranscodingsTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTranscodingsTemplateInvoker) Invoke() (*model.CreateTranscodingsTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTranscodingsTemplateResponse), nil
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

type DeleteDomainMappingInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDomainMappingInvoker) Invoke() (*model.DeleteDomainMappingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDomainMappingResponse), nil
	}
}

type DeleteRecordCallbackConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteRecordCallbackConfigInvoker) Invoke() (*model.DeleteRecordCallbackConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteRecordCallbackConfigResponse), nil
	}
}

type DeleteRecordRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteRecordRuleInvoker) Invoke() (*model.DeleteRecordRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteRecordRuleResponse), nil
	}
}

type DeleteStreamForbiddenInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteStreamForbiddenInvoker) Invoke() (*model.DeleteStreamForbiddenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteStreamForbiddenResponse), nil
	}
}

type DeleteTranscodingsTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTranscodingsTemplateInvoker) Invoke() (*model.DeleteTranscodingsTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTranscodingsTemplateResponse), nil
	}
}

type ListLiveSampleLogsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListLiveSampleLogsInvoker) Invoke() (*model.ListLiveSampleLogsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListLiveSampleLogsResponse), nil
	}
}

type ListLiveStreamsOnlineInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListLiveStreamsOnlineInvoker) Invoke() (*model.ListLiveStreamsOnlineResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListLiveStreamsOnlineResponse), nil
	}
}

type ListRecordCallbackConfigsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRecordCallbackConfigsInvoker) Invoke() (*model.ListRecordCallbackConfigsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRecordCallbackConfigsResponse), nil
	}
}

type ListRecordContentsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRecordContentsInvoker) Invoke() (*model.ListRecordContentsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRecordContentsResponse), nil
	}
}

type ListRecordRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRecordRulesInvoker) Invoke() (*model.ListRecordRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRecordRulesResponse), nil
	}
}

type ListStreamForbiddenInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListStreamForbiddenInvoker) Invoke() (*model.ListStreamForbiddenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListStreamForbiddenResponse), nil
	}
}

type RunRecordInvoker struct {
	*invoker.BaseInvoker
}

func (i *RunRecordInvoker) Invoke() (*model.RunRecordResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RunRecordResponse), nil
	}
}

type ShowDomainInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainInvoker) Invoke() (*model.ShowDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainResponse), nil
	}
}

type ShowRecordCallbackConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRecordCallbackConfigInvoker) Invoke() (*model.ShowRecordCallbackConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRecordCallbackConfigResponse), nil
	}
}

type ShowRecordRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRecordRuleInvoker) Invoke() (*model.ShowRecordRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRecordRuleResponse), nil
	}
}

type ShowTranscodingsTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTranscodingsTemplateInvoker) Invoke() (*model.ShowTranscodingsTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTranscodingsTemplateResponse), nil
	}
}

type UpdateDomainInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainInvoker) Invoke() (*model.UpdateDomainResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainResponse), nil
	}
}

type UpdateRecordCallbackConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateRecordCallbackConfigInvoker) Invoke() (*model.UpdateRecordCallbackConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateRecordCallbackConfigResponse), nil
	}
}

type UpdateRecordRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateRecordRuleInvoker) Invoke() (*model.UpdateRecordRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateRecordRuleResponse), nil
	}
}

type UpdateStreamForbiddenInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateStreamForbiddenInvoker) Invoke() (*model.UpdateStreamForbiddenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateStreamForbiddenResponse), nil
	}
}

type UpdateTranscodingsTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTranscodingsTemplateInvoker) Invoke() (*model.UpdateTranscodingsTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTranscodingsTemplateResponse), nil
	}
}
