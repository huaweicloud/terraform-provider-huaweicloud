package v3

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
)

type IamClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewIamClient(hcClient *httpclient.HcHttpClient) *IamClient {
	return &IamClient{HcClient: hcClient}
}

func IamClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder().WithCredentialsType("global.Credentials,basic.Credentials,v3.IamCredentials")
	return builder
}

// AssociateAgencyWithAllProjectsPermission 为委托授予所有项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)为委托授予所有项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) AssociateAgencyWithAllProjectsPermission(request *model.AssociateAgencyWithAllProjectsPermissionRequest) (*model.AssociateAgencyWithAllProjectsPermissionResponse, error) {
	requestDef := GenReqDefForAssociateAgencyWithAllProjectsPermission()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociateAgencyWithAllProjectsPermissionResponse), nil
	}
}

// AssociateAgencyWithAllProjectsPermissionInvoker 为委托授予所有项目服务权限
func (c *IamClient) AssociateAgencyWithAllProjectsPermissionInvoker(request *model.AssociateAgencyWithAllProjectsPermissionRequest) *AssociateAgencyWithAllProjectsPermissionInvoker {
	requestDef := GenReqDefForAssociateAgencyWithAllProjectsPermission()
	return &AssociateAgencyWithAllProjectsPermissionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AssociateAgencyWithDomainPermission 为委托授予全局服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)为委托授予全局服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) AssociateAgencyWithDomainPermission(request *model.AssociateAgencyWithDomainPermissionRequest) (*model.AssociateAgencyWithDomainPermissionResponse, error) {
	requestDef := GenReqDefForAssociateAgencyWithDomainPermission()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociateAgencyWithDomainPermissionResponse), nil
	}
}

// AssociateAgencyWithDomainPermissionInvoker 为委托授予全局服务权限
func (c *IamClient) AssociateAgencyWithDomainPermissionInvoker(request *model.AssociateAgencyWithDomainPermissionRequest) *AssociateAgencyWithDomainPermissionInvoker {
	requestDef := GenReqDefForAssociateAgencyWithDomainPermission()
	return &AssociateAgencyWithDomainPermissionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AssociateAgencyWithProjectPermission 为委托授予项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)为委托授予项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) AssociateAgencyWithProjectPermission(request *model.AssociateAgencyWithProjectPermissionRequest) (*model.AssociateAgencyWithProjectPermissionResponse, error) {
	requestDef := GenReqDefForAssociateAgencyWithProjectPermission()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociateAgencyWithProjectPermissionResponse), nil
	}
}

// AssociateAgencyWithProjectPermissionInvoker 为委托授予项目服务权限
func (c *IamClient) AssociateAgencyWithProjectPermissionInvoker(request *model.AssociateAgencyWithProjectPermissionRequest) *AssociateAgencyWithProjectPermissionInvoker {
	requestDef := GenReqDefForAssociateAgencyWithProjectPermission()
	return &AssociateAgencyWithProjectPermissionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AssociateRoleToAgencyOnEnterpriseProject 基于委托为企业项目授权
//
// 该接口可以基于委托为企业项目授权
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) AssociateRoleToAgencyOnEnterpriseProject(request *model.AssociateRoleToAgencyOnEnterpriseProjectRequest) (*model.AssociateRoleToAgencyOnEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForAssociateRoleToAgencyOnEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociateRoleToAgencyOnEnterpriseProjectResponse), nil
	}
}

// AssociateRoleToAgencyOnEnterpriseProjectInvoker 基于委托为企业项目授权
func (c *IamClient) AssociateRoleToAgencyOnEnterpriseProjectInvoker(request *model.AssociateRoleToAgencyOnEnterpriseProjectRequest) *AssociateRoleToAgencyOnEnterpriseProjectInvoker {
	requestDef := GenReqDefForAssociateRoleToAgencyOnEnterpriseProject()
	return &AssociateRoleToAgencyOnEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AssociateRoleToGroupOnEnterpriseProject 基于用户组为企业项目授权
//
// 该接口用于基于用户组为企业项目授权。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) AssociateRoleToGroupOnEnterpriseProject(request *model.AssociateRoleToGroupOnEnterpriseProjectRequest) (*model.AssociateRoleToGroupOnEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForAssociateRoleToGroupOnEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociateRoleToGroupOnEnterpriseProjectResponse), nil
	}
}

// AssociateRoleToGroupOnEnterpriseProjectInvoker 基于用户组为企业项目授权
func (c *IamClient) AssociateRoleToGroupOnEnterpriseProjectInvoker(request *model.AssociateRoleToGroupOnEnterpriseProjectRequest) *AssociateRoleToGroupOnEnterpriseProjectInvoker {
	requestDef := GenReqDefForAssociateRoleToGroupOnEnterpriseProject()
	return &AssociateRoleToGroupOnEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AssociateRoleToUserOnEnterpriseProject 基于用户为企业项目授权
//
// 基于用户为企业项目授权。
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) AssociateRoleToUserOnEnterpriseProject(request *model.AssociateRoleToUserOnEnterpriseProjectRequest) (*model.AssociateRoleToUserOnEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForAssociateRoleToUserOnEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociateRoleToUserOnEnterpriseProjectResponse), nil
	}
}

// AssociateRoleToUserOnEnterpriseProjectInvoker 基于用户为企业项目授权
func (c *IamClient) AssociateRoleToUserOnEnterpriseProjectInvoker(request *model.AssociateRoleToUserOnEnterpriseProjectRequest) *AssociateRoleToUserOnEnterpriseProjectInvoker {
	requestDef := GenReqDefForAssociateRoleToUserOnEnterpriseProject()
	return &AssociateRoleToUserOnEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CheckAllProjectsPermissionForAgency 检查委托下是否具有所有项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)检查委托是否具有所有项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CheckAllProjectsPermissionForAgency(request *model.CheckAllProjectsPermissionForAgencyRequest) (*model.CheckAllProjectsPermissionForAgencyResponse, error) {
	requestDef := GenReqDefForCheckAllProjectsPermissionForAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CheckAllProjectsPermissionForAgencyResponse), nil
	}
}

// CheckAllProjectsPermissionForAgencyInvoker 检查委托下是否具有所有项目服务权限
func (c *IamClient) CheckAllProjectsPermissionForAgencyInvoker(request *model.CheckAllProjectsPermissionForAgencyRequest) *CheckAllProjectsPermissionForAgencyInvoker {
	requestDef := GenReqDefForCheckAllProjectsPermissionForAgency()
	return &CheckAllProjectsPermissionForAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CheckDomainPermissionForAgency 查询委托是否拥有全局服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询委托是否拥有全局服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CheckDomainPermissionForAgency(request *model.CheckDomainPermissionForAgencyRequest) (*model.CheckDomainPermissionForAgencyResponse, error) {
	requestDef := GenReqDefForCheckDomainPermissionForAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CheckDomainPermissionForAgencyResponse), nil
	}
}

// CheckDomainPermissionForAgencyInvoker 查询委托是否拥有全局服务权限
func (c *IamClient) CheckDomainPermissionForAgencyInvoker(request *model.CheckDomainPermissionForAgencyRequest) *CheckDomainPermissionForAgencyInvoker {
	requestDef := GenReqDefForCheckDomainPermissionForAgency()
	return &CheckDomainPermissionForAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CheckProjectPermissionForAgency 查询委托是否拥有项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询委托是否拥有项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CheckProjectPermissionForAgency(request *model.CheckProjectPermissionForAgencyRequest) (*model.CheckProjectPermissionForAgencyResponse, error) {
	requestDef := GenReqDefForCheckProjectPermissionForAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CheckProjectPermissionForAgencyResponse), nil
	}
}

// CheckProjectPermissionForAgencyInvoker 查询委托是否拥有项目服务权限
func (c *IamClient) CheckProjectPermissionForAgencyInvoker(request *model.CheckProjectPermissionForAgencyRequest) *CheckProjectPermissionForAgencyInvoker {
	requestDef := GenReqDefForCheckProjectPermissionForAgency()
	return &CheckProjectPermissionForAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAgency 创建委托
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)创建委托。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateAgency(request *model.CreateAgencyRequest) (*model.CreateAgencyResponse, error) {
	requestDef := GenReqDefForCreateAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAgencyResponse), nil
	}
}

// CreateAgencyInvoker 创建委托
func (c *IamClient) CreateAgencyInvoker(request *model.CreateAgencyRequest) *CreateAgencyInvoker {
	requestDef := GenReqDefForCreateAgency()
	return &CreateAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAgencyCustomPolicy 创建委托自定义策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)创建委托自定义策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateAgencyCustomPolicy(request *model.CreateAgencyCustomPolicyRequest) (*model.CreateAgencyCustomPolicyResponse, error) {
	requestDef := GenReqDefForCreateAgencyCustomPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAgencyCustomPolicyResponse), nil
	}
}

// CreateAgencyCustomPolicyInvoker 创建委托自定义策略
func (c *IamClient) CreateAgencyCustomPolicyInvoker(request *model.CreateAgencyCustomPolicyRequest) *CreateAgencyCustomPolicyInvoker {
	requestDef := GenReqDefForCreateAgencyCustomPolicy()
	return &CreateAgencyCustomPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateBindingDevice 绑定MFA设备
//
// 该接口可以用于绑定MFA设备。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateBindingDevice(request *model.CreateBindingDeviceRequest) (*model.CreateBindingDeviceResponse, error) {
	requestDef := GenReqDefForCreateBindingDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateBindingDeviceResponse), nil
	}
}

// CreateBindingDeviceInvoker 绑定MFA设备
func (c *IamClient) CreateBindingDeviceInvoker(request *model.CreateBindingDeviceRequest) *CreateBindingDeviceInvoker {
	requestDef := GenReqDefForCreateBindingDevice()
	return &CreateBindingDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateCloudServiceCustomPolicy 创建云服务自定义策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)创建云服务自定义策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateCloudServiceCustomPolicy(request *model.CreateCloudServiceCustomPolicyRequest) (*model.CreateCloudServiceCustomPolicyResponse, error) {
	requestDef := GenReqDefForCreateCloudServiceCustomPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateCloudServiceCustomPolicyResponse), nil
	}
}

// CreateCloudServiceCustomPolicyInvoker 创建云服务自定义策略
func (c *IamClient) CreateCloudServiceCustomPolicyInvoker(request *model.CreateCloudServiceCustomPolicyRequest) *CreateCloudServiceCustomPolicyInvoker {
	requestDef := GenReqDefForCreateCloudServiceCustomPolicy()
	return &CreateCloudServiceCustomPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateLoginToken 获取自定义代理登录票据
//
// 该接口用于用于获取自定义代理登录票据logintoken。logintoken是系统颁发给自定义代理用户的登录票据，承载用户的身份、session等信息。调用自定义代理URL登录云服务控制台时，可以使用本接口获取的logintoken进行认证。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// &gt; - logintoken的有效期为10分钟。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateLoginToken(request *model.CreateLoginTokenRequest) (*model.CreateLoginTokenResponse, error) {
	requestDef := GenReqDefForCreateLoginToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateLoginTokenResponse), nil
	}
}

// CreateLoginTokenInvoker 获取自定义代理登录票据
func (c *IamClient) CreateLoginTokenInvoker(request *model.CreateLoginTokenRequest) *CreateLoginTokenInvoker {
	requestDef := GenReqDefForCreateLoginToken()
	return &CreateLoginTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateMetadata 导入Metadata文件
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)导入Metadata文件。
//
// 账号在使用联邦认证功能前，需要先将Metadata文件导入到IAM中。Metadata文件是SAML 2.0协议约定的接口文件，包含访问接口地址和证书信息，请找企业管理员获取企业IdP的Metadata文件。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateMetadata(request *model.CreateMetadataRequest) (*model.CreateMetadataResponse, error) {
	requestDef := GenReqDefForCreateMetadata()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateMetadataResponse), nil
	}
}

// CreateMetadataInvoker 导入Metadata文件
func (c *IamClient) CreateMetadataInvoker(request *model.CreateMetadataRequest) *CreateMetadataInvoker {
	requestDef := GenReqDefForCreateMetadata()
	return &CreateMetadataInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateMfaDevice 创建MFA设备
//
// 该接口可以用于创建MFA设备。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateMfaDevice(request *model.CreateMfaDeviceRequest) (*model.CreateMfaDeviceResponse, error) {
	requestDef := GenReqDefForCreateMfaDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateMfaDeviceResponse), nil
	}
}

// CreateMfaDeviceInvoker 创建MFA设备
func (c *IamClient) CreateMfaDeviceInvoker(request *model.CreateMfaDeviceRequest) *CreateMfaDeviceInvoker {
	requestDef := GenReqDefForCreateMfaDevice()
	return &CreateMfaDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateOpenIdConnectConfig 创建OpenId Connect身份提供商配置
//
// 创建OpenId Connect身份提供商配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateOpenIdConnectConfig(request *model.CreateOpenIdConnectConfigRequest) (*model.CreateOpenIdConnectConfigResponse, error) {
	requestDef := GenReqDefForCreateOpenIdConnectConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateOpenIdConnectConfigResponse), nil
	}
}

// CreateOpenIdConnectConfigInvoker 创建OpenId Connect身份提供商配置
func (c *IamClient) CreateOpenIdConnectConfigInvoker(request *model.CreateOpenIdConnectConfigRequest) *CreateOpenIdConnectConfigInvoker {
	requestDef := GenReqDefForCreateOpenIdConnectConfig()
	return &CreateOpenIdConnectConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTokenWithIdToken 获取联邦认证token(OpenId Connect Id token方式)
//
// 获取联邦认证token(OpenId Connect Id token方式)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateTokenWithIdToken(request *model.CreateTokenWithIdTokenRequest) (*model.CreateTokenWithIdTokenResponse, error) {
	requestDef := GenReqDefForCreateTokenWithIdToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTokenWithIdTokenResponse), nil
	}
}

