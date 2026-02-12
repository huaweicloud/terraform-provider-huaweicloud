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

// @API HSS GET /v5/{project_id}/plugins/code/{code}
func DataSourcePluginInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluginInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_arch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourcePluginInfoSchema(),
			},
		},
	}
}

func dataSourcePluginInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildPluginInfoQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("plugin_version"); ok {
		queryParams = fmt.Sprintf("%s&plugin_version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("agent_version"); ok {
		queryParams = fmt.Sprintf("%s&agent_version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("plugin_arch"); ok {
		queryParams = fmt.Sprintf("%s&plugin_arch=%v", queryParams, v)
	}
	if v, ok := d.GetOk("plugin_os_type"); ok {
		queryParams = fmt.Sprintf("%s&plugin_os_type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourcePluginInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/plugins/code/{code}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{code}", d.Get("code").(string))
	requestPath += buildPluginInfoQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS plugin info: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		dataListResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataListResp) == 0 {
			break
		}

		result = append(result, dataListResp...)
		offset += len(dataListResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenPluginInfoDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPluginInfoDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":                utils.PathSearch("name", v, nil),
			"id":                  utils.PathSearch("id", v, nil),
			"version":             utils.PathSearch("version", v, nil),
			"agent_version":       utils.PathSearch("agent_version", v, nil),
			"arch":                utils.PathSearch("arch", v, nil),
			"os_type":             utils.PathSearch("os_type", v, nil),
			"version_description": utils.PathSearch("version_description", v, nil),
			"size":                utils.PathSearch("size", v, nil),
			"cpu_limit":           utils.PathSearch("cpu_limit", v, nil),
			"memory_limit":        utils.PathSearch("memory_limit", v, nil),
			"update_time":         utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
