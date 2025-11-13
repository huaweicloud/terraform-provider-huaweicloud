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

// @API HSS GET /v5/{project_id}/antivirus/handle-history
func DataSourceAntivirusHandleHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAntivirusHandleHistoryRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"malware_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"handle_method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
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
						"file_path": {
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
						"notes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAntivirusHandleHistoryQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("malware_name"); ok {
		queryParams = fmt.Sprintf("%s&malware_name=%v", queryParams, v)
	}

	if v, ok := d.GetOk("file_path"); ok {
		queryParams = fmt.Sprintf("%s&file_path=%v", queryParams, v)
	}

	if v, ok := d.GetOk("severity_list"); ok {
		severityList := v.([]interface{})
		for _, severity := range severityList {
			queryParams = fmt.Sprintf("%s&severity_list=%v", queryParams, severity)
		}
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

	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}

	if v, ok := d.GetOk("handle_method"); ok {
		queryParams = fmt.Sprintf("%s&handle_method=%v", queryParams, v)
	}

	if v, ok := d.GetOk("user_name"); ok {
		queryParams = fmt.Sprintf("%s&user_name=%v", queryParams, v)
	}

	if v, ok := d.GetOk("event_type"); ok {
		queryParams = fmt.Sprintf("%s&event_type=%v", queryParams, v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceAntivirusHandleHistoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		result  = make([]interface{}, 0)
		offset  = 0
		httpUrl = "v5/{project_id}/antivirus/handle-history"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAntivirusHandleHistoryQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving antivirus handle history: %s", err)
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
		d.Set("data_list", flattenAntivirusHandleHistory(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAntivirusHandleHistory(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"result_id":     utils.PathSearch("result_id", v, nil),
			"malware_type":  utils.PathSearch("malware_type", v, nil),
			"malware_name":  utils.PathSearch("malware_name", v, nil),
			"severity":      utils.PathSearch("severity", v, nil),
			"file_path":     utils.PathSearch("file_path", v, nil),
			"host_name":     utils.PathSearch("host_name", v, nil),
			"private_ip":    utils.PathSearch("private_ip", v, nil),
			"public_ip":     utils.PathSearch("public_ip", v, nil),
			"asset_value":   utils.PathSearch("asset_value", v, nil),
			"occur_time":    utils.PathSearch("occur_time", v, nil),
			"handle_status": utils.PathSearch("handle_status", v, nil),
			"handle_method": utils.PathSearch("handle_method", v, nil),
			"notes":         utils.PathSearch("notes", v, nil),
			"handle_time":   utils.PathSearch("handle_time", v, nil),
			"user_name":     utils.PathSearch("user_name", v, nil),
		})
	}

	return rst
}
