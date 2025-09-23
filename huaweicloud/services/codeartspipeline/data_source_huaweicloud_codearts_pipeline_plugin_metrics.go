package codeartspipeline

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsPipeline POST /v1/{domain_id}/agent-plugin/plugin-metrics
func DataSourceCodeArtsPipelinePluginMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelinePluginMetricsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"plugin_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the plugin name.`,
			},
			"plugin_attribution": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the extension attribute, official or custom.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the display name.`,
			},
			"version_attribution": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version attribute. The value can be draft or formal.`,
			},
			"metrics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the plugin list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the version.`,
						},
						"output_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the output key.`,
						},
						"output_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the output value.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelinePluginMetricsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v1/{domain_id}/agent-plugin/plugin-metrics"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"X-Language": "en-us"},
		JSONBody:         buildPipelineCodeArtsPipelinePluginMetricsQueryParams(d),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving pipeline plugin metrics: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	metrics := utils.PathSearch("[0].data", getRespBody, make([]interface{}, 0)).([]interface{})
	rst := make([]map[string]interface{}, 0, len(metrics))
	for _, metric := range metrics {
		rst = append(rst, map[string]interface{}{
			"version":      utils.PathSearch("version", metric, nil),
			"output_key":   utils.PathSearch("output_key", metric, nil),
			"output_value": utils.PathSearch("output_value", metric, nil),
		})
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("metrics", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPipelineCodeArtsPipelinePluginMetricsQueryParams(d *schema.ResourceData) interface{} {
	bodyParams := map[string]interface{}{
		"plugin_name":         d.Get("plugin_name"),
		"plugin_attribution":  utils.ValueIgnoreEmpty(d.Get("plugin_attribution")),
		"version":             utils.ValueIgnoreEmpty(d.Get("version")),
		"display_name":        utils.ValueIgnoreEmpty(d.Get("display_name")),
		"version_attribution": utils.ValueIgnoreEmpty(d.Get("version_attribution")),
	}

	return []interface{}{bodyParams}
}
