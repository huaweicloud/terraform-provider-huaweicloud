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

var gaussDBInstanceLtsLogAssociateNonUpdatableParams = []string{"instance_id", "log_type"}

// @API GaussDB POST /v3/{project_id}/instances/logs/lts-config
// @API GaussDB GET /v3/{project_id}/instances/logs/lts-config
// @API GaussDB DELETE /v3/{project_id}/instances/logs/lts-config
func ResourceGaussdbInstanceLtsLogAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussdbInstanceLtsLogAssociateCreate,
		ReadContext:   resourceGaussdbInstanceLtsLogAssociateRead,
		UpdateContext: resourceGaussdbInstanceLtsLogAssociateUpdate,
		DeleteContext: resourceGaussdbInstanceLtsLogAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGaussdbInstanceLtsLogAssociateImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(gaussDBInstanceLtsLogAssociateNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
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

func resourceGaussdbInstanceLtsLogAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/logs/lts-config"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussdbInstanceLtsLogAssociateBodyParams(d))

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
		return diag.Errorf("error creating GaussDB instance LTS log associate: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)
	resourceId := fmt.Sprintf("%s/%s", instanceId, logType)
	d.SetId(resourceId)

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating GaussDB instance LTS log associate, job_id is not found in the response")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 2, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGaussdbInstanceLtsLogAssociateRead(ctx, d, meta)
}

func resourceGaussdbInstanceLtsLogAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/logs/lts-config"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetGaussdbInstanceLtsLogAssociateQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB instance LTS log associate")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)
	listPath := fmt.Sprintf("instance_lts_configs[?instance.id=='%s']|[0].lts_configs", instanceId)
	ltsList := utils.PathSearch(listPath, getRespBody, []interface{}{}).([]interface{})

	filterPath := fmt.Sprintf("[?log_type=='%s']|[0]", logType)
	matchedConfig := utils.PathSearch(filterPath, ltsList, nil)
	if matchedConfig == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB instance LTS log associate")
	}

	enabled := utils.PathSearch("enabled", matchedConfig, false).(bool)
	if !enabled {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "GaussDB instance LTS log associate is disabled")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("log_type", utils.PathSearch("log_type", matchedConfig, nil)),
		d.Set("lts_group_id", utils.PathSearch("lts_group_id", matchedConfig, nil)),
		d.Set("lts_stream_id", utils.PathSearch("lts_stream_id", matchedConfig, nil)),
		d.Set("enabled", utils.PathSearch("enabled", matchedConfig, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussdbInstanceLtsLogAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/logs/lts-config"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	updateOpt.JSONBody = utils.RemoveNil(buildCreateGaussdbInstanceLtsLogAssociateBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error updating GaussDB instance LTS log associate: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error updating GaussDB instance LTS log associate, job_id is not found in the response")
	}

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 2, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGaussdbInstanceLtsLogAssociateRead(ctx, d, meta)
}

func resourceGaussdbInstanceLtsLogAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/logs/lts-config"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteGaussdbInstanceLtsLogAssociateBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return deleteResp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200823"),
			"error deleting GaussDB instance LTS log associate")
	}

	deleteRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting GaussDB instance LTS log associate, job_id is not found in the response")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 2, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGaussdbInstanceLtsLogAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<log_type>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("log_type", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func buildCreateGaussdbInstanceLtsLogAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"list": []map[string]interface{}{
			{
				"instance_id":   d.Get("instance_id"),
				"log_type":      d.Get("log_type"),
				"lts_group_id":  d.Get("lts_group_id"),
				"lts_stream_id": d.Get("lts_stream_id"),
			},
		},
	}
}

func buildDeleteGaussdbInstanceLtsLogAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"list": []map[string]interface{}{
			{
				"instance_id":   d.Get("instance_id"),
				"log_type":      d.Get("log_type"),
				"lts_group_id":  d.Get("lts_group_id"),
				"lts_stream_id": d.Get("lts_stream_id"),
			},
		},
	}
}

func buildGetGaussdbInstanceLtsLogAssociateQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&instance_id=%v", res, d.Get("instance_id"))
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
