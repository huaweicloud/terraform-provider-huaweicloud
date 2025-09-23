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

// @API HSS GET /v5/{project_id}/host-management/hosts
func DataSourceHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detect_result": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_id": {
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
						"protect_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_charging_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detect_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asset_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"open_time": {
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
						"asset_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vulnerability_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"baseline_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"intrusion_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildHostsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=20"
	queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&host_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		queryParams = fmt.Sprintf("%s&os_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("agent_status"); ok {
		queryParams = fmt.Sprintf("%s&agent_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_status"); ok {
		queryParams = fmt.Sprintf("%s&protect_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_version"); ok {
		queryParams = fmt.Sprintf("%s&version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_charging_mode"); ok {
		queryParams = fmt.Sprintf("%s&charging_mode=%v", queryParams, convertChargingModeRequest(v.(string)))
	}
	if v, ok := d.GetOk("detect_result"); ok {
		queryParams = fmt.Sprintf("%s&detect_result=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_group_id"); ok {
		queryParams = fmt.Sprintf("%s&policy_group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
		httpUrl = "v5/{project_id}/host-management/hosts"
		offset  = 0
		result  = make([]interface{}, 0)
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildHostsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS hosts, %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		hostsResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(hostsResp) == 0 {
			break
		}

		result = append(result, hostsResp...)
		offset += len(hostsResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("hosts", flattenHosts(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHosts(hostsResp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(hostsResp))
	for _, v := range hostsResp {
		openTime := utils.PathSearch("open_time", v, float64(0)).(float64)
		result = append(result, map[string]interface{}{
			"id":                     utils.PathSearch("host_id", v, nil),
			"name":                   utils.PathSearch("host_name", v, nil),
			"status":                 utils.PathSearch("host_status", v, nil),
			"os_type":                utils.PathSearch("os_type", v, nil),
			"agent_id":               utils.PathSearch("agent_id", v, nil),
			"agent_status":           utils.PathSearch("agent_status", v, nil),
			"protect_status":         utils.PathSearch("protect_status", v, nil),
			"protect_version":        utils.PathSearch("version", v, nil),
			"protect_charging_mode":  flattenChargingMode(utils.PathSearch("charging_mode", v, "").(string)),
			"quota_id":               utils.PathSearch("resource_id", v, nil),
			"detect_result":          utils.PathSearch("detect_result", v, nil),
			"group_id":               utils.PathSearch("group_id", v, nil),
			"policy_group_id":        utils.PathSearch("policy_group_id", v, nil),
			"asset_value":            utils.PathSearch("asset_value", v, nil),
			"open_time":              utils.FormatTimeStampRFC3339(int64(openTime)/1000, false),
			"private_ip":             utils.PathSearch("private_ip", v, nil),
			"public_ip":              utils.PathSearch("public_ip", v, nil),
			"asset_risk_num":         utils.PathSearch("asset", v, nil),
			"vulnerability_risk_num": utils.PathSearch("vulnerability", v, nil),
			"baseline_risk_num":      utils.PathSearch("baseline", v, nil),
			"intrusion_risk_num":     utils.PathSearch("intrusion", v, nil),
			"enterprise_project_id":  utils.PathSearch("enterprise_project_id", v, nil),
		})
	}

	return result
}
