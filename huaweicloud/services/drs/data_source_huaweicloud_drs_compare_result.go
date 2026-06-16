package drs

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

// @API DRS POST /v3/{project_id}/jobs/query-compare-result
func DataSourceDrsCompareResult() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsCompareResultRead,

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
			"current_page": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"per_page": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"object_level_compare_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"line_compare_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_compare_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_level_compare_results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     objectLevelCompareResultsSchema(),
			},
			"line_compare_results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     lineCompareResultsSchema(),
			},
			"content_compare_results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     contentCompareResultsSchema(),
			},
			"compare_task_list_results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     compareTaskListResultsSchema(),
			},
		},
	}
}

func objectLevelCompareResultsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"compare_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_compare_overview": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     objectCompareOverviewSchema(),
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func objectCompareOverviewSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_compare_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"diff_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func lineCompareResultsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"compare_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"line_compare_overview": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     lineCompareOverviewSchema(),
			},
			"line_compare_overview_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"line_compare_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     lineCompareDetailsSchema(),
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func lineCompareOverviewSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"line_compare_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func lineCompareDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"line_compare_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     lineCompareDetailSchema(),
			},
			"line_compare_detail_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func lineCompareDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"diff_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"line_compare_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func contentCompareResultsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"compare_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_compare_overview": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     contentCompareOverviewSchema(),
			},
			"content_compare_overview_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content_compare_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     contentCompareDetailsSchema(),
			},
			"content_compare_diffs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     contentCompareDiffsSchema(),
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func contentCompareOverviewSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_compare_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func contentCompareDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_compare_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     contentCompareDetailSchema(),
			},
			"content_compare_detail_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content_uncompare_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     contentCompareDetailSchema(),
			},
			"content_uncompare_detail_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func contentCompareDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"diff_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"line_compare_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_compare_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func contentCompareDiffsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_compare_diff": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     contentCompareDiffSchema(),
			},
			"content_compare_diff_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func contentCompareDiffSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_select_sql": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_select_sql": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_key_value": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_key_value": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func compareTaskListResultsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"compare_task_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     compareTaskListSchema(),
			},
			"compare_task_list_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func compareTaskListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"compare_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_task_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCompareResultBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"job_id":       d.Get("job_id").(string),
		"current_page": d.Get("current_page").(int),
		"per_page":     d.Get("per_page").(int),
	}

	if v, ok := d.GetOk("object_level_compare_id"); ok {
		body["object_level_compare_id"] = v.(string)
	}
	if v, ok := d.GetOk("line_compare_id"); ok {
		body["line_compare_id"] = v.(string)
	}
	if v, ok := d.GetOk("content_compare_id"); ok {
		body["content_compare_id"] = v.(string)
	}

	return body
}

func dataSourceDrsCompareResultRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/query-compare-result"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildCompareResultBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS compare result: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("object_level_compare_results", flattenObjectLevelCompareResults(respBody)),
		d.Set("line_compare_results", flattenLineCompareResults(respBody)),
		d.Set("content_compare_results", flattenContentCompareResults(respBody)),
		d.Set("compare_task_list_results", flattenCompareTaskListResults(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenObjectLevelCompareResults(respBody interface{}) []interface{} {
	resultMap := utils.PathSearch("object_level_compare_results", respBody, nil)
	if resultMap == nil {
		return nil
	}

	objectCompareOverview := utils.PathSearch("object_compare_overview", resultMap, make([]interface{}, 0))
	return []interface{}{
		map[string]interface{}{
			"compare_task_id":         utils.PathSearch("compare_task_id", resultMap, nil),
			"object_compare_overview": flattenObjectCompareOverview(objectCompareOverview),
			"error_code":              utils.PathSearch("error_code", resultMap, nil),
			"error_msg":               utils.PathSearch("error_msg", resultMap, nil),
		},
	}
}

func flattenObjectCompareOverview(overview interface{}) []interface{} {
	overviewRaw, ok := overview.([]interface{})
	if !ok || len(overviewRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(overviewRaw))
	for _, item := range overviewRaw {
		result = append(result, map[string]interface{}{
			"object_type":           utils.PathSearch("object_type", item, nil),
			"object_compare_result": utils.PathSearch("object_compare_result", item, nil),
			"target_count":          utils.PathSearch("target_count", item, nil),
			"source_count":          utils.PathSearch("source_count", item, nil),
			"diff_count":            utils.PathSearch("diff_count", item, nil),
		})
	}
	return result
}

func flattenLineCompareResults(respBody interface{}) []interface{} {
	resultMap := utils.PathSearch("line_compare_results", respBody, nil)
	if resultMap == nil {
		return nil
	}

	lineCompareOverview := utils.PathSearch("line_compare_overview", resultMap, make([]interface{}, 0))
	lineCompareDetails := utils.PathSearch("line_compare_details", resultMap, make([]interface{}, 0))

	return []interface{}{
		map[string]interface{}{
			"compare_task_id":             utils.PathSearch("compare_task_id", resultMap, nil),
			"line_compare_overview":       flattenLineCompareOverview(lineCompareOverview),
			"line_compare_overview_count": utils.PathSearch("line_compare_overview_count", resultMap, nil),
			"line_compare_details":        flattenLineCompareDetails(lineCompareDetails),
			"error_code":                  utils.PathSearch("error_code", resultMap, nil),
			"error_msg":                   utils.PathSearch("error_msg", resultMap, nil),
		},
	}
}

func flattenLineCompareOverview(overview interface{}) []interface{} {
	overviewRaw, ok := overview.([]interface{})
	if !ok || len(overviewRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(overviewRaw))
	for _, item := range overviewRaw {
		result = append(result, map[string]interface{}{
			"source_db_name":      utils.PathSearch("source_db_name", item, nil),
			"target_db_name":      utils.PathSearch("target_db_name", item, nil),
			"line_compare_result": utils.PathSearch("line_compare_result", item, nil),
		})
	}
	return result
}

func flattenLineCompareDetails(details interface{}) []interface{} {
	detailsRaw, ok := details.([]interface{})
	if !ok || len(detailsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(detailsRaw))
	for _, item := range detailsRaw {
		lineCompareDetail := utils.PathSearch("line_compare_detail", item, make([]interface{}, 0))
		result = append(result, map[string]interface{}{
			"source_db_name":            utils.PathSearch("source_db_name", item, nil),
			"line_compare_detail":       flattenLineCompareDetail(lineCompareDetail),
			"line_compare_detail_count": utils.PathSearch("line_compare_detail_count", item, nil),
		})
	}
	return result
}

func flattenLineCompareDetail(detail interface{}) []interface{} {
	detailRaw, ok := detail.([]interface{})
	if !ok || len(detailRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(detailRaw))
	for _, item := range detailRaw {
		result = append(result, map[string]interface{}{
			"source_table_name":   utils.PathSearch("source_table_name", item, nil),
			"target_table_name":   utils.PathSearch("target_table_name", item, nil),
			"source_row_num":      utils.PathSearch("source_row_num", item, nil),
			"target_row_num":      utils.PathSearch("target_row_num", item, nil),
			"diff_row_num":        utils.PathSearch("diff_row_num", item, nil),
			"line_compare_result": utils.PathSearch("line_compare_result", item, nil),
			"message":             utils.PathSearch("message", item, nil),
		})
	}
	return result
}

func flattenContentCompareResults(respBody interface{}) []interface{} {
	resultMap := utils.PathSearch("content_compare_results", respBody, nil)
	if resultMap == nil {
		return nil
	}

	contentCompareOverview := utils.PathSearch("content_compare_overview", resultMap, make([]interface{}, 0))
	contentCompareDetails := utils.PathSearch("content_compare_details", resultMap, make([]interface{}, 0))
	contentCompareDiffs := utils.PathSearch("content_compare_diffs", resultMap, make([]interface{}, 0))

	return []interface{}{
		map[string]interface{}{
			"compare_task_id":                utils.PathSearch("compare_task_id", resultMap, nil),
			"content_compare_overview":       flattenContentCompareOverview(contentCompareOverview),
			"content_compare_overview_count": utils.PathSearch("content_compare_overview_count", resultMap, nil),
			"content_compare_details":        flattenContentCompareDetails(contentCompareDetails),
			"content_compare_diffs":          flattenContentCompareDiffs(contentCompareDiffs),
			"error_code":                     utils.PathSearch("error_code", resultMap, nil),
			"error_msg":                      utils.PathSearch("error_msg", resultMap, nil),
		},
	}
}

func flattenContentCompareOverview(overview interface{}) []interface{} {
	overviewRaw, ok := overview.([]interface{})
	if !ok || len(overviewRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(overviewRaw))
	for _, item := range overviewRaw {
		result = append(result, map[string]interface{}{
			"source_db_name":         utils.PathSearch("source_db_name", item, nil),
			"target_db_name":         utils.PathSearch("target_db_name", item, nil),
			"content_compare_result": utils.PathSearch("content_compare_result", item, nil),
		})
	}
	return result
}

func flattenContentCompareDetails(details interface{}) []interface{} {
	detailsRaw, ok := details.([]interface{})
	if !ok || len(detailsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(detailsRaw))
	for _, item := range detailsRaw {
		contentCompareDetail := utils.PathSearch("content_compare_detail", item, make([]interface{}, 0))
		contentUncompareDetail := utils.PathSearch("content_uncompare_detail", item, make([]interface{}, 0))

		result = append(result, map[string]interface{}{
			"source_db_name":                 utils.PathSearch("source_db_name", item, nil),
			"content_compare_detail":         flattenContentCompareDetail(contentCompareDetail),
			"content_compare_detail_count":   utils.PathSearch("content_compare_detail_count", item, nil),
			"content_uncompare_detail":       flattenContentCompareDetail(contentUncompareDetail),
			"content_uncompare_detail_count": utils.PathSearch("content_uncompare_detail_count", item, nil),
		})
	}
	return result
}

func flattenContentCompareDetail(detail interface{}) []interface{} {
	detailRaw, ok := detail.([]interface{})
	if !ok || len(detailRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(detailRaw))
	for _, item := range detailRaw {
		result = append(result, map[string]interface{}{
			"source_db_name":         utils.PathSearch("source_db_name", item, nil),
			"target_db_name":         utils.PathSearch("target_db_name", item, nil),
			"source_table_name":      utils.PathSearch("source_table_name", item, nil),
			"target_table_name":      utils.PathSearch("target_table_name", item, nil),
			"source_row_num":         utils.PathSearch("source_row_num", item, nil),
			"target_row_num":         utils.PathSearch("target_row_num", item, nil),
			"diff_row_num":           utils.PathSearch("diff_row_num", item, nil),
			"line_compare_result":    utils.PathSearch("line_compare_result", item, nil),
			"content_compare_result": utils.PathSearch("content_compare_result", item, nil),
			"message":                utils.PathSearch("message", item, nil),
		})
	}
	return result
}

func flattenContentCompareDiffs(diffs interface{}) []interface{} {
	diffsRaw, ok := diffs.([]interface{})
	if !ok || len(diffsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(diffsRaw))
	for _, item := range diffsRaw {
		contentCompareDiff := utils.PathSearch("content_compare_diff", item, make([]interface{}, 0))

		result = append(result, map[string]interface{}{
			"source_db_name":             utils.PathSearch("source_db_name", item, nil),
			"source_table_name":          utils.PathSearch("source_table_name", item, nil),
			"content_compare_diff":       flattenContentCompareDiff(contentCompareDiff),
			"content_compare_diff_count": utils.PathSearch("content_compare_diff_count", item, nil),
		})
	}
	return result
}

func flattenContentCompareDiff(diff interface{}) []interface{} {
	diffRaw, ok := diff.([]interface{})
	if !ok || len(diffRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(diffRaw))
	for _, item := range diffRaw {
		result = append(result, map[string]interface{}{
			"target_select_sql": utils.PathSearch("target_select_sql", item, nil),
			"source_select_sql": utils.PathSearch("source_select_sql", item, nil),
			"source_key_value":  utils.PathSearch("source_key_value", item, nil),
			"target_key_value":  utils.PathSearch("target_key_value", item, nil),
		})
	}
	return result
}

func flattenCompareTaskListResults(respBody interface{}) []interface{} {
	resultMap := utils.PathSearch("compare_task_list_results", respBody, nil)
	if resultMap == nil {
		return nil
	}

	compareTaskList := utils.PathSearch("compare_task_list", resultMap, make([]interface{}, 0))

	return []interface{}{
		map[string]interface{}{
			"compare_task_list":       flattenCompareTaskList(compareTaskList),
			"compare_task_list_count": utils.PathSearch("compare_task_list_count", resultMap, nil),
			"error_msg":               utils.PathSearch("error_msg", resultMap, nil),
			"error_code":              utils.PathSearch("error_code", resultMap, nil),
		},
	}
}

func flattenCompareTaskList(taskList interface{}) []interface{} {
	taskListRaw, ok := taskList.([]interface{})
	if !ok || len(taskListRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(taskListRaw))
	for _, item := range taskListRaw {
		result = append(result, map[string]interface{}{
			"compare_task_id":     utils.PathSearch("compare_task_id", item, nil),
			"compare_type":        utils.PathSearch("compare_type", item, nil),
			"compare_task_status": utils.PathSearch("compare_task_status", item, nil),
			"create_time":         utils.PathSearch("create_time", item, nil),
			"end_time":            utils.PathSearch("end_time", item, nil),
		})
	}
	return result
}
