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

// @API HSS GET /v5/{project_id}/antivirus/task
func DataSourceAntivirusVirusScanTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAntivirusVirusScanTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"whether_paid_task": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_task_status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scan_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"task_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"success_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"fail_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cancel_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_info_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"asset_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"run_duration": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"scan_progress": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"virus_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"scan_file_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"host_task_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fail_reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"deleted": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"whether_using_quota": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"agent_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"agent_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protect_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"rescan": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"whether_paid_task": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAntivirusVirusScanTasksQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	queryParams = fmt.Sprintf("%s&whether_paid_task=%v", queryParams, d.Get("whether_paid_task").(bool))

	if v, ok := d.GetOk("task_name"); ok {
		queryParams = fmt.Sprintf("%s&task_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("last_days"); ok {
		queryParams = fmt.Sprintf("%s&last_days=%v", queryParams, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		queryParams = fmt.Sprintf("%s&begin_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("task_status"); ok {
		queryParams = fmt.Sprintf("%s&task_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		queryParams = fmt.Sprintf("%s&public_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_task_status"); ok {
		for _, raw := range v.([]interface{}) {
			queryParams = fmt.Sprintf("%s&host_task_status=%v", queryParams, raw)
		}
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceAntivirusVirusScanTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		offset   = 0
		totalNum float64
		result   = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/antivirus/task"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAntivirusVirusScanTasksQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS antivirus virus scan tasks: %s", err)
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

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		if int(totalNum) == len(result) {
			break
		}

		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenAntivirusVirusScanTasksDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAntivirusVirusScanTasksDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"task_id":          utils.PathSearch("task_id", v, nil),
			"task_name":        utils.PathSearch("task_name", v, nil),
			"scan_type":        utils.PathSearch("scan_type", v, nil),
			"start_type":       utils.PathSearch("start_type", v, nil),
			"action":           utils.PathSearch("action", v, nil),
			"start_time":       utils.PathSearch("start_time", v, nil),
			"task_status":      utils.PathSearch("task_status", v, nil),
			"host_num":         utils.PathSearch("host_num", v, nil),
			"success_host_num": utils.PathSearch("success_host_num", v, nil),
			"fail_host_num":    utils.PathSearch("fail_host_num", v, nil),
			"cancel_host_num":  utils.PathSearch("cancel_host_num", v, nil),
			"host_info_list": flattenAntivirusVirusScanTasksHostInfoList(
				utils.PathSearch("host_info_list", v, make([]interface{}, 0)).([]interface{})),
			"rescan":            utils.PathSearch("rescan", v, nil),
			"whether_paid_task": utils.PathSearch("whether_paid_task", v, nil),
		})
	}

	return rst
}

func flattenAntivirusVirusScanTasksHostInfoList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_id":             utils.PathSearch("host_id", v, nil),
			"host_name":           utils.PathSearch("host_name", v, nil),
			"private_ip":          utils.PathSearch("private_ip", v, nil),
			"public_ip":           utils.PathSearch("public_ip", v, nil),
			"asset_value":         utils.PathSearch("asset_value", v, nil),
			"start_time":          utils.PathSearch("start_time", v, nil),
			"run_duration":        utils.PathSearch("run_duration", v, nil),
			"scan_progress":       utils.PathSearch("scan_progress", v, nil),
			"virus_num":           utils.PathSearch("virus_num", v, nil),
			"scan_file_num":       utils.PathSearch("scan_file_num", v, nil),
			"host_task_status":    utils.PathSearch("host_task_status", v, nil),
			"fail_reason":         utils.PathSearch("fail_reason", v, nil),
			"deleted":             utils.PathSearch("deleted", v, nil),
			"whether_using_quota": utils.PathSearch("whether_using_quota", v, nil),
			"agent_id":            utils.PathSearch("agent_id", v, nil),
			"os_type":             utils.PathSearch("os_type", v, nil),
			"host_status":         utils.PathSearch("host_status", v, nil),
			"agent_status":        utils.PathSearch("agent_status", v, nil),
			"protect_status":      utils.PathSearch("protect_status", v, nil),
			"os_name":             utils.PathSearch("os_name", v, nil),
			"os_version":          utils.PathSearch("os_version", v, nil),
		})
	}

	return rst
}
