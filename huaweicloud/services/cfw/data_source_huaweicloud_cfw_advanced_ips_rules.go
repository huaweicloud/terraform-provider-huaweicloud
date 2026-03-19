package cfw

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

// @API CFW GET /v1/{project_id}/advanced-ips-rules
func DataSourceAdvancedIpsRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedIpsRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"advanced_ips_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ips_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ips_rule_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"param": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildQueryAdvancedIpsRulesQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	queryParam := fmt.Sprintf("?object_id=%s", d.Get("object_id").(string))

	if v := cfg.GetEnterpriseProjectID(d); v != "" {
		queryParam += fmt.Sprintf("&enterprise_project_id=%s", v)
	}

	return queryParam
}

func dataSourceAdvancedIpsRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/advanced-ips-rules"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildQueryAdvancedIpsRulesQueryParams(cfg, d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW advanced ips rules: %s", err)
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

	records := utils.PathSearch("data.advanced_ips_rules", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("advanced_ips_rules", flattenAdvancedIpsRulesResponse(records)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAdvancedIpsRulesResponse(respArray []interface{}) []map[string]interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"action":        utils.PathSearch("action", v, nil),
			"ips_rule_id":   utils.PathSearch("ips_rule_id", v, nil),
			"ips_rule_type": utils.PathSearch("ips_rule_type", v, nil),
			"param":         utils.PathSearch("param", v, nil),
			"status":        utils.PathSearch("status", v, nil),
		})
	}

	return rst
}
