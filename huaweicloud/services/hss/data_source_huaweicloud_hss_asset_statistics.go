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

// @API HSS GET /v5/{project_id}/asset/statistics
func DataSourceAssetStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetStatisticsRead,
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
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type. The default value is host.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the resource belongs.",
			},
			"account_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of server accounts.",
			},
			"port_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of open ports.",
			},
			"process_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of processes.",
			},
			"app_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of applications.",
			},
			"auto_launch_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of auto launch startup processes.",
			},
			"web_framework_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of web frameworks.",
			},
			"web_site_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of web sites.",
			},
			"jar_package_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of JAR packages.",
			},
			"kernel_module_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of kernel modules.",
			},
			"web_service_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of web services.",
			},
			"web_app_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of web applications.",
			},
			"database_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of databases.",
			},
		},
	}
}

func buildAssetStatisticsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	queryParams := "?limit=100"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	return queryParams
}

func dataSourceAssetStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	requestPath := client.Endpoint + "v5/{project_id}/asset/statistics"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAssetStatisticsQueryParams(d, cfg)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS asset statistics: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("account_num", utils.PathSearch("account_num", respBody, nil)),
		d.Set("port_num", utils.PathSearch("port_num", respBody, nil)),
		d.Set("process_num", utils.PathSearch("process_num", respBody, nil)),
		d.Set("app_num", utils.PathSearch("app_num", respBody, nil)),
		d.Set("auto_launch_num", utils.PathSearch("auto_launch_num", respBody, nil)),
		d.Set("web_framework_num", utils.PathSearch("web_framework_num", respBody, nil)),
		d.Set("web_site_num", utils.PathSearch("web_site_num", respBody, nil)),
		d.Set("jar_package_num", utils.PathSearch("jar_package_num", respBody, nil)),
		d.Set("kernel_module_num", utils.PathSearch("kernel_module_num", respBody, nil)),
		d.Set("web_service_num", utils.PathSearch("web_service_num", respBody, nil)),
		d.Set("web_app_num", utils.PathSearch("web_app_num", respBody, nil)),
		d.Set("database_num", utils.PathSearch("database_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
