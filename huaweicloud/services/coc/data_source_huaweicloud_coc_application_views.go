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

// @API COC GET /v1/application-view/search
func DataSourceCocApplicationViews() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocApplicationViewsRead,

		Schema: map[string]*schema.Schema{
			"name_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_collection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Elem:     applicationViewsDataSchema(),
				Computed: true,
			},
		},
	}
}

func applicationViewsDataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceCocApplicationViewsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/application-view/search"
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

	var marker string
	res := make([]map[string]interface{}, 0)
	for {
		getPath := basePath + buildGetApplicationViewsParams(d, marker)
		getResp, err := client.Request("GET", getPath, &getOpt)

		if err != nil {
			return diag.Errorf("error retrieving COC application views: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		views, nextMarker := flattenCocGetApplicationViews(getRespBody)
		if len(views) < 1 {
			break
		}
		res = append(res, views...)
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

func buildGetApplicationViewsParams(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("name_like"); ok {
		res = fmt.Sprintf("%s&name_like=%v", res, v)
	}
	if v, ok := d.GetOk("code_list"); ok {
		res += buildQueryStringParams("code_list", v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("is_collection"); ok {
		res = fmt.Sprintf("%s&is_collection=%v", res, v)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func buildQueryStringParams(queryKey string, queryValues []interface{}) string {
	rst := ""
	for _, val := range queryValues {
		if queryValue, ok := val.(string); ok && queryValue != "" {
			rst += fmt.Sprintf("&%s=%s", queryKey, queryValue)
		}
	}

	return rst
}

func flattenCocGetApplicationViews(resp interface{}) ([]map[string]interface{}, string) {
	if resp == nil {
		return nil, ""
	}
	viewsJson := utils.PathSearch("data", resp, make([]interface{}, 0))
	viewsArray := viewsJson.([]interface{})
	if len(viewsArray) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(viewsArray))
	marker := ""
	for _, view := range viewsArray {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", view, nil),
			"name":              utils.PathSearch("name", view, nil),
			"code":              utils.PathSearch("code", view, nil),
			"type":              utils.PathSearch("type", view, nil),
			"parent_id":         utils.PathSearch("parent_id", view, nil),
			"component_id":      utils.PathSearch("component_id", view, nil),
			"application_id":    utils.PathSearch("application_id", view, nil),
			"path":              utils.PathSearch("path", view, nil),
			"vendor":            utils.PathSearch("vendor", view, nil),
			"related_domain_id": utils.PathSearch("related_domain_id", view, nil),
		})
		marker = utils.PathSearch("id", view, "").(string)
	}
	return result, marker
}