// CreateTokenWithIdTokenInvoker 获取联邦认证token(OpenId Connect Id token方式)
func (c *IamClient) CreateTokenWithIdTokenInvoker(request *model.CreateTokenWithIdTokenRequest) *CreateTokenWithIdTokenInvoker {
	requestDef := GenReqDefForCreateTokenWithIdToken()
	return &CreateTokenWithIdTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateUnscopedTokenWithIdToken 获取联邦认证unscoped token(OpenId Connect Id token方式)
//
// 获取联邦认证token(OpenId Connect Id token方式)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateUnscopedTokenWithIdToken(request *model.CreateUnscopedTokenWithIdTokenRequest) (*model.CreateUnscopedTokenWithIdTokenResponse, error) {
	requestDef := GenReqDefForCreateUnscopedTokenWithIdToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateUnscopedTokenWithIdTokenResponse), nil
	}
}

// CreateUnscopedTokenWithIdTokenInvoker 获取联邦认证unscoped token(OpenId Connect Id token方式)
func (c *IamClient) CreateUnscopedTokenWithIdTokenInvoker(request *model.CreateUnscopedTokenWithIdTokenRequest) *CreateUnscopedTokenWithIdTokenInvoker {
	requestDef := GenReqDefForCreateUnscopedTokenWithIdToken()
	return &CreateUnscopedTokenWithIdTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAgency 删除委托
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)删除委托。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) DeleteAgency(request *model.DeleteAgencyRequest) (*model.DeleteAgencyResponse, error) {
	requestDef := GenReqDefForDeleteAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAgencyResponse), nil
	}
}

// DeleteAgencyInvoker 删除委托
func (c *IamClient) DeleteAgencyInvoker(request *model.DeleteAgencyRequest) *DeleteAgencyInvoker {
	requestDef := GenReqDefForDeleteAgency()
	return &DeleteAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteBindingDevice 解绑MFA设备
//
// 该接口可以用于解绑MFA设备
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) DeleteBindingDevice(request *model.DeleteBindingDeviceRequest) (*model.DeleteBindingDeviceResponse, error) {
	requestDef := GenReqDefForDeleteBindingDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteBindingDeviceResponse), nil
	}
}

// DeleteBindingDeviceInvoker 解绑MFA设备
func (c *IamClient) DeleteBindingDeviceInvoker(request *model.DeleteBindingDeviceRequest) *DeleteBindingDeviceInvoker {
	requestDef := GenReqDefForDeleteBindingDevice()
	return &DeleteBindingDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteCustomPolicy 删除自定义策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)删除自定义策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) DeleteCustomPolicy(request *model.DeleteCustomPolicyRequest) (*model.DeleteCustomPolicyResponse, error) {
	requestDef := GenReqDefForDeleteCustomPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteCustomPolicyResponse), nil
	}
}

// DeleteCustomPolicyInvoker 删除自定义策略
func (c *IamClient) DeleteCustomPolicyInvoker(request *model.DeleteCustomPolicyRequest) *DeleteCustomPolicyInvoker {
	requestDef := GenReqDefForDeleteCustomPolicy()
	return &DeleteCustomPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDomainGroupInheritedRole 移除用户组的所有项目服务权限
//
// 该接口可以用于移除用户组的所有项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) DeleteDomainGroupInheritedRole(request *model.DeleteDomainGroupInheritedRoleRequest) (*model.DeleteDomainGroupInheritedRoleResponse, error) {
	requestDef := GenReqDefForDeleteDomainGroupInheritedRole()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDomainGroupInheritedRoleResponse), nil
	}
}

// DeleteDomainGroupInheritedRoleInvoker 移除用户组的所有项目服务权限
func (c *IamClient) DeleteDomainGroupInheritedRoleInvoker(request *model.DeleteDomainGroupInheritedRoleRequest) *DeleteDomainGroupInheritedRoleInvoker {
	requestDef := GenReqDefForDeleteDomainGroupInheritedRole()
	return &DeleteDomainGroupInheritedRoleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteMfaDevice 删除MFA设备
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)删除MFA设备。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) DeleteMfaDevice(request *model.DeleteMfaDeviceRequest) (*model.DeleteMfaDeviceResponse, error) {
	requestDef := GenReqDefForDeleteMfaDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteMfaDeviceResponse), nil
	}
}

// DeleteMfaDeviceInvoker 删除MFA设备
func (c *IamClient) DeleteMfaDeviceInvoker(request *model.DeleteMfaDeviceRequest) *DeleteMfaDeviceInvoker {
	requestDef := GenReqDefForDeleteMfaDevice()
	return &DeleteMfaDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneAddUserToGroup 添加IAM用户到用户组
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)添加IAM用户到用户组。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneAddUserToGroup(request *model.KeystoneAddUserToGroupRequest) (*model.KeystoneAddUserToGroupResponse, error) {
	requestDef := GenReqDefForKeystoneAddUserToGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneAddUserToGroupResponse), nil
	}
}

// KeystoneAddUserToGroupInvoker 添加IAM用户到用户组
func (c *IamClient) KeystoneAddUserToGroupInvoker(request *model.KeystoneAddUserToGroupRequest) *KeystoneAddUserToGroupInvoker {
	requestDef := GenReqDefForKeystoneAddUserToGroup()
	return &KeystoneAddUserToGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneAssociateGroupWithDomainPermission 为用户组授予全局服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)为用户组授予全局服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneAssociateGroupWithDomainPermission(request *model.KeystoneAssociateGroupWithDomainPermissionRequest) (*model.KeystoneAssociateGroupWithDomainPermissionResponse, error) {
	requestDef := GenReqDefForKeystoneAssociateGroupWithDomainPermission()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneAssociateGroupWithDomainPermissionResponse), nil
	}
}

// KeystoneAssociateGroupWithDomainPermissionInvoker 为用户组授予全局服务权限
func (c *IamClient) KeystoneAssociateGroupWithDomainPermissionInvoker(request *model.KeystoneAssociateGroupWithDomainPermissionRequest) *KeystoneAssociateGroupWithDomainPermissionInvoker {
	requestDef := GenReqDefForKeystoneAssociateGroupWithDomainPermission()
	return &KeystoneAssociateGroupWithDomainPermissionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneAssociateGroupWithProjectPermission 为用户组授予项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)为用户组授予项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneAssociateGroupWithProjectPermission(request *model.KeystoneAssociateGroupWithProjectPermissionRequest) (*model.KeystoneAssociateGroupWithProjectPermissionResponse, error) {
	requestDef := GenReqDefForKeystoneAssociateGroupWithProjectPermission()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneAssociateGroupWithProjectPermissionResponse), nil
	}
}

// KeystoneAssociateGroupWithProjectPermissionInvoker 为用户组授予项目服务权限
func (c *IamClient) KeystoneAssociateGroupWithProjectPermissionInvoker(request *model.KeystoneAssociateGroupWithProjectPermissionRequest) *KeystoneAssociateGroupWithProjectPermissionInvoker {
	requestDef := GenReqDefForKeystoneAssociateGroupWithProjectPermission()
	return &KeystoneAssociateGroupWithProjectPermissionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCheckDomainPermissionForGroup 查询用户组是否拥有全局服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询用户组是否拥有全局服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCheckDomainPermissionForGroup(request *model.KeystoneCheckDomainPermissionForGroupRequest) (*model.KeystoneCheckDomainPermissionForGroupResponse, error) {
	requestDef := GenReqDefForKeystoneCheckDomainPermissionForGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCheckDomainPermissionForGroupResponse), nil
	}
}

// KeystoneCheckDomainPermissionForGroupInvoker 查询用户组是否拥有全局服务权限
func (c *IamClient) KeystoneCheckDomainPermissionForGroupInvoker(request *model.KeystoneCheckDomainPermissionForGroupRequest) *KeystoneCheckDomainPermissionForGroupInvoker {
	requestDef := GenReqDefForKeystoneCheckDomainPermissionForGroup()
	return &KeystoneCheckDomainPermissionForGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCheckProjectPermissionForGroup 查询用户组是否拥有项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询用户组是否拥有项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCheckProjectPermissionForGroup(request *model.KeystoneCheckProjectPermissionForGroupRequest) (*model.KeystoneCheckProjectPermissionForGroupResponse, error) {
	requestDef := GenReqDefForKeystoneCheckProjectPermissionForGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCheckProjectPermissionForGroupResponse), nil
	}
}

// KeystoneCheckProjectPermissionForGroupInvoker 查询用户组是否拥有项目服务权限
func (c *IamClient) KeystoneCheckProjectPermissionForGroupInvoker(request *model.KeystoneCheckProjectPermissionForGroupRequest) *KeystoneCheckProjectPermissionForGroupInvoker {
	requestDef := GenReqDefForKeystoneCheckProjectPermissionForGroup()
	return &KeystoneCheckProjectPermissionForGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCheckUserInGroup 查询IAM用户是否在用户组中
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户是否在用户组中。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCheckUserInGroup(request *model.KeystoneCheckUserInGroupRequest) (*model.KeystoneCheckUserInGroupResponse, error) {
	requestDef := GenReqDefForKeystoneCheckUserInGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCheckUserInGroupResponse), nil
	}
}

// KeystoneCheckUserInGroupInvoker 查询IAM用户是否在用户组中
func (c *IamClient) KeystoneCheckUserInGroupInvoker(request *model.KeystoneCheckUserInGroupRequest) *KeystoneCheckUserInGroupInvoker {
	requestDef := GenReqDefForKeystoneCheckUserInGroup()
	return &KeystoneCheckUserInGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCheckroleForGroup 查询用户组是否拥有所有项目指定权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询用户组是否拥有所有项目指定权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCheckroleForGroup(request *model.KeystoneCheckroleForGroupRequest) (*model.KeystoneCheckroleForGroupResponse, error) {
	requestDef := GenReqDefForKeystoneCheckroleForGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCheckroleForGroupResponse), nil
	}
}

// KeystoneCheckroleForGroupInvoker 查询用户组是否拥有所有项目指定权限
func (c *IamClient) KeystoneCheckroleForGroupInvoker(request *model.KeystoneCheckroleForGroupRequest) *KeystoneCheckroleForGroupInvoker {
	requestDef := GenReqDefForKeystoneCheckroleForGroup()
	return &KeystoneCheckroleForGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateGroup 创建用户组
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)创建用户组。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateGroup(request *model.KeystoneCreateGroupRequest) (*model.KeystoneCreateGroupResponse, error) {
	requestDef := GenReqDefForKeystoneCreateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateGroupResponse), nil
	}
}

// KeystoneCreateGroupInvoker 创建用户组
func (c *IamClient) KeystoneCreateGroupInvoker(request *model.KeystoneCreateGroupRequest) *KeystoneCreateGroupInvoker {
	requestDef := GenReqDefForKeystoneCreateGroup()
	return &KeystoneCreateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateIdentityProvider 注册身份提供商
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)注册身份提供商。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateIdentityProvider(request *model.KeystoneCreateIdentityProviderRequest) (*model.KeystoneCreateIdentityProviderResponse, error) {
	requestDef := GenReqDefForKeystoneCreateIdentityProvider()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateIdentityProviderResponse), nil
	}
}

// KeystoneCreateIdentityProviderInvoker 注册身份提供商
func (c *IamClient) KeystoneCreateIdentityProviderInvoker(request *model.KeystoneCreateIdentityProviderRequest) *KeystoneCreateIdentityProviderInvoker {
	requestDef := GenReqDefForKeystoneCreateIdentityProvider()
	return &KeystoneCreateIdentityProviderInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateMapping 注册映射
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)注册映射。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateMapping(request *model.KeystoneCreateMappingRequest) (*model.KeystoneCreateMappingResponse, error) {
	requestDef := GenReqDefForKeystoneCreateMapping()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateMappingResponse), nil
	}
}

// KeystoneCreateMappingInvoker 注册映射
func (c *IamClient) KeystoneCreateMappingInvoker(request *model.KeystoneCreateMappingRequest) *KeystoneCreateMappingInvoker {
	requestDef := GenReqDefForKeystoneCreateMapping()
	return &KeystoneCreateMappingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateProject 创建项目
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)创建项目。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateProject(request *model.KeystoneCreateProjectRequest) (*model.KeystoneCreateProjectResponse, error) {
	requestDef := GenReqDefForKeystoneCreateProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateProjectResponse), nil
	}
}

// KeystoneCreateProjectInvoker 创建项目
func (c *IamClient) KeystoneCreateProjectInvoker(request *model.KeystoneCreateProjectRequest) *KeystoneCreateProjectInvoker {
	requestDef := GenReqDefForKeystoneCreateProject()
	return &KeystoneCreateProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateProtocol 注册协议
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)注册协议（将协议关联到某一身份提供商）。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateProtocol(request *model.KeystoneCreateProtocolRequest) (*model.KeystoneCreateProtocolResponse, error) {
	requestDef := GenReqDefForKeystoneCreateProtocol()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateProtocolResponse), nil
	}
}

// KeystoneCreateProtocolInvoker 注册协议
func (c *IamClient) KeystoneCreateProtocolInvoker(request *model.KeystoneCreateProtocolRequest) *KeystoneCreateProtocolInvoker {
	requestDef := GenReqDefForKeystoneCreateProtocol()
	return &KeystoneCreateProtocolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateScopedToken 获取联邦认证scoped token
//
// 该接口可以用于通过联邦认证方式获取scoped token。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateScopedToken(request *model.KeystoneCreateScopedTokenRequest) (*model.KeystoneCreateScopedTokenResponse, error) {
	requestDef := GenReqDefForKeystoneCreateScopedToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateScopedTokenResponse), nil
	}
}

