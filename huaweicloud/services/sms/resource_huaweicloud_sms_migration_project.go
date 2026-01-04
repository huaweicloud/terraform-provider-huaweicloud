package sms

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SMS POST /v3/migprojects
// @API SMS GET /v3/migprojects/{mig_project_id}
// @API SMS PUT /v3/migprojects/{mig_project_id}
// @API SMS DELETE /v3/migprojects/{mig_project_id}
// @API SMS PUT /v3/migprojects/{mig_project_id}/default
func ResourceMigrationProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrationProjectCreate,
		ReadContext:   resourceMigrationProjectRead,
		UpdateContext: resourceMigrationProjectUpdate,
		DeleteContext: resourceMigrationProjectDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			oldValue, newValue := d.GetChange("is_default")
			oldBool := oldValue.(bool)
			newBool := newValue.(bool)
			if oldBool && !newBool {
				return errors.New("the is_default param cannot be changed from true to false")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"use_public_ip": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"exist_server": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"syncing": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"start_target_server": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"speed_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enterprise_project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_network_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceMigrationProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	createHttpUrl := "v3/migprojects"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateAndUpdateMigrationProjectBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SMS migration project: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening creating migration project response: %s", err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SMS migration project: can not found migration project id in return")
	}

	d.SetId(id)

	if d.Get("is_default").(bool) {
		if err := defaultMigrateProject(client, id); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMigrationProjectRead(ctx, d, meta)
}

func defaultMigrateProject(client *golangsdk.ServiceClient, migrateProjectID string) error {
	httpUrl := "v3/migprojects/{mig_project_id}/default"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{mig_project_id}", migrateProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         map[string]interface{}{},
	}

	_, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return err
	}
	return nil
}

func resourceMigrationProjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	migrationProject, err := GetMigrationProject(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving migration project")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", migrationProject, nil)),
		d.Set("description", utils.PathSearch("description", migrationProject, nil)),
		d.Set("is_default", utils.PathSearch("isdefault", migrationProject, nil)),
		d.Set("region", utils.PathSearch("region", migrationProject, nil)),
		d.Set("start_target_server", utils.PathSearch("start_target_server", migrationProject, nil)),
		d.Set("speed_limit", utils.PathSearch("speed_limit", migrationProject, nil)),
		d.Set("use_public_ip", utils.PathSearch("use_public_ip", migrationProject, nil)),
		d.Set("exist_server", utils.PathSearch("exist_server", migrationProject, nil)),
		d.Set("type", utils.PathSearch("type", migrationProject, nil)),
		d.Set("enterprise_project", utils.PathSearch("enterprise_project", migrationProject, nil)),
		d.Set("syncing", utils.PathSearch("syncing", migrationProject, nil)),
		d.Set("start_network_check", utils.PathSearch("start_network_check", migrationProject, false).(bool)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetMigrationProject(client *golangsdk.ServiceClient, migrationProjectId string) (interface{}, error) {
	getHttpUrl := "v3/migprojects/{mig_project_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{mig_project_id}", migrationProjectId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourceMigrationProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	changeList := []string{
		"name", "description", "region", "start_target_server", "speed_limit", "use_public_ip",
		"exist_server", "type", "enterprise_project", "syncing", "start_network_check",
	}
	if d.HasChanges(changeList...) {
		updateHttpUrl := "v3/migprojects/{mig_project_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{mig_project_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(buildCreateAndUpdateMigrationProjectBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating migration project: %s", err)
		}
	}

	if d.HasChange("is_default") {
		oldValue, newValue := d.GetChange("is_default")
		oldBool := oldValue.(bool)
		newBool := newValue.(bool)
		if !oldBool && newBool {
			if err := defaultMigrateProject(client, d.Id()); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceMigrationProjectRead(ctx, d, meta)
}

func buildCreateAndUpdateMigrationProjectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                d.Get("name"),
		"description":         d.Get("description"),
		"region":              d.Get("region"),
		"start_target_server": utils.ValueIgnoreEmpty(d.Get("start_target_server")),
		"speed_limit":         d.Get("speed_limit"),
		"use_public_ip":       d.Get("use_public_ip"),
		"exist_server":        d.Get("exist_server"),
		"type":                d.Get("type"),
		"enterprise_project":  utils.ValueIgnoreEmpty(d.Get("enterprise_project")),
		"syncing":             d.Get("syncing"),
		"start_network_check": utils.ValueIgnoreEmpty(d.Get("start_network_check")),
	}

	return bodyParams
}

func resourceMigrationProjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	deleteHttpUrl := "v3/migprojects/{mig_project_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{mig_project_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: []interface{}{d.Id()},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMS migration project")
	}

	return nil
}
