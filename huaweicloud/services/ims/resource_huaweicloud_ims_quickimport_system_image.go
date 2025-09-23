package ims

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var quickImportSystemImageNonUpdatableParams = []string{"os_version", "image_url", "min_disk", "type", "architecture",
	"license_type"}

// @API IMS POST /v2/cloudimages/quickimport/action
// @API IMS GET /v1/{project_id}/jobs/{job_id}
// @API IMS GET /v2/cloudimages
// @API IMS GET /v2/{project_id}/images/{image_id}/tags
// @API IMS PATCH /v2/cloudimages/{image_id}
// @API IMS POST /v2/{project_id}/images/{image_id}/tags/action
// @API IMS DELETE /v2/images/{image_id}
func ResourceQuickImportSystemImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQuickImportSystemImageCreate,
		ReadContext:   resourceQuickImportSystemImageRead,
		UpdateContext: resourceQuickImportSystemImageUpdate,
		DeleteContext: resourceQuickImportSystemImageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(quickImportSystemImageNonUpdatableParams),
			config.MergeDefaultTags(),
		),

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
			"os_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"min_disk": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// The `description` field can be left blank, so the `Computed` attribute is not used.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hw_firmware_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			// Attributes.
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
			"min_ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_ram": {
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
			"__image_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__data_origin": {
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
			"__system__cmkid": {
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
		},
	}
}

func buildCreateQuickImportSystemImageBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":       d.Get("name"),
		"os_version": d.Get("os_version"),
		"image_url":  d.Get("image_url"),
		"min_disk":   d.Get("min_disk"),
		// The `description` field can be left blank.
		"description":           d.Get("description"),
		"hw_firmware_type":      utils.ValueIgnoreEmpty(d.Get("hw_firmware_type")),
		"type":                  utils.ValueIgnoreEmpty(d.Get("type")),
		"architecture":          utils.ValueIgnoreEmpty(d.Get("architecture")),
		"license_type":          utils.ValueIgnoreEmpty(d.Get("license_type")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"image_tags":            utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}
}

func resourceQuickImportSystemImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/cloudimages/quickimport/action"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateQuickImportSystemImageBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IMS quick import system image: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating IMS quick import system image: job ID is not found in API response")
	}

	imageId, err := waitForQuickImportSystemImageJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for IMS quick import system image to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceQuickImportSystemImageRead(ctx, d, meta)
}

func waitForQuickImportSystemImageJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    quickImportSystemImageJobStatusRefreshFunc(jobId, client),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for IMS quick import system image job (%s) to succeed: %s", jobId, err)
	}

	return getQuickImportSystemImageIdByJob(client, jobId)
}

func quickImportSystemImageJobStatusRefreshFunc(jobId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return getResp, "ERROR", fmt.Errorf("error retrieving IMS quick import system image job: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return getRespBody, "ERROR", err
		}

		status := utils.PathSearch("status", getRespBody, "").(string)
		if status == "SUCCESS" {
			return "SUCCESS", "COMPLETED", nil
		}

		if status == "FAIL" {
			return getRespBody, "COMPLETED", errors.New("the quick import system image job execution has failed")
		}

		if status == "" {
			return getRespBody, "ERROR", errors.New("status field is not found in API response")
		}

		return getRespBody, "PENDING", nil
	}
}

func getQuickImportSystemImageIdByJob(client *golangsdk.ServiceClient, jobId string) (string, error) {
	getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return "", fmt.Errorf("error retrieving IMS quick import system image job: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return "", err
	}

	imageId := utils.PathSearch("entities.image_id", getRespBody, "").(string)
	if imageId == "" {
		return "", errors.New("the image ID is not found in API response")
	}

	return imageId, nil
}

func resourceQuickImportSystemImageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/cloudimages"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath += fmt.Sprintf("?id=%s", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving IMS quick import system image: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	image := utils.PathSearch("images[0]", getRespBody, nil)
	// If the list API return empty, then process `CheckDeleted` logic.
	if image == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving IMS quick import system image")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", image, nil)),
		d.Set("os_version", utils.PathSearch("__os_version", image, nil)),
		d.Set("min_disk", utils.PathSearch("min_disk", image, nil)),
		d.Set("description", utils.PathSearch("__description", image, nil)),
		d.Set("hw_firmware_type", utils.PathSearch("hw_firmware_type", image, nil)),
		d.Set("tags", flattenIMSImageTags(client, d.Id())),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", image, nil)),
		d.Set("file", utils.PathSearch("file", image, nil)),
		d.Set("self", utils.PathSearch("self", image, nil)),
		d.Set("schema", utils.PathSearch("schema", image, nil)),
		d.Set("status", utils.PathSearch("status", image, nil)),
		d.Set("visibility", utils.PathSearch("visibility", image, nil)),
		d.Set("protected", utils.PathSearch("protected", image, nil)),
		d.Set("container_format", utils.PathSearch("container_format", image, nil)),
		d.Set("min_ram", utils.PathSearch("min_ram", image, nil)),
		d.Set("max_ram", utils.PathSearch("max_ram", image, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", image, nil)),
		d.Set("__os_bit", utils.PathSearch("__os_bit", image, nil)),
		d.Set("disk_format", utils.PathSearch("disk_format", image, nil)),
		d.Set("__isregistered", utils.PathSearch("__isregistered", image, nil)),
		d.Set("__platform", utils.PathSearch("__platform", image, nil)),
		d.Set("__os_type", utils.PathSearch("__os_type", image, nil)),
		d.Set("virtual_env_type", utils.PathSearch("virtual_env_type", image, nil)),
		d.Set("__image_source_type", utils.PathSearch("__image_source_type", image, nil)),
		d.Set("__imagetype", utils.PathSearch("__imagetype", image, nil)),
		d.Set("created_at", utils.PathSearch("created_at", image, nil)),
		d.Set("__originalimagename", utils.PathSearch("__originalimagename", image, nil)),
		d.Set("__backup_id", utils.PathSearch("__backup_id", image, nil)),
		d.Set("__productcode", utils.PathSearch("__productcode", image, nil)),
		d.Set("__image_size", utils.PathSearch("__image_size", image, nil)),
		d.Set("__data_origin", utils.PathSearch("__data_origin", image, nil)),
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
		d.Set("hw_vif_multiqueue_enabled", utils.PathSearch("hw_vif_multiqueue_enabled", image, nil)),
		d.Set("__support_arm", utils.PathSearch("__support_arm", image, nil)),
		d.Set("__support_agent_list", utils.PathSearch("__support_agent_list", image, nil)),
		d.Set("__system__cmkid", utils.PathSearch("__system__cmkid", image, nil)),
		d.Set("__account_code", utils.PathSearch("__account_code", image, nil)),
		d.Set("__support_amd", utils.PathSearch("__support_amd", image, nil)),
		d.Set("__support_kvm_hi1822_hisriov", utils.PathSearch("__support_kvm_hi1822_hisriov", image, nil)),
		d.Set("__support_kvm_hi1822_hivirtionet", utils.PathSearch("__support_kvm_hi1822_hivirtionet", image, nil)),
		d.Set("os_shutdown_timeout", utils.PathSearch("os_shutdown_timeout", image, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// When querying tags, even if there is an error when calling the API, it does not affect the resources and only prints
// log reminders.
func flattenIMSImageTags(client *golangsdk.ServiceClient, imageId string) map[string]interface{} {
	getPath := client.Endpoint + "v2/{project_id}/images/{image_id}/tags"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{image_id}", fmt.Sprintf("%v", imageId))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		log.Printf("[WARN] failed to retrieve IMS image tags: %s", err)

		return nil
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		log.Printf("[WARN] failed to flatten IMS image tags: %s", err)

		return nil
	}

	return utils.FlattenTagsToMap(utils.PathSearch("tags", getRespBody, nil))
}

func buildCreateOrDeleteImageTagsBodyParams(tagsMap map[string]interface{}, action string) map[string]interface{} {
	tagsList := make([]map[string]interface{}, 0, len(tagsMap))

	for k, v := range tagsMap {
		tagsList = append(tagsList, map[string]interface{}{
			"key":   k,
			"value": v,
		})
	}

	return map[string]interface{}{
		"action": action,
		"tags":   tagsList,
	}
}

func updateIMSImageTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		oRaw, nRaw = d.GetChange("tags")
		oMap       = oRaw.(map[string]interface{})
		nMap       = nRaw.(map[string]interface{})
	)

	updateTagsPath := client.Endpoint + "v2/{project_id}/images/{image_id}/tags/action"
	updateTagsPath = strings.ReplaceAll(updateTagsPath, "{project_id}", client.ProjectID)
	updateTagsPath = strings.ReplaceAll(updateTagsPath, "{image_id}", d.Id())

	// Delete old tags.
	if len(oMap) > 0 {
		updateTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			OkCodes:          []int{200, 202, 204},
			JSONBody:         buildCreateOrDeleteImageTagsBodyParams(oMap, "delete"),
		}

		_, err := client.Request("POST", updateTagsPath, &updateTagsOpt)
		if err != nil {
			return fmt.Errorf("failed to delete old tags: %s", err)
		}
	}

	// Create new tags.
	if len(nMap) > 0 {
		updateTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			OkCodes:          []int{200, 202, 204},
			JSONBody:         buildCreateOrDeleteImageTagsBodyParams(nMap, "create"),
		}

		_, err := client.Request("POST", updateTagsPath, &updateTagsOpt)
		if err != nil {
			return fmt.Errorf("failed to create new tags: %s", err)
		}
	}

	return nil
}