// KeystoneCreateScopedTokenInvoker 获取联邦认证scoped token
func (c *IamClient) KeystoneCreateScopedTokenInvoker(request *model.KeystoneCreateScopedTokenRequest) *KeystoneCreateScopedTokenInvoker {
	requestDef := GenReqDefForKeystoneCreateScopedToken()
	return &KeystoneCreateScopedTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneDeleteGroup 删除用户组
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)删除用户组。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneDeleteGroup(request *model.KeystoneDeleteGroupRequest) (*model.KeystoneDeleteGroupResponse, error) {
	requestDef := GenReqDefForKeystoneDeleteGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneDeleteGroupResponse), nil
	}
}

// KeystoneDeleteGroupInvoker 删除用户组
func (c *IamClient) KeystoneDeleteGroupInvoker(request *model.KeystoneDeleteGroupRequest) *KeystoneDeleteGroupInvoker {
	requestDef := GenReqDefForKeystoneDeleteGroup()
	return &KeystoneDeleteGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneDeleteIdentityProvider 删除身份提供商
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html) 删除身份提供商。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneDeleteIdentityProvider(request *model.KeystoneDeleteIdentityProviderRequest) (*model.KeystoneDeleteIdentityProviderResponse, error) {
	requestDef := GenReqDefForKeystoneDeleteIdentityProvider()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneDeleteIdentityProviderResponse), nil
	}
}

// KeystoneDeleteIdentityProviderInvoker 删除身份提供商
func (c *IamClient) KeystoneDeleteIdentityProviderInvoker(request *model.KeystoneDeleteIdentityProviderRequest) *KeystoneDeleteIdentityProviderInvoker {
	requestDef := GenReqDefForKeystoneDeleteIdentityProvider()
	return &KeystoneDeleteIdentityProviderInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneDeleteMapping 删除映射
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)删除映射。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneDeleteMapping(request *model.KeystoneDeleteMappingRequest) (*model.KeystoneDeleteMappingResponse, error) {
	requestDef := GenReqDefForKeystoneDeleteMapping()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneDeleteMappingResponse), nil
	}
}

// KeystoneDeleteMappingInvoker 删除映射
func (c *IamClient) KeystoneDeleteMappingInvoker(request *model.KeystoneDeleteMappingRequest) *KeystoneDeleteMappingInvoker {
	requestDef := GenReqDefForKeystoneDeleteMapping()
	return &KeystoneDeleteMappingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneDeleteProtocol 删除协议
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)删除协议。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneDeleteProtocol(request *model.KeystoneDeleteProtocolRequest) (*model.KeystoneDeleteProtocolResponse, error) {
	requestDef := GenReqDefForKeystoneDeleteProtocol()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneDeleteProtocolResponse), nil
	}
}

// KeystoneDeleteProtocolInvoker 删除协议
func (c *IamClient) KeystoneDeleteProtocolInvoker(request *model.KeystoneDeleteProtocolRequest) *KeystoneDeleteProtocolInvoker {
	requestDef := GenReqDefForKeystoneDeleteProtocol()
	return &KeystoneDeleteProtocolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListAllProjectPermissionsForGroup 查询用户组的所有项目权限列表
//
// 该接口可以用于管理员查询用户组所有项目服务权限列表。 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListAllProjectPermissionsForGroup(request *model.KeystoneListAllProjectPermissionsForGroupRequest) (*model.KeystoneListAllProjectPermissionsForGroupResponse, error) {
	requestDef := GenReqDefForKeystoneListAllProjectPermissionsForGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListAllProjectPermissionsForGroupResponse), nil
	}
}

// KeystoneListAllProjectPermissionsForGroupInvoker 查询用户组的所有项目权限列表
func (c *IamClient) KeystoneListAllProjectPermissionsForGroupInvoker(request *model.KeystoneListAllProjectPermissionsForGroupRequest) *KeystoneListAllProjectPermissionsForGroupInvoker {
	requestDef := GenReqDefForKeystoneListAllProjectPermissionsForGroup()
	return &KeystoneListAllProjectPermissionsForGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListAuthDomains 查询IAM用户可以访问的账号详情
//
// 该接口可以用于查询IAM用户可以用访问的账号详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListAuthDomains(request *model.KeystoneListAuthDomainsRequest) (*model.KeystoneListAuthDomainsResponse, error) {
	requestDef := GenReqDefForKeystoneListAuthDomains()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListAuthDomainsResponse), nil
	}
}

// KeystoneListAuthDomainsInvoker 查询IAM用户可以访问的账号详情
func (c *IamClient) KeystoneListAuthDomainsInvoker(request *model.KeystoneListAuthDomainsRequest) *KeystoneListAuthDomainsInvoker {
	requestDef := GenReqDefForKeystoneListAuthDomains()
	return &KeystoneListAuthDomainsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListAuthProjects 查询IAM用户可以访问的项目列表
//
// 该接口可以用于查询IAM用户可以访问的项目列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListAuthProjects(request *model.KeystoneListAuthProjectsRequest) (*model.KeystoneListAuthProjectsResponse, error) {
	requestDef := GenReqDefForKeystoneListAuthProjects()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListAuthProjectsResponse), nil
	}
}

// KeystoneListAuthProjectsInvoker 查询IAM用户可以访问的项目列表
func (c *IamClient) KeystoneListAuthProjectsInvoker(request *model.KeystoneListAuthProjectsRequest) *KeystoneListAuthProjectsInvoker {
	requestDef := GenReqDefForKeystoneListAuthProjects()
	return &KeystoneListAuthProjectsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListDomainPermissionsForGroup 查询全局服务中的用户组权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询全局服务中的用户组权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListDomainPermissionsForGroup(request *model.KeystoneListDomainPermissionsForGroupRequest) (*model.KeystoneListDomainPermissionsForGroupResponse, error) {
	requestDef := GenReqDefForKeystoneListDomainPermissionsForGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListDomainPermissionsForGroupResponse), nil
	}
}

// KeystoneListDomainPermissionsForGroupInvoker 查询全局服务中的用户组权限
func (c *IamClient) KeystoneListDomainPermissionsForGroupInvoker(request *model.KeystoneListDomainPermissionsForGroupRequest) *KeystoneListDomainPermissionsForGroupInvoker {
	requestDef := GenReqDefForKeystoneListDomainPermissionsForGroup()
	return &KeystoneListDomainPermissionsForGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListEndpoints 查询终端节点列表
//
// 该接口可以用于查询终端节点列表。终端节点用来提供服务访问入口。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListEndpoints(request *model.KeystoneListEndpointsRequest) (*model.KeystoneListEndpointsResponse, error) {
	requestDef := GenReqDefForKeystoneListEndpoints()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListEndpointsResponse), nil
	}
}

// KeystoneListEndpointsInvoker 查询终端节点列表
func (c *IamClient) KeystoneListEndpointsInvoker(request *model.KeystoneListEndpointsRequest) *KeystoneListEndpointsInvoker {
	requestDef := GenReqDefForKeystoneListEndpoints()
	return &KeystoneListEndpointsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListFederationDomains 查询联邦用户可以访问的账号列表
//
// 该接口用于查询联邦用户可以访问的账号列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
// &gt; - 推荐使用[查询IAM用户可以访问的账号详情](https://apiexplorer.developer.huaweicloud.com/apiexplorer/doc?product&#x3D;IAM&amp;api&#x3D;KeystoneQueryAccessibleDomainDetailsToUser)，该接口可以返回相同的响应格式。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListFederationDomains(request *model.KeystoneListFederationDomainsRequest) (*model.KeystoneListFederationDomainsResponse, error) {
	requestDef := GenReqDefForKeystoneListFederationDomains()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListFederationDomainsResponse), nil
	}
}

// KeystoneListFederationDomainsInvoker 查询联邦用户可以访问的账号列表
func (c *IamClient) KeystoneListFederationDomainsInvoker(request *model.KeystoneListFederationDomainsRequest) *KeystoneListFederationDomainsInvoker {
	requestDef := GenReqDefForKeystoneListFederationDomains()
	return &KeystoneListFederationDomainsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListGroups 查询用户组列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询用户组列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListGroups(request *model.KeystoneListGroupsRequest) (*model.KeystoneListGroupsResponse, error) {
	requestDef := GenReqDefForKeystoneListGroups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListGroupsResponse), nil
	}
}

// KeystoneListGroupsInvoker 查询用户组列表
func (c *IamClient) KeystoneListGroupsInvoker(request *model.KeystoneListGroupsRequest) *KeystoneListGroupsInvoker {
	requestDef := GenReqDefForKeystoneListGroups()
	return &KeystoneListGroupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListIdentityProviders 查询身份提供商列表
//
// 该接口可以用于查询身份提供商列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListIdentityProviders(request *model.KeystoneListIdentityProvidersRequest) (*model.KeystoneListIdentityProvidersResponse, error) {
	requestDef := GenReqDefForKeystoneListIdentityProviders()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListIdentityProvidersResponse), nil
	}
}

// KeystoneListIdentityProvidersInvoker 查询身份提供商列表
func (c *IamClient) KeystoneListIdentityProvidersInvoker(request *model.KeystoneListIdentityProvidersRequest) *KeystoneListIdentityProvidersInvoker {
	requestDef := GenReqDefForKeystoneListIdentityProviders()
	return &KeystoneListIdentityProvidersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListMappings 查询映射列表
//
// 该接口可以用于查询映射列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListMappings(request *model.KeystoneListMappingsRequest) (*model.KeystoneListMappingsResponse, error) {
	requestDef := GenReqDefForKeystoneListMappings()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListMappingsResponse), nil
	}
}

// KeystoneListMappingsInvoker 查询映射列表
func (c *IamClient) KeystoneListMappingsInvoker(request *model.KeystoneListMappingsRequest) *KeystoneListMappingsInvoker {
	requestDef := GenReqDefForKeystoneListMappings()
	return &KeystoneListMappingsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListPermissions 查询权限列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询权限列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListPermissions(request *model.KeystoneListPermissionsRequest) (*model.KeystoneListPermissionsResponse, error) {
	requestDef := GenReqDefForKeystoneListPermissions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListPermissionsResponse), nil
	}
}

// KeystoneListPermissionsInvoker 查询权限列表
func (c *IamClient) KeystoneListPermissionsInvoker(request *model.KeystoneListPermissionsRequest) *KeystoneListPermissionsInvoker {
	requestDef := GenReqDefForKeystoneListPermissions()
	return &KeystoneListPermissionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListProjectPermissionsForGroup 查询项目服务中的用户组权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询项目服务中的用户组权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListProjectPermissionsForGroup(request *model.KeystoneListProjectPermissionsForGroupRequest) (*model.KeystoneListProjectPermissionsForGroupResponse, error) {
	requestDef := GenReqDefForKeystoneListProjectPermissionsForGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListProjectPermissionsForGroupResponse), nil
	}
}

// KeystoneListProjectPermissionsForGroupInvoker 查询项目服务中的用户组权限
func (c *IamClient) KeystoneListProjectPermissionsForGroupInvoker(request *model.KeystoneListProjectPermissionsForGroupRequest) *KeystoneListProjectPermissionsForGroupInvoker {
	requestDef := GenReqDefForKeystoneListProjectPermissionsForGroup()
	return &KeystoneListProjectPermissionsForGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListProjects 查询指定条件下的项目列表
//
// 该接口可以用于查询指定条件下的项目列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListProjects(request *model.KeystoneListProjectsRequest) (*model.KeystoneListProjectsResponse, error) {
	requestDef := GenReqDefForKeystoneListProjects()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListProjectsResponse), nil
	}
}

// KeystoneListProjectsInvoker 查询指定条件下的项目列表
func (c *IamClient) KeystoneListProjectsInvoker(request *model.KeystoneListProjectsRequest) *KeystoneListProjectsInvoker {
	requestDef := GenReqDefForKeystoneListProjects()
	return &KeystoneListProjectsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListProjectsForUser 查询指定IAM用户的项目列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询指定IAM用户的项目列表，或IAM用户查询自己的项目列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListProjectsForUser(request *model.KeystoneListProjectsForUserRequest) (*model.KeystoneListProjectsForUserResponse, error) {
	requestDef := GenReqDefForKeystoneListProjectsForUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListProjectsForUserResponse), nil
	}
}

// KeystoneListProjectsForUserInvoker 查询指定IAM用户的项目列表
func (c *IamClient) KeystoneListProjectsForUserInvoker(request *model.KeystoneListProjectsForUserRequest) *KeystoneListProjectsForUserInvoker {
	requestDef := GenReqDefForKeystoneListProjectsForUser()
	return &KeystoneListProjectsForUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListProtocols 查询协议列表
//
// 该接口可以用于查询协议列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListProtocols(request *model.KeystoneListProtocolsRequest) (*model.KeystoneListProtocolsResponse, error) {
	requestDef := GenReqDefForKeystoneListProtocols()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListProtocolsResponse), nil
	}
}

// KeystoneListProtocolsInvoker 查询协议列表
func (c *IamClient) KeystoneListProtocolsInvoker(request *model.KeystoneListProtocolsRequest) *KeystoneListProtocolsInvoker {
	requestDef := GenReqDefForKeystoneListProtocols()
	return &KeystoneListProtocolsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListRegions 查询区域列表
//
// 该接口可以用于查询区域列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListRegions(request *model.KeystoneListRegionsRequest) (*model.KeystoneListRegionsResponse, error) {
	requestDef := GenReqDefForKeystoneListRegions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListRegionsResponse), nil
	}
}

// KeystoneListRegionsInvoker 查询区域列表
func (c *IamClient) KeystoneListRegionsInvoker(request *model.KeystoneListRegionsRequest) *KeystoneListRegionsInvoker {
	requestDef := GenReqDefForKeystoneListRegions()
	return &KeystoneListRegionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListServices 查询服务列表
//
// 该接口可以用于查询服务列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListServices(request *model.KeystoneListServicesRequest) (*model.KeystoneListServicesResponse, error) {
	requestDef := GenReqDefForKeystoneListServices()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListServicesResponse), nil
	}
}

