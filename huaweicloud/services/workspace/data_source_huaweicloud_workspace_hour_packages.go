package workspace

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

// @API Workspace GET /v2/{project_id}/products/hour-packages
func DataSourceHourPackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHourPackagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to obtain the hour packages.`,
			},
			"desktop_resource_spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The specification code of desktop resource to be queried.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The specification code of hour package to be queried.`,
			},
			"hour_packages": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        hourPackageSchema(),
				Description: `The list of hour package information that matched filter parameters.`,
			},
		},
	}
}

func hourPackageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cloud_service_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of cloud service.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of resource.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of hour package.`,
			},
			"desktop_resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of desktop resource.`,
			},
			"descriptions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        hourPackageDescriptionSchema(),
				Description: `The descriptions of hour package.`,
			},
			"package_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The duration of hour package.`,
			},
			"domain_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of domain IDs supported by the hour package.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of hour package.`,
			},
		},
	}
}

func hourPackageDescriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"zh_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese description of hour package.`,
			},
			"en_us": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English description of hour package.`,
			},
		},
	}
}

func flattenHourPackages(packages []interface{}) []interface{} {
	if len(packages) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(packages))
	for _, item := range packages {
		result = append(result, map[string]interface{}{
			"cloud_service_type":         utils.PathSearch("cloud_service_type", item, nil),
			"resource_type":              utils.PathSearch("resource_type", item, nil),
			"resource_spec_code":         utils.PathSearch("resource_spec_code", item, nil),
			"desktop_resource_spec_code": utils.PathSearch("desktop_resource_spec_code", item, nil),
			"descriptions":               flattenHourPackageDescriptions(utils.PathSearch("descriptions", item, nil).(map[string]interface{})),
			"package_duration":           utils.PathSearch("package_duration", item, nil),
			"domain_ids":                 utils.PathSearch("domain_ids", item, nil),
			"status":                     utils.PathSearch("status", item, nil),
		})
	}
	return result
}

func flattenHourPackageDescriptions(descriptions map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"zh_cn": utils.PathSearch("zh_cn", descriptions, nil),
			"en_us": utils.PathSearch("en_us", descriptions, nil),
		},
	}
}

func buildHourPackagesParams(d *schema.ResourceData) string {
	res := "?"

	if v, ok := d.GetOk("desktop_resource_spec_code"); ok {
		res = fmt.Sprintf("%sdesktop_resource_spec_code=%v&", res, v)
	}
	if v, ok := d.GetOk("resource_spec_code"); ok {
		res = fmt.Sprintf("%sresource_spec_code=%v", res, v)
	}

	return res
}

func queryHourPackages(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/products/hour-packages"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildHourPackagesParams(d)

	requestOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, requestOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	packages := utils.PathSearch("hour_packages", respBody, make([]interface{}, 0)).([]interface{})

	return packages, nil
}

func dataSourceHourPackagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	packages, err := queryHourPackages(client, d)
	if err != nil {
		return diag.Errorf("error querying hour packages: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("hour_packages", flattenHourPackages(packages)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
