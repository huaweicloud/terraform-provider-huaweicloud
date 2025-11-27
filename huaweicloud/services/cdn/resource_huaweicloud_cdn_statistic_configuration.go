package cdn

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var statisticConfigurationNonUpdatableParams = []string{
	"config_type",
	"resource_type",
	"resource_name",
	"config_info",
}

// @API CDN POST /v1/cdn/statistics/stats-configs
func ResourceStatisticConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatisticConfigurationCreate,
		ReadContext:   resourceStatisticConfigurationRead,
		UpdateContext: resourceStatisticConfigurationUpdate,
		DeleteContext: resourceStatisticConfigurationDelete,

		CustomizeDiff: config.FlexibleForceNew(statisticConfigurationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource type.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource name, which can be an account or domain name.`,
			},
			"config_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether to enable the top URL statistics configuration.`,
									},
								},
							},
							Description: `The top URL statistics configuration.`,
						},
						"ua": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether to enable the top UA statistics configuration.`,
									},
								},
							},
							Description: `The top UA statistics configuration.`,
						},
					},
				},
				Description: `The statistics configuration information.`,
			},

			// Optional parameters.
			"config_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The configuration category.`,
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

func buildStatisticConfigurationConfigInfoUrl(urlConfigs []interface{}) map[string]interface{} {
	if len(urlConfigs) < 1 {
		return nil
	}

	urlConfig := urlConfigs[0]
	return map[string]interface{}{
		"enable": utils.PathSearch("enable", urlConfig, false),
	}
}

func buildStatisticConfigurationConfigInfoUa(uaConfigs []interface{}) map[string]interface{} {
	if len(uaConfigs) < 1 {
		return nil
	}

	uaConfig := uaConfigs[0]
	return map[string]interface{}{
		"enable": utils.PathSearch("enable", uaConfig, false),
	}
}

func buildStatisticConfigurationConfigInfo(configInfos []interface{}) map[string]interface{} {
	if len(configInfos) < 1 {
		return nil
	}

	configInfo := configInfos[0]
	return utils.RemoveNil(map[string]interface{}{
		"url": buildStatisticConfigurationConfigInfoUrl(utils.PathSearch("url", configInfo, make([]interface{}, 0)).([]interface{})),
		"ua":  buildStatisticConfigurationConfigInfoUa(utils.PathSearch("ua", configInfo, make([]interface{}, 0)).([]interface{})),
	})
}

func buildStatisticConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"resource_type": d.Get("resource_type"),
		"resource_name": d.Get("resource_name"),
		"config_info":   buildStatisticConfigurationConfigInfo(d.Get("config_info").([]interface{})),
		"config_type":   utils.ValueIgnoreEmpty(d.Get("config_type")),
	})
}

func resourceStatisticConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	httpUrl := "v1/cdn/statistics/stats-configs"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildStatisticConfigurationBodyParams(d),
		OkCodes:  []int{204},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CDN statistic configuration: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceStatisticConfigurationRead(ctx, d, meta)
}

func resourceStatisticConfigurationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceStatisticConfigurationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceStatisticConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to set CDN statistic configuration. Deleting this resource
	will not clear the corresponding request record, but will only remove the resource information from the tf state
    file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