// KeystoneListServicesInvoker 查询服务列表
func (c *IamClient) KeystoneListServicesInvoker(request *model.KeystoneListServicesRequest) *KeystoneListServicesInvoker {
	requestDef := GenReqDefForKeystoneListServices()
	return &KeystoneListServicesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListVersions 查询版本信息列表
//
// 该接口用于查询Keystone API的版本信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListVersions(request *model.KeystoneListVersionsRequest) (*model.KeystoneListVersionsResponse, error) {
	requestDef := GenReqDefForKeystoneListVersions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListVersionsResponse), nil
	}
}

// KeystoneListVersionsInvoker 查询版本信息列表
func (c *IamClient) KeystoneListVersionsInvoker(request *model.KeystoneListVersionsRequest) *KeystoneListVersionsInvoker {
	requestDef := GenReqDefForKeystoneListVersions()
	return &KeystoneListVersionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneRemoveDomainPermissionFromGroup 移除用户组的全局服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)移除用户组的全局服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneRemoveDomainPermissionFromGroup(request *model.KeystoneRemoveDomainPermissionFromGroupRequest) (*model.KeystoneRemoveDomainPermissionFromGroupResponse, error) {
	requestDef := GenReqDefForKeystoneRemoveDomainPermissionFromGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneRemoveDomainPermissionFromGroupResponse), nil
	}
}

// KeystoneRemoveDomainPermissionFromGroupInvoker 移除用户组的全局服务权限
func (c *IamClient) KeystoneRemoveDomainPermissionFromGroupInvoker(request *model.KeystoneRemoveDomainPermissionFromGroupRequest) *KeystoneRemoveDomainPermissionFromGroupInvoker {
	requestDef := GenReqDefForKeystoneRemoveDomainPermissionFromGroup()
	return &KeystoneRemoveDomainPermissionFromGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneRemoveProjectPermissionFromGroup 移除用户组的项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)移除用户组的项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneRemoveProjectPermissionFromGroup(request *model.KeystoneRemoveProjectPermissionFromGroupRequest) (*model.KeystoneRemoveProjectPermissionFromGroupResponse, error) {
	requestDef := GenReqDefForKeystoneRemoveProjectPermissionFromGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneRemoveProjectPermissionFromGroupResponse), nil
	}
}

// KeystoneRemoveProjectPermissionFromGroupInvoker 移除用户组的项目服务权限
func (c *IamClient) KeystoneRemoveProjectPermissionFromGroupInvoker(request *model.KeystoneRemoveProjectPermissionFromGroupRequest) *KeystoneRemoveProjectPermissionFromGroupInvoker {
	requestDef := GenReqDefForKeystoneRemoveProjectPermissionFromGroup()
	return &KeystoneRemoveProjectPermissionFromGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneRemoveUserFromGroup 移除用户组中的IAM用户
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)移除用户组中的IAM用户。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneRemoveUserFromGroup(request *model.KeystoneRemoveUserFromGroupRequest) (*model.KeystoneRemoveUserFromGroupResponse, error) {
	requestDef := GenReqDefForKeystoneRemoveUserFromGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneRemoveUserFromGroupResponse), nil
	}
}

// KeystoneRemoveUserFromGroupInvoker 移除用户组中的IAM用户
func (c *IamClient) KeystoneRemoveUserFromGroupInvoker(request *model.KeystoneRemoveUserFromGroupRequest) *KeystoneRemoveUserFromGroupInvoker {
	requestDef := GenReqDefForKeystoneRemoveUserFromGroup()
	return &KeystoneRemoveUserFromGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowCatalog 查询服务目录
//
// 该接口可以用于查询请求头中X-Auth-Token对应的服务目录。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowCatalog(request *model.KeystoneShowCatalogRequest) (*model.KeystoneShowCatalogResponse, error) {
	requestDef := GenReqDefForKeystoneShowCatalog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowCatalogResponse), nil
	}
}

// KeystoneShowCatalogInvoker 查询服务目录
func (c *IamClient) KeystoneShowCatalogInvoker(request *model.KeystoneShowCatalogRequest) *KeystoneShowCatalogInvoker {
	requestDef := GenReqDefForKeystoneShowCatalog()
	return &KeystoneShowCatalogInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowEndpoint 查询终端节点详情
//
// 该接口可以用于查询终端节点详情。终端节点用来提供服务访问入口。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowEndpoint(request *model.KeystoneShowEndpointRequest) (*model.KeystoneShowEndpointResponse, error) {
	requestDef := GenReqDefForKeystoneShowEndpoint()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowEndpointResponse), nil
	}
}

// KeystoneShowEndpointInvoker 查询终端节点详情
func (c *IamClient) KeystoneShowEndpointInvoker(request *model.KeystoneShowEndpointRequest) *KeystoneShowEndpointInvoker {
	requestDef := GenReqDefForKeystoneShowEndpoint()
	return &KeystoneShowEndpointInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowGroup 查询用户组详情
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询用户组详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowGroup(request *model.KeystoneShowGroupRequest) (*model.KeystoneShowGroupResponse, error) {
	requestDef := GenReqDefForKeystoneShowGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowGroupResponse), nil
	}
}

// KeystoneShowGroupInvoker 查询用户组详情
func (c *IamClient) KeystoneShowGroupInvoker(request *model.KeystoneShowGroupRequest) *KeystoneShowGroupInvoker {
	requestDef := GenReqDefForKeystoneShowGroup()
	return &KeystoneShowGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowIdentityProvider 查询身份提供商详情
//
// 该接口可以用于查询身份提供商详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowIdentityProvider(request *model.KeystoneShowIdentityProviderRequest) (*model.KeystoneShowIdentityProviderResponse, error) {
	requestDef := GenReqDefForKeystoneShowIdentityProvider()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowIdentityProviderResponse), nil
	}
}

// KeystoneShowIdentityProviderInvoker 查询身份提供商详情
func (c *IamClient) KeystoneShowIdentityProviderInvoker(request *model.KeystoneShowIdentityProviderRequest) *KeystoneShowIdentityProviderInvoker {
	requestDef := GenReqDefForKeystoneShowIdentityProvider()
	return &KeystoneShowIdentityProviderInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowMapping 查询映射详情
//
// 该接口可以用于查询映射详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowMapping(request *model.KeystoneShowMappingRequest) (*model.KeystoneShowMappingResponse, error) {
	requestDef := GenReqDefForKeystoneShowMapping()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowMappingResponse), nil
	}
}

// KeystoneShowMappingInvoker 查询映射详情
func (c *IamClient) KeystoneShowMappingInvoker(request *model.KeystoneShowMappingRequest) *KeystoneShowMappingInvoker {
	requestDef := GenReqDefForKeystoneShowMapping()
	return &KeystoneShowMappingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowPermission 查询权限详情
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询权限详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowPermission(request *model.KeystoneShowPermissionRequest) (*model.KeystoneShowPermissionResponse, error) {
	requestDef := GenReqDefForKeystoneShowPermission()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowPermissionResponse), nil
	}
}

// KeystoneShowPermissionInvoker 查询权限详情
func (c *IamClient) KeystoneShowPermissionInvoker(request *model.KeystoneShowPermissionRequest) *KeystoneShowPermissionInvoker {
	requestDef := GenReqDefForKeystoneShowPermission()
	return &KeystoneShowPermissionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowProject 查询项目详情
//
// 该接口可以用于查询项目详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowProject(request *model.KeystoneShowProjectRequest) (*model.KeystoneShowProjectResponse, error) {
	requestDef := GenReqDefForKeystoneShowProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowProjectResponse), nil
	}
}

// KeystoneShowProjectInvoker 查询项目详情
func (c *IamClient) KeystoneShowProjectInvoker(request *model.KeystoneShowProjectRequest) *KeystoneShowProjectInvoker {
	requestDef := GenReqDefForKeystoneShowProject()
	return &KeystoneShowProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowProtocol 查询协议详情
//
// 该接口可以用于查询协议详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowProtocol(request *model.KeystoneShowProtocolRequest) (*model.KeystoneShowProtocolResponse, error) {
	requestDef := GenReqDefForKeystoneShowProtocol()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowProtocolResponse), nil
	}
}

// KeystoneShowProtocolInvoker 查询协议详情
func (c *IamClient) KeystoneShowProtocolInvoker(request *model.KeystoneShowProtocolRequest) *KeystoneShowProtocolInvoker {
	requestDef := GenReqDefForKeystoneShowProtocol()
	return &KeystoneShowProtocolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowRegion 查询区域详情
//
// 该接口可以用于查询区域详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowRegion(request *model.KeystoneShowRegionRequest) (*model.KeystoneShowRegionResponse, error) {
	requestDef := GenReqDefForKeystoneShowRegion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowRegionResponse), nil
	}
}

// KeystoneShowRegionInvoker 查询区域详情
func (c *IamClient) KeystoneShowRegionInvoker(request *model.KeystoneShowRegionRequest) *KeystoneShowRegionInvoker {
	requestDef := GenReqDefForKeystoneShowRegion()
	return &KeystoneShowRegionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowSecurityCompliance 查询账号密码强度策略
//
// 该接口可以用于查询账号密码强度策略，查询结果包括密码强度策略的正则表达式及其描述。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowSecurityCompliance(request *model.KeystoneShowSecurityComplianceRequest) (*model.KeystoneShowSecurityComplianceResponse, error) {
	requestDef := GenReqDefForKeystoneShowSecurityCompliance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowSecurityComplianceResponse), nil
	}
}

// KeystoneShowSecurityComplianceInvoker 查询账号密码强度策略
func (c *IamClient) KeystoneShowSecurityComplianceInvoker(request *model.KeystoneShowSecurityComplianceRequest) *KeystoneShowSecurityComplianceInvoker {
	requestDef := GenReqDefForKeystoneShowSecurityCompliance()
	return &KeystoneShowSecurityComplianceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowSecurityComplianceByOption 按条件查询账号密码强度策略
//
// 该接口可以用于按条件查询账号密码强度策略，查询结果包括密码强度策略的正则表达式及其描述。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowSecurityComplianceByOption(request *model.KeystoneShowSecurityComplianceByOptionRequest) (*model.KeystoneShowSecurityComplianceByOptionResponse, error) {
	requestDef := GenReqDefForKeystoneShowSecurityComplianceByOption()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowSecurityComplianceByOptionResponse), nil
	}
}

// KeystoneShowSecurityComplianceByOptionInvoker 按条件查询账号密码强度策略
func (c *IamClient) KeystoneShowSecurityComplianceByOptionInvoker(request *model.KeystoneShowSecurityComplianceByOptionRequest) *KeystoneShowSecurityComplianceByOptionInvoker {
	requestDef := GenReqDefForKeystoneShowSecurityComplianceByOption()
	return &KeystoneShowSecurityComplianceByOptionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowService 查询服务详情
//
// 该接口可以用于查询服务详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowService(request *model.KeystoneShowServiceRequest) (*model.KeystoneShowServiceResponse, error) {
	requestDef := GenReqDefForKeystoneShowService()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowServiceResponse), nil
	}
}

// KeystoneShowServiceInvoker 查询服务详情
func (c *IamClient) KeystoneShowServiceInvoker(request *model.KeystoneShowServiceRequest) *KeystoneShowServiceInvoker {
	requestDef := GenReqDefForKeystoneShowService()
	return &KeystoneShowServiceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowVersion 查询版本信息
//
// 该接口用于查询Keystone API的3.0版本的信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowVersion(request *model.KeystoneShowVersionRequest) (*model.KeystoneShowVersionResponse, error) {
	requestDef := GenReqDefForKeystoneShowVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowVersionResponse), nil
	}
}

// KeystoneShowVersionInvoker 查询版本信息
func (c *IamClient) KeystoneShowVersionInvoker(request *model.KeystoneShowVersionRequest) *KeystoneShowVersionInvoker {
	requestDef := GenReqDefForKeystoneShowVersion()
	return &KeystoneShowVersionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneUpdateGroup 更新用户组
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)更新用户组信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneUpdateGroup(request *model.KeystoneUpdateGroupRequest) (*model.KeystoneUpdateGroupResponse, error) {
	requestDef := GenReqDefForKeystoneUpdateGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneUpdateGroupResponse), nil
	}
}

// KeystoneUpdateGroupInvoker 更新用户组
func (c *IamClient) KeystoneUpdateGroupInvoker(request *model.KeystoneUpdateGroupRequest) *KeystoneUpdateGroupInvoker {
	requestDef := GenReqDefForKeystoneUpdateGroup()
	return &KeystoneUpdateGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneUpdateIdentityProvider 更新身份提供商
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)更新身份提供商。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneUpdateIdentityProvider(request *model.KeystoneUpdateIdentityProviderRequest) (*model.KeystoneUpdateIdentityProviderResponse, error) {
	requestDef := GenReqDefForKeystoneUpdateIdentityProvider()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneUpdateIdentityProviderResponse), nil
	}
}

// KeystoneUpdateIdentityProviderInvoker 更新身份提供商
func (c *IamClient) KeystoneUpdateIdentityProviderInvoker(request *model.KeystoneUpdateIdentityProviderRequest) *KeystoneUpdateIdentityProviderInvoker {
	requestDef := GenReqDefForKeystoneUpdateIdentityProvider()
	return &KeystoneUpdateIdentityProviderInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneUpdateMapping 更新映射
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)更新映射。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneUpdateMapping(request *model.KeystoneUpdateMappingRequest) (*model.KeystoneUpdateMappingResponse, error) {
	requestDef := GenReqDefForKeystoneUpdateMapping()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneUpdateMappingResponse), nil
	}
}

// KeystoneUpdateMappingInvoker 更新映射
func (c *IamClient) KeystoneUpdateMappingInvoker(request *model.KeystoneUpdateMappingRequest) *KeystoneUpdateMappingInvoker {
	requestDef := GenReqDefForKeystoneUpdateMapping()
	return &KeystoneUpdateMappingInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneUpdateProject 修改项目信息
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改项目信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneUpdateProject(request *model.KeystoneUpdateProjectRequest) (*model.KeystoneUpdateProjectResponse, error) {
	requestDef := GenReqDefForKeystoneUpdateProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneUpdateProjectResponse), nil
	}
}

