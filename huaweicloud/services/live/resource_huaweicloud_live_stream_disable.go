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

func ResourceLiveStreamDisable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLiveStreamDisableCreate,
		UpdateContext: resourceLiveStreamDisableUpdate,
		ReadContext:   resourceLiveStreamDisableRead,
		DeleteContext: resourceLiveStreamDisableDelete,
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
			"stream_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the stream name(not *).`,
			},
			"resume_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the time to resume stream push.`,
			},
		},
	}
}

func resourceLiveStreamDisableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createLiveStreamDisable: create Live stream disable
	var (
		createLiveStreamDisableHttpUrl = "v1/{project_id}/stream/blocks"
		createLiveStreamDisableProduct = "live"
	)
	createLiveStreamDisableClient, err := cfg.NewServiceClient(createLiveStreamDisableProduct, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	createLiveStreamDisablePath := createLiveStreamDisableClient.Endpoint + createLiveStreamDisableHttpUrl
	createLiveStreamDisablePath = strings.ReplaceAll(createLiveStreamDisablePath, "{project_id}",
		createLiveStreamDisableClient.ProjectID)

	createLiveStreamDisableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	createLiveStreamDisableOpt.JSONBody = utils.RemoveNil(buildLiveStreamDisableBodyParams(d))
	_, err = createLiveStreamDisableClient.Request("POST", createLiveStreamDisablePath, &createLiveStreamDisableOpt)
	if err != nil {
		return diag.Errorf("error creating Live stream disable: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	appName := d.Get("app_name").(string)
	streamName := d.Get("stream_name").(string)

	d.SetId(domainName + "/" + appName + "/" + streamName)

	return resourceLiveStreamDisableRead(ctx, d, meta)
}

func buildLiveStreamDisableBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":      utils.ValueIngoreEmpty(d.Get("domain_name")),
		"app_name":    utils.ValueIngoreEmpty(d.Get("app_name")),
		"stream_name": utils.ValueIngoreEmpty(d.Get("stream_name")),
		"resume_time": utils.ValueIngoreEmpty(d.Get("resume_time")),
	}
	return bodyParams
}

func resourceLiveStreamDisableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getLiveStreamDisable: Query Live stream disable
	var (
		getLiveStreamDisableHttpUrl = "v1/{project_id}/stream/blocks"
		getLiveStreamDisableProduct = "live"
	)
	getLiveStreamDisableClient, err := cfg.NewServiceClient(getLiveStreamDisableProduct, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return diag.Errorf("invalid id format, must be <domain_name>/<app_name>/<stream_name>")
	}
	domainName := parts[0]
	appName := parts[1]
	streamName := parts[2]

	getLiveStreamDisablePath := getLiveStreamDisableClient.Endpoint + getLiveStreamDisableHttpUrl
	getLiveStreamDisablePath = strings.ReplaceAll(getLiveStreamDisablePath, "{project_id}",
		getLiveStreamDisableClient.ProjectID)

	getLiveStreamDisableQueryParams := buildGetLiveStreamDisableQueryParams(domainName, appName, streamName)
	getLiveStreamDisablePath += getLiveStreamDisableQueryParams

	getLiveStreamDisableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getLiveStreamDisableResp, err := getLiveStreamDisableClient.Request("GET", getLiveStreamDisablePath,
		&getLiveStreamDisableOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live stream disable")
	}

	getLiveStreamDisableRespBody, err := utils.FlattenResponse(getLiveStreamDisableResp)
	if err != nil {
		return diag.FromErr(err)
	}

	blocks := utils.PathSearch("blocks", getLiveStreamDisableRespBody, make([]interface{}, 0)).([]interface{})
	if len(blocks) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("domain_name", domainName),
		d.Set("app_name", utils.PathSearch("app_name", blocks[0], nil)),
		d.Set("stream_name", utils.PathSearch("stream_name", blocks[0], nil)),
		d.Set("resume_time", utils.PathSearch("resume_time", blocks[0], nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetLiveStreamDisableQueryParams(domainName, appName, streamName string) string {
	return fmt.Sprintf("?domain=%s&app_name=%s&stream_name=%s", domainName, appName, streamName)
}

func resourceLiveStreamDisableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateLiveStreamDisableHasChanges := []string{
		"resume_time",
	}

	if d.HasChanges(updateLiveStreamDisableHasChanges...) {
		// updateLiveStreamDisable: update Live stream disable
		var (
			updateLiveStreamDisableHttpUrl = "v1/{project_id}/stream/blocks"
			updateLiveStreamDisableProduct = "live"
		)
		updateLiveStreamDisableClient, err := cfg.NewServiceClient(updateLiveStreamDisableProduct, region)
		if err != nil {
			return diag.Errorf("error creating Live Client: %s", err)
		}

		updateLiveStreamDisablePath := updateLiveStreamDisableClient.Endpoint + updateLiveStreamDisableHttpUrl
		updateLiveStreamDisablePath = strings.ReplaceAll(updateLiveStreamDisablePath, "{project_id}",
			updateLiveStreamDisableClient.ProjectID)

		updateLiveStreamDisableOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateLiveStreamDisableOpt.JSONBody = utils.RemoveNil(buildLiveStreamDisableBodyParams(d))
		_, err = updateLiveStreamDisableClient.Request("PUT", updateLiveStreamDisablePath, &updateLiveStreamDisableOpt)
		if err != nil {
			return diag.Errorf("error updating Live stream disable: %s", err)
		}
	}
	return resourceLiveStreamDisableRead(ctx, d, meta)
}

func resourceLiveStreamDisableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteLiveStreamDisable: Delete Live stream disable
	var (
		deleteLiveStreamDisableHttpUrl = "v1/{project_id}/stream/blocks"
		deleteLiveStreamDisableProduct = "live"
	)
	deleteLiveStreamDisableClient, err := cfg.NewServiceClient(deleteLiveStreamDisableProduct, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	deleteLiveStreamDisablePath := deleteLiveStreamDisableClient.Endpoint + deleteLiveStreamDisableHttpUrl
	deleteLiveStreamDisablePath = strings.ReplaceAll(deleteLiveStreamDisablePath, "{project_id}",
		deleteLiveStreamDisableClient.ProjectID)

	deleteLiveStreamDisableParams := buildDeleteLiveStreamDisableParams(d)
	deleteLiveStreamDisablePath += deleteLiveStreamDisableParams

	deleteLiveStreamDisableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteLiveStreamDisableClient.Request("DELETE", deleteLiveStreamDisablePath,
		&deleteLiveStreamDisableOpt)
	if err != nil {
		return diag.Errorf("error deleting Live stream disable: %s", err)
	}

	return nil
}

func buildDeleteLiveStreamDisableParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("domain_name"); ok {
		res = fmt.Sprintf("%s&domain=%v", res, v)
	}

	if v, ok := d.GetOk("app_name"); ok {
		res = fmt.Sprintf("%s&app_name=%v", res, v)
	}

	if v, ok := d.GetOk("stream_name"); ok {
		res = fmt.Sprintf("%s&stream_name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
