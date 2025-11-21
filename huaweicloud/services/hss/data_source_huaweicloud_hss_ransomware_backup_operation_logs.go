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

// @API HSS GET /v5/{project_id}/ransomware/backup/operation-logs
func DataSourceRansomwareBackupOperationLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRansomwareBackupOperationLogsRead,

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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_days": {
				Type:     schema.TypeInt,
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
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"started_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ended_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
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

func buildRansomwareBackupOperationLogsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("resource_name"); ok {
		queryParams = fmt.Sprintf("%s&resource_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("last_days"); ok {
		queryParams = fmt.Sprintf("%s&last_days=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceRansomwareBackupOperationLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/ransomware/backup/operation-logs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildRansomwareBackupOperationLogsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS ransomware backup operation logs: %s", err)
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
		d.Set("data_list", flattenRansomwareBackupOperationLogsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRansomwareBackupOperationLogsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_id":     utils.PathSearch("host_id", v, nil),
			"host_name":   utils.PathSearch("host_name", v, nil),
			"backup_name": utils.PathSearch("backup_name", v, nil),
			"process":     utils.PathSearch("process", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"started_at":  utils.PathSearch("started_at", v, nil),
			"ended_at":    utils.PathSearch("ended_at", v, nil),
			"error_info":  flattenRansomwareBackupOperationLogsErrorInfo(utils.PathSearch("error_info", v, nil)),
		})
	}

	return rst
}

func flattenRansomwareBackupOperationLogsErrorInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"code":    utils.PathSearch("code", resp, nil),
			"message": utils.PathSearch("message", resp, nil),
		},
	}
}
