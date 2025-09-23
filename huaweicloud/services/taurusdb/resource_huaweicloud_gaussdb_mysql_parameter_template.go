// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GaussDB
// ---------------------------------------------------------------

package taurusdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL POST /v3/{project_id}/configurations
// @API GaussDBforMySQL POST /v3/{project_id}/configurations/{configuration_id}/copy
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/configurations/{configuration_id}/copy
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
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				RequiredWith:  []string{"datastore_version"},
				ConflictsWith: []string{"source_configuration_id", "instance_id"},
				Description:   `Specifies the DB engine.`,
			},
			"datastore_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"datastore_engine"},
				Description:  `Specifies the DB version.`,
			},
			"source_configuration_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"datastore_engine", "instance_id"},
				Description:   `Specifies the source parameter template ID.`,
			},
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				RequiredWith:  []string{"instance_configuration_id"},
				ConflictsWith: []string{"datastore_engine", "source_configuration_id"},
				Description:   `Specifies the ID of the GaussDB MySQL instance.`,
			},
			"instance_configuration_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"instance_id"},
				Description:  `Specifies the parameter template ID of the GaussDB MySQL instance.`,
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
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	var id string
	isCopy := false
	if _, ok := d.GetOk("source_configuration_id"); ok {
		id, err = copyParameterTemplate(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
		isCopy = true
	} else if _, ok = d.GetOk("instance_id"); ok {
		id, err = copyInstanceParameterTemplate(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
		isCopy = true
	} else {
		id, err = createParameterTemplate(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(id)

	if _, ok := d.GetOk("parameter_values"); ok && isCopy {
		err = updateParameterTemplate(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceParameterTemplateRead(ctx, d, meta)
}

func createParameterTemplate(d *schema.ResourceData, client *golangsdk.ServiceClient) (string, error) {
	httpUrl := "v3/{project_id}/configurations"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateParameterTemplateBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return "", fmt.Errorf("error creating GaussDB MySQL parameter template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return "", err
	}

	id := utils.PathSearch("configurations.id", createRespBody, "").(string)
	if id == "" {
		return "", fmt.Errorf("error creating GaussDB MySQL parameter template: ID is not found in API response")
	}
	return id, nil
}

func copyParameterTemplate(d *schema.ResourceData, client *golangsdk.ServiceClient) (string, error) {
	httpUrl := "v3/{project_id}/configurations/{configuration_id}/copy"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{configuration_id}", d.Get("source_configuration_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCopyParameterTemplateBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return "", fmt.Errorf("error creating GaussDB MySQL parameter template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return "", err
	}

	id := utils.PathSearch("configuration_id", createRespBody, "").(string)
	if id == "" {
		return "", fmt.Errorf("error creating GaussDB MySQL parameter template: ID is not found in API response")
	}
	return id, nil
}

func copyInstanceParameterTemplate(d *schema.ResourceData, client *golangsdk.ServiceClient) (string, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/configurations/{configuration_id}/copy"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{configuration_id}", d.Get("instance_configuration_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCopyParameterTemplateBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return "", fmt.Errorf("error creating GaussDB MySQL parameter template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return "", err
	}

	id := utils.PathSearch("configuration_id", createRespBody, "").(string)
	if id == "" {
		return "", fmt.Errorf("error creating GaussDB MySQL parameter template: ID is not found in API response")
	}
	return id, nil
}

func buildCreateParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             d.Get("name"),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"datastore":        buildCreateParameterTemplateDatastoreChildBody(d),
		"parameter_values": utils.ValueIgnoreEmpty(d.Get("parameter_values")),
	}
	return bodyParams
}

func buildCopyParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
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
			product = "gaussdb"
		)
		client, err := cfg.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating GaussDB client: %s", err)
		}

		err = updateParameterTemplate(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceParameterTemplateRead(ctx, d, meta)
}

func updateParameterTemplate(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/configurations/{configuration_id}"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{configuration_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateParameterTemplateBodyParams(d))
	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL parameter template: %s", err)
	}
	return nil
}

func buildUpdateParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             d.Get("name"),
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
