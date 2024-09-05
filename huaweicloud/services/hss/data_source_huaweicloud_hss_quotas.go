package hss

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

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

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		epsId     = cfg.GetEnterpriseProjectID(d, "all_granted_eps")
		limit     = int32(20)
		offset    int32
		allQuotas []hssv5model.QuotaResourcesResponseInfo
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	for {
		request := hssv5model.ListQuotasDetailRequest{
			Region:              &region,
			EnterpriseProjectId: utils.String(epsId),
			Limit:               utils.Int32(limit),
			Offset:              utils.Int32(offset),
			Category:            utils.StringIgnoreEmpty(d.Get("category").(string)),
			Version:             utils.StringIgnoreEmpty(d.Get("version").(string)),
			QuotaStatus:         utils.StringIgnoreEmpty(d.Get("status").(string)),
			UsedStatus:          utils.StringIgnoreEmpty(d.Get("used_status").(string)),
			HostName:            utils.StringIgnoreEmpty(d.Get("host_name").(string)),
			ResourceId:          utils.StringIgnoreEmpty(d.Get("quota_id").(string)),
			ChargingMode:        utils.StringIgnoreEmpty(convertChargingModeRequest(d.Get("charging_mode").(string))),
		}

		listResp, listErr := client.ListQuotasDetail(&request)
		if listErr != nil {
			return diag.Errorf("error querying HSS quotas: %s", listErr)
		}

		if listResp == nil || listResp.DataList == nil {
			break
		}
		if len(*listResp.DataList) == 0 {
			break
		}

		allQuotas = append(allQuotas, *listResp.DataList...)
		offset += limit
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenQuotas(allQuotas)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQuotas(quotas []hssv5model.QuotaResourcesResponseInfo) []interface{} {
	if len(quotas) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(quotas))
	for _, v := range quotas {
		expireTime := ""
		if v.ExpireTime != nil {
			expireTime = utils.FormatTimeStampRFC3339(*(v.ExpireTime)/1000, false)
		}

		rst = append(rst, map[string]interface{}{
			"id":                      v.ResourceId,
			"version":                 v.Version,
			"status":                  v.QuotaStatus,
			"used_status":             v.UsedStatus,
			"host_id":                 v.HostId,
			"host_name":               v.HostName,
			"charging_mode":           convertChargingMode(v.ChargingMode),
			"expire_time":             expireTime,
			"shared_quota":            v.SharedQuota,
			"enterprise_project_id":   v.EnterpriseProjectId,
			"enterprise_project_name": v.EnterpriseProjectName,
			"tags":                    flattenTags(v.Tags),
		})
	}

	return rst
}

func flattenTags(tags *[]hssv5model.TagInfo) map[string]interface{} {
	if tags == nil {
		return nil
	}

	rst := make(map[string]interface{})
	for _, tag := range *tags {
		if tag.Key != nil {
			rst[*tag.Key] = tag.Value
		}
	}

	return rst
}
