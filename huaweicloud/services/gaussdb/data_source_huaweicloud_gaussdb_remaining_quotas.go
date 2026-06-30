package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/enterprise-projects/remaining-quotas
func DataSourceGaussDbRemainingQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbRemainingQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"eps_tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"eps_remaining_quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbRemainingQuotasQuotaSchema(),
			},
		},
	}
}

func gaussDbRemainingQuotasQuotaSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"eps_tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_eps_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_eps_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mem_eps_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_eps_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_eps_remaining_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_eps_remaining_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mem_eps_remaining_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_eps_remaining_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbRemainingQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/enterprise-projects/remaining-quotas"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getOpt.JSONBody = utils.RemoveNil(buildGaussDbRemainingQuotasBodyParams(d))
	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB remaining quotas: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("eps_remaining_quotas", flattenGaussDbRemainingQuotas(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaussDbRemainingQuotasBodyParams(d *schema.ResourceData) map[string]interface{} {
	epsTags := d.Get("eps_tags").([]interface{})
	tags := make([]string, len(epsTags))
	for i, v := range epsTags {
		tags[i] = v.(string)
	}
	return map[string]interface{}{
		"eps_tags": tags,
	}
}

func flattenGaussDbRemainingQuotas(resp interface{}) []interface{} {
	quotas := utils.PathSearch("eps_remaining_quotas", resp, make([]interface{}, 0)).([]interface{})
	if len(quotas) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(quotas))
	for _, quota := range quotas {
		result = append(result, map[string]interface{}{
			"eps_tag":                      utils.PathSearch("eps_tag", quota, nil),
			"instance_eps_quota":           utils.PathSearch("instance_eps_quota", quota, nil),
			"cpu_eps_quota":                utils.PathSearch("cpu_eps_quota", quota, nil),
			"mem_eps_quota":                utils.PathSearch("mem_eps_quota", quota, nil),
			"volume_eps_quota":             utils.PathSearch("volume_eps_quota", quota, nil),
			"instance_eps_remaining_quota": utils.PathSearch("instance_eps_remaining_quota", quota, nil),
			"cpu_eps_remaining_quota":      utils.PathSearch("cpu_eps_remaining_quota", quota, nil),
			"mem_eps_remaining_quota":      utils.PathSearch("mem_eps_remaining_quota", quota, nil),
			"volume_eps_remaining_quota":   utils.PathSearch("volume_eps_remaining_quota", quota, nil),
		})
	}

	return result
}
