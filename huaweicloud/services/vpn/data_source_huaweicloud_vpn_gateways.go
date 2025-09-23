// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product VPN
// ---------------------------------------------------------------

package vpn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN GET /v5/{project_id}/vpn-gateways
func DataSourceGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGatewaysRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the gateway.`,
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the gateway.`,
			},
			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the network type of the gateway.`,
			},
			"attachment_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the attachment type of the gateway.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID of the gateway.`,
			},
			"gateways": {
				Type:        schema.TypeList,
				Elem:        gatewayGatewaysSchema(),
				Computed:    true,
				Description: `The list of gateways.`,
			},
		},
	}
}

func gatewayGatewaysSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the gateway`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the gateway.`,
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The network type of the gateway.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the gateway.`,
			},
			"attachment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attachment type.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC to which the VPN gateway is connected.`,
			},
			"er_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the ER to which the VPN gateway is connected.`,
			},
			"er_attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ER attachment ID.`,
			},
			"local_subnets": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The local subnets.`,
			},
			"connect_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VPC network segment used by the VPN gateway`,
			},
			"bgp_asn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The ASN number of BGP`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor of the VPN gateway.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The availability zone IDs.`,
			},
			"connection_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The max number of connections.`,
			},
			"used_connection_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of used connections.`,
			},
			"used_connection_group": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of used connection groups.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID`,
			},
			"eips": {
				Type:     schema.TypeList,
				Elem:     gatewayGatewaysResponseEipSchema(),
				Computed: true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"access_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the access VPC.`,
			},
			"access_subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the access subnet.`,
			},
			"access_private_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of private access IPs.`,
			},
			"ha_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HA mode.`,
			},
		},
	}
	return &sc
}

func gatewayGatewaysResponseEipSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth_billing_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The bandwidth billing info.`,
			},
			"bandwidth_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The bandwidth ID.`,
			},
			"bandwidth_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The bandwidth name.`,
			},
			"bandwidth_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Bandwidth size in Mbit/s.`,
			},
			"billing_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The billing info.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The charge mode of the bandwidth. The value can be **bandwidth** and **traffic**.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP ID.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP address.`,
			},
			"ip_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The public IP version.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The EIP type. The value can be **5_bgp** and **5_sbgp**.`,
			},
		},
	}
	return &sc
}

func dataSourceGatewaysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGateways: Query the VPN gateways
	var (
		getGatewaysHttpUrl = "v5/{project_id}/vpn-gateways"
		getGatewaysProduct = "vpn"
	)
	getGatewaysClient, err := cfg.NewServiceClient(getGatewaysProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getGatewaysPath := getGatewaysClient.Endpoint + getGatewaysHttpUrl
	getGatewaysPath = strings.ReplaceAll(getGatewaysPath, "{project_id}", getGatewaysClient.ProjectID)

	getGatewaysqueryParams := buildGetGatewaysQueryParams(d, cfg)
	getGatewaysPath += getGatewaysqueryParams

	getGatewaysOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getGatewaysResp, err := getGatewaysClient.Request("GET", getGatewaysPath, &getGatewaysOpt)

	if err != nil {
		return diag.Errorf("error retrieving VPN gateway: %s", err)
	}

	getGatewaysRespBody, err := utils.FlattenResponse(getGatewaysResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("gateways", filterGetGatewaysResponseBodyGateways(
			flattenGetGatewaysResponseBodyGateways(getGatewaysRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGatewaysResponseBodyGateways(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("vpn_gateways", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"network_type":           utils.PathSearch("network_type", v, nil),
			"status":                 utils.PathSearch("status", v, nil),
			"attachment_type":        utils.PathSearch("attachment_type", v, nil),
			"vpc_id":                 utils.PathSearch("vpc_id", v, nil),
			"er_id":                  utils.PathSearch("er_id", v, nil),
			"er_attachment_id":       utils.PathSearch("er_attachment_id", v, nil),
			"local_subnets":          utils.PathSearch("local_subnets", v, nil),
			"connect_subnet":         utils.PathSearch("connect_subnet", v, nil),
			"bgp_asn":                utils.PathSearch("bgp_asn", v, nil),
			"flavor":                 utils.PathSearch("flavor", v, nil),
			"availability_zones":     utils.PathSearch("availability_zone_ids", v, nil),
			"connection_number":      utils.PathSearch("connection_number", v, nil),
			"used_connection_number": utils.PathSearch("used_connection_number", v, nil),
			"used_connection_group":  utils.PathSearch("used_connection_group", v, nil),
			"enterprise_project_id":  utils.PathSearch("enterprise_project_id", v, nil),
			"eips":                   flattenGatewaysResponseEips(v),
			"created_at":             utils.PathSearch("created_at", v, nil),
			"updated_at":             utils.PathSearch("updated_at", v, nil),
			"access_vpc_id":          utils.PathSearch("access_vpc_id", v, nil),
			"access_subnet_id":       utils.PathSearch("access_subnet_id", v, nil),
			"access_private_ips":     utils.PathSearch("access_private_ips", v, nil),
			"ha_mode":                utils.PathSearch("ha_mode", v, nil),
		})
	}
	return rst
}

func flattenGatewaysResponseEips(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("[eip1,eip2]", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"bandwidth_billing_info": utils.PathSearch("bandwidth_billing_info", v, nil),
			"bandwidth_id":           utils.PathSearch("bandwidth_id", v, nil),
			"bandwidth_name":         utils.PathSearch("bandwidth_name", v, nil),
			"bandwidth_size":         utils.PathSearch("bandwidth_size", v, nil),
			"billing_info":           utils.PathSearch("billing_info", v, nil),
			"charge_mode":            utils.PathSearch("charge_mode", v, nil),
			"id":                     utils.PathSearch("id", v, nil),
			"ip_address":             utils.PathSearch("ip_address", v, nil),
			"ip_version":             utils.PathSearch("ip_version", v, nil),
			"type":                   utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func filterGetGatewaysResponseBodyGateways(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("gateway_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("network_type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("network_type", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("attachment_type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("attachment_type", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildGetGatewaysQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	enterpriseProjectID := cfg.GetEnterpriseProjectID(d)
	if enterpriseProjectID != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", enterpriseProjectID)
	}
	return ""
}
