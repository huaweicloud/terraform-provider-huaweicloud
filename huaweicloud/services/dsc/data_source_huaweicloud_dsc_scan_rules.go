package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/scan-rules
func DataSourceScanRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScanRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scan_rules_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     scanRulesSchema(),
			},
		},
	}
}

func scanRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_desc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceScanRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scan-rules"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC scan rules: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("scan_rules_list", flattenScanRulesList(
			utils.PathSearch("scan_rules_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenScanRulesList(scanRulesList []interface{}) []interface{} {
	if len(scanRulesList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(scanRulesList))
	for _, v := range scanRulesList {
		rst = append(rst, map[string]interface{}{
			"category":   utils.PathSearch("category", v, nil),
			"project_id": utils.PathSearch("project_id", v, nil),
			"rule_desc":  utils.PathSearch("rule_desc", v, nil),
			"rule_id":    utils.PathSearch("rule_id", v, nil),
			"rule_name":  utils.PathSearch("rule_name", v, nil),
			"rule_type":  utils.PathSearch("rule_type", v, nil),
		})
	}

	return rst
}
