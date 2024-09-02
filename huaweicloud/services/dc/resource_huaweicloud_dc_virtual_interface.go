package dc

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dc/v3/interfaces"

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

func buildVirtualInterfaceCreateOpts(d *schema.ResourceData, cfg *config.Config) interfaces.CreateOpts {
	return interfaces.CreateOpts{
		VgwId:               d.Get("vgw_id").(string),
		Type:                d.Get("type").(string),
		RouteMode:           d.Get("route_mode").(string),
		Vlan:                d.Get("vlan").(int),
		Bandwidth:           d.Get("bandwidth").(int),
		RemoteEpGroup:       utils.ExpandToStringList(d.Get("remote_ep_group").([]interface{})),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		DirectConnectId:     d.Get("direct_connect_id").(string),
		ServiceType:         d.Get("service_type").(string),
		LocalGatewayV4Ip:    d.Get("local_gateway_v4_ip").(string),
		RemoteGatewayV4Ip:   d.Get("remote_gateway_v4_ip").(string),
		AddressFamily:       d.Get("address_family").(string),
		LocalGatewayV6Ip:    d.Get("local_gateway_v6_ip").(string),
		RemoteGatewayV6Ip:   d.Get("remote_gateway_v6_ip").(string),
		BgpAsn:              d.Get("asn").(int),
		BgpMd5:              d.Get("bgp_md5").(string),
		EnableBfd:           d.Get("enable_bfd").(bool),
		EnableNqa:           d.Get("enable_nqa").(bool),
		LagId:               d.Get("lag_id").(string),
		ResourceTenantId:    d.Get("resource_tenant_id").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
}

func resourceVirtualInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	opts := buildVirtualInterfaceCreateOpts(d, cfg)
	resp, err := interfaces.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating virtual interface: %s", err)
	}
	d.SetId(resp.ID)

	// create tags
	if err := utils.CreateResourceTags(client, d, "dc-vif", d.Id()); err != nil {
		return diag.Errorf("error setting tags of DC virtual interface %s: %s", d.Id(), err)
	}

	return resourceVirtualInterfaceRead(ctx, d, meta)
}

func resourceVirtualInterfaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	interfaceId := d.Id()
	resp, err := interfaces.Get(client, interfaceId)
	if err != nil {
		// When the interface does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving DC virtual interface")
	}
	log.Printf("[DEBUG] The response of virtual interface is: %#v", resp)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("vgw_id", resp.VgwId),
		d.Set("type", resp.Type),
		d.Set("route_mode", resp.RouteMode),
		d.Set("vlan", resp.Vlan),
		d.Set("bandwidth", resp.Bandwidth),
		d.Set("remote_ep_group", resp.RemoteEpGroup),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("direct_connect_id", resp.DirectConnectId),
		d.Set("service_type", resp.ServiceType),
		d.Set("local_gateway_v4_ip", resp.LocalGatewayV4Ip),
		d.Set("remote_gateway_v4_ip", resp.RemoteGatewayV4Ip),
		d.Set("address_family", resp.AddressFamily),
		d.Set("local_gateway_v6_ip", resp.LocalGatewayV6Ip),
		d.Set("remote_gateway_v6_ip", resp.RemoteGatewayV6Ip),
		d.Set("asn", resp.BgpAsn),
		d.Set("bgp_md5", resp.BgpMd5),
		d.Set("enable_bfd", resp.EnableBfd),
		d.Set("enable_nqa", resp.EnableNqa),
		d.Set("lag_id", resp.LagId),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("device_id", resp.DeviceId),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedAt),
		d.Set("vif_peers", flattenVifPeers(resp.VifPeers)),
		utils.SetResourceTagsToState(d, client, "dc-vif", d.Id()),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving virtual interface fields: %s", err)
	}
	return nil
}

func flattenVifPeers(vifPeers []interfaces.VifPeer) []interface{} {
	if vifPeers == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(vifPeers))
	for _, v := range vifPeers {
		rst = append(rst, map[string]interface{}{
			"id":                v.ID,
			"name":              v.Name,
			"description":       v.Description,
			"address_family":    v.AddressFamily,
			"local_gateway_ip":  v.LocalGatewayIp,
			"remote_gateway_ip": v.RemoteGatewayIp,
			"route_mode":        v.RouteMode,
			"bgp_asn":           v.BgpAsn,
			"bgp_md5":           v.BgpMd5,
			"device_id":         v.DeviceId,
			"enable_bfd":        v.EnableBfd,
			"enable_nqa":        v.EnableNqa,
			"bgp_route_limit":   v.BgpRouteLimit,
			"bgp_status":        v.BgpStatus,
			"status":            v.Status,
			"vif_id":            v.VifId,
			"receive_route_num": v.ReceiveRouteNum,
			"remote_ep_group":   v.RemoteEpGroup,
		})
	}
	return rst
}

func closeVirtualInterfaceNetworkDetection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		interfaceId = d.Id()
		opts        = interfaces.UpdateOpts{}
	)

	// At the same time, only one of BFD and NQA is enabled.
	if d.HasChange("enable_bfd") && !d.Get("enable_bfd").(bool) {
		opts.EnableBfd = utils.Bool(false)
	} else if d.HasChange("enable_nqa") && !d.Get("enable_nqa").(bool) {
		opts.EnableNqa = utils.Bool(false)
	}
	if reflect.DeepEqual(opts, interfaces.UpdateOpts{}) {
		return nil
	}

	_, err := interfaces.Update(client, interfaceId, opts)
	if err != nil {
		return fmt.Errorf("error closing network detection of the virtual interface (%s): %s", interfaceId, err)
	}
	return nil
}

func openVirtualInterfaceNetworkDetection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		interfaceId     = d.Id()
		detectionOpened = false
		opts            = interfaces.UpdateOpts{}
	)

	if d.HasChange("enable_bfd") && d.Get("enable_bfd").(bool) {
		detectionOpened = true
		opts.EnableBfd = utils.Bool(true)
	}
	if d.HasChange("enable_nqa") && d.Get("enable_nqa").(bool) {
		// The enable requests of BFD and NQA cannot be sent at the same time.
		if detectionOpened {
			return fmt.Errorf("BFD and NQA cannot be enabled at the same time")
		}
		opts.EnableNqa = utils.Bool(true)
	}
	if reflect.DeepEqual(opts, interfaces.UpdateOpts{}) {
		return nil
	}

	_, err := interfaces.Update(client, interfaceId, opts)
	if err != nil {
		return fmt.Errorf("error opening network detection of the virtual interface (%s): %s", interfaceId, err)
	}
	return nil
}

func resourceVirtualInterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	if d.HasChanges("name", "description", "bandwidth", "remote_ep_group") {
		var (
			interfaceId = d.Id()

			opts = interfaces.UpdateOpts{
				Name:          d.Get("name").(string),
				Description:   utils.String(d.Get("description").(string)),
				Bandwidth:     d.Get("bandwidth").(int),
				RemoteEpGroup: utils.ExpandToStringList(d.Get("remote_ep_group").([]interface{})),
			}
		)

		_, err := interfaces.Update(client, interfaceId, opts)
		if err != nil {
			return diag.Errorf("error closing network detection of the virtual interface (%s): %s", interfaceId, err)
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
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	interfaceId := d.Id()
	err = interfaces.Delete(client, interfaceId)
	if err != nil {
		// When the interface does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting DC virtual interface")
	}

	return nil
}
