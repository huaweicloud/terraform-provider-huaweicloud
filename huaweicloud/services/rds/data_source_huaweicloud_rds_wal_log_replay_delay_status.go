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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/replay-delay/show
func DataSourceRdsWalLogReplayDelayStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsWalLogReplayDelayStatusRead,
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
			"cur_delay_time_mills": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"delay_time_value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_delay_time_mills": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cur_log_replay_paused": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"latest_receive_log": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_replay_log": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsWalLogReplayDelayStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	url := "v3/{project_id}/instances/{instance_id}/replay-delay/show"
	getUrl := client.Endpoint + url
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
		return diag.Errorf("error retrieving RDS WAL log replay delay status: %s", err)
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
		d.Set("cur_delay_time_mills", utils.PathSearch("cur_delay_time_mills", body, nil)),
		d.Set("delay_time_value_range", utils.PathSearch("delay_time_value_range", body, nil)),
		d.Set("real_delay_time_mills", utils.PathSearch("real_delay_time_mills", body, nil)),
		d.Set("cur_log_replay_paused", utils.PathSearch("cur_log_replay_paused", body, nil)),
		d.Set("latest_receive_log", utils.PathSearch("latest_receive_log", body, nil)),
		d.Set("latest_replay_log", utils.PathSearch("latest_replay_log", body, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
