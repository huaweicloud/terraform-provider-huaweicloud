// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM POST /v1/{project_id}/instances/{instance_id}/users
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
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]\w*$`),
						"An account name starts with a letter, consists of 1 to 32 characters,"+
							"and can contain only letters, digits, and underscores (_)"),
					validation.StringLenBetween(1, 32),
				),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the DDM account password.`,
			},
			"permissions": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"CREATE", "DROP", "ALTER", "INDEX", "INSERT", "DELETE", "UPDATE", "SELECT",
					}, false),
				},
				Required:    true,
				Description: `Specifies the basic permissions of the DDM account.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the DDM account.`,
			},
			"schemas": {
				Type:        schema.TypeList,
				Elem:        AccountSchemaSchema(),
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

func AccountSchemaSchema() *schema.Resource {
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
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createAccountPath := createAccountClient.Endpoint + createAccountHttpUrl
	createAccountPath = strings.ReplaceAll(createAccountPath, "{project_id}", createAccountClient.ProjectID)
	createAccountPath = strings.ReplaceAll(createAccountPath, "{instance_id}", fmt.Sprintf("%v", instanceID))

	createAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createAccountOpt.JSONBody = utils.RemoveNil(buildCreateAccountBodyParams(d))

	var createAccountResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		createAccountResp, err = createAccountClient.Request("POST", createAccountPath, &createAccountOpt)
		isRetry, err := handleOperationError(err, "creating", "account")
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	createAccountRespBody, err := utils.FlattenResponse(createAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	accounts, err := jmespath.Search("users", createAccountRespBody)
	if err != nil {
		return diag.Errorf("error creating DDM account, users is not found in API response %s", err)
	}
	accountName, err := jmespath.Search("name", accounts.([]interface{})[0])
	if err != nil {
		return diag.Errorf("error creating DDM account, name is not found in API response %s", err)
	}

	d.SetId(instanceID + "/" + accountName.(string))

	return resourceDdmAccountRead(ctx, d, meta)
}

func buildCreateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           utils.ValueIgnoreEmpty(d.Get("name")),
		"password":       utils.ValueIgnoreEmpty(d.Get("password")),
		"base_authority": utils.ValueIgnoreEmpty(d.Get("permissions")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"databases":      buildCreateAccountSchemasChildBody(d),
	}
	params := map[string]interface{}{
		"users": []interface{}{bodyParams},
	}
	return params
}

func buildCreateAccountSchemasChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("schemas").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0)
	for _, param := range rawParams {
		perm := make(map[string]interface{})
		perm["name"] = utils.PathSearch("name", param, nil)
		params = append(params, perm)
	}
	return params
}

func resourceDdmAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAccountHasChanges := []string{
		"permissions",
		"description",
		"schemas",
	}

	if d.HasChanges(updateAccountHasChanges...) {
		// updateAccount: update DDM account
		err := updateAccount(ctx, d, cfg, region)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password") {
		// updateAccountPassword: update DDM account password
		err := updateAccountPassword(ctx, d, cfg, region)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDdmAccountRead(ctx, d, meta)
}

func updateAccount(ctx context.Context, d *schema.ResourceData, cfg *config.Config, region string) error {
	var (
		updateAccountHttpUrl = "v1/{project_id}/instances/{instance_id}/users/{username}"
		updateAccountProduct = "ddm"
	)
	updateAccountClient, err := cfg.NewServiceClient(updateAccountProduct, region)
	if err != nil {
		return fmt.Errorf("error creating DDM Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid id format, must be <instance_id>/<account_name>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	updateAccountPath := updateAccountClient.Endpoint + updateAccountHttpUrl
	updateAccountPath = strings.ReplaceAll(updateAccountPath, "{project_id}", updateAccountClient.ProjectID)
	updateAccountPath = strings.ReplaceAll(updateAccountPath, "{instance_id}", fmt.Sprintf("%v", instanceID))
	updateAccountPath = strings.ReplaceAll(updateAccountPath, "{username}", fmt.Sprintf("%v", accountName))

	updateAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	updateAccountOpt.JSONBody = utils.RemoveNil(buildUpdateAccountBodyParams(d))

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err = updateAccountClient.Request("PUT", updateAccountPath, &updateAccountOpt)
		isRetry, err := handleOperationError(err, "updating", "account")
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return err
}

func updateAccountPassword(ctx context.Context, d *schema.ResourceData, cfg *config.Config, region string) error {
	var (
		updateAccountPasswordHttpUrl = "v2/{project_id}/instances/{instance_id}/users/{username}/password"
		updateAccountPasswordProduct = "ddm"
	)
	updateAccountPasswordClient, err := cfg.NewServiceClient(updateAccountPasswordProduct, region)
	if err != nil {
		return fmt.Errorf("error creating DDM Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid id format, must be <instance_id>/<account_name>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	updateAccountPasswordPath := updateAccountPasswordClient.Endpoint + updateAccountPasswordHttpUrl
	updateAccountPasswordPath = strings.ReplaceAll(updateAccountPasswordPath, "{project_id}",
		updateAccountPasswordClient.ProjectID)
	updateAccountPasswordPath = strings.ReplaceAll(updateAccountPasswordPath, "{instance_id}",
		fmt.Sprintf("%v", instanceID))
	updateAccountPasswordPath = strings.ReplaceAll(updateAccountPasswordPath, "{username}",
		fmt.Sprintf("%v", accountName))

	updateAccountPasswordOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	updateAccountPasswordOpt.JSONBody = utils.RemoveNil(buildUpdateAccountPasswordBodyParams(d))

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err = updateAccountPasswordClient.Request("POST", updateAccountPasswordPath, &updateAccountPasswordOpt)
		isRetry, err := handleOperationError(err, "updating", "account password")
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return err
}

func buildUpdateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"base_authority": utils.ValueIgnoreEmpty(d.Get("permissions")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"databases":      buildUpdateAccountSchemasChildBody(d),
	}
	bodyParams := map[string]interface{}{
		"user": params,
	}
	return bodyParams
}

func buildUpdateAccountSchemasChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("schemas").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0)
	for _, param := range rawParams {
		perm := make(map[string]interface{})
		perm["name"] = utils.PathSearch("name", param, nil)
		params = append(params, perm)
	}
	return params
}

func buildUpdateAccountPasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password": utils.ValueIgnoreEmpty(d.Get("password")),
	}
	return bodyParams
}

func resourceDdmAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAccount: Query DDM account
	var (
		getAccountHttpUrl = "v1/{project_id}/instances/{instance_id}/users"
		getAccountProduct = "ddm"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<account_name>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{project_id}", getAccountClient.ProjectID)
	getAccountPath = strings.ReplaceAll(getAccountPath, "{instance_id}", fmt.Sprintf("%v", instanceID))

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

	accounts := utils.PathSearch("users", getAccountRespBody, nil)
	if accounts == nil {
		log.Printf("[WARN] failed to get DDM account, user %s is not found in API response", accountName)
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	for _, account := range accounts.([]interface{}) {
		name := utils.PathSearch("name", account, nil)
		if accountName != name {
			continue
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

	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
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
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<account_name>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	deleteAccountPath := deleteAccountClient.Endpoint + deleteAccountHttpUrl
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{project_id}", deleteAccountClient.ProjectID)
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{instance_id}", fmt.Sprintf("%v", instanceID))
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{username}", fmt.Sprintf("%v", accountName))

	deleteAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = deleteAccountClient.Request("DELETE", deleteAccountPath, &deleteAccountOpt)
		isRetry, err := handleOperationError(err, "deleting", "account")
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
