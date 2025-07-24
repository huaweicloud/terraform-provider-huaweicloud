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

// @API HSS GET /v5/{project_id}/antivirus/policy
func DataSourceAntivirusCustomScanPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAntivirusCustomScanPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
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
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scan_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scan_period_date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scan_hour": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scan_minute": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"next_start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scan_dir": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ignore_dir": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invalidate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"host_num": {
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
								},
							},
						},
						"whether_paid_task": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"file_type_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
		},
	}
}

func buildAntivirusCustomScanPoliciesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("policy_name"); ok {
		queryParams = fmt.Sprintf("%s&policy_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceAntivirusCustomScanPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/antivirus/policy"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAntivirusCustomScanPoliciesQueryParams(d, epsId)
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS antivirus custom scan policies: %s", err)
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
		d.Set("data_list", flattenAntivirusCustomScanPoliciesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAntivirusCustomScanPoliciesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"policy_id":        utils.PathSearch("policy_id", v, nil),
			"policy_name":      utils.PathSearch("policy_name", v, nil),
			"start_type":       utils.PathSearch("start_type", v, nil),
			"scan_period":      utils.PathSearch("scan_period", v, nil),
			"scan_period_date": utils.PathSearch("scan_period_date", v, nil),
			"scan_time":        utils.PathSearch("scan_time", v, nil),
			"scan_hour":        utils.PathSearch("scan_hour", v, nil),
			"scan_minute":      utils.PathSearch("scan_minute", v, nil),
			"next_start_time":  utils.PathSearch("next_start_time", v, nil),
			"scan_dir":         utils.PathSearch("scan_dir", v, nil),
			"ignore_dir":       utils.PathSearch("ignore_dir", v, nil),
			"action":           utils.PathSearch("action", v, nil),
			"invalidate":       utils.PathSearch("invalidate", v, nil),
			"host_num":         utils.PathSearch("host_num", v, nil),
			"host_info_list": flattenAntivirusCustomScanPoliciesHostInfoList(utils.PathSearch(
				"host_info_list", v, make([]interface{}, 0)).([]interface{})),
			"whether_paid_task": utils.PathSearch("whether_paid_task", v, nil),
			"file_type_list": utils.ExpandToInt32List(utils.PathSearch(
				"file_type_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenAntivirusCustomScanPoliciesHostInfoList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_id":     utils.PathSearch("host_id", v, nil),
			"host_name":   utils.PathSearch("host_name", v, nil),
			"private_ip":  utils.PathSearch("private_ip", v, nil),
			"public_ip":   utils.PathSearch("public_ip", v, nil),
			"asset_value": utils.PathSearch("asset_value", v, nil),
		})
	}

	return rst
}