// KeystoneUpdateProjectInvoker 修改项目信息
func (c *IamClient) KeystoneUpdateProjectInvoker(request *model.KeystoneUpdateProjectRequest) *KeystoneUpdateProjectInvoker {
	requestDef := GenReqDefForKeystoneUpdateProject()
	return &KeystoneUpdateProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneUpdateProtocol 更新协议
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)更新协议。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneUpdateProtocol(request *model.KeystoneUpdateProtocolRequest) (*model.KeystoneUpdateProtocolResponse, error) {
	requestDef := GenReqDefForKeystoneUpdateProtocol()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneUpdateProtocolResponse), nil
	}
}

// KeystoneUpdateProtocolInvoker 更新协议
func (c *IamClient) KeystoneUpdateProtocolInvoker(request *model.KeystoneUpdateProtocolRequest) *KeystoneUpdateProtocolInvoker {
	requestDef := GenReqDefForKeystoneUpdateProtocol()
	return &KeystoneUpdateProtocolInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAgencies 查询指定条件下的委托列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询指定条件下的委托列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListAgencies(request *model.ListAgenciesRequest) (*model.ListAgenciesResponse, error) {
	requestDef := GenReqDefForListAgencies()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAgenciesResponse), nil
	}
}

// ListAgenciesInvoker 查询指定条件下的委托列表
func (c *IamClient) ListAgenciesInvoker(request *model.ListAgenciesRequest) *ListAgenciesInvoker {
	requestDef := GenReqDefForListAgencies()
	return &ListAgenciesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAllProjectsPermissionsForAgency 查询委托下的所有项目服务权限列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询委托所有项目服务权限列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListAllProjectsPermissionsForAgency(request *model.ListAllProjectsPermissionsForAgencyRequest) (*model.ListAllProjectsPermissionsForAgencyResponse, error) {
	requestDef := GenReqDefForListAllProjectsPermissionsForAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAllProjectsPermissionsForAgencyResponse), nil
	}
}

// ListAllProjectsPermissionsForAgencyInvoker 查询委托下的所有项目服务权限列表
func (c *IamClient) ListAllProjectsPermissionsForAgencyInvoker(request *model.ListAllProjectsPermissionsForAgencyRequest) *ListAllProjectsPermissionsForAgencyInvoker {
	requestDef := GenReqDefForListAllProjectsPermissionsForAgency()
	return &ListAllProjectsPermissionsForAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListCustomPolicies 查询自定义策略列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询自定义策略列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListCustomPolicies(request *model.ListCustomPoliciesRequest) (*model.ListCustomPoliciesResponse, error) {
	requestDef := GenReqDefForListCustomPolicies()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListCustomPoliciesResponse), nil
	}
}

// ListCustomPoliciesInvoker 查询自定义策略列表
func (c *IamClient) ListCustomPoliciesInvoker(request *model.ListCustomPoliciesRequest) *ListCustomPoliciesInvoker {
	requestDef := GenReqDefForListCustomPolicies()
	return &ListCustomPoliciesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDomainPermissionsForAgency 查询全局服务中的委托权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询全局服务中的委托权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListDomainPermissionsForAgency(request *model.ListDomainPermissionsForAgencyRequest) (*model.ListDomainPermissionsForAgencyResponse, error) {
	requestDef := GenReqDefForListDomainPermissionsForAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDomainPermissionsForAgencyResponse), nil
	}
}

// ListDomainPermissionsForAgencyInvoker 查询全局服务中的委托权限
func (c *IamClient) ListDomainPermissionsForAgencyInvoker(request *model.ListDomainPermissionsForAgencyRequest) *ListDomainPermissionsForAgencyInvoker {
	requestDef := GenReqDefForListDomainPermissionsForAgency()
	return &ListDomainPermissionsForAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListEnterpriseProjectsForGroup 查询用户组关联的企业项目
//
// 该接口可用于查询用户组所关联的企业项目。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListEnterpriseProjectsForGroup(request *model.ListEnterpriseProjectsForGroupRequest) (*model.ListEnterpriseProjectsForGroupResponse, error) {
	requestDef := GenReqDefForListEnterpriseProjectsForGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListEnterpriseProjectsForGroupResponse), nil
	}
}

// ListEnterpriseProjectsForGroupInvoker 查询用户组关联的企业项目
func (c *IamClient) ListEnterpriseProjectsForGroupInvoker(request *model.ListEnterpriseProjectsForGroupRequest) *ListEnterpriseProjectsForGroupInvoker {
	requestDef := GenReqDefForListEnterpriseProjectsForGroup()
	return &ListEnterpriseProjectsForGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListEnterpriseProjectsForUser 查询用户直接关联的企业项目
//
// 该接口可用于查询用户所关联的企业项目。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListEnterpriseProjectsForUser(request *model.ListEnterpriseProjectsForUserRequest) (*model.ListEnterpriseProjectsForUserResponse, error) {
	requestDef := GenReqDefForListEnterpriseProjectsForUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListEnterpriseProjectsForUserResponse), nil
	}
}

// ListEnterpriseProjectsForUserInvoker 查询用户直接关联的企业项目
func (c *IamClient) ListEnterpriseProjectsForUserInvoker(request *model.ListEnterpriseProjectsForUserRequest) *ListEnterpriseProjectsForUserInvoker {
	requestDef := GenReqDefForListEnterpriseProjectsForUser()
	return &ListEnterpriseProjectsForUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListGroupsForEnterpriseProject 查询企业项目关联的用户组
//
// 该接口可用于查询企业项目关联的用户组。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListGroupsForEnterpriseProject(request *model.ListGroupsForEnterpriseProjectRequest) (*model.ListGroupsForEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForListGroupsForEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListGroupsForEnterpriseProjectResponse), nil
	}
}

// ListGroupsForEnterpriseProjectInvoker 查询企业项目关联的用户组
func (c *IamClient) ListGroupsForEnterpriseProjectInvoker(request *model.ListGroupsForEnterpriseProjectRequest) *ListGroupsForEnterpriseProjectInvoker {
	requestDef := GenReqDefForListGroupsForEnterpriseProject()
	return &ListGroupsForEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProjectPermissionsForAgency 查询项目服务中的委托权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询项目服务中的委托权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListProjectPermissionsForAgency(request *model.ListProjectPermissionsForAgencyRequest) (*model.ListProjectPermissionsForAgencyResponse, error) {
	requestDef := GenReqDefForListProjectPermissionsForAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProjectPermissionsForAgencyResponse), nil
	}
}

// ListProjectPermissionsForAgencyInvoker 查询项目服务中的委托权限
func (c *IamClient) ListProjectPermissionsForAgencyInvoker(request *model.ListProjectPermissionsForAgencyRequest) *ListProjectPermissionsForAgencyInvoker {
	requestDef := GenReqDefForListProjectPermissionsForAgency()
	return &ListProjectPermissionsForAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRolesForGroupOnEnterpriseProject 查询企业项目关联用户组的权限
//
// 该接口可用于查询企业项目已关联用户组的权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListRolesForGroupOnEnterpriseProject(request *model.ListRolesForGroupOnEnterpriseProjectRequest) (*model.ListRolesForGroupOnEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForListRolesForGroupOnEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRolesForGroupOnEnterpriseProjectResponse), nil
	}
}

// ListRolesForGroupOnEnterpriseProjectInvoker 查询企业项目关联用户组的权限
func (c *IamClient) ListRolesForGroupOnEnterpriseProjectInvoker(request *model.ListRolesForGroupOnEnterpriseProjectRequest) *ListRolesForGroupOnEnterpriseProjectInvoker {
	requestDef := GenReqDefForListRolesForGroupOnEnterpriseProject()
	return &ListRolesForGroupOnEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRolesForUserOnEnterpriseProject 查询企业项目直接关联用户的权限
//
// 该接口可用于查询企业项目直接关联用户的权限。
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListRolesForUserOnEnterpriseProject(request *model.ListRolesForUserOnEnterpriseProjectRequest) (*model.ListRolesForUserOnEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForListRolesForUserOnEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRolesForUserOnEnterpriseProjectResponse), nil
	}
}

// ListRolesForUserOnEnterpriseProjectInvoker 查询企业项目直接关联用户的权限
func (c *IamClient) ListRolesForUserOnEnterpriseProjectInvoker(request *model.ListRolesForUserOnEnterpriseProjectRequest) *ListRolesForUserOnEnterpriseProjectInvoker {
	requestDef := GenReqDefForListRolesForUserOnEnterpriseProject()
	return &ListRolesForUserOnEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUsersForEnterpriseProject 查询企业项目直接关联用户
//
// 该接口可用于查询企业项目直接关联的用户。
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListUsersForEnterpriseProject(request *model.ListUsersForEnterpriseProjectRequest) (*model.ListUsersForEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForListUsersForEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUsersForEnterpriseProjectResponse), nil
	}
}

// ListUsersForEnterpriseProjectInvoker 查询企业项目直接关联用户
func (c *IamClient) ListUsersForEnterpriseProjectInvoker(request *model.ListUsersForEnterpriseProjectRequest) *ListUsersForEnterpriseProjectInvoker {
	requestDef := GenReqDefForListUsersForEnterpriseProject()
	return &ListUsersForEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveAllProjectsPermissionFromAgency 移除委托下的所有项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)移除委托的所有项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) RemoveAllProjectsPermissionFromAgency(request *model.RemoveAllProjectsPermissionFromAgencyRequest) (*model.RemoveAllProjectsPermissionFromAgencyResponse, error) {
	requestDef := GenReqDefForRemoveAllProjectsPermissionFromAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveAllProjectsPermissionFromAgencyResponse), nil
	}
}

// RemoveAllProjectsPermissionFromAgencyInvoker 移除委托下的所有项目服务权限
func (c *IamClient) RemoveAllProjectsPermissionFromAgencyInvoker(request *model.RemoveAllProjectsPermissionFromAgencyRequest) *RemoveAllProjectsPermissionFromAgencyInvoker {
	requestDef := GenReqDefForRemoveAllProjectsPermissionFromAgency()
	return &RemoveAllProjectsPermissionFromAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveDomainPermissionFromAgency 移除委托的全局服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)移除委托的全局服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) RemoveDomainPermissionFromAgency(request *model.RemoveDomainPermissionFromAgencyRequest) (*model.RemoveDomainPermissionFromAgencyResponse, error) {
	requestDef := GenReqDefForRemoveDomainPermissionFromAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveDomainPermissionFromAgencyResponse), nil
	}
}

// RemoveDomainPermissionFromAgencyInvoker 移除委托的全局服务权限
func (c *IamClient) RemoveDomainPermissionFromAgencyInvoker(request *model.RemoveDomainPermissionFromAgencyRequest) *RemoveDomainPermissionFromAgencyInvoker {
	requestDef := GenReqDefForRemoveDomainPermissionFromAgency()
	return &RemoveDomainPermissionFromAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveProjectPermissionFromAgency 移除委托的项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)移除委托的项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) RemoveProjectPermissionFromAgency(request *model.RemoveProjectPermissionFromAgencyRequest) (*model.RemoveProjectPermissionFromAgencyResponse, error) {
	requestDef := GenReqDefForRemoveProjectPermissionFromAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveProjectPermissionFromAgencyResponse), nil
	}
}

// RemoveProjectPermissionFromAgencyInvoker 移除委托的项目服务权限
func (c *IamClient) RemoveProjectPermissionFromAgencyInvoker(request *model.RemoveProjectPermissionFromAgencyRequest) *RemoveProjectPermissionFromAgencyInvoker {
	requestDef := GenReqDefForRemoveProjectPermissionFromAgency()
	return &RemoveProjectPermissionFromAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RevokeRoleFromAgencyOnEnterpriseProject 删除企业项目关联委托的权限
//
// 该接口可以删除企业项目委托上的授权
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) RevokeRoleFromAgencyOnEnterpriseProject(request *model.RevokeRoleFromAgencyOnEnterpriseProjectRequest) (*model.RevokeRoleFromAgencyOnEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForRevokeRoleFromAgencyOnEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeRoleFromAgencyOnEnterpriseProjectResponse), nil
	}
}

// RevokeRoleFromAgencyOnEnterpriseProjectInvoker 删除企业项目关联委托的权限
func (c *IamClient) RevokeRoleFromAgencyOnEnterpriseProjectInvoker(request *model.RevokeRoleFromAgencyOnEnterpriseProjectRequest) *RevokeRoleFromAgencyOnEnterpriseProjectInvoker {
	requestDef := GenReqDefForRevokeRoleFromAgencyOnEnterpriseProject()
	return &RevokeRoleFromAgencyOnEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RevokeRoleFromGroupOnEnterpriseProject 删除企业项目关联用户组的权限
//
// 该接口用于删除企业项目关联用户组的权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) RevokeRoleFromGroupOnEnterpriseProject(request *model.RevokeRoleFromGroupOnEnterpriseProjectRequest) (*model.RevokeRoleFromGroupOnEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForRevokeRoleFromGroupOnEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeRoleFromGroupOnEnterpriseProjectResponse), nil
	}
}

// RevokeRoleFromGroupOnEnterpriseProjectInvoker 删除企业项目关联用户组的权限
func (c *IamClient) RevokeRoleFromGroupOnEnterpriseProjectInvoker(request *model.RevokeRoleFromGroupOnEnterpriseProjectRequest) *RevokeRoleFromGroupOnEnterpriseProjectInvoker {
	requestDef := GenReqDefForRevokeRoleFromGroupOnEnterpriseProject()
	return &RevokeRoleFromGroupOnEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RevokeRoleFromUserOnEnterpriseProject 删除企业项目直接关联用户的权限
//
// 删除企业项目直接关联用户的权限。
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) RevokeRoleFromUserOnEnterpriseProject(request *model.RevokeRoleFromUserOnEnterpriseProjectRequest) (*model.RevokeRoleFromUserOnEnterpriseProjectResponse, error) {
	requestDef := GenReqDefForRevokeRoleFromUserOnEnterpriseProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeRoleFromUserOnEnterpriseProjectResponse), nil
	}
}

