package dew

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

// @API DEW GET /v1/{project_id}/dew/cpcs/apps
func DataSourceCpcsApps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCpcsAppsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application name.`,
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VPC name.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort attribute.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort direction.`,
			},
			"apps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the applications.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application ID.`,
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application name.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC ID to which the application belongs.`,
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC name to which the application belongs.`,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subnet ID to which the application belongs.`,
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subnet name to which the application belongs.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The account ID.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application description.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The creation time of the application.`,
						},
					},
				},
			},
		},
	}
}

func buildDataSourceCpcsAppsQueryParams(d *schema.ResourceData, pageNum int) string {
	rst := ""

	if v, ok := d.GetOk("app_name"); ok {
		rst += fmt.Sprintf("&app_name=%v", v)
	}

	if v, ok := d.GetOk("vpc_name"); ok {
		rst += fmt.Sprintf("&vpc_name=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	if pageNum > 0 {
		rst += fmt.Sprintf("&page_num=%d", pageNum)
	}

	if len(rst) > 0 {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceCpcsAppsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/apps"
		product = "kms"
		pageNum = 0
		allApps = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithPageNum := requestPath + buildDataSourceCpcsAppsQueryParams(d, pageNum)
		resp, err := client.Request("GET", requestPathWithPageNum, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DEW CPCS applications: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.Errorf("error flattening DEW CPCS applications response: %s", err)
		}

		results := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
		if len(results) == 0 {
			break
		}

		allApps = append(allApps, results...)
		pageNum += len(results)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(generateId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("apps", flattenCpcsAppsResponseBody(allApps)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCpcsAppsResponseBody(apps []interface{}) []interface{} {
	if len(apps) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(apps))
	for _, app := range apps {
		result = append(result, map[string]interface{}{
			"app_id":      utils.PathSearch("app_id", app, nil),
			"app_name":    utils.PathSearch("app_name", app, nil),
			"vpc_id":      utils.PathSearch("vpc_id", app, nil),
			"vpc_name":    utils.PathSearch("vpc_name", app, nil),
			"subnet_id":   utils.PathSearch("subnet_id", app, nil),
			"subnet_name": utils.PathSearch("subnet_name", app, nil),
			"domain_id":   utils.PathSearch("domain_id", app, nil),
			"description": utils.PathSearch("description", app, nil),
			"create_time": utils.PathSearch("create_time", app, nil),
		})
	}

	return result
}
