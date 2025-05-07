package cci

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var imageSnapshotNonUpdatableParams = []string{"name", "annotations", "labels", "finalizers", "building_config",
	"building_config.*.auto_create_eip", "building_config.*.auto_create_eip_attribute",
	"building_config.*.auto_create_eip_attribute.*.bandwidth_charge_mode",
	"building_config.*.auto_create_eip_attribute.*.bandwidth_id",
	"building_config.*.auto_create_eip_attribute.*.bandwidth_size",
	"building_config.*.auto_create_eip_attribute.*.ip_version", "building_config.*.auto_create_eip_attribute.*.type",
	"building_config.*.eip_id", "building_config.*.namespace", "image_snapshot_size", "images", "images.*.image",
	"registries", "registries.*.image_pull_secret", "registries.*.insecure_skip_verify", "registries.*.plain_http",
	"registries.*.server", "ttl_days_after_created",
}

// @API CCI POST /apis/cci/v2/imagesnapshots
// @API CCI GET /apis/cci/v2/imagesnapshots/{name}
// @API CCI DELETE /apis/cci/v2/imagesnapshots/{name}
func ResourceV2ImageSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2ImageSnapshotCreate,
		ReadContext:   resourceV2ImageSnapshotRead,
		UpdateContext: resourceV2ImageSnapshotUpdate,
		DeleteContext: resourceV2ImageSnapshotDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(imageSnapshotNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"finalizers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"building_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_create_eip": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"auto_create_eip_attribute": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bandwidth_charge_mode": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"bandwidth_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"bandwidth_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"ip_version": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"eip_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"image_snapshot_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"images": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"registries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_pull_secret": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"insecure_skip_verify": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"plain_http": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ttl_days_after_created": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expire_date_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"images": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"digest": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size_bytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"last_updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phase": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceV2ImageSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createImageSnapshotHttpUrl := "apis/cci/v2/imagesnapshots"
	createImageSnapshotPath := client.Endpoint + createImageSnapshotHttpUrl
	createImageSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createImageSnapshotOpt.JSONBody = utils.RemoveNil(buildCreateV2ImageSnapshotParams(d))

	resp, err := client.Request("POST", createImageSnapshotPath, &createImageSnapshotOpt)
	if err != nil {
		return diag.Errorf("error creating CCI image snapshot: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if name == "" {
		return diag.Errorf("unable to find CCI image snapshot name from API response")
	}
	d.SetId(name)

	return resourceV2ImageSnapshotRead(ctx, d, meta)
}

func buildCreateV2ImageSnapshotParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":        d.Get("name"),
			"annotations": d.Get("annotations"),
			"labels":      d.Get("labels"),
		},
		"spec": buildImageSnapshotSpecParams(d),
	}

	return bodyParams
}

func buildImageSnapshotSpecParams(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"buildingConfig":      buildImageSnapshotSpecBuildingConfigParams(d.Get("building_config")),
		"imageSnapshotSize":   utils.ValueIgnoreEmpty(d.Get("image_snapshot_size")),
		"images":              buildImageSnapshotSpecImagesParams(d.Get("images")),
		"registries":          buildImageSnapshotSpecRegistriesParams(d.Get("registries")),
		"ttlDaysAfterCreated": utils.ValueIgnoreEmpty(d.Get("ttl_days_after_created")),
	}

	return rst
}

func buildImageSnapshotSpecBuildingConfigParams(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		if v, ok := rawArray[0].(map[string]interface{}); ok {
			rst := map[string]interface{}{
				"autoCreateEIP":          v["auto_create_eip"],
				"autoCreateEIPAttribute": buildImageSnapshotSpecBcAutoCreateEIPAttributeParams(v["auto_create_eip_attribute"]),
				"eipID":                  utils.ValueIgnoreEmpty(v["eip_id"]),
				"namespace":              utils.ValueIgnoreEmpty(v["namespace"]),
			}
			return rst
		}
	}

	return nil
}

func buildImageSnapshotSpecBcAutoCreateEIPAttributeParams(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		if v, ok := rawArray[0].(map[string]interface{}); ok {
			rst := map[string]interface{}{
				"bandwidthChargeMode": utils.ValueIgnoreEmpty(v["bandwidth_charge_mode"]),
				"bandwidthId":         utils.ValueIgnoreEmpty(v["bandwidth_id"]),
				"bandwidthSize":       utils.ValueIgnoreEmpty(v["bandwidth_size"]),
				"ipVersion":           utils.ValueIgnoreEmpty(v["ip_version"]),
				"type":                utils.ValueIgnoreEmpty(v["type"]),
			}
			return rst
		}
	}

	return nil
}

func buildImageSnapshotSpecImagesParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, 0, len(rawArray))
		for _, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst = append(rst, map[string]interface{}{
					"image": utils.ValueIgnoreEmpty(raw["image"]),
				})
			}
		}
		return rst
	}

	return nil
}

func buildImageSnapshotSpecRegistriesParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, 0, len(rawArray))
		for _, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst = append(rst, map[string]interface{}{
					"imagePullSecret":    utils.ValueIgnoreEmpty(raw["image_pull_secret"]),
					"insecureSkipVerify": utils.ValueIgnoreEmpty(raw["insecure_skip_verify"]),
					"plainHTTP":          utils.ValueIgnoreEmpty(raw["plain_http"]),
					"server":             utils.ValueIgnoreEmpty(raw["server"]),
				})
			}
		}
		return rst
	}

	return nil
}

