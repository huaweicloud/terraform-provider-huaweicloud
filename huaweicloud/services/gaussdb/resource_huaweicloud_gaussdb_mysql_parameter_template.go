// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GaussDB
// ---------------------------------------------------------------

package gaussdb

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL POST /v3/{project_id}/configurations
// @API GaussDBforMySQL PUT /v3/{project_id}/configurations/{configuration_id}
// @API GaussDBforMySQL GET /v3/{project_id}/configurations/{configuration_id}
// @API GaussDBforMySQL DELETE /v3/{project_id}/configurations/{configuration_id}
func ResourceGaussDBMysqlTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParameterTemplateCreate,
		UpdateContext: resourceParameterTemplateUpdate,
		ReadContext:   resourceParameterTemplateRead,
		DeleteContext: resourceParameterTemplateDelete,
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the parameter template description.`,
			},
			"datastore_engine": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the DB engine.`,
			},
			"datastore_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the DB version.`,
				RequiredWith: []string{
					"datastore_engine",
				},
			},
			"parameter_values": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the mapping between parameter names and parameter values.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time in the "yyyy-MM-ddTHH:mm:ssZ" format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the update time in the "yyyy-MM-ddTHH:mm:ssZ" format.`,
			},
		},
	}
}

func resourceParameterTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGaussDBMysqlParameterTemplate: create a GaussDB MySQL parameter Template
	var (
		createParameterTemplateHttpUrl = "v3/{project_id}/configurations"
		createParameterTemplateProduct = "gaussdb"
	)
	createParameterTemplateClient, err := cfg.NewServiceClient(createParameterTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	createParameterTemplatePath := createParameterTemplateClient.Endpoint + createParameterTemplateHttpUrl
	createParameterTemplatePath = strings.ReplaceAll(createParameterTemplatePath, "{project_id}",
		createParameterTemplateClient.ProjectID)

	createParameterTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createParameterTemplateOpt.JSONBody = utils.RemoveNil(buildCreateParameterTemplateBodyParams(d))
	createParameterTemplateResp, err := createParameterTemplateClient.Request("POST",
		createParameterTemplatePath, &createParameterTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB MySQL Parameter Template: %s", err)
	}

	createParameterTemplateRespBody, err := utils.FlattenResponse(createParameterTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("configurations.id", createParameterTemplateRespBody)
	if err != nil {
		return diag.Errorf("error creating GaussDB MySQL Parameter Template: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceParameterTemplateRead(ctx, d, meta)
}

func buildCreateParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             utils.ValueIgnoreEmpty(d.Get("name")),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"datastore":        buildCreateParameterTemplateDatastoreChildBody(d),
		"parameter_values": utils.ValueIgnoreEmpty(d.Get("parameter_values")),
	}
	return bodyParams
}

func buildCreateParameterTemplateDatastoreChildBody(d *schema.ResourceData) map[string]interface{} {
	datastoreEngine := d.Get("datastore_engine").(string)
	if datastoreEngine == "" {
		return nil
	}
	params := map[string]interface{}{
		"type":    utils.ValueIgnoreEmpty(datastoreEngine),
		"version": utils.ValueIgnoreEmpty(d.Get("datastore_version")),
	}
	return params
}

func resourceParameterTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getParameterTemplate: Query the GaussDB MySQL parameter Template
	var (
		getParameterTemplateHttpUrl = "v3/{project_id}/configurations/{configuration_id}"
		getParameterTemplateProduct = "gaussdb"
	)
	getParameterTemplateClient, err := cfg.NewServiceClient(getParameterTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	getParameterTemplatePath := getParameterTemplateClient.Endpoint + getParameterTemplateHttpUrl
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{project_id}",
		getParameterTemplateClient.ProjectID)
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{configuration_id}", d.Id())

	getParameterTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getParameterTemplateResp, err := getParameterTemplateClient.Request("GET", getParameterTemplatePath,
		&getParameterTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDBMysqlTemplate")
	}

	getParameterTemplateRespBody, err := utils.FlattenResponse(getParameterTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	configurations := utils.PathSearch("configurations", getParameterTemplateRespBody, nil)
	if configurations == nil {
		log.Printf("[WARN] failed to get GaussDB MySQL ParameterTemplate by ID(%s)", d.Id())
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", configurations, nil)),
		d.Set("description", utils.PathSearch("description", configurations, nil)),
		d.Set("datastore_engine", utils.PathSearch("datastore.type", configurations, nil)),
		d.Set("datastore_version", utils.PathSearch("datastore.version", configurations, nil)),
		d.Set("created_at", utils.PathSearch("created", configurations, nil)),
		d.Set("updated_at", utils.PathSearch("updated", configurations, nil)),
	)

	parameterValuesRaw := d.Get("parameter_values").(map[string]interface{})
	parameterValuesRes := utils.PathSearch("parameter_values", getParameterTemplateRespBody,
		make(map[string]interface{})).(map[string]interface{})
	parameterValues := make(map[string]interface{})
	for key := range parameterValuesRaw {
		if v, ok := parameterValuesRes[key]; ok {
			parameterValues[key] = v
		}
	}
	mErr = multierror.Append(mErr, d.Set("parameter_values", parameterValues))

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceParameterTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateGaussDBMysqlTemplateHasChanges := []string{
		"name",
		"description",
		"parameter_values",
	}

	if d.HasChanges(updateGaussDBMysqlTemplateHasChanges...) {
		// updateParameterTemplate: update the GaussDB MySQL parameter Template
		var (
			updateParameterTemplateHttpUrl = "v3/{project_id}/configurations/{configuration_id}"
			updateParameterTemplateProduct = "gaussdb"
		)
		updateParameterTemplateClient, err := cfg.NewServiceClient(updateParameterTemplateProduct, region)
		if err != nil {
			return diag.Errorf("error creating GaussDB Client: %s", err)
		}

		updateParameterTemplatePath := updateParameterTemplateClient.Endpoint + updateParameterTemplateHttpUrl
		updateParameterTemplatePath = strings.ReplaceAll(updateParameterTemplatePath,
			"{project_id}", updateParameterTemplateClient.ProjectID)
		updateParameterTemplatePath = strings.ReplaceAll(updateParameterTemplatePath,
			"{configuration_id}", d.Id())

		updateParameterTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateParameterTemplateOpt.JSONBody = utils.RemoveNil(buildUpdateParameterTemplateBodyParams(d))
		_, err = updateParameterTemplateClient.Request("PUT", updateParameterTemplatePath,
			&updateParameterTemplateOpt)
		if err != nil {
			return diag.Errorf("error updating GaussDB MySQL ParameterTemplate: %s", err)
		}
	}
	return resourceParameterTemplateRead(ctx, d, meta)
}

func buildUpdateParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             utils.ValueIgnoreEmpty(d.Get("name")),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"parameter_values": utils.ValueIgnoreEmpty(d.Get("parameter_values")),
	}
	return bodyParams
}

func resourceParameterTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteParameterTemplate: delete the GaussDB MySQL parameter Template
	var (
		deleteParameterTemplateHttpUrl = "v3/{project_id}/configurations/{configuration_id}"
		deleteParameterTemplateProduct = "gaussdb"
	)
	deleteParameterTemplateClient, err := cfg.NewServiceClient(deleteParameterTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	deleteParameterTemplatePath := deleteParameterTemplateClient.Endpoint + deleteParameterTemplateHttpUrl
	deleteParameterTemplatePath = strings.ReplaceAll(deleteParameterTemplatePath, "{project_id}",
		deleteParameterTemplateClient.ProjectID)
	deleteParameterTemplatePath = strings.ReplaceAll(deleteParameterTemplatePath, "{configuration_id}", d.Id())

	deleteGaussDBMysqlTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteParameterTemplateClient.Request("DELETE", deleteParameterTemplatePath,
		&deleteGaussDBMysqlTemplateOpt)
	if err != nil {
		return diag.Errorf("error deleting GaussDB MySQL ParameterTemplate: %s", err)
	}

	return nil
}
