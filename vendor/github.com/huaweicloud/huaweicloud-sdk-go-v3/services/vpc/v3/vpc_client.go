package v3

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
)

type VpcClient struct {
	HcClient *http_client.HcHttpClient
}

func NewVpcClient(hcClient *http_client.HcHttpClient) *VpcClient {
	return &VpcClient{HcClient: hcClient}
}

func VpcClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// BatchCreateSubNetworkInterface 批量创建辅助弹性网卡
//
// 批量创建辅助弹性网卡
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) BatchCreateSubNetworkInterface(request *model.BatchCreateSubNetworkInterfaceRequest) (*model.BatchCreateSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForBatchCreateSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchCreateSubNetworkInterfaceResponse), nil
	}
}

// BatchCreateSubNetworkInterfaceInvoker 批量创建辅助弹性网卡
func (c *VpcClient) BatchCreateSubNetworkInterfaceInvoker(request *model.BatchCreateSubNetworkInterfaceRequest) *BatchCreateSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForBatchCreateSubNetworkInterface()
	return &BatchCreateSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSecurityGroup 创建安全组
//
// 创建安全组
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) CreateSecurityGroup(request *model.CreateSecurityGroupRequest) (*model.CreateSecurityGroupResponse, error) {
	requestDef := GenReqDefForCreateSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSecurityGroupResponse), nil
	}
}

// CreateSecurityGroupInvoker 创建安全组
func (c *VpcClient) CreateSecurityGroupInvoker(request *model.CreateSecurityGroupRequest) *CreateSecurityGroupInvoker {
	requestDef := GenReqDefForCreateSecurityGroup()
	return &CreateSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSecurityGroupRule 创建安全组规则
//
// 创建安全组规则
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) CreateSecurityGroupRule(request *model.CreateSecurityGroupRuleRequest) (*model.CreateSecurityGroupRuleResponse, error) {
	requestDef := GenReqDefForCreateSecurityGroupRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSecurityGroupRuleResponse), nil
	}
}

// CreateSecurityGroupRuleInvoker 创建安全组规则
func (c *VpcClient) CreateSecurityGroupRuleInvoker(request *model.CreateSecurityGroupRuleRequest) *CreateSecurityGroupRuleInvoker {
	requestDef := GenReqDefForCreateSecurityGroupRule()
	return &CreateSecurityGroupRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSubNetworkInterface 创建辅助弹性网卡
//
// 创建辅助弹性网卡
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) CreateSubNetworkInterface(request *model.CreateSubNetworkInterfaceRequest) (*model.CreateSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForCreateSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSubNetworkInterfaceResponse), nil
	}
}

// CreateSubNetworkInterfaceInvoker 创建辅助弹性网卡
func (c *VpcClient) CreateSubNetworkInterfaceInvoker(request *model.CreateSubNetworkInterfaceRequest) *CreateSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForCreateSubNetworkInterface()
	return &CreateSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSecurityGroup 删除安全组
//
// 删除安全组
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) DeleteSecurityGroup(request *model.DeleteSecurityGroupRequest) (*model.DeleteSecurityGroupResponse, error) {
	requestDef := GenReqDefForDeleteSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSecurityGroupResponse), nil
	}
}

// DeleteSecurityGroupInvoker 删除安全组
func (c *VpcClient) DeleteSecurityGroupInvoker(request *model.DeleteSecurityGroupRequest) *DeleteSecurityGroupInvoker {
	requestDef := GenReqDefForDeleteSecurityGroup()
	return &DeleteSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSecurityGroupRule 删除安全组规则
//
// 删除安全组规则
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) DeleteSecurityGroupRule(request *model.DeleteSecurityGroupRuleRequest) (*model.DeleteSecurityGroupRuleResponse, error) {
	requestDef := GenReqDefForDeleteSecurityGroupRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSecurityGroupRuleResponse), nil
	}
}

// DeleteSecurityGroupRuleInvoker 删除安全组规则
func (c *VpcClient) DeleteSecurityGroupRuleInvoker(request *model.DeleteSecurityGroupRuleRequest) *DeleteSecurityGroupRuleInvoker {
	requestDef := GenReqDefForDeleteSecurityGroupRule()
	return &DeleteSecurityGroupRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSubNetworkInterface 删除辅助弹性网卡
//
// 删除辅助弹性网卡
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) DeleteSubNetworkInterface(request *model.DeleteSubNetworkInterfaceRequest) (*model.DeleteSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForDeleteSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSubNetworkInterfaceResponse), nil
	}
}

// DeleteSubNetworkInterfaceInvoker 删除辅助弹性网卡
func (c *VpcClient) DeleteSubNetworkInterfaceInvoker(request *model.DeleteSubNetworkInterfaceRequest) *DeleteSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForDeleteSubNetworkInterface()
	return &DeleteSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSecurityGroupRules 查询安全组规则列表
