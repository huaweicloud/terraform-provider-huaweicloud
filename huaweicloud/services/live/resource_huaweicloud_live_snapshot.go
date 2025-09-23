// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Live
// ---------------------------------------------------------------

package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live DELETE /v1/{project_id}/stream/snapshot
// @API Live GET /v1/{project_id}/stream/snapshot
// @API Live POST /v1/{project_id}/stream/snapshot
// @API Live PUT /v1/{project_id}/stream/snapshot
func ResourceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotCreate,
		UpdateContext: resourceSnapshotUpdate,
		ReadContext:   resourceSnapshotRead,
		DeleteContext: resourceSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ingest domain name.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the application name.`,
			},
			"frequency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the screenshot frequency.`,
			},
			"storage_mode": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the store mode in OBS bucket.`,
			},
			"storage_bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the bucket name of the OBS.`,
			},
			"storage_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the path of OBS object.`,
			},
			"call_back_enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to enable callback notifications.`,
			},
			"call_back_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the notification server address. `,
			},
			"call_back_auth_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the callback authentication key value.`,
			},
		},
	}
}

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/stream/snapshot"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateLiveSnapshotBodyParams(d, region)),
	}
	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating Live snapshot: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	appName := d.Get("app_name").(string)

	d.SetId(domainName + "/" + appName)

	return resourceSnapshotRead(ctx, d, meta)
}

func buildCreateLiveSnapshotBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":            utils.ValueIgnoreEmpty(d.Get("domain_name")),
		"app_name":          utils.ValueIgnoreEmpty(d.Get("app_name")),
		"time_interval":     utils.ValueIgnoreEmpty(d.Get("frequency")),
		"object_write_mode": utils.ValueIgnoreEmpty(d.Get("storage_mode")),
		"obs_location":      buildLiveSnapshotObsLocationChildBody(d, region),
		"call_back_enable":  utils.ValueIgnoreEmpty(d.Get("call_back_enabled")),
		"call_back_url":     utils.ValueIgnoreEmpty(d.Get("call_back_url")),
		"auth_key":          utils.ValueIgnoreEmpty(d.Get("call_back_auth_key")),
	}
	return bodyParams
}

func buildLiveSnapshotObsLocationChildBody(d *schema.ResourceData, region string) map[string]interface{} {
	params := map[string]interface{}{
		"bucket":   utils.ValueIgnoreEmpty(d.Get("storage_bucket")),
		"location": region,
		"object":   utils.ValueIgnoreEmpty(d.Get("storage_path")),
	}
	return params
}

func resourceSnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/stream/snapshot"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, want '<domain_name>/<app_name>', but got '%s'", d.Id())
	}
	domainName := parts[0]
	appName := parts[1]

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildGetSnapshotQueryParams(domainName, appName)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live snapshot")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshot := utils.PathSearch("snapshot_config_list|[0]", respBody, nil)
	if snapshot == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("domain", snapshot, nil)),
		d.Set("app_name", utils.PathSearch("app_name", snapshot, nil)),
		d.Set("frequency", utils.PathSearch("time_interval", snapshot, nil)),
		d.Set("storage_mode", utils.PathSearch("object_write_mode", snapshot, nil)),
		d.Set("storage_bucket", utils.PathSearch("obs_location.bucket", snapshot, nil)),
		d.Set("storage_path", utils.PathSearch("obs_location.object", snapshot, nil)),
		d.Set("call_back_enabled", utils.PathSearch("call_back_enable", snapshot, nil)),
		d.Set("call_back_url", utils.PathSearch("call_back_url", snapshot, nil)),
		d.Set("call_back_auth_key", utils.PathSearch("auth_key", snapshot, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSnapshotQueryParams(domainName, appName string) string {
	return fmt.Sprintf("?domain=%s&app_name=%s", domainName, appName)
}

func resourceSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSnapshotHasChanges := []string{
		"frequency",
		"storage_mode",
		"storage_bucket",
		"storage_path",
		"call_back_enabled",
		"call_back_url",
		"call_back_auth_key",
	}

	if d.HasChanges(updateSnapshotHasChanges...) {
		var (
			httpUrl = "v1/{project_id}/stream/snapshot"
			product = "live"
		)
		client, err := cfg.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating Live client: %s", err)
		}

		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateLiveSnapshotBodyParams(d, region)),
		}
		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating Live snapshot: %s", err)
		}
	}
	return resourceSnapshotRead(ctx, d, meta)
}

func buildUpdateLiveSnapshotBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":            utils.ValueIgnoreEmpty(d.Get("domain_name")),
		"app_name":          utils.ValueIgnoreEmpty(d.Get("app_name")),
		"time_interval":     utils.ValueIgnoreEmpty(d.Get("frequency")),
		"object_write_mode": utils.ValueIgnoreEmpty(d.Get("storage_mode")),
		"obs_location":      buildLiveSnapshotObsLocationChildBody(d, region),
		"call_back_enable":  utils.ValueIgnoreEmpty(d.Get("call_back_enabled")),
		"call_back_url":     utils.ValueIgnoreEmpty(d.Get("call_back_url")),
		"auth_key":          utils.ValueIgnoreEmpty(d.Get("call_back_auth_key")),
	}
	return bodyParams
}

func resourceSnapshotDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/stream/snapshot"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDeleteSnapshotParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting Live snapshot: %s", err)
	}

	return nil
}

func buildDeleteSnapshotParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?domain=%v&app_name=%v", d.Get("domain_name"), d.Get("app_name"))
}
