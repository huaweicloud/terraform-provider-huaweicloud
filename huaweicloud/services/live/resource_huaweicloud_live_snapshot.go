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
func ResourceLiveSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLiveSnapshotCreate,
		UpdateContext: resourceLiveSnapshotUpdate,
		ReadContext:   resourceLiveSnapshotRead,
		DeleteContext: resourceLiveSnapshotDelete,
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

func resourceLiveSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createLiveSnapshot: create Live snapshot
	var (
		createLiveSnapshotHttpUrl = "v1/{project_id}/stream/snapshot"
		createLiveSnapshotProduct = "live"
	)
	createLiveSnapshotClient, err := cfg.NewServiceClient(createLiveSnapshotProduct, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	createLiveSnapshotPath := createLiveSnapshotClient.Endpoint + createLiveSnapshotHttpUrl
	createLiveSnapshotPath = strings.ReplaceAll(createLiveSnapshotPath, "{project_id}",
		createLiveSnapshotClient.ProjectID)

	createLiveSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createLiveSnapshotOpt.JSONBody = utils.RemoveNil(buildCreateLiveSnapshotBodyParams(d, region))
	_, err = createLiveSnapshotClient.Request("POST", createLiveSnapshotPath, &createLiveSnapshotOpt)
	if err != nil {
		return diag.Errorf("error creating Live snapshot: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	appName := d.Get("app_name").(string)

	d.SetId(domainName + "/" + appName)

	return resourceLiveSnapshotRead(ctx, d, meta)
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

func resourceLiveSnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getLiveSnapshot: Query Live snapshot
	var (
		getLiveSnapshotHttpUrl = "v1/{project_id}/stream/snapshot"
		getLiveSnapshotProduct = "live"
	)
	getLiveSnapshotClient, err := cfg.NewServiceClient(getLiveSnapshotProduct, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <domain_name>/<app_name>")
	}
	domainName := parts[0]
	appName := parts[1]

	getLiveSnapshotPath := getLiveSnapshotClient.Endpoint + getLiveSnapshotHttpUrl
	getLiveSnapshotPath = strings.ReplaceAll(getLiveSnapshotPath, "{project_id}", getLiveSnapshotClient.ProjectID)

	getLiveSnapshotQueryParams := buildGetLiveSnapshotQueryParams(domainName, appName)
	getLiveSnapshotPath += getLiveSnapshotQueryParams

	getLiveSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getLiveSnapshotResp, err := getLiveSnapshotClient.Request("GET", getLiveSnapshotPath, &getLiveSnapshotOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live snapshot")
	}

	getLiveSnapshotRespBody, err := utils.FlattenResponse(getLiveSnapshotResp)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshot := utils.PathSearch("snapshot_config_list|[0]", getLiveSnapshotRespBody, nil)
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

func buildGetLiveSnapshotQueryParams(domainName, appName string) string {
	return fmt.Sprintf("?domain=%s&app_name=%s", domainName, appName)
}

func resourceLiveSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateLiveSnapshotHasChanges := []string{
		"frequency",
		"storage_mode",
		"storage_bucket",
		"storage_path",
		"call_back_enabled",
		"call_back_url",
		"call_back_auth_key",
	}

	if d.HasChanges(updateLiveSnapshotHasChanges...) {
		// updateLiveSnapshot: update Live snapshot
		var (
			updateLiveSnapshotHttpUrl = "v1/{project_id}/stream/snapshot"
			updateLiveSnapshotProduct = "live"
		)
		updateLiveSnapshotClient, err := cfg.NewServiceClient(updateLiveSnapshotProduct, region)
		if err != nil {
			return diag.Errorf("error creating Live Client: %s", err)
		}

		updateLiveSnapshotPath := updateLiveSnapshotClient.Endpoint + updateLiveSnapshotHttpUrl
		updateLiveSnapshotPath = strings.ReplaceAll(updateLiveSnapshotPath, "{project_id}",
			updateLiveSnapshotClient.ProjectID)

		updateLiveSnapshotOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateLiveSnapshotOpt.JSONBody = utils.RemoveNil(buildUpdateLiveSnapshotBodyParams(d, region))
		_, err = updateLiveSnapshotClient.Request("PUT", updateLiveSnapshotPath, &updateLiveSnapshotOpt)
		if err != nil {
			return diag.Errorf("error updating Live snapshot: %s", err)
		}
	}
	return resourceLiveSnapshotRead(ctx, d, meta)
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

func resourceLiveSnapshotDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteLiveSnapshot: Delete Live snapshot
	var (
		deleteLiveSnapshotHttpUrl = "v1/{project_id}/stream/snapshot"
		deleteLiveSnapshotProduct = "live"
	)
	deleteLiveSnapshotClient, err := cfg.NewServiceClient(deleteLiveSnapshotProduct, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	deleteLiveSnapshotPath := deleteLiveSnapshotClient.Endpoint + deleteLiveSnapshotHttpUrl
	deleteLiveSnapshotPath = strings.ReplaceAll(deleteLiveSnapshotPath, "{project_id}",
		deleteLiveSnapshotClient.ProjectID)

	deleteLiveSnapshotParams := buildDeleteLiveSnapshotParams(d)
	deleteLiveSnapshotPath += deleteLiveSnapshotParams

	deleteLiveSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteLiveSnapshotClient.Request("DELETE", deleteLiveSnapshotPath, &deleteLiveSnapshotOpt)
	if err != nil {
		return diag.Errorf("error deleting Live snapshot: %s", err)
	}

	return nil
}

func buildDeleteLiveSnapshotParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?domain=%v&app_name=%v", d.Get("domain_name"), d.Get("app_name"))
}
