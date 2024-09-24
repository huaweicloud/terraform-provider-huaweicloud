package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
)

type AssociateAgencyWithAllProjectsPermissionInvoker struct {
	*invoker.BaseInvoker
}

func (i *AssociateAgencyWithAllProjectsPermissionInvoker) Invoke() (*model.AssociateAgencyWithAllProjectsPermissionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AssociateAgencyWithAllProjectsPermissionResponse), nil
	}
}

type AssociateAgencyWithDomainPermissionInvoker struct {
	*invoker.BaseInvoker
}

func (i *AssociateAgencyWithDomainPermissionInvoker) Invoke() (*model.AssociateAgencyWithDomainPermissionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AssociateAgencyWithDomainPermissionResponse), nil
	}
}

type AssociateAgencyWithProjectPermissionInvoker struct {
	*invoker.BaseInvoker
}

func (i *AssociateAgencyWithProjectPermissionInvoker) Invoke() (*model.AssociateAgencyWithProjectPermissionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AssociateAgencyWithProjectPermissionResponse), nil
	}
}

type AssociateRoleToAgencyOnEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *AssociateRoleToAgencyOnEnterpriseProjectInvoker) Invoke() (*model.AssociateRoleToAgencyOnEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AssociateRoleToAgencyOnEnterpriseProjectResponse), nil
	}
}

type AssociateRoleToGroupOnEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *AssociateRoleToGroupOnEnterpriseProjectInvoker) Invoke() (*model.AssociateRoleToGroupOnEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AssociateRoleToGroupOnEnterpriseProjectResponse), nil
	}
}

type AssociateRoleToUserOnEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *AssociateRoleToUserOnEnterpriseProjectInvoker) Invoke() (*model.AssociateRoleToUserOnEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AssociateRoleToUserOnEnterpriseProjectResponse), nil
	}
}

type CheckAllProjectsPermissionForAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CheckAllProjectsPermissionForAgencyInvoker) Invoke() (*model.CheckAllProjectsPermissionForAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CheckAllProjectsPermissionForAgencyResponse), nil
	}
}

type CheckDomainPermissionForAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CheckDomainPermissionForAgencyInvoker) Invoke() (*model.CheckDomainPermissionForAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CheckDomainPermissionForAgencyResponse), nil
	}
}

type CheckProjectPermissionForAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CheckProjectPermissionForAgencyInvoker) Invoke() (*model.CheckProjectPermissionForAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CheckProjectPermissionForAgencyResponse), nil
	}
}

type CreateAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAgencyInvoker) Invoke() (*model.CreateAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAgencyResponse), nil
	}
}

type CreateAgencyCustomPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAgencyCustomPolicyInvoker) Invoke() (*model.CreateAgencyCustomPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAgencyCustomPolicyResponse), nil
	}
}

type CreateBindingDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateBindingDeviceInvoker) Invoke() (*model.CreateBindingDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateBindingDeviceResponse), nil
	}
}

type CreateCloudServiceCustomPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateCloudServiceCustomPolicyInvoker) Invoke() (*model.CreateCloudServiceCustomPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateCloudServiceCustomPolicyResponse), nil
	}
}

type CreateLoginTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateLoginTokenInvoker) Invoke() (*model.CreateLoginTokenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateLoginTokenResponse), nil
	}
}

type CreateMetadataInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateMetadataInvoker) Invoke() (*model.CreateMetadataResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateMetadataResponse), nil
	}
}

type CreateMfaDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateMfaDeviceInvoker) Invoke() (*model.CreateMfaDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateMfaDeviceResponse), nil
	}
}

type CreateOpenIdConnectConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateOpenIdConnectConfigInvoker) Invoke() (*model.CreateOpenIdConnectConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateOpenIdConnectConfigResponse), nil
	}
}

type CreateTokenWithIdTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTokenWithIdTokenInvoker) Invoke() (*model.CreateTokenWithIdTokenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTokenWithIdTokenResponse), nil
	}
}

type CreateUnscopedTokenWithIdTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateUnscopedTokenWithIdTokenInvoker) Invoke() (*model.CreateUnscopedTokenWithIdTokenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateUnscopedTokenWithIdTokenResponse), nil
	}
}

type DeleteAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAgencyInvoker) Invoke() (*model.DeleteAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAgencyResponse), nil
	}
}

type DeleteBindingDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteBindingDeviceInvoker) Invoke() (*model.DeleteBindingDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteBindingDeviceResponse), nil
	}
}

type DeleteCustomPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteCustomPolicyInvoker) Invoke() (*model.DeleteCustomPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteCustomPolicyResponse), nil
	}
}

type DeleteDomainGroupInheritedRoleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDomainGroupInheritedRoleInvoker) Invoke() (*model.DeleteDomainGroupInheritedRoleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDomainGroupInheritedRoleResponse), nil
	}
}

type DeleteMfaDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteMfaDeviceInvoker) Invoke() (*model.DeleteMfaDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteMfaDeviceResponse), nil
	}
}

type KeystoneAddUserToGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneAddUserToGroupInvoker) Invoke() (*model.KeystoneAddUserToGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneAddUserToGroupResponse), nil
	}
}

type KeystoneAssociateGroupWithDomainPermissionInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneAssociateGroupWithDomainPermissionInvoker) Invoke() (*model.KeystoneAssociateGroupWithDomainPermissionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneAssociateGroupWithDomainPermissionResponse), nil
	}
}

type KeystoneAssociateGroupWithProjectPermissionInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneAssociateGroupWithProjectPermissionInvoker) Invoke() (*model.KeystoneAssociateGroupWithProjectPermissionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneAssociateGroupWithProjectPermissionResponse), nil
	}
}

type KeystoneCheckDomainPermissionForGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCheckDomainPermissionForGroupInvoker) Invoke() (*model.KeystoneCheckDomainPermissionForGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCheckDomainPermissionForGroupResponse), nil
	}
}

type KeystoneCheckProjectPermissionForGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCheckProjectPermissionForGroupInvoker) Invoke() (*model.KeystoneCheckProjectPermissionForGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCheckProjectPermissionForGroupResponse), nil
	}
}

type KeystoneCheckUserInGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCheckUserInGroupInvoker) Invoke() (*model.KeystoneCheckUserInGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCheckUserInGroupResponse), nil
	}
}

type KeystoneCheckroleForGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCheckroleForGroupInvoker) Invoke() (*model.KeystoneCheckroleForGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCheckroleForGroupResponse), nil
	}
}

type KeystoneCreateGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateGroupInvoker) Invoke() (*model.KeystoneCreateGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateGroupResponse), nil
	}
}

type KeystoneCreateIdentityProviderInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateIdentityProviderInvoker) Invoke() (*model.KeystoneCreateIdentityProviderResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateIdentityProviderResponse), nil
	}
}

type KeystoneCreateMappingInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateMappingInvoker) Invoke() (*model.KeystoneCreateMappingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateMappingResponse), nil
	}
}

type KeystoneCreateProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateProjectInvoker) Invoke() (*model.KeystoneCreateProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateProjectResponse), nil
	}
}

type KeystoneCreateProtocolInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateProtocolInvoker) Invoke() (*model.KeystoneCreateProtocolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateProtocolResponse), nil
	}
}

type KeystoneCreateScopedTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateScopedTokenInvoker) Invoke() (*model.KeystoneCreateScopedTokenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateScopedTokenResponse), nil
	}
}

type KeystoneDeleteGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneDeleteGroupInvoker) Invoke() (*model.KeystoneDeleteGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneDeleteGroupResponse), nil
	}
}

type KeystoneDeleteIdentityProviderInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneDeleteIdentityProviderInvoker) Invoke() (*model.KeystoneDeleteIdentityProviderResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneDeleteIdentityProviderResponse), nil
	}
}

type KeystoneDeleteMappingInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneDeleteMappingInvoker) Invoke() (*model.KeystoneDeleteMappingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneDeleteMappingResponse), nil
	}
}

type KeystoneDeleteProtocolInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneDeleteProtocolInvoker) Invoke() (*model.KeystoneDeleteProtocolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneDeleteProtocolResponse), nil
	}
}

type KeystoneListAllProjectPermissionsForGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListAllProjectPermissionsForGroupInvoker) Invoke() (*model.KeystoneListAllProjectPermissionsForGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListAllProjectPermissionsForGroupResponse), nil
	}
}

type KeystoneListAuthDomainsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListAuthDomainsInvoker) Invoke() (*model.KeystoneListAuthDomainsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListAuthDomainsResponse), nil
	}
}

type KeystoneListAuthProjectsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListAuthProjectsInvoker) Invoke() (*model.KeystoneListAuthProjectsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListAuthProjectsResponse), nil
	}
}

type KeystoneListDomainPermissionsForGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListDomainPermissionsForGroupInvoker) Invoke() (*model.KeystoneListDomainPermissionsForGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListDomainPermissionsForGroupResponse), nil
	}
}

type KeystoneListEndpointsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListEndpointsInvoker) Invoke() (*model.KeystoneListEndpointsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListEndpointsResponse), nil
	}
}

type KeystoneListFederationDomainsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListFederationDomainsInvoker) Invoke() (*model.KeystoneListFederationDomainsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListFederationDomainsResponse), nil
	}
}

type KeystoneListGroupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListGroupsInvoker) Invoke() (*model.KeystoneListGroupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListGroupsResponse), nil
	}
}

type KeystoneListIdentityProvidersInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListIdentityProvidersInvoker) Invoke() (*model.KeystoneListIdentityProvidersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListIdentityProvidersResponse), nil
	}
}

type KeystoneListMappingsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListMappingsInvoker) Invoke() (*model.KeystoneListMappingsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListMappingsResponse), nil
	}
}

type KeystoneListPermissionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListPermissionsInvoker) Invoke() (*model.KeystoneListPermissionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListPermissionsResponse), nil
	}
}

type KeystoneListProjectPermissionsForGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListProjectPermissionsForGroupInvoker) Invoke() (*model.KeystoneListProjectPermissionsForGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListProjectPermissionsForGroupResponse), nil
	}
}

type KeystoneListProjectsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListProjectsInvoker) Invoke() (*model.KeystoneListProjectsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListProjectsResponse), nil
	}
}

type KeystoneListProjectsForUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListProjectsForUserInvoker) Invoke() (*model.KeystoneListProjectsForUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListProjectsForUserResponse), nil
	}
}

type KeystoneListProtocolsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListProtocolsInvoker) Invoke() (*model.KeystoneListProtocolsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListProtocolsResponse), nil
	}
}

type KeystoneListRegionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListRegionsInvoker) Invoke() (*model.KeystoneListRegionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListRegionsResponse), nil
	}
}

type KeystoneListServicesInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListServicesInvoker) Invoke() (*model.KeystoneListServicesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListServicesResponse), nil
	}
}

type KeystoneListVersionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListVersionsInvoker) Invoke() (*model.KeystoneListVersionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListVersionsResponse), nil
	}
}

type KeystoneRemoveDomainPermissionFromGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneRemoveDomainPermissionFromGroupInvoker) Invoke() (*model.KeystoneRemoveDomainPermissionFromGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneRemoveDomainPermissionFromGroupResponse), nil
	}
}

type KeystoneRemoveProjectPermissionFromGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneRemoveProjectPermissionFromGroupInvoker) Invoke() (*model.KeystoneRemoveProjectPermissionFromGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneRemoveProjectPermissionFromGroupResponse), nil
	}
}

type KeystoneRemoveUserFromGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneRemoveUserFromGroupInvoker) Invoke() (*model.KeystoneRemoveUserFromGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneRemoveUserFromGroupResponse), nil
	}
}

type KeystoneShowCatalogInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowCatalogInvoker) Invoke() (*model.KeystoneShowCatalogResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowCatalogResponse), nil
	}
}

type KeystoneShowEndpointInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowEndpointInvoker) Invoke() (*model.KeystoneShowEndpointResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowEndpointResponse), nil
	}
}

type KeystoneShowGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowGroupInvoker) Invoke() (*model.KeystoneShowGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowGroupResponse), nil
	}
}

type KeystoneShowIdentityProviderInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowIdentityProviderInvoker) Invoke() (*model.KeystoneShowIdentityProviderResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowIdentityProviderResponse), nil
	}
}

type KeystoneShowMappingInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowMappingInvoker) Invoke() (*model.KeystoneShowMappingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowMappingResponse), nil
	}
}

type KeystoneShowPermissionInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowPermissionInvoker) Invoke() (*model.KeystoneShowPermissionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowPermissionResponse), nil
	}
}

type KeystoneShowProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowProjectInvoker) Invoke() (*model.KeystoneShowProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowProjectResponse), nil
	}
}

type KeystoneShowProtocolInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowProtocolInvoker) Invoke() (*model.KeystoneShowProtocolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowProtocolResponse), nil
	}
}

type KeystoneShowRegionInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowRegionInvoker) Invoke() (*model.KeystoneShowRegionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowRegionResponse), nil
	}
}

type KeystoneShowSecurityComplianceInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowSecurityComplianceInvoker) Invoke() (*model.KeystoneShowSecurityComplianceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowSecurityComplianceResponse), nil
	}
}

type KeystoneShowSecurityComplianceByOptionInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowSecurityComplianceByOptionInvoker) Invoke() (*model.KeystoneShowSecurityComplianceByOptionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowSecurityComplianceByOptionResponse), nil
	}
}

type KeystoneShowServiceInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowServiceInvoker) Invoke() (*model.KeystoneShowServiceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowServiceResponse), nil
	}
}

type KeystoneShowVersionInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowVersionInvoker) Invoke() (*model.KeystoneShowVersionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowVersionResponse), nil
	}
}

type KeystoneUpdateGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneUpdateGroupInvoker) Invoke() (*model.KeystoneUpdateGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneUpdateGroupResponse), nil
	}
}

type KeystoneUpdateIdentityProviderInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneUpdateIdentityProviderInvoker) Invoke() (*model.KeystoneUpdateIdentityProviderResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneUpdateIdentityProviderResponse), nil
	}
}

type KeystoneUpdateMappingInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneUpdateMappingInvoker) Invoke() (*model.KeystoneUpdateMappingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneUpdateMappingResponse), nil
	}
}

type KeystoneUpdateProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneUpdateProjectInvoker) Invoke() (*model.KeystoneUpdateProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneUpdateProjectResponse), nil
	}
}

type KeystoneUpdateProtocolInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneUpdateProtocolInvoker) Invoke() (*model.KeystoneUpdateProtocolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneUpdateProtocolResponse), nil
	}
}

type ListAgenciesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAgenciesInvoker) Invoke() (*model.ListAgenciesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAgenciesResponse), nil
	}
}

type ListAllProjectsPermissionsForAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAllProjectsPermissionsForAgencyInvoker) Invoke() (*model.ListAllProjectsPermissionsForAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAllProjectsPermissionsForAgencyResponse), nil
	}
}

type ListCustomPoliciesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListCustomPoliciesInvoker) Invoke() (*model.ListCustomPoliciesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListCustomPoliciesResponse), nil
	}
}

type ListDomainPermissionsForAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDomainPermissionsForAgencyInvoker) Invoke() (*model.ListDomainPermissionsForAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDomainPermissionsForAgencyResponse), nil
	}
}

type ListEnterpriseProjectsForGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListEnterpriseProjectsForGroupInvoker) Invoke() (*model.ListEnterpriseProjectsForGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListEnterpriseProjectsForGroupResponse), nil
	}
}

type ListEnterpriseProjectsForUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListEnterpriseProjectsForUserInvoker) Invoke() (*model.ListEnterpriseProjectsForUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListEnterpriseProjectsForUserResponse), nil
	}
}

type ListGroupsForEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListGroupsForEnterpriseProjectInvoker) Invoke() (*model.ListGroupsForEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListGroupsForEnterpriseProjectResponse), nil
	}
}

type ListProjectPermissionsForAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListProjectPermissionsForAgencyInvoker) Invoke() (*model.ListProjectPermissionsForAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListProjectPermissionsForAgencyResponse), nil
	}
}

type ListRolesForGroupOnEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRolesForGroupOnEnterpriseProjectInvoker) Invoke() (*model.ListRolesForGroupOnEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRolesForGroupOnEnterpriseProjectResponse), nil
	}
}

type ListRolesForUserOnEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRolesForUserOnEnterpriseProjectInvoker) Invoke() (*model.ListRolesForUserOnEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRolesForUserOnEnterpriseProjectResponse), nil
	}
}

type ListUsersForEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListUsersForEnterpriseProjectInvoker) Invoke() (*model.ListUsersForEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListUsersForEnterpriseProjectResponse), nil
	}
}

type RemoveAllProjectsPermissionFromAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *RemoveAllProjectsPermissionFromAgencyInvoker) Invoke() (*model.RemoveAllProjectsPermissionFromAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RemoveAllProjectsPermissionFromAgencyResponse), nil
	}
}

type RemoveDomainPermissionFromAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *RemoveDomainPermissionFromAgencyInvoker) Invoke() (*model.RemoveDomainPermissionFromAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RemoveDomainPermissionFromAgencyResponse), nil
	}
}

type RemoveProjectPermissionFromAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *RemoveProjectPermissionFromAgencyInvoker) Invoke() (*model.RemoveProjectPermissionFromAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RemoveProjectPermissionFromAgencyResponse), nil
	}
}

type RevokeRoleFromAgencyOnEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *RevokeRoleFromAgencyOnEnterpriseProjectInvoker) Invoke() (*model.RevokeRoleFromAgencyOnEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RevokeRoleFromAgencyOnEnterpriseProjectResponse), nil
	}
}

type RevokeRoleFromGroupOnEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *RevokeRoleFromGroupOnEnterpriseProjectInvoker) Invoke() (*model.RevokeRoleFromGroupOnEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RevokeRoleFromGroupOnEnterpriseProjectResponse), nil
	}
}

type RevokeRoleFromUserOnEnterpriseProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *RevokeRoleFromUserOnEnterpriseProjectInvoker) Invoke() (*model.RevokeRoleFromUserOnEnterpriseProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RevokeRoleFromUserOnEnterpriseProjectResponse), nil
	}
}

type ShowAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAgencyInvoker) Invoke() (*model.ShowAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAgencyResponse), nil
	}
}

type ShowCustomPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowCustomPolicyInvoker) Invoke() (*model.ShowCustomPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCustomPolicyResponse), nil
	}
}

type ShowDomainApiAclPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainApiAclPolicyInvoker) Invoke() (*model.ShowDomainApiAclPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainApiAclPolicyResponse), nil
	}
}

type ShowDomainConsoleAclPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainConsoleAclPolicyInvoker) Invoke() (*model.ShowDomainConsoleAclPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainConsoleAclPolicyResponse), nil
	}
}

type ShowDomainLoginPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainLoginPolicyInvoker) Invoke() (*model.ShowDomainLoginPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainLoginPolicyResponse), nil
	}
}

type ShowDomainPasswordPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainPasswordPolicyInvoker) Invoke() (*model.ShowDomainPasswordPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainPasswordPolicyResponse), nil
	}
}

type ShowDomainProtectPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainProtectPolicyInvoker) Invoke() (*model.ShowDomainProtectPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainProtectPolicyResponse), nil
	}
}

type ShowDomainQuotaInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainQuotaInvoker) Invoke() (*model.ShowDomainQuotaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainQuotaResponse), nil
	}
}

type ShowDomainRoleAssignmentsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainRoleAssignmentsInvoker) Invoke() (*model.ShowDomainRoleAssignmentsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainRoleAssignmentsResponse), nil
	}
}

type ShowMetadataInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowMetadataInvoker) Invoke() (*model.ShowMetadataResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowMetadataResponse), nil
	}
}

type ShowOpenIdConnectConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowOpenIdConnectConfigInvoker) Invoke() (*model.ShowOpenIdConnectConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowOpenIdConnectConfigResponse), nil
	}
}

type ShowProjectDetailsAndStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowProjectDetailsAndStatusInvoker) Invoke() (*model.ShowProjectDetailsAndStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowProjectDetailsAndStatusResponse), nil
	}
}

type ShowProjectQuotaInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowProjectQuotaInvoker) Invoke() (*model.ShowProjectQuotaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowProjectQuotaResponse), nil
	}
}

type UpdateAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAgencyInvoker) Invoke() (*model.UpdateAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAgencyResponse), nil
	}
}

type UpdateAgencyCustomPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAgencyCustomPolicyInvoker) Invoke() (*model.UpdateAgencyCustomPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAgencyCustomPolicyResponse), nil
	}
}

type UpdateCloudServiceCustomPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateCloudServiceCustomPolicyInvoker) Invoke() (*model.UpdateCloudServiceCustomPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateCloudServiceCustomPolicyResponse), nil
	}
}

type UpdateDomainApiAclPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainApiAclPolicyInvoker) Invoke() (*model.UpdateDomainApiAclPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainApiAclPolicyResponse), nil
	}
}

type UpdateDomainConsoleAclPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainConsoleAclPolicyInvoker) Invoke() (*model.UpdateDomainConsoleAclPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainConsoleAclPolicyResponse), nil
	}
}

type UpdateDomainGroupInheritRoleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainGroupInheritRoleInvoker) Invoke() (*model.UpdateDomainGroupInheritRoleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainGroupInheritRoleResponse), nil
	}
}

type UpdateDomainLoginPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainLoginPolicyInvoker) Invoke() (*model.UpdateDomainLoginPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainLoginPolicyResponse), nil
	}
}

type UpdateDomainPasswordPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainPasswordPolicyInvoker) Invoke() (*model.UpdateDomainPasswordPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainPasswordPolicyResponse), nil
	}
}

type UpdateDomainProtectPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDomainProtectPolicyInvoker) Invoke() (*model.UpdateDomainProtectPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDomainProtectPolicyResponse), nil
	}
}

type UpdateOpenIdConnectConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateOpenIdConnectConfigInvoker) Invoke() (*model.UpdateOpenIdConnectConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateOpenIdConnectConfigResponse), nil
	}
}

type UpdateProjectStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateProjectStatusInvoker) Invoke() (*model.UpdateProjectStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateProjectStatusResponse), nil
	}
}

type CreatePermanentAccessKeyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePermanentAccessKeyInvoker) Invoke() (*model.CreatePermanentAccessKeyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePermanentAccessKeyResponse), nil
	}
}

type CreateTemporaryAccessKeyByAgencyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTemporaryAccessKeyByAgencyInvoker) Invoke() (*model.CreateTemporaryAccessKeyByAgencyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTemporaryAccessKeyByAgencyResponse), nil
	}
}

type CreateTemporaryAccessKeyByTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTemporaryAccessKeyByTokenInvoker) Invoke() (*model.CreateTemporaryAccessKeyByTokenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTemporaryAccessKeyByTokenResponse), nil
	}
}

type DeletePermanentAccessKeyInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeletePermanentAccessKeyInvoker) Invoke() (*model.DeletePermanentAccessKeyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeletePermanentAccessKeyResponse), nil
	}
}

type ListPermanentAccessKeysInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPermanentAccessKeysInvoker) Invoke() (*model.ListPermanentAccessKeysResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPermanentAccessKeysResponse), nil
	}
}

type ShowPermanentAccessKeyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPermanentAccessKeyInvoker) Invoke() (*model.ShowPermanentAccessKeyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPermanentAccessKeyResponse), nil
	}
}

type UpdatePermanentAccessKeyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePermanentAccessKeyInvoker) Invoke() (*model.UpdatePermanentAccessKeyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePermanentAccessKeyResponse), nil
	}
}

type CreateUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateUserInvoker) Invoke() (*model.CreateUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateUserResponse), nil
	}
}

type KeystoneCreateUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateUserInvoker) Invoke() (*model.KeystoneCreateUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateUserResponse), nil
	}
}

type KeystoneDeleteUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneDeleteUserInvoker) Invoke() (*model.KeystoneDeleteUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneDeleteUserResponse), nil
	}
}

type KeystoneListGroupsForUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListGroupsForUserInvoker) Invoke() (*model.KeystoneListGroupsForUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListGroupsForUserResponse), nil
	}
}

type KeystoneListUsersInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListUsersInvoker) Invoke() (*model.KeystoneListUsersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListUsersResponse), nil
	}
}

type KeystoneListUsersForGroupByAdminInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneListUsersForGroupByAdminInvoker) Invoke() (*model.KeystoneListUsersForGroupByAdminResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneListUsersForGroupByAdminResponse), nil
	}
}

type KeystoneShowUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneShowUserInvoker) Invoke() (*model.KeystoneShowUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneShowUserResponse), nil
	}
}

type KeystoneUpdateUserByAdminInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneUpdateUserByAdminInvoker) Invoke() (*model.KeystoneUpdateUserByAdminResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneUpdateUserByAdminResponse), nil
	}
}

type KeystoneUpdateUserPasswordInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneUpdateUserPasswordInvoker) Invoke() (*model.KeystoneUpdateUserPasswordResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneUpdateUserPasswordResponse), nil
	}
}

type ListUserLoginProtectsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListUserLoginProtectsInvoker) Invoke() (*model.ListUserLoginProtectsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListUserLoginProtectsResponse), nil
	}
}

type ListUserMfaDevicesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListUserMfaDevicesInvoker) Invoke() (*model.ListUserMfaDevicesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListUserMfaDevicesResponse), nil
	}
}

type ShowUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowUserInvoker) Invoke() (*model.ShowUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUserResponse), nil
	}
}

type ShowUserLoginProtectInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowUserLoginProtectInvoker) Invoke() (*model.ShowUserLoginProtectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUserLoginProtectResponse), nil
	}
}

type ShowUserMfaDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowUserMfaDeviceInvoker) Invoke() (*model.ShowUserMfaDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUserMfaDeviceResponse), nil
	}
}

type UpdateLoginProtectInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateLoginProtectInvoker) Invoke() (*model.UpdateLoginProtectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateLoginProtectResponse), nil
	}
}

type UpdateUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateUserInvoker) Invoke() (*model.UpdateUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateUserResponse), nil
	}
}

type UpdateUserInformationInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateUserInformationInvoker) Invoke() (*model.UpdateUserInformationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateUserInformationResponse), nil
	}
}

type KeystoneCreateAgencyTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateAgencyTokenInvoker) Invoke() (*model.KeystoneCreateAgencyTokenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateAgencyTokenResponse), nil
	}
}

type KeystoneCreateUserTokenByPasswordInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateUserTokenByPasswordInvoker) Invoke() (*model.KeystoneCreateUserTokenByPasswordResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateUserTokenByPasswordResponse), nil
	}
}

type KeystoneCreateUserTokenByPasswordAndMfaInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneCreateUserTokenByPasswordAndMfaInvoker) Invoke() (*model.KeystoneCreateUserTokenByPasswordAndMfaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneCreateUserTokenByPasswordAndMfaResponse), nil
	}
}

type KeystoneValidateTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *KeystoneValidateTokenInvoker) Invoke() (*model.KeystoneValidateTokenResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeystoneValidateTokenResponse), nil
	}
}
