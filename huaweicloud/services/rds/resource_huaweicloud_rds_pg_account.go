// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

var pgAccountNonUpdatableParams = []string{"instance_id", "name"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/db-users/{user_name}/comment
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_user
// @API RDS GET /v3/{project_id}/instances
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db_user/{user_name}
// @API RDS GET /v3/{project_id}/instances/{instance_id}/db_user/detail
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_user/resetpwd
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db-user-role
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db-user-role
func ResourcePgAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgAccountCreate,
		UpdateContext: resourcePgAccountUpdate,
		ReadContext:   resourcePgAccountRead,
		DeleteContext: resourcePgAccountDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(pgAccountNonUpdatableParams),

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
				Description: `Specifies the ID of the RDS PostgreSQL instance.`,
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the remarks of the DB account.`,
			},
			"attributes": {
				Type:        schema.TypeList,
				Elem:        pgAccountAttributesSchema(),
				Computed:    true,
				Description: `Indicates the permission attributes of the account.`,
			},
			"memberof": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `schema: Deprecated`,
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

func pgAccountAttributesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rol_super": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether a user has the super-user permission.`,
			},
			"rol_inherit": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `Indicates whether a user automatically inherits the permissions of the role to which the
user belongs.`,
			},
			"rol_create_role": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether a user can create other sub-users.`,
			},
			"rol_create_db": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether a user can create a database.`,
			},
			"rol_can_login": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether a user can log in to the database.`,
			},
			"rol_conn_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of concurrent connections to a DB instance.`,
			},
			"rol_replication": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user is a replication role.`,
			},
			"rol_bypass_rls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether a user bypasses each row-level security policy.`,
			},
		},
	}
	return &sc
}

func resourcePgAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPgAccount: create RDS PostgreSQL account.
	var (
		createPgAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user"
		createPgAccountProduct = "rds"
	)
	createPgAccountClient, err := cfg.NewServiceClient(createPgAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPgAccountPath := createPgAccountClient.Endpoint + createPgAccountHttpUrl
	createPgAccountPath = strings.ReplaceAll(createPgAccountPath, "{project_id}", createPgAccountClient.ProjectID)
	createPgAccountPath = strings.ReplaceAll(createPgAccountPath, "{instance_id}", instanceId)

	createPgAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestBody := buildCreatePgAccountBodyParams(d)
	log.Printf("[DEBUG] Create RDS PostgreSQL account options: %#v", requestBody)
	requestBody["password"] = d.Get("password")
	createPgAccountOpt.JSONBody = utils.RemoveNil(requestBody)

	retryFunc := func() (interface{}, bool, error) {
		_, err = createPgAccountClient.Request("POST", createPgAccountPath, &createPgAccountOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(createPgAccountClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS PostgreSQL account: %s", err)
	}

	accountName := d.Get("name").(string)
	d.SetId(instanceId + "/" + accountName)

	if err = updatePgAccountMemberOf(ctx, d, createPgAccountClient); err != nil {
		return diag.FromErr(err)
	}

	return resourcePgAccountRead(ctx, d, meta)
}

func buildCreatePgAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":    d.Get("name"),
		"comment": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourcePgAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPgAccount: query RDS PostgreSQL account
	var (
		getPgAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		getPgAccountProduct = "rds"
	)
	getPgAccountClient, err := cfg.NewServiceClient(getPgAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getPgAccountPath := getPgAccountClient.Endpoint + getPgAccountHttpUrl
	getPgAccountPath = strings.ReplaceAll(getPgAccountPath, "{project_id}", getPgAccountClient.ProjectID)
	getPgAccountPath = strings.ReplaceAll(getPgAccountPath, "{instance_id}", instanceId)

	getPgAccountResp, err := pagination.ListAllItems(
		getPgAccountClient,
		"page",
		getPgAccountPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS PostgreSQL account")
	}

	getPgAccountRespJson, err := json.Marshal(getPgAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getPgAccountRespBody interface{}
	err = json.Unmarshal(getPgAccountRespJson, &getPgAccountRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getPgAccountRespBody, nil)

	if account == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", utils.PathSearch("name", account, nil)),
		d.Set("description", utils.PathSearch("comment", account, nil)),
		d.Set("attributes", flattenPgAccountAttributesBody(account)),
		d.Set("memberof", utils.PathSearch("memberof", account, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPgAccountAttributesBody(account interface{}) []interface{} {
	attributes := utils.PathSearch("attributes", account, nil)
	if attributes == nil {
		return nil
	}
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"rol_super":       utils.PathSearch("rolsuper", attributes, nil),
		"rol_inherit":     utils.PathSearch("rolinherit", attributes, nil),
		"rol_create_role": utils.PathSearch("rolcreaterole", attributes, nil),
		"rol_create_db":   utils.PathSearch("rolcreatedb", attributes, nil),
		"rol_can_login":   utils.PathSearch("rolcanlogin", attributes, nil),
		"rol_conn_limit":  utils.PathSearch("rolconnlimit", attributes, nil),
		"rol_replication": utils.PathSearch("rolreplication", attributes, nil),
		"rol_bypass_rls":  utils.PathSearch("rolbypassrls", attributes, nil),
	})
	return rst
}

func resourcePgAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updatePgAccountClient, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if err = updatePgAccountMemberOf(ctx, d, updatePgAccountClient); err != nil {
		return diag.FromErr(err)
	}

	if err = updatePgAccountPassword(ctx, d, updatePgAccountClient); err != nil {
		return diag.FromErr(err)
	}

	if err = updatePgAccountDescription(d, updatePgAccountClient); err != nil {
		return diag.FromErr(err)
	}
	return resourcePgAccountRead(ctx, d, meta)
}

func updatePgAccountMemberOf(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("memberof") {
		return nil
	}

	oldRaws, newRaws := d.GetChange("memberof")
	addMemberOf := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	deleteMemberOf := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))

	if deleteMemberOf.Len() > 0 {
		requestBody := buildUpdatePgAccountMemberOfBodyParams(d.Get("name").(string), deleteMemberOf.List())
		err := updateMemberOf(ctx, d, client, "DELETE", requestBody)
		if err != nil {
			return err
		}
	}

	if addMemberOf.Len() > 0 {
		requestBody := buildUpdatePgAccountMemberOfBodyParams(d.Get("name").(string), addMemberOf.List())
		err := updateMemberOf(ctx, d, client, "POST", requestBody)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildUpdatePgAccountMemberOfBodyParams(user string, memberOf []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user":  user,
		"roles": memberOf,
	}
	return bodyParams
}

func updateMemberOf(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, httpMethod string,
	requestBody map[string]interface{}) error {
	// updatePgAccount: update RDS PostgreSQL account memberOf
	updatePgAccountMemberOfHttpUrl := "v3/{project_id}/instances/{instance_id}/db-user-role"

	instanceId := d.Get("instance_id").(string)
	updatePgAccountMemberOfPath := client.Endpoint + updatePgAccountMemberOfHttpUrl
	updatePgAccountMemberOfPath = strings.ReplaceAll(updatePgAccountMemberOfPath, "{project_id}", client.ProjectID)
	updatePgAccountMemberOfPath = strings.ReplaceAll(updatePgAccountMemberOfPath, "{instance_id}", instanceId)

	updatePgAccountMemberOfOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request(httpMethod, updatePgAccountMemberOfPath, &updatePgAccountMemberOfOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating PostgreSQL account member: %s", err)
	}
	return nil
}

func updatePgAccountPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("password") {
		return nil
	}
	// updatePgAccount: update RDS PostgreSQL account password
	updatePgAccountHttpUrl := "v3/{project_id}/instances/{instance_id}/db_user/resetpwd"

	instanceId := d.Get("instance_id").(string)
	updatePgAccountPath := client.Endpoint + updatePgAccountHttpUrl
	updatePgAccountPath = strings.ReplaceAll(updatePgAccountPath, "{project_id}", client.ProjectID)
	updatePgAccountPath = strings.ReplaceAll(updatePgAccountPath, "{instance_id}", instanceId)

	updatePgAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updatePgAccountOpt.JSONBody = utils.RemoveNil(buildUpdatePgAccountPasswordBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("POST", updatePgAccountPath, &updatePgAccountOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating PostgreSQL account password: %s", err)
	}
	return nil
}

func updatePgAccountDescription(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("description") {
		return nil
	}
	// updatePglAccount: update RDS PostgreSQL account password
	updatePgAccountHttpUrl := "v3/{project_id}/instances/{instance_id}/db-users/{user_name}/comment"

	updatePgAccountPath := client.Endpoint + updatePgAccountHttpUrl
	updatePgAccountPath = strings.ReplaceAll(updatePgAccountPath, "{project_id}", client.ProjectID)
	updatePgAccountPath = strings.ReplaceAll(updatePgAccountPath, "{instance_id}",
		d.Get("instance_id").(string))
	updatePgAccountPath = strings.ReplaceAll(updatePgAccountPath, "{user_name}", d.Get("name").(string))

	updatePgAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestBody := buildUpdatePgAccountDescriptionBodyParams(d)
	log.Printf("[DEBUG] Update RDS PostgreSQL account description options: %#v", requestBody)
	updatePgAccountOpt.JSONBody = utils.RemoveNil(requestBody)

	_, err := client.Request("PUT", updatePgAccountPath, &updatePgAccountOpt)
	if err != nil {
		return fmt.Errorf("error updating PostgreSQL account description: %s", err)
	}
	return nil
}

func buildUpdatePgAccountPasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     d.Get("name"),
		"password": d.Get("password"),
	}
	return bodyParams
}

func buildUpdatePgAccountDescriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"comment": d.Get("description"),
	}
	return bodyParams
}

func resourcePgAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deletePgAccount: delete RDS PostgreSQL account
	var (
		deletePgAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/{user_name}"
		deletePgAccountProduct = "rds"
	)
	deletePgAccountClient, err := cfg.NewServiceClient(deletePgAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deletePgAccountPath := deletePgAccountClient.Endpoint + deletePgAccountHttpUrl
	deletePgAccountPath = strings.ReplaceAll(deletePgAccountPath, "{project_id}", deletePgAccountClient.ProjectID)
	deletePgAccountPath = strings.ReplaceAll(deletePgAccountPath, "{instance_id}", instanceId)
	deletePgAccountPath = strings.ReplaceAll(deletePgAccountPath, "{user_name}", d.Get("name").(string))

	deletePgAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err = deletePgAccountClient.Request("DELETE", deletePgAccountPath, &deletePgAccountOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deletePgAccountClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting RDS PostgreSQL account: %s", err)
	}

	return nil
}