// RevokeRoleFromUserOnEnterpriseProjectInvoker 删除企业项目直接关联用户的权限
func (c *IamClient) RevokeRoleFromUserOnEnterpriseProjectInvoker(request *model.RevokeRoleFromUserOnEnterpriseProjectRequest) *RevokeRoleFromUserOnEnterpriseProjectInvoker {
	requestDef := GenReqDefForRevokeRoleFromUserOnEnterpriseProject()
	return &RevokeRoleFromUserOnEnterpriseProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAgency 查询委托详情
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询委托详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowAgency(request *model.ShowAgencyRequest) (*model.ShowAgencyResponse, error) {
	requestDef := GenReqDefForShowAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAgencyResponse), nil
	}
}

// ShowAgencyInvoker 查询委托详情
func (c *IamClient) ShowAgencyInvoker(request *model.ShowAgencyRequest) *ShowAgencyInvoker {
	requestDef := GenReqDefForShowAgency()
	return &ShowAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowCustomPolicy 查询自定义策略详情
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询自定义策略详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowCustomPolicy(request *model.ShowCustomPolicyRequest) (*model.ShowCustomPolicyResponse, error) {
	requestDef := GenReqDefForShowCustomPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCustomPolicyResponse), nil
	}
}

// ShowCustomPolicyInvoker 查询自定义策略详情
func (c *IamClient) ShowCustomPolicyInvoker(request *model.ShowCustomPolicyRequest) *ShowCustomPolicyInvoker {
	requestDef := GenReqDefForShowCustomPolicy()
	return &ShowCustomPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainApiAclPolicy 查询账号接口访问策略
//
// 该接口可以用于查询账号接口访问控制策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowDomainApiAclPolicy(request *model.ShowDomainApiAclPolicyRequest) (*model.ShowDomainApiAclPolicyResponse, error) {
	requestDef := GenReqDefForShowDomainApiAclPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainApiAclPolicyResponse), nil
	}
}

// ShowDomainApiAclPolicyInvoker 查询账号接口访问策略
func (c *IamClient) ShowDomainApiAclPolicyInvoker(request *model.ShowDomainApiAclPolicyRequest) *ShowDomainApiAclPolicyInvoker {
	requestDef := GenReqDefForShowDomainApiAclPolicy()
	return &ShowDomainApiAclPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainConsoleAclPolicy 查询账号控制台访问策略
//
// 该接口可以用于查询账号控制台访问控制策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowDomainConsoleAclPolicy(request *model.ShowDomainConsoleAclPolicyRequest) (*model.ShowDomainConsoleAclPolicyResponse, error) {
	requestDef := GenReqDefForShowDomainConsoleAclPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainConsoleAclPolicyResponse), nil
	}
}

// ShowDomainConsoleAclPolicyInvoker 查询账号控制台访问策略
func (c *IamClient) ShowDomainConsoleAclPolicyInvoker(request *model.ShowDomainConsoleAclPolicyRequest) *ShowDomainConsoleAclPolicyInvoker {
	requestDef := GenReqDefForShowDomainConsoleAclPolicy()
	return &ShowDomainConsoleAclPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainLoginPolicy 查询账号登录策略
//
// 该接口可以用于查询账号登录策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowDomainLoginPolicy(request *model.ShowDomainLoginPolicyRequest) (*model.ShowDomainLoginPolicyResponse, error) {
	requestDef := GenReqDefForShowDomainLoginPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainLoginPolicyResponse), nil
	}
}

// ShowDomainLoginPolicyInvoker 查询账号登录策略
func (c *IamClient) ShowDomainLoginPolicyInvoker(request *model.ShowDomainLoginPolicyRequest) *ShowDomainLoginPolicyInvoker {
	requestDef := GenReqDefForShowDomainLoginPolicy()
	return &ShowDomainLoginPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainPasswordPolicy 查询账号密码策略
//
// 该接口可以用于查询账号密码策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowDomainPasswordPolicy(request *model.ShowDomainPasswordPolicyRequest) (*model.ShowDomainPasswordPolicyResponse, error) {
	requestDef := GenReqDefForShowDomainPasswordPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainPasswordPolicyResponse), nil
	}
}

// ShowDomainPasswordPolicyInvoker 查询账号密码策略
func (c *IamClient) ShowDomainPasswordPolicyInvoker(request *model.ShowDomainPasswordPolicyRequest) *ShowDomainPasswordPolicyInvoker {
	requestDef := GenReqDefForShowDomainPasswordPolicy()
	return &ShowDomainPasswordPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainProtectPolicy 查询账号操作保护策略
//
// 该接口可以用于查询账号操作保护策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowDomainProtectPolicy(request *model.ShowDomainProtectPolicyRequest) (*model.ShowDomainProtectPolicyResponse, error) {
	requestDef := GenReqDefForShowDomainProtectPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainProtectPolicyResponse), nil
	}
}

// ShowDomainProtectPolicyInvoker 查询账号操作保护策略
func (c *IamClient) ShowDomainProtectPolicyInvoker(request *model.ShowDomainProtectPolicyRequest) *ShowDomainProtectPolicyInvoker {
	requestDef := GenReqDefForShowDomainProtectPolicy()
	return &ShowDomainProtectPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainQuota 查询账号配额
//
// 该接口可以用于查询账号配额。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowDomainQuota(request *model.ShowDomainQuotaRequest) (*model.ShowDomainQuotaResponse, error) {
	requestDef := GenReqDefForShowDomainQuota()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainQuotaResponse), nil
	}
}

// ShowDomainQuotaInvoker 查询账号配额
func (c *IamClient) ShowDomainQuotaInvoker(request *model.ShowDomainQuotaRequest) *ShowDomainQuotaInvoker {
	requestDef := GenReqDefForShowDomainQuota()
	return &ShowDomainQuotaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainRoleAssignments 查询指定账号中的授权记录
//
// 该接口用于查询指定账号中的授权记录。
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowDomainRoleAssignments(request *model.ShowDomainRoleAssignmentsRequest) (*model.ShowDomainRoleAssignmentsResponse, error) {
	requestDef := GenReqDefForShowDomainRoleAssignments()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainRoleAssignmentsResponse), nil
	}
}

// ShowDomainRoleAssignmentsInvoker 查询指定账号中的授权记录
func (c *IamClient) ShowDomainRoleAssignmentsInvoker(request *model.ShowDomainRoleAssignmentsRequest) *ShowDomainRoleAssignmentsInvoker {
	requestDef := GenReqDefForShowDomainRoleAssignments()
	return &ShowDomainRoleAssignmentsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowMetadata 查询Metadata文件
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询身份提供商导入到IAM中的Metadata文件。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowMetadata(request *model.ShowMetadataRequest) (*model.ShowMetadataResponse, error) {
	requestDef := GenReqDefForShowMetadata()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowMetadataResponse), nil
	}
}

// ShowMetadataInvoker 查询Metadata文件
func (c *IamClient) ShowMetadataInvoker(request *model.ShowMetadataRequest) *ShowMetadataInvoker {
	requestDef := GenReqDefForShowMetadata()
	return &ShowMetadataInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowOpenIdConnectConfig 查询OpenId Connect身份提供商配置
//
// 查询OpenId Connect身份提供商配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowOpenIdConnectConfig(request *model.ShowOpenIdConnectConfigRequest) (*model.ShowOpenIdConnectConfigResponse, error) {
	requestDef := GenReqDefForShowOpenIdConnectConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowOpenIdConnectConfigResponse), nil
	}
}

// ShowOpenIdConnectConfigInvoker 查询OpenId Connect身份提供商配置
func (c *IamClient) ShowOpenIdConnectConfigInvoker(request *model.ShowOpenIdConnectConfigRequest) *ShowOpenIdConnectConfigInvoker {
	requestDef := GenReqDefForShowOpenIdConnectConfig()
	return &ShowOpenIdConnectConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowProjectDetailsAndStatus 查询项目详情与状态
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询项目详情与状态。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowProjectDetailsAndStatus(request *model.ShowProjectDetailsAndStatusRequest) (*model.ShowProjectDetailsAndStatusResponse, error) {
	requestDef := GenReqDefForShowProjectDetailsAndStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowProjectDetailsAndStatusResponse), nil
	}
}

// ShowProjectDetailsAndStatusInvoker 查询项目详情与状态
func (c *IamClient) ShowProjectDetailsAndStatusInvoker(request *model.ShowProjectDetailsAndStatusRequest) *ShowProjectDetailsAndStatusInvoker {
	requestDef := GenReqDefForShowProjectDetailsAndStatus()
	return &ShowProjectDetailsAndStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowProjectQuota 查询项目配额
//
// 该接口可以用于查询项目配额。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowProjectQuota(request *model.ShowProjectQuotaRequest) (*model.ShowProjectQuotaResponse, error) {
	requestDef := GenReqDefForShowProjectQuota()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowProjectQuotaResponse), nil
	}
}

// ShowProjectQuotaInvoker 查询项目配额
func (c *IamClient) ShowProjectQuotaInvoker(request *model.ShowProjectQuotaRequest) *ShowProjectQuotaInvoker {
	requestDef := GenReqDefForShowProjectQuota()
	return &ShowProjectQuotaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAgency 修改委托
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改委托。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateAgency(request *model.UpdateAgencyRequest) (*model.UpdateAgencyResponse, error) {
	requestDef := GenReqDefForUpdateAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAgencyResponse), nil
	}
}

// UpdateAgencyInvoker 修改委托
func (c *IamClient) UpdateAgencyInvoker(request *model.UpdateAgencyRequest) *UpdateAgencyInvoker {
	requestDef := GenReqDefForUpdateAgency()
	return &UpdateAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAgencyCustomPolicy 修改委托自定义策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改委托自定义策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateAgencyCustomPolicy(request *model.UpdateAgencyCustomPolicyRequest) (*model.UpdateAgencyCustomPolicyResponse, error) {
	requestDef := GenReqDefForUpdateAgencyCustomPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAgencyCustomPolicyResponse), nil
	}
}

// UpdateAgencyCustomPolicyInvoker 修改委托自定义策略
func (c *IamClient) UpdateAgencyCustomPolicyInvoker(request *model.UpdateAgencyCustomPolicyRequest) *UpdateAgencyCustomPolicyInvoker {
	requestDef := GenReqDefForUpdateAgencyCustomPolicy()
	return &UpdateAgencyCustomPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateCloudServiceCustomPolicy 修改云服务自定义策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改云服务自定义策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateCloudServiceCustomPolicy(request *model.UpdateCloudServiceCustomPolicyRequest) (*model.UpdateCloudServiceCustomPolicyResponse, error) {
	requestDef := GenReqDefForUpdateCloudServiceCustomPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCloudServiceCustomPolicyResponse), nil
	}
}

// UpdateCloudServiceCustomPolicyInvoker 修改云服务自定义策略
func (c *IamClient) UpdateCloudServiceCustomPolicyInvoker(request *model.UpdateCloudServiceCustomPolicyRequest) *UpdateCloudServiceCustomPolicyInvoker {
	requestDef := GenReqDefForUpdateCloudServiceCustomPolicy()
	return &UpdateCloudServiceCustomPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainApiAclPolicy 修改账号接口访问策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改账号接口访问策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateDomainApiAclPolicy(request *model.UpdateDomainApiAclPolicyRequest) (*model.UpdateDomainApiAclPolicyResponse, error) {
	requestDef := GenReqDefForUpdateDomainApiAclPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainApiAclPolicyResponse), nil
	}
}

// UpdateDomainApiAclPolicyInvoker 修改账号接口访问策略
func (c *IamClient) UpdateDomainApiAclPolicyInvoker(request *model.UpdateDomainApiAclPolicyRequest) *UpdateDomainApiAclPolicyInvoker {
	requestDef := GenReqDefForUpdateDomainApiAclPolicy()
	return &UpdateDomainApiAclPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainConsoleAclPolicy 修改账号控制台访问策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改账号控制台访问策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateDomainConsoleAclPolicy(request *model.UpdateDomainConsoleAclPolicyRequest) (*model.UpdateDomainConsoleAclPolicyResponse, error) {
	requestDef := GenReqDefForUpdateDomainConsoleAclPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainConsoleAclPolicyResponse), nil
	}
}

// UpdateDomainConsoleAclPolicyInvoker 修改账号控制台访问策略
func (c *IamClient) UpdateDomainConsoleAclPolicyInvoker(request *model.UpdateDomainConsoleAclPolicyRequest) *UpdateDomainConsoleAclPolicyInvoker {
	requestDef := GenReqDefForUpdateDomainConsoleAclPolicy()
	return &UpdateDomainConsoleAclPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainGroupInheritRole 为用户组授予所有项目服务权限
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)为用户组授予所有项目服务权限。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateDomainGroupInheritRole(request *model.UpdateDomainGroupInheritRoleRequest) (*model.UpdateDomainGroupInheritRoleResponse, error) {
	requestDef := GenReqDefForUpdateDomainGroupInheritRole()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainGroupInheritRoleResponse), nil
	}
}

// UpdateDomainGroupInheritRoleInvoker 为用户组授予所有项目服务权限
func (c *IamClient) UpdateDomainGroupInheritRoleInvoker(request *model.UpdateDomainGroupInheritRoleRequest) *UpdateDomainGroupInheritRoleInvoker {
	requestDef := GenReqDefForUpdateDomainGroupInheritRole()
	return &UpdateDomainGroupInheritRoleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainLoginPolicy 修改账号登录策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改账号登录策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateDomainLoginPolicy(request *model.UpdateDomainLoginPolicyRequest) (*model.UpdateDomainLoginPolicyResponse, error) {
	requestDef := GenReqDefForUpdateDomainLoginPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainLoginPolicyResponse), nil
	}
}

