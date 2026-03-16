package rds

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
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

var (
	mysqlDatabasePrivilegeNonUpdatableParams = []string{
		"instance_id",
		"db_name",
	}
	objSliceParamKeysForMysqlDatabasePrivilege = []string{
		"users",
	}
)

// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db_privilege
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_privilege
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/db_user
func ResourceMysqlDatabasePrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMysqlDatabasePrivilegeCreateAndUpdate,
		ReadContext:   resourceMysqlDatabasePrivilegeRead,
		UpdateContext: resourceMysqlDatabasePrivilegeCreateAndUpdate,
		DeleteContext: resourceMysqlDatabasePrivilegeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: mysqlDatabasePrivilegeImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(mysqlDatabasePrivilegeNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The region where the database and users (accounts) are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the MySQL instance.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The database name to which the users (accounts) are privileged.`,
			},
			"users": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The username of the database account.`,
						},
						"readonly": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether the user has read-only permission.`,
						},
					},
				},
				DiffSuppressFunc: utils.SuppressObjectSliceDiffs(),
				Description:      `The user (account) permissions with the database.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Internal attributes.
			"users_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the username of the database account.`,
						},
						"readonly": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Specifies the read-only permission.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'users'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func findDeleteUserPrivilegesFromDatabase(oirginUsers, rawConfigUsers []interface{}) []interface{} {
	if len(oirginUsers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(oirginUsers))
	for _, user := range oirginUsers {
		if utils.PathSearch(fmt.Sprintf("length([?name == '%v'])", utils.PathSearch("name", user, "")), rawConfigUsers, float64(0)).(float64) < 1 {
			// If the new user list does not contain this user, it is considered that this user is no longer privileged.
			result = append(result, user)
		} else if !utils.PathSearch(fmt.Sprintf("[?name=='%v']|[0].readonly == `%v`",
			utils.PathSearch("name", user, ""), utils.PathSearch("readonly", user, "")), rawConfigUsers, false).(bool) {
			// If the read-only permission of this user is different from the new user list, it is considered that this user needs to be updated.
			result = append(result, user)
		}
	}
	return result
}

func buildDeleteUserPrivilegesFromDatabaseRequestBodyUsers(users []interface{}) []interface{} {
	result := make([]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(utils.PathSearch("name", user, nil)),
		})
	}
	return result
}

func buildDeleteUserPrivilegesFromDatabaseBodyParams(dbName string, users []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": dbName,
		"users":   buildDeleteUserPrivilegesFromDatabaseRequestBodyUsers(users),
	}
	return bodyParams
}

func deleteMysqlDatabasePrivilege(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	users []interface{}, timeout time.Duration) error {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/db_privilege"
		instanceId = d.Get("instance_id").(string)
		dbName     = d.Get("db_name").(string)

		start = 0
		// A single request supports a maximum of 50 users.
		end = int(math.Min(50, float64(len(users))))
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)

	for start < end {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         utils.RemoveNil(buildDeleteUserPrivilegesFromDatabaseBodyParams(dbName, users[start:end])),
		}

		retryFunc := func() (interface{}, bool, error) {
			_, err := client.Request("DELETE", deletePath, &opt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      timeout,
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error deleting user privileges from MySQL database: %s", err)
		}
		start += 50
		end = int(math.Min(float64(end+50), float64(len(users))))
	}
	return nil
}

func findAddUserPrivilegesToDatabase(rawConfigUsers, remoteStateUsers []interface{}) []interface{} {
	if len(rawConfigUsers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(rawConfigUsers))
	for _, user := range rawConfigUsers {
		if utils.PathSearch(fmt.Sprintf("length([?name == '%v'])", utils.PathSearch("name", user, "")), remoteStateUsers, float64(0)).(float64) < 1 {
			// If the old user list does not contain this user, it is considered that this user is a newly privileged user.
			result = append(result, user)
		} else if !utils.PathSearch(fmt.Sprintf("[?name=='%v']|[0].readonly == `%v`",
			utils.PathSearch("name", user, ""), utils.PathSearch("readonly", user, "")), remoteStateUsers, false).(bool) {
			// If the read-only permission of this user is different from the old user list, it is considered that this user needs to be updated.
			result = append(result, user)
		}
	}
	return result
}

func buildAddUserPrivilegesToDatabaseRequestBodyUsers(users []interface{}) []interface{} {
	result := make([]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"name":     utils.ValueIgnoreEmpty(utils.PathSearch("name", user, nil)),
			"readonly": utils.ValueIgnoreEmpty(utils.PathSearch("readonly", user, nil)),
		})
	}
	return result
}

func buildAddUserPrivilegesToDatabaseBodyParams(dbName string, users []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": dbName,
		"users":   buildAddUserPrivilegesToDatabaseRequestBodyUsers(users),
	}
	return bodyParams
}

