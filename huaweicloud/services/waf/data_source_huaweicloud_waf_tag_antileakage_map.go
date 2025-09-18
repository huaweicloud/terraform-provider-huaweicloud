package waf

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

// @API WAF GET /v1/{project_id}/waf/tag/antileakage/map
func DataSourceTagAntileakageMap() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagAntileakageMapRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Required: true,
			},
			"leakagemap": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sensitive": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"code": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"locale": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id_card": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sensitive": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// The API does not return this parameter, and named `responseCode` in API document
						"responsecode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTagAntileakageMapRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/tag/antileakage/map"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?lang=%v", listPath, d.Get("lang"))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving sensitive information options: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
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
		d.Set("leakagemap", flattenTagAntileakageMapLeakage(utils.PathSearch("leakagemap", listRespBody, nil))),
		d.Set("locale", flattenTagAntileakageMapLocale(utils.PathSearch("locale", listRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTagAntileakageMapLeakage(contentType interface{}) []interface{} {
	if contentType == nil {
		return nil
	}

	result := map[string]interface{}{
		"sensitive": utils.PathSearch("sensitive", contentType, nil),
		"code":      utils.PathSearch("code", contentType, nil),
	}

	return []interface{}{result}
}

func flattenTagAntileakageMapLocale(optionType interface{}) []interface{} {
	if optionType == nil {
		return nil
	}

	result := map[string]interface{}{
		"code":         utils.PathSearch("code", optionType, nil),
		"id_card":      utils.PathSearch("id_card", optionType, nil),
		"sensitive":    utils.PathSearch("sensitive", optionType, nil),
		"phone":        utils.PathSearch("phone", optionType, nil),
		"responsecode": utils.PathSearch("responseCode", optionType, nil),
		"email":        utils.PathSearch("email", optionType, nil),
	}

	return []interface{}{result}
}