//
// 查询安全组规则列表
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ListSecurityGroupRules(request *model.ListSecurityGroupRulesRequest) (*model.ListSecurityGroupRulesResponse, error) {
	requestDef := GenReqDefForListSecurityGroupRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSecurityGroupRulesResponse), nil
	}
}

// ListSecurityGroupRulesInvoker 查询安全组规则列表
func (c *VpcClient) ListSecurityGroupRulesInvoker(request *model.ListSecurityGroupRulesRequest) *ListSecurityGroupRulesInvoker {
	requestDef := GenReqDefForListSecurityGroupRules()
	return &ListSecurityGroupRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSecurityGroups 查询安全组列表
//
// 查询某租户下的安全组列表
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ListSecurityGroups(request *model.ListSecurityGroupsRequest) (*model.ListSecurityGroupsResponse, error) {
	requestDef := GenReqDefForListSecurityGroups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSecurityGroupsResponse), nil
	}
}

// ListSecurityGroupsInvoker 查询安全组列表
func (c *VpcClient) ListSecurityGroupsInvoker(request *model.ListSecurityGroupsRequest) *ListSecurityGroupsInvoker {
	requestDef := GenReqDefForListSecurityGroups()
	return &ListSecurityGroupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSubNetworkInterfaces 查询租户下辅助弹性网卡列表
//
// 查询辅助弹性网卡列表，单次查询最多返回2000条数据
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ListSubNetworkInterfaces(request *model.ListSubNetworkInterfacesRequest) (*model.ListSubNetworkInterfacesResponse, error) {
	requestDef := GenReqDefForListSubNetworkInterfaces()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSubNetworkInterfacesResponse), nil
	}
}

// ListSubNetworkInterfacesInvoker 查询租户下辅助弹性网卡列表
func (c *VpcClient) ListSubNetworkInterfacesInvoker(request *model.ListSubNetworkInterfacesRequest) *ListSubNetworkInterfacesInvoker {
	requestDef := GenReqDefForListSubNetworkInterfaces()
	return &ListSubNetworkInterfacesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// MigrateSubNetworkInterface 迁移辅助弹性网卡
//
// 批量迁移辅助弹性网卡
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) MigrateSubNetworkInterface(request *model.MigrateSubNetworkInterfaceRequest) (*model.MigrateSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForMigrateSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.MigrateSubNetworkInterfaceResponse), nil
	}
}

// MigrateSubNetworkInterfaceInvoker 迁移辅助弹性网卡
func (c *VpcClient) MigrateSubNetworkInterfaceInvoker(request *model.MigrateSubNetworkInterfaceRequest) *MigrateSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForMigrateSubNetworkInterface()
	return &MigrateSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSecurityGroup 查询安全组
//
// 查询单个安全组详情
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ShowSecurityGroup(request *model.ShowSecurityGroupRequest) (*model.ShowSecurityGroupResponse, error) {
	requestDef := GenReqDefForShowSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSecurityGroupResponse), nil
	}
}

// ShowSecurityGroupInvoker 查询安全组
func (c *VpcClient) ShowSecurityGroupInvoker(request *model.ShowSecurityGroupRequest) *ShowSecurityGroupInvoker {
	requestDef := GenReqDefForShowSecurityGroup()
	return &ShowSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSecurityGroupRule 查询安全组规则
//
// 查询单个安全组规则
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ShowSecurityGroupRule(request *model.ShowSecurityGroupRuleRequest) (*model.ShowSecurityGroupRuleResponse, error) {
	requestDef := GenReqDefForShowSecurityGroupRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSecurityGroupRuleResponse), nil
	}
}

// ShowSecurityGroupRuleInvoker 查询安全组规则
func (c *VpcClient) ShowSecurityGroupRuleInvoker(request *model.ShowSecurityGroupRuleRequest) *ShowSecurityGroupRuleInvoker {
	requestDef := GenReqDefForShowSecurityGroupRule()
	return &ShowSecurityGroupRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSubNetworkInterface 查询租户下辅助弹性网卡
//
// 查询辅助弹性网卡详情
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ShowSubNetworkInterface(request *model.ShowSubNetworkInterfaceRequest) (*model.ShowSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForShowSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSubNetworkInterfaceResponse), nil
	}
}

// ShowSubNetworkInterfaceInvoker 查询租户下辅助弹性网卡
func (c *VpcClient) ShowSubNetworkInterfaceInvoker(request *model.ShowSubNetworkInterfaceRequest) *ShowSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForShowSubNetworkInterface()
	return &ShowSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSubNetworkInterfacesQuantity 查询租户下辅助弹性网卡数目
