package vpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v3/{project_id}/vpc/sub-network-interfaces
// @API VPC GET /v3/{project_id}/vpc/sub-network-interfaces/{sub_network_interface_id}
// @API VPC PUT /v3/{project_id}/vpc/sub-network-interfaces/{sub_network_interface_id}
// @API VPC DELETE /v3/{project_id}/vpc/sub-network-interfaces/{sub_network_interface_id}
func ResourceSubNetworkInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubNetworkInterfaceCreate,
		ReadContext:   resourceSubNetworkInterfaceRead,
		UpdateContext: resourceSubNetworkInterfaceUpdate,
		DeleteContext: resourceSubNetworkInterfaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv6_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateSubNetworkInterfaceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sub_network_interface": map[string]interface{}{
			"virsubnet_id":       d.Get("subnet_id"),
			"parent_id":          d.Get("parent_id"),
			"security_groups":    utils.ValueIgnoreEmpty(d.Get("security_group_ids")),
			"description":        utils.ValueIgnoreEmpty(d.Get("description")),
			"vlan_id":            utils.ValueIgnoreEmpty(d.Get("vlan_id")),
			"private_ip_address": utils.ValueIgnoreEmpty(d.Get("ip_address")),
			"ipv6_enable":        utils.ValueIgnoreEmpty(d.Get("ipv6_enable")),
			"ipv6_ip_address":    utils.ValueIgnoreEmpty(d.Get("ipv6_ip_address")),
		},
	}
	return bodyParams
}

func resourceSubNetworkInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	createSubNetworkInterfaceHttpUrl := "vpc/sub-network-interfaces"
	createSubNetworkInterfacePath := client.ResourceBaseURL() + createSubNetworkInterfaceHttpUrl

	createSubNetworkInterfaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createSubNetworkInterfaceOpt.JSONBody = utils.RemoveNil(buildCreateSubNetworkInterfaceBodyParams(d))
	createSubNetworkInterfaceResp, err := client.Request("POST", createSubNetworkInterfacePath, &createSubNetworkInterfaceOpt)
	if err != nil {
		return diag.Errorf("error creating supplementary network interface: %s", err)
	}

	createcreateSubNetworkInterfaceRespBody, err := utils.FlattenResponse(createSubNetworkInterfaceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("sub_network_interface.id", createcreateSubNetworkInterfaceRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating supplementary network interface: %s is not found in API response", id)
	}

	d.SetId(id)

	return resourceSubNetworkInterfaceRead(ctx, d, meta)
}

func getSubNetworkInterfaceInfo(d *schema.ResourceData, meta interface{}) (*http.Response, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v3 client: %s", err)
	}

	getSubNetworkInterfaceHttpUrl := "vpc/sub-network-interfaces/" + d.Id()
	getSubNetworkInterfacePath := client.ResourceBaseURL() + getSubNetworkInterfaceHttpUrl

	getSubNetworkInterfaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	resp, err := client.Request("GET", getSubNetworkInterfacePath, &getSubNetworkInterfaceOpt)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func resourceSubNetworkInterfaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	getSubNetworkInterfaceResp, err := getSubNetworkInterfaceInfo(d, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC supplementary network interface")
	}

	getSubNetworkInterfaceRespBody, err := utils.FlattenResponse(getSubNetworkInterfaceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("subnet_id", utils.PathSearch("sub_network_interface.virsubnet_id", getSubNetworkInterfaceRespBody, nil)),
		d.Set("parent_id", utils.PathSearch("sub_network_interface.parent_id", getSubNetworkInterfaceRespBody, nil)),
		d.Set("security_group_ids", utils.PathSearch("sub_network_interface.security_groups", getSubNetworkInterfaceRespBody, nil)),
		d.Set("description", utils.PathSearch("sub_network_interface.description", getSubNetworkInterfaceRespBody, nil)),
		d.Set("vlan_id", fmt.Sprintf("%v", utils.PathSearch("sub_network_interface.vlan_id", getSubNetworkInterfaceRespBody, nil))),
		d.Set("ip_address", utils.PathSearch("sub_network_interface.private_ip_address", getSubNetworkInterfaceRespBody, nil)),
		d.Set("ipv6_enable", utils.PathSearch("sub_network_interface.ipv6_enable", getSubNetworkInterfaceRespBody, nil)),
		d.Set("ipv6_ip_address", utils.PathSearch("sub_network_interface.ipv6_ip_address", getSubNetworkInterfaceRespBody, nil)),
		d.Set("mac_address", utils.PathSearch("sub_network_interface.mac_address", getSubNetworkInterfaceRespBody, nil)),
		d.Set("parent_device_id", utils.PathSearch("sub_network_interface.parent_device_id", getSubNetworkInterfaceRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("sub_network_interface.vpc_id", getSubNetworkInterfaceRespBody, nil)),
		d.Set("status", utils.PathSearch("sub_network_interface.state", getSubNetworkInterfaceRespBody, nil)),
		d.Set("created_at", utils.PathSearch("sub_network_interface.created_at", getSubNetworkInterfaceRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateSubNetworkInterfaceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sub_network_interface": map[string]interface{}{
			"description":     d.Get("description"),
			"security_groups": utils.ValueIgnoreEmpty(d.Get("security_group_ids")),
		},
	}
	return bodyParams
}

func resourceSubNetworkInterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 Client: %s", err)
	}

	if d.HasChanges("security_group_ids", "description") {
		updateSubNetworkInterfaceHttpUrl := "vpc/sub-network-interfaces/" + d.Id()
		updateSubNetworkInterfacePath := client.ResourceBaseURL() + updateSubNetworkInterfaceHttpUrl

		updateSubNetworkInterfaceOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateSubNetworkInterfaceOpts.JSONBody = utils.RemoveNil(buildUpdateSubNetworkInterfaceBodyParams(d))
		_, err = client.Request("PUT", updateSubNetworkInterfacePath, &updateSubNetworkInterfaceOpts)
		if err != nil {
			return diag.Errorf("error updating VPC supplementary network interface: %s", err)
		}
	}
	return resourceSubNetworkInterfaceRead(ctx, d, meta)
}

func resourceSubNetworkInterfaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 Client: %s", err)
	}

	deleteSubNetworkInterfaceHttpUrl := "vpc/sub-network-interfaces/" + d.Id()
	deleteSubNetworkInterfacePath := client.ResourceBaseURL() + deleteSubNetworkInterfaceHttpUrl

	deleteSubNetworkInterfaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deleteSubNetworkInterfacePath, &deleteSubNetworkInterfaceOpt)
	if err != nil {
		return diag.Errorf("error deleting supplementary network interface: %s", err)
	}

	return nil
}
