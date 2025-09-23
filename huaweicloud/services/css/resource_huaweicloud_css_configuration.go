// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CSS
// ---------------------------------------------------------------

package css

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/ymls/update
// @API CSS GET /v1.0/{project_id}/clusters/{id}/ymls/joblists
// @API CSS GET /v1.0/{project_id}/clusters/{id}/ymls/template
// @API CSS POST /v1.0/{project_id}/clusters/{id}/ymls/update
func ResourceCssConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssConfigurationUpdate,
		UpdateContext: resourceCssConfigurationUpdate,
		ReadContext:   resourceCssConfigurationRead,
		DeleteContext: resourceCssConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The CSS cluster ID.`,
			},
			"http_cors_allow_credetials": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Whether to return the Access-Control-Allow-Credentials of the header during cross-domain access.`,
			},
			"http_cors_allow_origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Origin IP address allowed for cross-domain access, for example, **122.122.122.122:9200**.`,
			},
			"http_cors_max_age": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Cache duration of the browser. The cache is automatically cleared after the time range you specify.`,
			},
			"http_cors_allow_headers": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Headers allowed for cross-domain access.`,
			},
			"http_cors_enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Whether to allow cross-domain access.`,
			},
			"http_cors_allow_methods": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Methods allowed for cross-domain access.`,
			},
			"reindex_remote_whitelist": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Configured for migrating data from the current cluster to the target cluster through the reindex API.`,
			},
			"indices_queries_cache_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Cache size in the query phase. Value range: **1** to **100**.`,
			},
			"thread_pool_force_merge_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Queue size in the force merge thread pool.`,
			},
		},
	}
}

func resourceCssConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateConfigurationHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ymls/update"
		updateConfigurationProduct = "css"
	)
	updateConfigurationClient, err := cfg.NewServiceClient(updateConfigurationProduct, region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	updateConfigurationPath := updateConfigurationClient.Endpoint + updateConfigurationHttpUrl
	updateConfigurationPath = strings.ReplaceAll(updateConfigurationPath, "{project_id}", updateConfigurationClient.ProjectID)
	updateConfigurationPath = strings.ReplaceAll(updateConfigurationPath, "{cluster_id}", d.Get("cluster_id").(string))

	updateConfigurationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	updateConfigurationOpt.JSONBody = utils.RemoveNil(buildUpdateConfigurationBodyParams(d))
	updateConfigurationResp, err := updateConfigurationClient.Request("POST", updateConfigurationPath, &updateConfigurationOpt)
	if err != nil {
		return diag.Errorf("error creating CSS configuration: %s", err)
	}

	_, err = utils.FlattenResponse(updateConfigurationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("cluster_id").(string))

	err = configurationWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the CSS configuration (%s) update to complete: %s", d.Id(), err)
	}
	return resourceCssConfigurationRead(ctx, d, meta)
}

func buildUpdateConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"edit": map[string]interface{}{
			"modify": map[string]interface{}{
				"elasticsearch.yml": map[string]interface{}{
					"http.cors.allow-credentials":  utils.ValueIgnoreEmpty(d.Get("http_cors_allow_credetials")),
					"http.cors.allow-origin":       utils.ValueIgnoreEmpty(d.Get("http_cors_allow_origin")),
					"http.cors.max-age":            utils.ValueIgnoreEmpty(d.Get("http_cors_max_age")),
					"http.cors.allow-headers":      utils.ValueIgnoreEmpty(d.Get("http_cors_allow_headers")),
					"http.cors.enabled":            utils.ValueIgnoreEmpty(d.Get("http_cors_enabled")),
					"http.cors.allow-methods":      utils.ValueIgnoreEmpty(d.Get("http_cors_allow_methods")),
					"reindex.remote.whitelist":     utils.ValueIgnoreEmpty(d.Get("reindex_remote_whitelist")),
					"indices.queries.cache.size":   utils.ValueIgnoreEmpty(d.Get("indices_queries_cache_size")),
					"thread_pool.force_merge.size": utils.ValueIgnoreEmpty(d.Get("thread_pool_force_merge_size")),
				},
			},
		},
	}
	return bodyParams
}

func resourceCssConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getConfiguration: Query the CSS configuration.
	var (
		getConfigurationHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ymls/template"
		getConfigurationProduct = "css"
	)
	getConfigurationClient, err := cfg.NewServiceClient(getConfigurationProduct, region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	getConfigurationPath := getConfigurationClient.Endpoint + getConfigurationHttpUrl
	getConfigurationPath = strings.ReplaceAll(getConfigurationPath, "{project_id}", getConfigurationClient.ProjectID)
	getConfigurationPath = strings.ReplaceAll(getConfigurationPath, "{cluster_id}", d.Id())

	getConfigurationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getConfigurationResp, err := getConfigurationClient.Request("GET", getConfigurationPath, &getConfigurationOpt)
	if err != nil {
		// The cluster does not exist, http code is 403, key/value of error code is errCode/CSS.0015
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error retrieving CSS configuration")
	}

	getConfigurationRespBody, err := utils.FlattenResponse(getConfigurationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("cluster_id", d.Id()),
		d.Set("http_cors_allow_credetials", utils.PathSearch(`configurations."http.cors.allow-credentials".value`, getConfigurationRespBody, nil)),
		d.Set("http_cors_allow_origin", utils.PathSearch(`configurations."http.cors.allow-origin".value`, getConfigurationRespBody, nil)),
		d.Set("http_cors_max_age", utils.PathSearch(`configurations."http.cors.max-age".value`, getConfigurationRespBody, nil)),
		d.Set("http_cors_allow_headers", utils.PathSearch(`configurations."http.cors.allow-headers".value`, getConfigurationRespBody, nil)),
		d.Set("http_cors_enabled", utils.PathSearch(`configurations."http.cors.enabled".value`, getConfigurationRespBody, nil)),
		d.Set("http_cors_allow_methods", utils.PathSearch(`configurations."http.cors.allow-methods".value`, getConfigurationRespBody, nil)),
		d.Set("reindex_remote_whitelist", utils.PathSearch(`configurations."reindex.remote.whitelist".value`, getConfigurationRespBody, nil)),
		d.Set("indices_queries_cache_size", utils.PathSearch(`configurations."indices.queries.cache.size".value`, getConfigurationRespBody, nil)),
		d.Set("thread_pool_force_merge_size", utils.PathSearch(`configurations."thread_pool.force_merge.size".value`, getConfigurationRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCssConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteConfiguration: delete CSS configuration
	var (
		deleteConfigurationHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ymls/update"
		deleteConfigurationProduct = "css"
	)
	deleteConfigurationClient, err := cfg.NewServiceClient(deleteConfigurationProduct, region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	deleteConfigurationPath := deleteConfigurationClient.Endpoint + deleteConfigurationHttpUrl
	deleteConfigurationPath = strings.ReplaceAll(deleteConfigurationPath, "{project_id}", deleteConfigurationClient.ProjectID)
	deleteConfigurationPath = strings.ReplaceAll(deleteConfigurationPath, "{cluster_id}", d.Id())

	deleteConfigurationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	deleteConfigurationOpt.JSONBody = buildDeleteConfigurationBodyParams()
	_, err = deleteConfigurationClient.Request("POST", deleteConfigurationPath, &deleteConfigurationOpt)
	if err != nil {
		// The cluster does not exist, http code is 403, key/value of error code is errCode/CSS.0015
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error deleting CSS configuration")
	}

	err = configurationWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the CSS configuration (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

// Reset to default value.
func buildDeleteConfigurationBodyParams() map[string]interface{} {
	bodyParams := map[string]interface{}{
		"edit": map[string]interface{}{
			"reset": map[string]interface{}{
				"elasticsearch.yml": map[string]interface{}{
					"http.cors.allow-credentials":  "",
					"http.cors.allow-origin":       "",
					"http.cors.max-age":            "",
					"http.cors.allow-headers":      "",
					"http.cors.enabled":            "",
					"http.cors.allow-methods":      "",
					"reindex.remote.whitelist":     "",
					"indices.queries.cache.size":   "",
					"thread_pool.force_merge.size": "",
				},
			},
		},
	}
	return bodyParams
}

func configurationWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				configurationWaitingHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ymls/joblists"
				configurationWaitingProduct = "css"
			)
			client, err := cfg.NewServiceClient(configurationWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CSS client: %s", err)
			}

			configWaitingPath := client.Endpoint + configurationWaitingHttpUrl
			configWaitingPath = strings.ReplaceAll(configWaitingPath, "{project_id}", client.ProjectID)
			configWaitingPath = strings.ReplaceAll(configWaitingPath, "{cluster_id}", d.Id())

			requestOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			configurationWaitingResp, err := client.Request("GET", configWaitingPath, &requestOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return configurationWaitingResp, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(configurationWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}

			// If the status of a job is running, it is considered not completed yet.
			actionProgressRaw := utils.PathSearch(`length(configList[?status == 'running'])`, respBody, 0.0)
			if actionProgressRaw.(float64) > 0 {
				return respBody, "PENDING", nil
			}

			return respBody, "COMPLETED", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
