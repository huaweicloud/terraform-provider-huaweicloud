package waf

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

// @API WAF GET /v1/{project_id}/premium-waf/pools
func DataSourceWafPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafPoolsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detail": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     wafPoolItemSchema(),
			},
		},
	}
}

func wafPoolItemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourcePoolIdNameEntrySchema(),
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourcePoolIdNameEntrySchema(),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcePoolIdNameEntrySchema() *schema.Resource {
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
			"service_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildWafPoolsQueryParams(d *schema.ResourceData, cfg *config.Config, page, pageSize int) string {
	queryParams := fmt.Sprintf("?page=%d&pagesize=%d", page, pageSize)

	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}

	typeInput := d.Get("type").([]interface{})
	if len(typeInput) > 0 {
		for _, v := range utils.ExpandToStringList(typeInput) {
			queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
		}
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		queryParams = fmt.Sprintf("%s&vpc_id=%v", queryParams, v)
	}

	if v := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "detail"); v != nil {
		queryParams = fmt.Sprintf("%s&detail=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceWafPoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/premium-waf/pools"
		product  = "waf"
		page     = 1
		pageSize = 100
		result   = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildWafPoolsQueryParams(d, cfg, page, pageSize)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving WAF pools: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(items) == 0 {
			break
		}

		result = append(result, items...)

		page++
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", flattenWafPoolsItems(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWafPoolsItems(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(items))
	for _, v := range items {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"region":      utils.PathSearch("region", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"vpc_id":      utils.PathSearch("vpc_id", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"hosts": flattenDataSourcePoolIdNameEntries(
				utils.PathSearch("hosts", v, make([]interface{}, 0))),
			"instances": flattenDataSourcePoolIdNameEntries(
				utils.PathSearch("instances", v, make([]interface{}, 0))),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"create_time":           utils.PathSearch("create_time", v, nil),
		})
	}

	return rst
}

func flattenDataSourcePoolIdNameEntries(raw interface{}) []map[string]interface{} {
	entries, ok := raw.([]interface{})
	if !ok || len(entries) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(entries))
	for _, v := range entries {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"service_ip": utils.PathSearch("service_ip", v, nil),
		})
	}

	return rst
}
