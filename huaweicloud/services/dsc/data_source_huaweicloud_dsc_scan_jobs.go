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

// @API DSC GET /v1/{project_id}/sdg/scan/job
func DataSourceScanJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScanJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_new": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     scanJobsSchema(),
			},
		},
	}
}

func scanJobsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"scan_templates": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cycle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_run_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_scan_risk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_nlp": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"open": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_level_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_level_color": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"asset_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     assetInfosSchema(),
			},
		},
	}
}

func assetInfosSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asset_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asset_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildScanJobsQueryParams(d *schema.ResourceData) string {
	params := ""
	if v, ok := d.GetOk("content"); ok {
		params += fmt.Sprintf("&content=%s", v.(string))
	}

	if v, ok := d.GetOk("is_new"); ok {
		params += fmt.Sprintf("&is_new=%s", v.(string))
	}

	if params != "" {
		params = "?" + params[1:]
	}

	return params
}

func dataSourceScanJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sdg/scan/job"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildScanJobsQueryParams(d)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	// The API's pagination parameters are not taking effect, so pagination functionality is not supported.
	resp, err := client.Request("GET", listPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC scan jobs: %s", err)
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
		d.Set("jobs", flattenScanJobs(utils.PathSearch(
			"tasks", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenScanJobs(tasks []interface{}) []interface{} {
	if len(tasks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tasks))
	for _, v := range tasks {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"rule_groups":          utils.PathSearch("rule_groups", v, nil),
			"scan_templates":       utils.PathSearch("scan_templates", v, nil),
			"cycle":                utils.PathSearch("cycle", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"last_run_time":        utils.PathSearch("last_run_time", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"last_scan_risk":       utils.PathSearch("last_scan_risk", v, nil),
			"use_nlp":              utils.PathSearch("use_nlp", v, nil),
			"open":                 utils.PathSearch("open", v, nil),
			"topic_urn":            utils.PathSearch("topic_urn", v, nil),
			"start_time":           utils.PathSearch("start_time", v, nil),
			"security_level_name":  utils.PathSearch("security_level_name", v, nil),
			"security_level_color": utils.PathSearch("security_level_color", v, nil),
			"asset_infos": flattenAssetInfos(utils.PathSearch(
				"asset_infos", v, make([]interface{}, 0)).([]interface{})),
		}))
	}

	return rst
}

func flattenAssetInfos(assetInfos []interface{}) []interface{} {
	if len(assetInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(assetInfos))
	for _, v := range assetInfos {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"asset_id":   utils.PathSearch("asset_id", v, nil),
			"asset_type": utils.PathSearch("asset_type", v, nil),
		}))
	}

	return rst
}
