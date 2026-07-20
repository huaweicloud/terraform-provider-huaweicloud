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

// @API DSC GET /v1/{project_id}/scan-jobs/{job_id}/classification-results
func DataSourceDscClassificationResults() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscClassificationResultsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the scan job ID.",
			},
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the keyword for fuzzy search on object names.",
			},
			"asset_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the asset type for filtering.",
			},
			"asset_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the asset ID for filtering.",
			},
			"security_level_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the security level IDs for filtering.",
			},
			"classification_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The classification results list.",
				Elem:        dscClassificationResultsEntitySchema(),
			},
		},
	}
}

func dscClassificationResultsEntitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The result ID.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project ID.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The job ID.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task ID.",
			},
			"ins_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID.",
			},
			"asset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset ID.",
			},
			"asset_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset name.",
			},
			"asset_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset type.",
			},
			"object_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The object name.",
			},
			"object_full_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The object full path.",
			},
			"security_level_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level name.",
			},
			"security_level_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level ID.",
			},
			"security_level_color": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The security level color.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time.",
			},
			"scan_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The scan time.",
			},
			"match_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The match information list.",
				Elem:        dscClassificationMatchInfoResultSchema(),
			},
		},
	}
}

func dscClassificationMatchInfoResultSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template ID.",
			},
			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template name.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule ID.",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule name.",
			},
			"security_level_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level name.",
			},
			"security_level_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level ID.",
			},
			"security_level_color": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The security level color.",
			},
			"classification_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The classification name.",
			},
			"classification_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The classification ID.",
			},
			"matched_detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The matched detail.",
			},
			"matched_examples": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The matched examples list.",
				Elem:        dscMatchedExamplesResultSchema(),
			},
		},
	}
}

func dscMatchedExamplesResultSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"line_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The line number of the match.",
			},
			"matched_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The matched content.",
			},
		},
	}
}

func buildDscClassificationResultsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("keyword"); ok {
		queryParams = fmt.Sprintf("%s&keyword=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_type"); ok {
		queryParams = fmt.Sprintf("%s&asset_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_id"); ok {
		queryParams = fmt.Sprintf("%s&asset_id=%v", queryParams, v)
	}
	securityLevelIds := d.Get("security_level_ids").([]interface{})
	for _, id := range securityLevelIds {
		queryParams = fmt.Sprintf("%s&security_level_ids=%v", queryParams, id)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceDscClassificationResultsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		jobId   = d.Get("job_id").(string)
		httpUrl = "v1/{project_id}/scan-jobs/{job_id}/classification-results"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	currentPath := requestPath + buildDscClassificationResultsQueryParams(d)

	resp, err := client.Request("GET", currentPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC classification results: %s", err)
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
		d.Set("classification_list", flattenDscClassificationResults(
			utils.PathSearch("classification_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscClassificationResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, v := range results {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"project_id":           utils.PathSearch("project_id", v, nil),
			"job_id":               utils.PathSearch("job_id", v, nil),
			"task_id":              utils.PathSearch("task_id", v, nil),
			"ins_id":               utils.PathSearch("ins_id", v, nil),
			"asset_id":             utils.PathSearch("asset_id", v, nil),
			"asset_name":           utils.PathSearch("asset_name", v, nil),
			"asset_type":           utils.PathSearch("asset_type", v, nil),
			"object_name":          utils.PathSearch("object_name", v, nil),
			"object_full_path":     utils.PathSearch("object_full_path", v, nil),
			"security_level_name":  utils.PathSearch("security_level_name", v, nil),
			"security_level_id":    utils.PathSearch("security_level_id", v, nil),
			"security_level_color": utils.PathSearch("security_level_color", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"update_time":          utils.PathSearch("update_time", v, nil),
			"scan_time":            utils.PathSearch("scan_time", v, nil),
			"match_info": flattenDscClassificationMatchInfoResults(
				utils.PathSearch("match_info", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscClassificationMatchInfoResults(matchInfos []interface{}) []interface{} {
	if len(matchInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(matchInfos))
	for _, v := range matchInfos {
		rst = append(rst, map[string]interface{}{
			"template_id":          utils.PathSearch("template_id", v, nil),
			"template_name":        utils.PathSearch("template_name", v, nil),
			"rule_id":              utils.PathSearch("rule_id", v, nil),
			"rule_name":            utils.PathSearch("rule_name", v, nil),
			"security_level_name":  utils.PathSearch("security_level_name", v, nil),
			"security_level_id":    utils.PathSearch("security_level_id", v, nil),
			"security_level_color": utils.PathSearch("security_level_color", v, nil),
			"classification_name":  utils.PathSearch("classification_name", v, nil),
			"classification_id":    utils.PathSearch("classification_id", v, nil),
			"matched_detail":       utils.PathSearch("matched_detail", v, nil),
			"matched_examples": flattenDscMatchedExamplesResults(
				utils.PathSearch("matched_examples", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscMatchedExamplesResults(examples []interface{}) []interface{} {
	if len(examples) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(examples))
	for _, v := range examples {
		rst = append(rst, map[string]interface{}{
			"line_number":     utils.PathSearch("line_number", v, nil),
			"matched_content": utils.PathSearch("matched_content", v, nil),
		})
	}

	return rst
}