func addPrivilegesToDatabase(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	users []interface{}, timeout time.Duration) error {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/db_privilege"
		instanceId = d.Get("instance_id").(string)
		dbName     = d.Get("db_name").(string)

		start = 0
		// A single request supports a maximum of 50 users.
		end = int(math.Min(50, float64(len(users))))
	)

	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{project_id}", client.ProjectID)
	addPath = strings.ReplaceAll(addPath, "{instance_id}", instanceId)

	for start < end {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         utils.RemoveNil(buildAddUserPrivilegesToDatabaseBodyParams(dbName, users[start:end])),
		}

		retryFunc := func() (interface{}, bool, error) {
			_, err := client.Request("POST", addPath, &opt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      timeout,
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error creating user privileges for MySQL database: %s", err)
		}
		start += 50
		end = int(math.Min(float64(end+50), float64(len(users))))
	}
	return nil
}

func resourceMysqlDatabasePrivilegeCreateAndUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		rawConfigUsers = utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "users").([]interface{})
		originUsers    = d.Get("users_origin").([]interface{})
		instanceId     = d.Get("instance_id").(string)
		dbName         = d.Get("db_name").(string)
	)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(fmt.Sprintf("%s/%s", instanceId, dbName))
	}

	remoteStateUsers, err := ListMysqlDatabasePrivileges(client, instanceId, dbName, originUsers, true)
	if err != nil {
		return diag.Errorf("error getting remoteMySQL database privileges before update: %s", err)
	}

	deleteUsers := findDeleteUserPrivilegesFromDatabase(originUsers, rawConfigUsers)
	addUsers := findAddUserPrivilegesToDatabase(rawConfigUsers, remoteStateUsers)

	if len(deleteUsers) > 0 {
		err = deleteMysqlDatabasePrivilege(ctx, client, d, deleteUsers, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(addUsers) > 0 {
		err = addPrivilegesToDatabase(ctx, client, d, addUsers, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshObjectParamOriginValues(d, objSliceParamKeysForMysqlDatabasePrivilege)
	if err != nil {
		// Don't report an error if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}
	return resourceMysqlDatabasePrivilegeRead(ctx, d, meta)
}

func orderMysqlDatabasePrivilegesByUsersOrigin(databasePrivileges, usersOrigin []interface{}) []interface{} {
	if len(usersOrigin) < 1 {
		return databasePrivileges
	}

	sortedDatabasePrivileges := make([]interface{}, 0, len(databasePrivileges))
	databasePrivilegesCopy := databasePrivileges
	for _, userOrigin := range usersOrigin {
		userNameOrigin := utils.PathSearch("name", userOrigin, "").(string)
		for index, databasePrivilege := range databasePrivilegesCopy {
			if utils.PathSearch("name", databasePrivilege, "").(string) == userNameOrigin {
				// Add the found database privilege to the sorted database privileges list.
				sortedDatabasePrivileges = append(sortedDatabasePrivileges, databasePrivilegesCopy[index])
				// Remove the processed database privilege from the original array.
				databasePrivilegesCopy = append(databasePrivilegesCopy[:index], databasePrivilegesCopy[index+1:]...)
				break
			}
		}
	}

	return sortedDatabasePrivileges
}

func ListMysqlDatabasePrivileges(client *golangsdk.ServiceClient, instanceId, dbName string, usersOrigin []interface{},
	ignoreNotFound ...bool) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/database/db_user?db-name={db_name}&limit={limit}"
		limit   = 100
		page    = 1
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{db_name}", dbName)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	result := make([]interface{}, 0)
	for {
		listPathWithPage := fmt.Sprintf("%s&page=%d", listPath, page)
		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		users := utils.PathSearch("users", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, users...)
		if len(users) < limit {
			break
		}
		page++
	}

	parsedUsers := orderMysqlDatabasePrivilegesByUsersOrigin(result, usersOrigin)
	if len(parsedUsers) < 1 && (len(ignoreNotFound) < 1 || !ignoreNotFound[0]) {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/instances/{instance_id}/database/db_user",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the database privileges (%#v) for database (%s) do not exist", usersOrigin, dbName)),
			},
		}
	}
	return parsedUsers, nil
}

func flattenMysqlDatabasePrivilegeUsers(users []interface{}) []interface{} {
	if len(users) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"name":     utils.PathSearch("name", user, nil),
			"readonly": utils.PathSearch("readonly", user, nil),
		})
	}
	return result
}

func resourceMysqlDatabasePrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		instanceId  = d.Get("instance_id").(string)
		dbName      = d.Get("db_name").(string)
		usersOrigin = d.Get("users_origin").([]interface{})
	)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	users, err := ListMysqlDatabasePrivileges(client, instanceId, dbName, usersOrigin)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("users", flattenMysqlDatabasePrivilegeUsers(users)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMysqlDatabasePrivilegeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		users  = d.Get("users").([]interface{})
	)

	// The value of users_origin is empty only when the resource is imported and the terraform apply command is not executed.
	// In this case, all information obtained from the remote service is used to remove user relationships from the database.
	if originUsers, ok := d.GetOk("users_origin"); ok && len(originUsers.([]interface{})) > 0 {
		log.Printf("[DEBUG] Find the custom users configuration, according to it to remove users from the database (%v)", d.Id())
		users = originUsers.([]interface{})
	}

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	err = deleteMysqlDatabasePrivilege(ctx, client, d, users, d.Timeout(schema.TimeoutDelete))
	return diag.FromErr(err)
}

func mysqlDatabasePrivilegeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resource ID format for privilege management, want '<instance_id>/<db_name>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("db_name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
