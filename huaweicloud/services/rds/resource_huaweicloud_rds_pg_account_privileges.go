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

var pgAccountPrivilegesNonUpdatableParams = []string{"instance_id", "user_name"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/db-user-privilege
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/db_user/detail
func ResourcePgAccountPrivileges() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgAccountPrivilegesCreate,
		UpdateContext: resourcePgAccountPrivilegesUpdate,
		ReadContext:   resourcePgAccountPrivilegesRead,
		DeleteContext: resourcePgAccountPrivilegesDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePgAccountPrivilegesImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(pgAccountPrivilegesNonUpdatableParams),

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
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the username of the account.`,
			},
			"role_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of role privileges.`,
			},
			"system_role_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of system role privileges.`,
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

func resourcePgAccountPrivilegesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	username := d.Get("user_name").(string)
	rolePrivileges := d.Get("role_privileges").(*schema.Set)
	systemRolePrivileges := d.Get("system_role_privileges").(*schema.Set)
	if rolePrivileges.Len() > 0 {
		requestBody := buildUpdatePgAccountPrivilegesBodyParams(username, "ROLE", rolePrivileges.List())
		err = updateAccountPrivileges(ctx, d, client, schema.TimeoutCreate, requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if rolePrivileges.Len() > 0 {
		requestBody := buildUpdatePgAccountPrivilegesBodyParams(username, "SYSTEM_ROLE", systemRolePrivileges.List())
		err = updateAccountPrivileges(ctx, d, client, schema.TimeoutCreate, requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("instance_id").(string), username))

	return resourcePgAccountPrivilegesRead(ctx, d, meta)
}

func resourcePgAccountPrivilegesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getPgAccountResp, err := pagination.ListAllItems(
		client,
		"page",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS PostgreSQL account privileges")
	}

	respJson, err := json.Marshal(getPgAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	username := d.Get("user_name").(string)
	attributes := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0].attributes", username), respBody, nil)

	if attributes == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	rolePrivileges := make([]string, 0)
	if utils.PathSearch("rolcreatedb", attributes, false).(bool) {
		rolePrivileges = append(rolePrivileges, "CREATEDB")
	}
	if utils.PathSearch("rolcreaterole", attributes, false).(bool) {
		rolePrivileges = append(rolePrivileges, "CREATEROLE")
	}
	if utils.PathSearch("rolreplication", attributes, false).(bool) {
		rolePrivileges = append(rolePrivileges, "REPLICATION")
	}
	if utils.PathSearch("rolcanlogin", attributes, false).(bool) {
		rolePrivileges = append(rolePrivileges, "LOGIN")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("user_name", username),
		d.Set("role_privileges", rolePrivileges),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePgAccountPrivilegesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	username := d.Get("user_name").(string)
	oldRaws, newRaws := d.GetChange("role_privileges")
	addSet := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	deleteSet := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	if deleteSet.Len() > 0 {
		deleteRolePrivileges := buildDeleteRolePrivileges(deleteSet.List())
		requestBody := buildUpdatePgAccountPrivilegesBodyParams(username, "RECYCLING_ROLE", deleteRolePrivileges)
		err = updateAccountPrivileges(ctx, d, client, schema.TimeoutUpdate, requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		requestBody := buildUpdatePgAccountPrivilegesBodyParams(username, "ROLE", addSet.List())
		err = updateAccountPrivileges(ctx, d, client, schema.TimeoutUpdate, requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	oldRaws, newRaws = d.GetChange("system_role_privileges")
	addSet = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	deleteSet = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	if deleteSet.Len() > 0 {
		requestBody := buildUpdatePgAccountPrivilegesBodyParams(username, "RECYCLING_SYSTEM_ROLE", deleteSet.List())
		err = updateAccountPrivileges(ctx, d, client, schema.TimeoutUpdate, requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		requestBody := buildUpdatePgAccountPrivilegesBodyParams(username, "SYSTEM_ROLE", addSet.List())
		err = updateAccountPrivileges(ctx, d, client, schema.TimeoutUpdate, requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePgAccountPrivilegesRead(ctx, d, meta)
}

func resourcePgAccountPrivilegesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	username := d.Get("user_name").(string)
	rolePrivileges := d.Get("role_privileges").(*schema.Set)
	systemRolePrivileges := d.Get("system_role_privileges").(*schema.Set)
	if rolePrivileges.Len() > 0 {
		deleteRolePrivileges := buildDeleteRolePrivileges(rolePrivileges.List())
		requestBody := buildUpdatePgAccountPrivilegesBodyParams(username, "RECYCLING_ROLE", deleteRolePrivileges)
		err = updateAccountPrivileges(ctx, d, client, schema.TimeoutDelete, requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if rolePrivileges.Len() > 0 {
		requestBody := buildUpdatePgAccountPrivilegesBodyParams(username, "RECYCLING_SYSTEM_ROLE", systemRolePrivileges.List())
		err = updateAccountPrivileges(ctx, d, client, schema.TimeoutDelete, requestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func buildDeleteRolePrivileges(privileges []interface{}) []interface{} {
	res := make([]interface{}, 0, len(privileges))
	for _, v := range privileges {
		res = append(res, fmt.Sprintf("NO%s", v.(string)))
	}
	return res
}

func buildUpdatePgAccountPrivilegesBodyParams(username, authorizationType string, privileges []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_name":          username,
		"authorization_type": authorizationType,
		"privileges":         privileges,
	}
	return bodyParams
}

func updateAccountPrivileges(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string,
	requestBody map[string]interface{}) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/db-user-privilege"

	instanceId := d.Get("instance_id").(string)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("POST", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating PostgreSQL account privileges: %s", err)
	}
	return nil
}

func resourcePgAccountPrivilegesImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<user_name>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("user_name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
