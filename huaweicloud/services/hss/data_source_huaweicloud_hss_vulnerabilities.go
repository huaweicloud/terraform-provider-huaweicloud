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

// @API HSS GET /v5/{project_id}/vulnerability/vulnerabilities
func DataSourceVulnerabilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVulnerabilitiesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vul_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vul_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"handle_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repair_priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cve_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
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
						"vul_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vul_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repair_necessity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhandle_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"solution_detail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cve_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cve_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cvss": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"patch_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repair_priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hosts_num": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"important": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"common": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"test": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"repair_success_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"fixed_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ignored_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"verify_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"repair_priority_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repair_priority": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildVulnerabilitiesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=10"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("vul_id"); ok {
		queryParams = fmt.Sprintf("%s&vul_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("vul_name"); ok {
		queryParams = fmt.Sprintf("%s&vul_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("handle_status"); ok {
		queryParams = fmt.Sprintf("%s&handle_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("repair_priority"); ok {
		queryParams = fmt.Sprintf("%s&repair_priority=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cve_id"); ok {
		queryParams = fmt.Sprintf("%s&cve_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("label_list"); ok {
		queryParams = fmt.Sprintf("%s&label_list=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_name"); ok {
		queryParams = fmt.Sprintf("%s&group_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceVulnerabilitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/vulnerability/vulnerabilities"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildVulnerabilitiesQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS vulnerabilities: %s", err)
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
		d.Set("data_list", flattenVulnerabilitiesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVulnerabilitiesDataList(dataResp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"vul_id":            utils.PathSearch("vul_id", v, nil),
			"vul_name":          utils.PathSearch("vul_name", v, nil),
			"type":              utils.PathSearch("type", v, nil),
			"repair_necessity":  utils.PathSearch("repair_necessity", v, nil),
			"severity_level":    utils.PathSearch("severity_level", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"label_list":        utils.PathSearch("label_list", v, nil),
			"host_num":          utils.PathSearch("host_num", v, nil),
			"unhandle_host_num": utils.PathSearch("unhandle_host_num", v, nil),
			"solution_detail":   utils.PathSearch("solution_detail", v, nil),
			"url":               utils.PathSearch("url", v, nil),
			"host_id_list":      utils.PathSearch("host_id_list", v, nil),
			"cve_list": flattenVulnerabilitiesCveList(utils.PathSearch("cve_list", v,
				make([]interface{}, 0)).([]interface{})),
			"patch_url":          utils.PathSearch("patch_url", v, nil),
			"repair_priority":    utils.PathSearch("repair_priority", v, nil),
			"hosts_num":          flattenVulnerabilitiesHostsNum(utils.PathSearch("hosts_num", v, nil)),
			"repair_success_num": utils.PathSearch("repair_success_num", v, nil),
			"fixed_num":          utils.PathSearch("fixed_num", v, nil),
			"ignored_num":        utils.PathSearch("ignored_num", v, nil),
			"verify_num":         utils.PathSearch("verify_num", v, nil),
			"repair_priority_list": flattenRepairPriorityList(utils.PathSearch("repair_priority_list", v,
				make([]interface{}, 0)).([]interface{})),
			"scan_time": utils.PathSearch("scan_time", v, nil),
		})
	}

	return result
}

func flattenVulnerabilitiesCveList(rawCveList []interface{}) []interface{} {
	result := make([]interface{}, 0, len(rawCveList))
	for _, v := range rawCveList {
		result = append(result, map[string]interface{}{
			"cve_id": utils.PathSearch("cve_id", v, nil),
			"cvss":   utils.PathSearch("cvss", v, nil),
		})
	}

	return result
}

func flattenVulnerabilitiesHostsNum(rawHostsNum interface{}) []map[string]interface{} {
	result := map[string]interface{}{
		"important": utils.PathSearch("important", rawHostsNum, nil),
		"common":    utils.PathSearch("common", rawHostsNum, nil),
		"test":      utils.PathSearch("test", rawHostsNum, nil),
	}

	return []map[string]interface{}{result}
}

func flattenRepairPriorityList(rawRepairPriorityList []interface{}) []interface{} {
	result := make([]interface{}, 0, len(rawRepairPriorityList))
	for _, v := range rawRepairPriorityList {
		result = append(result, map[string]interface{}{
			"repair_priority": utils.PathSearch("repair_priority", v, nil),
			"host_num":        utils.PathSearch("host_num", v, nil),
		})
	}

	return result
}
