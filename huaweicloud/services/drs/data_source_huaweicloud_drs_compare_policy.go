package drs

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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/compare-policy
func DataSourceDrsComparePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsComparePolicyRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interval_hour": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"next_compare_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDrsComparePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/compare-policy"
		jobId   = d.Get("job_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS compare policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("interval_hour", utils.PathSearch("interval_hour", respBody, nil)),
		d.Set("period", utils.PathSearch("period", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("begin_time", utils.PathSearch("begin_time", respBody, nil)),
		d.Set("end_time", utils.PathSearch("end_time", respBody, nil)),
		d.Set("compare_type", utils.PathSearch("compare_type", respBody, nil)),
		d.Set("next_compare_time", utils.PathSearch("next_compare_time", respBody, nil)),
		d.Set("compare_policy", utils.PathSearch("compare_policy", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
