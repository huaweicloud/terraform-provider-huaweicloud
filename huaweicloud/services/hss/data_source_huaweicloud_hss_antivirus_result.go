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

// @API HSS GET /v5/{project_id}/antivirus/result
func DataSourceAntivirusResult() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAntivirusResultRead,

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
			"handle_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"malware_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"manual_isolate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"malware_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"malware_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_owner": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_attr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_ctime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_mtime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"resource_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip": {
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
									"asset_value": {
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
						"event_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"occur_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"handle_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memo": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operate_accept_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"operate_detail_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"keyword": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"isolate_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAntivirusResultQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
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
	if v, ok := d.GetOk("handle_status"); ok {
		queryParams = fmt.Sprintf("%s&handle_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity_list"); ok {
		severityList := v.([]interface{})
		for _, severity := range severityList {
			queryParams = fmt.Sprintf("%s&severity_list=%v", queryParams, severity)
		}
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}
	if v, ok := d.GetOk("malware_name"); ok {
		queryParams = fmt.Sprintf("%s&malware_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_path"); ok {
		queryParams = fmt.Sprintf("%s&file_path=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_hash"); ok {
		queryParams = fmt.Sprintf("%s&file_hash=%v", queryParams, v)
	}
	if v, ok := d.GetOk("task_name"); ok {
		queryParams = fmt.Sprintf("%s&task_name=%v", queryParams, v)
	}
	queryParams = fmt.Sprintf("%s&manual_isolate=%v", queryParams, d.Get("manual_isolate").(bool))

	return queryParams
}

func dataSourceAntivirusResultRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		result  = make([]interface{}, 0)
		offset  = 0
		httpUrl = "v5/{project_id}/antivirus/result"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAntivirusResultQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS antivirus result: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

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
		d.Set("data_list", flattenAntivirusResultDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAntivirusResultDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"result_id":           utils.PathSearch("result_id", v, nil),
			"malware_type":        utils.PathSearch("malware_type", v, nil),
			"malware_name":        utils.PathSearch("malware_name", v, nil),
			"severity":            utils.PathSearch("severity", v, nil),
			"task_id":             utils.PathSearch("task_id", v, nil),
			"task_name":           utils.PathSearch("task_name", v, nil),
			"file_info":           flattenAntivirusResultFileInfo(utils.PathSearch("file_info", v, nil)),
			"resource_info":       flattenAntivirusResultResourceInfo(utils.PathSearch("resource_info", v, nil)),
			"event_type":          utils.PathSearch("event_type", v, nil),
			"occur_time":          utils.PathSearch("occur_time", v, nil),
			"handle_status":       utils.PathSearch("handle_status", v, nil),
			"handle_method":       utils.PathSearch("handle_method", v, nil),
			"memo":                utils.PathSearch("memo", v, nil),
			"operate_accept_list": utils.ExpandToStringList(utils.PathSearch("operate_accept_list", v, make([]interface{}, 0)).([]interface{})),
			"operate_detail_list": flattenAntivirusResultOperateDetailList(
				utils.PathSearch("operate_detail_list", v, make([]interface{}, 0)).([]interface{})),
			"isolate_tag": utils.PathSearch("isolate_tag", v, nil),
		})
	}

	return rst
}

func flattenAntivirusResultFileInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson, ok := resp.(map[string]interface{})
	if !ok {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"file_path":  utils.PathSearch("file_path", curJson, nil),
			"file_hash":  utils.PathSearch("file_hash", curJson, nil),
			"file_size":  utils.PathSearch("file_size", curJson, nil),
			"file_owner": utils.PathSearch("file_owner", curJson, nil),
			"file_attr":  utils.PathSearch("file_attr", curJson, nil),
			"file_ctime": utils.PathSearch("file_ctime", curJson, nil),
			"file_mtime": utils.PathSearch("file_mtime", curJson, nil),
		},
	}
}

func flattenAntivirusResultResourceInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson, ok := resp.(map[string]interface{})
	if !ok {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"host_name":      utils.PathSearch("host_name", curJson, nil),
			"host_id":        utils.PathSearch("host_id", curJson, nil),
			"agent_id":       utils.PathSearch("agent_id", curJson, nil),
			"private_ip":     utils.PathSearch("private_ip", curJson, nil),
			"public_ip":      utils.PathSearch("public_ip", curJson, nil),
			"os_type":        utils.PathSearch("os_type", curJson, nil),
			"host_status":    utils.PathSearch("host_status", curJson, nil),
			"agent_status":   utils.PathSearch("agent_status", curJson, nil),
			"protect_status": utils.PathSearch("protect_status", curJson, nil),
			"asset_value":    utils.PathSearch("asset_value", curJson, nil),
			"os_name":        utils.PathSearch("os_name", curJson, nil),
			"os_version":     utils.PathSearch("os_version", curJson, nil),
		},
	}
}

func flattenAntivirusResultOperateDetailList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"keyword": utils.PathSearch("keyword", v, nil),
			"hash":    utils.PathSearch("hash", v, nil),
		})
	}

	return rst
}
