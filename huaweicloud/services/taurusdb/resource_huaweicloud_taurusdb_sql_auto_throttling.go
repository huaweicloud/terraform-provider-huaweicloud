package taurusdb

import (
	"context"
	"errors"
	"fmt"
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

var autoThrottlingRuleNoneUpdatableParams = []string{"instance_id", "node_id"}

// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/auto-sql-limiting
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/auto-sql-limiting
// @API TaurusDB DELETE /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/auto-sql-limiting
// @API TaurusDB GET /v3/{project_id}/jobs
func ResourceTaurusDBSqlAutoThrottling() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBSqlAutoThrottlingCreate,
		ReadContext:   resourceTaurusDBSqlAutoThrottlingRead,
		UpdateContext: resourceTaurusDBSqlAutoThrottlingUpdate,
		DeleteContext: resourceTaurusDBSqlAutoThrottlingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTaurusDBSqlAutoThrottlingImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(autoThrottlingRuleNoneUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"node_id": {
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
			"condition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu_usage": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"active_sessions": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"clear_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_concurrency": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"retain_sql_rule": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
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

func resourceTaurusDBSqlAutoThrottlingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/auto-sql-limiting"
		product = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateSqlAutoThrottlingBodyParams(d))

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating TaurusDB SQL auto throttling: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)
	d.SetId(fmt.Sprintf("%s/%s", instanceId, nodeId))

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for TaurusDB SQL auto throttling job (%s) to complete: %s", jobId, err)
	}

	return resourceTaurusDBSqlAutoThrottlingRead(ctx, d, meta)
}

func buildCreateSqlAutoThrottlingBodyParams(d *schema.ResourceData) map[string]interface{} {
	rule := map[string]interface{}{
		"node_id":         d.Get("node_id"),
		"start_time":      d.Get("start_time"),
		"end_time":        d.Get("end_time"),
		"condition":       d.Get("condition"),
		"cpu_usage":       d.Get("cpu_usage"),
		"active_sessions": d.Get("active_sessions"),
		"clear_time":      d.Get("clear_time"),
		"duration":        d.Get("duration"),
		"max_concurrency": d.Get("max_concurrency"),
		"retain_sql_rule": d.Get("retain_sql_rule"),
	}

	bodyParams := map[string]interface{}{
		"auto_sql_limiting_rules": []map[string]interface{}{rule},
	}
	return bodyParams
}

func resourceTaurusDBSqlAutoThrottlingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/auto-sql-limiting"
		product = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getOpt.JSONBody = map[string]interface{}{
		"node_ids": []string{nodeId},
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving TaurusDB SQL auto throttling")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	rules := utils.PathSearch("auto_sql_limiting_rules", getRespBody, make([]interface{}, 0)).([]interface{})
	if len(rules) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving TaurusDB SQL auto throttling")
	}

	ruleKey := fmt.Sprintf("auto_sql_limiting_rules[?node_id=='%s']", nodeId)
	filteredRules := utils.PathSearch(ruleKey, getRespBody, make([]interface{}, 0)).([]interface{})
	if len(filteredRules) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving TaurusDB SQL auto throttling")
	}

	ruleMap := filteredRules[0].(map[string]interface{})

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", getRespBody, nil)),
		d.Set("node_id", ruleMap["node_id"]),
		d.Set("start_time", ruleMap["start_time"]),
		d.Set("end_time", ruleMap["end_time"]),
		d.Set("condition", ruleMap["condition"]),
		d.Set("cpu_usage", ruleMap["cpu_usage"]),
		d.Set("active_sessions", ruleMap["active_sessions"]),
		d.Set("clear_time", ruleMap["clear_time"]),
		d.Set("duration", ruleMap["duration"]),
		d.Set("max_concurrency", ruleMap["max_concurrency"]),
		d.Set("retain_sql_rule", d.Get("retain_sql_rule")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTaurusDBSqlAutoThrottlingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/auto-sql-limiting"
		product = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateOpt.JSONBody = utils.RemoveNil(buildCreateSqlAutoThrottlingBodyParams(d))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating TaurusDB SQL auto throttling: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for TaurusDB SQL auto throttling job (%s) to complete: %s", jobId, err)
	}

	return resourceTaurusDBSqlAutoThrottlingRead(ctx, d, meta)
}

func resourceTaurusDBSqlAutoThrottlingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/auto-sql-limiting"
		product = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{node_id}", nodeId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting TaurusDB SQL auto throttling: %s", err)
	}
	return nil
}

func resourceTaurusDBSqlAutoThrottlingImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<node_id>")
	}
	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("node_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
