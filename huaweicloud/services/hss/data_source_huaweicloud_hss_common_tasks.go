package hss

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

// @API HSS POST /v5/{project_id}/common/tasks/batch-query
func DataSourceCommonTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHssCommonTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Body parameters
			"task_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the task type.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the task ID to query.",
			},
			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the task name to fuzzy match.",
			},
			"start_create_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the start time of task creation time range query.",
			},
			"end_create_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the end time of task creation time range query.",
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the task trigger type.",
			},
			"task_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the task status.",
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sort key.",
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sort direction.",
			},
			"cluster_scan_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specifies the cluster scan information.",
				Elem:        buildHssClusterScanInfoSchema(),
			},
			"iac_scan_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specifies the IAC scan information.",
				Elem:        buildHssIacScanInfoSchema(),
			},
			// Query parameters
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
			// Attributes
			"data_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of tasks.",
				Elem:        buildCommonTasksDataListSchema(),
			},
		},
	}
}

func buildCommonTasksDataListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task ID.",
			},
			"task_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task type.",
			},
			"task_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task name.",
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task trigger type.",
			},
			"task_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task status.",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The start time of the task.",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The end time of the task.",
			},
			"estimated_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The estimated remaining time in minutes.",
			},
			"cluster_scan_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        buildHssClusterScanInfoAttributeSchema(),
				Description: "The cluster scan information.",
			},
			"iac_scan_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        buildHssIacScanInfoAttributeSchema(),
				Description: "The IAC scan information.",
			},
		},
	}
}

func buildHssIacScanInfoAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The file type.",
			},
			"scan_file_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of scanned files.",
			},
			"success_file_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of successfully scanned files.",
			},
			"failed_file_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of failed scanned files.",
			},
		},
	}
}

func buildHssClusterScanInfoAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scan_type_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of scan type.",
			},
			"scanning_cluster_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of clusters being scanned.",
			},
			"success_cluster_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of successfully scanned clusters.",
			},
			"failed_cluster_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of failed scanned clusters.",
			},
		},
	}
}

func buildHssIacScanInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the file type.",
			},
		},
	}
}

func buildHssClusterScanInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scan_type_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the scan type list.",
			},
		},
	}
}

func buildHssClusterScanInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("cluster_scan_info").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"scan_type_list": rawMap["scan_type_list"],
	}
}

func buildHssIacScanInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("iac_scan_info").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"file_type": rawMap["file_type"],
	}
}

func buildHssTasksBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"task_type":         d.Get("task_type"),
		"task_id":           utils.ValueIgnoreEmpty(d.Get("task_id")),
		"task_name":         utils.ValueIgnoreEmpty(d.Get("task_name")),
		"start_create_time": utils.ValueIgnoreEmpty(d.Get("start_create_time")),
		"end_create_time":   utils.ValueIgnoreEmpty(d.Get("end_create_time")),
		"trigger_type":      utils.ValueIgnoreEmpty(d.Get("trigger_type")),
		"task_status":       utils.ValueIgnoreEmpty(d.Get("task_status")),
		"sort_key":          utils.ValueIgnoreEmpty(d.Get("sort_key")),
		"sort_dir":          utils.ValueIgnoreEmpty(d.Get("sort_dir")),
		"cluster_scan_info": buildHssClusterScanInfoBodyParams(d),
		"iac_scan_info":     buildHssIacScanInfoBodyParams(d),
	}
}

func buildHssCommonTasksQueryParams(epsId string, offset int) string {
	rst := ""
	if epsId != "" {
		rst += fmt.Sprintf("&enterprise_project_id=%s", epsId)
	}
	if offset > 0 {
		rst += fmt.Sprintf("&offset=%d", offset)
	}

	if len(rst) > 0 {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceHssCommonTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/common/tasks/batch-query"
		result  = make([]interface{}, 0)
		offset  = 0
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildHssTasksBodyParams(d)),
	}

	for {
		requestWithOffset := requestPath + buildHssCommonTasksQueryParams(epsId, offset)
		resp, err := client.Request("POST", requestWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS tasks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data_list", flattenHssTasksDataList(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHssTasksDataList(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"task_id":           utils.PathSearch("task_id", v, nil),
			"task_type":         utils.PathSearch("task_type", v, nil),
			"task_name":         utils.PathSearch("task_name", v, nil),
			"trigger_type":      utils.PathSearch("trigger_type", v, nil),
			"task_status":       utils.PathSearch("task_status", v, nil),
			"start_time":        utils.PathSearch("start_time", v, nil),
			"end_time":          utils.PathSearch("end_time", v, nil),
			"estimated_time":    utils.PathSearch("estimated_time", v, nil),
			"cluster_scan_info": flattenClusterScanInfo(v),
			"iac_scan_info":     flattenIacScanInfo(v),
		})
	}

	return rst
}

func flattenClusterScanInfo(respBody interface{}) []interface{} {
	clusterScanInfo := utils.PathSearch("cluster_scan_info", respBody, nil)
	if clusterScanInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"scan_type_list":       utils.PathSearch("scan_type_list", clusterScanInfo, nil),
			"scanning_cluster_num": utils.PathSearch("scanning_cluster_num", clusterScanInfo, nil),
			"success_cluster_num":  utils.PathSearch("success_cluster_num", clusterScanInfo, nil),
			"failed_cluster_num":   utils.PathSearch("failed_cluster_num", clusterScanInfo, nil),
		},
	}
}

func flattenIacScanInfo(respBody interface{}) []interface{} {
	iacScanInfo := utils.PathSearch("iac_scan_info", respBody, nil)
	if iacScanInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"file_type":        utils.PathSearch("file_type", iacScanInfo, nil),
			"scan_file_num":    utils.PathSearch("scan_file_num", iacScanInfo, nil),
			"success_file_num": utils.PathSearch("success_file_num", iacScanInfo, nil),
			"failed_file_num":  utils.PathSearch("failed_file_num", iacScanInfo, nil),
		},
	}
}
