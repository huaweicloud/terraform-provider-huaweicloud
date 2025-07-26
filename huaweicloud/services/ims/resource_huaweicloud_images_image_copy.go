package ims

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IMS POST /v1/cloudimages/{image_id}/copy
// @API IMS POST /v1/cloudimages/{image_id}/cross_region_copy
// @API IMS GET /v1/{project_id}jobs/{job_id}
// @API IMS PATCH /v2/cloudimages/{image_id}
// @API IMS POST /v2/{project_id}/images/{image_id}/tags/action
// @API IMS GET /v2/cloudimages
// @API IMS GET /v2/{project_id}/images/{image_id}/tags
// @API IMS DELETE /v2/images/{image_id}
func ResourceImsImageCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsImageCopyCreate,
		UpdateContext: resourceImsImageCopyUpdate,
		ReadContext:   resourceImsImageCopyRead,
		DeleteContext: resourceImsImageCopyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the source image.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the copy image.`,
			},
			"target_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the target region name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the copy image.`,
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the master key used for encrypting an image.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project id of the image.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the agency name.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the vault.`,
			},
			"max_ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			// following are additional attributes
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__os_bit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__isregistered": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"virtual_env_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__image_source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__imagetype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__originalimagename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__productcode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__lazyloading": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__image_displayname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__os_feature_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_xen": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_largememory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_diskintensive": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_highperformance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_xen_gpu_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm_gpu_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_xen_hana": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm_infiniband": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__system_support_market": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"__is_offshelved": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__root_origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__sequence_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_fc_inject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hw_firmware_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hw_vif_multiqueue_enabled": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_arm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_agent_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__account_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_amd": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm_hi1822_hisriov": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__support_kvm_hi1822_hivirtionet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_shutdown_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Deprecated attributes
			"checksum": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: utils.SchemaDesc("checksum is deprecated", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateWithinRegionCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"description":           d.Get("description"),
		"cmk_id":                utils.ValueIgnoreEmpty(d.Get("kms_key_id")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
	}

	return bodyParams
}

func buildCreateCrossRegionCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"description":  d.Get("description"),
		"region":       utils.ValueIgnoreEmpty(d.Get("target_region")),
		"project_name": utils.ValueIgnoreEmpty(d.Get("target_region")),
		"agency_name":  utils.ValueIgnoreEmpty(d.Get("agency_name")),
		"vault_id":     utils.ValueIgnoreEmpty(d.Get("vault_id")),
	}

	return bodyParams
}

func resourceImsImageCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                     = meta.(*config.Config)
		sourceRegion            = cfg.GetRegion(d)
		targetRegion            = d.Get("target_region").(string)
		sourceImageId           = d.Get("source_image_id").(string)
		product                 = "ims"
		withinRegionCopyHttpUrl = "v1/cloudimages/{image_id}/copy"
		crossRegionCopyHttpUrl  = "v1/cloudimages/{image_id}/cross_region_copy"
		jobId                   string
	)

	client, err := cfg.NewServiceClient(product, sourceRegion)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	if targetRegion == "" || targetRegion == sourceRegion {
		createPath := client.Endpoint + withinRegionCopyHttpUrl
		createPath = strings.ReplaceAll(createPath, "{image_id}", sourceImageId)
		createOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         utils.RemoveNil(buildCreateWithinRegionCopyBodyParams(d)),
		}

		createResp, err := client.Request("POST", createPath, &createOpt)
		if err != nil {
			return diag.Errorf("error creating IMS image copy within region: %s", err)
		}

		createRespBody, err := utils.FlattenResponse(createResp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobId = utils.PathSearch("job_id", createRespBody, "").(string)
		if jobId == "" {
			return diag.Errorf("error creating IMS image copy within region: job ID is not found in API response")
		}

	} else {
		createPath := client.Endpoint + crossRegionCopyHttpUrl
		createPath = strings.ReplaceAll(createPath, "{image_id}", sourceImageId)
		createOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         utils.RemoveNil(buildCreateCrossRegionCopyBodyParams(d)),
		}

		createResp, err := client.Request("POST", createPath, &createOpt)
		if err != nil {
			return diag.Errorf("error creating IMS image copy cross region: %s", err)
		}

		createRespBody, err := utils.FlattenResponse(createResp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobId = utils.PathSearch("job_id", createRespBody, "").(string)
		if jobId == "" {
			return diag.Errorf("error creating IMS image copy cross region: job ID is not found in API response")
		}
	}

	// Wait for the copy image to become available.
	imageId, err := waitForCreateImageCopyJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for IMS image copy to complete: %s", err)
	}

	d.SetId(imageId)

	copiedRegionClient, err := getCopiedRegionClient(d, cfg)
	if err != nil {
		return diag.Errorf("error creating IMS copied region client: %s", err)
	}

	// Set `max_ram` and `min_ram` attributes.
	updateOpts := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("max_ram"); ok {
		updateOpts = append(updateOpts, map[string]interface{}{
			"op":    "replace",
			"path":  "/max_ram",
			"value": v,
		})
	}

	if v, ok := d.GetOk("min_ram"); ok {
		updateOpts = append(updateOpts, map[string]interface{}{
			"op":    "replace",
			"path":  "/min_ram",
			"value": v,
		})
	}

	if len(updateOpts) > 0 {
		updatePath := copiedRegionClient.Endpoint + "v2/cloudimages/{image_id}"
		updatePath = strings.ReplaceAll(updatePath, "{image_id}", imageId)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         updateOpts,
		}

		_, err = copiedRegionClient.Request("PATCH", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error setting IMS image attributes in creation operation: %s", err)
		}
	}

	// Set tags.
	if _, ok := d.GetOk("tags"); ok {
		err = updateIMSImageTags(copiedRegionClient, d)
		if err != nil {
			return diag.Errorf("error setting IMS image tags field: %s", err)
		}
	}

	return resourceImsImageCopyRead(ctx, d, meta)
}

func waitForCreateImageCopyJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      imageCopyJobStatusRefreshFunc(jobId, client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	getRespBody, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for IMS image copy job (%s) to succeed: %s", jobId, err)
	}

	imageId := utils.PathSearch("entities.image_id", getRespBody, "").(string)
	if imageId == "" {
		return "", errors.New("the image ID is not found in API response")
	}

	return imageId, nil
}

func imageCopyJobStatusRefreshFunc(jobId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return getResp, "ERROR", fmt.Errorf("error retrieving IMS image copy job: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return getRespBody, "ERROR", err
		}

		status := utils.PathSearch("status", getRespBody, "").(string)
		if status == "SUCCESS" {
			return getRespBody, "COMPLETED", nil
		}

		if status == "FAIL" {
			return getRespBody, "COMPLETED", errors.New("the image copy job execution failed")
		}

		if status == "" {
			return getRespBody, "ERROR", errors.New("status field is not found in API response")
		}

		return getRespBody, "PENDING", nil
	}
}

func getImageCopy(client *golangsdk.ServiceClient, imageId string) (interface{}, error) {
	// If the `enterprise_project_id` is not filled, the list API will query images under all enterprise projects.
	// So there's no need to fill `enterprise_project_id` here.
	getPath := client.Endpoint + "v2/cloudimages"
	getPath += fmt.Sprintf("?id=%s", imageId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS image copy: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("images[0]", getRespBody, nil), nil
}

func resourceImsImageCopyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := getCopiedRegionClient(d, cfg)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	image, err := getImageCopy(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// If the list API return empty, then process `CheckDeleted` logic.
	if image == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS image copy")
	}

	dataOrigin := utils.PathSearch("__data_origin", image, "").(string)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", image, nil)),
		d.Set("description", utils.PathSearch("__description", image, nil)),
		d.Set("kms_key_id", utils.PathSearch("__system__cmkid", image, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", image, nil)),
		d.Set("max_ram", flattenMaxRAM(utils.PathSearch("max_ram", image, "").(string))),
		d.Set("min_ram", utils.PathSearch("min_ram", image, nil)),
		d.Set("tags", flattenIMSImageTags(client, d.Id())),
		d.Set("instance_id", flattenSpecificValueFormDataOrigin(dataOrigin, "instance")),
		d.Set("file", utils.PathSearch("file", image, nil)),
		d.Set("self", utils.PathSearch("self", image, nil)),
		d.Set("schema", utils.PathSearch("schema", image, nil)),
		d.Set("status", utils.PathSearch("status", image, nil)),
		d.Set("visibility", utils.PathSearch("visibility", image, nil)),
		d.Set("protected", utils.PathSearch("protected", image, nil)),
		d.Set("container_format", utils.PathSearch("container_format", image, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", image, nil)),
		d.Set("__os_bit", utils.PathSearch("__os_bit", image, nil)),
		d.Set("os_version", utils.PathSearch("__os_version", image, nil)),
		d.Set("disk_format", utils.PathSearch("disk_format", image, nil)),
		d.Set("__isregistered", utils.PathSearch("__isregistered", image, nil)),
		d.Set("__platform", utils.PathSearch("__platform", image, nil)),
		d.Set("__os_type", utils.PathSearch("__os_type", image, nil)),
		d.Set("min_disk", utils.PathSearch("min_disk", image, nil)),
		d.Set("virtual_env_type", utils.PathSearch("virtual_env_type", image, nil)),
		d.Set("__image_source_type", utils.PathSearch("__image_source_type", image, nil)),
		d.Set("__imagetype", utils.PathSearch("__imagetype", image, nil)),
		d.Set("created_at", utils.PathSearch("created_at", image, nil)),
		d.Set("__originalimagename", utils.PathSearch("__originalimagename", image, nil)),
		d.Set("__backup_id", utils.PathSearch("__backup_id", image, nil)),
		d.Set("__productcode", utils.PathSearch("__productcode", image, nil)),
		d.Set("image_size", utils.PathSearch("__image_size", image, nil)),
		d.Set("data_origin", dataOrigin),
		d.Set("__lazyloading", utils.PathSearch("__lazyloading", image, nil)),
		d.Set("active_at", utils.PathSearch("active_at", image, nil)),
		d.Set("__image_displayname", utils.PathSearch("__image_displayname", image, nil)),
		d.Set("__os_feature_list", utils.PathSearch("__os_feature_list", image, nil)),
		d.Set("__support_kvm", utils.PathSearch("__support_kvm", image, nil)),
		d.Set("__support_xen", utils.PathSearch("__support_xen", image, nil)),
		d.Set("__support_largememory", utils.PathSearch("__support_largememory", image, nil)),
		d.Set("__support_diskintensive", utils.PathSearch("__support_diskintensive", image, nil)),
		d.Set("__support_highperformance", utils.PathSearch("__support_highperformance", image, nil)),
		d.Set("__support_xen_gpu_type", utils.PathSearch("__support_xen_gpu_type", image, nil)),
		d.Set("__support_kvm_gpu_type", utils.PathSearch("__support_kvm_gpu_type", image, nil)),
		d.Set("__support_xen_hana", utils.PathSearch("__support_xen_hana", image, nil)),
		d.Set("__support_kvm_infiniband", utils.PathSearch("__support_kvm_infiniband", image, nil)),
		d.Set("__system_support_market", utils.PathSearch("__system_support_market", image, nil)),
		d.Set("__is_offshelved", utils.PathSearch("__is_offshelved", image, nil)),
		d.Set("__root_origin", utils.PathSearch("__root_origin", image, nil)),
		d.Set("__sequence_num", utils.PathSearch("__sequence_num", image, nil)),
		d.Set("__support_fc_inject", utils.PathSearch("__support_fc_inject", image, nil)),
		d.Set("hw_firmware_type", utils.PathSearch("hw_firmware_type", image, nil)),
		d.Set("hw_vif_multiqueue_enabled", utils.PathSearch("hw_vif_multiqueue_enabled", image, nil)),
		d.Set("__support_arm", utils.PathSearch("__support_arm", image, nil)),
		d.Set("__support_agent_list", utils.PathSearch("__support_agent_list", image, nil)),
		d.Set("__account_code", utils.PathSearch("__account_code", image, nil)),
		d.Set("__support_amd", utils.PathSearch("__support_amd", image, nil)),
		d.Set("__support_kvm_hi1822_hisriov", utils.PathSearch("__support_kvm_hi1822_hisriov", image, nil)),
		d.Set("__support_kvm_hi1822_hivirtionet", utils.PathSearch("__support_kvm_hi1822_hivirtionet", image, nil)),
		d.Set("os_shutdown_timeout", utils.PathSearch("os_shutdown_timeout", image, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceImsImageCopyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/cloudimages/{image_id}"
		imageId = d.Id()
	)

	client, err := getCopiedRegionClient(d, cfg)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{image_id}", imageId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	if d.HasChange("name") {
		bodyParams := []map[string]interface{}{
			{
				"op":    "replace",
				"path":  "/name",
				"value": d.Get("name"),
			},
		}

		updateOpt.JSONBody = bodyParams
		_, err = client.Request("PATCH", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating IMS image copy name field: %s", err)
		}
	}

	if d.HasChange("description") {
		bodyParams := []map[string]interface{}{
			{
				"op":    "replace",
				"path":  "/__description",
				"value": d.Get("description"),
			},
		}

		updateOpt.JSONBody = bodyParams
		_, err = client.Request("PATCH", updatePath, &updateOpt)
		if err != nil {
			err = processUpdateDescriptionError(d, client, err)
			if err != nil {
				return diag.Errorf("error updating IMS image copy description field: %s", err)
			}
		}
	}

	if d.HasChange("max_ram") {
		bodyParams := []map[string]interface{}{
			{
				"op":    "replace",
				"path":  "/max_ram",
				"value": d.Get("max_ram"),
			},
		}

		updateOpt.JSONBody = bodyParams
		_, err = client.Request("PATCH", updatePath, &updateOpt)
		if err != nil {
			err = processUpdateMaxRAMError(d, client, err)
			if err != nil {
				return diag.Errorf("error updating IMS image copy max_ram field: %s", err)
			}
		}
	}

	if d.HasChange("min_ram") {
		bodyParams := []map[string]interface{}{
			{
				"op":    "replace",
				"path":  "/min_ram",
				"value": d.Get("min_ram"),
			},
		}

		updateOpt.JSONBody = bodyParams
		_, err = client.Request("PATCH", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating IMS image copy min_ram field: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateIMSImageTags(client, d)
		if err != nil {
			return diag.Errorf("error updating IMS image copy tags field: %s", err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   imageId,
			ResourceType: "images",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceImsImageCopyRead(ctx, d, meta)
}

func resourceImsImageCopyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/images/{image_id}"
		imageId = d.Id()
	)

	client, err := getCopiedRegionClient(d, cfg)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	// Before deleting, call the query API first, if the query result is empty, then process `CheckDeleted` logic.
	image, err := getImageCopy(client, imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	if image == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS image copy")
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{image_id}", imageId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting IMS image copy: %s", err)
	}

	// Because the delete API always return `204` status code,
	// so we need to call the list query API to check if the image has been successfully deleted.
	err = waitForImageCopyDeleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for IMS image copy to be deleted: %s", err)
	}

	return nil
}

func waitForImageCopyDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			image, err := getImageCopy(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			if image == nil {
				return "SUCCESS", "COMPLETED", nil
			}

			return image, "PENDING", nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func getCopiedRegionClient(d *schema.ResourceData, cfg *config.Config) (*golangsdk.ServiceClient, error) {
	imageRegion := cfg.GetRegion(d)
	if v, ok := d.GetOk("target_region"); ok {
		imageRegion = v.(string)
	}

	imsClient, err := cfg.NewServiceClient("ims", imageRegion)
	if err != nil {
		return nil, err
	}

	return imsClient, nil
}
