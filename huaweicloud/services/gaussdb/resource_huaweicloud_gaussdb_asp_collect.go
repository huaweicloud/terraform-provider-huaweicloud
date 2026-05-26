package gaussdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gaussDbAspCollectNonUpdatableParams = []string{
	"instance_id",
	"start_time",
	"end_time",
}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/asp/collect
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/asp
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceGaussDbAspCollect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDbAspCollectCreate,
		UpdateContext: resourceGaussDbAspCollectUpdate,
		ReadContext:   resourceGaussDbAspCollectRead,
		DeleteContext: resourceGaussDbAspCollectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGaussDbAspImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(gaussDbAspCollectNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

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
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"file_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"download_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"obs_bucket": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aspObsBucketSchema(),
			},
		},
	}
}

func aspObsBucketSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGaussDbAspCollectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/asp/collect"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateAspCollectBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB ASP collect: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId != "" {
		if err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 10, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error creating GaussDB ASP collect: %s", err)
		}
	}

	d.SetId(jobId)

	return resourceGaussDbAspCollectRead(ctx, d, meta)
}

func buildCreateAspCollectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time": d.Get("start_time"),
		"end_time":   d.Get("end_time"),
	}
	return bodyParams
}

func resourceGaussDbAspCollectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/asp"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB ASP collect")
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	asp := utils.PathSearch(fmt.Sprintf("asp[?job_id=='%s']|[0]", d.Id()), getRespBody, nil)
	if asp == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB ASP collect")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", d.Get("instance_id")),
		d.Set("start_time", utils.PathSearch("start_time", asp, nil)),
		d.Set("end_time", utils.PathSearch("end_time", asp, nil)),
		d.Set("file_size", utils.PathSearch("file_size", asp, nil)),
		d.Set("download_url", utils.PathSearch("download_url", asp, nil)),
		d.Set("status", utils.PathSearch("status", asp, nil)),
		d.Set("file_path", utils.PathSearch("file_path", asp, nil)),
		d.Set("file_name", utils.PathSearch("file_name", asp, nil)),
		d.Set("obs_bucket", flattenGaussDbAspCollectObsBucket(asp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGaussDbAspCollectObsBucket(resp interface{}) []interface{} {
	obsBucket := utils.PathSearch("obs_bucket", resp, nil)
	if obsBucket == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"name":      utils.PathSearch("name", obsBucket, nil),
			"type":      utils.PathSearch("type", obsBucket, nil),
			"url":       utils.PathSearch("url", obsBucket, nil),
			"port":      utils.PathSearch("port", obsBucket, nil),
			"domain_id": utils.PathSearch("domain_id", obsBucket, nil),
		},
	}
}

func resourceGaussDbAspCollectUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbAspCollectDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB ASP collect resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceGaussDbAspImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
