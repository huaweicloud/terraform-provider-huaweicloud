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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/replication/publications/{publication_id}/monitor
func DataSourceRdsPublicationMonitor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsPublicationMonitorRead,

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
			"publication_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"worst_latency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"best_latency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"average_latency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_dist_sync": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replicated_transactions": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replication_rate_trans": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsPublicationMonitorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/replication/publications/{publication_id}/monitor"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{publication_id}", d.Get("publication_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS publication monitor: %s", err)
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
		d.Set("worst_latency", utils.PathSearch("worst_latency", getRespBody, nil)),
		d.Set("best_latency", utils.PathSearch("best_latency", getRespBody, nil)),
		d.Set("average_latency", utils.PathSearch("average_latency", getRespBody, nil)),
		d.Set("last_dist_sync", utils.PathSearch("last_dist_sync", getRespBody, nil)),
		d.Set("replicated_transactions", utils.PathSearch("replicated_transactions", getRespBody, nil)),
		d.Set("replication_rate_trans", utils.PathSearch("replication_rate_trans", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
