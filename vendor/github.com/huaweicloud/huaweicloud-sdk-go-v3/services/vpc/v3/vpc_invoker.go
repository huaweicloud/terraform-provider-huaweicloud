package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
)

type BatchCreateSubNetworkInterfaceInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreateSubNetworkInterfaceInvoker) Invoke() (*model.BatchCreateSubNetworkInterfaceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreateSubNetworkInterfaceResponse), nil
	}
}

type CreateSecurityGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSecurityGroupInvoker) Invoke() (*model.CreateSecurityGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSecurityGroupResponse), nil
	}
}

type CreateSecurityGroupRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSecurityGroupRuleInvoker) Invoke() (*model.CreateSecurityGroupRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSecurityGroupRuleResponse), nil
	}
}

type CreateSubNetworkInterfaceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSubNetworkInterfaceInvoker) Invoke() (*model.CreateSubNetworkInterfaceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSubNetworkInterfaceResponse), nil
	}
}

type DeleteSecurityGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSecurityGroupInvoker) Invoke() (*model.DeleteSecurityGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSecurityGroupResponse), nil
	}
}

type DeleteSecurityGroupRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSecurityGroupRuleInvoker) Invoke() (*model.DeleteSecurityGroupRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSecurityGroupRuleResponse), nil
	}
}

type DeleteSubNetworkInterfaceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSubNetworkInterfaceInvoker) Invoke() (*model.DeleteSubNetworkInterfaceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSubNetworkInterfaceResponse), nil
	}
}

type ListSecurityGroupRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSecurityGroupRulesInvoker) Invoke() (*model.ListSecurityGroupRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSecurityGroupRulesResponse), nil
	}
}

type ListSecurityGroupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSecurityGroupsInvoker) Invoke() (*model.ListSecurityGroupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSecurityGroupsResponse), nil
	}
}

type ListSubNetworkInterfacesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSubNetworkInterfacesInvoker) Invoke() (*model.ListSubNetworkInterfacesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSubNetworkInterfacesResponse), nil
	}
}

type MigrateSubNetworkInterfaceInvoker struct {
	*invoker.BaseInvoker
}

func (i *MigrateSubNetworkInterfaceInvoker) Invoke() (*model.MigrateSubNetworkInterfaceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.MigrateSubNetworkInterfaceResponse), nil
	}
}

type ShowSecurityGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowSecurityGroupInvoker) Invoke() (*model.ShowSecurityGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowSecurityGroupResponse), nil
	}
}

type ShowSecurityGroupRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowSecurityGroupRuleInvoker) Invoke() (*model.ShowSecurityGroupRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowSecurityGroupRuleResponse), nil
	}
}

type ShowSubNetworkInterfaceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowSubNetworkInterfaceInvoker) Invoke() (*model.ShowSubNetworkInterfaceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowSubNetworkInterfaceResponse), nil
	}
}

type ShowSubNetworkInterfacesQuantityInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowSubNetworkInterfacesQuantityInvoker) Invoke() (*model.ShowSubNetworkInterfacesQuantityResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowSubNetworkInterfacesQuantityResponse), nil
	}
}

type UpdateSecurityGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateSecurityGroupInvoker) Invoke() (*model.UpdateSecurityGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateSecurityGroupResponse), nil
	}
}

type UpdateSubNetworkInterfaceInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateSubNetworkInterfaceInvoker) Invoke() (*model.UpdateSubNetworkInterfaceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateSubNetworkInterfaceResponse), nil
	}
}

type CreateAddressGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAddressGroupInvoker) Invoke() (*model.CreateAddressGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAddressGroupResponse), nil
	}
}

type DeleteAddressGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAddressGroupInvoker) Invoke() (*model.DeleteAddressGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAddressGroupResponse), nil
	}
}

type DeleteIpAddressGroupForceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteIpAddressGroupForceInvoker) Invoke() (*model.DeleteIpAddressGroupForceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteIpAddressGroupForceResponse), nil
	}
}

type ListAddressGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAddressGroupInvoker) Invoke() (*model.ListAddressGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAddressGroupResponse), nil
	}
}

type ShowAddressGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAddressGroupInvoker) Invoke() (*model.ShowAddressGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAddressGroupResponse), nil
	}
}

type UpdateAddressGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAddressGroupInvoker) Invoke() (*model.UpdateAddressGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAddressGroupResponse), nil
	}
}

type AddVpcExtendCidrInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddVpcExtendCidrInvoker) Invoke() (*model.AddVpcExtendCidrResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddVpcExtendCidrResponse), nil
	}
}

type CreateVpcInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateVpcInvoker) Invoke() (*model.CreateVpcResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateVpcResponse), nil
	}
}

type DeleteVpcInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteVpcInvoker) Invoke() (*model.DeleteVpcResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteVpcResponse), nil
	}
}

type ListVpcsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListVpcsInvoker) Invoke() (*model.ListVpcsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListVpcsResponse), nil
	}
}

type RemoveVpcExtendCidrInvoker struct {
	*invoker.BaseInvoker
}

func (i *RemoveVpcExtendCidrInvoker) Invoke() (*model.RemoveVpcExtendCidrResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RemoveVpcExtendCidrResponse), nil
	}
}

type ShowVpcInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowVpcInvoker) Invoke() (*model.ShowVpcResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowVpcResponse), nil
	}
}

type UpdateVpcInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateVpcInvoker) Invoke() (*model.UpdateVpcResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateVpcResponse), nil
	}
}
