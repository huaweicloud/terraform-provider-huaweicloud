// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"encoding/json"
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

var sqlServerAccountNonUpdatableParams = []string{"instance_id", "name"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_user
// @API RDS GET /v3/{project_id}/instances
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db_user/{user_name}
// @API RDS GET /v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_user/resetpwd
func ResourceSQLServerAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSQLServerAccountCreate,
		UpdateContext: resourceSQLServerAccountUpdate,
		ReadContext:   resourceSQLServerAccountRead,
		DeleteContext: resourceSQLServerAccountDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(sqlServerAccountNonUpdatableParams),

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
				Description: `Specifies the ID of the RDS SQLServer instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the username of the DB account.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the password of the DB account.`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DB user status.`,
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

func resourceSQLServerAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSQLServerAccount: create RDS SQLServer account.
	var (
		createSQLServerAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user"
		createSQLServerAccountProduct = "rds"
	)
	createSQLServerAccountClient, err := cfg.NewServiceClient(createSQLServerAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createSQLServerAccountPath := createSQLServerAccountClient.Endpoint + createSQLServerAccountHttpUrl
	createSQLServerAccountPath = strings.ReplaceAll(createSQLServerAccountPath, "{project_id}",
		createSQLServerAccountClient.ProjectID)
	createSQLServerAccountPath = strings.ReplaceAll(createSQLServerAccountPath, "{instance_id}", instanceId)

	createSQLServerAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createSQLServerAccountOpt.JSONBody = utils.RemoveNil(buildCreateSQLServerAccountBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		_, err = createSQLServerAccountClient.Request("POST", createSQLServerAccountPath, &createSQLServerAccountOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(createSQLServerAccountClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS SQLServer account: %s", err)
	}

	accountName := d.Get("name").(string)
	d.SetId(instanceId + "/" + accountName)

	return resourceSQLServerAccountRead(ctx, d, meta)
}

func buildCreateSQLServerAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     d.Get("name"),
		"password": d.Get("password"),
	}
	return bodyParams
}

func resourceSQLServerAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSQLServerAccount: query RDS SQLServer account
	var (
		getSQLServerAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		getSQLServerAccountProduct = "rds"
	)
	getSQLServerAccountClient, err := cfg.NewServiceClient(getSQLServerAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and account name from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getSQLServerAccountPath := getSQLServerAccountClient.Endpoint + getSQLServerAccountHttpUrl
	getSQLServerAccountPath = strings.ReplaceAll(getSQLServerAccountPath, "{project_id}",
		getSQLServerAccountClient.ProjectID)
	getSQLServerAccountPath = strings.ReplaceAll(getSQLServerAccountPath, "{instance_id}", instanceId)

	getSQLServerAccountResp, err := pagination.ListAllItems(
		getSQLServerAccountClient,
		"page",
		getSQLServerAccountPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS SQLServer account")
	}

	getSQLServerAccountRespJson, err := json.Marshal(getSQLServerAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getSQLServerAccountRespBody interface{}
	err = json.Unmarshal(getSQLServerAccountRespJson, &getSQLServerAccountRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getSQLServerAccountRespBody, nil)

	if account == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", utils.PathSearch("name", account, nil)),
		d.Set("state", utils.PathSearch("state", account, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSQLServerAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSQLServerAccountChanges := []string{
		"password",
	}

	if d.HasChanges(updateSQLServerAccountChanges...) {
		// updateSQLServerAccount: update RDS SQLServer account password
		var (
			updateSQLServerAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/resetpwd"
			updateSQLServerAccountProduct = "rds"
		)
		updateSQLServerAccountClient, err := cfg.NewServiceClient(updateSQLServerAccountProduct, region)
		if err != nil {
			return diag.Errorf("error creating RDS client: %s", err)
		}

		instanceId := d.Get("instance_id").(string)
		updateSQLServerAccountPath := updateSQLServerAccountClient.Endpoint + updateSQLServerAccountHttpUrl
		updateSQLServerAccountPath = strings.ReplaceAll(updateSQLServerAccountPath, "{project_id}",
			updateSQLServerAccountClient.ProjectID)
		updateSQLServerAccountPath = strings.ReplaceAll(updateSQLServerAccountPath, "{instance_id}", instanceId)

		updateSQLServerAccountOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		updateSQLServerAccountOpt.JSONBody = utils.RemoveNil(buildUpdateSQLServerAccountBodyParams(d))
		retryFunc := func() (interface{}, bool, error) {
			_, err = updateSQLServerAccountClient.Request("POST", updateSQLServerAccountPath, &updateSQLServerAccountOpt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(updateSQLServerAccountClient, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		_, err = updateSQLServerAccountClient.Request("POST", updateSQLServerAccountPath, &updateSQLServerAccountOpt)
		if err != nil {
			return diag.Errorf("error updating RDS SQLServer account: %s", err)
		}
	}
	return resourceSQLServerAccountRead(ctx, d, meta)
}

func buildUpdateSQLServerAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(d.Get("name")),
		"password": utils.ValueIgnoreEmpty(d.Get("password")),
	}
	return bodyParams
}

func resourceSQLServerAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSQLServerAccount: delete RDS SQLServer account
	var (
		deleteSQLServerAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/{user_name}"
		deleteSQLServerAccountProduct = "rds"
	)
	deleteSQLServerAccountClient, err := cfg.NewServiceClient(deleteSQLServerAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deleteSQLServerAccountPath := deleteSQLServerAccountClient.Endpoint + deleteSQLServerAccountHttpUrl
	deleteSQLServerAccountPath = strings.ReplaceAll(deleteSQLServerAccountPath, "{project_id}",
		deleteSQLServerAccountClient.ProjectID)
	deleteSQLServerAccountPath = strings.ReplaceAll(deleteSQLServerAccountPath, "{instance_id}", instanceId)
	deleteSQLServerAccountPath = strings.ReplaceAll(deleteSQLServerAccountPath, "{user_name}",
		d.Get("name").(string))

	deleteSQLServerAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err = deleteSQLServerAccountClient.Request("DELETE", deleteSQLServerAccountPath, &deleteSQLServerAccountOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deleteSQLServerAccountClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting RDS SQLServer account: %s", err)
	}

	return nil
}
