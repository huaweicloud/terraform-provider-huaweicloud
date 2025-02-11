package iotda

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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
									"sign": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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

func buildObsLocationParams(rawParams interface{}) map[string]interface{} {
	if rawParams == nil {
		return nil
	}

	rawObsLocation := rawParams.([]interface{})
	if len(rawObsLocation) < 1 {
		return nil
	}

	obsLocation := rawObsLocation[0].(map[string]interface{})
	obsLocationParams := map[string]interface{}{
		"region_name": obsLocation["region"],
		"bucket_name": obsLocation["bucket_name"],
		"object_key":  obsLocation["object_key"],
		"sign":        utils.ValueIgnoreEmpty(obsLocation["sign"]),
	}

	return obsLocationParams
}

func buildFileLocationParam(rawFileLocation []interface{}) map[string]interface{} {
	if len(rawFileLocation) == 0 {
		return nil
	}

	obsLocation := rawFileLocation[0].(map[string]interface{})
	fileLocationParam := map[string]interface{}{
		"obs_location": buildObsLocationParams(obsLocation["obs_location"]),
	}

	return fileLocationParam
}

func buildUpgradePackageCreateParams(d *schema.ResourceData) map[string]interface{} {
	packageParams := map[string]interface{}{
		"app_id":                  d.Get("space_id"),
		"package_type":            d.Get("type"),
		"product_id":              d.Get("product_id"),
		"version":                 d.Get("version"),
		"file_location":           buildFileLocationParam(d.Get("file_location").([]interface{})),
		"support_source_versions": utils.ExpandToStringList(d.Get("support_source_versions").([]interface{})),
		"description":             utils.ValueIgnoreEmpty(d.Get("description")),
		"custom_info":             utils.ValueIgnoreEmpty(d.Get("custom_info")),
	}

	return packageParams
}

func resourceUpgradePackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/ota-upgrades/packages"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpgradePackageCreateParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA OTA upgrade package: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	packageId := utils.PathSearch("package_id", respBody, "").(string)
	if packageId == "" {
		return diag.Errorf("error creating IoTDA OTA upgrade package: ID is not found in API response")
	}

	d.SetId(packageId)

	return resourceUpgradePackageRead(ctx, d, meta)
}

func resourceUpgradePackageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/ota-upgrades/packages/{package_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{package_id}", d.Id())
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		// When the resource does not exist, query API will return `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA OTA upgrade package")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("space_id", utils.PathSearch("app_id", getRespBody, nil)),
		d.Set("type", utils.PathSearch("package_type", getRespBody, nil)),
		d.Set("product_id", utils.PathSearch("product_id", getRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getRespBody, nil)),
		d.Set("file_location", flattenFileLocation(utils.PathSearch("file_location", getRespBody, nil))),
		d.Set("support_source_versions", utils.PathSearch("support_source_versions", getRespBody, make([]interface{}, 0))),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("custom_info", utils.PathSearch("custom_info", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFileLocation(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	fileLocation := map[string]interface{}{
		"obs_location": flattenObsLocation(utils.PathSearch("obs_location", resp, nil)),
	}

	return []map[string]interface{}{fileLocation}
}

func flattenObsLocation(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	obsLocation := map[string]interface{}{
		"region":      utils.PathSearch("region_name", resp, nil),
		"bucket_name": utils.PathSearch("bucket_name", resp, nil),
		"object_key":  utils.PathSearch("object_key", resp, nil),
		"sign":        utils.PathSearch("sign", resp, nil),
	}

	return []map[string]interface{}{obsLocation}
}

func resourceUpgradePackageDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/ota-upgrades/packages/{package_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{package_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// When the resource does not exist, delete API will return `404`.
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA OTA upgrade package")
	}

	return nil
}
