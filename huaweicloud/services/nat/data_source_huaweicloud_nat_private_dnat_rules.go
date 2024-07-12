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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// listDnatRules: Query the private DNAT rule list
	var (
		listDnatRulesHttpUrl = "v3/{project_id}/private-nat/dnat-rules"
		listDnatRulesProduct = "nat"
	)
	listDnatRulesClient, err := cfg.NewServiceClient(listDnatRulesProduct, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	listDnatRulesPath := listDnatRulesClient.Endpoint + listDnatRulesHttpUrl
	listDnatRulesPath = strings.ReplaceAll(listDnatRulesPath, "{project_id}", listDnatRulesClient.ProjectID)

	listDnatRulesQueryParams := buildListDnatRulesQueryParams(d, cfg)
	listDnatRulesPath += listDnatRulesQueryParams

	listDnatRulesResp, err := pagination.ListAllItems(
		listDnatRulesClient,
		"marker",
		listDnatRulesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving private DNAT rules %s", err)
	}

	listDnatRulesRespJson, err := json.Marshal(listDnatRulesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listDnatRulesRespBody interface{}
	err = json.Unmarshal(listDnatRulesRespJson, &listDnatRulesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("rules", filterListDnatRulesResponseBody(flattenListDnatRuleResponseBody(listDnatRulesRespBody), d)),
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
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("protocol"); ok &&
			fmt.Sprint(param) != utils.PathSearch("protocol", v, nil) {
			continue
		}

		if param, ok := d.GetOk("internal_service_port"); ok &&
			fmt.Sprint(param) != utils.PathSearch("internal_service_port", v, nil) {
			continue
		}

		if param, ok := d.GetOk("transit_service_port"); ok &&
			fmt.Sprint(param) != utils.PathSearch("transit_service_port", v, nil) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListDnatRulesQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)

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
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}
	if res != "" {
		res = "?" + res[1:]
	}

	return res
}
