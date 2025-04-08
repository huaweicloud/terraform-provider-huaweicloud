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

// @API IMS POST /v2/cloudimages/action
// @API IMS GET /v1/{project_id}/jobs/{job_id}
// @API IMS GET /v2/cloudimages
// @API IMS GET /v2/{project_id}/images/{image_id}/tags
// @API IMS PATCH /v2/cloudimages/{image_id}
// @API IMS POST /v2/{project_id}/images/{image_id}/tags/action
// @API IMS DELETE /v2/images/{image_id}
func ResourceObsSystemImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsSystemImageCreate,
		ReadContext:   resourceObsSystemImageRead,
		UpdateContext: resourceObsSystemImageUpdate,
		DeleteContext: resourceObsSystemImageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
			"image_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"min_disk": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			// The `description` field can be left blank, so the `Computed` attribute is not used.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_config": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"cmk_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"is_quick_import": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateObsSystemImageBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"image_url":             d.Get("image_url"),
		"min_disk":              d.Get("min_disk"),
		"description":           d.Get("description"),
		"type":                  utils.ValueIgnoreEmpty(d.Get("type")),
		"cmk_id":                utils.ValueIgnoreEmpty(d.Get("cmk_id")),
		"os_type":               utils.ValueIgnoreEmpty(d.Get("os_type")),
		"os_version":            utils.ValueIgnoreEmpty(d.Get("os_version")),
		"architecture":          utils.ValueIgnoreEmpty(d.Get("architecture")),
		"max_ram":               utils.ValueIgnoreEmpty(d.Get("max_ram")),
		"min_ram":               utils.ValueIgnoreEmpty(d.Get("min_ram")),
		"image_tags":            utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
	}

	if d.Get("is_config").(bool) {
		bodyParams["is_config"] = true
	}

	if d.Get("is_quick_import").(bool) {
		bodyParams["is_quick_import"] = true
	}

	return bodyParams
}

func resourceObsSystemImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/cloudimages/action"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateObsSystemImageBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IMS OBS system image: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating IMS OBS system image: job ID is not found in API response")
	}

	imageId, err := waitForCreateObsSystemImageJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for IMS OBS system image to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceObsSystemImageRead(ctx, d, meta)
}

func waitForCreateObsSystemImageJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      obsSystemImageJobStatusRefreshFunc(jobId, client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for IMS OBS system image job (%s) to succeed: %s", jobId, err)
	}

	return getObsSystemImageIdByJobId(client, jobId)
}

func obsSystemImageJobStatusRefreshFunc(jobId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return getResp, "ERROR", fmt.Errorf("error retrieving IMS OBS system image job: %s", err)
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
			return getRespBody, "COMPLETED", errors.New("the OBS system image creation job execution failed")
		}

		if status == "" {
			return getRespBody, "ERROR", errors.New("status field is not found in API response")
		}

		return getRespBody, "PENDING", nil
	}
}

func getObsSystemImageIdByJobId(client *golangsdk.ServiceClient, jobId string) (string, error) {
	getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return "", fmt.Errorf("error retrieving IMS OBS system image job: %s", err)
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

func getObsSystemImage(client *golangsdk.ServiceClient, imageId string) (interface{}, error) {
	// If the `enterprise_project_id` is not filled, the list API will query images under all enterprise projects.
	// So there's no need to fill `enterprise_project_id` here.
	getPath := client.Endpoint + "v2/cloudimages"
	getPath += fmt.Sprintf("?id=%s", imageId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS OBS system image: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("images[0]", getRespBody, nil), nil
}

func resourceObsSystemImageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	image, err := getObsSystemImage(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// If the list API return empty, then process `CheckDeleted` logic.
	if image == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving IMS OBS system image")
	}

	dataOrigin := utils.PathSearch("__data_origin", image, "").(string)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", image, nil)),
		d.Set("image_url", flattenSpecificValueFormDataOrigin(dataOrigin, "file")),
		d.Set("min_disk", utils.PathSearch("min_disk", image, nil)),
		d.Set("description", utils.PathSearch("__description", image, nil)),
		d.Set("os_type", utils.PathSearch("__os_type", image, nil)),
		d.Set("os_version", utils.PathSearch("__os_version", image, nil)),
		d.Set("cmk_id", utils.PathSearch("__system__cmkid", image, nil)),
		d.Set("max_ram", flattenMaxRAM(utils.PathSearch("max_ram", image, "").(string))),
		d.Set("min_ram", utils.PathSearch("min_ram", image, nil)),
		d.Set("architecture", flattenArchitecture(utils.PathSearch("__support_arm", image, "").(string))),
		d.Set("tags", flattenIMSImageTags(client, d.Id())),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", image, nil)),
		d.Set("status", utils.PathSearch("status", image, nil)),
		d.Set("visibility", utils.PathSearch("visibility", image, nil)),
		d.Set("image_size", utils.PathSearch("__image_size", image, nil)),
		d.Set("disk_format", utils.PathSearch("disk_format", image, nil)),
		d.Set("data_origin", dataOrigin),
		d.Set("active_at", utils.PathSearch("active_at", image, nil)),
		d.Set("created_at", utils.PathSearch("created_at", image, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", image, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenArchitecture(supportArm string) string {
	if supportArm == "true" {
		return "arm"
	}

	return "x86"
}

func resourceObsSystemImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.Errorf("error updating IMS OBS system image name field: %s", err)
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
				return diag.Errorf("error updating IMS OBS system image description field: %s", err)
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
				return diag.Errorf("error updating IMS OBS system image max_ram field: %s", err)
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
			return diag.Errorf("error updating IMS OBS system image min_ram field: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateIMSImageTags(client, d)
		if err != nil {
			return diag.Errorf("error updating IMS OBS system image tags field: %s", err)
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

	return resourceObsSystemImageRead(ctx, d, meta)
}

func resourceObsSystemImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	image, err := getObsSystemImage(client, imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	if image == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS OBS system image")
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{image_id}", imageId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting IMS OBS system image: %s", err)
	}

	// Because the delete API always return `204` status code,
	// so we need to call the list query API to check if the image has been successfully deleted.
	err = waitForObsSystemImageDeleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for IMS OBS system image to be deleted: %s", err)
	}

	return nil
}

func waitForObsSystemImageDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			image, err := getObsSystemImage(client, d.Id())
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
