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

// @API DSC GET /v1/{project_id}/metadata/catalog/statistics
func DataSourceCatalogStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCatalogStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"label_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     catalogStatisticsObsStaticInfoSchema(),
			},
			"column": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     catalogStatisticsColumnStaticInfoSchema(),
			},
			"database": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     catalogStatisticsStaticInfoSchema(),
			},
			"file": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     catalogStatisticsObsStaticInfoSchema(),
			},
			"table": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     catalogStatisticsStaticInfoSchema(),
			},
		},
	}
}

func catalogStatisticsObsStaticInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"risk_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func catalogStatisticsColumnStaticInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"risk_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"week_on_week": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func catalogStatisticsStaticInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"risk_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"week_on_week": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCatalogStatisticsQueryParams(d *schema.ResourceData) string {
	queryParams := ""
	if v, ok := d.GetOk("label_id"); ok {
		queryParams = fmt.Sprintf("%s&label_id=%v", queryParams, v)
	}

	if v, ok := d.GetOk("type_id"); ok {
		queryParams = fmt.Sprintf("%s&type_id=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceCatalogStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/metadata/catalog/statistics"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	currentPath := requestPath + buildCatalogStatisticsQueryParams(d)
	resp, err := client.Request("GET", currentPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC catalog statistics: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bucket", flattenCatalogStatisticsObsStaticInfo(
			utils.PathSearch("bucket", respBody, nil))),
		d.Set("column", flattenCatalogStatisticsColumnStaticInfo(
			utils.PathSearch("column", respBody, nil))),
		d.Set("database", flattenCatalogStatisticsStaticInfo(
			utils.PathSearch("database", respBody, nil))),
		d.Set("file", flattenCatalogStatisticsObsStaticInfo(
			utils.PathSearch("file", respBody, nil))),
		d.Set("table", flattenCatalogStatisticsStaticInfo(
			utils.PathSearch("table", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCatalogStatisticsObsStaticInfo(obj interface{}) []interface{} {
	if obj == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"risk_num": utils.PathSearch("risk_num", obj, nil),
			"total":    utils.PathSearch("total", obj, nil),
		},
	}
}

func flattenCatalogStatisticsColumnStaticInfo(obj interface{}) []interface{} {
	if obj == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"risk_num":     utils.PathSearch("risk_num", obj, nil),
			"total":        utils.PathSearch("total", obj, nil),
			"week_on_week": utils.PathSearch("week_on_week", obj, nil),
		},
	}
}

func flattenCatalogStatisticsStaticInfo(obj interface{}) []interface{} {
	if obj == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"risk_num":     utils.PathSearch("risk_num", obj, nil),
			"total":        utils.PathSearch("total", obj, nil),
			"week_on_week": utils.PathSearch("week_on_week", obj, nil),
		},
	}
}
