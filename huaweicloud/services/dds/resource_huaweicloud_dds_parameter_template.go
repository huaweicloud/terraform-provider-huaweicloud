// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDS
// ---------------------------------------------------------------

package dds

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

// @API DDS POST /v3/{project_id}/configurations
// @API DDS DELETE /v3/{project_id}/configurations/{config_id}
// @API DDS GET /v3/{project_id}/configurations/{config_id}
// @API DDS PUT /v3/{project_id}/configurations/{config_id}
func ResourceDdsParameterTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsParameterTemplateCreate,
		UpdateContext: resourceDdsParameterTemplateUpdate,
		ReadContext:   resourceDdsParameterTemplateRead,
		DeleteContext: resourceDdsParameterTemplateDelete,
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
				Description: `Specifies the parameter template name.`,
			},
			"node_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the node type of parameter template.`,
			},
			"node_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the database version.`,
			},
			"parameter_values": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the mapping between parameter names and parameter values.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter template description.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Elem:        ParameterTemplateParameterSchema(),
				Computed:    true,
				Description: `Indicates the parameters defined by users based on the default parameter templates.`,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ParameterTemplateParameterSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the parameter name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the parameter value.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the parameter description.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the parameter type.`,
			},
			"value_range": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the value range.`,
			},
			"restart_required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the instance needs to be restarted.`,
			},
			"readonly": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the parameter is read-only.`,
			},
		},
	}
	return &sc
}

func resourceDdsParameterTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createParameterTemplate: create DDS parameter template
	var (
		createParameterTemplateHttpUrl = "v3/{project_id}/configurations"
		createParameterTemplateProduct = "dds"
	)
	createParameterTemplateClient, err := cfg.NewServiceClient(createParameterTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS Client: %s", err)
	}

	createParameterTemplatePath := createParameterTemplateClient.Endpoint + createParameterTemplateHttpUrl
	createParameterTemplatePath = strings.ReplaceAll(createParameterTemplatePath, "{project_id}",
		createParameterTemplateClient.ProjectID)

	createParameterTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createParameterTemplateOpt.JSONBody = utils.RemoveNil(buildCreateParameterTemplateBodyParams(d))
	createParameterTemplateResp, err := createParameterTemplateClient.Request("POST",
		createParameterTemplatePath, &createParameterTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating DDS parameter template: %s", err)
	}

	createParameterTemplateRespBody, err := utils.FlattenResponse(createParameterTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("configuration.id", createParameterTemplateRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find ID from the API response")
	}
	d.SetId(id)

	return resourceDdsParameterTemplateRead(ctx, d, meta)
}

func buildCreateParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	datastoreParams := map[string]interface{}{
		"node_type": utils.ValueIgnoreEmpty(d.Get("node_type")),
		"version":   utils.ValueIgnoreEmpty(d.Get("node_version")),
	}
	bodyParams := map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(d.Get("name")),
		// this param can be empty
		"parameter_values": d.Get("parameter_values"),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"datastore":        datastoreParams,
	}
	return bodyParams
}

func resourceDdsParameterTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS Client: %s", err)
	}

	if err := updateParameterTemplate(client, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceDdsParameterTemplateRead(ctx, d, meta)
}

func updateParameterTemplate(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v3/{project_id}/configurations/{config_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{config_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateOpt.JSONBody = utils.RemoveNil(buildUpdateParameterTemplateBodyParams(d))
	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating DDS parameter template: %s", err)
	}

	return nil
}

func buildUpdateParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             utils.ValueIgnoreEmpty(d.Get("name")),
		"parameter_values": utils.ValueIgnoreEmpty(d.Get("parameter_values")),
		"description":      d.Get("description"),
	}
	return bodyParams
}

func resourceDdsParameterTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getParameterTemplate: Query DDS parameter template
	var (
		getParameterTemplateHttpUrl = "v3/{project_id}/configurations/{config_id}"
		getParameterTemplateProduct = "dds"
	)
	getParameterTemplateClient, err := cfg.NewServiceClient(getParameterTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS Client: %s", err)
	}

	getParameterTemplatePath := getParameterTemplateClient.Endpoint + getParameterTemplateHttpUrl
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{project_id}",
		getParameterTemplateClient.ProjectID)
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{config_id}", d.Id())

	getParameterTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getParameterTemplateResp, err := getParameterTemplateClient.Request("GET",
		getParameterTemplatePath, &getParameterTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200057"),
			"error retrieving DDS parameter template")
	}

	getParameterTemplateRespBody, err := utils.FlattenResponse(getParameterTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getParameterTemplateRespBody, nil)),
		d.Set("node_version", utils.PathSearch("datastore_version",
			getParameterTemplateRespBody, nil)),
		d.Set("description", utils.PathSearch("description",
			getParameterTemplateRespBody, nil)),
		d.Set("parameters", flattenGetParameterTemplateResponseBodyParameter(getParameterTemplateRespBody)),
		d.Set("created_at", utils.PathSearch("created", getParameterTemplateRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated", getParameterTemplateRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetParameterTemplateResponseBodyParameter(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("parameters", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":             utils.PathSearch("name", v, nil),
			"value":            utils.PathSearch("value", v, nil),
			"description":      utils.PathSearch("description", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"value_range":      utils.PathSearch("value_range", v, nil),
			"restart_required": utils.PathSearch("restart_required", v, nil),
			"readonly":         utils.PathSearch("readonly", v, nil),
		})
	}
	return rst
}

func resourceDdsParameterTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteParameterTemplate: Delete DDS parameter template
	var (
		deleteParameterTemplateHttpUrl = "v3/{project_id}/configurations/{config_id}"
		deleteParameterTemplateProduct = "dds"
	)
	deleteParameterTemplateClient, err := cfg.NewServiceClient(deleteParameterTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS Client: %s", err)
	}

	deleteParameterTemplatePath := deleteParameterTemplateClient.Endpoint + deleteParameterTemplateHttpUrl
	deleteParameterTemplatePath = strings.ReplaceAll(deleteParameterTemplatePath, "{project_id}",
		deleteParameterTemplateClient.ProjectID)
	deleteParameterTemplatePath = strings.ReplaceAll(deleteParameterTemplatePath, "{config_id}", d.Id())

	deleteParameterTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	_, err = deleteParameterTemplateClient.Request("DELETE", deleteParameterTemplatePath,
		&deleteParameterTemplateOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200057"),
			"error deleting DDS parameter template")
	}

	return nil
}
