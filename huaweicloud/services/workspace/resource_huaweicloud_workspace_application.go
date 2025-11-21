package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/app-center/apps
// @API Workspace GET /v1/{project_id}/app-center/apps
// @API Workspace PATCH /v1/{project_id}/app-center/apps/{id}
// @API Workspace DELETE /v1/{project_id}/app-center/apps/{id}
func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the application.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version of the application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The description of the application.`,
			},
			"authorization_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The authorization type of the application.`,
			},
			"application_file_store": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        workspaceApplicationFileStore(),
				Description: `The file store configuration of the application.`,
			},
			"install_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The installation type of the application.`,
			},
			"support_os": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The supported operating system of the application.`,
			},
			"catalog_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The catalog ID of the application.`,
			},

			// Optional parameters.
			"application_icon_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The icon URL of the application.`,
			},
			"install_command": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The installation command of the application.`,
			},
			"uninstall_command": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The uninstallation command of the application.`,
			},
			"install_info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The installation information of the application.`,
			},
			"reserve_obs_file": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(_, _, _ string, _ *schema.ResourceData) bool {
					return true
				},
				Description: `Whether to delete the installation package in the OBS bucket.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the application.`,
			},
			"application_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source of the application.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the application, in UTC format.`,
			},
			"catalog": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The catalog name of the application.`,
			},
		},
	}
}

func workspaceApplicationFileStore() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"store_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The store type of the application file.`,
			},
			"bucket_store": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        workspaceApplicationBucketStore(),
				Description: `The OBS bucket store configuration.`,
			},
			"file_link": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The external file link.`,
			},
		},
	}
}

func workspaceApplicationBucketStore() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the OBS bucket.`,
			},
			"bucket_file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The file path in the OBS bucket.`,
			},
		},
	}
}

func buildCreateApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required.
		"name":               d.Get("name"),
		"version":            d.Get("version"),
		"description":        d.Get("description"),
		"authorization_type": d.Get("authorization_type"),
		"app_file_store":     buildApplicationFileStore(d.Get("application_file_store.0")),
		"install_type":       d.Get("install_type"),
		"support_os":         d.Get("support_os"),
		"catalog_id":         d.Get("catalog_id"),
		// Optional.
		"app_icon_url":      utils.ValueIgnoreEmpty(d.Get("application_icon_url")),
		"install_command":   utils.ValueIgnoreEmpty(d.Get("install_command")),
		"uninstall_command": utils.ValueIgnoreEmpty(d.Get("uninstall_command")),
		"install_info":      utils.ValueIgnoreEmpty(d.Get("install_info")),
	}
}

func buildApplicationFileStore(fileStore interface{}) map[string]interface{} {
	if fileStore == nil {
		return nil
	}

	return map[string]interface{}{
		"store_type": utils.PathSearch("store_type", fileStore, nil),
		"bucket_store": utils.ValueIgnoreEmpty(buildApplicationBucketStore(
			utils.PathSearch("bucket_store", fileStore, make([]interface{}, 0)).([]interface{}))),
		"file_link": utils.ValueIgnoreEmpty(utils.PathSearch("file_link", fileStore, nil)),
	}
}

func buildApplicationBucketStore(bucketStores []interface{}) map[string]interface{} {
	if len(bucketStores) == 0 {
		return nil
	}

	return map[string]interface{}{
		"bucket_name":      utils.PathSearch("bucket_name", bucketStores[0], nil),
		"bucket_file_path": utils.PathSearch("bucket_file_path", bucketStores[0], nil),
	}
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/app-center/apps"
		region  = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateApplicationBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating Workspace application: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	applicationId := utils.PathSearch("id", respBody, "").(string)
	if applicationId == "" {
		return diag.Errorf("unable to find application ID from API response")
	}
	d.SetId(applicationId)

	return resourceApplicationRead(ctx, d, meta)
}

func buildApplicationQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	return res
}

func listApplications(client *golangsdk.ServiceClient, d ...*schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-center/apps?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	if len(d) != 0 {
		listPathWithLimit += buildApplicationQueryParams(d[0])
	}

	opt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%v", listPathWithLimit, strconv.Itoa(offset))
		requestResp, err := client.Request("GET", listPathWithOffset, opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		applications := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, applications...)
		if len(applications) < limit {
			break
		}
		offset += len(applications)
	}

	return result, nil
}

