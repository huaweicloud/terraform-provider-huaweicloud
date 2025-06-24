package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/recovery-time
func DataSourceRdsWalLogRecoveryTimeWindow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsWalLogRecoveryTimeWindowRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"recovery_min_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recovery_max_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsWalLogRecoveryTimeWindowRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/recovery-time"
	getUrl := client.Endpoint + httpUrl
	getUrl = strings.ReplaceAll(getUrl, "{project_id}", client.ProjectID)
	getUrl = strings.ReplaceAll(getUrl, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getUrl, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS WAL log recovery time window: %s", err)
	}

	body, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("recovery_min_time", utils.PathSearch("recovery_min_time", body, nil)),
		d.Set("recovery_max_time", utils.PathSearch("recovery_max_time", body, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
