package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var imageBatchScanNonUpdatableParams = []string{
	"image_type", "repo_type", "image_info_list", "image_info_list.*.namespace", "image_info_list.*.image_name",
	"image_info_list.*.image_version", "image_info_list.*.instance_id", "image_info_list.*.instance_url",
	"operate_all", "namespace", "image_name", "image_version", "scan_status", "latest_version", "image_size",
	"start_latest_update_time", "end_latest_update_time", "start_latest_scan_time", "end_latest_scan_time",
	"enterprise_project_id",
}

// @API HSS POST /v5/{project_id}/image/batch-scan
func ResourceImageBatchScan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageBatchScanCreate,
		ReadContext:   resourceImageBatchScanRead,
		UpdateContext: resourceImageBatchScanUpdate,
		DeleteContext: resourceImageBatchScanDelete,

		CustomizeDiff: config.FlexibleForceNew(imageBatchScanNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_info_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"image_version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"operate_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scan_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"latest_version": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"image_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start_latest_update_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_latest_update_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start_latest_scan_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_latest_scan_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildImageBatchScanBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"image_type": d.Get("image_type").(string),
	}

	if d.Get("operate_all").(bool) {
		bodyParams["operate_all"] = true
	}

	if v, ok := d.GetOk("image_info_list"); ok {
		imageInfoList := v.([]interface{})
		if len(imageInfoList) == 0 {
			return nil
		}

		imageInfos := make([]map[string]interface{}, len(imageInfoList))
		for i, info := range imageInfoList {
			infoMap := info.(map[string]interface{})
			imageInfo := map[string]interface{}{
				"namespace":     infoMap["namespace"],
				"image_name":    infoMap["image_name"],
				"image_version": infoMap["image_version"],
				"instance_id":   infoMap["instance_id"],
				"instance_url":  infoMap["instance_url"],
			}

			imageInfos[i] = imageInfo
		}
		bodyParams["image_info_list"] = imageInfos
	}

	optionalParams := []string{
		"repo_type", "namespace", "image_name", "image_version", "scan_status", "latest_version", "image_size",
		"start_latest_update_time", "end_latest_update_time", "start_latest_scan_time", "end_latest_scan_time",
	}

	for _, param := range optionalParams {
		if v, ok := d.GetOk(param); ok {
			bodyParams[param] = v
		}
	}

	return bodyParams
}

func resourceImageBatchScanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/image/batch-scan"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildImageBatchScanBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch scanning images: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceImageBatchScanRead(ctx, d, meta)
}

func resourceImageBatchScanRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageBatchScanUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageBatchScanDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch scan images. Deleting this resource
    will not clear the corresponding scan records, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