// UpdateDomainLoginPolicyInvoker 修改账号登录策略
func (c *IamClient) UpdateDomainLoginPolicyInvoker(request *model.UpdateDomainLoginPolicyRequest) *UpdateDomainLoginPolicyInvoker {
	requestDef := GenReqDefForUpdateDomainLoginPolicy()
	return &UpdateDomainLoginPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainPasswordPolicy 修改账号密码策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改账号密码策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateDomainPasswordPolicy(request *model.UpdateDomainPasswordPolicyRequest) (*model.UpdateDomainPasswordPolicyResponse, error) {
	requestDef := GenReqDefForUpdateDomainPasswordPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainPasswordPolicyResponse), nil
	}
}

// UpdateDomainPasswordPolicyInvoker 修改账号密码策略
func (c *IamClient) UpdateDomainPasswordPolicyInvoker(request *model.UpdateDomainPasswordPolicyRequest) *UpdateDomainPasswordPolicyInvoker {
	requestDef := GenReqDefForUpdateDomainPasswordPolicy()
	return &UpdateDomainPasswordPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDomainProtectPolicy 修改账号操作保护策略
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改账号操作保护策略。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateDomainProtectPolicy(request *model.UpdateDomainProtectPolicyRequest) (*model.UpdateDomainProtectPolicyResponse, error) {
	requestDef := GenReqDefForUpdateDomainProtectPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDomainProtectPolicyResponse), nil
	}
}

// UpdateDomainProtectPolicyInvoker 修改账号操作保护策略
func (c *IamClient) UpdateDomainProtectPolicyInvoker(request *model.UpdateDomainProtectPolicyRequest) *UpdateDomainProtectPolicyInvoker {
	requestDef := GenReqDefForUpdateDomainProtectPolicy()
	return &UpdateDomainProtectPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateOpenIdConnectConfig 修改OpenId Connect身份提供商配置
//
// 修改OpenId Connect身份提供商配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateOpenIdConnectConfig(request *model.UpdateOpenIdConnectConfigRequest) (*model.UpdateOpenIdConnectConfigResponse, error) {
	requestDef := GenReqDefForUpdateOpenIdConnectConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateOpenIdConnectConfigResponse), nil
	}
}

// UpdateOpenIdConnectConfigInvoker 修改OpenId Connect身份提供商配置
func (c *IamClient) UpdateOpenIdConnectConfigInvoker(request *model.UpdateOpenIdConnectConfigRequest) *UpdateOpenIdConnectConfigInvoker {
	requestDef := GenReqDefForUpdateOpenIdConnectConfig()
	return &UpdateOpenIdConnectConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateProjectStatus 设置项目状态
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)设置项目状态。项目状态包括：正常、冻结。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateProjectStatus(request *model.UpdateProjectStatusRequest) (*model.UpdateProjectStatusResponse, error) {
	requestDef := GenReqDefForUpdateProjectStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateProjectStatusResponse), nil
	}
}

// UpdateProjectStatusInvoker 设置项目状态
func (c *IamClient) UpdateProjectStatusInvoker(request *model.UpdateProjectStatusRequest) *UpdateProjectStatusInvoker {
	requestDef := GenReqDefForUpdateProjectStatus()
	return &UpdateProjectStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePermanentAccessKey 创建永久访问密钥
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)给IAM用户创建永久访问密钥，或IAM用户给自己创建永久访问密钥。
//
// 访问密钥（Access Key ID/Secret Access Key，简称AK/SK），是您通过开发工具（API、CLI、SDK）访问华为云时的身份凭证，不用于登录控制台。系统通过AK识别访问用户的身份，通过SK进行签名验证，通过加密签名验证可以确保请求的机密性、完整性和请求者身份的正确性。在控制台创建访问密钥的方式请参见：[访问密钥](https://support.huaweicloud.com/usermanual-ca/zh-cn_topic_0046606340.html) 。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreatePermanentAccessKey(request *model.CreatePermanentAccessKeyRequest) (*model.CreatePermanentAccessKeyResponse, error) {
	requestDef := GenReqDefForCreatePermanentAccessKey()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePermanentAccessKeyResponse), nil
	}
}

// CreatePermanentAccessKeyInvoker 创建永久访问密钥
func (c *IamClient) CreatePermanentAccessKeyInvoker(request *model.CreatePermanentAccessKeyRequest) *CreatePermanentAccessKeyInvoker {
	requestDef := GenReqDefForCreatePermanentAccessKey()
	return &CreatePermanentAccessKeyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTemporaryAccessKeyByAgency 通过委托获取临时访问密钥
//
// 该接口可以用于通过委托来获取临时访问密钥（临时AK/SK）和securitytoken。
//
// 临时AK/SK和securitytoken是系统颁发给IAM用户的临时访问令牌，有效期为15分钟至24小时，过期后需要重新获取。临时AK/SK和securitytoken遵循权限最小化原则。鉴权时，临时AK/SK和securitytoken必须同时使用，请求头中需要添加“x-security-token”字段，使用方法详情请参考：[API签名参考](https://support.huaweicloud.com/devg-apisign/api-sign-provide.html) 。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateTemporaryAccessKeyByAgency(request *model.CreateTemporaryAccessKeyByAgencyRequest) (*model.CreateTemporaryAccessKeyByAgencyResponse, error) {
	requestDef := GenReqDefForCreateTemporaryAccessKeyByAgency()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTemporaryAccessKeyByAgencyResponse), nil
	}
}

// CreateTemporaryAccessKeyByAgencyInvoker 通过委托获取临时访问密钥
func (c *IamClient) CreateTemporaryAccessKeyByAgencyInvoker(request *model.CreateTemporaryAccessKeyByAgencyRequest) *CreateTemporaryAccessKeyByAgencyInvoker {
	requestDef := GenReqDefForCreateTemporaryAccessKeyByAgency()
	return &CreateTemporaryAccessKeyByAgencyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTemporaryAccessKeyByToken 通过token获取临时访问密钥
//
// 该接口可以用于通过token来获取临时AK/SK和securitytoken。
//
// 临时AK/SK和securitytoken是系统颁发给IAM用户的临时访问令牌，有效期为15分钟至24小时，过期后需要重新获取。临时AK/SK和securitytoken遵循权限最小化原则。鉴权时，临时AK/SK和securitytoken必须同时使用，请求头中需要添加“x-security-token”字段，使用方法详情请参考：[API签名参考](https://support.huaweicloud.com/devg-apisign/api-sign-provide.html)。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateTemporaryAccessKeyByToken(request *model.CreateTemporaryAccessKeyByTokenRequest) (*model.CreateTemporaryAccessKeyByTokenResponse, error) {
	requestDef := GenReqDefForCreateTemporaryAccessKeyByToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTemporaryAccessKeyByTokenResponse), nil
	}
}

// CreateTemporaryAccessKeyByTokenInvoker 通过token获取临时访问密钥
func (c *IamClient) CreateTemporaryAccessKeyByTokenInvoker(request *model.CreateTemporaryAccessKeyByTokenRequest) *CreateTemporaryAccessKeyByTokenInvoker {
	requestDef := GenReqDefForCreateTemporaryAccessKeyByToken()
	return &CreateTemporaryAccessKeyByTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeletePermanentAccessKey 删除指定永久访问密钥
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)删除IAM用户的指定永久访问密钥，或IAM用户删除自己的指定永久访问密钥。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) DeletePermanentAccessKey(request *model.DeletePermanentAccessKeyRequest) (*model.DeletePermanentAccessKeyResponse, error) {
	requestDef := GenReqDefForDeletePermanentAccessKey()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeletePermanentAccessKeyResponse), nil
	}
}

// DeletePermanentAccessKeyInvoker 删除指定永久访问密钥
func (c *IamClient) DeletePermanentAccessKeyInvoker(request *model.DeletePermanentAccessKeyRequest) *DeletePermanentAccessKeyInvoker {
	requestDef := GenReqDefForDeletePermanentAccessKey()
	return &DeletePermanentAccessKeyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPermanentAccessKeys 查询所有永久访问密钥
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户的所有永久访问密钥，或IAM用户查询自己的所有永久访问密钥。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListPermanentAccessKeys(request *model.ListPermanentAccessKeysRequest) (*model.ListPermanentAccessKeysResponse, error) {
	requestDef := GenReqDefForListPermanentAccessKeys()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPermanentAccessKeysResponse), nil
	}
}

// ListPermanentAccessKeysInvoker 查询所有永久访问密钥
func (c *IamClient) ListPermanentAccessKeysInvoker(request *model.ListPermanentAccessKeysRequest) *ListPermanentAccessKeysInvoker {
	requestDef := GenReqDefForListPermanentAccessKeys()
	return &ListPermanentAccessKeysInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowPermanentAccessKey 查询指定永久访问密钥
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户的指定永久访问密钥，或IAM用户查询自己的指定永久访问密钥。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowPermanentAccessKey(request *model.ShowPermanentAccessKeyRequest) (*model.ShowPermanentAccessKeyResponse, error) {
	requestDef := GenReqDefForShowPermanentAccessKey()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowPermanentAccessKeyResponse), nil
	}
}

// ShowPermanentAccessKeyInvoker 查询指定永久访问密钥
func (c *IamClient) ShowPermanentAccessKeyInvoker(request *model.ShowPermanentAccessKeyRequest) *ShowPermanentAccessKeyInvoker {
	requestDef := GenReqDefForShowPermanentAccessKey()
	return &ShowPermanentAccessKeyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePermanentAccessKey 修改指定永久访问密钥
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改IAM用户的指定永久访问密钥，或IAM用户修改自己的指定永久访问密钥。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdatePermanentAccessKey(request *model.UpdatePermanentAccessKeyRequest) (*model.UpdatePermanentAccessKeyResponse, error) {
	requestDef := GenReqDefForUpdatePermanentAccessKey()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePermanentAccessKeyResponse), nil
	}
}

// UpdatePermanentAccessKeyInvoker 修改指定永久访问密钥
func (c *IamClient) UpdatePermanentAccessKeyInvoker(request *model.UpdatePermanentAccessKeyRequest) *UpdatePermanentAccessKeyInvoker {
	requestDef := GenReqDefForUpdatePermanentAccessKey()
	return &UpdatePermanentAccessKeyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateUser 管理员创建IAM用户（推荐）
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)创建IAM用户。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) CreateUser(request *model.CreateUserRequest) (*model.CreateUserResponse, error) {
	requestDef := GenReqDefForCreateUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateUserResponse), nil
	}
}

// CreateUserInvoker 管理员创建IAM用户（推荐）
func (c *IamClient) CreateUserInvoker(request *model.CreateUserRequest) *CreateUserInvoker {
	requestDef := GenReqDefForCreateUser()
	return &CreateUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateUser 管理员创建IAM用户
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)创建IAM用户。IAM用户首次登录时需要修改密码。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateUser(request *model.KeystoneCreateUserRequest) (*model.KeystoneCreateUserResponse, error) {
	requestDef := GenReqDefForKeystoneCreateUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateUserResponse), nil
	}
}

// KeystoneCreateUserInvoker 管理员创建IAM用户
func (c *IamClient) KeystoneCreateUserInvoker(request *model.KeystoneCreateUserRequest) *KeystoneCreateUserInvoker {
	requestDef := GenReqDefForKeystoneCreateUser()
	return &KeystoneCreateUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneDeleteUser 管理员删除IAM用户
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)删除指定IAM用户。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneDeleteUser(request *model.KeystoneDeleteUserRequest) (*model.KeystoneDeleteUserResponse, error) {
	requestDef := GenReqDefForKeystoneDeleteUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneDeleteUserResponse), nil
	}
}

// KeystoneDeleteUserInvoker 管理员删除IAM用户
func (c *IamClient) KeystoneDeleteUserInvoker(request *model.KeystoneDeleteUserRequest) *KeystoneDeleteUserInvoker {
	requestDef := GenReqDefForKeystoneDeleteUser()
	return &KeystoneDeleteUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListGroupsForUser 查询IAM用户所属用户组
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户所属用户组，或IAM用户查询自己所属用户组。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListGroupsForUser(request *model.KeystoneListGroupsForUserRequest) (*model.KeystoneListGroupsForUserResponse, error) {
	requestDef := GenReqDefForKeystoneListGroupsForUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListGroupsForUserResponse), nil
	}
}

// KeystoneListGroupsForUserInvoker 查询IAM用户所属用户组
func (c *IamClient) KeystoneListGroupsForUserInvoker(request *model.KeystoneListGroupsForUserRequest) *KeystoneListGroupsForUserInvoker {
	requestDef := GenReqDefForKeystoneListGroupsForUser()
	return &KeystoneListGroupsForUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListUsers 管理员查询IAM用户列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListUsers(request *model.KeystoneListUsersRequest) (*model.KeystoneListUsersResponse, error) {
	requestDef := GenReqDefForKeystoneListUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListUsersResponse), nil
	}
}

// KeystoneListUsersInvoker 管理员查询IAM用户列表
func (c *IamClient) KeystoneListUsersInvoker(request *model.KeystoneListUsersRequest) *KeystoneListUsersInvoker {
	requestDef := GenReqDefForKeystoneListUsers()
	return &KeystoneListUsersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneListUsersForGroupByAdmin 管理员查询用户组所包含的IAM用户
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询用户组中所包含的IAM用户。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneListUsersForGroupByAdmin(request *model.KeystoneListUsersForGroupByAdminRequest) (*model.KeystoneListUsersForGroupByAdminResponse, error) {
	requestDef := GenReqDefForKeystoneListUsersForGroupByAdmin()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneListUsersForGroupByAdminResponse), nil
	}
}

