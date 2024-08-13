package fgs

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"

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
				Description: "The region in which to configure the asynchronous invocation.",
			},
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The function URN to which the asynchronous invocation belongs.",
			},
			"max_async_event_age_in_seconds": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The maximum validity period of a message.",
			},
			"max_async_retry_attempts": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "The maximum number of retry attempts to be made if asynchronous invocation " +
					"fails.",
			},
			"on_success": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        destinationConfigSchemaResource(),
				Description: "The target to be invoked when a function is successfully executed.",
			},
			"on_failure": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     destinationConfigSchemaResource(),
				Description: "The target to be invoked when a function fails to be executed due to a " +
					"system error or an internal error.",
			},
			"enable_async_status_log": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable asynchronous invocation status persistence.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the asynchronous invocation.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the asynchronous invocation.",
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
				Description: "The object type.",
			},
			"param": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parameters (in JSON format) corresponding to the target service.",
			},
		},
	}
}

func modifyAsyncInvokeConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		functionUrn = d.Get("function_urn").(string)
		opts        = function.AsyncInvokeConfigOpts{
			MaxAsyncEventAgeInSeconds: d.Get("max_async_event_age_in_seconds").(int),
			MaxAsyncRetryAttempts:     d.Get("max_async_retry_attempts").(int),
			EnableAsyncStatusLog:      utils.Bool(d.Get("enable_async_status_log").(bool)),
		}
		destinationConfig = function.DestinationConfig{}
	)

	if successConfigs, ok := d.GetOk("on_success"); ok {
		raws := successConfigs.([]interface{})
		cfgDetails := raws[0].(map[string]interface{})
		destinationConfig.OnSuccess = function.DestinationConfigDetails{
			Destination: cfgDetails["destination"].(string),
			Param:       cfgDetails["param"].(string),
		}
	}
	if failureConfigs, ok := d.GetOk("on_failure"); ok {
		raws := failureConfigs.([]interface{})
		cfgDetails := raws[0].(map[string]interface{})
		destinationConfig.OnFailure = function.DestinationConfigDetails{
			Destination: cfgDetails["destination"].(string),
			Param:       cfgDetails["param"].(string),
		}
	}
	if destinationConfig != (function.DestinationConfig{}) {
		opts.DestinationConfig = destinationConfig
	}
	_, err := function.UpdateAsyncInvokeConfig(client, functionUrn, opts)
	if err != nil {
		return fmt.Errorf("error modifying the async invoke configuration: %s", err)
	}
	return nil
}

func resourceAsyncInvokeConfigurationCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph V2 client: %s", err)
	}

	err = modifyAsyncInvokeConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error creating the configuration of the asynchronous invocation: %s", err)
	}
	d.SetId(d.Get("function_urn").(string))

	return resourceAsyncInvokeConfigurationRead(ctx, d, meta)
}

func flattenDestinationConfig(destConfig function.DestinationConfigDetails) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"destination": destConfig.Destination,
			"param":       destConfig.Param,
		},
	}
}

func resourceAsyncInvokeConfigurationRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.FgsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph V2 client: %s", err)
	}

	resp, err := function.GetAsyncInvokeConfig(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "asynchronous invocation configuration")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("max_async_event_age_in_seconds", resp.MaxAsyncEventAgeInSeconds),
		d.Set("max_async_retry_attempts", resp.MaxAsyncRetryAttempts),
		d.Set("on_success", flattenDestinationConfig(resp.DestinationConfig.OnSuccess)),
		d.Set("on_failure", flattenDestinationConfig(resp.DestinationConfig.OnFailure)),
		d.Set("enable_async_status_log", resp.EnableAsyncStatusLog),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving asynchronous invocation configuration fields: %s", mErr)
	}

	return nil
}

func resourceAsyncInvokeConfigurationUpdate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph V2 client: %s", err)
	}

	err = modifyAsyncInvokeConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error updating the configuration of the asynchronous invocation: %s", err)
	}

	return resourceAsyncInvokeConfigurationRead(ctx, d, meta)
}

func resourceAsyncInvokeConfigurationDelete(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph V2 client: %s", err)
	}

	err = function.DeleteAsyncInvokeConfig(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting the configuration of the asynchronous invocation")
	}
	return nil
}

func resourceAsyncInvokeConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	return []*schema.ResourceData{d}, d.Set("function_urn", d.Id())
}
