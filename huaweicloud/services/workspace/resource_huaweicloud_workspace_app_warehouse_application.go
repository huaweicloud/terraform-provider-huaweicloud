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

// @API Workspace POST /v1/{project_id}/app-warehouse/apps
// @API Workspace GET /v1/{project_id}/app-warehouse/apps
// @API Workspace PATCH /v1/{project_id}/app-warehouse/apps/{id}
// @API Workspace DELETE /v1/{project_id}/app-warehouse/apps/{id}
func ResourceAppWarehouseApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppWarehouseApplicationCreate,
		ReadContext:   resourceAppWarehouseApplicationRead,
		UpdateContext: resourceAppWarehouseApplicationUpdate,
		DeleteContext: resourceAppWarehouseApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the application.`,
			},
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The category of the application.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The operating system type of the application.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version of the application.`,
			},
			"version_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version name of the application.`,
			},
			"file_store_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The storage path of the OBS bucket where the application is located.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the application.`,
			},
			"icon": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The icon of the application.`,
			},
			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The record ID of the application.`,
			},
		},
	}
}

func resourceAppWarehouseApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	httpUrl := "v1/{project_id}/app-warehouse/apps"
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateWarehouseApplicationBodyParams(d)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating warehouse application: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	applicationId := utils.PathSearch("app_id", respBody, "").(string)
	if applicationId == "" {
		return diag.Errorf("unable to find application ID from API response")
	}
	d.SetId(applicationId)
	return resourceAppWarehouseApplicationRead(ctx, d, meta)
}

func buildCreateWarehouseApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"app_name":           d.Get("name"),
		"app_category":       d.Get("category"),
		"os_type":            d.Get("os_type"),
		"version_id":         d.Get("version"),
		"version_name":       d.Get("version_name"),
		"appfile_store_path": d.Get("file_store_path"),
		"app_description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"app_icon":           utils.ValueIgnoreEmpty(d.Get("icon")),
	}
}

func resourceAppWarehouseApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	application, err := GetWarehouseApplicationById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace APP warehouse application")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("app_name", application, nil)),
		d.Set("category", utils.PathSearch("app_category", application, nil)),
		d.Set("os_type", utils.PathSearch("os_type", application, nil)),
		d.Set("version", utils.PathSearch("version_id", application, nil)),
		d.Set("version_name", utils.PathSearch("version_name", application, nil)),
		d.Set("file_store_path", utils.PathSearch("appfile_store_path", application, nil)),
		d.Set("description", utils.PathSearch("app_description", application, nil)),
		d.Set("record_id", utils.PathSearch("id", application, nil)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// GetWarehouseApplicationById is a method used to get application detail by application ID.
func GetWarehouseApplicationById(client *golangsdk.ServiceClient, applicationId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/app-warehouse/apps?app_id={app_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{app_id}", applicationId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving application (%s) from warehouse: %s", applicationId, err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	application := utils.PathSearch("items|[0]", respBody, nil)
	if application == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return application, nil
}

func resourceAppWarehouseApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	if d.HasChanges("name", "category", "os_type", "version", "version_name", "description", "icon") {
		httpUrl := "v1/{project_id}/app-warehouse/apps/{id}"
		applicationId := d.Id()
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{id}", applicationId)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateWarehouseApplicationBodyParams(d)),
		}

		_, err := client.Request("PATCH", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating warehouse application (%s): %s", applicationId, err)
		}
	}

	return resourceAppWarehouseApplicationRead(ctx, d, meta)
}

func buildUpdateWarehouseApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"app_name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"app_category":    utils.ValueIgnoreEmpty(d.Get("category")),
		"os_type":         utils.ValueIgnoreEmpty(d.Get("os_type")),
		"version_id":      utils.ValueIgnoreEmpty(d.Get("version")),
		"version_name":    utils.ValueIgnoreEmpty(d.Get("version_name")),
		"app_description": d.Get("description"),
		"app_icon":        d.Get("icon"),
	}
}

func resourceAppWarehouseApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	httpUrl := "v1/{project_id}/app-warehouse/apps/{id}"
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	// Although the deletion result of the main region shows that the interface returns a 200 status code when
	// deleting a non-existent warehouse application, in order to avoid the possible return of a 404 status code in the
	// future, the CheckDeleted design is retained here.
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "application of Workspace APP warehouse")
	}

	return nil
}
