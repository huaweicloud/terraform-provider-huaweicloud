package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/baseline/security-checks/policy-detail
func DataSourceBaselineSecurityChecksDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineSecurityChecksDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_os": {
				Type:     schema.TypeString,
				Required: true,
			},
			"standard": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The API is `bool` string, set to `string` type here to support both true and false scenarios.
			"checked": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_rule_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"checked": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rule_params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_param_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rule_desc": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"default_value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"range_min": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"range_max": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildCheckedParam(queryParams, checked string) string {
	if checked == "false" {
		return fmt.Sprintf("%s&checked=%v", queryParams, false)
	}
	if checked == "true" {
		return fmt.Sprintf("%s&checked=%v", queryParams, true)
	}

	return queryParams
}

func buildBaselineSecurityChecksDetailsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	queryParams = fmt.Sprintf("%s&support_os=%v", queryParams, d.Get("support_os"))
	queryParams = fmt.Sprintf("%s&standard=%v", queryParams, d.Get("standard"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("check_type"); ok {
		queryParams = fmt.Sprintf("%s&check_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("tag"); ok {
		queryParams = fmt.Sprintf("%s&tag=%v", queryParams, v)
	}
	if v, ok := d.GetOk("check_rule_name"); ok {
		queryParams = fmt.Sprintf("%s&check_rule_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("level"); ok {
		queryParams = fmt.Sprintf("%s&level=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("checked"); ok {
		queryParams = buildCheckedParam(queryParams, v.(string))
	}

	return queryParams
}

func dataSourceBaselineSecurityChecksDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v5/{project_id}/baseline/security-checks/default-policy/details"
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBaselineSecurityChecksDetailsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS baseline security checks default policy details: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		checkDetailsResp := utils.PathSearch("check_details", respBody, make([]interface{}, 0)).([]interface{})
		if len(checkDetailsResp) == 0 {
			break
		}

		result = append(result, checkDetailsResp...)
		offset += len(checkDetailsResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("check_details", flattenCheckDetails(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCheckDetails(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"key":             utils.PathSearch("key", v, nil),
			"check_rule_id":   utils.PathSearch("check_rule_id", v, nil),
			"check_rule_name": utils.PathSearch("check_rule_name", v, nil),
			"check_rule_type": utils.PathSearch("check_rule_type", v, nil),
			"check_type":      utils.PathSearch("check_type", v, nil),
			"severity":        utils.PathSearch("severity", v, nil),
			"level":           utils.PathSearch("level", v, nil),
			"checked":         utils.PathSearch("checked", v, nil),
			"rule_params": flattenCheckDetailsRuleParams(
				utils.PathSearch("rule_params", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenCheckDetailsRuleParams(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"rule_param_id": utils.PathSearch("rule_param_id", v, nil),
			"rule_desc":     utils.PathSearch("rule_desc", v, nil),
			"default_value": utils.PathSearch("default_value", v, nil),
			"range_min":     utils.PathSearch("range_min", v, nil),
			"range_max":     utils.PathSearch("range_max", v, nil),
		})
	}

	return result
}
