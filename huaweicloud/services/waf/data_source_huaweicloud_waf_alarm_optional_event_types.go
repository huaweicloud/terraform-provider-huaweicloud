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
func DataSourceAlarmOptionalEventTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmOptionalEventTypesRead,

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

func dataSourceAlarmOptionalEventTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/tag/threat/map"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving WAF alarm optional event types: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("threats", utils.ExpandToStringList(utils.PathSearch("threats", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("locale", flattenAlarmOptionalEventTypesLocale(utils.PathSearch("locale", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAlarmOptionalEventTypesLocale(rawLocale interface{}) []interface{} {
	if rawLocale == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cmdi":                    utils.PathSearch("cmdi", rawLocale, nil),
			"llm_prompt_injection":    utils.PathSearch("llm_prompt_injection", rawLocale, nil),
			"anticrawler":             utils.PathSearch("anticrawler", rawLocale, nil),
			"custom_custom":           utils.PathSearch("custom_custom", rawLocale, nil),
			"third_bot_river":         utils.PathSearch("third_bot_river", rawLocale, nil),
			"robot":                   utils.PathSearch("robot", rawLocale, nil),
			"custom_idc_ip":           utils.PathSearch("custom_idc_ip", rawLocale, nil),
			"antiscan_dir_traversal":  utils.PathSearch("antiscan_dir_traversal", rawLocale, nil),
			"advanced_bot":            utils.PathSearch("advanced_bot", rawLocale, nil),
			"xss":                     utils.PathSearch("xss", rawLocale, nil),
			"antiscan_high_freq_scan": utils.PathSearch("antiscan_high_freq_scan", rawLocale, nil),
			"webshell":                utils.PathSearch("webshell", rawLocale, nil),
			"cc":                      utils.PathSearch("cc", rawLocale, nil),
			"botm":                    utils.PathSearch("botm", rawLocale, nil),
			"illegal":                 utils.PathSearch("illegal", rawLocale, nil),
			"llm_prompt_sensitive":    utils.PathSearch("llm_prompt_sensitive", rawLocale, nil),
			"sqli":                    utils.PathSearch("sqli", rawLocale, nil),
			"lfi":                     utils.PathSearch("lfi", rawLocale, nil),
			"antitamper":              utils.PathSearch("antitamper", rawLocale, nil),
			"custom_geoip":            utils.PathSearch("custom_geoip", rawLocale, nil),
			"rfi":                     utils.PathSearch("rfi", rawLocale, nil),
			"vuln":                    utils.PathSearch("vuln", rawLocale, nil),
			"llm_response_sensitive":  utils.PathSearch("llm_response_sensitive", rawLocale, nil),
			"custom_whiteblackip":     utils.PathSearch("custom_whiteblackip", rawLocale, nil),
			"leakage":                 utils.PathSearch("leakage", rawLocale, nil),
		},
	}
}
