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

// @API NAT GET /v2/{project_id}/dnat_rules
func DataSourceDnatRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnatRulesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the DNAT rules are located.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the DNAT rule.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the NAT gateway to which the DNAT rule belongs.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol type of the DNAT rule.",
			},
			"port_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port ID of the backend instance to which the DNAT rule belongs.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The private IP address of the backend instance to which the DNAT rule belongs.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the DNAT rule belongs.",
			},
			"internal_service_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port of the backend instance to which the DNAT rule belongs.",
			},
			"external_service_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port of the EIP associated with the DNAT rule belongs.",
			},
			"floating_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the EIP associated with the DNAT rule.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of EIP associated with the DNAT rule.",
			},
			"global_eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the global EIP associated with the DNAT rule.",
			},
			"global_eip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the global EIP associated with the DNAT rule.",
			},
			"rules": {
				Type:        schema.TypeList,
				Elem:        dnatRuleSchema(),
				Computed:    true,
				Description: "The list of the DNAT rules.",
			},
		},
	}
}

func dnatRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DNAT rule.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the NAT gateway to which the DNAT rule belongs.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol type of the private DNAT rule.",
			},
			"port_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The port ID of the backend instance to which the DNAT rule belongs.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private IP address of the backend instance to which the DNAT rule belongs.",
			},
			"internal_service_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port of the backend instance to which the DNAT rule belongs.",
			},
			"external_service_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port of the EIP associated with the DNAT rule belongs",
			},
			"floating_ip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the EIP associated with the DNAT rule.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the EIP associated with the DNAT rule.",
			},
			"global_eip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the global EIP associated with the DNAT rule.",
			},
			"global_eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the global EIP associated with the DNAT rule.",
			},
			"internal_service_port_range": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The port range of the backend instance to which the DNAT rule belongs.",
			},
			"external_service_port_range": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The port range of the EIP associated with the DNAT rule belongs",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the DNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the DNAT rule.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the DNAT rule.",
			},
		},
	}
	return &sc
}

func dataSourceDnatRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// listDnatRules: Query the DNAT rule list
	var (
		listDnatRulesHttpUrl = "v2/{project_id}/dnat_rules"
		listDnatRulesProduct = "nat"
	)
	listDnatRulesClient, err := cfg.NewServiceClient(listDnatRulesProduct, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	listDnatRulesPath := listDnatRulesClient.Endpoint + listDnatRulesHttpUrl
	listDnatRulesPath = strings.ReplaceAll(listDnatRulesPath, "{project_id}", listDnatRulesClient.ProjectID)

	listDnatRulesQueryParams := buildListDnatRuleQueryParams(d)
	listDnatRulesPath += listDnatRulesQueryParams

	listDnatRulesResp, err := pagination.ListAllItems(
		listDnatRulesClient,
		"marker",
		listDnatRulesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving DNAT rules %s", err)
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
		d.Set("rules", filterListDnatRuleResponseBody(flattenListDnatRulesResponseBody(listDnatRulesRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDnatRulesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("dnat_rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                          utils.PathSearch("id", v, nil),
			"gateway_id":                  utils.PathSearch("nat_gateway_id", v, nil),
			"protocol":                    utils.PathSearch("protocol", v, nil),
			"port_id":                     utils.PathSearch("port_id", v, nil),
			"private_ip":                  utils.PathSearch("private_ip", v, nil),
			"internal_service_port":       utils.PathSearch("internal_service_port", v, nil),
			"external_service_port":       utils.PathSearch("external_service_port", v, nil),
			"floating_ip_id":              utils.PathSearch("floating_ip_id", v, nil),
			"floating_ip_address":         utils.PathSearch("floating_ip_address", v, nil),
			"global_eip_id":               utils.PathSearch("global_eip_id", v, nil),
			"global_eip_address":          utils.PathSearch("global_eip_address", v, nil),
			"internal_service_port_range": utils.PathSearch("internal_service_port_range", v, nil),
			"external_service_port_range": utils.PathSearch("external_service_port_range", v, nil),
			"description":                 utils.PathSearch("description", v, nil),
			"status":                      utils.PathSearch("status", v, nil),
			"created_at":                  utils.PathSearch("created_at", v, nil),
		})
	}
	return rst
}

func filterListDnatRuleResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("global_eip_id"); ok &&
			fmt.Sprint(param) != utils.PathSearch("global_eip_id", v, nil) {
			continue
		}

		if param, ok := d.GetOk("global_eip_address"); ok &&
			fmt.Sprint(param) != utils.PathSearch("global_eip_address", v, nil) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListDnatRuleQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("rule_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		res = fmt.Sprintf("%s&nat_gateway_id=%v", res, v)
	}
	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&protocol=%v", res, v)
	}
	if v, ok := d.GetOk("port_id"); ok {
		res = fmt.Sprintf("%s&port_id=%v", res, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		res = fmt.Sprintf("%s&private_ip=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("floating_ip_id"); ok {
		res = fmt.Sprintf("%s&floating_ip_id=%v", res, v)
	}
	if v, ok := d.GetOk("floating_ip_address"); ok {
		res = fmt.Sprintf("%s&floating_ip_address=%v", res, v)
	}
	if v, ok := d.GetOk("internal_service_port"); ok {
		res = fmt.Sprintf("%s&internal_service_port=%v", res, v)
	}
	if v, ok := d.GetOk("external_service_port"); ok {
		res = fmt.Sprintf("%s&external_service_port=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}

	return res
}
