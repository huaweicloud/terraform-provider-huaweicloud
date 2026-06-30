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

// @API GaussDB GET /v3.5/{project_id}/instances/{instance_id}/disaster-recovery/monitor
func DataSourceGaussDbInstanceDrStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbInstanceDrStatusRead,

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
			"disaster_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rpo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rto": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rpo_threshold": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rto_threshold": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switchover_progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failover_progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbInstanceDrStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.5/{project_id}/instances/{instance_id}/disaster-recovery/monitor"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += "?disaster_type=" + d.Get("disaster_type").(string)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB instance DR status: %s", err)
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
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("rpo", utils.PathSearch("rpo", getRespBody, nil)),
		d.Set("rto", utils.PathSearch("rto", getRespBody, nil)),
		d.Set("rpo_threshold", utils.PathSearch("rpo_threshold", getRespBody, nil)),
		d.Set("rto_threshold", utils.PathSearch("rto_threshold", getRespBody, nil)),
		d.Set("switchover_progress", utils.PathSearch("switchover_progress", getRespBody, nil)),
		d.Set("failover_progress", utils.PathSearch("failover_progress", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
