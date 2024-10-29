package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

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

func dataSourceIotdaUpgradePackagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		upgradePackages []model.OtaPackageInfo
		limit           = int32(50)
		offset          int32
	)

	for {
		listOpts := model.ListOtaPackageInfoRequest{
			AppId:       utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			PackageType: d.Get("type").(string),
			ProductId:   utils.StringIgnoreEmpty(d.Get("product_id").(string)),
			Version:     utils.StringIgnoreEmpty(d.Get("version").(string)),
			Limit:       utils.Int32(limit),
			Offset:      &offset,
		}

		listResp, listErr := client.ListOtaPackageInfo(&listOpts)
		if listErr != nil {
			return diag.Errorf("error retrieving upgrade packages: %s", listErr)
		}

		if listResp == nil || listResp.Packages == nil {
			break
		}

		if len(*listResp.Packages) == 0 {
			break
		}

		upgradePackages = append(upgradePackages, *listResp.Packages...)
		//nolint:gosec
		offset += int32(len(*listResp.Packages))
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("packages", flattenPackages(upgradePackages)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPackages(packages []model.OtaPackageInfo) []interface{} {
	if len(packages) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(packages))
	for _, v := range packages {
		rst = append(rst, map[string]interface{}{
			"id":                      v.PackageId,
			"space_id":                v.AppId,
			"product_id":              v.ProductId,
			"type":                    v.PackageType,
			"version":                 v.Version,
			"description":             v.Description,
			"support_source_versions": v.SupportSourceVersions,
			"custom_info":             v.CustomInfo,
			"created_at":              v.CreateTime,
		})
	}

	return rst
}
