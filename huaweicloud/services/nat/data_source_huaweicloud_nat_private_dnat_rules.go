package nat

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

// @API NAT GET /v3/{project_id}/private-nat/dnat-rules
func DataSourcePrivateDnatRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateDnatRulesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the private DNAT rules are located.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the private DNAT rule.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the private NAT gateway to which the private DNAT rules belong.",
			},
			"backend_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the backend instance to which the private DNAT rules belong.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol type of the private DNAT rules.",
			},
			"internal_service_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port of the backend instance to which the private DNAT rule belongs.",
			},
			"backend_interface_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The network interface ID of the backend instance to which the private DNAT rule belongs.",
			},
			"transit_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the transit IP associated with the private DNAT rules.",
			},
			"transit_service_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port of the transit IP associated with the private DNAT rule.",
			},
			"backend_private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The private IP address of the backend instance to which the private DNAT rule belongs.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the private DNAT rules belong.",
			},
			"description": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The description of the private DNAT rule.",
			},
			"external_ip_address": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The transit IP address used to the private DNAT rule.",
			},
			"rules": {
				Type:        schema.TypeList,
				Elem:        dnatRulesSchema(),
				Computed:    true,
				Description: "The list of the private DNAT rules.",
			},
		},
	}
}

func dnatRulesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the private DNAT rule.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the private NAT gateway to which the private DNAT rule belongs.",
			},
			"backend_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the backend instance to which the private DNAT rules belong.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol type of the private DNAT rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the private DNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the private DNAT rule.",
			},
			"internal_service_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The port of the backend instance to which the private DNAT rule belongs.",
			},
			"backend_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network interface ID of the backend instance to which the private DNAT rule belongs.",
			},
			"transit_ip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the transit IP associated with the private DNAT rule.",
			},
			"transit_service_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the transit IP associated with the private DNAT rule.",
			},
			"backend_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private IP address of the backend instance to which the private DNAT rule belongs.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the private DNAT rule.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the private DNAT rule.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the private DNAT rule belongs.",
			},
		},
	}
	return &sc
}

func dataSourcePrivateDnatRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/dnat-rules"
		product = "nat"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListDnatRulesQueryParams(d, cfg)
	resp, err := pagination.ListAllItems(
		client,
		"marker",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving private DNAT rules %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("rules", filterListDnatRulesResponseBody(flattenListDnatRuleResponseBody(respBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDnatRuleResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("dnat_rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"gateway_id":            utils.PathSearch("gateway_id", v, nil),
			"backend_type":          utils.PathSearch("type", v, nil),
			"protocol":              utils.PathSearch("protocol", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"internal_service_port": utils.PathSearch("internal_service_port", v, nil),
			"backend_interface_id":  utils.PathSearch("network_interface_id", v, nil),
			"transit_ip_id":         utils.PathSearch("transit_ip_id", v, nil),
			"transit_service_port":  utils.PathSearch("transit_service_port", v, nil),
			"backend_private_ip":    utils.PathSearch("private_ip_address", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"updated_at":            utils.PathSearch("updated_at", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
		})
	}
	return rst
}

func filterListDnatRulesResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	var (
		protocol            = d.Get("protocol").(string)
		internalServicePort = d.Get("internal_service_port").(string)
		transitServicePort  = d.Get("transit_service_port").(string)
		rst                 = make([]interface{}, 0, len(all))
	)

	for _, v := range all {
		if protocol != "" && protocol != utils.PathSearch("protocol", v, nil) {
			continue
		}

		if internalServicePort != "" && internalServicePort != utils.PathSearch("internal_service_port", v, nil) {
			continue
		}

		if transitServicePort != "" && transitServicePort != utils.PathSearch("transit_service_port", v, nil) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListDnatRulesQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)
	descriptionList := d.Get("description").([]interface{})
	externalIpAddresses := d.Get("external_ip_address").([]interface{})

	if v, ok := d.GetOk("rule_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		res = fmt.Sprintf("%s&gateway_id=%v", res, v)
	}
	if v, ok := d.GetOk("backend_type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("backend_interface_id"); ok {
		res = fmt.Sprintf("%s&network_interface_id=%v", res, v)
	}
	if v, ok := d.GetOk("transit_ip_id"); ok {
		res = fmt.Sprintf("%s&transit_ip_id=%v", res, v)
	}
	if v, ok := d.GetOk("backend_private_ip"); ok {
		res = fmt.Sprintf("%s&private_ip_address=%v", res, v)
	}
	if len(descriptionList) > 0 {
		for _, v := range descriptionList {
			res = fmt.Sprintf("%s&description=%v", res, v)
		}
	}
	if len(externalIpAddresses) > 0 {
		for _, v := range externalIpAddresses {
			res = fmt.Sprintf("%s&external_ip_address=%v", res, v)
		}
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}
	if res != "" {
		res = "?" + res[1:]
	}

	return res
}
