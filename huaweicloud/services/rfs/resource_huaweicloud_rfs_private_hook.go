package rfs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/private-hooks
// @API RFS GET /v1/private-hooks/{hook_name}/metadata
// @API RFS GET /v1/private-hooks/{hook_name}/versions/{hook_version}/metadata
// @API RFS POST /v1/private-hooks/{hook_name}/versions
// @API RFS DELETE /v1/private-hooks/{hook_name}/versions/{hook_version}
// @API RFS PATCH /v1/private-hooks/{hook_name}/metadata
// @API RFS DELETE /v1/private-hooks/{hook_name}
func ResourcePrivateHook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateHookCreate,
		ReadContext:   resourcePrivateHookRead,
		UpdateContext: resourcePrivateHookUpdate,
		DeleteContext: resourcePrivateHookDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePrivateHookImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the private hook is located.`,
			},
			// Arguments
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the private hook.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version of the private hook.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the private hook.`,
			},
			"version_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the private hook version.`,
			},
			"policy_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  `The OBS address of the policy file.`,
				ExactlyOneOf: []string{"policy_body"},
			},
			"policy_body": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The policy content of the private hook.`,
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_stacks": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The target resource stack for the private hook to take effect.`,
						},
						"failure_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The behavior after private hook verification fails.`,
						},
					},
				},
				Description: `The configuration of the private hook.`,
			},
			"keep_old_version": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether keeping old version while updating hook version.`,
			},

			// Attributes
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation of the private hook, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update of the private hook, in RFC3339 format.`,
			},
		},
	}
}

func buildPrivateHookCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"hook_name":                d.Get("name"),
		"hook_version":             d.Get("version"),
		"hook_description":         utils.ValueIgnoreEmpty(d.Get("description")),
		"hook_version_description": utils.ValueIgnoreEmpty(d.Get("version_description")),
		"policy_uri":               utils.ValueIgnoreEmpty(d.Get("policy_uri")),
		"policy_body":              utils.ValueIgnoreEmpty(d.Get("policy_body")),
		"configuration":            utils.ValueIgnoreEmpty(buildPrivateHookConfigurationBodyParams(d.Get("configuration").([]interface{}))),
	}
}

func buildPrivateHookConfigurationBodyParams(configurations []interface{}) map[string]interface{} {
	if len(configurations) < 1 {
		return nil
	}

	configuration := configurations[0]
	return map[string]interface{}{
		"target_stacks": utils.ValueIgnoreEmpty(utils.PathSearch("target_stacks", configuration, nil)),
		"failure_mode":  utils.ValueIgnoreEmpty(utils.PathSearch("failure_mode", configuration, nil)),
	}
}

func resourcePrivateHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/private-hooks"
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}
	createPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
		JSONBody: utils.RemoveNil(buildPrivateHookCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating private hook: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// Considering that enterprise projects may be supported in the future, hook_id is used as the resource ID to
	// facilitate calling public method.
	resourceId := utils.PathSearch("hook_id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the hook ID from the API response: %#v", respBody)
	}
	d.SetId(resourceId)

	return resourcePrivateHookRead(ctx, d, meta)
}

func flattenPrivateHookConfiguration(configuration interface{}) []map[string]interface{} {
	if configuration == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"target_stacks": utils.PathSearch("target_stacks", configuration, nil),
			"failure_mode":  utils.PathSearch("failure_mode", configuration, nil),
		},
	}
}

func queryPrivateHookByName(client *golangsdk.ServiceClient, name string) (interface{}, error) {
	httpUrl := "v1/private-hooks/{hook_name}/metadata"

	// Generate a random UUID as the RFS request ID.
	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate RFS request ID: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{hook_name}", name)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func queryPrivateHookVersionMetadata(client *golangsdk.ServiceClient, name, version string) (interface{}, error) {
	httpUrl := "v1/private-hooks/{hook_name}/versions/{hook_version}/metadata"

	// Generate a random UUID as the RFS request ID.
	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate RFS request ID: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{hook_name}", name)
	getPath = strings.ReplaceAll(getPath, "{hook_version}", version)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourcePrivateHookRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		hookName = d.Get("name").(string)
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	respBody, err := queryPrivateHookByName(client, hookName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying private hook (%s)", hookName))
	}

	currentVersion := utils.PathSearch("default_version", respBody, "").(string)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("hook_name", respBody, nil)),
		d.Set("description", utils.PathSearch("hook_description", respBody, nil)),
		d.Set("version", currentVersion),
		d.Set("configuration", flattenPrivateHookConfiguration(utils.PathSearch("configuration", respBody, nil))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			respBody, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_time",
			respBody, "").(string))/1000, false)),
	)

	versionMetadata, err := queryPrivateHookVersionMetadata(client, hookName, currentVersion)
	if err != nil {
		log.Printf("[ERROR] error querying the information about the hook version (%s): %s", currentVersion, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("version_description", utils.PathSearch("hook_version_description", versionMetadata, nil)))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving stack hook (%s) fields: %s", hookName, err)
	}
	return nil
}

