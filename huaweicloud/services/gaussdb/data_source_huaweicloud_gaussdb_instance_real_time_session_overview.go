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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/session-overview
func DataSourceGaussDBInstanceRealTimeSessionOverview() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBInstanceRealTimeSessionOverviewRead,

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
			"active_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slow_sql_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDBInstanceRealTimeSessionOverviewRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/session-overview"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &opts)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB instance real time session overview: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
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
		d.Set("active_num", utils.PathSearch("active_num", respBody, nil)),
		d.Set("total_num", utils.PathSearch("total_num", respBody, nil)),
		d.Set("slow_sql_num", utils.PathSearch("slow_sql_num", respBody, nil)),
		d.Set("lock_num", utils.PathSearch("lock_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
