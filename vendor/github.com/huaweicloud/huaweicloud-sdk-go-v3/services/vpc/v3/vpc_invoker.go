package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
)

type AddSecurityGroupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddSecurityGroupsInvoker) Invoke() (*model.AddSecurityGroupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddSecurityGroupsResponse), nil
	}
}

type AddSourcesToTrafficMirrorSessionInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddSourcesToTrafficMirrorSessionInvoker) Invoke() (*model.AddSourcesToTrafficMirrorSessionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddSourcesToTrafficMirrorSessionResponse), nil
	}
}

type BatchCreatePortTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreatePortTagsInvoker) Invoke() (*model.BatchCreatePortTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreatePortTagsResponse), nil
	}
}

type BatchCreateSecurityGroupRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreateSecurityGroupRulesInvoker) Invoke() (*model.BatchCreateSecurityGroupRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreateSecurityGroupRulesResponse), nil
	}
}

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

type BatchDeletePortTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeletePortTagsInvoker) Invoke() (*model.BatchDeletePortTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeletePortTagsResponse), nil
	}
}

type CountPortsByTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CountPortsByTagsInvoker) Invoke() (*model.CountPortsByTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CountPortsByTagsResponse), nil
	}
}

type CreatePortTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePortTagInvoker) Invoke() (*model.CreatePortTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePortTagResponse), nil
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

type CreateTrafficMirrorFilterInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTrafficMirrorFilterInvoker) Invoke() (*model.CreateTrafficMirrorFilterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTrafficMirrorFilterResponse), nil
	}
}

type CreateTrafficMirrorFilterRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTrafficMirrorFilterRuleInvoker) Invoke() (*model.CreateTrafficMirrorFilterRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTrafficMirrorFilterRuleResponse), nil
	}
}

type CreateTrafficMirrorSessionInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTrafficMirrorSessionInvoker) Invoke() (*model.CreateTrafficMirrorSessionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTrafficMirrorSessionResponse), nil
	}
}

type DeletePortTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeletePortTagInvoker) Invoke() (*model.DeletePortTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeletePortTagResponse), nil
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

type DeleteTrafficMirrorFilterInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTrafficMirrorFilterInvoker) Invoke() (*model.DeleteTrafficMirrorFilterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTrafficMirrorFilterResponse), nil
	}
}

type DeleteTrafficMirrorFilterRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTrafficMirrorFilterRuleInvoker) Invoke() (*model.DeleteTrafficMirrorFilterRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTrafficMirrorFilterRuleResponse), nil
	}
}

type DeleteTrafficMirrorSessionInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTrafficMirrorSessionInvoker) Invoke() (*model.DeleteTrafficMirrorSessionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTrafficMirrorSessionResponse), nil
	}
}

type ListPortTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPortTagsInvoker) Invoke() (*model.ListPortTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPortTagsResponse), nil
	}
}

type ListPortsByTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPortsByTagsInvoker) Invoke() (*model.ListPortsByTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPortsByTagsResponse), nil
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

type ListTrafficMirrorFilterRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTrafficMirrorFilterRulesInvoker) Invoke() (*model.ListTrafficMirrorFilterRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTrafficMirrorFilterRulesResponse), nil
	}
}

type ListTrafficMirrorFiltersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTrafficMirrorFiltersInvoker) Invoke() (*model.ListTrafficMirrorFiltersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTrafficMirrorFiltersResponse), nil
	}
}

type ListTrafficMirrorSessionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTrafficMirrorSessionsInvoker) Invoke() (*model.ListTrafficMirrorSessionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTrafficMirrorSessionsResponse), nil
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

type RemoveSecurityGroupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *RemoveSecurityGroupsInvoker) Invoke() (*model.RemoveSecurityGroupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RemoveSecurityGroupsResponse), nil
	}
}

type RemoveSourcesFromTrafficMirrorSessionInvoker struct {
	*invoker.BaseInvoker
}

func (i *RemoveSourcesFromTrafficMirrorSessionInvoker) Invoke() (*model.RemoveSourcesFromTrafficMirrorSessionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RemoveSourcesFromTrafficMirrorSessionResponse), nil
	}
}

type ShowPortTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPortTagsInvoker) Invoke() (*model.ShowPortTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPortTagsResponse), nil
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

type ShowTrafficMirrorFilterInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTrafficMirrorFilterInvoker) Invoke() (*model.ShowTrafficMirrorFilterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTrafficMirrorFilterResponse), nil
	}
}

type ShowTrafficMirrorFilterRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTrafficMirrorFilterRuleInvoker) Invoke() (*model.ShowTrafficMirrorFilterRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTrafficMirrorFilterRuleResponse), nil
	}
}

type ShowTrafficMirrorSessionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTrafficMirrorSessionInvoker) Invoke() (*model.ShowTrafficMirrorSessionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTrafficMirrorSessionResponse), nil
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

type UpdateTrafficMirrorFilterInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTrafficMirrorFilterInvoker) Invoke() (*model.UpdateTrafficMirrorFilterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTrafficMirrorFilterResponse), nil
	}
}

type UpdateTrafficMirrorFilterRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTrafficMirrorFilterRuleInvoker) Invoke() (*model.UpdateTrafficMirrorFilterRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTrafficMirrorFilterRuleResponse), nil
	}
}

type UpdateTrafficMirrorSessionInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTrafficMirrorSessionInvoker) Invoke() (*model.UpdateTrafficMirrorSessionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTrafficMirrorSessionResponse), nil
	}
}

type AddFirewallRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddFirewallRulesInvoker) Invoke() (*model.AddFirewallRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddFirewallRulesResponse), nil
	}
}

