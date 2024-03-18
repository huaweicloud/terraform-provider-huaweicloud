package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/ota-upgrades/packages
// @API IoTDA GET /v5/iot/{project_id}/ota-upgrades/packages/{package_id}
// @API IoTDA DELETE /v5/iot/{project_id}/ota-upgrades/packages/{package_id}
func ResourceUpgradePackage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUpgradePackageCreate,
		ReadContext:   resourceUpgradePackageRead,
		DeleteContext: resourceUpgradePackageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file_location": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"obs_location": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"bucket_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"object_key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"support_source_versions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"custom_info": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildObsLocationParam(raw interface{}) *model.ObsLocation {
	obsLocationMap := raw.([]interface{})[0].(map[string]interface{})
	obsLocationParam := model.ObsLocation{
		RegionName: obsLocationMap["region"].(string),
		BucketName: obsLocationMap["bucket_name"].(string),
		ObjectKey:  obsLocationMap["object_key"].(string),
	}

	return &obsLocationParam
}

func buildFileLocationParam(raw []interface{}) *model.FileLocation {
	if raw[0] == nil {
		return nil
	}

	obsLocationParam := buildObsLocationParam(raw[0].(map[string]interface{})["obs_location"])
	fileLocationParam := model.FileLocation{
		ObsLocation: obsLocationParam,
	}

	return &fileLocationParam
}

func buildUpgradePackageCreateParams(d *schema.ResourceData) *model.CreateOtaPackageRequest {
	req := model.CreateOtaPackageRequest{
		Body: &model.CreateOtaPackage{
			AppId:                 d.Get("space_id").(string),
			PackageType:           d.Get("type").(string),
			ProductId:             d.Get("product_id").(string),
			Version:               d.Get("version").(string),
			FileLocation:          buildFileLocationParam(d.Get("file_location").([]interface{})),
			SupportSourceVersions: utils.ExpandToStringListPointer(d.Get("support_source_versions").([]interface{})),
			Description:           utils.StringIgnoreEmpty(d.Get("description").(string)),
			CustomInfo:            utils.StringIgnoreEmpty(d.Get("custom_info").(string)),
		},
	}

	return &req
}

func resourceUpgradePackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := buildUpgradePackageCreateParams(d)
	resp, err := client.CreateOtaPackage(createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA OTA upgrade package: %s", err)
	}

	if resp == nil || resp.PackageId == nil {
		return diag.Errorf("error creating IoTDA OTA upgrade package: ID is not found in API response")
	}

	d.SetId(*resp.PackageId)

	return resourceUpgradePackageRead(ctx, d, meta)
}

func resourceUpgradePackageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	getOpts := &model.ShowOtaPackageRequest{
		PackageId: d.Id(),
	}
	resp, err := client.ShowOtaPackage(getOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying IoTDA OTA upgrade package")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("space_id", resp.AppId),
		d.Set("type", resp.PackageType),
		d.Set("product_id", resp.ProductId),
		d.Set("version", resp.Version),
		d.Set("file_location", flattenFileLocation(resp.FileLocation)),
		d.Set("support_source_versions", resp.SupportSourceVersions),
		d.Set("description", resp.Description),
		d.Set("custom_info", resp.CustomInfo),
		d.Set("created_at", resp.CreateTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFileLocation(resp *model.FileLocation) []interface{} {
	if resp == nil {
		return nil
	}

	obsLocation := flattenObsLocation(resp.ObsLocation)
	fileLocation := []interface{}{
		map[string]interface{}{
			"obs_location": obsLocation,
		},
	}

	return fileLocation
}

func flattenObsLocation(resp *model.ObsLocation) []interface{} {
	if resp == nil {
		return nil
	}

	obsLocation := []interface{}{
		map[string]interface{}{
			"region":      resp.RegionName,
			"bucket_name": resp.BucketName,
			"object_key":  resp.ObjectKey,
		},
	}

	return obsLocation
}

func resourceUpgradePackageDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpts := model.DeleteOtaPackageRequest{
		PackageId: d.Id(),
	}
	_, err = client.DeleteOtaPackage(&deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA OTA upgrade package: %s", err)
	}

	return nil
}
