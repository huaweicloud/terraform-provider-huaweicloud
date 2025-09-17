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

// @API WAF GET /v1/{project_id}/waf/overviews/attack/url
func DataSourceOverviewsAttackUrl() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOverviewsAttackUrlRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"from": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"top": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hosts": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildOverviewsAttackUrlQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?from=%v&to=%v", d.Get("from"), d.Get("to"))

	if v, ok := d.GetOk("top"); ok {
		queryParams = fmt.Sprintf("%s&top=%v", queryParams, v)
	}

	if v, ok := d.GetOk("hosts"); ok {
		queryParams = fmt.Sprintf("%s&hosts=%v", queryParams, v)
	}

	if v, ok := d.GetOk("instances"); ok {
		queryParams = fmt.Sprintf("%s&instances=%v", queryParams, v)
	}

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceOverviewsAttackUrlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/overviews/attack/url"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildOverviewsAttackUrlQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving attacked URLs: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	items := utils.PathSearch("items", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", flattenOverviewsAttackUrl(items)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOverviewsAttackUrl(subscriptions []interface{}) []interface{} {
	result := make([]interface{}, 0, len(subscriptions))
	for _, v := range subscriptions {
		result = append(result, map[string]interface{}{
			"key":  utils.PathSearch("key", v, nil),
			"num":  utils.PathSearch("num", v, nil),
			"host": utils.PathSearch("host", v, nil),
		})
	}

	return result
}
