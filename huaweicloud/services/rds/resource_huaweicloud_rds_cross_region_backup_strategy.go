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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/backups/offsite-policy
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/backups/offsite-policy
// @API RDS GET /v3/{project_id}/instances
func ResourceBackupStrategy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupStrategyCreate,
		UpdateContext: resourceBackupStrategyUpdate,
		ReadContext:   resourceBackupStrategyRead,
		DeleteContext: resourceBackupStrategyDelete,
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
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  `Specifies the backup type.`,
				ValidateFunc: validation.StringInSlice([]string{"auto", "all"}, true),
			},
			"keep_days": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `Specifies the number of days to retain the generated backup files.`,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"destination_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the target region ID for the cross-region backup policy.`,
			},
			"destination_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the target project ID for the cross-region backup policy.`,
			},
		},
	}
}

func resourceBackupStrategyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createBackupStrategy: create RDS cross-region backup strategy.
	var (
		createBackupStrategyProduct = "rds"
	)
	createBackupStrategyClient, err := cfg.NewServiceClient(createBackupStrategyProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	err = updateCrossRegionBackupStrategy(ctx, createBackupStrategyClient, d.Timeout(schema.TimeoutCreate), instanceID,
		buildCreateBackupStrategyBodyParams(d))
	if err != nil {
		return diag.Errorf("error creating RDS cross region backup strategy: %s", err)
	}

	d.SetId(instanceID)

	return resourceBackupStrategyRead(ctx, d, meta)
}

func buildCreateBackupStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_para": map[string]interface{}{
			"backup_type":            d.Get("backup_type"),
			"keep_days":              d.Get("keep_days"),
			"destination_region":     d.Get("destination_region"),
			"destination_project_id": d.Get("destination_project_id"),
		},
	}
	return bodyParams
}

func resourceBackupStrategyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getBackupStrategy: Query the RDS cross-region backup strategy
	var (
		getBackupStrategyHttpUrl = "v3/{project_id}/instances/{instance_id}/backups/offsite-policy"
		getBackupStrategyProduct = "rds"
	)
	getBackupStrategyClient, err := cfg.NewServiceClient(getBackupStrategyProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instance, err := GetRdsInstanceByID(getBackupStrategyClient, d.Id())
	if err != nil {
		return diag.Errorf("error getting RDS instance: %s", err)
	}
	instanceId := utils.PathSearch("id", instance, "").(string)
	if instanceId == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	getBackupStrategyPath := getBackupStrategyClient.Endpoint + getBackupStrategyHttpUrl
	getBackupStrategyPath = strings.ReplaceAll(getBackupStrategyPath, "{project_id}",
		getBackupStrategyClient.ProjectID)
	getBackupStrategyPath = strings.ReplaceAll(getBackupStrategyPath, "{instance_id}", d.Id())

	getBackupStrategyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getBackupStrategyResp, err := getBackupStrategyClient.Request("GET", getBackupStrategyPath, &getBackupStrategyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS cross region backup strategy")
	}

	getBackupStrategyRespBody, err := utils.FlattenResponse(getBackupStrategyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyPara := utils.PathSearch("policy_para", getBackupStrategyRespBody, nil)
	if policyPara == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	backupStrategies := policyPara.([]interface{})
	if len(backupStrategies) == 0 || utils.PathSearch("keep_days", backupStrategies[0], float64(0)).(float64) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", d.Id()),
		d.Set("keep_days", utils.PathSearch("keep_days", backupStrategies[0], nil)),
		d.Set("destination_region", utils.PathSearch("destination_region", backupStrategies[0], nil)),
		d.Set("destination_project_id", utils.PathSearch("destination_project_id", backupStrategies[0], nil)),
	)

	if len(backupStrategies) == 1 {
		mErr = multierror.Append(mErr, d.Set("backup_type", "auto"))
	} else {
		mErr = multierror.Append(mErr, d.Set("backup_type", "all"))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBackupStrategyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateBackupStrategyChanges := []string{
		"backup_type",
		"keep_days",
	}

	if d.HasChanges(updateBackupStrategyChanges...) {
		// updateBackupStrategy: update a RDS cross-region backup strategy.
		var (
			updateBackupStrategyProduct = "rds"
		)
		updateBackupStrategyClient, err := cfg.NewServiceClient(updateBackupStrategyProduct, region)
		if err != nil {
			return diag.Errorf("error creating RDS client: %s", err)
		}

		if d.HasChange("keep_days") {
			err = updateCrossRegionBackupStrategy(ctx, updateBackupStrategyClient, d.Timeout(schema.TimeoutUpdate),
				d.Get("instance_id").(string), buildUpdateBackupStrategyKeepDaysBodyParams(d))
			if err != nil {
				return diag.Errorf("error updating RDS cross region backup strategy: %s", err)
			}
		}
		// If backup_type is changed from `all` to `auto`, we should close the incremental backup
		if d.HasChange("backup_type") {
			err = updateCrossRegionBackupStrategy(ctx, updateBackupStrategyClient, d.Timeout(schema.TimeoutUpdate),
				d.Get("instance_id").(string), buildUpdateBackupStrategyBackupTypeBodyParams(d))
			if err != nil {
				return diag.Errorf("error updating RDS cross region backup strategy: %s", err)
			}
		}
	}
	return resourceBackupStrategyRead(ctx, d, meta)
}

func buildUpdateBackupStrategyKeepDaysBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_para": map[string]interface{}{
			"backup_type":            utils.ValueIgnoreEmpty(d.Get("backup_type")),
			"keep_days":              utils.ValueIgnoreEmpty(d.Get("keep_days")),
			"destination_region":     utils.ValueIgnoreEmpty(d.Get("destination_region")),
			"destination_project_id": utils.ValueIgnoreEmpty(d.Get("destination_project_id")),
		},
	}
	return bodyParams
}

func buildUpdateBackupStrategyBackupTypeBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"backup_type":            "incremental",
		"destination_region":     utils.ValueIgnoreEmpty(d.Get("destination_region")),
		"destination_project_id": utils.ValueIgnoreEmpty(d.Get("destination_project_id")),
	}
	// If backup_type is changed to auto, it indicates the original value is all, we should close incremental backup,
	// so we should set keep_days to 0, on the contrary, if backup_type is changed to all, it indicates the original
	// value is auto, so we should set keep_days to the value of the input param keep_days
	if d.Get("backup_type") == "auto" {
		params["keep_days"] = 0
	} else {
		params["keep_days"] = d.Get("keep_days").(int)
	}
	bodyParams := map[string]interface{}{
		"policy_para": params,
	}
	return bodyParams
}

func resourceBackupStrategyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteBackupStrategy: delete a RDS cross-region backup strategy.
	var (
		deleteBackupStrategyProduct = "rds"
	)
	deleteBackupStrategyClient, err := cfg.NewServiceClient(deleteBackupStrategyProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	err = updateCrossRegionBackupStrategy(ctx, deleteBackupStrategyClient, d.Timeout(schema.TimeoutDelete),
		d.Get("instance_id").(string), buildDeleteBackupStrategyBodyParams(d))
	if err != nil {
		return diag.Errorf("error deleting RDS cross region backup strategy: %s", err)
	}

	return nil
}

func buildDeleteBackupStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_para": map[string]interface{}{
			"backup_type":            "all",
			"keep_days":              0,
			"destination_region":     utils.ValueIgnoreEmpty(d.Get("destination_region")),
			"destination_project_id": utils.ValueIgnoreEmpty(d.Get("destination_project_id")),
		},
	}
	return bodyParams
}

func updateCrossRegionBackupStrategy(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceID string, params map[string]interface{}) error {
	var (
		updateBackupStrategyHttpUrl = "v3/{project_id}/instances/{instance_id}/backups/offsite-policy"
	)

	updateBackupStrategyPath := client.Endpoint + updateBackupStrategyHttpUrl
	updateBackupStrategyPath = strings.ReplaceAll(updateBackupStrategyPath, "{project_id}", client.ProjectID)
	updateBackupStrategyPath = strings.ReplaceAll(updateBackupStrategyPath, "{instance_id}", instanceID)

	deleteBackupStrategyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	log.Printf("[DEBUG] Update RDS instance cross region backup strategy params: %#v", params)
	deleteBackupStrategyOpt.JSONBody = utils.RemoveNil(params)

	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("PUT", updateBackupStrategyPath, &deleteBackupStrategyOpt)
		retry, err := handleCrossRegionBackupStrategyError(err)
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