//
// 查询辅助弹性网卡数目
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ShowSubNetworkInterfacesQuantity(request *model.ShowSubNetworkInterfacesQuantityRequest) (*model.ShowSubNetworkInterfacesQuantityResponse, error) {
	requestDef := GenReqDefForShowSubNetworkInterfacesQuantity()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSubNetworkInterfacesQuantityResponse), nil
	}
}

// ShowSubNetworkInterfacesQuantityInvoker 查询租户下辅助弹性网卡数目
func (c *VpcClient) ShowSubNetworkInterfacesQuantityInvoker(request *model.ShowSubNetworkInterfacesQuantityRequest) *ShowSubNetworkInterfacesQuantityInvoker {
	requestDef := GenReqDefForShowSubNetworkInterfacesQuantity()
	return &ShowSubNetworkInterfacesQuantityInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateSecurityGroup 更新安全组
//
// 更新安全组
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) UpdateSecurityGroup(request *model.UpdateSecurityGroupRequest) (*model.UpdateSecurityGroupResponse, error) {
	requestDef := GenReqDefForUpdateSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateSecurityGroupResponse), nil
	}
}

// UpdateSecurityGroupInvoker 更新安全组
func (c *VpcClient) UpdateSecurityGroupInvoker(request *model.UpdateSecurityGroupRequest) *UpdateSecurityGroupInvoker {
	requestDef := GenReqDefForUpdateSecurityGroup()
	return &UpdateSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateSubNetworkInterface 更新辅助弹性网卡
//
// 更新辅助弹性网卡
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) UpdateSubNetworkInterface(request *model.UpdateSubNetworkInterfaceRequest) (*model.UpdateSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForUpdateSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateSubNetworkInterfaceResponse), nil
	}
}

// UpdateSubNetworkInterfaceInvoker 更新辅助弹性网卡
func (c *VpcClient) UpdateSubNetworkInterfaceInvoker(request *model.UpdateSubNetworkInterfaceRequest) *UpdateSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForUpdateSubNetworkInterface()
	return &UpdateSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAddressGroup 创建地址组
//
// 创建地址组
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) CreateAddressGroup(request *model.CreateAddressGroupRequest) (*model.CreateAddressGroupResponse, error) {
	requestDef := GenReqDefForCreateAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAddressGroupResponse), nil
	}
}

// CreateAddressGroupInvoker 创建地址组
func (c *VpcClient) CreateAddressGroupInvoker(request *model.CreateAddressGroupRequest) *CreateAddressGroupInvoker {
	requestDef := GenReqDefForCreateAddressGroup()
	return &CreateAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAddressGroup 删除地址组
//
// 删除地址组，非强制删除，删除前请确保未被其他资源引用
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) DeleteAddressGroup(request *model.DeleteAddressGroupRequest) (*model.DeleteAddressGroupResponse, error) {
	requestDef := GenReqDefForDeleteAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAddressGroupResponse), nil
	}
}

// DeleteAddressGroupInvoker 删除地址组
func (c *VpcClient) DeleteAddressGroupInvoker(request *model.DeleteAddressGroupRequest) *DeleteAddressGroupInvoker {
	requestDef := GenReqDefForDeleteAddressGroup()
	return &DeleteAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteIpAddressGroupForce 强制删除地址组
//
// 强制删除地址组，删除的地址组与安全组规则关联时，会删除地址组与关联的安全组规则。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) DeleteIpAddressGroupForce(request *model.DeleteIpAddressGroupForceRequest) (*model.DeleteIpAddressGroupForceResponse, error) {
	requestDef := GenReqDefForDeleteIpAddressGroupForce()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteIpAddressGroupForceResponse), nil
	}
}

// DeleteIpAddressGroupForceInvoker 强制删除地址组
func (c *VpcClient) DeleteIpAddressGroupForceInvoker(request *model.DeleteIpAddressGroupForceRequest) *DeleteIpAddressGroupForceInvoker {
	requestDef := GenReqDefForDeleteIpAddressGroupForce()
	return &DeleteIpAddressGroupForceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAddressGroup 查询地址组列表
//
// 查询地址组列表，根据过滤条件进行过滤。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ListAddressGroup(request *model.ListAddressGroupRequest) (*model.ListAddressGroupResponse, error) {
	requestDef := GenReqDefForListAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAddressGroupResponse), nil
	}
}

// ListAddressGroupInvoker 查询地址组列表
func (c *VpcClient) ListAddressGroupInvoker(request *model.ListAddressGroupRequest) *ListAddressGroupInvoker {
	requestDef := GenReqDefForListAddressGroup()
	return &ListAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAddressGroup 查询地址组
//
// 查询地址组详情。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ShowAddressGroup(request *model.ShowAddressGroupRequest) (*model.ShowAddressGroupResponse, error) {
	requestDef := GenReqDefForShowAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAddressGroupResponse), nil
	}
}

