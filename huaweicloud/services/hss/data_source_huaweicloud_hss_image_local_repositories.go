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

// @API HSS GET /v5/{project_id}/image/local-repositories
func DataSourceImageLocalRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageLocalRepositoriesRead,

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
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scan_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_image_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start_latest_update_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_latest_update_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start_latest_scan_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_latest_scan_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"has_vul": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pod_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pod_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"has_container": {
				Type:     schema.TypeBool,
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
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_digest": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scan_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"latest_update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"latest_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vul_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unsafe_setting_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"malicious_file_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"container_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"component_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scan_failed_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"non_scan_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildImageLocalRepositoriesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("image_name"); ok {
		queryParams = fmt.Sprintf("%s&image_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_version"); ok {
		queryParams = fmt.Sprintf("%s&image_version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("scan_status"); ok {
		queryParams = fmt.Sprintf("%s&scan_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("local_image_type"); ok {
		queryParams = fmt.Sprintf("%s&local_image_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_size"); ok {
		queryParams = fmt.Sprintf("%s&image_size=%v", queryParams, v)
	}
	if v, ok := d.GetOk("start_latest_update_time"); ok {
		queryParams = fmt.Sprintf("%s&start_latest_update_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_latest_update_time"); ok {
		queryParams = fmt.Sprintf("%s&end_latest_update_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("start_latest_scan_time"); ok {
		queryParams = fmt.Sprintf("%s&start_latest_scan_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_latest_scan_time"); ok {
		queryParams = fmt.Sprintf("%s&end_latest_scan_time=%v", queryParams, v)
	}

	queryParams = fmt.Sprintf("%s&has_vul=%v", queryParams, d.Get("has_vul").(bool))

	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_ip"); ok {
		queryParams = fmt.Sprintf("%s&host_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("container_id"); ok {
		queryParams = fmt.Sprintf("%s&container_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("container_name"); ok {
		queryParams = fmt.Sprintf("%s&container_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("pod_id"); ok {
		queryParams = fmt.Sprintf("%s&pod_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("pod_name"); ok {
		queryParams = fmt.Sprintf("%s&pod_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("app_name"); ok {
		queryParams = fmt.Sprintf("%s&app_name=%v", queryParams, v)
	}

	queryParams = fmt.Sprintf("%s&has_container=%v", queryParams, d.Get("has_container").(bool))

	return queryParams
}

func dataSourceImageLocalRepositoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/image/local-repositories"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildImageLocalRepositoriesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS image local repositories: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenImageLocalRepositoriesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageLocalRepositoriesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"image_name":         utils.PathSearch("image_name", v, nil),
			"image_id":           utils.PathSearch("image_id", v, nil),
			"image_digest":       utils.PathSearch("image_digest", v, nil),
			"image_version":      utils.PathSearch("image_version", v, nil),
			"local_image_type":   utils.PathSearch("local_image_type", v, nil),
			"scan_status":        utils.PathSearch("scan_status", v, nil),
			"image_size":         utils.PathSearch("image_size", v, nil),
			"latest_update_time": utils.PathSearch("latest_update_time", v, nil),
			"latest_scan_time":   utils.PathSearch("latest_scan_time", v, nil),
			"vul_num":            utils.PathSearch("vul_num", v, nil),
			"unsafe_setting_num": utils.PathSearch("unsafe_setting_num", v, nil),
			"malicious_file_num": utils.PathSearch("malicious_file_num", v, nil),
			"host_num":           utils.PathSearch("host_num", v, nil),
			"container_num":      utils.PathSearch("container_num", v, nil),
			"component_num":      utils.PathSearch("component_num", v, nil),
			"scan_failed_desc":   utils.PathSearch("scan_failed_desc", v, nil),
			"severity_level":     utils.PathSearch("severity_level", v, nil),
			"host_name":          utils.PathSearch("host_name", v, nil),
			"host_id":            utils.PathSearch("host_id", v, nil),
			"agent_id":           utils.PathSearch("agent_id", v, nil),
			"non_scan_reason":    utils.PathSearch("non_scan_reason", v, nil),
		})
	}

	return rst
}
