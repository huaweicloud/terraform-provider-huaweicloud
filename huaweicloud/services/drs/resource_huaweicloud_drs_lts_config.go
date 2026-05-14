package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS PUT /v3/{project_id}/jobs/{job_id}/lts-log-switch
// @API DRS GET /v3/{project_id}/jobs/{job_id}/lts-log-switch
func ResourceDrsLtsConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsConfigCreate,
		ReadContext:   resourceLtsConfigRead,
		UpdateContext: resourceLtsConfigUpdate,
		DeleteContext: resourceLtsConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"job_id"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The API supports editing `log_group_id` and `log_stream_id` to be empty, so Computed is not needed.
			"log_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func updateLtsConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"log_group_id":  d.Get("log_group_id"),
		"log_stream_id": d.Get("log_stream_id"),
		"lts_enabled":   true,
	}

	return map[string]interface{}{
		"job": bodyParam,
	}
}

func updateLtsConfig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/{project_id}/jobs/{job_id}/lts-log-switch"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         updateLtsConfigBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceLtsConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		jobId  = d.Get("job_id").(string)
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	if err := updateLtsConfig(client, d); err != nil {
		return diag.Errorf("error configing LTS to DRS in create operation: %s", err)
	}

	d.SetId(jobId)

	return resourceLtsConfigRead(ctx, d, meta)
}

func GetLtsConfig(client *golangsdk.ServiceClient, jobId string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{project_id}/jobs/{job_id}/lts-log-switch"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	ltsEnabled := utils.PathSearch("job.lts_enabled", respBody, false).(bool)
	if !ltsEnabled {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func resourceLtsConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		jobId  = d.Id()
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	respBody, err := GetLtsConfig(client, jobId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DRS LTS config")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("job_id", d.Id()),
		d.Set("log_group_id", utils.PathSearch("job.log_group_id", respBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("job.log_stream_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLtsConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	if err := updateLtsConfig(client, d); err != nil {
		return diag.Errorf("error configing LTS to DRS in update operation: %s", err)
	}

	return resourceLtsConfigRead(ctx, d, meta)
}

func deleteLtsConfigBodyParams() map[string]interface{} {
	bodyParam := map[string]interface{}{
		"lts_enabled": false,
	}

	return map[string]interface{}{
		"job": bodyParam,
	}
}

func resourceLtsConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/jobs/{job_id}/lts-log-switch"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         deleteLtsConfigBodyParams(),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DRS LTS config: %s", err)
	}

	return nil
}
