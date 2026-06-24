package taurusdb

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

var htapStarRocksUserNoneUpdatableParams = []string{
	"instance_id", "user_name",
}

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/users
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/users
// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/starrocks/users/permission
// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/starrocks/users/password
// @API TaurusDB DELETE /v3/{project_id}/instances/{instance_id}/starrocks/users
func ResourceTaurusDBHtapStarrocksUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksUserCreate,
		UpdateContext: resourceTaurusDBHtapStarrocksUserUpdate,
		ReadContext:   resourceTaurusDBHtapStarrocksUserRead,
		DeleteContext: resourceTaurusDBHtapStarrocksUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceHtapStarrocksUserImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(htapStarRocksUserNoneUpdatableParams),

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
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"databases": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dml": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ddl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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

func resourceTaurusDBHtapStarrocksUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	userName := d.Get("user_name").(string)

	createHttpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/users"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildCreateStarrocksUserBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating TaurusDB HTAP StarRocks user: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	result := utils.PathSearch("result", respBody, "").(string)
	if result != "SUCCESS" {
		return diag.Errorf("error creating TaurusDB HTAP StarRocks user: result is %s", result)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, userName))

	return resourceTaurusDBHtapStarrocksUserRead(ctx, d, meta)
}

func buildCreateStarrocksUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_name": d.Get("user_name"),
		"password":  d.Get("password"),
		"databases": utils.ExpandToStringList(d.Get("databases").(*schema.Set).List()),
	}
	// set param ddl only when it's set by user
	rawConfig := d.GetRawConfig()
	if dmlVal := rawConfig.GetAttr("dml"); !dmlVal.IsNull() {
		bodyParams["dml"] = d.Get("dml").(int)
	}
	// set param ddl only when it's set by user
	if ddlVal := rawConfig.GetAttr("ddl"); !ddlVal.IsNull() {
		bodyParams["ddl"] = d.Get("ddl").(int)
	}
	return bodyParams
}

func resourceTaurusDBHtapStarrocksUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceID = d.Get("instance_id").(string)
		userName   = d.Get("user_name").(string)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	// Query the user by user_name filter
	users, err := QueryHtapStarrocksUsers(client, instanceID, userName)
	if err != nil {
		err = common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200076")
		err = common.ConvertExpected403ErrInto404Err(err, "error_name", "DBS.200044")
		return common.CheckDeletedDiag(d, err, "error retrieving TaurusDB HTAP StarRocks user")
	}

	// No user found, return 404
	if len(users) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving TaurusDB HTAP StarRocks user")
	}

	user := users[0]
	databasesRaw := utils.PathSearch("databases", user, make([]interface{}, 0)).([]interface{})

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("user_name", userName),
		d.Set("databases", utils.ExpandToStringList(databasesRaw)),
		d.Set("dml", int(utils.PathSearch("dml", user, float64(0)).(float64))),
		d.Set("ddl", int(utils.PathSearch("ddl", user, float64(0)).(float64))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTaurusDBHtapStarrocksUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceID = d.Get("instance_id").(string)
		userName   = d.Get("user_name").(string)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	if d.HasChange("password") {
		if err := updateStarrocksUserPassword(ctx, d, client, instanceID, userName); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("databases", "dml", "ddl") {
		if err := updateStarrocksUserPermission(ctx, d, client, instanceID, userName); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaurusDBHtapStarrocksUserRead(ctx, d, meta)
}

func updateStarrocksUserPassword(_ context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID, userName string) error {
	updateHttpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/users/password"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceID)

	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	updateOpts.JSONBody = map[string]interface{}{
		"user_name": userName,
		"password":  d.Get("password"),
	}

	resp, err := client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return fmt.Errorf("error updating TaurusDB HTAP StarRocks user password: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	result := utils.PathSearch("result", respBody, "").(string)
	if result != "SUCCESS" {
		return fmt.Errorf("error updating TaurusDB HTAP StarRocks user password: result is %s", result)
	}

	return nil
}

func updateStarrocksUserPermission(_ context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID, userName string) error {
	updateHttpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/users/permission"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceID)

	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildUpdateStarrocksUserBodyParams(d, userName),
	}

	resp, err := client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return fmt.Errorf("error updating TaurusDB HTAP StarRocks user permission: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	result := utils.PathSearch("result", respBody, "").(string)
	if result != "SUCCESS" {
		return fmt.Errorf("error updating TaurusDB HTAP StarRocks user permission: result is %s", result)
	}

	return nil
}

func buildUpdateStarrocksUserBodyParams(d *schema.ResourceData, userName string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_name": userName,
	}
	if d.HasChange("databases") {
		bodyParams["databases"] = utils.ExpandToStringList(d.Get("databases").(*schema.Set).List())
	}
	if d.HasChange("dml") {
		// set param ddl only when it's set by user
		rawConfig := d.GetRawConfig()
		if dmlVal := rawConfig.GetAttr("dml"); !dmlVal.IsNull() {
			bodyParams["dml"] = d.Get("dml").(int)
		}
	}
	if d.HasChange("ddl") {
		// set param ddl only when it's set by user
		rawConfig := d.GetRawConfig()
		if ddlVal := rawConfig.GetAttr("ddl"); !ddlVal.IsNull() {
			bodyParams["ddl"] = d.Get("ddl").(int)
		}
	}
	return bodyParams
}

func resourceTaurusDBHtapStarrocksUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceID = d.Get("instance_id").(string)
		userName   = d.Get("user_name").(string)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	deleteHttpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/users"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)
	deletePath = fmt.Sprintf("%s?user_name=%s", deletePath, userName)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting TaurusDB HTAP StarRocks user")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	result := utils.PathSearch("result", respBody, "").(string)
	if result != "SUCCESS" {
		return diag.Errorf("error deleting TaurusDB HTAP StarRocks user: result is %s", result)
	}

	return nil
}

func resourceHtapStarrocksUserImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<user_name>, but got `%s`", importId)
	}

	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("user_name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
