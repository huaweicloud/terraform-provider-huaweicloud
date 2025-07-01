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

var mysqlAccountNonUpdatableParams = []string{"instance_id", "name", "hosts"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/db-users/{user_name}/comment
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_user
// @API RDS GET /v3/{project_id}/instances
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db_user/{user_name}
// @API RDS GET /v3/{project_id}/instances/{instance_id}/db_user/detail
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_user/resetpwd
func ResourceMysqlAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMysqlAccountCreate,
		UpdateContext: resourceMysqlAccountUpdate,
		ReadContext:   resourceMysqlAccountRead,
		DeleteContext: resourceMysqlAccountDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(mysqlAccountNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS Mysql instance.`,
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
			"hosts": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the IP addresses that are allowed to access your DB instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies remarks of the DB account.`,
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

func resourceMysqlAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createMysqlAccount: create RDS Mysql account.
	var (
		createMysqlAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user"
		createMysqlAccountProduct = "rds"
	)
	createMysqlAccountClient, err := cfg.NewServiceClient(createMysqlAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createMysqlAccountPath := createMysqlAccountClient.Endpoint + createMysqlAccountHttpUrl
	createMysqlAccountPath = strings.ReplaceAll(createMysqlAccountPath, "{project_id}",
		createMysqlAccountClient.ProjectID)
	createMysqlAccountPath = strings.ReplaceAll(createMysqlAccountPath, "{instance_id}", instanceId)

	createMysqlAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestBody := buildCreateMysqlAccountBodyParams(d)
	log.Printf("[DEBUG] Create RDS Mysql account options: %#v", createMysqlAccountOpt)
	requestBody["password"] = utils.ValueIgnoreEmpty(d.Get("password"))
	createMysqlAccountOpt.JSONBody = utils.RemoveNil(requestBody)

	retryFunc := func() (interface{}, bool, error) {
		_, err = createMysqlAccountClient.Request("POST", createMysqlAccountPath, &createMysqlAccountOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(createMysqlAccountClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS Mysql account: %s", err)
	}

	accountName := d.Get("name").(string)
	d.SetId(instanceId + "/" + accountName)

	return resourceMysqlAccountRead(ctx, d, meta)
}

func buildCreateMysqlAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":    utils.ValueIgnoreEmpty(d.Get("name")),
		"comment": utils.ValueIgnoreEmpty(d.Get("description")),
		"hosts":   utils.ValueIgnoreEmpty(d.Get("hosts")),
	}
	return bodyParams
}

func resourceMysqlAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getMysqlAccount: query RDS Mysql account
	var (
		getMysqlAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		getMysqlAccountProduct = "rds"
	)
	getMysqlAccountClient, err := cfg.NewServiceClient(getMysqlAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getMysqlAccountPath := getMysqlAccountClient.Endpoint + getMysqlAccountHttpUrl
	getMysqlAccountPath = strings.ReplaceAll(getMysqlAccountPath, "{project_id}", getMysqlAccountClient.ProjectID)
	getMysqlAccountPath = strings.ReplaceAll(getMysqlAccountPath, "{instance_id}", instanceId)

	getMysqlAccountResp, err := pagination.ListAllItems(
		getMysqlAccountClient,
		"page",
		getMysqlAccountPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS Mysql account")
	}

	getMysqlAccountRespJson, err := json.Marshal(getMysqlAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getMysqlAccountRespBody interface{}
	err = json.Unmarshal(getMysqlAccountRespJson, &getMysqlAccountRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getMysqlAccountRespBody, nil)

	if account == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", utils.PathSearch("name", account, nil)),
		d.Set("description", utils.PathSearch("comment", account, nil)),
		d.Set("hosts", utils.PathSearch("hosts", account, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMysqlAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateMysqlAccountClient, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if err = updateMysqlAccountPassword(ctx, d, updateMysqlAccountClient); err != nil {
		return diag.FromErr(err)
	}

	if err = updateMysqlAccountDescription(d, updateMysqlAccountClient); err != nil {
		return diag.FromErr(err)
	}
	return resourceMysqlAccountRead(ctx, d, meta)
}

func updateMysqlAccountPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("password") {
		return nil
	}
	// updateMysqlAccount: update RDS Mysql account password
	updateMysqlAccountHttpUrl := "v3/{project_id}/instances/{instance_id}/db_user/resetpwd"

	updateMysqlAccountPath := client.Endpoint + updateMysqlAccountHttpUrl
	updateMysqlAccountPath = strings.ReplaceAll(updateMysqlAccountPath, "{project_id}", client.ProjectID)
	updateMysqlAccountPath = strings.ReplaceAll(updateMysqlAccountPath, "{instance_id}",
		d.Get("instance_id").(string))

	updateMysqlAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateMysqlAccountOpt.JSONBody = utils.RemoveNil(buildUpdateMysqlAccountPasswordBodyParams(d))

	instanceId := d.Get("instance_id").(string)
	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("POST", updateMysqlAccountPath, &updateMysqlAccountOpt)
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
		return fmt.Errorf("error updating Mysql account password: %s", err)
	}
	return nil
}

func updateMysqlAccountDescription(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("description") {
		return nil
	}
	// updateMysqlAccount: update RDS Mysql account password
	updateMysqlAccountHttpUrl := "v3/{project_id}/instances/{instance_id}/db-users/{user_name}/comment"

	updateMysqlAccountPath := client.Endpoint + updateMysqlAccountHttpUrl
	updateMysqlAccountPath = strings.ReplaceAll(updateMysqlAccountPath, "{project_id}", client.ProjectID)
	updateMysqlAccountPath = strings.ReplaceAll(updateMysqlAccountPath, "{instance_id}",
		d.Get("instance_id").(string))
	updateMysqlAccountPath = strings.ReplaceAll(updateMysqlAccountPath, "{user_name}", d.Get("name").(string))

	updateMysqlAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateMysqlAccountOpt.JSONBody = utils.RemoveNil(buildUpdateMysqlAccountDescriptionBodyParams(d))

	log.Printf("[DEBUG] Update RDS Mysql account description options: %#v", updateMysqlAccountOpt)
	_, err := client.Request("PUT", updateMysqlAccountPath, &updateMysqlAccountOpt)
	if err != nil {
		return fmt.Errorf("error updating Mysql account description: %s", err)
	}
	return nil
}

func buildUpdateMysqlAccountPasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(d.Get("name")),
		"password": utils.ValueIgnoreEmpty(d.Get("password")),
	}
	return bodyParams
}

func buildUpdateMysqlAccountDescriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"comment": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceMysqlAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteMysqlAccount: delete RDS Mysql account
	var (
		deleteMysqlAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/{user_name}"
		deleteMysqlAccountProduct = "rds"
	)
	deleteMysqlAccountClient, err := cfg.NewServiceClient(deleteMysqlAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	deleteMysqlAccountPath := deleteMysqlAccountClient.Endpoint + deleteMysqlAccountHttpUrl
	deleteMysqlAccountPath = strings.ReplaceAll(deleteMysqlAccountPath, "{project_id}",
		deleteMysqlAccountClient.ProjectID)
	deleteMysqlAccountPath = strings.ReplaceAll(deleteMysqlAccountPath, "{instance_id}",
		d.Get("instance_id").(string))
	deleteMysqlAccountPath = strings.ReplaceAll(deleteMysqlAccountPath, "{user_name}", d.Get("name").(string))

	deleteMysqlAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	log.Printf("[DEBUG] Delete RDS Mysql account options: %#v", deleteMysqlAccountOpt)

	instanceId := d.Get("instance_id").(string)
	retryFunc := func() (interface{}, bool, error) {
		_, err = deleteMysqlAccountClient.Request("DELETE", deleteMysqlAccountPath, &deleteMysqlAccountOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deleteMysqlAccountClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting RDS Mysql account: %s", err)
	}

	return nil
}