func resourceV2ImageSnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	name := d.Get("name").(string)
	resp, err := GetImageSnapshot(client, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 image snapshot")
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("labels", utils.PathSearch("metadata.labels", resp, nil)),
		d.Set("finalizers", utils.PathSearch("metadata.finalizers", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("building_config", flattenImageSnapshotSpecBuildingConfig(resp)),
		d.Set("image_snapshot_size", utils.PathSearch("spec.imageSnapshotSize", resp, nil)),
		d.Set("images", flattenImageSnapshotSpecImages(resp)),
		d.Set("registries", flattenImageSnapshotSpecRegistries(resp)),
		d.Set("ttl_days_after_created", utils.PathSearch("spec.ttlDaysAfterCreated", resp, nil)),
		d.Set("status", flattenImageSnapshotStatus(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageSnapshotSpecBuildingConfig(resp interface{}) []interface{} {
	curJson := utils.PathSearch("spec.buildingConfig", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"auto_create_eip":           utils.PathSearch("autoCreateEIP", curJson, nil),
			"auto_create_eip_attribute": flattenImageSnapshotSpecBuildingConfigAutoCreateEipAttribute(curJson),
			"eip_id":                    utils.PathSearch("eipID", curJson, nil),
			"namespace":                 utils.PathSearch("namespace", curJson, nil),
		},
	}

	return rst
}

func flattenImageSnapshotSpecBuildingConfigAutoCreateEipAttribute(resp interface{}) []interface{} {
	curJson := utils.PathSearch("autoCreateEIPAttribute", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"bandwidth_charge_mode": utils.PathSearch("bandwidthChargeMode", curJson, nil),
			"bandwidth_id":          utils.PathSearch("bandwidthID", curJson, nil),
			"bandwidth_size":        utils.PathSearch("bandwidthSize", curJson, nil),
			"ip_version":            utils.PathSearch("ipVersion", curJson, nil),
			"type":                  utils.PathSearch("type", curJson, nil),
		},
	}

	return rst
}

func flattenImageSnapshotSpecImages(resp interface{}) []interface{} {
	curJson := utils.PathSearch("spec.images", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	if len(curArray) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"image": utils.PathSearch("image", v, nil),
		})
	}
	return rst
}

func flattenImageSnapshotSpecRegistries(resp interface{}) []interface{} {
	curJson := utils.PathSearch("spec.registries", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	if len(curArray) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"image_pull_secret":    utils.PathSearch("imagePullSecret", v, nil),
			"insecure_skip_verify": utils.PathSearch("insecureSkipVerify", v, nil),
			"plain_http":           utils.PathSearch("plainHttp", v, nil),
			"server":               utils.PathSearch("server", v, nil),
		})
	}

	return rst
}

func flattenImageSnapshotStatus(resp interface{}) []interface{} {
	curJson := utils.PathSearch("status", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"expire_date_time":  utils.PathSearch("expireDateTime", curJson, nil),
			"images":            flattenImageSnapshotStatusImages(curJson),
			"last_updated_time": utils.PathSearch("lastUpdatedTime", curJson, nil),
			"message":           utils.PathSearch("message", curJson, nil),
			"phase":             utils.PathSearch("phase", curJson, nil),
			"reason":            utils.PathSearch("reason", curJson, nil),
			"snapshot_id":       utils.PathSearch("snapshotID", curJson, nil),
			"snapshot_name":     utils.PathSearch("snapshotName", curJson, nil),
		},
	}

	return rst
}

func flattenImageSnapshotStatusImages(resp interface{}) []interface{} {
	curJson := utils.PathSearch("images", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	if len(curArray) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"digest":     utils.PathSearch("digest", v, nil),
			"image":      utils.PathSearch("image", v, nil),
			"size_bytes": utils.PathSearch("sizeBytes", v, nil),
		})
	}

	return rst
}

func resourceV2ImageSnapshotUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2ImageSnapshotDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	name := d.Get("name").(string)
	deleteImageSnapshotHttpUrl := "apis/cci/v2/imagesnapshots/{name}"
	deleteImageSnapshotPath := client.Endpoint + deleteImageSnapshotHttpUrl
	deleteImageSnapshotPath = strings.ReplaceAll(deleteImageSnapshotPath, "{name}", name)
	deleteImageSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteImageSnapshotPath, &deleteImageSnapshotOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 image snapshot: %s", err)
	}

	return nil
}

func GetImageSnapshot(client *golangsdk.ServiceClient, name string) (interface{}, error) {
	getImageSnapshotHttpUrl := "apis/cci/v2/imagesnapshots/{name}"
	getImageSnapshotPath := client.Endpoint + getImageSnapshotHttpUrl
	getImageSnapshotPath = strings.ReplaceAll(getImageSnapshotPath, "{name}", name)
	getImageSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getImageSnapshotResp, err := client.Request("GET", getImageSnapshotPath, &getImageSnapshotOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getImageSnapshotResp)
}
