package fgs

import (
	"context"
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

var functionTracingConfigurationNonUpdatableParams = []string{"function_urn"}

// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/tracing
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/tracing
func ResourceFunctionTracingConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionTracingConfigurationCreate,
		ReadContext:   resourceFunctionTracingConfigurationRead,
		UpdateContext: resourceFunctionTracingConfigurationUpdate,
		DeleteContext: resourceFunctionTracingConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceFunctionTracingConfigurationImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(functionTracingConfigurationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the function tracing configuration is located.`,
			},

			// Required parameters.
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The URN of the function to which the tracing configuration belongs.`,
			},
			"tracing_ak": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The APM access key for tracing configuration.`,
			},
			"tracing_sk": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The APM secret key for tracing configuration.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildFunctionTracingConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"tracing_ak": d.Get("tracing_ak"),
		"tracing_sk": d.Get("tracing_sk"),
	}
}

func resourceFunctionTracingConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/tracing"
	functionUrn := d.Get("function_urn").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{function_urn}", functionUrn)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildFunctionTracingConfigurationBodyParams(d),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph function tracing configuration: %s", err)
	}

	d.SetId(functionUrn)

	return resourceFunctionTracingConfigurationRead(ctx, d, meta)
}

func GetFunctionTracingConfiguration(client *golangsdk.ServiceClient, functionUrn string) (interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/tracing"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{function_urn}", functionUrn)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceFunctionTracingConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	respBody, err := GetFunctionTracingConfiguration(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving FunctionGraph function tracing configuration")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tracing_ak", utils.PathSearch("tracing_ak", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceFunctionTracingConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/tracing"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildFunctionTracingConfigurationBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating FunctionGraph function tracing configuration: %s", err)
	}

	return resourceFunctionTracingConfigurationRead(ctx, d, meta)
}

func resourceFunctionTracingConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/tracing"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{function_urn}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         make(map[string]interface{}),
	}

	_, err = client.Request("PUT", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting FunctionGraph function tracing configuration: %s", err)
	}

	return nil
}

func resourceFunctionTracingConfigurationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("function_urn", d.Id())
}
