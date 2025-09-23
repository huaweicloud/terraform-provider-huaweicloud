package iotda

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

// @API IoTDA GET /v5/iot/{project_id}/ota-upgrades/packages
func DataSourceUpgradePackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIotdaUpgradePackagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of the upgrade package.`,
			},
			"space_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource space ID to which the upgrade packages belong.`,
			},
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the product ID associated with the upgrade package.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version number of the upgrade package.`,
			},
			"packages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the upgrade packages.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the upgrade package.`,
						},
						"space_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource space ID to which the upgrade package belongs.`,
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product ID associated with the upgrade package.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the upgrade package.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the upgrade package.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version number of the upgrade package.`,
						},
						"support_source_versions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of source versions that support the upgrade of this version package.`,
						},
						"custom_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The custom information pushed to the device.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the software and firmware packages are uploaded to the IoT platform.`,
						},
					},
				},
			},
		},
	}
}

func buildUpgradePackagesQueryParams(d *schema.ResourceData) string {
	rst := fmt.Sprintf("&package_type=%v", d.Get("type"))

	if v, ok := d.GetOk("space_id"); ok {
		rst += fmt.Sprintf("&app_id=%v", v)
	}

	if v, ok := d.GetOk("product_id"); ok {
		rst += fmt.Sprintf("&product_id=%v", v)
	}

	if v, ok := d.GetOk("version"); ok {
		rst += fmt.Sprintf("&version=%v", v)
	}

	return rst
}

func dataSourceIotdaUpgradePackagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v5/iot/{project_id}/ota-upgrades/packages?limit=50"
		product     = "iotda"
		allPackages []interface{}
		offset      = 0
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildUpgradePackagesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving upgrade packages: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		packages := utils.PathSearch("packages", respBody, make([]interface{}, 0)).([]interface{})
		if len(packages) == 0 {
			break
		}

		allPackages = append(allPackages, packages...)
		offset += len(packages)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("packages", flattenPackages(allPackages)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPackages(packages []interface{}) []interface{} {
	if len(packages) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(packages))
	for _, v := range packages {
		rst = append(rst, map[string]interface{}{
			"id":                      utils.PathSearch("package_id", v, nil),
			"space_id":                utils.PathSearch("app_id", v, nil),
			"product_id":              utils.PathSearch("product_id", v, nil),
			"type":                    utils.PathSearch("package_type", v, nil),
			"version":                 utils.PathSearch("version", v, nil),
			"description":             utils.PathSearch("description", v, nil),
			"support_source_versions": utils.PathSearch("support_source_versions", v, nil),
			"custom_info":             utils.PathSearch("custom_info", v, nil),
			"created_at":              utils.PathSearch("create_time", v, nil),
		})
	}

	return rst
}