type AssociateSubnetFirewallInvoker struct {
	*invoker.BaseInvoker
}

func (i *AssociateSubnetFirewallInvoker) Invoke() (*model.AssociateSubnetFirewallResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AssociateSubnetFirewallResponse), nil
	}
}

type BatchCreateFirewallTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreateFirewallTagsInvoker) Invoke() (*model.BatchCreateFirewallTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreateFirewallTagsResponse), nil
	}
}

type BatchDeleteFirewallTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteFirewallTagsInvoker) Invoke() (*model.BatchDeleteFirewallTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteFirewallTagsResponse), nil
	}
}

type CountFirewallsByTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CountFirewallsByTagsInvoker) Invoke() (*model.CountFirewallsByTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CountFirewallsByTagsResponse), nil
	}
}

type CreateFirewallInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateFirewallInvoker) Invoke() (*model.CreateFirewallResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateFirewallResponse), nil
	}
}

type CreateFirewallTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateFirewallTagInvoker) Invoke() (*model.CreateFirewallTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateFirewallTagResponse), nil
	}
}

type DeleteFirewallInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteFirewallInvoker) Invoke() (*model.DeleteFirewallResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteFirewallResponse), nil
	}
}

type DeleteFirewallTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteFirewallTagInvoker) Invoke() (*model.DeleteFirewallTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteFirewallTagResponse), nil
	}
}

type DisassociateSubnetFirewallInvoker struct {
	*invoker.BaseInvoker
}

func (i *DisassociateSubnetFirewallInvoker) Invoke() (*model.DisassociateSubnetFirewallResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DisassociateSubnetFirewallResponse), nil
	}
}

type ListFirewallInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListFirewallInvoker) Invoke() (*model.ListFirewallResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListFirewallResponse), nil
	}
}

type ListFirewallTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListFirewallTagsInvoker) Invoke() (*model.ListFirewallTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListFirewallTagsResponse), nil
	}
}

type ListFirewallsByTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListFirewallsByTagsInvoker) Invoke() (*model.ListFirewallsByTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListFirewallsByTagsResponse), nil
	}
}

type RemoveFirewallRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *RemoveFirewallRulesInvoker) Invoke() (*model.RemoveFirewallRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RemoveFirewallRulesResponse), nil
	}
}

type ShowFirewallInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowFirewallInvoker) Invoke() (*model.ShowFirewallResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowFirewallResponse), nil
	}
}

type ShowFirewallTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowFirewallTagsInvoker) Invoke() (*model.ShowFirewallTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowFirewallTagsResponse), nil
	}
}

type UpdateFirewallInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateFirewallInvoker) Invoke() (*model.UpdateFirewallResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateFirewallResponse), nil
	}
}

type UpdateFirewallRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateFirewallRulesInvoker) Invoke() (*model.UpdateFirewallRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateFirewallRulesResponse), nil
	}
}

type AddClouddcnSubnetsTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddClouddcnSubnetsTagsInvoker) Invoke() (*model.AddClouddcnSubnetsTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddClouddcnSubnetsTagsResponse), nil
	}
}

type BatchCreateClouddcnSubnetsTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreateClouddcnSubnetsTagsInvoker) Invoke() (*model.BatchCreateClouddcnSubnetsTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreateClouddcnSubnetsTagsResponse), nil
	}
}

type BatchDeleteClouddcnSubnetsTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteClouddcnSubnetsTagsInvoker) Invoke() (*model.BatchDeleteClouddcnSubnetsTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteClouddcnSubnetsTagsResponse), nil
	}
}

type CreateClouddcnSubnetInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateClouddcnSubnetInvoker) Invoke() (*model.CreateClouddcnSubnetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateClouddcnSubnetResponse), nil
	}
}

type DeleteClouddcnSubnetInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteClouddcnSubnetInvoker) Invoke() (*model.DeleteClouddcnSubnetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteClouddcnSubnetResponse), nil
	}
}

type DeleteClouddcnSubnetsTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteClouddcnSubnetsTagInvoker) Invoke() (*model.DeleteClouddcnSubnetsTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteClouddcnSubnetsTagResponse), nil
	}
}

type ListClouddcnSubnetsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClouddcnSubnetsInvoker) Invoke() (*model.ListClouddcnSubnetsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClouddcnSubnetsResponse), nil
	}
}

type ListClouddcnSubnetsCountFilterTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClouddcnSubnetsCountFilterTagsInvoker) Invoke() (*model.ListClouddcnSubnetsCountFilterTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClouddcnSubnetsCountFilterTagsResponse), nil
	}
}

type ListClouddcnSubnetsFilterTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClouddcnSubnetsFilterTagsInvoker) Invoke() (*model.ListClouddcnSubnetsFilterTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClouddcnSubnetsFilterTagsResponse), nil
	}
}

type ListClouddcnSubnetsTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClouddcnSubnetsTagsInvoker) Invoke() (*model.ListClouddcnSubnetsTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClouddcnSubnetsTagsResponse), nil
	}
}

type ShowClouddcnSubnetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClouddcnSubnetInvoker) Invoke() (*model.ShowClouddcnSubnetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClouddcnSubnetResponse), nil
	}
}

type ShowClouddcnSubnetsTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClouddcnSubnetsTagsInvoker) Invoke() (*model.ShowClouddcnSubnetsTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClouddcnSubnetsTagsResponse), nil
	}
}

type UpdateClouddcnSubnetInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateClouddcnSubnetInvoker) Invoke() (*model.UpdateClouddcnSubnetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateClouddcnSubnetResponse), nil
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
