package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/scan-jobs/{job_id}/hit-rules
func DataSourceDscHitRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscHitRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_level_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hit_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscHitRuleSchema(),
			},
		},
	}
}

func dscHitRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"top_objects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildDscHitRulesQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	if v, ok := d.GetOk("keyword"); ok {
		queryParams = fmt.Sprintf("%s&keyword=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_type"); ok {
		queryParams = fmt.Sprintf("%s&asset_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_id"); ok {
		queryParams = fmt.Sprintf("%s&asset_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("security_level_ids"); ok {
		ids := v.([]interface{})
		for _, id := range ids {
			queryParams = fmt.Sprintf("%s&security_level_ids=%v", queryParams, id)
		}
	}
	return queryParams
}

func dataSourceDscHitRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/scan-jobs/{job_id}/hit-rules"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	for {
		currentPath := requestPath + buildDscHitRulesQueryParams(d, limit, offset)
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC hit rules: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		hitRulesList := utils.PathSearch("hit_rule_list", requestRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, hitRulesList...)
		if len(hitRulesList) < limit {
			break
		}
		offset += len(hitRulesList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("hit_rules", flattenDscHitRules(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscHitRules(hitRulesList []interface{}) []interface{} {
	if len(hitRulesList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(hitRulesList))
	for _, v := range hitRulesList {
		rst = append(rst, map[string]interface{}{
			"rule_id":     utils.PathSearch("rule_id", v, nil),
			"rule_name":   utils.PathSearch("rule_name", v, nil),
			"top_objects": utils.ExpandToStringList(utils.PathSearch("top_objects", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}
