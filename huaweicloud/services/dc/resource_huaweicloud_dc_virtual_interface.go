package dc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC DELETE /v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}
// @API DC GET /v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}
// @API DC PUT /v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}
// @API DC POST /v3/{project_id}/dcaas/virtual-interfaces
func ResourceVirtualInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualInterfaceCreate,
		ReadContext:   resourceVirtualInterfaceRead,
		UpdateContext: resourceVirtualInterfaceUpdate,
		DeleteContext: resourceVirtualInterfaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the virtual interface is located.",
			},
			"direct_connect_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the direct connection associated with the virtual interface.",
			},
			"vgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the virtual gateway to which the virtual interface is connected.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the virtual interface.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the virtual interface.",
			},
			"route_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The route mode of the virtual interface.",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The VLAN for constom side.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ingress bandwidth size of the virtual interface.",
			},
			"remote_ep_group": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The CIDR list of remote subnets.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the virtual interface.",
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The service type of the virtual interface.",
			},
			"local_gateway_v4_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"remote_gateway_v4_ip"},
				Description:  "The IPv4 address of the virtual interface in cloud side.",
			},
			"remote_gateway_v4_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"remote_gateway_v6_ip"},
				Description:   "The IPv4 address of the virtual interface in client side.",
			},
			"address_family": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The address family type of the virtual interface.",
			},
			"local_gateway_v6_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"remote_gateway_v6_ip"},
				ExactlyOneOf: []string{"local_gateway_v4_ip"},
				Description:  "The IPv6 address of the virtual interface in cloud side.",
			},
			"remote_gateway_v6_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The IPv6 address of the virtual interface in client side.",
			},
			"asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The local BGP ASN in client side.",
			},
			"bgp_md5": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The (MD5) password for the local BGP.",
			},
			"enable_bfd": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the Bidirectional Forwarding Detection (BFD) function.",
			},
			"enable_nqa": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the Network Quality Analysis (NQA) function.",
			},
			"lag_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the link aggregation group (LAG) associated with the virtual interface.",
			},
			"resource_tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The project ID of another tenant which is used to create virtual interface across tenant.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The enterprise project ID to which the virtual interface belongs.",
			},
			"tags": common.TagsSchema(),

			// Attributes
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The attributed device ID.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the virtual interface.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the virtual interface.",
			},
			"vif_peers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        vifPeersSchema(),
				Description: "The peer information of the virtual interface.",
			},
		},
	}
}

func vifPeersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VIF peer resource ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the virtual interface peer.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the virtual interface peer.`,
			},
			"address_family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The address family type of the virtual interface, which can be IPv4 or IPv6.`,
			},
			"local_gateway_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The address of the virtual interface peer used on the cloud.`,
			},
			"remote_gateway_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The address of the virtual interface peer used in the on-premises data center.`,
			},
			"route_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The routing mode, which can be static or bgp.`,
			},
			"bgp_asn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The ASN of the BGP peer.`,
			},
			"bgp_md5": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The MD5 password of the BGP peer.`,
			},
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the device that the virtual interface peer belongs to.`,
			},
			"enable_bfd": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable BFD.`,
			},
			"enable_nqa": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable NQA.`,
			},
			"bgp_route_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The BGP route configuration.`,
			},
			"bgp_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The BGP protocol status of the virtual interface peer.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the virtual interface peer.`,
			},
			"vif_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the virtual interface corresponding to the virtual interface peer.`,
			},
			"receive_route_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of received BGP routes if BGP routing is used.`,
			},
			"remote_ep_group": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The remote subnet list, which records the CIDR blocks used in the on-premises data center.`,
			},
		},
	}
	return &sc
}

// The default value of this field is false when it is created, so the false value is treated as nil when it is created.
func buildCreateVirtualInterfaceEnableBfd(d *schema.ResourceData) interface{} {
	if d.Get("enable_bfd").(bool) {
		return true
	}

	return nil
}

// The default value of this field is false when it is created, so the false value is treated as nil when it is created.
func buildCreateVirtualInterfaceEnableNqa(d *schema.ResourceData) interface{} {
	if d.Get("enable_nqa").(bool) {
		return true
	}

	return nil
}

func buildCreateVirtualInterfaceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vgw_id":                d.Get("vgw_id"),
		"type":                  d.Get("type"),
		"route_mode":            d.Get("route_mode"),
		"vlan":                  d.Get("vlan"),
		"bandwidth":             d.Get("bandwidth"),
		"remote_ep_group":       d.Get("remote_ep_group"),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"direct_connect_id":     utils.ValueIgnoreEmpty(d.Get("direct_connect_id")),
		"service_type":          utils.ValueIgnoreEmpty(d.Get("service_type")),
		"local_gateway_v4_ip":   utils.ValueIgnoreEmpty(d.Get("local_gateway_v4_ip")),
		"remote_gateway_v4_ip":  utils.ValueIgnoreEmpty(d.Get("remote_gateway_v4_ip")),
		"address_family":        utils.ValueIgnoreEmpty(d.Get("address_family")),
		"local_gateway_v6_ip":   utils.ValueIgnoreEmpty(d.Get("local_gateway_v6_ip")),
		"remote_gateway_v6_ip":  utils.ValueIgnoreEmpty(d.Get("remote_gateway_v6_ip")),
		"bgp_asn":               utils.ValueIgnoreEmpty(d.Get("asn")),
		"bgp_md5":               utils.ValueIgnoreEmpty(d.Get("bgp_md5")),
		"enable_bfd":            buildCreateVirtualInterfaceEnableBfd(d),
		"enable_nqa":            buildCreateVirtualInterfaceEnableNqa(d),
		"lag_id":                utils.ValueIgnoreEmpty(d.Get("lag_id")),
		"resource_tenant_id":    utils.ValueIgnoreEmpty(d.Get("resource_tenant_id")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}

	return map[string]interface{}{
		"virtual_interface": bodyParams,
	}
}

func resourceVirtualInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dcaas/virtual-interfaces"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateVirtualInterfaceBodyParams(d, cfg)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DC virtual interface: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("virtual_interface.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating DC virtual interface: ID is not found in API response")
	}
	d.SetId(id)

	// create tags
	if err := utils.CreateResourceTags(client, d, "dc-vif", d.Id()); err != nil {
		return diag.Errorf("error creating tags of DC virtual interface (%s): %s", d.Id(), err)
	}

	return resourceVirtualInterfaceRead(ctx, d, meta)
}

func resourceVirtualInterfaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{interfaceId}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// When the interface does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving DC virtual interface")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	interfaceResp := utils.PathSearch("virtual_interface", respBody, nil)
	if interfaceResp == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("vgw_id", utils.PathSearch("vgw_id", interfaceResp, nil)),
		d.Set("type", utils.PathSearch("type", interfaceResp, nil)),
		d.Set("route_mode", utils.PathSearch("route_mode", interfaceResp, nil)),
		d.Set("vlan", utils.PathSearch("vlan", interfaceResp, nil)),
		d.Set("bandwidth", utils.PathSearch("bandwidth", interfaceResp, nil)),
		d.Set("remote_ep_group", utils.PathSearch("remote_ep_group", interfaceResp, nil)),
		d.Set("name", utils.PathSearch("name", interfaceResp, nil)),
		d.Set("description", utils.PathSearch("description", interfaceResp, nil)),
		d.Set("direct_connect_id", utils.PathSearch("direct_connect_id", interfaceResp, nil)),
		d.Set("service_type", utils.PathSearch("service_type", interfaceResp, nil)),
		d.Set("local_gateway_v4_ip", utils.PathSearch("local_gateway_v4_ip", interfaceResp, nil)),
		d.Set("remote_gateway_v4_ip", utils.PathSearch("remote_gateway_v4_ip", interfaceResp, nil)),
		d.Set("address_family", utils.PathSearch("address_family", interfaceResp, nil)),
		d.Set("local_gateway_v6_ip", utils.PathSearch("local_gateway_v6_ip", interfaceResp, nil)),
		d.Set("remote_gateway_v6_ip", utils.PathSearch("remote_gateway_v6_ip", interfaceResp, nil)),
		d.Set("asn", utils.PathSearch("bgp_asn", interfaceResp, nil)),
		d.Set("bgp_md5", utils.PathSearch("bgp_md5", interfaceResp, nil)),
		d.Set("enable_bfd", utils.PathSearch("enable_bfd", interfaceResp, nil)),
		d.Set("enable_nqa", utils.PathSearch("enable_nqa", interfaceResp, nil)),
		d.Set("lag_id", utils.PathSearch("lag_id", interfaceResp, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", interfaceResp, nil)),
		d.Set("device_id", utils.PathSearch("device_id", interfaceResp, nil)),
		d.Set("status", utils.PathSearch("status", interfaceResp, nil)),
		d.Set("created_at", utils.PathSearch("create_time", interfaceResp, nil)),
		d.Set("vif_peers", flattenVifPeersAttribute(utils.PathSearch("vif_peers", interfaceResp, make([]interface{}, 0)).([]interface{}))),
		utils.SetResourceTagsToState(d, client, "dc-vif", d.Id()),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVifPeersAttribute(peersArray []interface{}) []interface{} {
	if len(peersArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(peersArray))
	for _, v := range peersArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"address_family":    utils.PathSearch("address_family", v, nil),
			"local_gateway_ip":  utils.PathSearch("local_gateway_ip", v, nil),
			"remote_gateway_ip": utils.PathSearch("remote_gateway_ip", v, nil),
			"route_mode":        utils.PathSearch("route_mode", v, nil),
			"bgp_asn":           utils.PathSearch("bgp_asn", v, nil),
			"bgp_md5":           utils.PathSearch("bgp_md5", v, nil),
			"device_id":         utils.PathSearch("device_id", v, nil),
			"enable_bfd":        utils.PathSearch("enable_bfd", v, nil),
			"enable_nqa":        utils.PathSearch("enable_nqa", v, nil),
			"bgp_route_limit":   utils.PathSearch("bgp_route_limit", v, nil),
			"bgp_status":        utils.PathSearch("bgp_status", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"vif_id":            utils.PathSearch("vif_id", v, nil),
			"receive_route_num": utils.PathSearch("receive_route_num", v, nil),
			"remote_ep_group":   utils.PathSearch("remote_ep_group", v, nil),
		})
	}
	return rst
}

func buildUpdateVirtualInterfaceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":            utils.ValueIgnoreEmpty(d.Get("name")),
		"description":     d.Get("description"),
		"bandwidth":       utils.ValueIgnoreEmpty(d.Get("bandwidth")),
		"remote_ep_group": utils.ValueIgnoreEmpty(d.Get("remote_ep_group")),
	}

	return map[string]interface{}{
		"virtual_interface": bodyParams,
	}
}

func closeVirtualInterfaceNetworkDetection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	bodyParams := make(map[string]interface{})
	// At the same time, only one of BFD and NQA is enabled.
	if d.HasChange("enable_bfd") && !d.Get("enable_bfd").(bool) {
		bodyParams["enable_bfd"] = false
	} else if d.HasChange("enable_nqa") && !d.Get("enable_nqa").(bool) {
		bodyParams["enable_nqa"] = false
	}

	if len(bodyParams) == 0 {
		return nil
	}

	requestPath := client.Endpoint + "v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{interfaceId}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"virtual_interface": bodyParams,
		},
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error closing network detection of the virtual interface (%s): %s", d.Id(), err)
	}
	return nil
}

func openVirtualInterfaceNetworkDetection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	bodyParams := make(map[string]interface{})
	detectionOpened := false
	if d.HasChange("enable_bfd") && d.Get("enable_bfd").(bool) {
		detectionOpened = true
		bodyParams["enable_bfd"] = true
	}

	if d.HasChange("enable_nqa") && d.Get("enable_nqa").(bool) {
		// The enable requests of BFD and NQA cannot be sent at the same time.
		if detectionOpened {
			return errors.New("BFD and NQA cannot be enabled at the same time")
		}
		bodyParams["enable_nqa"] = true
	}

	if len(bodyParams) == 0 {
		return nil
	}

	requestPath := client.Endpoint + "v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{interfaceId}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"virtual_interface": bodyParams,
		},
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error opening network detection of the virtual interface (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceVirtualInterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	if d.HasChanges("name", "description", "bandwidth", "remote_ep_group") {
		requestPath := client.Endpoint + "v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{interfaceId}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateVirtualInterfaceBodyParams(d)),
		}

		_, err := client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating DC virtual interface: %s", err)
		}
	}

	if d.HasChanges("enable_bfd", "enable_nqa") {
		// BFD and NQA cannot be enabled at the same time.
		// When BFD (NQA) is enabled and NQA (BFD) is disabled, we need to disable BFD (NQA) first, and then enable NQA (BFD).
		// If the disable and enable requests are sent at the same time, an error will be reported.
		if err = closeVirtualInterfaceNetworkDetection(client, d); err != nil {
			return diag.FromErr(err)
		}
		if err = openVirtualInterfaceNetworkDetection(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// update tags
	tagErr := utils.UpdateResourceTags(client, d, "dc-vif", d.Id())
	if tagErr != nil {
		return diag.Errorf("error updating tags of DC virtual interface %s: %s", d.Id(), tagErr)
	}

	return resourceVirtualInterfaceRead(ctx, d, meta)
}

func resourceVirtualInterfaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{interfaceId}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		// When the interface does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting DC virtual interface")
	}

	return nil
}
