package gaussdb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gaussdbPluginNonUpdatableParams = []string{"instance_id", "plugin_name", "url", "sha_256"}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/kernel-plugin
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/plugins
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceGaussDbInstancePlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginCreate,
		ReadContext:   resourcePluginRead,
		UpdateContext: resourcePluginUpdate,
		DeleteContext: resourcePluginDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePluginImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(gaussdbPluginNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plugin_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sha_256": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plugin_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildPluginInstallBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"plugin_name": d.Get("plugin_name"),
		"url":         d.Get("url"),
		"sha_256":     d.Get("sha_256"),
	}
}

func resourcePluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/kernel-plugin"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildPluginInstallBodyParams(d)),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error installing GaussDB instance plugin: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error installing GaussDB instance plugin, jobId is not found in the response")
	}

	pluginName := d.Get("plugin_name").(string)
	d.SetId(fmt.Sprintf("%s/%s", instanceID, pluginName))

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 2, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePluginRead(ctx, d, meta)
}

func resourcePluginRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/plugins?plugin_name={plugin_name}"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	pluginName := d.Get("plugin_name").(string)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getPath = strings.ReplaceAll(getPath, "{plugin_name}", pluginName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB kernel plugins")
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	plugin := utils.PathSearch(fmt.Sprintf("plugins[?plugin_name=='%s' && installed] | [0]", pluginName),
		respBody, nil)
	if plugin == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
			fmt.Sprintf("GaussDB kernel plugin %s not found or not installed on instance %s", pluginName,
				instanceID))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("plugin_name", utils.PathSearch("plugin_name", plugin, nil)),
		d.Set("installed", utils.PathSearch("installed", plugin, nil)),
		d.Set("port", utils.PathSearch("port", plugin, nil)),
		d.Set("plugin_version", utils.PathSearch("plugin_version", plugin, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePluginUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePluginDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB instance plugin resource is not supported. The plugin resource is only removed " +
		"from the state, the GaussDB instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourcePluginImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, want '<instance_id>/<plugin_name>'")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("plugin_name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
