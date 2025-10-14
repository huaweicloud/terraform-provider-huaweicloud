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

// @API HSS GET /v5/{project_id}/webtamper/{host_id}/policy
func DataSourceWebTamperPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebTamperPolicyRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_dir_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protect_dir_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protect_dir_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protect_dir": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exclude_child_dir": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exclue_file_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_backup_dir": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protect_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"exclude_file_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_timing_off": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"timing_off_config_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"week_off_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"timing_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"enable_rasp_protect": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rasp_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_privileged_process": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"privileged_child_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"privileged_process_path_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildWebTamperPolicyQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceWebTamperPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		hostId  = d.Get("host_id").(string)
		httpUrl = "v5/{project_id}/webtamper/{host_id}/policy"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{host_id}", hostId)
	requestPath += buildWebTamperPolicyQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS web tamper policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("protect_dir_num", utils.PathSearch("protect_dir_num", respBody, nil)),
		d.Set("protect_dir_info", flattenProtectDirInfo(respBody)),
		d.Set("enable_timing_off", utils.PathSearch("enable_timing_off", respBody, nil)),
		d.Set("timing_off_config_info", flattenTimingOffConfigInfo(respBody)),
		d.Set("enable_rasp_protect", utils.PathSearch("enable_rasp_protect", respBody, nil)),
		d.Set("rasp_path", utils.PathSearch("rasp_path", respBody, nil)),
		d.Set("enable_privileged_process", utils.PathSearch("enable_privileged_process", respBody, nil)),
		d.Set("privileged_child_status", utils.PathSearch("privileged_child_status", respBody, nil)),
		d.Set("privileged_process_path_list", utils.PathSearch("privileged_process_path_list", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenProtectDirInfo(resp interface{}) []interface{} {
	dirInfoResp := utils.PathSearch("protect_dir_info", resp, nil)
	if dirInfoResp == nil {
		return nil
	}

	dirListResp := utils.PathSearch("protect_dir_list", dirInfoResp, make([]interface{}, 0)).([]interface{})
	dirListResult := make([]map[string]interface{}, 0, len(dirListResp))
	for _, v := range dirListResp {
		dirListResult = append(dirListResult, map[string]interface{}{
			"protect_dir":       utils.PathSearch("protect_dir", v, nil),
			"exclude_child_dir": utils.PathSearch("exclude_child_dir", v, nil),
			"exclue_file_path":  utils.PathSearch("exclue_file_path", v, nil),
			"local_backup_dir":  utils.PathSearch("local_backup_dir", v, nil),
			"protect_status":    utils.PathSearch("protect_status", v, nil),
			"error":             utils.PathSearch("error", v, nil),
		})
	}

	return []interface{}{
		map[string]interface{}{
			"protect_dir_list":  dirListResult,
			"exclude_file_type": utils.PathSearch("exclude_file_type", dirInfoResp, nil),
			"protect_mode":      utils.PathSearch("protect_mode", dirInfoResp, nil),
		},
	}
}

func flattenTimingOffConfigInfo(resp interface{}) []interface{} {
	infoResp := utils.PathSearch("timing_off_config_info", resp, nil)
	if infoResp == nil {
		return nil
	}

	rangeListResp := utils.PathSearch("timing_range_list", infoResp, make([]interface{}, 0)).([]interface{})
	rangeListResult := make([]map[string]interface{}, 0, len(rangeListResp))
	for _, v := range rangeListResp {
		rangeListResult = append(rangeListResult, map[string]interface{}{
			"time_range":  utils.PathSearch("time_range", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}

	return []interface{}{
		map[string]interface{}{
			"week_off_list":     flattenWeekOffList(utils.PathSearch("week_off_list", infoResp, make([]interface{}, 0)).([]interface{})),
			"timing_range_list": rangeListResult,
		},
	}
}

func flattenWeekOffList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		if num, ok := v.(float64); ok {
			result = append(result, int(num))
		} else if num, ok := v.(int); ok {
			result = append(result, num)
		}
	}

	return result
}
