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

// @API DSC GET /v1/{project_id}/metadata/catalog/statical-chart
func DataSourceCatalogStaticalChart() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCatalogStaticalChartRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"label_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detection_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     detectionRuleSchema(),
			},
			"sensitive_col_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sensitiveColInfoSchema(),
			},
			"total_column_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func detectionRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hit_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func sensitiveColInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"color_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"level_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sensitive_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCatalogStaticalChartQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("label_id"); ok {
		res += fmt.Sprintf("&label_id=%s", v.(string))
	}
	if v, ok := d.GetOk("type_id"); ok {
		res += fmt.Sprintf("&type_id=%s", v.(string))
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func dataSourceCatalogStaticalChartRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/metadata/catalog/statical-chart"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildCatalogStaticalChartQueryParams(d)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC catalog statical chart: %s", err)
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
		d.Set("detection_rules", flattenDetectionRules(utils.PathSearch(
			"detection_rules", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("sensitive_col_infos", flattenSensitiveColInfos(utils.PathSearch(
			"sensitive_col_infos", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("total_column_number", utils.PathSearch("total_column_number", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDetectionRules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"hit_number": utils.PathSearch("hit_number", v, nil),
			"rule_name":  utils.PathSearch("rule_name", v, nil),
		}))
	}

	return rst
}

func flattenSensitiveColInfos(infos []interface{}) []interface{} {
	if len(infos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(infos))
	for _, v := range infos {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"color_number":     utils.PathSearch("color_number", v, nil),
			"level_name":       utils.PathSearch("level_name", v, nil),
			"sensitive_number": utils.PathSearch("sensitive_number", v, nil),
		}))
	}

	return rst
}
