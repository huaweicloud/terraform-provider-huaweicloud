// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DCS
// ---------------------------------------------------------------

package dcs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS POST /v2/{project_id}/config-templates
// @API DCS PUT /v2/{project_id}/config-templates/{template_id}
// @API DCS GET /v2/{project_id}/config-templates/{template_id}
// @API DCS DELETE /v2/{project_id}/config-templates/{template_id}
func ResourceCustomTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomTemplateCreate,
		UpdateContext: resourceCustomTemplateUpdate,
		ReadContext:   resourceCustomTemplateRead,
		DeleteContext: resourceCustomTemplateDelete,
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
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the source template.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the source template.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the template.`,
			},
			"params": {
				Type:        schema.TypeSet,
				Elem:        customTemplateParamSchema(),
				Required:    true,
				Description: `Specifies the list of the template params.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the template.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the template.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the cache engine.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the cache engine version.`,
			},
			"cache_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DCS instance type.`,
			},
			"product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the product edition.`,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the storage type.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the custom template is created.`,
			},
		},
	}
}

func customTemplateParamSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"param_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Indicates the name of the param.`,
			},
			"param_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Indicates the value of the param.`,
			},
			"param_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the param.`,
			},
			"default_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the default value of the param.`,
			},
			"value_range": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the value range of the param.`,
			},
			"value_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the value type of the param.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the param.`,
			},
			"need_restart": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the DCS instance need restart.`,
			},
		},
	}
	return &sc
}

func resourceCustomTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createCustomTemplate: create DCS custom template
	var (
		createCustomTemplateHttpUrl = "v2/{project_id}/config-templates"
		createCustomTemplateProduct = "dcs"
	)
	createCustomTemplateClient, err := cfg.NewServiceClient(createCustomTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	createCustomTemplatePath := createCustomTemplateClient.Endpoint + createCustomTemplateHttpUrl
	createCustomTemplatePath = strings.ReplaceAll(createCustomTemplatePath, "{project_id}",
		createCustomTemplateClient.ProjectID)

	createCustomTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createCustomTemplateOpt.JSONBody = utils.RemoveNil(buildCreateCustomTemplateBodyParams(d))
	createCustomTemplateResp, err := createCustomTemplateClient.Request("POST", createCustomTemplatePath,
		&createCustomTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating DCS custom template: %s", err)
	}

	createCustomTemplateRespBody, err := utils.FlattenResponse(createCustomTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateId := utils.PathSearch("id", createCustomTemplateRespBody, "").(string)
	if templateId == "" {
		return diag.Errorf("unable to find the DCS custom template ID from the API response")
	}
	d.SetId(templateId)

	return resourceCustomTemplateRead(ctx, d, meta)
}

func buildCreateCustomTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_id": d.Get("template_id"),
		"type":        d.Get("source_type"),
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"params":      buildCreateCustomTemplateRequestBodyParam(d.Get("params").(*schema.Set).List()),
	}
	return bodyParams
}

func buildCreateCustomTemplateRequestBodyParam(rawParams interface{}) map[string]string {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make(map[string]string)
		for _, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[raw["param_name"].(string)] = raw["param_value"].(string)
			}
		}
		return rst
	}
	return nil
}

func resourceCustomTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCustomTemplate: Query DCS custom template
	var (
		getCustomTemplateHttpUrl = "v2/{project_id}/config-templates/{template_id}?type=user"
		getCustomTemplateProduct = "dcs"
	)
	getCustomTemplateClient, err := cfg.NewServiceClient(getCustomTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getCustomTemplatePath := getCustomTemplateClient.Endpoint + getCustomTemplateHttpUrl
	getCustomTemplatePath = strings.ReplaceAll(getCustomTemplatePath, "{project_id}", getCustomTemplateClient.ProjectID)
	getCustomTemplatePath = strings.ReplaceAll(getCustomTemplatePath, "{template_id}", d.Id())

	getCustomTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getCustomTemplateResp, err := getCustomTemplateClient.Request("GET", getCustomTemplatePath, &getCustomTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DCS.4988"),
			"error retrieving DCS custom template")
	}

	getCustomTemplateRespBody, err := utils.FlattenResponse(getCustomTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getCustomTemplateRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getCustomTemplateRespBody, nil)),
		d.Set("engine", utils.PathSearch("engine", getCustomTemplateRespBody, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", getCustomTemplateRespBody, nil)),
		d.Set("cache_mode", utils.PathSearch("cache_mode", getCustomTemplateRespBody, nil)),
		d.Set("product_type", utils.PathSearch("product_type", getCustomTemplateRespBody, nil)),
		d.Set("storage_type", utils.PathSearch("storage_type", getCustomTemplateRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getCustomTemplateRespBody, nil)),
		d.Set("params", flattenGetCustomTemplateResponseBodyParam(d, getCustomTemplateRespBody)),
		d.Set("created_at", utils.PathSearch("created_at", getCustomTemplateRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetCustomTemplateResponseBodyParam(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	paramsMap := buildParamsMap(d)
	curJson := utils.PathSearch("params", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		paramName := utils.PathSearch("param_name", v, "").(string)
		if !paramsMap[paramName] {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"param_id":      utils.PathSearch("param_id", v, nil),
			"param_name":    paramName,
			"param_value":   utils.PathSearch("param_value", v, nil),
			"default_value": utils.PathSearch("default_value", v, nil),
			"value_range":   utils.PathSearch("value_range", v, nil),
			"value_type":    utils.PathSearch("value_type", v, nil),
			"description":   utils.PathSearch("description", v, nil),
			"need_restart":  utils.PathSearch("need_restart", v, nil),
		})
	}
	return rst
}

func buildParamsMap(d *schema.ResourceData) map[string]bool {
	params := d.Get("params").(*schema.Set).List()
	paramsMap := make(map[string]bool)
	for _, param := range params {
		if v, ok := param.(map[string]interface{}); ok {
			paramsMap[v["param_name"].(string)] = true
		}
	}
	return paramsMap
}

func resourceCustomTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateCustomTemplateChanges := []string{
		"name",
		"description",
		"params",
	}

	if d.HasChanges(updateCustomTemplateChanges...) {
		// updateCustomTemplate: update DCS custom template
		var (
			updateCustomTemplateHttpUrl = "v2/{project_id}/config-templates/{template_id}"
			updateCustomTemplateProduct = "dcs"
		)
		updateCustomTemplateClient, err := cfg.NewServiceClient(updateCustomTemplateProduct, region)
		if err != nil {
			return diag.Errorf("error creating DCS client: %s", err)
		}

		updateCustomTemplatePath := updateCustomTemplateClient.Endpoint + updateCustomTemplateHttpUrl
		updateCustomTemplatePath = strings.ReplaceAll(updateCustomTemplatePath, "{project_id}",
			updateCustomTemplateClient.ProjectID)
		updateCustomTemplatePath = strings.ReplaceAll(updateCustomTemplatePath, "{template_id}", d.Id())

		updateCustomTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}

		updateCustomTemplateOpt.JSONBody = utils.RemoveNil(buildUpdateCustomTemplateBodyParams(d))
		_, err = updateCustomTemplateClient.Request("PUT", updateCustomTemplatePath, &updateCustomTemplateOpt)
		if err != nil {
			return diag.Errorf("error updating DCS custom template: %s", err)
		}
	}
	return resourceCustomTemplateRead(ctx, d, meta)
}

func buildUpdateCustomTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"params":      buildUpdateCustomTemplateRequestBodyParam(d.Get("params").(*schema.Set).List()),
	}
	return bodyParams
}

func buildUpdateCustomTemplateRequestBodyParam(rawParams interface{}) map[string]string {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make(map[string]string)
		for _, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[raw["param_name"].(string)] = raw["param_value"].(string)
			}
		}
		return rst
	}
	return nil
}

func resourceCustomTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteCustomTemplate: Delete DCS custom template
	var (
		deleteCustomTemplateHttpUrl = "v2/{project_id}/config-templates/{template_id}"
		deleteCustomTemplateProduct = "dcs"
	)
	deleteCustomTemplateClient, err := cfg.NewServiceClient(deleteCustomTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	deleteCustomTemplatePath := deleteCustomTemplateClient.Endpoint + deleteCustomTemplateHttpUrl
	deleteCustomTemplatePath = strings.ReplaceAll(deleteCustomTemplatePath, "{project_id}",
		deleteCustomTemplateClient.ProjectID)
	deleteCustomTemplatePath = strings.ReplaceAll(deleteCustomTemplatePath, "{template_id}", d.Id())

	deleteCustomTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteCustomTemplateClient.Request("DELETE", deleteCustomTemplatePath, &deleteCustomTemplateOpt)
	if err != nil {
		return diag.Errorf("error deleting DCS custom template: %s", err)
	}

	return nil
}