// GetApplicationById is a method is used to get the application.
func GetApplicationById(client *golangsdk.ServiceClient, applicationId string) (interface{}, error) {
	applications, err := listApplications(client)
	if err != nil {
		return nil, err
	}

	application := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", applicationId), applications, nil)
	if application == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return application, nil
}

func resourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		applicationId = d.Id()
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	application, err := GetApplicationById(client, applicationId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying Workspace application (%s)", applicationId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", application, nil)),
		d.Set("version", utils.PathSearch("version", application, nil)),
		d.Set("description", utils.PathSearch("description", application, nil)),
		d.Set("authorization_type", utils.PathSearch("authorization_type", application, nil)),
		d.Set("application_file_store", flattenApplicationFileStore(utils.PathSearch("app_file_store", application, nil))),
		d.Set("install_type", utils.PathSearch("install_type", application, nil)),
		d.Set("support_os", utils.PathSearch("support_os", application, nil)),
		d.Set("catalog_id", utils.PathSearch("catalog_id", application, nil)),
		d.Set("application_icon_url", utils.PathSearch("app_icon_url", application, nil)),
		d.Set("install_command", utils.PathSearch("install_command", application, nil)),
		d.Set("uninstall_command", utils.PathSearch("uninstall_command", application, nil)),
		d.Set("install_info", utils.PathSearch("install_info", application, nil)),
		d.Set("reserve_obs_file", true),
		d.Set("status", utils.PathSearch("status", application, nil)),
		d.Set("application_source", utils.PathSearch("application_source", application, nil)),
		d.Set("create_time", utils.PathSearch("create_time", application, nil)),
		d.Set("catalog", utils.PathSearch("catalog", application, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApplicationFileStore(fileStore interface{}) []interface{} {
	if fileStore == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"store_type": utils.PathSearch("store_type", fileStore, nil),
			"bucket_store": utils.ValueIgnoreEmpty(flattenApplicationBucketStore(
				utils.PathSearch("bucket_store", fileStore, nil))),
			"file_link": utils.ValueIgnoreEmpty(utils.PathSearch("file_link", fileStore, nil)),
		},
	}
}

func flattenApplicationBucketStore(bucketStore interface{}) []interface{} {
	if bucketStore == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"bucket_name":      utils.PathSearch("bucket_name", bucketStore, nil),
			"bucket_file_path": utils.PathSearch("bucket_file_path", bucketStore, nil),
		},
	}
}

func buildUpdateApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":               d.Get("name"),
		"version":            d.Get("version"),
		"description":        d.Get("description"),
		"authorization_type": d.Get("authorization_type"),
		"app_file_store":     buildApplicationFileStore(d.Get("application_file_store.0")),
		"install_type":       d.Get("install_type"),
		"support_os":         d.Get("support_os"),
		"catalog_id":         d.Get("catalog_id"),
		"app_icon_url":       utils.ValueIgnoreEmpty(d.Get("application_icon_url")),
		"install_command":    utils.ValueIgnoreEmpty(d.Get("install_command")),
		"uninstall_command":  utils.ValueIgnoreEmpty(d.Get("uninstall_command")),
		"install_info":       utils.ValueIgnoreEmpty(d.Get("install_info")),
		"status":             utils.ValueIgnoreEmpty(d.Get("status")),
	}
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		appplicationId = d.Id()
		updatePath     = "v1/{project_id}/app-center/apps/{application_id}"
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	updatePath = client.Endpoint + updatePath
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{application_id}", appplicationId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateApplicationBodyParams(d),
	}

	_, err = client.Request("PATCH", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Workspace application: %s", err)
	}

	return resourceApplicationRead(ctx, d, meta)
}

func buildDeleteApplicationQueryParams(d *schema.ResourceData) string {
	res := ""
	v := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "reserve_obs_file")
	if v != nil {
		res = fmt.Sprintf("%s?reserve_obs_file=%v", res, v)
	}
	return res
}

func resourceApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		deletePath    = "v1/{project_id}/app-center/apps/{application_id}"
		applicationId = d.Id()
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	deletePath = client.Endpoint + deletePath
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{application_id}", applicationId)
	deletePath += buildDeleteApplicationQueryParams(d)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// Deleting a non-existent application returns a 200 status code.
		return diag.Errorf("error deleting Workspace application: %s", err)
	}

	return nil
}