func isHookVersionExist(client *golangsdk.ServiceClient, hookName, version string) bool {
	if resp, err := queryPrivateHookVersionMetadata(client, hookName, version); err == nil {
		log.Printf("[DEBUG] find the specified version (%s): %#v", version, resp)
		return true
	}
	return false
}

func buildPrivateHookVersionCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"hook_id":                  d.Id(),
		"hook_version":             d.Get("version"),
		"hook_version_description": utils.ValueIgnoreEmpty(d.Get("version_description")),
		"policy_uri":               utils.ValueIgnoreEmpty(d.Get("policy_uri")),
		"policy_body":              utils.ValueIgnoreEmpty(d.Get("policy_body")),
	}
}

func createNewHookVersion(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl = "v1/private-hooks/{hook_name}/versions"
		name    = d.Get("name").(string)
	)
	// Generate a random UUID as the RFS request ID.
	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("unable to generate RFS request ID: %s", err)
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{hook_name}", name)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
		JSONBody: buildPrivateHookVersionCreateBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &opt)
	return err
}

func deleteOldHookVersion(client *golangsdk.ServiceClient, hookName, hookVersion string) error {
	var (
		httpUrl = "v1/private-hooks/{hook_name}/versions/{hook_version}"
	)
	// Generate a random UUID as the RFS request ID.
	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("unable to generate RFS request ID: %s", err)
	}
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{hook_name}", hookName)
	deletePath = strings.ReplaceAll(deletePath, "{hook_version}", hookVersion)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	return err
}

func updatePrivateHookMetadata(client *golangsdk.ServiceClient, name string, requestBody map[string]interface{}) error {
	httpUrl := "v1/private-hooks/{hook_name}/metadata"
	// Generate a random UUID as the RFS request ID.
	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("unable to generate RFS request ID: %s", err)
	}
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{hook_name}", name)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
		JSONBody: requestBody,
	}

	_, err = client.Request("PATCH", updatePath, &opt)
	return err
}

func resourcePrivateHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                    = meta.(*config.Config)
		region                 = cfg.GetRegion(d)
		hookName               = d.Get("name").(string)
		isVersionChanged       = d.HasChange("version")
		oldVersion, newVersion interface{}
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestBody := map[string]interface{}{
		"hook_description": d.Get("description"),
		"configuration":    utils.ValueIgnoreEmpty(buildPrivateHookConfigurationBodyParams(d.Get("configuration").([]interface{}))),
	}
	// Before update the version via update metadata API, we should check the corresponding version is exist.
	// If not, create it and using this version to update the hook metadata.
	if d.HasChange("version") {
		oldVersion, newVersion = d.GetChange("version")
		if !isHookVersionExist(client, hookName, newVersion.(string)) {
			err = createNewHookVersion(client, d)
			if err != nil {
				return diag.Errorf("error creating hook version: %s", err)
			}
		}
		requestBody["default_version"] = newVersion.(string)
	} else if d.HasChanges("version_description", "policy_uri", "policy_body") {
		return diag.Errorf("Unable to update version description, policy URI or policy body without updating version")
	}

	err = updatePrivateHookMetadata(client, hookName, requestBody)
	if err != nil {
		return diag.Errorf("error updating private hook (%s): %s", hookName, err)
	}

	// Since a hook can only retain up to 200 versions, by default, the old version needs to be deleted after the update.
	// Whenever a version is updated, it must have a history value.
	if _, ok := d.GetOk("keep_old_version"); !ok && isVersionChanged {
		err = deleteOldHookVersion(client, hookName, oldVersion.(string))
		if err != nil {
			log.Printf("unable to delete the old version (%s): %s", oldVersion, err)
		}
	}
	return resourcePrivateHookRead(ctx, d, meta)
}

func resourcePrivateHookDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/private-hooks/{hook_name}"
		hookName = d.Get("name").(string)
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{hook_name}", hookName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
	}

	// When deleting the private hook, all historical versions will be deleted together.
	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return diag.Errorf("error deleting private hook (%s): %s", hookName, err)
	}
	return nil
}

func resourcePrivateHookImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		hookName = d.Id()
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	respBody, err := queryPrivateHookByName(client, hookName)
	if err != nil {
		return nil, err
	}

	hookId := utils.PathSearch("hook_id", respBody, "").(string)
	if hookId == "" {
		return nil, fmt.Errorf("unable to find the hook ID from the RFS query response: %#v", respBody)
	}

	d.SetId(hookId)
	return []*schema.ResourceData{d}, d.Set("name", hookName)
}
