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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/replication/subscriptions/{subscription_id}/monitor
func DataSourceRdsSubscriptionMonitor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsSubscriptionMonitorRead,

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
			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"agent_not_running": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pending_cmd_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_dist_sync": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"estimated_process_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsSubscriptionMonitorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/replication/subscriptions/{subscription_id}/monitor"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{subscription_id}", d.Get("subscription_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS subscription monitor: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("latency", utils.PathSearch("latency", getRespBody, nil)),
		d.Set("agent_not_running", utils.PathSearch("agent_not_running", getRespBody, nil)),
		d.Set("pending_cmd_count", utils.PathSearch("pending_cmd_count", getRespBody, nil)),
		d.Set("last_dist_sync", utils.PathSearch("last_dist_sync", getRespBody, nil)),
		d.Set("estimated_process_time", utils.PathSearch("estimated_process_time", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
