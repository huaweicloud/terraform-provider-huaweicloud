package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/app-groups/{app_group_id}/apps
// @API Workspace GET /v1/{project_id}/app-groups/{app_group_id}/apps
// @API Workspace PATCH /v1/{project_id}/app-groups/{app_group_id}/apps/{app_id}
// @API Workspace POST /v1/{project_id}/app-groups/{app_group_id}/apps/batch-unpublish
func ResourceAppApplicationPublishment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppApplicationPublishmentCreate,
		ReadContext:   resourceAppApplicationPublishmentRead,
		UpdateContext: resourceAppApplicationPublishmentUpdate,
		DeleteContext: resourceAppApplicationPublishmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppApplicationPublishmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"app_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The APP group ID to which the application belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the application.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the application.`,
			},
			"execute_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The location where the application file is installed.`,
			},
			"sandbox_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to run in sandbox mode.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version of the application.`,
			},
			"publisher": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The publisher of the application.`,
			},
			"work_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The working directory of the application.`,
			},
			"command_param": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The command line parameter used to start the application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the application.`,
			},
			"icon_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The path where the application icon is located.`,
			},
			"icon_index": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `The icon index of the application.`,
			},
			"source_image_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of image IDs corresponding to the server instance to which the application belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The current status of the application.`,
			},
			"published_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The release time of the application, in RFC3339 format.`,
			},
		},
	}
}

func resourceAppApplicationPublishmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	httpUrl := "v1/{project_id}/app-groups/{app_group_id}/apps"
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	appGroupId := d.Get("app_group_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{app_group_id}", appGroupId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAppApplicationPublishmentBodyParams(d)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error publishing application: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	applicationId := utils.PathSearch("items|[0].id", respBody, "").(string)
	d.SetId(applicationId)

	if d.HasChange("status") {
		updateStateOpt := map[string]interface{}{
			"state": d.Get("status"),
		}
		if err = updateAppPublishedApplication(client, appGroupId, applicationId, updateStateOpt); err != nil {
			return diag.Errorf("error updating status of the application (%s): %s", d.Get("name").(string), err)
		}
	}

	return resourceAppApplicationPublishmentRead(ctx, d, meta)
}

func buildCreateAppApplicationPublishmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	rest := make([]map[string]interface{}, 0)
	params := map[string]interface{}{
		"name":             d.Get("name"),
		"source_type":      d.Get("type"),
		"execute_path":     d.Get("execute_path"),
		"sandbox_enable":   utils.ValueIgnoreEmpty(d.Get("sandbox_enable")),
		"version":          utils.ValueIgnoreEmpty(d.Get("version")),
		"publisher":        utils.ValueIgnoreEmpty(d.Get("publisher")),
		"work_path":        utils.ValueIgnoreEmpty(d.Get("work_path")),
		"command_param":    utils.ValueIgnoreEmpty(d.Get("command_param")),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"icon_path":        utils.ValueIgnoreEmpty(d.Get("icon_path")),
		"icon_index":       utils.ValueIgnoreEmpty(d.Get("icon_index")),
		"source_image_ids": utils.ExpandToStringList(d.Get("source_image_ids").([]interface{})),
	}

	rest = append(rest, params)
	return map[string]interface{}{
		"items": rest,
	}
}

func resourceAppApplicationPublishmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	application, err := GetAppPublishedApplicationByName(client, d.Get("app_group_id").(string), d.Get("name").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace application")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("app_group_id", utils.PathSearch("app_group_id", application, nil)),
		d.Set("name", utils.PathSearch("name", application, nil)),
		d.Set("type", utils.PathSearch("source_type", application, nil)),
		d.Set("execute_path", utils.PathSearch("execute_path", application, nil)),
		d.Set("sandbox_enable", utils.PathSearch("sandbox_enable", application, false)),
		d.Set("version", utils.PathSearch("version", application, nil)),
		d.Set("publisher", utils.PathSearch("publisher", application, nil)),
		d.Set("work_path", utils.PathSearch("work_path", application, nil)),
		d.Set("command_param", utils.PathSearch("command_param", application, nil)),
		d.Set("description", utils.PathSearch("description", application, nil)),
		d.Set("icon_path", utils.PathSearch("icon_path", application, nil)),
		d.Set("icon_index", utils.PathSearch("icon_index", application, nil)),
		d.Set("status", utils.PathSearch("state", application, nil)),
		d.Set("published_at", utils.FormatTimeStampRFC3339(
			utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("publish_at", application, "").(string))/1000, false)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func GetAppPublishedApplicationByName(client *golangsdk.ServiceClient, appGroupId, appName string) (interface{}, error) {
	httpUrl := "v1/{project_id}/app-groups/{app_group_id}/apps?name={app_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{app_group_id}", appGroupId)
	// Fuzzy matching for name parameter.
	getPath = strings.ReplaceAll(getPath, "{app_name}", appName)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	express := fmt.Sprintf("items[?name=='%s']|[0]", appName)
	application := utils.PathSearch(express, respBody, nil)
	// When the application or application group does not exist, the status code is `200`.
	if application == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/app-groups/{app_group_id}/apps",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the application (%s) does not exist", appName)),
			},
		}
	}
	return application, nil
}

func resourceAppApplicationPublishmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	if d.HasChanges("name", "sandbox_enable", "version", "execute_path", "work_path", "command_param", "description", "status") {
		appName := d.Get("name")
		updateOpt := map[string]interface{}{
			"name":           appName,
			"execute_path":   utils.ValueIgnoreEmpty(d.Get("execute_path")),
			"sandbox_enable": d.Get("sandbox_enable"),
			"version":        d.Get("version"),
			"work_path":      d.Get("work_path"),
			"command_param":  d.Get("command_param"),
			"description":    d.Get("description"),
			"state":          utils.ValueIgnoreEmpty(d.Get("status")),
		}

		appGroupId := d.Get("app_group_id").(string)
		if err := updateAppPublishedApplication(client, appGroupId, d.Id(), updateOpt); err != nil {
			return diag.Errorf("error updating application (%s): %s", appName.(string), err)
		}
	}

	return resourceAppApplicationPublishmentRead(ctx, d, meta)
}

func updateAppPublishedApplication(client *golangsdk.ServiceClient, appGroupId, appId string, params map[string]interface{}) error {
	httpUrl := "v1/{project_id}/app-groups/{app_group_id}/apps/{app_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{app_group_id}", appGroupId)
	updatePath = strings.ReplaceAll(updatePath, "{app_id}", appId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(params),
	}

	_, err := client.Request("PATCH", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	return nil
}

func resourceAppApplicationPublishmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	httpUrl := "v1/{project_id}/app-groups/{app_group_id}/apps/batch-unpublish"
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{app_group_id}", d.Get("app_group_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"ids": []string{d.Id()},
		},
	}
	// In any case, the interface response status code is `200`.
	appName := d.Get("name").(string)
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting application (%s): %s", appName, err)
	}

	return nil
}

func resourceAppApplicationPublishmentImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<app_group_id>/<name>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("app_group_id", parts[0]),
		d.Set("name", parts[1]),
	)

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	application, err := GetAppPublishedApplicationByName(client, parts[0], parts[1])
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error getting published application from API response: %s", err)
	}

	applicationId := utils.PathSearch("id", application, "").(string)
	if applicationId == "" {
		return []*schema.ResourceData{d}, fmt.Errorf("unable to find application ID from API response")
	}

	d.SetId(applicationId)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
