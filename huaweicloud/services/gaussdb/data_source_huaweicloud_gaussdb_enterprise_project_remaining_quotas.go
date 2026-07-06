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
func DataSourceGaussDBEnterpriseProjectRemainingQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBEnterpriseProjectRemainingQuotasRead,

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
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eps_remaining_quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     epsRemainingQuotaSchema(),
			},
		},
	}
}

func epsRemainingQuotaSchema() *schema.Resource {
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

func dataSourceGaussDBEnterpriseProjectRemainingQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/enterprise-projects/remaining-quotas"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	opts.JSONBody = buildGetEnterpriseProjectRemainingQuotasBodyParams(d)

	requestResp, err := client.Request("POST", requestPath, &opts)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB enterprise project remaining quotas: %s", err)
	}

	respJson, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("job_id", utils.PathSearch("job_id", respJson, nil)),
		d.Set("eps_remaining_quotas", flattenGetEpsRemainingQuotasBody(respJson)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetEnterpriseProjectRemainingQuotasBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"eps_tags": d.Get("eps_tags"),
	}
}

func flattenGetEpsRemainingQuotasBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("eps_remaining_quotas", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"eps_tag":                      utils.PathSearch("eps_tag", v, nil),
			"instance_eps_quota":           utils.PathSearch("instance_eps_quota", v, nil),
			"cpu_eps_quota":                utils.PathSearch("cpu_eps_quota", v, nil),
			"mem_eps_quota":                utils.PathSearch("mem_eps_quota", v, nil),
			"volume_eps_quota":             utils.PathSearch("volume_eps_quota", v, nil),
			"instance_eps_remaining_quota": utils.PathSearch("instance_eps_remaining_quota", v, nil),
			"cpu_eps_remaining_quota":      utils.PathSearch("cpu_eps_remaining_quota", v, nil),
			"mem_eps_remaining_quota":      utils.PathSearch("mem_eps_remaining_quota", v, nil),
			"volume_eps_remaining_quota":   utils.PathSearch("volume_eps_remaining_quota", v, nil),
		})
	}
	return res
}
