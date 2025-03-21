package dc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC GET /v3/{project_id}/dcaas/virtual-interfaces
func DataSourceDCVirtualInterfaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDCVirtualInterfacesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_interface_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the virtual interface.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the virtual interface.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the status of the virtual interface.",
			},
			"direct_connect_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the direct connection associated with the virtual interface.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
			"vgw_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the virtual gateway for the virtual interface.",
			},
			"virtual_interfaces": {
				Type:        schema.TypeList,
				Elem:        virtualInterfaceSchema(),
				Computed:    true,
				Description: "Specifies the virtual interface list.",
			},
		},
	}
}

func virtualInterfaceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the virtual interface.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the virtual interface.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bandwidth of the virtual interface.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the virtual interface.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the virtual interface.",
			},
			"direct_connect_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the direct connection associated with the virtual interface.",
			},
			"service_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of access gateway with the virtual interface.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the virtual interface.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the virtual interface.",
			},
			"vgw_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the virtual gateway for the virtual interface.",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The VLAN connected to the user gateway of the virtual interface.",
			},
			"enable_nqa": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Does the enable nqa functionality of virtual interface.",
			},
			"enable_bfd": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Does the enable bfd functionality of virtual interface.",
			},
			"lag_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link aggregation group ID associated with vif of the virtual interface.",
			},
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The belong device ID of the virtual interface.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project that the virtual interface belongs to.",
			},
			"local_gateway_v4_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud side gateway IPv4 interface address of the virtual interface.",
			},
			"remote_gateway_v4_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The customer side gateway IPv4 interface address of the virtual interface.",
			},
			"address_family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address cluster type of the interface.",
			},
			"local_gateway_v6_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud side gateway IPv6 interface address of the virtual interface.",
			},
			"remote_gateway_v6_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The customer side gateway IPv6 interface address of the virtual interface.",
			},
			"remote_ep_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of remote subnets, recording the cidrs on the tenant side.",
			},
			"route_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The route mode of the virtual interface.",
			},
			"asn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The (ASN) number for the local BGP.",
			},
			"bgp_md5": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The (MD5) password for the local BGP.",
			},
			"vif_peers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        datasourceVifPeersSchema(),
				Description: "The peer information of the virtual interface.",
			},
		},
	}
	return &sc
}

func datasourceVifPeersSchema() *schema.Resource {
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

func resourceDCVirtualInterfacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listVirtualInterfaces: Query the List of DC virtual interfaces.
	var (
		listVirtualInterfacesHttpUrl = "v3/{project_id}/dcaas/virtual-interfaces"
		listVirtualInterfacesProduct = "dc"
	)
	listVirtualInterfacesClient, err := cfg.NewServiceClient(listVirtualInterfacesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	listVirtualInterfacesPath := listVirtualInterfacesClient.Endpoint + listVirtualInterfacesHttpUrl
	listVirtualInterfacesPath = strings.ReplaceAll(listVirtualInterfacesPath, "{project_id}", listVirtualInterfacesClient.ProjectID)

	listVirtualInterfacesQueryParams := buildListVirtualInterfacesQueryParams(d, cfg.GetEnterpriseProjectID(d))
	listVirtualInterfacesPath += listVirtualInterfacesQueryParams

	listVirtualInterfacesResp, err := pagination.ListAllItems(
		listVirtualInterfacesClient,
		"marker",
		listVirtualInterfacesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving DC virtual interfaces: %s", err)
	}

	listVirtualInterfacesRespJson, err := json.Marshal(listVirtualInterfacesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listVirtualInterfacesRespBody interface{}
	err = json.Unmarshal(listVirtualInterfacesRespJson, &listVirtualInterfacesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("virtual_interfaces", filterListVirtualInterfacesBody(
			flattenListVirtualInterfacesBody(listVirtualInterfacesRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListVirtualInterfacesBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("virtual_interfaces", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"bandwidth":             utils.PathSearch("bandwidth", v, nil),
			"created_at":            utils.PathSearch("create_time", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"direct_connect_id":     utils.PathSearch("direct_connect_id", v, nil),
			"service_type":          utils.PathSearch("service_type", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"type":                  utils.PathSearch("type", v, nil),
			"vgw_id":                utils.PathSearch("vgw_id", v, nil),
			"vlan":                  utils.PathSearch("vlan", v, nil),
			"enable_nqa":            utils.PathSearch("enable_nqa", v, nil),
			"enable_bfd":            utils.PathSearch("enable_bfd", v, nil),
			"lag_id":                utils.PathSearch("lag_id", v, nil),
			"device_id":             utils.PathSearch("device_id", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"local_gateway_v4_ip":   utils.PathSearch("local_gateway_v4_ip", v, nil),
			"remote_gateway_v4_ip":  utils.PathSearch("remote_gateway_v4_ip", v, nil),
			"address_family":        utils.PathSearch("address_family", v, nil),
			"local_gateway_v6_ip":   utils.PathSearch("local_gateway_v6_ip", v, nil),
			"remote_gateway_v6_ip":  utils.PathSearch("remote_gateway_v6_ip", v, nil),
			"remote_ep_group":       utils.PathSearch("remote_ep_group", v, nil),
			"route_mode":            utils.PathSearch("route_mode", v, nil),
			"asn":                   utils.PathSearch("bgp_asn", v, nil),
			"bgp_md5":               utils.PathSearch("bgp_md5", v, nil),
			"vif_peers":             flattenVifPeersAttributes(v),
		})
	}
	return rst
}

func flattenVifPeersAttributes(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("vif_peers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
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

func filterListVirtualInterfacesBody(all []interface{}, d *schema.ResourceData) []interface{} {
	name := d.Get("name").(string)
	if name == "" {
		return all
	}

	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if name != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListVirtualInterfacesQueryParams(d *schema.ResourceData, enterpriseProjectId string) string {
	res := ""
	if v, ok := d.GetOk("virtual_interface_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if v, ok := d.GetOk("direct_connect_id"); ok {
		res = fmt.Sprintf("%s&direct_connect_id=%v", res, v)
	}

	if v, ok := d.GetOk("vgw_id"); ok {
		res = fmt.Sprintf("%s&vgw_id=%v", res, v)
	}

	if enterpriseProjectId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
