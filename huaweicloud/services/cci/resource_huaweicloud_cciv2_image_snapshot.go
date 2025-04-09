package cci

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var imageSnapshotNonUpdatableParams = []string{"name"}

// @API CCI POST /apis/cci/v2/imagesnapshots
// @API CCI GET /apis/cci/v2/imagesnapshots/{name}
// @API CCI PUT /apis/cci/v2/imagesnapshots/{name}
// @API CCI DELETE /apis/cci/v2/imagesnapshots/{name}
func ResourceV2ImageSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2ImageSnapshotCreate,
		ReadContext:   resourceV2ImageSnapshotRead,
		UpdateContext: resourceV2ImageSnapshotUpdate,
		DeleteContext: resourceV2ImageSnapshotDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2ImageSnapshotImportState,
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
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CCI Image Snapshot.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the annotations of the CCI Image Snapshot.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the annotations of the CCI Image Snapshot.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI Image Snapshot.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI Image Snapshot.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the CCI Image Snapshot.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the CCI Image Snapshot.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the CCI Image Snapshot.`,
			},
			"building_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_create_eip": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Specifies whether to auto create EIP.`,
						},
						"auto_create_eip_attribute": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bandwidth_charge_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the bandwidth charge mode of EIP.`,
									},
									"bandwidth_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the ID of EIP.`,
									},
									"bandwidth_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the bandwidth size of EIP.`,
									},
									"ip_version": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the IP version used by pod.`,
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the type of EIP.`,
									},
								},
							},
							Description: `Specifies whether to auto create EIP.`,
						},
						"eip_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the EIP ID.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the namespace.`,
						},
					},
				},
				Description: `Specifies the building config of the CCI Image Snapshot.`,
			},
			"image_snapshot_size": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the size of the CCI Image Snapshot.`,
			},
			"images": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the image name.`,
						},
					},
				},
				Description: `The images list of references to images to make image snapshot.`,
			},
			"registries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_pull_secret": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the image pull secret.`,
						},
						"insecure_skip_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Specifies whether to allow connections to SSL sites without certs.`,
						},
						"plain_http": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Specifies whether the server uses http protocol.`,
						},
						"server": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the image repository server.`,
						},
					},
				},
				Description: `Specifies the registries list.`,
			},
			"ttl_days_after_created": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The TTL days after created.`,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expire_date_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expire date time.`,
						},
						"images": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"digest": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The image digest.`,
									},
									"image": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The image name.`,
									},
									"size_bytes": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The size of the image in bytes.`,
									},
								},
							},
							Description: `The status.`,
						},
						"last_updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last updated time.`,
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The message.`,
						},
						"phase": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The phase.`,
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The reason.`,
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The snapshot ID.`,
						},
						"snapshot_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The snapshot name.`,
						},
					},
				},
				Description: `The status.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
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
		return diag.Errorf("error creating CCI ImageSnapshot: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if name == "" {
		return diag.Errorf("unable to find CCI ImageSnapshot name or namespace from API response")
	}
	d.SetId(name)

	return resourceV2ImageSnapshotRead(ctx, d, meta)
}

func buildCreateV2ImageSnapshotParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
		"spec":   map[string]interface{}{},
		"status": map[string]interface{}{},
	}

	return bodyParams
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
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 ImageSnapshot")
	}

	mErr := multierror.Append(
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("building_config", flattenSpecBuildingConfig(utils.PathSearch("spec.buildingConfig", resp, nil))),
		d.Set("image_snapshot_size", utils.PathSearch("spec.imageSnapshotSize", resp, nil)),
		d.Set("images", flattenSpecImages(utils.PathSearch("spec.images", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("registries", flattenSpecRegistries(utils.PathSearch("spec.registries", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("ttl_days_after_created", utils.PathSearch("spec.ipFamilyPolicy", resp, nil)),
		d.Set("status", flattenImageSnapshotStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSpecBuildingConfig(buildingConfig interface{}) []interface{} {
	if buildingConfig == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"auto_create_eip":           utils.PathSearch("autoCreateEIP", buildingConfig, nil),
		"auto_create_eip_attribute": flattenAutoCreateEIPAttribute(utils.PathSearch("autoCreateEIPAttribute", buildingConfig, nil)),
		"eip_id":                    utils.PathSearch("eipID", buildingConfig, nil),
		"namespace":                 utils.PathSearch("namespace", buildingConfig, nil),
	})

	return rst
}

func flattenAutoCreateEIPAttribute(autoCreateEIPAttribute interface{}) []interface{} {
	if autoCreateEIPAttribute == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"bandwidth_charge_mode": utils.PathSearch("bandwidthChargeMode", autoCreateEIPAttribute, nil),
		"bandwidth_id":          utils.PathSearch("bandwidthID", autoCreateEIPAttribute, nil),
		"bandwidth_size":        utils.PathSearch("bandwidthSize", autoCreateEIPAttribute, nil),
		"ip_version":            utils.PathSearch("ipVersion", autoCreateEIPAttribute, nil),
		"type":                  utils.PathSearch("type", autoCreateEIPAttribute, nil),
	})

	return rst
}

func flattenSpecImages(images []interface{}) []interface{} {
	if len(images) == 0 {
		return nil
	}

	rst := make([]interface{}, len(images))
	for i, v := range images {
		rst[i] = map[string]interface{}{
			"image": utils.PathSearch("image", v, nil),
		}
	}

	return rst
}

func flattenSpecRegistries(registries []interface{}) []interface{} {
	if len(registries) == 0 {
		return nil
	}

	rst := make([]interface{}, len(registries))
	for i, v := range registries {
		rst[i] = map[string]interface{}{
			"image_pull_secret":    utils.PathSearch("imagePullSecret", v, nil),
			"insecure_skip_verify": utils.PathSearch("insecureSkipVerify", v, nil),
			"plain_http":           utils.PathSearch("plainHTTP", v, nil),
			"server":               utils.PathSearch("server", v, nil),
		}
	}

	return rst
}

func flattenImageSnapshotStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"expire_date_time":  utils.PathSearch("expireDateTime", status, nil),
			"images":            flattenImageSnapshotStatusImages(utils.PathSearch("images", status, make([]interface{}, 0)).([]interface{})),
			"last_updated_time": utils.PathSearch("lastUpdatedTime", status, nil),
			"message":           utils.PathSearch("message", status, nil),
			"phase":             utils.PathSearch("phase", status, nil),
			"reason":            utils.PathSearch("reason", status, nil),
			"snapshot_id":       utils.PathSearch("snapshotID", status, nil),
			"snapshot_name":     utils.PathSearch("snapshotName", status, nil),
		},
	}

	return rst
}

func flattenImageSnapshotStatusImages(images []interface{}) []interface{} {
	if len(images) == 0 {
		return nil
	}

	rst := make([]interface{}, len(images))
	for i, v := range images {
		rst[i] = map[string]interface{}{
			"digest":     utils.PathSearch("digest", v, nil),
			"image":      utils.PathSearch("image", v, nil),
			"size_bytes": utils.PathSearch("sizeBytes", v, nil),
		}
	}

	return rst
}

func resourceV2ImageSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	updateImageSnapshotHttpUrl := "apis/cci/v2/imagesnapshots/{name}"
	updateImageSnapshotPath := client.Endpoint + updateImageSnapshotHttpUrl
	updateImageSnapshotPath = strings.ReplaceAll(updateImageSnapshotPath, "{name}", d.Get("name").(string))
	updateImageSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateImageSnapshotOpt.JSONBody = utils.RemoveNil(buildUpdateV2ImageSnapshotParams(d))

	_, err = client.Request("PUT", updateImageSnapshotPath, &updateImageSnapshotOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 ImageSnapshot: %s", err)
	}
	return resourceV2ImageSnapshotRead(ctx, d, meta)
}

func buildUpdateV2ImageSnapshotParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       d.Get("kind"),
		"apiVersion": d.Get("api_version"),
		"metadata": map[string]interface{}{
			"name":              d.Get("name"),
			"namespace":         d.Get("namespace"),
			"uid":               d.Get("uid"),
			"resourceVersion":   d.Get("resource_version"),
			"creationTimestamp": d.Get("creation_timestamp"),
			"annotations":       d.Get("annotations"),
		},
		"spec":   map[string]interface{}{},
		"status": map[string]interface{}{},
	}

	return bodyParams
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
		return diag.Errorf("error deleting CCI v2 ImageSnapshot: %s", err)
	}

	return nil
}

func resourceV2ImageSnapshotImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())

	return []*schema.ResourceData{d}, nil
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
		return getImageSnapshotResp, err
	}

	return utils.FlattenResponse(getImageSnapshotResp)
}
