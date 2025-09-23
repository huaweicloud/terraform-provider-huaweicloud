package secmaster

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

// All body parameters defined in the current API documentation are optional, but leaving some fields blank in actual
// testing will result in an API error.
// Currently, the resource fields are designed to be consistent with the API documentation.

// After testing, some fields will report errors when editing. However, the lack of documentation makes it unclear which
// fields cannot be edited. Currently, fields that do not support editing rely on API errors, and resources are not
// subject to special restrictions.

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/catalogues
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/catalogues/search
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/catalogues/{catalogue_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/catalogues
func ResourceCatalogue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCatalogueCreate,
		UpdateContext: resourceCatalogueUpdate,
		ReadContext:   resourceCatalogueRead,
		DeleteContext: resourceCatalogueDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCatalogueImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"workspace_id"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			"parent_catalogue": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the name of the first-level directory.",
			},
			// Fields `parent_alias_en` and `parent_alias_zh` are misspelled in the API documentation.
			// The actual names defined in the current schema are valid.
			"parent_alias_en": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the English alias of the first-level directory.",
			},
			"parent_alias_zh": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the Chinese alias of the first-level directory.",
			},
			"second_catalogue": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the name of the second-level directory.",
			},
			"second_alias_en": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the English alias of the second-level directory.",
			},
			"second_alias_zh": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the Chinese alias of the second-level directory.",
			},
			"layout_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the ID of the layout.",
			},
			"layout_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the name of the layout.",
			},
			"catalogue_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the address of the directory.",
			},
			"publisher_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the name of the publisher.",
			},
			// The field `second_catalogue_code` is not returned in API response body. So `Computed` is not added.
			"second_catalogue_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the code of the second-level directory.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"is_card_area": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag indicating whether to display the card area.",
			},
			"is_display": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag indicating whether to display the directory.",
			},
			"is_landing_page": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag indicating whether it is a landing page.",
			},
			"is_navigation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag indicating whether to display the breadcrumb navigation.",
			},
			"catalogue_status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag indicating whether it is a built-in directory.",
			},
		},
	}
}

func buildCreateCatalogueBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"parent_catalogue":      utils.ValueIgnoreEmpty(d.Get("parent_catalogue")),
		"parent_alias_en":       utils.ValueIgnoreEmpty(d.Get("parent_alias_en")),
		"parent_alias_zh":       utils.ValueIgnoreEmpty(d.Get("parent_alias_zh")),
		"second_catalogue":      utils.ValueIgnoreEmpty(d.Get("second_catalogue")),
		"second_alias_en":       utils.ValueIgnoreEmpty(d.Get("second_alias_en")),
		"second_alias_zh":       utils.ValueIgnoreEmpty(d.Get("second_alias_zh")),
		"second_catalogue_code": utils.ValueIgnoreEmpty(d.Get("second_catalogue_code")),
		"layout_id":             utils.ValueIgnoreEmpty(d.Get("layout_id")),
		"layout_name":           utils.ValueIgnoreEmpty(d.Get("layout_name")),
		"catalogue_address":     utils.ValueIgnoreEmpty(d.Get("catalogue_address")),
		"publisher_name":        utils.ValueIgnoreEmpty(d.Get("publisher_name")),
	}
}

func resourceCatalogueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/catalogues"
		product     = "secmaster"
		workspaceID = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateCatalogueBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster catalogue: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster catalogue: ID is not found in API response")
	}
	d.SetId(id)

	return resourceCatalogueRead(ctx, d, meta)
}

// The service provides two query list APIs. Currently, only the `/search` API response body can query the catalogue ID.
func ReadCatalogueDetail(client *golangsdk.ServiceClient, workspaceID, catalogueID string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/catalogues/search"
		offset  = 0
		allData = make([]interface{}, 0)
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: make(map[string]interface{}),
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s?offset=%d", requestPath, offset)
		resp, err := client.Request("POST", requestPathWithOffset, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		data := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(data) == 0 {
			break
		}

		// There is a problem with the API's paging logic. The offset can always be queried for data.
		// It is necessary to compare the `total` value to safely jump out of the loop logic.
		allData = append(allData, data...)
		totalCount := int(utils.PathSearch("total", respBody, float64(0)).(float64))
		if len(allData) >= totalCount {
			break
		}

		offset += len(data)
	}

	dataResp := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", catalogueID), allData, nil)
	if dataResp == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return dataResp, nil
}

func resourceCatalogueRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := ReadCatalogueDetail(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster catalogue")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("parent_catalogue", utils.PathSearch("parent_catalogue", respBody, nil)),
		d.Set("second_catalogue", utils.PathSearch("second_catalogue", respBody, nil)),
		d.Set("catalogue_status", utils.PathSearch("catalogue_status", respBody, nil)),
		d.Set("catalogue_address", utils.PathSearch("catalogue_address", respBody, nil)),
		d.Set("layout_id", utils.PathSearch("layout_id", respBody, nil)),
		d.Set("layout_name", utils.PathSearch("layout_name", respBody, nil)),
		d.Set("publisher_name", utils.PathSearch("publisher_name", respBody, nil)),
		d.Set("is_card_area", utils.PathSearch("is_card_area", respBody, nil)),
		d.Set("is_display", utils.PathSearch("is_display", respBody, nil)),
		d.Set("is_landing_page", utils.PathSearch("is_landing_page", respBody, nil)),
		d.Set("is_navigation", utils.PathSearch("is_navigation", respBody, nil)),
		d.Set("parent_alias_en", utils.PathSearch("parent_alias_en", respBody, nil)),
		d.Set("parent_alias_zh", utils.PathSearch("parent_alias_zh", respBody, nil)),
		d.Set("second_alias_en", utils.PathSearch("second_alias_en", respBody, nil)),
		d.Set("second_alias_zh", utils.PathSearch("second_alias_zh", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateCatalogueBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"parent_catalogue":      utils.ValueIgnoreEmpty(d.Get("parent_catalogue")),
		"parent_alias_en":       utils.ValueIgnoreEmpty(d.Get("parent_alias_en")),
		"parent_alias_zh":       utils.ValueIgnoreEmpty(d.Get("parent_alias_zh")),
		"second_catalogue":      utils.ValueIgnoreEmpty(d.Get("second_catalogue")),
		"second_alias_en":       utils.ValueIgnoreEmpty(d.Get("second_alias_en")),
		"second_alias_zh":       utils.ValueIgnoreEmpty(d.Get("second_alias_zh")),
		"second_catalogue_code": utils.ValueIgnoreEmpty(d.Get("second_catalogue_code")),
		"layout_id":             utils.ValueIgnoreEmpty(d.Get("layout_id")),
		"layout_name":           utils.ValueIgnoreEmpty(d.Get("layout_name")),
		"catalogue_address":     utils.ValueIgnoreEmpty(d.Get("catalogue_address")),
		"publisher_name":        utils.ValueIgnoreEmpty(d.Get("publisher_name")),
	}
}

func resourceCatalogueUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/catalogues/{catalogue_id}"
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{catalogue_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateCatalogueBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster catalogue: %s", err)
	}

	return resourceCatalogueRead(ctx, d, meta)
}

func resourceCatalogueDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/catalogues"
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"batch_ids": []string{d.Id()},
			"is_delete": true,
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster catalogue: %s", err)
	}

	return nil
}

func resourceCatalogueImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want <workspace_id>/<catalogue_id>, but got %s", importId)
	}

	d.SetId(importIdParts[1])
	return []*schema.ResourceData{d}, d.Set("workspace_id", importIdParts[0])
}
