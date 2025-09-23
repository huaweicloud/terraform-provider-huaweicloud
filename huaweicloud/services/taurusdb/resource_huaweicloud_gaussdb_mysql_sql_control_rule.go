// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GaussDB
// ---------------------------------------------------------------

package taurusdb

import (
	"context"
	"fmt"
	"net/http"
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

// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/sql-filter/rules
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/sql-filter/rules
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/sql-filter/rules
func ResourceGaussDBSqlControlRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBSqlControlRuleCreate,
		UpdateContext: resourceGaussDBSqlControlRuleUpdate,
		ReadContext:   resourceGaussDBSqlControlRuleRead,
		DeleteContext: resourceGaussDBSqlControlRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the GaussDB MySQL instance.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies ID of the GaussDB MySQL node.`,
			},
			"sql_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies SQL statement type.`,
			},
			"pattern": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the concurrency control rule of SQL statements.`,
			},
			"max_concurrency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the maximum number of concurrent SQL statements.`,
			},
		},
	}
}

func resourceGaussDBSqlControlRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	// createGaussDBSqlControlRule: create a GaussDB MySQL Sql control rule
	err := dealGaussDBSqlControlRule(ctx, d, cfg, "creating", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)
	sqlType := d.Get("sql_type").(string)
	pattern := d.Get("pattern").(string)
	d.SetId(instanceID + "/" + nodeId + "/" + sqlType + "/" + pattern)

	return resourceGaussDBSqlControlRuleRead(ctx, d, meta)
}

func dealGaussDBSqlControlRule(ctx context.Context, d *schema.ResourceData, cfg *config.Config,
	operateMethod string, timeout time.Duration) error {
	region := cfg.GetRegion(d)
	var (
		gaussDBSqlControlRuleHttpUrl = "v3/{project_id}/instances/{instance_id}/sql-filter/rules"
		gaussDBSqlControlRuleProduct = "gaussdb"
	)
	gaussDBSqlControlRuleClient, err := cfg.NewServiceClient(gaussDBSqlControlRuleProduct, region)
	if err != nil {
		return fmt.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	gaussDBSqlControlRulePath := gaussDBSqlControlRuleClient.Endpoint + gaussDBSqlControlRuleHttpUrl
	gaussDBSqlControlRulePath = strings.ReplaceAll(gaussDBSqlControlRulePath, "{project_id}",
		gaussDBSqlControlRuleClient.ProjectID)
	gaussDBSqlControlRulePath = strings.ReplaceAll(gaussDBSqlControlRulePath, "{instance_id}", instanceID)

	gaussDBSqlControlRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	gaussDBSqlControlRuleOpt.JSONBody = utils.RemoveNil(buildGaussDBSqlControlRuleBodyParams(d))

	var gaussDBSqlControlRuleResp *http.Response
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		gaussDBSqlControlRuleResp, err = gaussDBSqlControlRuleClient.Request("PUT", gaussDBSqlControlRulePath,
			&gaussDBSqlControlRuleOpt)
		isRetry, err := handleGaussDBMysqlOperationError(err)
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error %s GaussDB MySQL SQL control rule: %s", operateMethod, err)
	}

	gaussDBSqlControlRuleRespBody, err := utils.FlattenResponse(gaussDBSqlControlRuleResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", gaussDBSqlControlRuleRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find the job ID of the GaussDB MySQL SQL control rule from the API response")
	}
	return waitForJobComplete(ctx, gaussDBSqlControlRuleClient, timeout, instanceID, jobId)
}

func buildGaussDBSqlControlRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	patternParams := map[string]interface{}{
		"pattern":         utils.ValueIgnoreEmpty(d.Get("pattern")),
		"max_concurrency": utils.ValueIgnoreEmpty(d.Get("max_concurrency")),
	}
	rulesParams := map[string]interface{}{
		"sql_type": utils.ValueIgnoreEmpty(d.Get("sql_type")),
		"patterns": []interface{}{patternParams},
	}
	nodeFilterRulesParams := map[string]interface{}{
		"node_id": utils.ValueIgnoreEmpty(d.Get("node_id")),
		"rules":   []interface{}{rulesParams},
	}
	bodyParams := map[string]interface{}{
		"sql_filter_rules": []interface{}{nodeFilterRulesParams},
	}
	return bodyParams
}

func resourceGaussDBSqlControlRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGaussDBSqlControlRule: Query the GaussDB MySQL SQL control rule
	var (
		getGaussDBSqlControlRuleHttpUrl = "v3/{project_id}/instances/{instance_id}/sql-filter/rules"
		getGaussDBSqlControlRuleProduct = "gaussdb"
	)
	getGaussDBSqlControlRuleClient, err := cfg.NewServiceClient(getGaussDBSqlControlRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 4)
	if len(parts) != 4 {
		return diag.Errorf("invalid id format, must be <instance_id>/<node_id>/<sql_type>/<pattern>")
	}
	instanceID := parts[0]
	nodeId := parts[1]
	sqlType := parts[2]
	pattern := parts[3]

	getGaussDBSqlControlRulePath := getGaussDBSqlControlRuleClient.Endpoint + getGaussDBSqlControlRuleHttpUrl
	getGaussDBSqlControlRulePath = strings.ReplaceAll(getGaussDBSqlControlRulePath, "{project_id}",
		getGaussDBSqlControlRuleClient.ProjectID)
	getGaussDBSqlControlRulePath = strings.ReplaceAll(getGaussDBSqlControlRulePath, "{instance_id}", instanceID)

	getGaussDBSqlControlRulePath += buildGetGaussDBSqlControlRuleQueryParams(nodeId)

	getGaussDBSqlControlRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getGaussDBSqlControlRuleResp, err := getGaussDBSqlControlRuleClient.Request("GET",
		getGaussDBSqlControlRulePath, &getGaussDBSqlControlRuleOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
			"error retrieving GaussDB MySQL SQL control rule")
	}

	getGaussDBSqlControlRuleRespBody, err := utils.FlattenResponse(getGaussDBSqlControlRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := fmt.Sprintf("sql_filter_rules[?sql_type=='%s']|[0].patterns[?pattern=='%s']|[0].max_concurrency",
		sqlType, pattern)
	maxConcurrency := utils.PathSearch(expression, getGaussDBSqlControlRuleRespBody, nil)
	if maxConcurrency == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("node_id", nodeId),
		d.Set("sql_type", sqlType),
		d.Set("pattern", pattern),
		d.Set("max_concurrency", maxConcurrency),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetGaussDBSqlControlRuleQueryParams(nodeId string) string {
	return fmt.Sprintf("?node_id=%v", nodeId)
}

func resourceGaussDBSqlControlRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	// updateGaussDBSqlControlRule: update the GaussDB MySQL Sql control rule
	err := dealGaussDBSqlControlRule(ctx, d, cfg, "updating", d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceGaussDBSqlControlRuleRead(ctx, d, meta)
}

func resourceGaussDBSqlControlRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGaussDBSqlControlRule: delete the GaussDB MySQL Sql control rule
	var (
		deleteGaussDBSqlControlRuleHttpUrl = "v3/{project_id}/instances/{instance_id}/sql-filter/rules"
		deleteGaussDBSqlControlRuleProduct = "gaussdb"
	)
	deleteGaussDBSqlControlRuleClient, err := cfg.NewServiceClient(deleteGaussDBSqlControlRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteGaussDBSqlControlRulePath := deleteGaussDBSqlControlRuleClient.Endpoint + deleteGaussDBSqlControlRuleHttpUrl
	deleteGaussDBSqlControlRulePath = strings.ReplaceAll(deleteGaussDBSqlControlRulePath, "{project_id}",
		deleteGaussDBSqlControlRuleClient.ProjectID)
	deleteGaussDBSqlControlRulePath = strings.ReplaceAll(deleteGaussDBSqlControlRulePath, "{instance_id}", instanceID)

	deleteGaussDBSqlControlRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteGaussDBSqlControlRuleOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDBSqlControlRuleBodyParams(d))

	var deleteGaussDBSqlControlRuleResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		deleteGaussDBSqlControlRuleResp, err = deleteGaussDBSqlControlRuleClient.Request("DELETE",
			deleteGaussDBSqlControlRulePath, &deleteGaussDBSqlControlRuleOpt)
		isRetry, err := handleGaussDBMysqlOperationError(err)
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error deleting GaussDB MySQL SQL control rule: %s", err)
	}

	deleteGaussDBSqlControlRulesRespBody, err := utils.FlattenResponse(deleteGaussDBSqlControlRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteGaussDBSqlControlRulesRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the GaussDB MySQL SQL control rule from the API response")
	}

	err = waitForJobComplete(ctx, deleteGaussDBSqlControlRuleClient, d.Timeout(schema.TimeoutDelete), instanceID, jobId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteGaussDBSqlControlRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	rulesParams := map[string]interface{}{
		"sql_type": utils.ValueIgnoreEmpty(d.Get("sql_type")),
		"patterns": []interface{}{utils.ValueIgnoreEmpty(d.Get("pattern"))},
	}
	nodeFilterRulesParams := map[string]interface{}{
		"node_id": utils.ValueIgnoreEmpty(d.Get("node_id")),
		"rules":   []interface{}{rulesParams},
	}
	bodyParams := map[string]interface{}{
		"sql_filter_rules": []interface{}{nodeFilterRulesParams},
	}
	return bodyParams
}
