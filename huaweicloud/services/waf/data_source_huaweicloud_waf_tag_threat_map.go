package waf

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

// @API WAF GET /v1/{project_id}/waf/tag/threat/map
func DataSourceTagThreatMap() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagThreatMapRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"threats": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"locale": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmdi": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"llm_prompt_injection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"anticrawler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_custom": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"third_bot_river": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"robot": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_idc_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"antiscan_dir_traversal": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"advanced_bot": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"xss": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"antiscan_high_freq_scan": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"webshell": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"botm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"illegal": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"llm_prompt_sensitive": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sqli": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lfi": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"antitamper": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_geoip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rfi": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vuln": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"llm_response_sensitive": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_whiteblackip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"leakage": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTagThreatMapRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/tag/threat/map"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving event types: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("threats", utils.PathSearch("threats", getRespBody, nil)),
		d.Set("locale", flattenTagThreatMap(utils.PathSearch("locale", getRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTagThreatMap(eventTypes interface{}) []interface{} {
	if eventTypes == nil {
		return nil
	}

	result := map[string]interface{}{
		"cmdi":                    utils.PathSearch("cmdi", eventTypes, nil),
		"llm_prompt_injection":    utils.PathSearch("llm_prompt_injection", eventTypes, nil),
		"anticrawler":             utils.PathSearch("anticrawler", eventTypes, nil),
		"custom_custom":           utils.PathSearch("custom_custom", eventTypes, nil),
		"third_bot_river":         utils.PathSearch("third_bot_river", eventTypes, nil),
		"robot":                   utils.PathSearch("robot", eventTypes, nil),
		"custom_idc_ip":           utils.PathSearch("custom_idc_ip", eventTypes, nil),
		"antiscan_dir_traversal":  utils.PathSearch("antiscan_dir_traversal", eventTypes, nil),
		"advanced_bot":            utils.PathSearch("advanced_bot", eventTypes, nil),
		"xss":                     utils.PathSearch("xss", eventTypes, nil),
		"antiscan_high_freq_scan": utils.PathSearch("antiscan_high_freq_scan", eventTypes, nil),
		"webshell":                utils.PathSearch("webshell", eventTypes, nil),
		"cc":                      utils.PathSearch("cc", eventTypes, nil),
		"botm":                    utils.PathSearch("botm", eventTypes, nil),
		"illegal":                 utils.PathSearch("illegal", eventTypes, nil),
		"llm_prompt_sensitive":    utils.PathSearch("llm_prompt_sensitive", eventTypes, nil),
		"sqli":                    utils.PathSearch("sqli", eventTypes, nil),
		"lfi":                     utils.PathSearch("lfi", eventTypes, nil),
		"antitamper":              utils.PathSearch("antitamper", eventTypes, nil),
		"custom_geoip":            utils.PathSearch("custom_geoip", eventTypes, nil),
		"rfi":                     utils.PathSearch("rfi", eventTypes, nil),
		"vuln":                    utils.PathSearch("vuln", eventTypes, nil),
		"llm_response_sensitive":  utils.PathSearch("llm_response_sensitive", eventTypes, nil),
		"custom_whiteblackip":     utils.PathSearch("custom_whiteblackip", eventTypes, nil),
		"leakage":                 utils.PathSearch("leakage", eventTypes, nil),
	}

	return []interface{}{result}
}
