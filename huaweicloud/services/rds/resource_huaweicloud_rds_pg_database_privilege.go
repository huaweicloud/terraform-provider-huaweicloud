package rds

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pgDatabasePrivilegeNonUpdatableParams = []string{"instance_id", "db_name"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_privilege
// @API RDS GET /v3/{project_id}/instances
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/db_privilege
func ResourcePgDatabasePrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgDatabasePrivilegeCreate,
		UpdateContext: resourcePgDatabasePrivilegeUpdate,
		ReadContext:   resourcePgDatabasePrivilegeRead,
		DeleteContext: resourcePgDatabasePrivilegeDelete,

		CustomizeDiff: config.FlexibleForceNew(pgDatabasePrivilegeNonUpdatableParams),

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
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"users": {
				Type:        schema.TypeSet,
				Elem:        pgDatabasePrivilegeCreateUserSchema(),
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

func pgDatabasePrivilegeCreateUserSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the username of the database account.`,
			},
			"readonly": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies the read-only permission.`,
			},
			"schema_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the schema.`,
			},
		},
	}
	return &sc
}

func resourcePgDatabasePrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	err = addPgDatabasePrivilege(ctx, d, schema.TimeoutCreate, client, d.Get("users").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	d.SetId(instanceId + "/" + dbName)

	return nil
}

func resourcePgDatabasePrivilegeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePgDatabasePrivilegeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	oldRaws, newRaws := d.GetChange("users")
	addUsers := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)).List()
	deleteUsers := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List()

	if len(deleteUsers) > 0 {
		err = deletePgDatabasePrivilege(ctx, d, schema.TimeoutUpdate, client, deleteUsers)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(addUsers) > 0 {
		err = addPgDatabasePrivilege(ctx, d, schema.TimeoutUpdate, client, addUsers)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourcePgDatabasePrivilegeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	product := "rds"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	err = deletePgDatabasePrivilege(ctx, d, schema.TimeoutDelete, client, d.Get("users").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func addPgDatabasePrivilege(ctx context.Context, d *schema.ResourceData, timeout string,
	client *golangsdk.ServiceClient, rawUsers interface{}) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/db_privilege"

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)

	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{project_id}", client.ProjectID)
	addPath = strings.ReplaceAll(addPath, "{instance_id}", d.Get("instance_id").(string))

	addOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	users := buildAddPgDatabasePrivilegeRequestBodyUser(rawUsers)
	start := 0
	end := int(math.Min(maxElementsPerRequest, float64(len(users))))
	for start < end {
		// A single request supports a maximum of 50 elements.
		subUsers := users[start:end]
		addOpt.JSONBody = utils.RemoveNil(buildPgDatabasePrivilegeBodyParams(dbName, subUsers))

		retryFunc := func() (interface{}, bool, error) {
			_, err := client.Request("POST", addPath, &addOpt)
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
			return fmt.Errorf("error adding RDS PostgreSQL database privilege: %s", err)
		}
		start += maxElementsPerRequest
		end = int(math.Min(float64(end+maxElementsPerRequest), float64(len(users))))
	}
	return nil
}

func buildPgDatabasePrivilegeBodyParams(dbName string, users interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": dbName,
		"users":   users,
	}
	return bodyParams
}

func buildAddPgDatabasePrivilegeRequestBodyUser(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":        raw["name"],
				"readonly":    raw["readonly"],
				"schema_name": raw["schema_name"],
			}
		}
		return rst
	}
	return nil
}

func deletePgDatabasePrivilege(ctx context.Context, d *schema.ResourceData, timeout string,
	client *golangsdk.ServiceClient, rawUsers interface{}) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/db_privilege"

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	users := buildDeletePgDatabasePrivilegeRequestBodyUser(rawUsers)
	start := 0
	end := int(math.Min(maxElementsPerRequest, float64(len(users))))
	for start < end {
		// A single request supports a maximum of 50 elements.
		subUsers := users[start:end]
		deleteOpt.JSONBody = utils.RemoveNil(buildPgDatabasePrivilegeBodyParams(dbName, subUsers))

		retryFunc := func() (interface{}, bool, error) {
			_, err := client.Request("DELETE", deletePath, &deleteOpt)
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
			return fmt.Errorf("error deleting RDS PostgreSQL database privilege: %s", err)
		}
		start += maxElementsPerRequest
		end = int(math.Min(float64(end+maxElementsPerRequest), float64(len(users))))
	}
	return nil
}

func buildDeletePgDatabasePrivilegeRequestBodyUser(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":        raw["name"],
				"schema_name": raw["schema_name"],
			}
		}
		return rst
	}
	return nil
}
