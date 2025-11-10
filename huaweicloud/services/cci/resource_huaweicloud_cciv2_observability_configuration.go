package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI PUT /v1/observabilityconfiguration
// @API CCI GET /v1/observabilityconfiguration
func ResourceV2ObservabilityConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2ObservabilityConfigurationCreateOrUpdate,
		UpdateContext: resourceV2ObservabilityConfigurationCreateOrUpdate,
		ReadContext:   resourceV2ObservabilityConfigurationRead,
		DeleteContext: resourceV2ObservabilityConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2SecretImportState,
		},

		Schema: map[string]*schema.Schema{
			"event": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV2ObservabilityConfigurationCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/observabilityconfiguration"
		product = "cci"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateObservabilityConfigurationParams(d))

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 observability configuration: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(cfg.DomainID)
	}

	return resourceV2ObservabilityConfigurationRead(ctx, d, meta)
}

func buildUpdateObservabilityConfigurationParams(d *schema.ResourceData) map[string]interface{} {
	eventRaw := d.Get("event").([]interface{})
	if len(eventRaw) == 0 {
		return nil
	}

	if event, ok := eventRaw[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"event": map[string]interface{}{
				"enable": event["enable"],
			},
		}
		return bodyParams
	}
	return nil
}

func resourceV2ObservabilityConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/observabilityconfiguration"
		product = "cci"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCI V2 observability configuration")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("event", flattenV2ObservabilityConfigurationEvent(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV2ObservabilityConfigurationEvent(resp interface{}) interface{} {
	event := utils.PathSearch("event", resp, nil)
	if event == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"enable": utils.PathSearch("enable", event, nil),
		},
	}
	return rst
}

func resourceV2ObservabilityConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting CCI v2 observability configuration is not supported. The restoration record is only removed" +
		" from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
