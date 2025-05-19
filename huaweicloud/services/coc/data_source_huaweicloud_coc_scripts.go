package coc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC GET /v1/job/scripts
func DataSourceCocScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocScriptsRead,

		Schema: map[string]*schema.Schema{
			"name_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"risk_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_created": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gmt_modified": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"script_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"usage_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"properties": dataDataProperties(),
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_tags": dataDataResourceTags(),
					},
				},
			},
		},
	}
}

func dataDataProperties() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"risk_level": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"reviewers": dataDataPropertiesReviewers(),
				"version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"protocol": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataDataPropertiesReviewers() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"reviewer_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"reviewer_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataDataResourceTags() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceCocScriptsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/job/scripts"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	basePath := client.Endpoint + httpUrl

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	marker := 0.0
	res := make([]map[string]interface{}, 0)
	for {
		getPath := basePath + buildGetScriptsParams(d, marker)
		getResp, err := client.Request("GET", getPath, &getOpt)

		if err != nil {
			return diag.Errorf("error retrieving COC scripts: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		scripts, nextMarker := flattenCocGetScripts(getRespBody)
		if len(scripts) < 1 {
			break
		}
		res = append(res, scripts...)
		marker = nextMarker
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("data", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetScriptsParams(d *schema.ResourceData, marker float64) string {
	res := "?limit=100"
	if v, ok := d.GetOk("name_like"); ok {
		res = fmt.Sprintf("%s&name_like=%v", res, v)
	}
	if v, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, v)
	}
	if v, ok := d.GetOk("risk_level"); ok {
		res = fmt.Sprintf("%s&risk_level=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if marker != 0 {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func flattenCocGetScripts(resp interface{}) ([]map[string]interface{}, float64) {
	scriptsJson := utils.PathSearch("data.data", resp, make([]interface{}, 0))
	scriptsArray := scriptsJson.([]interface{})
	if len(scriptsArray) == 0 {
		return nil, 0
	}

	result := make([]map[string]interface{}, 0, len(scriptsArray))
	var marker float64
	for _, script := range scriptsArray {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", script, nil),
			"name":                  utils.PathSearch("name", script, nil),
			"type":                  utils.PathSearch("type", script, nil),
			"creator":               utils.PathSearch("creator", script, nil),
			"creator_id":            utils.PathSearch("creator_id", script, nil),
			"operator":              utils.PathSearch("operator", script, nil),
			"gmt_created":           utils.PathSearch("gmt_created", script, nil),
			"gmt_modified":          utils.PathSearch("gmt_modified", script, nil),
			"status":                utils.PathSearch("status", script, nil),
			"script_uuid":           utils.PathSearch("script_uuid", script, nil),
			"usage_count":           utils.PathSearch("usage_count", script, nil),
			"properties":            flattenScriptProperties(utils.PathSearch("properties", script, nil)),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", script, nil),
		})
		marker = utils.PathSearch("id", script, float64(0)).(float64)
	}
	return result, marker
}

func flattenScriptProperties(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"risk_level": utils.PathSearch("risk_level", param, nil),
			"reviewers": flattenScriptPropertiesReviewers(
				utils.PathSearch("reviewers", param, make([]interface{}, 0)).([]interface{})),
			"version":  utils.PathSearch("version", param, nil),
			"protocol": utils.PathSearch("protocol", param, nil),
		},
	}

	return rst
}

func flattenScriptPropertiesReviewers(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"reviewer_id":   utils.PathSearch("reviewer_id", params, nil),
			"reviewer_name": utils.PathSearch("reviewer_name", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}
