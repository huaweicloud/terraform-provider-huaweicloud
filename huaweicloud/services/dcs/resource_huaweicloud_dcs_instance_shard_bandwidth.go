package dcs

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

var instanceShardBandwidthNonUpdatableParams = []string{"instance_id", "group_id"}

// @API DCS PUT /v2/{project_id}/instances/{instance_id}/bandwidths
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS GET /v2/{project_id}/jobs/{job_id}
// @API DCS GET /v2/{project_id}/instances/{instance_id}/bandwidths
func ResourceDcsInstanceShardBandwidth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsInstanceShardBandwidthCreate,
		ReadContext:   resourceDcsInstanceShardBandwidthRead,
		UpdateContext: resourceDcsInstanceShardBandwidthUpdate,
		DeleteContext: resourceDcsInstanceShardBandwidthDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDcsInstanceShardBandwidthImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(instanceShardBandwidthNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"assured_bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceDcsInstanceShardBandwidthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	createRespBody, err := updateShardBandwidth(ctx, d, client, schema.TimeoutCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("instance_id").(string), d.Get("group_id").(string)))

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating instance(%s) shard bandwidth: job_id is not found in API response",
			d.Get("instance_id").(string))
	}

	err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDcsInstanceShardBandwidthRead(ctx, d, meta)
}

func resourceDcsInstanceShardBandwidthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/bandwidths"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting DCS instance shard bandwidth")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	searchExpression := fmt.Sprintf("group_bandwidths[?group_id=='%s']|[0]", d.Get("group_id").(string))
	groupBandwidth := utils.PathSearch(searchExpression, getRespBody, nil)
	if groupBandwidth == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting DCS instance shard bandwidth")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("group_id", utils.PathSearch("group_id", groupBandwidth, nil)),
		d.Set("bandwidth", utils.PathSearch("bandwidth", groupBandwidth, nil)),
		d.Set("max_bandwidth", utils.PathSearch("max_bandwidth", groupBandwidth, nil)),
		d.Set("assured_bandwidth", utils.PathSearch("assured_bandwidth", groupBandwidth, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", groupBandwidth, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcsInstanceShardBandwidthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	updateRespBody, err := updateShardBandwidth(ctx, d, client, schema.TimeoutUpdate)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating instance(%s) shard bandwidth: job_id is not found in API response",
			d.Get("instance_id").(string))
	}

	err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDcsInstanceShardBandwidthRead(ctx, d, meta)
}

func updateShardBandwidth(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	timeout string) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/bandwidths"
	instanceId := d.Get("instance_id").(string)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)

	updateOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	updateOpt.JSONBody = utils.RemoveNil(buildInstanceShardBandwidthBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		r, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleOperationError(err)
		return r, retry, err
	}
	updateResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceId),
		WaitTarget:   []string{"RUNNING"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(updateResp.(*http.Response))
}

func buildInstanceShardBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"group_id":         d.Get("group_id"),
		"target_bandwidth": d.Get("bandwidth"),
	}
	bodyParams := map[string]interface{}{
		"group_bandwidths": []interface{}{params},
	}
	return bodyParams
}

func resourceDcsInstanceShardBandwidthDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS instance shard bandwidth resource is not supported. The resource is only removed from the" +
		"state, but it remains in the cloud. And the instance doesn't return to the state before restoration."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceDcsInstanceShardBandwidthImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import id, must be <instance_id>/<group_id>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("group_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
