package aad

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v2/aad/domains/waf-info/flow/request/peak
func DataSourceCcAttackProtectionQPS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcAttackProtectionQPSRead,

		Schema: map[string]*schema.Schema{
			"recent": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specified the time range for querying CC attack protection QPS data.`,
			},
			"domains": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specified the domain name to query.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specified the start time for querying CC attack protection QPS data.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specified the end time for querying CC attack protection QPS data.`,
			},
			"overseas_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specified the protection region.`,
			},
			"qps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The QPS value of CC attack protection.`,
			},
		},
	}
}

func buildCcAttackProtectionQPSQueryParams(d *schema.ResourceData) string {
	rst := fmt.Sprintf("?recent=%s", d.Get("recent").(string))

	if v, ok := d.GetOk("domains"); ok {
		rst += fmt.Sprintf("&domains=%s", v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		rst += fmt.Sprintf("&start_time=%s", v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		rst += fmt.Sprintf("&end_time=%s", v.(string))
	}

	if v, ok := d.GetOk("overseas_type"); ok {
		rst += fmt.Sprintf("&overseas_type=%s", v.(string))
	}

	return rst
}

func dataSourceCcAttackProtectionQPSRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v2/aad/domains/waf-info/flow/request/peak"
	)

	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildCcAttackProtectionQPSQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD CC attack protection QPS data: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("qps", utils.PathSearch("qps", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