func buildUpdateQuickImportSystemImageBodyParams(d *schema.ResourceData) []map[string]interface{} {
	updateFields := map[string]string{
		"name":             "/name",
		"description":      "/__description",
		"hw_firmware_type": "/hw_firmware_type",
	}

	var bodyParams []map[string]interface{}

	for field, path := range updateFields {
		if d.HasChange(field) {
			bodyParams = append(bodyParams, map[string]interface{}{
				"op":    "replace",
				"path":  path,
				"value": d.Get(field),
			})
		}
	}

	return bodyParams
}

func resourceQuickImportSystemImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/cloudimages/{image_id}"
		imageId = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	updateBodyParams := buildUpdateQuickImportSystemImageBodyParams(d)
	if len(updateBodyParams) > 0 {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{image_id}", imageId)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         updateBodyParams,
		}

		_, err = client.Request("PATCH", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating IMS quick import system image: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateIMSImageTags(client, d)
		if err != nil {
			return diag.Errorf("error updating IMS quick import system image tags: %s", err)
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

	return resourceQuickImportSystemImageRead(ctx, d, meta)
}

func resourceQuickImportSystemImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/images/{image_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{image_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting IMS image: %s", err)
	}

	// Because the delete API always return `204` status code,
	// so we need to call the list query API to check if the image has been successfully deleted.
	err = waitForQuickImportSystemImageDeleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for IMS quick import system image to be deleted: %s", err)
	}

	return nil
}

func waitForQuickImportSystemImageDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getPath := client.Endpoint + "v2/cloudimages"
			getPath += fmt.Sprintf("?id=%s", d.Id())
			getOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}

			getResp, err := client.Request("GET", getPath, &getOpt)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error retrieving IMS quick import system image: %s", err)
			}

			getRespBody, err := utils.FlattenResponse(getResp)
			if err != nil {
				return nil, "ERROR", err
			}

			image := utils.PathSearch("images[0]", getRespBody, nil)
			if image == nil {
				return "SUCCESS", "COMPLETED", nil
			}

			return image, "PENDING", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
