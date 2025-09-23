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

// @API NAT GET /v3/{project_id}/private-nat/snat-rules
func DataSourcePrivateSnatRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateSnatRulesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the private SNAT rules are located.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the private SNAT rule.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the private NAT gateway to which the private SNAT rules belong.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CIDR block of the private SNAT rule.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the subnet to which the private SNAT rule belongs.",
			},
			"transit_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the transit IP associated with the private SNAT rule.",
			},
			"transit_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the transit IP associated with the private SNAT rule.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the private SNAT rules belong.",
			},
			"description": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The description of the private SNAT rule.",
			},
			"rules": {
				Type:        schema.TypeList,
				Elem:        snatRulesSchema(),
				Computed:    true,
				Description: "The list of the private SNAT rules.",
			},
		},
	}
}

func snatRulesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the private SNAT rule.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the private NAT gateway to which the private SNAT rule belongs.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CIDR block of the private SNAT rule.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the subnet to which the private SNAT rule belongs.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the private SNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the private SNAT rule.",
			},
			// Deprecated
			"transit_ip_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The ID of the transit IP associated with the private SNAT rule`, utils.SchemaDescInput{Deprecated: true},
				),
			},
			// Deprecated
			"transit_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The IP address of the transit IP associated with the private SNAT rule`, utils.SchemaDescInput{Deprecated: true},
				),
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the private SNAT rule.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the private SNAT rule.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the private SNAT rule belongs.",
			},
			"transit_ip_associations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The transit IP list associate with the private SNAT rule.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transit_ip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the transit IP associated with the private SNAT rule.`,
						},
						"transit_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address of the transit IP associated with the private SNAT rule.`,
						},
					},
				},
			},
		},
	}
	return &sc
}

func dataSourcePrivateSnatRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/snat-rules"
		product = "nat"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListSnatRulesQueryParams(d, cfg)
	resp, err := pagination.ListAllItems(
		client,
		"marker",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving private SNAT rules %s", err)
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
		d.Set("rules", flattenListSnatRuleResponseBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListSnatRuleResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("snat_rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                      utils.PathSearch("id", v, nil),
			"gateway_id":              utils.PathSearch("gateway_id", v, nil),
			"cidr":                    utils.PathSearch("cidr", v, nil),
			"subnet_id":               utils.PathSearch("virsubnet_id", v, nil),
			"description":             utils.PathSearch("description", v, nil),
			"status":                  utils.PathSearch("status", v, nil),
			"transit_ip_id":           utils.PathSearch("transit_ip_associations[0].transit_ip_id", v, nil),
			"transit_ip_address":      utils.PathSearch("transit_ip_associations[0].transit_ip_address", v, nil),
			"transit_ip_associations": flattenTransitIpList(utils.PathSearch("transit_ip_associations", v, make([]interface{}, 0))),
			"created_at":              utils.PathSearch("created_at", v, nil),
			"updated_at":              utils.PathSearch("updated_at", v, nil),
			"enterprise_project_id":   utils.PathSearch("enterprise_project_id", v, nil),
		})
	}
	return rst
}

func flattenTransitIpList(transitIpAssociations interface{}) []map[string]interface{} {
	rawArray := transitIpAssociations.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		params := map[string]interface{}{
			"transit_ip_id":      utils.PathSearch("transit_ip_id", v, nil),
			"transit_ip_address": utils.PathSearch("transit_ip_address", v, nil),
		}
		rst[i] = params
	}

	return rst
}

func buildListSnatRulesQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)
	descriptionList := d.Get("description").([]interface{})

	if v, ok := d.GetOk("rule_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		res = fmt.Sprintf("%s&gateway_id=%v", res, v)
	}
	if v, ok := d.GetOk("cidr"); ok {
		res = fmt.Sprintf("%s&cidr=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&virsubnet_id=%v", res, v)
	}
	if v, ok := d.GetOk("transit_ip_id"); ok {
		res = fmt.Sprintf("%s&transit_ip_id=%v", res, v)
	}
	if v, ok := d.GetOk("transit_ip_address"); ok {
		res = fmt.Sprintf("%s&transit_ip_address=%v", res, v)
	}
	if len(descriptionList) > 0 {
		for _, v := range descriptionList {
			res = fmt.Sprintf("%s&description=%v", res, v)
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
