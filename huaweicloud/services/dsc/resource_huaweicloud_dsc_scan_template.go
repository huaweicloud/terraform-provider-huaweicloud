package dsc

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

var nonUpdatableParamsScanTemplate = []string{"action", "add_built_in_rules", "origin_template_id"}

// @API DSC POST /v1/{project_id}/scan-templates
// @API DSC GET /v1/{project_id}/scan-templates
// @API DSC PUT /v1/{project_id}/scan-templates/{template_id}
// @API DSC DELETE /v1/{project_id}/scan-templates/{template_id}
func ResourceScanTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScanTemplateCreate,
		ReadContext:   resourceScanTemplateRead,
		UpdateContext: resourceScanTemplateUpdate,
		DeleteContext: resourceScanTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsScanTemplate),

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
				Description: "The template name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The template description.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The operation type, such as creating or copying a template.",
			},
			"add_built_in_rules": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to add built-in rules when creating the template.",
			},
			"origin_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The origin template ID, used when copying a template.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the template is the default template.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template category.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time.",
			},
		},
	}
}

func buildCreateScanTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_template_name": d.Get("name"),
		"description":       d.Get("description"),
	}

	if v, ok := d.GetOk("action"); ok {
		bodyParams["action"] = v
	}
	if v, ok := d.GetOk("add_built_in_rules"); ok {
		bodyParams["add_built_in_rules"] = v
	}
	if v, ok := d.GetOk("origin_template_id"); ok {
		bodyParams["origin_template_id"] = v
	}

	return bodyParams
}

func resourceScanTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		name    = d.Get("name").(string)
		httpUrl = "v1/{project_id}/scan-templates"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateScanTemplateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DSC scan template: %s", err)
	}

	respBody, err := GetScanTemplateInfo(client, "name", name)
	if err != nil {
		return diag.Errorf("error retrieving DSC scan template after creation: %s", err)
	}

	templateId := utils.PathSearch("id", respBody, "").(string)
	if templateId == "" {
		return diag.Errorf("error creating DSC scan template: ID is not found in API response")
	}

	d.SetId(templateId)

	return resourceScanTemplateRead(ctx, d, meta)
}

func GetScanTemplateInfo(client *golangsdk.ServiceClient, filterField, filterValue string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/scan-templates"
		offset  = 0
		limit   = 1000
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		templates := utils.PathSearch("scan_templates_list", respBody, make([]interface{}, 0)).([]interface{})
		for _, template := range templates {
			fieldValue := utils.PathSearch(filterField, template, "").(string)
			if fieldValue == filterValue {
				return template, nil
			}
		}

		if len(templates) < limit {
			break
		}

		offset += len(templates)
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceScanTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	respBody, err := GetScanTemplateInfo(client, "id", d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DSC scan template")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("is_default", utils.PathSearch("is_default", respBody, nil)),
		d.Set("category", utils.PathSearch("category", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateScanTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_name": d.Get("name"),
		"template_desc": d.Get("description"),
	}

	if v, ok := d.GetOk("is_default"); ok {
		bodyParams["is_default"] = v
	}

	return bodyParams
}

func resourceScanTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scan-templates/{template_id}"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{template_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateScanTemplateBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating DSC scan template: %s", err)
	}

	return resourceScanTemplateRead(ctx, d, meta)
}

func resourceScanTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scan-templates/{template_id}"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{template_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DSC scan template")
	}

	return nil
}
