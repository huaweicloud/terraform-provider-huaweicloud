package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gaussdbInstanceRolePermissionNonUpdatableParams = []string{"instance_id", "db_name"}

// @API GaussDB POST /v3.1/{project_id}/instances/{instance_id}/db-privilege
func ResourceGaussDBInstanceRolePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBInstanceRolePermissionCreate,
		ReadContext:   resourceGaussDBInstanceRolePermissionRead,
		UpdateContext: resourceGaussDBInstanceRolePermissionUpdate,
		DeleteContext: resourceGaussDBInstanceRolePermissionDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussdbInstanceRolePermissionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     gaussDBInstanceRolePermissionUserSchema(),
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

func gaussDBInstanceRolePermissionUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"readonly": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_privilege_grantee": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGaussDBInstanceRolePermissionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/db-privilege"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBInstanceRolePermissionBodyParams(d))

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB instance role permission: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	resourceId := fmt.Sprintf("%s/%s", instanceID, dbName)
	d.SetId(resourceId)

	return nil
}

func buildCreateGaussDBInstanceRolePermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": d.Get("db_name"),
		"user":    buildCreateGaussDBInstanceDatabaseRolePermissionUserBodyParams(d.Get("user")),
	}
	return bodyParams
}

func buildCreateGaussDBInstanceDatabaseRolePermissionUserBodyParams(rawParams interface{}) map[string]interface{} {
	if rawParams == nil {
		return nil
	}
	rawArray := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	readonly := raw["readonly"].(string) == "true"

	bodyParams := map[string]interface{}{
		"name":     raw["name"],
		"readonly": readonly,
		"schema":   raw["schema"],
	}
	if v, ok := raw["default_privilege_grantee"]; ok && v.(string) != "" {
		bodyParams["default_privilege_grantee"] = v
	}
	return bodyParams
}

func resourceGaussDBInstanceRolePermissionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBInstanceRolePermissionUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/db-privilege"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	updateOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBInstanceRolePermissionBodyParams(d))

	_, err = client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating GaussDB instance role permission: %s", err)
	}

	return nil
}

func resourceGaussDBInstanceRolePermissionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB instance role permission resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
