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

// @API NAT GET /v2/{project_id}/snat_rules
func DataSourceSnatRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSnatRulesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the SNAT rules are located.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the SNAT rule.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the NAT gateway to which the SNAT rule belongs.",
			},
			"floating_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the EIP associated with SNAT rule.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP of the EIP associated with SNAT rule.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CIDR block to which the SNAT rule belongs.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the subnet to which the SNAT rule belongs.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The source type of the SNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the SNAT rule.",
			},
			"global_eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IDs of the global EIP associated with SNAT rule.",
			},
			"global_eip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IPs of the global EIP associated with SNAT rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the SNAT rule.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The creation time of the SNAT rule.",
			},
			"rules": {
				Type:        schema.TypeList,
				Elem:        snatRuleSchema(),
				Computed:    true,
				Description: "The list of the SNAT rules.",
			},
		},
	}
}

func snatRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the SNAT rule.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the NAT gateway to which the SNAT rule belongs.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CIDR block to which the SNAT rule belongs.",
			},
			"floating_ip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the EIP associated with SNAT rule.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP of the EIP associated with SNAT rule.",
			},
			"source_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The source type of the SNAT rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the SNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the SNAT rule.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the SNAT rule.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the subnet to which the SNAT rule belongs.",
			},
			"global_eip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the global EIP associated with SNAT rule.",
			},
			"global_eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP of the global EIP associated with SNAT rule.",
			},
			"freezed_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP of the frozen global EIP associated with SNAT rule",
			},
		},
	}
	return &sc
}

func dataSourceSnatRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/snat_rules"
		product = "nat"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListSnatRuleQueryParams(d)

	resp, err := pagination.ListAllItems(
		client,
		"marker",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving SNAT rules %s", err)
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
		d.Set("rules", filterListSnatRulesResponseBody(flattenListSnatRulesResponseBody(respBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListSnatRulesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("snat_rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"gateway_id":          utils.PathSearch("nat_gateway_id", v, nil),
			"cidr":                utils.PathSearch("cidr", v, nil),
			"source_type":         utils.PathSearch("source_type", v, nil),
			"floating_ip_id":      utils.PathSearch("floating_ip_id", v, nil),
			"floating_ip_address": utils.PathSearch("floating_ip_address", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"status":              utils.PathSearch("status", v, nil),
			"created_at":          utils.PathSearch("created_at", v, nil),
			"subnet_id":           utils.PathSearch("network_id", v, nil),
			"global_eip_id":       utils.PathSearch("global_eip_id", v, nil),
			"global_eip_address":  utils.PathSearch("global_eip_address", v, nil),
			"freezed_ip_address":  utils.PathSearch("freezed_ip_address", v, nil),
		})
	}
	return rst
}

func filterListSnatRulesResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	var (
		globalEipID      = d.Get("global_eip_id").(string)
		globalEipAddress = d.Get("global_eip_address").(string)
	)

	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if globalEipID != "" && globalEipID != utils.PathSearch("global_eip_id", v, nil) {
			continue
		}

		if globalEipAddress != "" && globalEipAddress != utils.PathSearch("global_eip_address", v, nil) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListSnatRuleQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("rule_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		res = fmt.Sprintf("%s&nat_gateway_id=%v", res, v)
	}
	if v, ok := d.GetOk("cidr"); ok {
		res = fmt.Sprintf("%s&cidr=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&network_id=%v", res, v)
	}
	if v, ok := d.GetOk("floating_ip_id"); ok {
		res = fmt.Sprintf("%s&floating_ip_id=%v", res, v)
	}
	if v, ok := d.GetOk("floating_ip_address"); ok {
		res = fmt.Sprintf("%s&floating_ip_address=%v", res, v)
	}
	if v, ok := d.GetOk("source_type"); ok {
		res = fmt.Sprintf("%s&source_type=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
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
