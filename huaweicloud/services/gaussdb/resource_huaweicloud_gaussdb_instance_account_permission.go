package gaussdb

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceAccountPermissionNonUpdatableParams = []string{"instance_id", "db_name"}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/db-privilege
func ResourceGaussdbInstanceAccountPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussdbInstanceAccountPermissionCreate,
		ReadContext:   resourceGaussdbInstanceAccountPermissionRead,
		UpdateContext: resourceGaussdbInstanceAccountPermissionUpdate,
		DeleteContext: resourceGaussdbInstanceAccountPermissionDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceAccountPermissionNonUpdatableParams),

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
			"users": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     instanceAccountPermissionUsersSchema(),
				Set:      instanceAccountPermissionUsersHash,
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

func instanceAccountPermissionUsersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"readonly": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"schema_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGaussdbInstanceAccountPermissionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-privilege"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	usersSet := d.Get("users").(*schema.Set)
	usersList := usersSet.List()

	start := 0
	end := int(math.Min(50, float64(len(usersList))))

	for start < end {
		createOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}
		createOpt.JSONBody = utils.RemoveNil(buildCreateGaussdbInstanceAccountPermissionBodyParams(d, usersList[start:end]))

		_, err = client.Request("POST", createPath, &createOpt)
		if err != nil {
			return diag.Errorf("error creating GaussDB instance account permission: %s", err)
		}
		start += 50
		end = int(math.Min(float64(end+50), float64(len(usersList))))
	}

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	resourceId := fmt.Sprintf("%s/%s", instanceId, dbName)
	d.SetId(resourceId)

	return nil
}

func resourceGaussdbInstanceAccountPermissionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussdbInstanceAccountPermissionUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-privilege"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	usersSet := d.Get("users").(*schema.Set)
	usersList := usersSet.List()

	start := 0
	end := int(math.Min(50, float64(len(usersList))))

	for start < end {
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
		updateOpt.JSONBody = utils.RemoveNil(buildCreateGaussdbInstanceAccountPermissionBodyParams(d, usersList[start:end]))

		_, err = client.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating GaussDB instance account permission: %s", err)
		}
		start += 50
		end = int(math.Min(float64(end+50), float64(len(usersList))))
	}

	return nil
}

func resourceGaussdbInstanceAccountPermissionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB instance account permission resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func buildCreateGaussdbInstanceAccountPermissionBodyParams(d *schema.ResourceData, users []interface{}) map[string]interface{} {
	usersParams := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		userMap := user.(map[string]interface{})
		usersParams = append(usersParams, map[string]interface{}{
			"name":        userMap["name"],
			"readonly":    userMap["readonly"].(string) == "true",
			"schema_name": userMap["schema_name"],
		})
	}
	return map[string]interface{}{
		"db_name": d.Get("db_name"),
		"users":   usersParams,
	}
}

func instanceAccountPermissionUsersHash(v interface{}) int {
	m := v.(map[string]interface{})
	key := fmt.Sprintf("%s-%s", m["name"].(string), m["schema_name"].(string))
	return schema.HashString(key)
}
