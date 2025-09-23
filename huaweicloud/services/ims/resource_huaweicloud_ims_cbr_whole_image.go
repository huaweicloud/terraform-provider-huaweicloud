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

// @API IMS POST /v1/cloudimages/wholeimages/action
// @API IMS GET /v1/{project_id}/jobs/{job_id}
// @API IMS GET /v2/cloudimages
// @API IMS GET /v2/{project_id}/images/{image_id}/tags
// @API IMS PATCH /v2/cloudimages/{image_id}
// @API IMS POST /v2/{project_id}/images/{image_id}/tags/action
// @API IMS DELETE /v2/images/{image_id}
func ResourceCbrWholeImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCbrWholeImageCreate,
		ReadContext:   resourceCbrWholeImageRead,
		UpdateContext: resourceCbrWholeImageUpdate,
		DeleteContext: resourceCbrWholeImageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

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
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The `description` field can be left blank, so the `Computed` attribute is not used.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"is_delete_backup": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Attributes
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
			"__productcode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__image_size": {
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

func buildCreateCbrWholeImageBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"backup_id":             d.Get("backup_id"),
		"description":           d.Get("description"),
		"max_ram":               utils.ValueIgnoreEmpty(d.Get("max_ram")),
		"min_ram":               utils.ValueIgnoreEmpty(d.Get("min_ram")),
		"image_tags":            utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"whole_image_type":      "CBR",
	}

	return bodyParams
}

func resourceCbrWholeImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v1/cloudimages/wholeimages/action"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateCbrWholeImageBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IMS CBR whole image: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating IMS CBR whole image: job ID is not found in API response")
	}

	imageId, err := waitForCreateCbrWholeImageJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for IMS CBR whole image to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceCbrWholeImageRead(ctx, d, meta)
}

func waitForCreateCbrWholeImageJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      cbrWholeImageJobStatusRefreshFunc(jobId, client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	getRespBody, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for IMS CBR whole image job (%s) to succeed: %s", jobId, err)
	}

	imageId := utils.PathSearch("entities.image_id", getRespBody, "").(string)
	if imageId == "" {
		return "", errors.New("the image ID is not found in API response")
	}

	return imageId, nil
}

func cbrWholeImageJobStatusRefreshFunc(jobId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return getResp, "ERROR", fmt.Errorf("error retrieving IMS CBR whole image job: %s", err)
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
			return getRespBody, "COMPLETED", errors.New("the CBR whole image creation job execution failed")
		}

		if status == "" {
			return getRespBody, "ERROR", errors.New("status field is not found in API response")
		}

		return getRespBody, "PENDING", nil
	}
}

func getCbrWholeImage(client *golangsdk.ServiceClient, imageId string) (interface{}, error) {
	// If the `enterprise_project_id` is not filled, the list API will query images under all enterprise projects.
	// So there's no need to fill `enterprise_project_id` here.
	getPath := client.Endpoint + "v2/cloudimages"
	getPath += fmt.Sprintf("?id=%s", imageId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS CBR whole image: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("images[0]", getRespBody, nil), nil
}

func resourceCbrWholeImageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	image, err := getCbrWholeImage(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// If the list API return empty, then process `CheckDeleted` logic.
	if image == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving IMS CBR whole image")
	}

	dataOrigin := utils.PathSearch("__data_origin", image, "").(string)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", image, nil)),
		d.Set("backup_id", utils.PathSearch("__backup_id", image, nil)),
		d.Set("description", utils.PathSearch("__description", image, nil)),
		d.Set("max_ram", flattenMaxRAM(utils.PathSearch("max_ram", image, "").(string))),
		d.Set("min_ram", utils.PathSearch("min_ram", image, nil)),
		d.Set("tags", flattenIMSImageTags(client, d.Id())),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", image, nil)),
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
		d.Set("__productcode", utils.PathSearch("__productcode", image, nil)),
		d.Set("__image_size", utils.PathSearch("__image_size", image, nil)),
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
		d.Set("__system__cmkid", utils.PathSearch("__system__cmkid", image, nil)),
		d.Set("__account_code", utils.PathSearch("__account_code", image, nil)),
		d.Set("__support_amd", utils.PathSearch("__support_amd", image, nil)),
		d.Set("__support_kvm_hi1822_hisriov", utils.PathSearch("__support_kvm_hi1822_hisriov", image, nil)),
		d.Set("__support_kvm_hi1822_hivirtionet", utils.PathSearch("__support_kvm_hi1822_hivirtionet", image, nil)),
		d.Set("os_shutdown_timeout", utils.PathSearch("os_shutdown_timeout", image, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCbrWholeImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.Errorf("error updating IMS CBR whole image name field: %s", err)
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
				return diag.Errorf("error updating IMS CBR whole image description field: %s", err)
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
				return diag.Errorf("error updating IMS CBR whole image max_ram field: %s", err)
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
			return diag.Errorf("error updating IMS CBR whole image min_ram field: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateIMSImageTags(client, d)
		if err != nil {
			return diag.Errorf("error updating IMS CBR whole image tags field: %s", err)
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

	return resourceCbrWholeImageRead(ctx, d, meta)
}

func resourceCbrWholeImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/images/{image_id}"
		imageId = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	// Before deleting, call the query API first, if the query result is empty, then process `CheckDeleted` logic.
	image, err := getCbrWholeImage(client, imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	if image == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS CBR whole image")
	}

	// For the whole image, need to use `delete_backup` to control whether to delete backup when deleting image.
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{image_id}", imageId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"delete_backup": d.Get("is_delete_backup").(bool),
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting IMS CBR whole image: %s", err)
	}

	// Because the delete API always return `204` status code,
	// so we need to call the list query API to check if the image has been successfully deleted.
	err = waitForCbrWholeImageDeleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for IMS CBR whole image to be deleted: %s", err)
	}

	return nil
}

func waitForCbrWholeImageDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			image, err := getCbrWholeImage(client, d.Id())
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
