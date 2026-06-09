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

var nonUpdatableParamsComponentTemplate = []string{"workspace_id"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/components/template
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/components/template/{template_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/components/template/{template_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/components/template/{template_id}
func ResourceComponentTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComponentTemplateCreate,
		ReadContext:   resourceComponentTemplateRead,
		UpdateContext: resourceComponentTemplateUpdate,
		DeleteContext: resourceComponentTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceComponentTemplateImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsComponentTemplate),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_config": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateComponentTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"component_id":  d.Get("component_id"),
		"template_name": d.Get("template_name"),
		"task_config":   d.Get("task_config"),
	}

	return bodyParams
}

func resourceComponentTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		workspaceId   = d.Get("workspace_id").(string)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/components/template"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         buildCreateComponentTemplateBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster component template: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateId := utils.PathSearch("data.id", respBody, "").(string)
	if templateId == "" {
		return diag.Errorf("error creating SecMaster component template: unable to find component template ID")
	}

	d.SetId(templateId)

	return resourceComponentTemplateRead(ctx, d, meta)
}

func resourceComponentTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetComponentTemplateInfo(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "code", "SecMaster.20040220"),
			"error retrieving SecMaster component template",
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("component_id", utils.PathSearch("data.component_id", respBody, nil)),
		d.Set("template_name", utils.PathSearch("data.template_name", respBody, nil)),
		d.Set("task_config", utils.PathSearch("data.task_config", respBody, nil)),
		d.Set("create_time", utils.PathSearch("data.create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("data.update_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetComponentTemplateInfo(client *golangsdk.ServiceClient, workspaceId, templateId string) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/soc/components/template/{template_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{template_id}", templateId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func buildUpdateComponentTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"component_id":  d.Get("component_id"),
		"template_name": d.Get("template_name"),
		"task_config":   d.Get("task_config"),
	}

	return bodyParams
}

func resourceComponentTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/components/template/{template_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{template_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         buildUpdateComponentTemplateBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster component template: %s", err)
	}

	return resourceComponentTemplateRead(ctx, d, meta)
}

func resourceComponentTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/components/template/{template_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{template_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster component template: %s", err)
	}

	return nil
}

func resourceComponentTemplateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
