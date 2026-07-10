package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/scan-templates/{template_id}/scan-rules
func DataSourceDscTemplateRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscTemplateRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"classification_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_level_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_used": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_rules_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscTemplateRulesSchema(),
			},
		},
	}
}

func dscTemplateRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"classification_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_level_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_level_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_level_color": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_used": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rule_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDscTemplateRulesQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := ""

	if v, ok := d.GetOk("keyword"); ok {
		queryParams = fmt.Sprintf("%s&keyword=%v", queryParams, v)
	}
	if v, ok := d.GetOk("classification_ids"); ok {
		classificationIds := v.([]interface{})
		for _, id := range classificationIds {
			queryParams = fmt.Sprintf("%s&classification_ids=%v", queryParams, id)
		}
	}
	if v, ok := d.GetOk("security_level_ids"); ok {
		securityLevelIds := v.([]interface{})
		for _, id := range securityLevelIds {
			queryParams = fmt.Sprintf("%s&security_level_ids=%v", queryParams, id)
		}
	}
	if v, ok := d.GetOk("is_used"); ok {
		queryParams = fmt.Sprintf("%s&is_used=%v", queryParams, v)
	}
	if v, ok := d.GetOk("rule_name"); ok {
		queryParams = fmt.Sprintf("%s&rule_name=%v", queryParams, v)
	}

	// In the API, this field is designated as the page number.
	if offset > 0 {
		queryParams = fmt.Sprintf("%s&offset=%v", queryParams, offset)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceDscTemplateRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/scan-templates/{template_id}/scan-rules"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", d.Get("template_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := requestPath + buildDscTemplateRulesQueryParams(d, offset)
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC template rules: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		templateRulesList := utils.PathSearch("template_rules_list", requestRespBody, make([]interface{}, 0)).([]interface{})
		if len(templateRulesList) == 0 {
			break
		}

		result = append(result, templateRulesList...)
		offset++
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("template_rules_list", flattenDscTemplateRulesRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscTemplateRulesRecords(templateRulesList []interface{}) []interface{} {
	if len(templateRulesList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(templateRulesList))
	for _, v := range templateRulesList {
		rst = append(rst, map[string]interface{}{
			"rule_id":              utils.PathSearch("rule_id", v, nil),
			"project_id":           utils.PathSearch("project_id", v, nil),
			"rule_name":            utils.PathSearch("rule_name", v, nil),
			"template_id":          utils.PathSearch("template_id", v, nil),
			"classification_id":    utils.PathSearch("classification_id", v, nil),
			"security_level_id":    utils.PathSearch("security_level_id", v, nil),
			"security_level_name":  utils.PathSearch("security_level_name", v, nil),
			"security_level_color": utils.PathSearch("security_level_color", v, nil),
			"is_used":              utils.PathSearch("is_used", v, nil),
			"rule_description":     utils.PathSearch("rule_description", v, nil),
		})
	}

	return rst
}
