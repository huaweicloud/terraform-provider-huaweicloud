package fgs

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

// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/async-invoke-config
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/async-invoke-config
// @API FunctionGraph DELETE /v2/{project_id}/fgs/functions/{function_urn}/async-invoke-config
func ResourceAsyncInvokeConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAsyncInvokeConfigurationCreate,
		ReadContext:   resourceAsyncInvokeConfigurationRead,
		UpdateContext: resourceAsyncInvokeConfigurationUpdate,
		DeleteContext: resourceAsyncInvokeConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAsyncInvokeConfigImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to configure the asynchronous invocation.`,
			},
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The function URN to which the asynchronous invocation belongs.`,
			},
			"max_async_event_age_in_seconds": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The maximum validity period of a message.`,
			},
			"max_async_retry_attempts": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The maximum number of retry attempts to be made if asynchronous invocation fails.`,
			},
			"on_success": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        destinationConfigSchemaResource(),
				Description: `The target to be invoked when a function is successfully executed.`,
			},
			"on_failure": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     destinationConfigSchemaResource(),
				Description: `The target to be invoked when a function fails to be executed due to a system error or an
internal error.`,
			},
			"enable_async_status_log": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable asynchronous invocation status persistence.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the asynchronous invocation.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the asynchronous invocation.`,
			},
		},
	}
}

func destinationConfigSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The object type.`,
			},
			"param": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The parameters (in JSON format) corresponding to the target service.`,
			},
		},
	}
}

func buildAsyncInvokeConfigurationDestinationConfig(destinationConfigs []interface{}) map[string]interface{} {
	if len(destinationConfigs) < 1 {
		return nil
	}

	destinationConfig := destinationConfigs[0]
	return map[string]interface{}{
		// Required parameters.
		"destination": utils.PathSearch("destination", destinationConfig, ""),
		"param":       utils.PathSearch("param", destinationConfig, ""),
	}
}

func buildCreateAsyncInvokeConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"max_async_event_age_in_seconds": d.Get("max_async_event_age_in_seconds").(int),
		"max_async_retry_attempts":       d.Get("max_async_retry_attempts").(int),
		// Optional parameters.
		"enable_async_status_log": d.Get("enable_async_status_log").(bool),
		// Empty structure will send if the map is not parsing destination configurations via RemoveNil function.
		"destination_config": utils.RemoveNil(map[string]interface{}{
			"on_success": buildAsyncInvokeConfigurationDestinationConfig(d.Get("on_success").([]interface{})),
			"on_failure": buildAsyncInvokeConfigurationDestinationConfig(d.Get("on_failure").([]interface{})),
		}),
	}
}

func modifyAsyncInvokeConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		functionUrn = d.Get("function_urn").(string)
		httpUrl     = "v2/{project_id}/fgs/functions/{function_urn}/async-invoke-config"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", functionUrn)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateAsyncInvokeConfigurationBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceAsyncInvokeConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	err = modifyAsyncInvokeConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error creating the configuration of the asynchronous invocation: %s", err)
	}
	d.SetId(d.Get("function_urn").(string))

	return resourceAsyncInvokeConfigurationRead(ctx, d, meta)
}

func GetAsyncIncokeConfigurations(client *golangsdk.ServiceClient, functionUrn string) (interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/async-invoke-config"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{function_urn}", functionUrn)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenDestinationItem(destConfig map[string]interface{}) []interface{} {
	if len(destConfig) < 1 {
		return nil
	}

	parsedConfig := utils.RemoveNil(map[string]interface{}{
		"destination": utils.ValueIgnoreEmpty(utils.PathSearch("destination", destConfig, nil)),
		"param":       utils.ValueIgnoreEmpty(utils.PathSearch("param", destConfig, nil)),
	})

	result := make([]interface{}, 0, 1)
	if len(parsedConfig) > 0 {
		result = append(result, parsedConfig)
	}
	return result
}

func resourceAsyncInvokeConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		functionUrn = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	resp, err := GetAsyncIncokeConfigurations(client, functionUrn)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "asynchronous invocation configuration")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("max_async_event_age_in_seconds", utils.PathSearch("max_async_event_age_in_seconds", resp, nil)),
		d.Set("max_async_retry_attempts", utils.PathSearch("max_async_retry_attempts", resp, nil)),
		d.Set("on_success", flattenDestinationItem(utils.PathSearch("destination_config.on_success",
			resp, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("on_failure", flattenDestinationItem(utils.PathSearch("destination_config.on_failure",
			resp, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("enable_async_status_log", utils.PathSearch("enable_async_status_log", resp, nil)),
		d.Set("created_at", utils.PathSearch("created_time", resp, nil)),
		d.Set("updated_at", utils.PathSearch("last_modified", resp, nil)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving asynchronous invocation configuration fields: %s", mErr)
	}
	return nil
}

func resourceAsyncInvokeConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	err = modifyAsyncInvokeConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error updating the configuration of the asynchronous invocation: %s", err)
	}

	return resourceAsyncInvokeConfigurationRead(ctx, d, meta)
}

func resourceAsyncInvokeConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/fgs/functions/{function_urn}/async-invoke-config"
		functionUrn = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{function_urn}", functionUrn)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting the configuration of the asynchronous invocation")
	}
	return nil
}

func resourceAsyncInvokeConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("function_urn", d.Id())
}
