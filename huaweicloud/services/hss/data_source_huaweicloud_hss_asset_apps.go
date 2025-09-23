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

// @API HSS GET /v5/{project_id}/asset/apps
func DataSourceAssetApps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetAppsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the host ID.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the host name.",
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the software name.",
			},
			"host_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the host IP address.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the software version.",
			},
			"install_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the installation directory.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type.",
			},
			"part_match": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to use fuzzy matching.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
			"data_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appSchema(),
				Description: "The software list.",
			},
		},
	}
}

func appSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The agent ID.",
			},
			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host ID.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host name.",
			},
			"host_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host IP address.",
			},
			"app_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The software name.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version number.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The latest update time, in milliseconds.",
			},
			"recent_scan_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The latest scanning time, in milliseconds.",
			},
			"container_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The container ID.",
			},
			"container_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The container name.",
			},
		},
	}
	return &sc
}

func buildAppsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	queryParams := "?limit=100"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("app_name"); ok {
		queryParams = fmt.Sprintf("%s&app_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_ip"); ok {
		queryParams = fmt.Sprintf("%s&host_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("version"); ok {
		queryParams = fmt.Sprintf("%s&version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("install_dir"); ok {
		queryParams = fmt.Sprintf("%s&install_dir=%v", queryParams, v)
	}
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	if d.Get("part_match").(bool) {
		queryParams = fmt.Sprintf("%s&part_match=true", queryParams)
	}
	return queryParams
}

func flattenApps(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"agent_id":         utils.PathSearch("agent_id", v, nil),
			"host_id":          utils.PathSearch("host_id", v, nil),
			"host_name":        utils.PathSearch("host_name", v, nil),
			"host_ip":          utils.PathSearch("host_ip", v, nil),
			"app_name":         utils.PathSearch("app_name", v, nil),
			"version":          utils.PathSearch("version", v, nil),
			"update_time":      utils.PathSearch("update_time", v, nil),
			"recent_scan_time": utils.PathSearch("recent_scan_time", v, nil),
			"container_id":     utils.PathSearch("container_id", v, nil),
			"container_name":   utils.PathSearch("container_name", v, nil),
		})
	}
	return rst
}

func dataSourceAssetAppsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/asset/apps"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAppsQueryParams(d, cfg)
	allApps := make([]interface{}, 0)
	offset := 0

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS apps: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		appsResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(appsResp) == 0 {
			break
		}
		allApps = append(allApps, appsResp...)
		offset += len(appsResp)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data_list", flattenApps(allApps)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
