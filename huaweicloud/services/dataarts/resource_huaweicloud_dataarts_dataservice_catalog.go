package dataarts

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

var CatalogResourceNotFoundCodes = []string{
	"DLM.4001", // Instance or workspace does not exist.
	"DLM.4205", // Catalog does not exist
}

// @API DataArtsStudio POST /v1/{project_id}/service/servicecatalogs
// @API DataArtsStudio GET /v1/{project_id}/service/servicecatalogs/{catalog_id}
// @API DataArtsStudio PUT /v1/{project_id}/service/servicecatalogs/{catalog_id}
// @API DataArtsStudio POST /v1/{project_id}/service/servicecatalogs/{catalog_id}/move
// @API DataArtsStudio POST /v1/{project_id}/service/servicecatalogs/batch-delete
func ResourceDatatServiceCatalog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatatServiceCatalogCreate,
		ReadContext:   resourceDatatServiceCatalogRead,
		UpdateContext: resourceDatatServiceCatalogUpdate,
		DeleteContext: resourceDatatServiceCatalogDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataServiceCatalogImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the catalog is located.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the workspace to which the catalog belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the catalog.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the catalog.`,
			},
			"dlm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The type of DLM engine.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0", // Root path.
				Description: `The ID of the parent catalog for current catalog.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The path of current catalog.`,
			},
			"catalog_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of sub-catalogs in the current catalog.`,
			},
			"api_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of APIs in the current catalog.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the catalog.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the catalog.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the catalog.`,
			},
			"update_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user who latest updated the catalog.`,
			},
		},
	}
}

func buildCreateDatatServiceCatalogBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"pid":         utils.ValueIgnoreEmpty(d.Get("parent_id")),
	}
	return bodyParams
}

func resourceDatatServiceCatalogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/service/servicecatalogs"
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"Dlm-Type":     d.Get("dlm_type").(string),
		},
		JSONBody: utils.RemoveNil(buildCreateDatatServiceCatalogBodyParams(d)),
	}

	createAppResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating catalog: %s", err)
	}

	createAppRespBody, err := utils.FlattenResponse(createAppResp)
	if err != nil {
		return diag.FromErr(err)
	}

	catalogId := utils.PathSearch("catalog_id", createAppRespBody, "").(string)
	if catalogId == "" {
		return diag.Errorf("unable to find the DataArts DataService catalog ID from the API response")
	}
	d.SetId(catalogId)

	return resourceDatatServiceCatalogRead(ctx, d, meta)
}

func resourceDatatServiceCatalogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/service/servicecatalogs/{catalog_id}"
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{catalog_id}", d.Id())

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"Dlm-Type":     d.Get("dlm_type").(string),
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", CatalogResourceNotFoundCodes...),
			"error retrieving catalog")
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("path", utils.PathSearch("path", respBody, nil)),
		d.Set("catalog_total", utils.PathSearch("catalog_total", respBody, nil)),
		d.Set("api_total", utils.PathSearch("api_total", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("create_user", utils.PathSearch("create_user", respBody, nil)),
		d.Set("update_user", utils.PathSearch("update_user", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDatatServiceCatalogBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
	return bodyParams
}

func updateCatalogBasicConfigs(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl   = "v1/{project_id}/service/servicecatalogs/{catalog_id}"
		catalogId = d.Id()
	)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{catalog_id}", catalogId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"Dlm-Type":     d.Get("dlm_type").(string),
		},
		JSONBody: buildUpdateDatatServiceCatalogBodyParams(d),
	}

	_, err := client.Request("PUT", createPath, &opt)
	if err != nil {
		return fmt.Errorf("error updating catalog (%s): %s", catalogId, err)
	}
	return nil
}

func updateCatalogParentId(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl   = "v1/{project_id}/service/servicecatalogs/{catalog_id}/move"
		catalogId = d.Id()
	)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{catalog_id}", catalogId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"Dlm-Type":     d.Get("dlm_type").(string),
		},
		JSONBody: map[string]interface{}{
			"target_pid": d.Get("parent_id"),
		},
	}

	_, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return fmt.Errorf("error moving current catalog to a new path: %s", err)
	}
	return nil
}

func resourceDatatServiceCatalogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if d.HasChanges("name", "description") {
		err = updateCatalogBasicConfigs(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("parent_id") {
		err = updateCatalogParentId(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDatatServiceCatalogRead(ctx, d, meta)
}

func resourceDatatServiceCatalogDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/service/servicecatalogs/batch-delete"
		workspaceId = d.Get("workspace_id").(string)
		catalogId   = d.Id()
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	// Due to the lack of an interface for deleting a single directory, the corresponding function can only be
	// implemented through the batch deletion interface, and the batch deletion API cannot be sent to the server
	// at the same time, which will cause an error to be returned.
	config.MutexKV.Lock(workspaceId)
	defer config.MutexKV.Unlock(workspaceId)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     d.Get("dlm_type").(string),
		},
		JSONBody: map[string]interface{}{
			"ids": []string{catalogId},
		},
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error deleting catalog (%s): %s", catalogId, err)
	}
	return nil
}

func resourceDataServiceCatalogImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 && len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be '<workspace_id>/<dlm_type>/<id>' or "+
			"'<workspace_id>/<id>', but got '%s'", importedId)
	}

	mErr := multierror.Append(nil, d.Set("workspace_id", parts[0]))
	if len(parts) == 2 {
		d.SetId(parts[1])
	}
	if len(parts) == 3 {
		mErr = multierror.Append(mErr, d.Set("dlm_type", parts[1]))
		d.SetId(parts[2])
	}

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
