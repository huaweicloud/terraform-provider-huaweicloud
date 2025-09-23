package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LIVE GET /v1/{project_id}/stream/snapshot
func DataSourceLiveSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLiveSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the domain name.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application name.`,
			},
			"snapshots": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The snapshot list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"call_back_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the server for receiving callback notifications.`,
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ingest domain name.`,
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application name.`,
						},
						"call_back_auth_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The callback authentication key value.`,
						},
						"frequency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The snapshot capturing frequency.`,
						},
						"storage_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The method for storing snapshots in an OBS bucket.`,
						},
						"call_back_enabled": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Whether to enable callback notification.`,
						},
						"storage_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS bucket name.`,
						},
						"storage_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The region where the OBS bucket is located.`,
						},
						"storage_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS object path.`,
						},
					},
				},
			},
		},
	}
}

func buildSnapshotsQueryParams(d *schema.ResourceData) string {
	rst := fmt.Sprintf("?domain=%s&limit=100", d.Get("domain_name").(string))
	if appName, ok := d.GetOk("app_name"); ok {
		rst = fmt.Sprintf("%s&app_name=%s", rst, appName.(string))
	}
	return rst
}

func dataSourceLiveSnapshotsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/stream/snapshot"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildSnapshotsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// The pagination query of this query API is not accurate, so pagination query is not supported for the time being.
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving Live snapshots: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshots := utils.PathSearch("snapshot_config_list", respBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("snapshots", flattenSnapshots(snapshots)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSnapshots(totalSnapshots []interface{}) []interface{} {
	if len(totalSnapshots) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(totalSnapshots))
	for _, snapshot := range totalSnapshots {
		result = append(result, map[string]interface{}{
			"domain_name":        utils.PathSearch("domain", snapshot, nil),
			"app_name":           utils.PathSearch("app_name", snapshot, nil),
			"frequency":          utils.PathSearch("time_interval", snapshot, nil),
			"storage_mode":       utils.PathSearch("object_write_mode", snapshot, nil),
			"storage_bucket":     utils.PathSearch("obs_location.bucket", snapshot, nil),
			"storage_location":   utils.PathSearch("obs_location.location", snapshot, nil),
			"storage_path":       utils.PathSearch("obs_location.object", snapshot, nil),
			"call_back_enabled":  utils.PathSearch("call_back_enable", snapshot, nil),
			"call_back_url":      utils.PathSearch("call_back_url", snapshot, nil),
			"call_back_auth_key": utils.PathSearch("auth_key", snapshot, nil),
		})
	}
	return result
}