// KeystoneListUsersForGroupByAdminInvoker 管理员查询用户组所包含的IAM用户
func (c *IamClient) KeystoneListUsersForGroupByAdminInvoker(request *model.KeystoneListUsersForGroupByAdminRequest) *KeystoneListUsersForGroupByAdminInvoker {
	requestDef := GenReqDefForKeystoneListUsersForGroupByAdmin()
	return &KeystoneListUsersForGroupByAdminInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneShowUser 查询IAM用户详情
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户详情，或IAM用户查询自己的用户详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneShowUser(request *model.KeystoneShowUserRequest) (*model.KeystoneShowUserResponse, error) {
	requestDef := GenReqDefForKeystoneShowUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneShowUserResponse), nil
	}
}

// KeystoneShowUserInvoker 查询IAM用户详情
func (c *IamClient) KeystoneShowUserInvoker(request *model.KeystoneShowUserRequest) *KeystoneShowUserInvoker {
	requestDef := GenReqDefForKeystoneShowUser()
	return &KeystoneShowUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneUpdateUserByAdmin 管理员修改IAM用户信息
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改IAM用户信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneUpdateUserByAdmin(request *model.KeystoneUpdateUserByAdminRequest) (*model.KeystoneUpdateUserByAdminResponse, error) {
	requestDef := GenReqDefForKeystoneUpdateUserByAdmin()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneUpdateUserByAdminResponse), nil
	}
}

// KeystoneUpdateUserByAdminInvoker 管理员修改IAM用户信息
func (c *IamClient) KeystoneUpdateUserByAdminInvoker(request *model.KeystoneUpdateUserByAdminRequest) *KeystoneUpdateUserByAdminInvoker {
	requestDef := GenReqDefForKeystoneUpdateUserByAdmin()
	return &KeystoneUpdateUserByAdminInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneUpdateUserPassword 修改IAM用户密码
//
// 该接口可以用于IAM用户修改自己的密码。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneUpdateUserPassword(request *model.KeystoneUpdateUserPasswordRequest) (*model.KeystoneUpdateUserPasswordResponse, error) {
	requestDef := GenReqDefForKeystoneUpdateUserPassword()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneUpdateUserPasswordResponse), nil
	}
}

// KeystoneUpdateUserPasswordInvoker 修改IAM用户密码
func (c *IamClient) KeystoneUpdateUserPasswordInvoker(request *model.KeystoneUpdateUserPasswordRequest) *KeystoneUpdateUserPasswordInvoker {
	requestDef := GenReqDefForKeystoneUpdateUserPassword()
	return &KeystoneUpdateUserPasswordInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUserLoginProtects 查询IAM用户的登录保护状态信息列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户的登录保护状态列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListUserLoginProtects(request *model.ListUserLoginProtectsRequest) (*model.ListUserLoginProtectsResponse, error) {
	requestDef := GenReqDefForListUserLoginProtects()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUserLoginProtectsResponse), nil
	}
}

// ListUserLoginProtectsInvoker 查询IAM用户的登录保护状态信息列表
func (c *IamClient) ListUserLoginProtectsInvoker(request *model.ListUserLoginProtectsRequest) *ListUserLoginProtectsInvoker {
	requestDef := GenReqDefForListUserLoginProtects()
	return &ListUserLoginProtectsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUserMfaDevices 查询IAM用户的MFA绑定信息列表
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户的MFA绑定信息列表。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ListUserMfaDevices(request *model.ListUserMfaDevicesRequest) (*model.ListUserMfaDevicesResponse, error) {
	requestDef := GenReqDefForListUserMfaDevices()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUserMfaDevicesResponse), nil
	}
}

// ListUserMfaDevicesInvoker 查询IAM用户的MFA绑定信息列表
func (c *IamClient) ListUserMfaDevicesInvoker(request *model.ListUserMfaDevicesRequest) *ListUserMfaDevicesInvoker {
	requestDef := GenReqDefForListUserMfaDevices()
	return &ListUserMfaDevicesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowUser 查询IAM用户详情（推荐）
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询IAM用户详情，或IAM用户查询自己的详情。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowUser(request *model.ShowUserRequest) (*model.ShowUserResponse, error) {
	requestDef := GenReqDefForShowUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowUserResponse), nil
	}
}

// ShowUserInvoker 查询IAM用户详情（推荐）
func (c *IamClient) ShowUserInvoker(request *model.ShowUserRequest) *ShowUserInvoker {
	requestDef := GenReqDefForShowUser()
	return &ShowUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowUserLoginProtect 查询指定IAM用户的登录保护状态信息
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询指定IAM用户的登录保护状态信息，或IAM用户查询自己的登录保护状态信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowUserLoginProtect(request *model.ShowUserLoginProtectRequest) (*model.ShowUserLoginProtectResponse, error) {
	requestDef := GenReqDefForShowUserLoginProtect()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowUserLoginProtectResponse), nil
	}
}

// ShowUserLoginProtectInvoker 查询指定IAM用户的登录保护状态信息
func (c *IamClient) ShowUserLoginProtectInvoker(request *model.ShowUserLoginProtectRequest) *ShowUserLoginProtectInvoker {
	requestDef := GenReqDefForShowUserLoginProtect()
	return &ShowUserLoginProtectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowUserMfaDevice 查询指定IAM用户的MFA绑定信息
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)查询指定IAM用户的MFA绑定信息，或IAM用户查询自己的MFA绑定信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) ShowUserMfaDevice(request *model.ShowUserMfaDeviceRequest) (*model.ShowUserMfaDeviceResponse, error) {
	requestDef := GenReqDefForShowUserMfaDevice()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowUserMfaDeviceResponse), nil
	}
}

// ShowUserMfaDeviceInvoker 查询指定IAM用户的MFA绑定信息
func (c *IamClient) ShowUserMfaDeviceInvoker(request *model.ShowUserMfaDeviceRequest) *ShowUserMfaDeviceInvoker {
	requestDef := GenReqDefForShowUserMfaDevice()
	return &ShowUserMfaDeviceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateLoginProtect 修改IAM用户登录保护状态信息
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改账号操作保护。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateLoginProtect(request *model.UpdateLoginProtectRequest) (*model.UpdateLoginProtectResponse, error) {
	requestDef := GenReqDefForUpdateLoginProtect()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateLoginProtectResponse), nil
	}
}

// UpdateLoginProtectInvoker 修改IAM用户登录保护状态信息
func (c *IamClient) UpdateLoginProtectInvoker(request *model.UpdateLoginProtectRequest) *UpdateLoginProtectInvoker {
	requestDef := GenReqDefForUpdateLoginProtect()
	return &UpdateLoginProtectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateUser 管理员修改IAM用户信息（推荐）
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)修改IAM用户信息 。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateUser(request *model.UpdateUserRequest) (*model.UpdateUserResponse, error) {
	requestDef := GenReqDefForUpdateUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateUserResponse), nil
	}
}

// UpdateUserInvoker 管理员修改IAM用户信息（推荐）
func (c *IamClient) UpdateUserInvoker(request *model.UpdateUserRequest) *UpdateUserInvoker {
	requestDef := GenReqDefForUpdateUser()
	return &UpdateUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateUserInformation 修改IAM用户信息（推荐）
//
// 该接口可以用于IAM用户修改自己的用户信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) UpdateUserInformation(request *model.UpdateUserInformationRequest) (*model.UpdateUserInformationResponse, error) {
	requestDef := GenReqDefForUpdateUserInformation()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateUserInformationResponse), nil
	}
}

// UpdateUserInformationInvoker 修改IAM用户信息（推荐）
func (c *IamClient) UpdateUserInformationInvoker(request *model.UpdateUserInformationRequest) *UpdateUserInformationInvoker {
	requestDef := GenReqDefForUpdateUserInformation()
	return &UpdateUserInformationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateAgencyToken 获取委托Token
//
// 该接口可以用于获取委托方的token。
//
// 例如：A账号希望B账号管理自己的某些资源，所以A账号创建了委托给B账号，则A账号为委托方，B账号为被委托方。那么B账号可以通过该接口获取委托token。B账号仅能使用该token管理A账号的委托资源，不能管理自己账号中的资源。如果B账号需要管理自己账号中的资源，则需要获取自己的用户token。
//
// token是系统颁发给用户的访问令牌，承载用户的身份、权限等信息。调用IAM以及其他云服务的接口时，可以使用本接口获取的token进行鉴权。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。如果使用全局区域的Endpoint调用，该token可以在所有区域使用；如果使用非全局区域的Endpoint调用，则该token仅在该区域生效，不能跨区域使用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// &gt; - token的有效期为24小时，建议进行缓存，避免频繁调用。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateAgencyToken(request *model.KeystoneCreateAgencyTokenRequest) (*model.KeystoneCreateAgencyTokenResponse, error) {
	requestDef := GenReqDefForKeystoneCreateAgencyToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateAgencyTokenResponse), nil
	}
}

// KeystoneCreateAgencyTokenInvoker 获取委托Token
func (c *IamClient) KeystoneCreateAgencyTokenInvoker(request *model.KeystoneCreateAgencyTokenRequest) *KeystoneCreateAgencyTokenInvoker {
	requestDef := GenReqDefForKeystoneCreateAgencyToken()
	return &KeystoneCreateAgencyTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateUserTokenByPassword 获取IAM用户Token（使用密码）
//
// 该接口可以用于通过用户名/密码的方式进行认证来获取IAM用户token。
//
// token是系统颁发给IAM用户的访问令牌，承载用户的身份、权限等信息。调用IAM以及其他云服务的接口时，可以使用本接口获取的IAM用户token进行鉴权。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。如果使用全局区域的Endpoint调用，该token可以在所有区域使用；如果使用非全局区域的Endpoint调用，则该token仅在该区域生效，不能跨区域使用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
// &gt; - token的有效期为24小时，建议进行缓存，避免频繁调用。
// &gt; - 通过Postman获取用户token示例请参见：[如何通过Postman获取用户token](https://support.huaweicloud.com/iam_faq/iam_01_034.html)。
// &gt; - 如果需要获取具有Security Administrator权限的token，请参见：[IAM 常见问题](https://support.huaweicloud.com/iam_faq/iam_01_0608.html)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateUserTokenByPassword(request *model.KeystoneCreateUserTokenByPasswordRequest) (*model.KeystoneCreateUserTokenByPasswordResponse, error) {
	requestDef := GenReqDefForKeystoneCreateUserTokenByPassword()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateUserTokenByPasswordResponse), nil
	}
}

// KeystoneCreateUserTokenByPasswordInvoker 获取IAM用户Token（使用密码）
func (c *IamClient) KeystoneCreateUserTokenByPasswordInvoker(request *model.KeystoneCreateUserTokenByPasswordRequest) *KeystoneCreateUserTokenByPasswordInvoker {
	requestDef := GenReqDefForKeystoneCreateUserTokenByPassword()
	return &KeystoneCreateUserTokenByPasswordInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneCreateUserTokenByPasswordAndMfa 获取IAM用户Token（使用密码+虚拟MFA）
//
// 该接口可以用于通过用户名/密码+虚拟MFA的方式进行认证，在IAM用户开启了的登录保护功能，并选择通过虚拟MFA验证时获取IAM用户token。
//
// token是系统颁发给用户的访问令牌，承载用户的身份、权限等信息。调用IAM以及其他云服务的接口时，可以使用本接口获取的token进行鉴权。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。如果使用全局区域的Endpoint调用，该token可以在所有区域使用；如果使用非全局区域的Endpoint调用，则该token仅在该区域生效，不能跨区域使用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
// &gt; - token的有效期为24小时，建议进行缓存，避免频繁调用。
// &gt; - 通过Postman获取用户token示例请参见：[如何通过Postman获取用户token](https://support.huaweicloud.com/iam_faq/iam_01_034.html)。
// &gt; - 如果需要获取具有Security Administrator权限的token，请参见：[IAM 常见问题](https://support.huaweicloud.com/iam_faq/iam_01_0608.html)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneCreateUserTokenByPasswordAndMfa(request *model.KeystoneCreateUserTokenByPasswordAndMfaRequest) (*model.KeystoneCreateUserTokenByPasswordAndMfaResponse, error) {
	requestDef := GenReqDefForKeystoneCreateUserTokenByPasswordAndMfa()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneCreateUserTokenByPasswordAndMfaResponse), nil
	}
}

// KeystoneCreateUserTokenByPasswordAndMfaInvoker 获取IAM用户Token（使用密码+虚拟MFA）
func (c *IamClient) KeystoneCreateUserTokenByPasswordAndMfaInvoker(request *model.KeystoneCreateUserTokenByPasswordAndMfaRequest) *KeystoneCreateUserTokenByPasswordAndMfaInvoker {
	requestDef := GenReqDefForKeystoneCreateUserTokenByPasswordAndMfa()
	return &KeystoneCreateUserTokenByPasswordAndMfaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// KeystoneValidateToken 校验Token的有效性
//
// 该接口可以用于[管理员](https://support.huaweicloud.com/usermanual-iam/iam_01_0001.html)校验本账号中IAM用户token的有效性，或IAM用户校验自己token的有效性。管理员仅能校验本账号中IAM用户token的有效性，不能校验其他账号中IAM用户token的有效性。如果被校验的token有效，则返回该token的详细信息。
//
// 该接口可以使用全局区域的Endpoint和其他区域的Endpoint调用。IAM的Endpoint请参见：[地区和终端节点](https://developer.huaweicloud.com/endpoint?IAM)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *IamClient) KeystoneValidateToken(request *model.KeystoneValidateTokenRequest) (*model.KeystoneValidateTokenResponse, error) {
	requestDef := GenReqDefForKeystoneValidateToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.KeystoneValidateTokenResponse), nil
	}
}

// KeystoneValidateTokenInvoker 校验Token的有效性
func (c *IamClient) KeystoneValidateTokenInvoker(request *model.KeystoneValidateTokenRequest) *KeystoneValidateTokenInvoker {
	requestDef := GenReqDefForKeystoneValidateToken()
	return &KeystoneValidateTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