// ShowAddressGroupInvoker 查询地址组
func (c *VpcClient) ShowAddressGroupInvoker(request *model.ShowAddressGroupRequest) *ShowAddressGroupInvoker {
	requestDef := GenReqDefForShowAddressGroup()
	return &ShowAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAddressGroup 更新地址组
//
// 更新地址组。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) UpdateAddressGroup(request *model.UpdateAddressGroupRequest) (*model.UpdateAddressGroupResponse, error) {
	requestDef := GenReqDefForUpdateAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAddressGroupResponse), nil
	}
}

// UpdateAddressGroupInvoker 更新地址组
func (c *VpcClient) UpdateAddressGroupInvoker(request *model.UpdateAddressGroupRequest) *UpdateAddressGroupInvoker {
	requestDef := GenReqDefForUpdateAddressGroup()
	return &UpdateAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddVpcExtendCidr 添加VPC扩展网段
//
// 添加VPC的扩展网段
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) AddVpcExtendCidr(request *model.AddVpcExtendCidrRequest) (*model.AddVpcExtendCidrResponse, error) {
	requestDef := GenReqDefForAddVpcExtendCidr()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddVpcExtendCidrResponse), nil
	}
}

// AddVpcExtendCidrInvoker 添加VPC扩展网段
func (c *VpcClient) AddVpcExtendCidrInvoker(request *model.AddVpcExtendCidrRequest) *AddVpcExtendCidrInvoker {
	requestDef := GenReqDefForAddVpcExtendCidr()
	return &AddVpcExtendCidrInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateVpc 创建VPC
//
// 创建虚拟私有云
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) CreateVpc(request *model.CreateVpcRequest) (*model.CreateVpcResponse, error) {
	requestDef := GenReqDefForCreateVpc()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateVpcResponse), nil
	}
}

// CreateVpcInvoker 创建VPC
func (c *VpcClient) CreateVpcInvoker(request *model.CreateVpcRequest) *CreateVpcInvoker {
	requestDef := GenReqDefForCreateVpc()
	return &CreateVpcInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteVpc 删除VPC
//
// 删除VPC
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) DeleteVpc(request *model.DeleteVpcRequest) (*model.DeleteVpcResponse, error) {
	requestDef := GenReqDefForDeleteVpc()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteVpcResponse), nil
	}
}

// DeleteVpcInvoker 删除VPC
func (c *VpcClient) DeleteVpcInvoker(request *model.DeleteVpcRequest) *DeleteVpcInvoker {
	requestDef := GenReqDefForDeleteVpc()
	return &DeleteVpcInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListVpcs 查询VPC列表
//
// 查询vpc列表
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ListVpcs(request *model.ListVpcsRequest) (*model.ListVpcsResponse, error) {
	requestDef := GenReqDefForListVpcs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListVpcsResponse), nil
	}
}

// ListVpcsInvoker 查询VPC列表
func (c *VpcClient) ListVpcsInvoker(request *model.ListVpcsRequest) *ListVpcsInvoker {
	requestDef := GenReqDefForListVpcs()
	return &ListVpcsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveVpcExtendCidr 移除VPC扩展网段
//
// 移除VPC扩展网段
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) RemoveVpcExtendCidr(request *model.RemoveVpcExtendCidrRequest) (*model.RemoveVpcExtendCidrResponse, error) {
	requestDef := GenReqDefForRemoveVpcExtendCidr()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveVpcExtendCidrResponse), nil
	}
}

// RemoveVpcExtendCidrInvoker 移除VPC扩展网段
func (c *VpcClient) RemoveVpcExtendCidrInvoker(request *model.RemoveVpcExtendCidrRequest) *RemoveVpcExtendCidrInvoker {
	requestDef := GenReqDefForRemoveVpcExtendCidr()
	return &RemoveVpcExtendCidrInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowVpc 查询VPC详情
//
// 查询vpc详情
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) ShowVpc(request *model.ShowVpcRequest) (*model.ShowVpcResponse, error) {
	requestDef := GenReqDefForShowVpc()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowVpcResponse), nil
	}
}

// ShowVpcInvoker 查询VPC详情
func (c *VpcClient) ShowVpcInvoker(request *model.ShowVpcRequest) *ShowVpcInvoker {
	requestDef := GenReqDefForShowVpc()
	return &ShowVpcInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateVpc 更新VPC
//
// 更新vpc
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *VpcClient) UpdateVpc(request *model.UpdateVpcRequest) (*model.UpdateVpcResponse, error) {
	requestDef := GenReqDefForUpdateVpc()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateVpcResponse), nil
	}
}

// UpdateVpcInvoker 更新VPC
func (c *VpcClient) UpdateVpcInvoker(request *model.UpdateVpcRequest) *UpdateVpcInvoker {
	requestDef := GenReqDefForUpdateVpc()
	return &UpdateVpcInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
