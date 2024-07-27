package rds

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS POST /v3/{project_id}/instances/{instance_id}/db_privilege
// @API RDS GET /v3/{project_id}/instances
func ResourcePgDatabasePrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgDatabasePrivilegeCreate,
		UpdateContext: resourcePgDatabasePrivilegeUpdate,
		ReadContext:   resourcePgDatabasePrivilegeRead,
		DeleteContext: resourcePgDatabasePrivilegeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
				Description: `Specifies the ID of the RDS PostgreSQL instance.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the database name.`,
			},
			"users": {
				Type:        schema.TypeSet,
				Elem:        pgDatabasePrivilegeCreateUserSchema(),
				Required:    true,
				Description: `Specifies the account that associated with the database`,
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

	err = updatePgDatabasePrivilege(ctx, d, schema.TimeoutCreate, client, d.Get("users").(*schema.Set).List())
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

	err = updatePgDatabasePrivilege(ctx, d, schema.TimeoutUpdate, client, d.Get("users").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updatePgDatabasePrivilege(ctx context.Context, d *schema.ResourceData, timeout string,
	client *golangsdk.ServiceClient, rawUsers interface{}) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/db_privilege"

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	users := buildUpdatePgDatabasePrivilegeRequestBodyCreateUser(rawUsers)
	start := 0
	end := int(math.Min(maxElementsPerRequest, float64(len(users))))
	for start < end {
		// A single request supports a maximum of 50 elements.
		subUsers := users[start:end]
		updateOpt.JSONBody = utils.RemoveNil(buildPgDatabasePrivilegeBodyParams(dbName, subUsers))

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
			return fmt.Errorf("error creating RDS PostgreSQL database privilege: %s", err)
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

func buildUpdatePgDatabasePrivilegeRequestBodyCreateUser(rawParams interface{}) []map[string]interface{} {
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

func resourcePgDatabasePrivilegeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS PostgreSQL database privilege resource is not supported. The resource is only removed " +
		"from the state, the RDS PostgreSQL database privilege remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
