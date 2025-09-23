// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var sqlAuditNonUpdatableParams = []string{"instance_id"}

// @API RDS GET /v3/{project_id}/instances/{instance_id}/auditlog-policy
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/auditlog-policy
func ResourceSQLAudit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSQLAuditCreate,
		UpdateContext: resourceSQLAuditUpdate,
		ReadContext:   resourceSQLAuditRead,
		DeleteContext: resourceSQLAuditDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(sqlAuditNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"keep_days": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `Specifies the number of days for storing audit logs.`,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"audit_types": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the list of audit types.`,
			},
			"reserve_auditlogs": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `Specifies whether the historical audit logs will be reserved for some time when SQL
audit is disabled.`,
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

func resourceSQLAuditCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSQLAudit: create RDS SQL audit
	var (
		createSQLAuditProduct = "rds"
	)
	createSQLAuditClient, err := cfg.NewServiceClient(createSQLAuditProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	params := buildUpdateSQLAuditBodyParams(d)
	log.Printf("[DEBUG] Create RDS SQL audit params: %#v", params)
	err = updateSQLAudit(ctx, d, createSQLAuditClient, d.Timeout(schema.TimeoutCreate), params)
	if err != nil {
		return diag.Errorf("error creating RDS SQL audit: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	d.SetId(instanceID)

	stateConf := &resource.StateChangeConf{
		Target:       []string{"COMPLETED"},
		Refresh:      rdsSQLAuditStateRefreshFunc(createSQLAuditClient, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RDS instance(%s) SQL audit creation completed: %s", instanceID, err)
	}

	return resourceSQLAuditRead(ctx, d, meta)
}

func resourceSQLAuditRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSQLAudit: Query the RDS SQL audit
	var (
		getSQLAuditProduct = "rds"
	)
	getSQLAuditClient, err := cfg.NewServiceClient(getSQLAuditProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instance, err := GetRdsInstanceByID(getSQLAuditClient, d.Id())
	if err != nil {
		return diag.Errorf("error getting RDS instance: %s", err)
	}
	instanceId := utils.PathSearch("id", instance, "").(string)
	if instanceId == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	getSQLAuditRespBody, err := getSQLAudit(getSQLAuditClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS SQL audit")
	}

	keepDays := utils.PathSearch("keep_days", getSQLAuditRespBody, float64(0)).(float64)
	if keepDays == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", d.Id()),
		d.Set("keep_days", keepDays),
		d.Set("audit_types", utils.PathSearch("audit_types", getSQLAuditRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSQLAuditUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateSQLAudit: update RDS SQL audit
	var (
		updateSQLAuditProduct = "rds"
	)
	updateSQLAuditClient, err := cfg.NewServiceClient(updateSQLAuditProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	updateSQLAuditChanges := []string{
		"keep_days",
		"audit_types",
	}

	if d.HasChanges(updateSQLAuditChanges...) {
		params := buildUpdateSQLAuditBodyParams(d)
		log.Printf("[DEBUG] Update RDS SQL audit params: %#v", params)
		err = updateSQLAudit(ctx, d, updateSQLAuditClient, d.Timeout(schema.TimeoutUpdate), params)
		if err != nil {
			return diag.Errorf("error updating RDS SQL audit: %s", err)
		}

		// lintignore:R018
		time.Sleep(10 * time.Second)
	}

	return resourceSQLAuditRead(ctx, d, meta)
}

func resourceSQLAuditDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSQLAudit: delete RDS SQL audit
	var (
		deleteSQLAuditProduct = "rds"
	)
	deleteSQLAuditClient, err := cfg.NewServiceClient(deleteSQLAuditProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	params := buildDeleteSQLAuditBodyParams(d)
	log.Printf("[DEBUG] Delete RDS SQL audit params: %#v", params)
	err = updateSQLAudit(ctx, d, deleteSQLAuditClient, d.Timeout(schema.TimeoutDelete), params)
	if err != nil {
		return diag.Errorf("error deleting RDS SQL audit: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Refresh:      rdsSQLAuditStateRefreshFunc(deleteSQLAuditClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RDS instance(%s) SQL audit to be deleted: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateSQLAuditBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keep_days":   utils.ValueIgnoreEmpty(d.Get("keep_days")),
		"audit_types": d.Get("audit_types").(*schema.Set).List(),
	}
	return bodyParams
}

func buildDeleteSQLAuditBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keep_days":         0,
		"reserve_auditlogs": utils.ValueIgnoreEmpty(d.Get("reserve_auditlogs")),
	}
	return bodyParams
}

func updateSQLAudit(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout time.Duration,
	params map[string]interface{}) error {
	// updateSQLAudit: update RDS SQL audit
	var (
		updateSQLAuditHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
	)

	instanceID := d.Get("instance_id").(string)
	updateSQLAuditPath := client.Endpoint + updateSQLAuditHttpUrl
	updateSQLAuditPath = strings.ReplaceAll(updateSQLAuditPath, "{project_id}", client.ProjectID)
	updateSQLAuditPath = strings.ReplaceAll(updateSQLAuditPath, "{instance_id}", instanceID)

	updateSQLAuditOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateSQLAuditOpt.JSONBody = utils.RemoveNil(params)
	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("PUT", updateSQLAuditPath, &updateSQLAuditOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	return err
}

func rdsSQLAuditStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getSQLAuditRespBody, err := getSQLAudit(client, instanceID)
		if err != nil {
			return nil, "ERROR", err
		}
		keepDays := utils.PathSearch("keep_days", getSQLAuditRespBody, float64(0)).(float64)
		if keepDays == 0 {
			return getSQLAuditRespBody, "DELETED", nil
		}
		return getSQLAuditRespBody, "COMPLETED", nil
	}
}

func getSQLAudit(client *golangsdk.ServiceClient, instanceID string) (interface{}, error) {
	// getSQLAudit: Query the RDS SQL audit
	var (
		getSQLAuditHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
	)

	getSQLAuditPath := client.Endpoint + getSQLAuditHttpUrl
	getSQLAuditPath = strings.ReplaceAll(getSQLAuditPath, "{project_id}", client.ProjectID)
	getSQLAuditPath = strings.ReplaceAll(getSQLAuditPath, "{instance_id}", instanceID)

	getSQLAuditOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getSQLAuditResp, err := client.Request("GET", getSQLAuditPath, &getSQLAuditOpt)

	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getSQLAuditResp)
}
