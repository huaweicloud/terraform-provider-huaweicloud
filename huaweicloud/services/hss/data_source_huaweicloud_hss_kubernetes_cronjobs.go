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

// @API HSS GET /v5/{project_id}/kubernetes/cronjobs
func DataSourceKubernetesCronJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesCronJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cronjob_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cronjob_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"running_jobs_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"schedule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"match_labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"val": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"execute_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildKubernetesCronJobsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("cronjob_name"); ok {
		queryParams = fmt.Sprintf("%s&cronjob_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("namespace_name"); ok {
		queryParams = fmt.Sprintf("%s&namespace_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		queryParams = fmt.Sprintf("%s&cluster_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceKubernetesCronJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		product        = "hss"
		epsId          = cfg.GetEnterpriseProjectID(d)
		result         = make([]interface{}, 0)
		offset         = 0
		totalNum       float64
		lastUpdateTime float64
		httpUrl        = "v5/{project_id}/kubernetes/cronjobs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildKubernetesCronJobsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS kubernetes cronjobs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		lastUpdateTime = utils.PathSearch("last_update_time", respBody, float64(0)).(float64)
		cronjobInfoList := utils.PathSearch("cronjob_info_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(cronjobInfoList) == 0 {
			break
		}

		result = append(result, cronjobInfoList...)
		offset += len(cronjobInfoList)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("last_update_time", lastUpdateTime),
		d.Set("cronjob_info_list", flattenKubernetesCronJobsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKubernetesCronJobsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":             utils.PathSearch("name", v, nil),
			"namespace_name":   utils.PathSearch("namespace_name", v, nil),
			"cluster_name":     utils.PathSearch("cluster_name", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"running_jobs_num": utils.PathSearch("running_jobs_num", v, nil),
			"schedule":         utils.PathSearch("schedule", v, nil),
			"image_name":       utils.PathSearch("image_name", v, nil),
			"match_labels":     flattenMatchLabels(utils.PathSearch("match_labels", v, make([]interface{}, 0)).([]interface{})),
			"execute_time":     utils.PathSearch("execute_time", v, nil),
			"create_time":      utils.PathSearch("create_time", v, nil),
		})
	}

	return rst
}

func flattenMatchLabels(labels []interface{}) []interface{} {
	if len(labels) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(labels))
	for _, label := range labels {
		result = append(result, map[string]interface{}{
			"key": utils.PathSearch("key", label, nil),
			"val": utils.PathSearch("val", label, nil),
		})
	}

	return result
}
