// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
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

func ResourceSQLAudit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSQLAuditCreate,
		UpdateContext: resourceSQLAuditUpdate,
		ReadContext:   resourceSQLAuditRead,
		DeleteContext: resourceSQLAuditDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				ForceNew:    true,
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
		},
	}
}

func resourceSQLAuditCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	params := buildSetSQLAuditBodyParams(d)
	err := updateSQLAudit(ctx, d, meta, d.Timeout(schema.TimeoutCreate), params)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	d.SetId(instanceID)

	return resourceSQLAuditRead(ctx, d, meta)
}

func resourceSQLAuditRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSQLAudit: Query the RDS SQL audit
	var (
		getSQLAuditHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
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
	if instance.Id == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	getSQLAuditPath := getSQLAuditClient.Endpoint + getSQLAuditHttpUrl
	getSQLAuditPath = strings.ReplaceAll(getSQLAuditPath, "{project_id}", getSQLAuditClient.ProjectID)
	getSQLAuditPath = strings.ReplaceAll(getSQLAuditPath, "{instance_id}", d.Id())

	getSQLAuditOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getSQLAuditResp, err := getSQLAuditClient.Request("GET", getSQLAuditPath, &getSQLAuditOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS SQL audit")
	}

	getSQLAuditRespBody, err := utils.FlattenResponse(getSQLAuditResp)
	if err != nil {
		return diag.FromErr(err)
	}

	keepDays := utils.PathSearch("keep_days", getSQLAuditRespBody, 0).(float64)
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
	updateSQLAuditChanges := []string{
		"keep_days",
		"audit_types",
	}

	if d.HasChanges(updateSQLAuditChanges...) {
		params := buildSetSQLAuditBodyParams(d)
		err := updateSQLAudit(ctx, d, meta, d.Timeout(schema.TimeoutUpdate), params)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSQLAuditRead(ctx, d, meta)
}

func resourceSQLAuditDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	params := buildDeleteSQLAuditBodyParams(d)
	err := updateSQLAudit(ctx, d, meta, d.Timeout(schema.TimeoutDelete), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildSetSQLAuditBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keep_days":   utils.ValueIngoreEmpty(d.Get("keep_days")),
		"audit_types": d.Get("audit_types").(*schema.Set).List(),
	}
	return bodyParams
}

func buildDeleteSQLAuditBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keep_days":         0,
		"reserve_auditlogs": utils.ValueIngoreEmpty(d.Get("reserve_auditlogs")),
	}
	return bodyParams
}

func updateSQLAudit(ctx context.Context, d *schema.ResourceData, meta interface{}, timeout time.Duration,
	params map[string]interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// setSQLAudit: set RDS SQL audit
	var (
		setSQLAuditHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		setSQLAuditProduct = "rds"
	)
	setSQLAuditClient, err := cfg.NewServiceClient(setSQLAuditProduct, region)
	if err != nil {
		return fmt.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	setSQLAuditPath := setSQLAuditClient.Endpoint + setSQLAuditHttpUrl
	setSQLAuditPath = strings.ReplaceAll(setSQLAuditPath, "{project_id}", setSQLAuditClient.ProjectID)
	setSQLAuditPath = strings.ReplaceAll(setSQLAuditPath, "{instance_id}", instanceID)

	setSQLAuditOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	setSQLAuditOpt.JSONBody = utils.RemoveNil(params)
	retryFunc := func() (interface{}, bool, error) {
		_, err = setSQLAuditClient.Request("PUT", setSQLAuditPath, &setSQLAuditOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(setSQLAuditClient, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating RDS SQL audit: %s", err)
	}

	// lintignore:R018
	time.Sleep(10 * time.Second)

	return nil
}
