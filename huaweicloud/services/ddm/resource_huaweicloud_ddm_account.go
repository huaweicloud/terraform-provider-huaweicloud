// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM POST /v1/{project_id}/instances/{instance_id}/users
// @API DDM GET /v1/{project_id}/instances/{instance_id}
// @API DDM PUT /v1/{project_id}/instances/{instance_id}/users/{username}
// @API DDM POST /v2/{project_id}/instances/{instance_id}/users/{username}/password
// @API DDM GET /v1/{project_id}/instances/{instance_id}/users
// @API DDM DELETE /v1/{project_id}/instances/{instance_id}/users/{username}
func ResourceDdmAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdmAccountCreate,
		UpdateContext: resourceDdmAccountUpdate,
		ReadContext:   resourceDdmAccountRead,
		DeleteContext: resourceDdmAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDdmAccountImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
				Description: `Specifies the ID of a DDM instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the DDM account.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the DDM account password.`,
			},
			"permissions": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the basic permissions of the DDM account.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the DDM account.`,
			},
			"schemas": {
				Type:        schema.TypeSet,
				Elem:        accountSchemaSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the schemas that associated with the account.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DDM account.`,
			},
		},
	}
}

func accountSchemaSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the associated schema.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the schema description.`,
			},
		},
	}
	return &sc
}

func resourceDdmAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAccount: create DDM account
	var (
		createAccountHttpUrl = "v1/{project_id}/instances/{instance_id}/users"
		createAccountProduct = "ddm"
	)
	createAccountClient, err := cfg.NewServiceClient(createAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createAccountPath := createAccountClient.Endpoint + createAccountHttpUrl
	createAccountPath = strings.ReplaceAll(createAccountPath, "{project_id}", createAccountClient.ProjectID)
	createAccountPath = strings.ReplaceAll(createAccountPath, "{instance_id}", instanceID)

	createAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccountOpt.JSONBody = utils.RemoveNil(buildCreateAccountBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := createAccountClient.Request("POST", createAccountPath, &createAccountOpt)
		retry, err := handleOperationError(err, "creating", "account")
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddmInstanceStatusRefreshFunc(instanceID, createAccountClient),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID + "/" + d.Get("name").(string))

	return resourceDdmAccountRead(ctx, d, meta)
}

func buildCreateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"password":       d.Get("password"),
		"base_authority": d.Get("permissions").(*schema.Set).List(),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"databases":      buildCreateAccountSchemasChildBody(d),
	}
	params := map[string]interface{}{
		"users": []interface{}{bodyParams},
	}
	return params
}

func buildCreateAccountSchemasChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("schemas").(*schema.Set)
	if rawParams.Len() == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0)
	for _, param := range rawParams.List() {
		perm := make(map[string]interface{})
		perm["name"] = utils.PathSearch("name", param, nil)
		params = append(params, perm)
	}
	return params
}

func resourceDdmAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	updateAccountHasChanges := []string{
		"permissions",
		"description",
		"schemas",
	}

	if d.HasChanges(updateAccountHasChanges...) {
		err = updateAccount(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password") {
		err = updateAccountPassword(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDdmAccountRead(ctx, d, meta)
}

func updateAccount(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		updateAccountHttpUrl = "v1/{project_id}/instances/{instance_id}/users/{username}"
	)

	instanceId := d.Get("instance_id").(string)
	updateAccountPath := client.Endpoint + updateAccountHttpUrl
	updateAccountPath = strings.ReplaceAll(updateAccountPath, "{project_id}", client.ProjectID)
	updateAccountPath = strings.ReplaceAll(updateAccountPath, "{instance_id}", instanceId)
	updateAccountPath = strings.ReplaceAll(updateAccountPath, "{username}", d.Get("name").(string))

	updateAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateAccountOpt.JSONBody = utils.RemoveNil(buildUpdateAccountBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updateAccountPath, &updateAccountOpt)
		retry, err := handleOperationError(err, "updating", "account")
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddmInstanceStatusRefreshFunc(instanceId, client),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	return err
}

func updateAccountPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		updateAccountPasswordHttpUrl = "v2/{project_id}/instances/{instance_id}/users/{username}/password"
	)

	instanceId := d.Get("instance_id").(string)
	updateAccountPasswordPath := client.Endpoint + updateAccountPasswordHttpUrl
	updateAccountPasswordPath = strings.ReplaceAll(updateAccountPasswordPath, "{project_id}", client.ProjectID)
	updateAccountPasswordPath = strings.ReplaceAll(updateAccountPasswordPath, "{instance_id}", instanceId)
	updateAccountPasswordPath = strings.ReplaceAll(updateAccountPasswordPath, "{username}", d.Get("name").(string))

	updateAccountPasswordOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateAccountPasswordOpt.JSONBody = utils.RemoveNil(buildUpdateAccountPasswordBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", updateAccountPasswordPath, &updateAccountPasswordOpt)
		retry, err := handleOperationError(err, "updating", "account password")
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddmInstanceStatusRefreshFunc(instanceId, client),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	return err
}

func buildUpdateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"base_authority": d.Get("permissions").(*schema.Set).List(),
		"description":    d.Get("description"),
		"databases":      buildUpdateAccountSchemasChildBody(d),
	}
	bodyParams := map[string]interface{}{
		"user": params,
	}
	return bodyParams
}

func buildUpdateAccountSchemasChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("schemas").(*schema.Set)
	if rawParams.Len() == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0)
	for _, param := range rawParams.List() {
		perm := make(map[string]interface{})
		perm["name"] = utils.PathSearch("name", param, nil)
		params = append(params, perm)
	}
	return params
}

func buildUpdateAccountPasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password": d.Get("password"),
	}
	return bodyParams
}

func resourceDdmAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getAccountHttpUrl = "v1/{project_id}/instances/{instance_id}/users"
		getAccountProduct = "ddm"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{project_id}", getAccountClient.ProjectID)
	getAccountPath = strings.ReplaceAll(getAccountPath, "{instance_id}", d.Get("instance_id").(string))

	getAccountResp, err := pagination.ListAllItems(
		getAccountClient,
		"offset",
		getAccountPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DDM account")
	}

	getAccountRespJson, err := json.Marshal(getAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var getAccountRespBody interface{}
	err = json.Unmarshal(getAccountRespJson, &getAccountRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	account := utils.PathSearch(fmt.Sprintf("users|[?name=='%s']|[0]", name), getAccountRespBody, nil)
	if account == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", name),
		d.Set("status", utils.PathSearch("status", account, nil)),
		d.Set("permissions", utils.PathSearch("base_authority", account, nil)),
		d.Set("description", utils.PathSearch("description", account, nil)),
		d.Set("schemas", flattenGetAccountResponseBodyDatabase(account)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAccountResponseBodyDatabase(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func resourceDdmAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAccount: Delete DDM account
	var (
		deleteAccountHttpUrl = "v1/{project_id}/instances/{instance_id}/users/{username}"
		deleteAccountProduct = "ddm"
	)
	deleteAccountClient, err := cfg.NewServiceClient(deleteAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deleteAccountPath := deleteAccountClient.Endpoint + deleteAccountHttpUrl
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{project_id}", deleteAccountClient.ProjectID)
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{instance_id}", instanceId)
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{username}", d.Get("name").(string))

	deleteAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := deleteAccountClient.Request("DELETE", deleteAccountPath, &deleteAccountOpt)
		retry, err := handleOperationError(err, "deleting", "account")
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddmInstanceStatusRefreshFunc(instanceId, deleteAccountClient),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DDM account")
	}

	return nil
}

func resourceDdmAccountImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <instance_id>/<name>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
