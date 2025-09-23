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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/billing/quotas-detail
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"used_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quota_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charging_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shared_quota": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": common.TagsComputedSchema(),
					},
				},
			},
		},
	}
}

func buildQuotasQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=20"
	queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	if v, ok := d.GetOk("version"); ok {
		queryParams = fmt.Sprintf("%s&version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&quota_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("used_status"); ok {
		queryParams = fmt.Sprintf("%s&used_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("quota_id"); ok {
		queryParams = fmt.Sprintf("%s&resource_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("charging_mode"); ok {
		queryParams = fmt.Sprintf("%s&charging_mode=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d, "all_granted_eps")
		product = "hss"
		httpUrl = "v5/{project_id}/billing/quotas-detail"
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
	getPath += buildQuotasQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS quotas, %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		quotasResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(quotasResp) == 0 {
			break
		}

		result = append(result, quotasResp...)
		offset += len(quotasResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenQuotas(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQuotas(quotasResp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(quotasResp))
	for _, v := range quotasResp {
		expireTime := utils.PathSearch("expire_time", v, float64(0)).(float64)
		result = append(result, map[string]interface{}{
			"id":                      utils.PathSearch("resource_id", v, nil),
			"version":                 utils.PathSearch("version", v, nil),
			"status":                  utils.PathSearch("quota_status", v, nil),
			"used_status":             utils.PathSearch("used_status", v, nil),
			"host_id":                 utils.PathSearch("host_id", v, nil),
			"host_name":               utils.PathSearch("host_name", v, nil),
			"charging_mode":           flattenChargingMode(utils.PathSearch("charging_mode", v, "").(string)),
			"expire_time":             utils.FormatTimeStampRFC3339(int64(expireTime)/1000, false),
			"shared_quota":            utils.PathSearch("shared_quota", v, nil),
			"enterprise_project_id":   utils.PathSearch("enterprise_project_id", v, nil),
			"enterprise_project_name": utils.PathSearch("enterprise_project_name", v, nil),
			"tags":                    utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
		})
	}

	return result
}
