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

// @API HSS GET /v5/{project_id}/image/tasks
func DataSourceImageTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageTasksRead,

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
			"global_image_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"task_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scan_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begin_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remain_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"task_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scan_scope": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rate_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_all": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"failed_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"success_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"risky_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sync_task_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failed_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failed_images": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"registry_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"registry_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"registry_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"failed_reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildImageTasksQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?type=%v&limit=200", d.Get("type"))

	if v, ok := d.GetOk("global_image_type"); ok {
		queryParams = fmt.Sprintf("%s&global_image_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("task_type"); ok {
		queryParams = fmt.Sprintf("%s&task_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("task_name"); ok {
		queryParams = fmt.Sprintf("%s&task_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("task_id"); ok {
		queryParams = fmt.Sprintf("%s&task_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("create_time"); ok {
		queryParams = fmt.Sprintf("%s&create_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("task_status"); ok {
		queryParams = fmt.Sprintf("%s&task_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("scan_scope"); ok {
		queryParams = fmt.Sprintf("%s&scan_scope=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceImageTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/image/tasks"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildImageTasksQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving image tasks: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenImageTasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageTasks(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		result = append(result, map[string]interface{}{
			"task_id":        utils.PathSearch("task_id", v, nil),
			"policy_id":      utils.PathSearch("policy_id", v, nil),
			"task_name":      utils.PathSearch("task_name", v, nil),
			"begin_time":     utils.PathSearch("begin_time", v, nil),
			"end_time":       utils.PathSearch("end_time", v, nil),
			"remain_min":     utils.PathSearch("remain_min", v, nil),
			"task_type":      utils.PathSearch("task_type", v, nil),
			"image_type":     utils.PathSearch("image_type", v, nil),
			"task_status":    utils.PathSearch("task_status", v, nil),
			"scan_scope":     utils.PathSearch("scan_scope", v, nil),
			"rate_limit":     utils.PathSearch("rate_limit", v, nil),
			"is_all":         utils.PathSearch("is_all", v, nil),
			"failed_num":     utils.PathSearch("failed_num", v, nil),
			"success_num":    utils.PathSearch("success_num", v, nil),
			"total_num":      utils.PathSearch("total_num", v, nil),
			"risky_num":      utils.PathSearch("risky_num", v, nil),
			"sync_task_type": utils.PathSearch("sync_task_type", v, nil),
			"failed_reason":  utils.PathSearch("failed_reason", v, nil),
			"failed_images":  flattenImageTasksImages(utils.PathSearch("failed_images", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenImageTasksImages(imageInfo []interface{}) []interface{} {
	if len(imageInfo) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(imageInfo))
	for _, v := range imageInfo {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"registry_id":   utils.PathSearch("registry_id", v, nil),
			"registry_name": utils.PathSearch("registry_name", v, nil),
			"image_name":    utils.PathSearch("image_name", v, nil),
			"image_version": utils.PathSearch("image_version", v, nil),
			"namespace":     utils.PathSearch("namespace", v, nil),
			"registry_type": utils.PathSearch("registry_type", v, nil),
			"failed_reason": utils.PathSearch("failed_reason", v, nil),
		})
	}

	return result
}
