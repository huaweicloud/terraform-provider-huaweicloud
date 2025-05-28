package aad

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v1/unblockservice/{domain_id}/block-statistics
func DataSourceAADBlockStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAADBlockStatisticsRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specified the account ID of IAM user.",
			},
			"total_unblocking_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total unblocking times.",
			},
			"manual_unblocking_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The manual unblocking times.",
			},
			"automatic_unblocking_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The automatic unblocking times.",
			},
			"current_blocked_ip_numbers": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The current blocked IP number.",
			},
		},
	}
}

func dataSourceAADBlockStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)

	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	httpUrl := "v1/unblockservice/{domain_id}/block-statistics"
	httpUrl = strings.ReplaceAll(httpUrl, "{domain_id}", d.Get("domain_id").(string))
	listPath := client.Endpoint + httpUrl
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving AAD block statistics: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("total_unblocking_times", utils.PathSearch("total_unblocking_times", respBody, nil)),
		d.Set("manual_unblocking_times", utils.PathSearch("manual_unblocking_times", respBody, nil)),
		d.Set("current_blocked_ip_numbers", utils.PathSearch("current_blocked_ip_numbers", respBody, nil)),
		d.Set("automatic_unblocking_times", utils.PathSearch("automatic_unblocking_times", respBody, nil)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error retrieving AAD block statistics: %s", mErr)
	}

	return nil
}
