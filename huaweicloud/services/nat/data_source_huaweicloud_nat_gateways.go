package nat

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

// @API NAT GET /v2/{project_id}/nat_gateways
func DataSourcePublicGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicGatewaysRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the NAT gateways are located.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the NAT gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The specification of the NAT gateways.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The current status of the NAT gateways.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the VPC to which the NAT gateways belong.",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The network ID of the downstream interface (the next hop of the DVR)" +
					"of the NAT gateways.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the NAT gateways belong.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the NAT gateway.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The creation time of the NAT gateway.",
			},
			"gateways": {
				Type:        schema.TypeList,
				Elem:        gatewayPublicGatewaysSchema(),
				Computed:    true,
				Description: "The list of the NAT gateway.",
			},
		},
	}
}

func gatewayPublicGatewaysSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the NAT gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The specification of the NAT gateway.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the NAT gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the NAT gateway.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the NAT gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the VPC to which the NAT gateway belongs.",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The network ID of the downstream interface (the next hop of the DVR)" +
					"of the NAT gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the NAT gateway belongs.",
			},
			"session_conf": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The session configuration of the NAT gateway.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tcp_session_expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The TCP session expiration time, in seconds.",
						},
						"udp_session_expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The UDP session expiration time, in seconds.",
						},
						"icmp_session_expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ICMP session expiration time, in seconds.",
						},
						"tcp_time_wait_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The duration of TIME_WAIT state when TCP connection is closed, in seconds.",
						},
					},
				},
			},
			"ngport_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private IP address of the NAT gateway.",
			},
			"billing_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The order information of the NAT gateway.",
			},
			"dnat_rules_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of DNAT rules on the NAT gateway.",
			},
			"snat_rule_public_ip_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of SNAT rules on the NAT gateway.",
			},
			"pps_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of packets that the NAT gateway can receive or send per second.",
			},
			"bps_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bandwidth that the NAT gateway can receive or send per second.",
			},
		},
	}
	return &sc
}

func dataSourcePublicGatewaysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v2/{project_id}/nat_gateways"
		product = "nat"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListPublicGatewaysQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving NAT gateways %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("gateways", flattenListGatewaysResponseBodyPublicGateways(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListGatewaysResponseBodyPublicGateways(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("nat_gateways", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                        utils.PathSearch("id", v, nil),
			"name":                      utils.PathSearch("name", v, nil),
			"spec":                      utils.PathSearch("spec", v, nil),
			"status":                    utils.PathSearch("status", v, nil),
			"description":               utils.PathSearch("description", v, nil),
			"created_at":                utils.PathSearch("created_at", v, nil),
			"vpc_id":                    utils.PathSearch("router_id", v, nil),
			"subnet_id":                 utils.PathSearch("internal_network_id", v, nil),
			"enterprise_project_id":     utils.PathSearch("enterprise_project_id", v, nil),
			"ngport_ip_address":         utils.PathSearch("ngport_ip_address", v, nil),
			"billing_info":              utils.PathSearch("billing_info", v, nil),
			"dnat_rules_limit":          utils.PathSearch("dnat_rules_limit", v, nil),
			"snat_rule_public_ip_limit": utils.PathSearch("snat_rule_public_ip_limit", v, nil),
			"pps_max":                   utils.PathSearch("pps_max", v, nil),
			"bps_max":                   utils.PathSearch("bps_max", v, nil),
			"session_conf":              flattenSessionConfig(utils.PathSearch("session_conf", v, nil)),
		})
	}
	return rst
}

func flattenSessionConfig(sessionConfig interface{}) []map[string]interface{} {
	if sessionConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"tcp_session_expire_time":  utils.PathSearch("tcp_session_expire_time", sessionConfig, nil),
			"udp_session_expire_time":  utils.PathSearch("udp_session_expire_time", sessionConfig, nil),
			"icmp_session_expire_time": utils.PathSearch("icmp_session_expire_time", sessionConfig, nil),
			"tcp_time_wait_time":       utils.PathSearch("tcp_time_wait_time", sessionConfig, nil),
		},
	}
}

func buildListPublicGatewaysQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	enterpriseProjectID := cfg.GetEnterpriseProjectID(d)

	if v, ok := d.GetOk("gateway_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("spec"); ok {
		res = fmt.Sprintf("%s&spec=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s&router_id=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&internal_network_id=%v", res, v)
	}
	if enterpriseProjectID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%s", res, enterpriseProjectID)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("created_at"); ok {
		res = fmt.Sprintf("%s&created_at=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
