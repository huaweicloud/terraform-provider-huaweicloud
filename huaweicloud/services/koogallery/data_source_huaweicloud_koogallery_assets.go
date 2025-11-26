package koogallery

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/koogallery/v1/assets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API KooGallery GET /api/mkp-openapi-public/v1/asset/deployed-object
func DataSourceKooGalleryAssets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKooGalleryAssetsRead,
		Description: "schema: Internal",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"asset_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deployed_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"asset_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deployed_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_deployed_object": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"create_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"architecture": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"software_pkg_deployed_object": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"package_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"internal_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKooGalleryAssetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.KooGalleryV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating KooGallery client: %s", err)
	}

	listOpts := assets.AssetsOpts{
		AssetId:      d.Get("asset_id").(string),
		DeployedType: d.Get("deployed_type").(string),
		Region:       d.Get("region").(string),
		AssetVersion: d.Get("asset_version").(string),
	}

	assetsList, err := assets.List(client, listOpts).Extract()
	if err != nil {
		return diag.Errorf(err.Error())
	}

	var resultAssets []interface{}
	var ids []string

	for _, asset := range assetsList {
		resultAssets = append(resultAssets, flattenAssets(asset))
		ids = append(ids, asset.AssetId)
	}

	mErr := d.Set("assets", resultAssets)
	if mErr != nil {
		return diag.Errorf("set assets data err:%s", mErr)
	}

	d.SetId(hashcode.Strings(ids))

	return nil
}

func flattenImgDeployedObj(obj assets.ImageDeployedObj) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"image_id":     obj.ImageId,
			"image_name":   obj.ImageName,
			"os_type":      obj.OsType,
			"create_time":  obj.CreateTime,
			"architecture": obj.Architecture,
		},
	}
}

func flattenSwPkgDeployedObj(obj assets.SoftwarePkgDeployedObj) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"package_name":  obj.PackageName,
			"internal_path": obj.InternalPath,
			"checksum":      obj.Checksum,
		},
	}
}

func flattenAssets(asset assets.Data) map[string]interface{} {
	return map[string]interface{}{
		"asset_id":                     asset.AssetId,
		"deployed_type":                asset.DeployedType,
		"version":                      asset.Version,
		"version_id":                   asset.VersionId,
		"region":                       asset.Region,
		"image_deployed_object":        flattenImgDeployedObj(asset.ImgDeployedObj),
		"software_pkg_deployed_object": flattenSwPkgDeployedObj(asset.SwPkgDeployedObj),
	}
}
