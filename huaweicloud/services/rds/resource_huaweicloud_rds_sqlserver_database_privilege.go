// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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

const maxElementsPerRequest = 50

var sqlServerDatabasePrivilegeNonUpdatableParams = []string{"instance_id", "db_name"}

// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/db_user
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db_privilege
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_privilege
// @API RDS GET /v3/{project_id}/instances
func ResourceSQLServerDatabasePrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSQLServerDatabasePrivilegeCreate,
		UpdateContext: resourceSQLServerDatabasePrivilegeUpdate,
		ReadContext:   resourceSQLServerDatabasePrivilegeRead,
		DeleteContext: resourceSQLServerDatabasePrivilegeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(sqlServerDatabasePrivilegeNonUpdatableParams),

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
				Description: `Specifies the ID of the RDS SQL Server instance.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"users": {
				Type:        schema.TypeSet,
				Elem:        sQLServerDatabasePrivilegeCreateUserSchema(),
				Required:    true,
				Description: `Specifies the account that associated with the database`,
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

func sQLServerDatabasePrivilegeCreateUserSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the username of the database account.`,
			},
			"readonly": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the read-only permission.`,
			},
		},
	}
	return &sc
}

func resourceSQLServerDatabasePrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSQLServerDatabasePrivilege: create RDS SQL Server database privilege.
	var (
		createSQLServerDatabasePrivilegeProduct = "rds"
	)
	createSQLServerDatabasePrivilegeClient, err := cfg.NewServiceClient(createSQLServerDatabasePrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if err = createSQLServerDatabasePrivilege(ctx, d, schema.TimeoutCreate, createSQLServerDatabasePrivilegeClient,
		d.Get("users").(*schema.Set).List()); err != nil {
		return diag.FromErr(err)
	}

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	d.SetId(instanceId + "/" + dbName)

	return resourceSQLServerDatabasePrivilegeRead(ctx, d, meta)
}

func resourceSQLServerDatabasePrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSQLServerDatabasePrivilege: query RDS SQL Server database privilege
	var (
		getSQLServerDatabasePrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/database/db_user"
		getSQLServerDatabasePrivilegeProduct = "rds"
	)
	getSQLServerDatabasePrivilegeClient, err := cfg.NewServiceClient(getSQLServerDatabasePrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<db_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getSQLServerDatabasePrivilegePath := getSQLServerDatabasePrivilegeClient.Endpoint + getSQLServerDatabasePrivilegeHttpUrl
	getSQLServerDatabasePrivilegePath = strings.ReplaceAll(getSQLServerDatabasePrivilegePath, "{project_id}",
		getSQLServerDatabasePrivilegeClient.ProjectID)
	getSQLServerDatabasePrivilegePath = strings.ReplaceAll(getSQLServerDatabasePrivilegePath, "{instance_id}", instanceId)

	getSQLServerDatabasePrivilegeQueryParams := buildGetSQLServerDatabasePrivilegeQueryParams(dbName)
	getSQLServerDatabasePrivilegePath += getSQLServerDatabasePrivilegeQueryParams

	getSQLServerDatabasePrivilegeResp, err := pagination.ListAllItems(
		getSQLServerDatabasePrivilegeClient,
		"page",
		getSQLServerDatabasePrivilegePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS SQL Server database privilege")
	}

	getSQLServerDatabasePrivilegeRespJson, err := json.Marshal(getSQLServerDatabasePrivilegeResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getSQLServerDatabasePrivilegeRespBody interface{}
	err = json.Unmarshal(getSQLServerDatabasePrivilegeRespJson, &getSQLServerDatabasePrivilegeRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// 'rdsuser' is system account, it can not be managed
	users := utils.PathSearch("users[?name != 'rdsuser']", getSQLServerDatabasePrivilegeRespBody,
		make([]interface{}, 0)).([]interface{})
	if len(users) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("db_name", dbName),
		d.Set("users", users),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSQLServerDatabasePrivilegeQueryParams(dbName string) string {
	return fmt.Sprintf("?db-name=%s&page=1&limit=100", dbName)
}

func resourceSQLServerDatabasePrivilegeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateSQLServerDatabasePrivilege: update RDS SQL Server database privilege.
	var (
		updateSQLServerDatabasePrivilegeProduct = "rds"
	)
	updateSQLServerDatabasePrivilegeClient, err := cfg.NewServiceClient(updateSQLServerDatabasePrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	oldRaws, newRaws := d.GetChange("users")
	createUsers := newRaws.(*schema.Set).List()
	deleteUsers := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List()

	if len(deleteUsers) > 0 {
		if err = deleteSQLServerDatabasePrivilege(ctx, d, schema.TimeoutUpdate, updateSQLServerDatabasePrivilegeClient,
			deleteUsers); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(createUsers) > 0 {
		if err = createSQLServerDatabasePrivilege(ctx, d, schema.TimeoutUpdate, updateSQLServerDatabasePrivilegeClient,
			createUsers); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSQLServerDatabasePrivilegeRead(ctx, d, meta)
}

func resourceSQLServerDatabasePrivilegeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSQLServerDatabasePrivilege: delete RDS SQL Server database privilege
	var (
		deleteSQLServerDatabasePrivilegeProduct = "rds"
	)
	deleteSQLServerDatabasePrivilegeClient, err := cfg.NewServiceClient(deleteSQLServerDatabasePrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if err = deleteSQLServerDatabasePrivilege(ctx, d, schema.TimeoutDelete, deleteSQLServerDatabasePrivilegeClient,
		d.Get("users").(*schema.Set).List()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createSQLServerDatabasePrivilege(ctx context.Context, d *schema.ResourceData, timeout string,
	client *golangsdk.ServiceClient, rawUsers interface{}) error {
	// createSQLServerDatabasePrivilege: create RDS SQL Server database privilege.
	createSQLServerDatabasePrivilegeHttpUrl := "v3/{project_id}/instances/{instance_id}/db_privilege"

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)

	createSQLServerDatabasePrivilegePath := client.Endpoint + createSQLServerDatabasePrivilegeHttpUrl
	createSQLServerDatabasePrivilegePath = strings.ReplaceAll(createSQLServerDatabasePrivilegePath, "{project_id}",
		client.ProjectID)
	createSQLServerDatabasePrivilegePath = strings.ReplaceAll(createSQLServerDatabasePrivilegePath, "{instance_id}",
		d.Get("instance_id").(string))

	createSQLServerDatabasePrivilegeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	users := buildCreateSQLServerDatabasePrivilegeRequestBodyCreateUser(rawUsers)
	start := 0
	end := int(math.Min(maxElementsPerRequest, float64(len(users))))
	for start < end {
		// A single request supports a maximum of 50 elements.
		subUsers := users[start:end]
		createSQLServerDatabasePrivilegeOpt.JSONBody = utils.RemoveNil(buildSQLServerDatabasePrivilegeBodyParams(dbName, subUsers))

		retryFunc := func() (interface{}, bool, error) {
			_, err := client.Request("POST", createSQLServerDatabasePrivilegePath, &createSQLServerDatabasePrivilegeOpt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(timeout),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error creating RDS SQL Server database privilege: %s", err)
		}
		start += maxElementsPerRequest
		end = int(math.Min(float64(end+maxElementsPerRequest), float64(len(users))))
	}
	return nil
}

func deleteSQLServerDatabasePrivilege(ctx context.Context, d *schema.ResourceData, timeout string,
	client *golangsdk.ServiceClient, rawUsers interface{}) error {
	// deleteSQLServerDatabasePrivilege: delete RDS SQL Server database privilege
	deleteSQLServerDatabasePrivilegeHttpUrl := "v3/{project_id}/instances/{instance_id}/db_privilege"

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)

	deleteSQLServerDatabasePrivilegePath := client.Endpoint + deleteSQLServerDatabasePrivilegeHttpUrl
	deleteSQLServerDatabasePrivilegePath = strings.ReplaceAll(deleteSQLServerDatabasePrivilegePath, "{project_id}",
		client.ProjectID)
	deleteSQLServerDatabasePrivilegePath = strings.ReplaceAll(deleteSQLServerDatabasePrivilegePath, "{instance_id}",
		instanceId)

	deleteSQLServerDatabasePrivilegeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	users := buildDeleteSQLServerDatabasePrivilegeRequestBodyDeleteUser(rawUsers)
	start := 0
	end := int(math.Min(maxElementsPerRequest, float64(len(users))))
	for start < end {
		// A single request supports a maximum of 50 elements.
		subUsers := users[start:end]
		deleteSQLServerDatabasePrivilegeOpt.JSONBody = utils.RemoveNil(buildSQLServerDatabasePrivilegeBodyParams(dbName, subUsers))

		retryFunc := func() (interface{}, bool, error) {
			_, err := client.Request("DELETE", deleteSQLServerDatabasePrivilegePath, &deleteSQLServerDatabasePrivilegeOpt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(timeout),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error deleting RDS SQL Server database privilege: %s", err)
		}
		start += maxElementsPerRequest
		end = int(math.Min(float64(end+maxElementsPerRequest), float64(len(users))))
	}
	return nil
}

func buildSQLServerDatabasePrivilegeBodyParams(dbName string, users interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": dbName,
		"users":   users,
	}
	return bodyParams
}

func buildCreateSQLServerDatabasePrivilegeRequestBodyCreateUser(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":     raw["name"],
				"readonly": utils.ValueIgnoreEmpty(raw["readonly"]),
			}
		}
		return rst
	}
	return nil
}

func buildDeleteSQLServerDatabasePrivilegeRequestBodyDeleteUser(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name": raw["name"],
			}
		}
		return rst
	}
	return nil
}
