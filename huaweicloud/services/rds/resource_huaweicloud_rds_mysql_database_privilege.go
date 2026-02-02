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

var mysqlDatabasePrivilegeNonUpdatableParams = []string{"instance_id", "db_name"}

// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db_privilege
// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_privilege
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/db_user
func ResourceMysqlDatabasePrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMysqlDatabasePrivilegeCreate,
		UpdateContext: resourceMysqlDatabasePrivilegeUpdate,
		ReadContext:   resourceMysqlDatabasePrivilegeRead,
		DeleteContext: resourceMysqlDatabasePrivilegeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(mysqlDatabasePrivilegeNonUpdatableParams),

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
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"users": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the account that associated with the database.`,
				Elem:        mysqlDatabasePrivilegeUserSchema(),
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

func mysqlDatabasePrivilegeUserSchema() *schema.Resource {
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

func resourceMysqlDatabasePrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createMysqlDatabasePrivilege: create RDS Mysql database privilege.
	createMysqlDatabasePrivilegeProduct := "rds"
	createMysqlDatabasePrivilegeClient, err := cfg.NewServiceClient(createMysqlDatabasePrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	err = createMysqlDatabasePrivilege(ctx, d, createMysqlDatabasePrivilegeClient, d.Get("users").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	d.SetId(instanceId + "/" + dbName)

	return resourceMysqlDatabasePrivilegeRead(ctx, d, meta)
}

func resourceMysqlDatabasePrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getMysqlDatabasePrivilege: query RDS Mysql database privilege
	var (
		getMysqlDatabasePrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/database/db_user"
		getMysqlDatabasePrivilegeProduct = "rds"
	)
	getMysqlDatabasePrivilegeClient, err := cfg.NewServiceClient(getMysqlDatabasePrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<db_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getMysqlDatabasePrivilegePath := getMysqlDatabasePrivilegeClient.Endpoint + getMysqlDatabasePrivilegeHttpUrl
	getMysqlDatabasePrivilegePath = strings.ReplaceAll(getMysqlDatabasePrivilegePath, "{project_id}",
		getMysqlDatabasePrivilegeClient.ProjectID)
	getMysqlDatabasePrivilegePath = strings.ReplaceAll(getMysqlDatabasePrivilegePath, "{instance_id}", instanceId)

	getMysqlDatabasePrivilegeQueryParams := buildGetMysqlDatabasePrivilegeQueryParams(dbName)
	getMysqlDatabasePrivilegePath += getMysqlDatabasePrivilegeQueryParams

	getMysqlDatabasePrivilegeResp, err := pagination.ListAllItems(
		getMysqlDatabasePrivilegeClient,
		"page",
		getMysqlDatabasePrivilegePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS Mysql database privilege")
	}

	getMysqlDatabasePrivilegeRespJson, err := json.Marshal(getMysqlDatabasePrivilegeResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getMysqlDatabasePrivilegeRespBody interface{}
	err = json.Unmarshal(getMysqlDatabasePrivilegeRespJson, &getMysqlDatabasePrivilegeRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	users := flattenGetMysqlDatabasePrivilegeResponseBodyGetUser(d, getMysqlDatabasePrivilegeRespBody)
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

func flattenGetMysqlDatabasePrivilegeResponseBodyGetUser(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	usersRaw := d.Get("users").(*schema.Set)
	userNames := make(map[string]bool)
	for _, userRaw := range usersRaw.List() {
		v := userRaw.(map[string]interface{})
		userNames[v["name"].(string)] = true
	}

	curJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		name := utils.PathSearch("name", v, "").(string)
		// for import, all users will be imported
		if len(userNames) > 0 && !userNames[name] {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"name":     name,
			"readonly": utils.PathSearch("readonly", v, nil),
		})
	}
	return rst
}

func buildGetMysqlDatabasePrivilegeQueryParams(dbName string) string {
	return fmt.Sprintf("?db-name=%s&page=1&limit=100", dbName)
}

func resourceMysqlDatabasePrivilegeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateMysqlDatabasePrivilege: update RDS Mysql database privilege.
	var (
		updateMysqlDatabasePrivilegeProduct = "rds"
	)
	updateMysqlDatabasePrivilegeClient, err := cfg.NewServiceClient(updateMysqlDatabasePrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	oldRaws, newRaws := d.GetChange("users")
	createUsers := newRaws.(*schema.Set).List()
	deleteUsers := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List()

	if len(deleteUsers) > 0 {
		err = deleteMysqlDatabasePrivilege(ctx, d, updateMysqlDatabasePrivilegeClient, deleteUsers)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(createUsers) > 0 {
		err = createMysqlDatabasePrivilege(ctx, d, updateMysqlDatabasePrivilegeClient, createUsers)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMysqlDatabasePrivilegeRead(ctx, d, meta)
}

func resourceMysqlDatabasePrivilegeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteMysqlDatabasePrivilege: delete RDS Mysql database privilege
	deleteMysqlDatabasePrivilegeProduct := "rds"
	deleteMysqlDatabasePrivilegeClient, err := cfg.NewServiceClient(deleteMysqlDatabasePrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	err = deleteMysqlDatabasePrivilege(ctx, d, deleteMysqlDatabasePrivilegeClient, d.Get("users").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createMysqlDatabasePrivilege(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	rawUsers interface{}) error {
	// createMysqlDatabasePrivilege: create RDS Mysql database privilege.
	createMysqlDatabasePrivilegeHttpUrl := "v3/{project_id}/instances/{instance_id}/db_privilege"

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)

	createMysqlDatabasePrivilegePath := client.Endpoint + createMysqlDatabasePrivilegeHttpUrl
	createMysqlDatabasePrivilegePath = strings.ReplaceAll(createMysqlDatabasePrivilegePath, "{project_id}",
		client.ProjectID)
	createMysqlDatabasePrivilegePath = strings.ReplaceAll(createMysqlDatabasePrivilegePath, "{instance_id}",
		instanceId)

	createMysqlDatabasePrivilegeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	users := buildCreateMysqlDatabasePrivilegeRequestBodyUser(rawUsers)
	start := 0
	end := int(math.Min(50, float64(len(users))))
	for start < end {
		// A single request supports a maximum of 50 elements.
		subUsers := users[start:end]
		createMysqlDatabasePrivilegeOpt.JSONBody = utils.RemoveNil(buildMysqlDatabasePrivilegeBodyParams(dbName, subUsers))
		log.Printf("[DEBUG] Create RDS Mysql database privilege options: %#v", createMysqlDatabasePrivilegeOpt)

		retryFunc := func() (interface{}, bool, error) {
			_, err := client.Request("POST", createMysqlDatabasePrivilegePath, &createMysqlDatabasePrivilegeOpt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutCreate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error creating RDS Mysql database privilege: %s", err)
		}
		start += 50
		end = int(math.Min(float64(end+50), float64(len(users))))
	}
	return nil
}

func deleteMysqlDatabasePrivilege(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	rawUsers interface{}) error {
	// deleteMysqlDatabasePrivilege: delete RDS Mysql database privilege
	deleteMysqlDatabasePrivilegeHttpUrl := "v3/{project_id}/instances/{instance_id}/db_privilege"

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)

	deleteMysqlDatabasePrivilegePath := client.Endpoint + deleteMysqlDatabasePrivilegeHttpUrl
	deleteMysqlDatabasePrivilegePath = strings.ReplaceAll(deleteMysqlDatabasePrivilegePath, "{project_id}",
		client.ProjectID)
	deleteMysqlDatabasePrivilegePath = strings.ReplaceAll(deleteMysqlDatabasePrivilegePath, "{instance_id}",
		instanceId)

	deleteMysqlDatabasePrivilegeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	users := buildDeleteMysqlDatabasePrivilegeRequestBodyDeleteUser(rawUsers)
	start := 0
	end := int(math.Min(50, float64(len(users))))
	for start < end {
		// A single request supports a maximum of 50 elements.
		subUsers := users[start:end]
		deleteMysqlDatabasePrivilegeOpt.JSONBody = utils.RemoveNil(buildMysqlDatabasePrivilegeBodyParams(dbName, subUsers))
		log.Printf("[DEBUG] Delete RDS Mysql database privilege options: %#v", deleteMysqlDatabasePrivilegeOpt)

		retryFunc := func() (interface{}, bool, error) {
			_, err := client.Request("DELETE", deleteMysqlDatabasePrivilegePath, &deleteMysqlDatabasePrivilegeOpt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error deleting RDS Mysql database privilege: %s", err)
		}
		start += 50
		end = int(math.Min(float64(end+50), float64(len(users))))
	}
	return nil
}

func buildMysqlDatabasePrivilegeBodyParams(dbName string, users interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": dbName,
		"users":   users,
	}
	return bodyParams
}

func buildCreateMysqlDatabasePrivilegeRequestBodyUser(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":     utils.ValueIgnoreEmpty(raw["name"]),
				"readonly": utils.ValueIgnoreEmpty(raw["readonly"]),
			}
		}
		return rst
	}
	return nil
}

func buildDeleteMysqlDatabasePrivilegeRequestBodyDeleteUser(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name": utils.ValueIgnoreEmpty(raw["name"]),
			}
		}
		return rst
	}
	return nil
}
