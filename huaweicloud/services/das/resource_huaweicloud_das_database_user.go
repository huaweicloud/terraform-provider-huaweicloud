package das

import (
	"context"
	"fmt"
	"strings"

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
	databaseUserNonUpdatableParams = []string{
		"instance_id",
	}

	databaseUserNotFoundCodes = []string{
		"DAS.220003", // Database user does not exist during query or deletion.
	}
)

// @API DAS POST /v3/{project_id}/instances/{instance_id}/db-users
// @API DAS GET /v3/{project_id}/instances/{instance_id}/db-users/{db_user_id}
// @API DAS PUT /v3/{project_id}/instances/{instance_id}/db-users/{db_user_id}
// @API DAS DELETE /v3/{project_id}/instances/{instance_id}/db-users/{db_user_id}
func ResourceDatabaseUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseUserCreate,
		ReadContext:   resourceDatabaseUserRead,
		UpdateContext: resourceDatabaseUserUpdate,
		DeleteContext: resourceDatabaseUserDelete,

		CustomizeDiff: config.FlexibleForceNew(databaseUserNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceDatabaseUserImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the database user is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance to which the database user belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the database user.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of the database user.",
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
		},
	}
}

func buildCreateDatabaseUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"db_username":      d.Get("name"),
		"db_user_password": d.Get("password"),
		"datastore_type":   "mysql",
	}
}

func createDatabaseUser(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/db-users"
		instanceId = d.Get("instance_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createDatabaseUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildCreateDatabaseUserBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createDatabaseUserOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	respBody, err := createDatabaseUser(client, d)
	if err != nil {
		return diag.Errorf("error creating DAS Database user: %s", err)
	}

	dbUserId := utils.PathSearch("db_user_id", respBody, "").(string)
	if dbUserId == "" {
		return diag.Errorf("unable to find the ID of the DAS Database user from the API response")
	}
	d.SetId(dbUserId)

	return resourceDatabaseUserRead(ctx, d, meta)
}

func GetDatabaseUserById(client *golangsdk.ServiceClient, instanceId, dbUserId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/db-users/{db_user_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{db_user_id}", dbUserId)
	getDatabaseUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getDatabaseUserOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceDatabaseUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		dbUserId   = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	resp, err := GetDatabaseUserById(client, instanceId, dbUserId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", databaseUserNotFoundCodes...),
			fmt.Sprintf("error retrieving Database user (%s)", dbUserId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("db_user.db_username", resp, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDatabaseUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"db_username":      d.Get("name"),
		"db_user_password": d.Get("password"),
	}
}

func updateDatabaseUser(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/db-users/{db_user_id}"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{db_user_id}", d.Id())

	updateDatabaseUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildUpdateDatabaseUserBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateDatabaseUserOpt)
	return err
}

func resourceDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS Client: %s", err)
	}

	err = updateDatabaseUser(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDatabaseUserRead(ctx, d, meta)
}

func deleteDatabaseUser(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/db-users/{db_user_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{db_user_id}", d.Id())

	deleteDatabaseUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteDatabaseUserOpt)
	return err
}

func resourceDatabaseUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		dbUserId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	err = deleteDatabaseUser(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", databaseUserNotFoundCodes...),
			fmt.Sprintf("error deleting DAS Database user (%s)", dbUserId),
		)
	}

	return nil
}

func resourceDatabaseUserImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', "+
			"but got '%s'", importedId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
