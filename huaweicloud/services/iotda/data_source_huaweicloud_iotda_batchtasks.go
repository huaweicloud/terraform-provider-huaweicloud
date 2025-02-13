package iotda

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

// @API IoTDA GET /v5/iot/{project_id}/batchtasks
func DataSourceBatchTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBatchTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"batchtasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"targets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"targets_filter": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"task_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"retry_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"retry_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_progress": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"processing": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"success": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fail": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"waitting": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fail_wait_retry": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"stopped": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"removed": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBatchTasksQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?limit=50&task_type=%v", d.Get("type").(string))
	if v, ok := d.GetOk("space_id"); ok {
		queryParams = fmt.Sprintf("%s&app_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceBatchTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		httpUrl   = "v5/iot/{project_id}/batchtasks"
		offset    = 0
		result    = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildBatchTasksQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving IoTDA batch tasks: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		batchTasksResp := utils.PathSearch("batchtasks", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(batchTasksResp) == 0 {
			break
		}

		result = append(result, batchTasksResp...)
		offset += len(batchTasksResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("batchtasks", flattenBatchTasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBatchTasks(batchTasksResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(batchTasksResp))
	for _, v := range batchTasksResp {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("task_id", v, nil),
			"name":           utils.PathSearch("task_name", v, nil),
			"type":           utils.PathSearch("task_type", v, nil),
			"targets":        utils.PathSearch("targets", v, nil),
			"targets_filter": utils.PathSearch("targets_filter", v, nil),
			"task_policy":    flattenDataSourceTaskPolicy(utils.PathSearch("task_policy", v, nil)),
			"status":         utils.PathSearch("status", v, nil),
			"status_desc":    utils.PathSearch("status_desc", v, nil),
			"task_progress":  flattenDataSourceTaskProgress(utils.PathSearch("task_progress", v, nil)),
			"created_at":     utils.PathSearch("create_time", v, nil),
		})
	}

	return rst
}

func flattenDataSourceTaskPolicy(taskPolicyResp interface{}) []interface{} {
	if taskPolicyResp == nil {
		return nil
	}

	rst := map[string]interface{}{
		"schedule_time":  utils.PathSearch("schedule_time", taskPolicyResp, nil),
		"retry_count":    utils.PathSearch("retry_count", taskPolicyResp, nil),
		"retry_interval": utils.PathSearch("retry_interval", taskPolicyResp, nil),
	}

	return []interface{}{rst}
}

func flattenDataSourceTaskProgress(taskProgressResp interface{}) []interface{} {
	if taskProgressResp == nil {
		return nil
	}

	rst := map[string]interface{}{
		"total":           utils.PathSearch("total", taskProgressResp, nil),
		"processing":      utils.PathSearch("processing", taskProgressResp, nil),
		"success":         utils.PathSearch("success", taskProgressResp, nil),
		"fail":            utils.PathSearch("fail", taskProgressResp, nil),
		"waitting":        utils.PathSearch("waitting", taskProgressResp, nil),
		"fail_wait_retry": utils.PathSearch("fail_wait_retry", taskProgressResp, nil),
		"stopped":         utils.PathSearch("stopped", taskProgressResp, nil),
		"removed":         utils.PathSearch("removed", taskProgressResp, nil),
	}

	return []interface{}{rst}
}
