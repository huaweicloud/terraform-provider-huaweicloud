package dcs

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

// @API DCS POST /v2/{project_id}/instances/{instance_id}/accounts
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS GET /v2/{project_id}/instances/{instance_id}/accounts
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/accounts/{account_id}/password/reset
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/accounts/{account_id}/role
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/accounts/{account_id}
// @API DCS DELETE /v2/{project_id}/instances/{instance_id}/accounts/{account_id}
// @API DCS GET /v2/{project_id}/instances/{instance_id}/tasks
func ResourceDcsAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsAccountCreate,
		UpdateContext: resourceDcsAccountUpdate,
		ReadContext:   resourceDcsAccountRead,
		DeleteContext: resourceDcsAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDcsAccountImportState,
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
				Description: "Specifies the ID of the DCS instance.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the account.`,
			},
			"account_role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the role of the account.`,
			},
			"account_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the password of the account.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the account.`,
			},
			"account_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the account.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the account.`,
			},
		},
	}
}

func buildCreateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"account_name":     d.Get("account_name"),
		"account_role":     d.Get("account_role"),
		"account_password": d.Get("account_password"),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceDcsAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// createAccount: create account
	var (
		createAccountHttpUrl = "v2/{project_id}/instances/{instance_id}/accounts"
		createAccountProduct = "dcs"
	)
	createAccountClient, err := cfg.NewServiceClient(createAccountProduct, region)

	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	createAccountPath := createAccountClient.Endpoint + createAccountHttpUrl
	createAccountPath = strings.ReplaceAll(createAccountPath, "{project_id}", createAccountClient.ProjectID)
	createAccountPath = strings.ReplaceAll(createAccountPath, "{instance_id}", instanceId)

	createAccountOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createAccountOpt.JSONBody = utils.RemoveNil(buildCreateAccountBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		createAccountResp, createErr := createAccountClient.Request("POST", createAccountPath, &createAccountOpt)
		retry, err := handleOperationError(createErr)
		return createAccountResp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(createAccountClient, instanceId),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error creating DCS account: %v", err)
	}

	_, err = utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"AVAILABLE"},
		Refresh:      accountStatusRefreshFunc(instanceId, accountName, createAccountClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}

	account, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for creating account (%s) to complete: %s", accountName, err)
	}

	accountId := utils.PathSearch("account_id", account, nil)
	if accountId == nil {
		return diag.Errorf("the account (%s) is not found", accountName)
	}

	d.SetId(accountId.(string))
	return resourceDcsAccountRead(ctx, d, meta)
}

func resourceDcsAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAccount: get account
	var (
		getAccountHttpUrl = "v2/{project_id}/instances/{instance_id}/accounts"
		getAccountProduct = "dcs"
	)

	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{project_id}", getAccountClient.ProjectID)
	getAccountPath = strings.ReplaceAll(getAccountPath, "{instance_id}", instanceId)

	getAccountOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	getAccountResp, err := getAccountClient.Request("GET", getAccountPath, &getAccountOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DCS accounts")
	}

	getAccountRespBody, err := utils.FlattenResponse(getAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	account := utils.PathSearch(fmt.Sprintf("accounts|[?account_id =='%s']|[0]", d.Id()), getAccountRespBody, nil)
	if account == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("account_name", utils.PathSearch("account_name", account, nil)),
		d.Set("account_type", utils.PathSearch("account_type", account, nil)),
		d.Set("account_role", utils.PathSearch("account_role", account, nil)),
		d.Set("status", utils.PathSearch("status", account, nil)),
		d.Set("description", utils.PathSearch("description", account, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcsAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// updateAccount: update account

	updateAccountProduct := "dcs"
	updateAccountClient, err := cfg.NewServiceClient(updateAccountProduct, region)

	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	if d.HasChange("account_password") {
		err := resetPassword(ctx, d, updateAccountClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("account_role") {
		err := updateRole(ctx, d, updateAccountClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("description") {
		err := updateDescription(ctx, d, updateAccountClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDcsAccountRead(ctx, d, meta)
}

func resetPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// resetPassword: reset password
	resetPasswordHttpUrl := "v2/{project_id}/instances/{instance_id}/accounts/{account_id}/password/reset"

	instanceId := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	resetPasswordPath := client.Endpoint + resetPasswordHttpUrl
	resetPasswordPath = strings.ReplaceAll(resetPasswordPath, "{project_id}", client.ProjectID)
	resetPasswordPath = strings.ReplaceAll(resetPasswordPath, "{instance_id}", instanceId)
	resetPasswordPath = strings.ReplaceAll(resetPasswordPath, "{account_id}", d.Id())

	resetPasswordOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	resetPasswordOpt.JSONBody = map[string]interface{}{"new_password": d.Get("account_password")}

	retryFunc := func() (interface{}, bool, error) {
		resetPasswordResp, updateErr := client.Request("PUT", resetPasswordPath, &resetPasswordOpt)
		retry, err := handleOperationError(updateErr)
		return resetPasswordResp, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceId),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error resetting password of the account(%s): %v", accountName, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"UPDATING"},
		Target:       []string{"AVAILABLE"},
		Refresh:      accountStatusRefreshFunc(instanceId, accountName, client),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for resetting password of the account (%s) to complete: %s", accountName, err)
	}

	return nil
}

func updateRole(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// updateRole: update role
	updateRoleHttpUrl := "v2/{project_id}/instances/{instance_id}/accounts/{account_id}/role"

	instanceId := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	updateRolePath := client.Endpoint + updateRoleHttpUrl
	updateRolePath = strings.ReplaceAll(updateRolePath, "{project_id}", client.ProjectID)
	updateRolePath = strings.ReplaceAll(updateRolePath, "{instance_id}", instanceId)
	updateRolePath = strings.ReplaceAll(updateRolePath, "{account_id}", d.Id())

	updateRoleOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	updateRoleOpt.JSONBody = map[string]interface{}{"account_role": d.Get("account_role")}

	retryFunc := func() (interface{}, bool, error) {
		updateRoleResp, updateErr := client.Request("PUT", updateRolePath, &updateRoleOpt)
		retry, err := handleOperationError(updateErr)
		return updateRoleResp, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceId),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error updating role of the account (%s): %v", accountName, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"UPDATING"},
		Target:       []string{"AVAILABLE"},
		Refresh:      accountStatusRefreshFunc(instanceId, accountName, client),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for updating role of the account (%s) to complete: %s", accountName, err)
	}

	return nil
}

func updateDescription(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// updateDescription: update description
	updateDescriptionHttpUrl := "v2/{project_id}/instances/{instance_id}/accounts/{account_id}"

	instanceId := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	updateDescriptionPath := client.Endpoint + updateDescriptionHttpUrl
	updateDescriptionPath = strings.ReplaceAll(updateDescriptionPath, "{project_id}", client.ProjectID)
	updateDescriptionPath = strings.ReplaceAll(updateDescriptionPath, "{instance_id}", instanceId)
	updateDescriptionPath = strings.ReplaceAll(updateDescriptionPath, "{account_id}", d.Id())

	updateDescriptionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	updateDescriptionOpt.JSONBody = map[string]interface{}{"description": d.Get("description")}

	retryFunc := func() (interface{}, bool, error) {
		updateDescriptionResp, err := client.Request("PUT", updateDescriptionPath, &updateDescriptionOpt)
		retry, err := handleOperationError(err)
		return updateDescriptionResp, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceId),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error updating description of the account (%s): %v", accountName, err)
	}

	return nil
}

func resourceDcsAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAccount: delete account
	var (
		deleteAccountHttpUrl = "v2/{project_id}/instances/{instance_id}/accounts/{account_id}"
		deleteAccountProduct = "dcs"
	)
	deleteAccountClient, err := cfg.NewServiceClient(deleteAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS Client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	deleteAccountPath := deleteAccountClient.Endpoint + deleteAccountHttpUrl
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{project_id}", deleteAccountClient.ProjectID)
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{instance_id}", instanceId)
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{account_id}", d.Id())

	deleteAccountOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	retryFunc := func() (interface{}, bool, error) {
		_, err = deleteAccountClient.Request("DELETE", deleteAccountPath, &deleteAccountOpt)
		retry, err := handleOperationError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(deleteAccountClient, instanceId),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error deleting the account (%s): %v", accountName, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EXECUTING"},
		Target:       []string{"SUCCESS"},
		Refresh:      accountTaskStatusRefreshFunc(instanceId, "DeleteAcl", deleteAccountClient),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for deleting the account (%s) to complete: %s", accountName, err)
	}

	return resourceDcsAccountRead(ctx, d, meta)
}

func accountStatusRefreshFunc(instanceId, accountName string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getAccountHttpUrl = "v2/{project_id}/instances/{instance_id}/accounts"
		)
		getAccountPath := client.Endpoint + getAccountHttpUrl
		getAccountPath = strings.ReplaceAll(getAccountPath, "{project_id}", client.ProjectID)
		getAccountPath = strings.ReplaceAll(getAccountPath, "{instance_id}", instanceId)

		getAccountOpt := golangsdk.RequestOpts{KeepResponseBody: true}

		getAccountResp, err := client.Request("GET", getAccountPath, &getAccountOpt)
		if err != nil {
			return nil, "QUERY ERROR", err
		}
		getAccountBody, err := utils.FlattenResponse(getAccountResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		account := utils.PathSearch(fmt.Sprintf("accounts|[?account_name == '%s']|[0]", accountName), getAccountBody, nil)
		if account == nil {
			return nil, "NOT FOUND", nil
		}
		status := utils.PathSearch("status", account, "").(string)
		return account, status, nil
	}
}

func accountTaskStatusRefreshFunc(instanceId, taskName string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getAccountTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/tasks"
		getAccountTaskPath := client.Endpoint + getAccountTaskHttpUrl
		getAccountTaskPath = strings.ReplaceAll(getAccountTaskPath, "{project_id}", client.ProjectID)
		getAccountTaskPath = strings.ReplaceAll(getAccountTaskPath, "{instance_id}", instanceId)

		reqOpt := golangsdk.RequestOpts{KeepResponseBody: true}
		getAccountTaskResp, err := client.Request("GET", getAccountTaskPath, &reqOpt)

		if err != nil {
			return nil, "QUERY ERROR", err
		}

		getAccountTaskRespBody, err := utils.FlattenResponse(getAccountTaskResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		task := utils.PathSearch(fmt.Sprintf("tasks|[?name=='%s']|[0]", taskName), getAccountTaskRespBody, nil)
		if task == nil {
			return nil, "NOT FOUND", nil
		}
		status := utils.PathSearch("status", task, "").(string)
		return task, status, nil
	}
}

func resourceDcsAccountImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
