package gaussdb

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

// @API GaussDB POST /v3/{project_id}/configurations
// @API GaussDB POST /v3/{project_id}/configurations/{config_id}/copy
// @API GaussDB GET /v3/{project_id}/configurations/{config_id}
// @API GaussDB DELETE /v3/{project_id}/configurations/{config_id}
func ResourceOpenGaussParameterTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussParameterTemplateCreate,
		ReadContext:   resourceOpenGaussParameterTemplateRead,
		DeleteContext: resourceOpenGaussParameterTemplateDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"engine_version": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				RequiredWith:  []string{"instance_mode"},
				ConflictsWith: []string{"source_configuration_id"},
			},
			"instance_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				RequiredWith:  []string{"engine_version"},
				ConflictsWith: []string{"source_configuration_id"},
			},
			"source_configuration_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"engine_version", "instance_mode"},
			},
			"parameters": {
				Type:          schema.TypeSet,
				Elem:          templateParametersSchema(),
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				RequiredWith:  []string{"engine_version", "instance_mode"},
				ConflictsWith: []string{"source_configuration_id"},
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

func templateParametersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"need_restart": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceOpenGaussParameterTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	var id string
	if _, ok := d.GetOk("source_configuration_id"); ok {
		id, err = copyParameterTemplate(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		id, err = createParameterTemplate(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(id)

	return resourceOpenGaussParameterTemplateRead(ctx, d, meta)
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
		return "", fmt.Errorf("error creating GaussDB OpenGauss parameter template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return "", err
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return "", fmt.Errorf("error creating GaussDB OpenGauss parameter template: ID is not found in API response")
	}
	return id, nil
}

func buildCreateParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             d.Get("name"),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"datastore":        buildCreateParameterTemplateDatastoreChildBody(d),
		"parameter_values": buildCreateTemplateParametersBodyParam(d),
	}
	return bodyParams
}

func buildCreateParameterTemplateDatastoreChildBody(d *schema.ResourceData) map[string]interface{} {
	datastoreEngine := d.Get("engine_version").(string)
	if datastoreEngine == "" {
		return nil
	}
	params := map[string]interface{}{
		"engine_version": utils.ValueIgnoreEmpty(datastoreEngine),
		"instance_mode":  utils.ValueIgnoreEmpty(d.Get("instance_mode")),
	}
	return params
}

func copyParameterTemplate(d *schema.ResourceData, client *golangsdk.ServiceClient) (string, error) {
	httpUrl := "v3/{project_id}/configurations/{config_id}/copy"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{config_id}", d.Get("source_configuration_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCopyParameterTemplateBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return "", fmt.Errorf("error creating GaussDB OpenGauss parameter template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return "", err
	}

	id := utils.PathSearch("config_id", createRespBody, "").(string)
	if id == "" {
		return "", fmt.Errorf("error creating GaussDB OpenGauss parameter template: config_id is not found in API response")
	}
	return id, nil
}

func buildCreateTemplateParametersBodyParam(d *schema.ResourceData) map[string]string {
	rawParameters := d.Get("parameters").(*schema.Set)
	if rawParameters.Len() == 0 {
		return nil
	}
	rst := make(map[string]string)
	for _, v := range rawParameters.List() {
		if raw, ok := v.(map[string]interface{}); ok {
			rst[raw["name"].(string)] = raw["value"].(string)
		}
	}
	return rst
}

func buildCopyParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceOpenGaussParameterTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB OpenGauss parameter template")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", getRespBody, nil)),
		d.Set("instance_mode", utils.PathSearch("instance_mode", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
		d.Set("parameters", flattenGaussDBOpenGaussResponseBodyParameters(d, getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGaussDBOpenGaussResponseBodyParameters(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	paramsMap := buildParamsMap(d)
	curJson := utils.PathSearch("configuration_parameters", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		paramName := utils.PathSearch("name", v, "").(string)
		if !paramsMap[paramName] {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"name":         paramName,
			"value":        utils.PathSearch("value", v, nil),
			"need_restart": utils.PathSearch("need_restart", v, nil),
			"readonly":     utils.PathSearch("readonly", v, nil),
			"value_range":  utils.PathSearch("value_range", v, nil),
			"data_type":    utils.PathSearch("data_type", v, nil),
			"description":  utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func buildParamsMap(d *schema.ResourceData) map[string]bool {
	params := d.Get("parameters").(*schema.Set).List()
	paramsMap := make(map[string]bool)
	for _, param := range params {
		if v, ok := param.(map[string]interface{}); ok {
			paramsMap[v["name"].(string)] = true
		}
	}
	return paramsMap
}

func resourceOpenGaussParameterTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{config_id}", d.Id())

	deleteGOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath,
		&deleteGOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting GaussDB OpenGauss parameter template")
	}

	return nil
}
