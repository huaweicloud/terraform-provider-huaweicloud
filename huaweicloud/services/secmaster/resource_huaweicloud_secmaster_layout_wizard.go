package secmaster

import (
	"context"
	"fmt"
	"log"
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

var layoutWizardNonUpdatableParams = []string{"workspace_id", "layout_id"}

// Omitting parameters `id`, `layout_id`, `project_id`, and `workspace_id` in the body of the API creation,
// because they have no meaning or have already been used in the URL.

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}/wizards
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards/{wizard_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards/{wizard_id}
func ResourceLayoutWizard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLayoutWizardCreate,
		UpdateContext: resourceLayoutWizardUpdate,
		ReadContext:   resourceLayoutWizardRead,
		DeleteContext: resourceLayoutWizardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLayoutWizardImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(layoutWizardNonUpdatableParams),

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
			// Query API no return.
			"layout_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the layout ID.",
			},
			// Optional in the API documentation, required in reality.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the wizard name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the description.",
			},
			"wizard_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the layout wizard information.",
			},
			// In the API documentation, it is of type `boolean`, here it has been changed to type `string`.
			"is_binding": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the button is bound.",
			},
			"binding_button": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the binding buttons.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"button_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the button ID.",
						},
						"button_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the button name.",
						},
					},
				},
			},
			"boa_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the BOA version.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator ID.",
			},
			"en_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The English description.",
			},
			"en_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The English name.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time.",
			},
			"is_built_in": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the wizard is built-in.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SecMaster version.",
			},
		},
	}
}

func convertStringToBool(stringValue string) interface{} {
	if stringValue == "" {
		return nil
	}

	boolValue, err := strconv.ParseBool(stringValue)
	if err != nil {
		log.Printf("[ERROR] error converting string %s to boolean: %s", stringValue, err)
		return nil
	}

	return boolValue
}

func buildLayoutWizardBindingButtonsBodyParams(buttons []interface{}) []map[string]interface{} {
	if len(buttons) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(buttons))
	for _, v := range buttons {
		raw, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		button := map[string]interface{}{
			"button_id":   raw["button_id"],
			"button_name": utils.ValueIgnoreEmpty(raw["button_name"]),
		}

		result = append(result, utils.RemoveNil(button))
	}

	return result
}

func buildCreateLayoutWizardBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           utils.ValueIgnoreEmpty(d.Get("name")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"wizard_json":    utils.ValueIgnoreEmpty(d.Get("wizard_json")),
		"is_binding":     convertStringToBool(d.Get("is_binding").(string)),
		"binding_button": buildLayoutWizardBindingButtonsBodyParams(d.Get("binding_button").([]interface{})),
		"boa_version":    utils.ValueIgnoreEmpty(d.Get("boa_version")),
	}

	return bodyParams
}

func resourceLayoutWizardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}/wizards"
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
		layoutId    = d.Get("layout_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{layout_id}", layoutId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildCreateLayoutWizardBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster layout wizard: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster layout wizard: ID is not found in API response")
	}

	d.SetId(id)

	return resourceLayoutWizardRead(ctx, d, meta)
}

func resourceLayoutWizardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards/{wizard_id}"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{wizard_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"SecMaster.20097006"), "error retrieving SecMaster layout wizard")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	wizard := utils.PathSearch("data", respBody, nil)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("workspace_id", utils.PathSearch("workspace_id", wizard, nil)),
		d.Set("name", utils.PathSearch("name", wizard, nil)),
		d.Set("description", utils.PathSearch("description", wizard, nil)),
		d.Set("wizard_json", utils.PathSearch("wizard_json", wizard, nil)),
		d.Set("is_binding", flattenBoolToString(utils.PathSearch("is_binding", wizard, nil))),
		d.Set("binding_button", flattenLayoutWizardBindingButtons(wizard)),
		d.Set("boa_version", utils.PathSearch("boa_version", wizard, nil)),
		d.Set("create_time", utils.PathSearch("create_time", wizard, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", wizard, nil)),
		d.Set("en_description", utils.PathSearch("en_description", wizard, nil)),
		d.Set("en_name", utils.PathSearch("en_name", wizard, nil)),
		d.Set("update_time", utils.PathSearch("update_time", wizard, nil)),
		d.Set("is_built_in", utils.PathSearch("is_built_in", wizard, nil)),
		d.Set("version", utils.PathSearch("version", wizard, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBoolToString(respBool interface{}) string {
	if respBool == nil {
		return ""
	}

	boolValue, ok := respBool.(bool)
	if !ok {
		return ""
	}

	if boolValue {
		return "true"
	}

	return "false"
}

func flattenLayoutWizardBindingButtons(wizard interface{}) []interface{} {
	buttons := utils.PathSearch("binding_button", wizard, make([]interface{}, 0)).([]interface{})
	if len(buttons) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(buttons))
	for _, v := range buttons {
		result = append(result, map[string]interface{}{
			"button_id":   utils.PathSearch("button_id", v, nil),
			"button_name": utils.PathSearch("button_name", v, nil),
		})
	}

	return result
}

func buildUpdateLayoutWizardBodyParams(d *schema.ResourceData) map[string]interface{} {
	pojo := map[string]interface{}{
		"id":             d.Id(),
		"name":           d.Get("name"),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"wizard_json":    utils.ValueIgnoreEmpty(d.Get("wizard_json")),
		"is_binding":     convertStringToBool(d.Get("is_binding").(string)),
		"binding_button": buildLayoutWizardBindingButtonsBodyParams(d.Get("binding_button").([]interface{})),
		"boa_version":    utils.ValueIgnoreEmpty(d.Get("boa_version")),
	}

	return map[string]interface{}{
		"layout_wizard_update_pojo_list": []interface{}{pojo},
	}
}

func resourceLayoutWizardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildUpdateLayoutWizardBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster layout wizard: %s", err)
	}

	return resourceLayoutWizardRead(ctx, d, meta)
}

func resourceLayoutWizardDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards/{wizard_id}"
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{wizard_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster layout wizard: %s", err)
	}

	return nil
}

func resourceLayoutWizardImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want '<workspace_id>/<id>', but got '%s'", importId)
	}

	d.SetId(importIdParts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", importIdParts[0])
}
