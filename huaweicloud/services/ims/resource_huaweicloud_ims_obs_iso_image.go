package ims

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/imageservice/v2/images"
	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/chnsz/golangsdk/openstack/ims/v2/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IMS POST /v2/cloudimages/action
// @API IMS GET /v1/{project_id}/jobs/{job_id}
// @API IMS GET /v2/cloudimages
// @API IMS GET /v2/{project_id}/images/{image_id}/tags
// @API IMS PATCH /v2/cloudimages/{image_id}
// @API IMS POST /v2/{project_id}/images/{image_id}/tags/action
// @API IMS DELETE /v2/images/{image_id}
func ResourceObsIsoImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsIsoImageCreate,
		ReadContext:   resourceObsIsoImageRead,
		UpdateContext: resourceObsIsoImageUpdate,
		DeleteContext: resourceImageDelete,

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
			"os_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The `description` field can be left blank, so the `Computed` attribute is not used.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"os_type": {
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

func resourceObsIsoImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	imageTags := buildCreateImageTagsParam(d)
	createOpts := &cloudimages.CreateByOBSOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		OsVersion:           d.Get("os_version").(string),
		ImageUrl:            d.Get("image_url").(string),
		MinDisk:             d.Get("min_disk").(int),
		IsConfig:            d.Get("is_config").(bool),
		CmkId:               d.Get("cmk_id").(string),
		ImageTags:           imageTags,
		Type:                "IsoImage",
		MaxRam:              d.Get("max_ram").(int),
		MinRam:              d.Get("min_ram").(int),
		Architecture:        d.Get("architecture").(string),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}
	createResp, err := cloudimages.CreateImageByOBS(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating IMS OBS ISO image: %s", err)
	}

	imageId, err := waitForCreateImageCompleted(client, d, createResp.JobID)
	if err != nil {
		return diag.Errorf("error waiting for IMS OBS ISO image to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceObsIsoImageRead(ctx, d, meta)
}

func resourceObsIsoImageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	imageList, err := GetImageList(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving IMS OBS ISO images: %s", err)
	}

	// If the list API return empty, then process `CheckDeleted` logic.
	if len(imageList) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS OBS ISO image")
	}

	image := imageList[0]
	imageTags := flattenImageTags(d, client)
	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("name", image.Name),
		d.Set("image_url", flattenSpecificValueFormDataOrigin(image.DataOrigin, "file")),
		d.Set("min_disk", image.MinDisk),
		d.Set("os_version", image.OsVersion),
		d.Set("description", image.Description),
		d.Set("cmk_id", image.SystemCmkid),
		d.Set("architecture", flattenArchitecture(image.SupportArm)),
		d.Set("max_ram", flattenMaxRAM(image.MaxRam)),
		d.Set("min_ram", image.MinRam),
		d.Set("tags", imageTags),
		d.Set("enterprise_project_id", image.EnterpriseProjectID),
		d.Set("status", image.Status),
		d.Set("visibility", image.Visibility),
		d.Set("image_size", image.ImageSize),
		d.Set("os_type", image.OsType),
		d.Set("disk_format", image.DiskFormat),
		d.Set("data_origin", image.DataOrigin),
		d.Set("active_at", image.ActiveAt),
		d.Set("created_at", image.CreatedAt.Format(time.RFC3339)),
		d.Set("updated_at", image.UpdatedAt.Format(time.RFC3339)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceObsIsoImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	err = updateImage(ctx, cfg, client, d)
	if err != nil {
		return diag.Errorf("error updating IMS OBS ISO image: %s", err)
	}

	return resourceObsIsoImageRead(ctx, d, meta)
}

// This method is being used by other resources and will be deleted in the future.
func buildCreateImageTagsParam(d *schema.ResourceData) []cloudimages.ImageTag {
	var imageTags []cloudimages.ImageTag

	rawTags := d.Get("tags").(map[string]interface{})
	for key, val := range rawTags {
		imageTag := cloudimages.ImageTag{
			Key:   key,
			Value: val.(string),
		}
		imageTags = append(imageTags, imageTag)
	}

	return imageTags
}

// This method is being used by other resources and will be deleted in the future.
func waitForCreateImageCompleted(client *golangsdk.ServiceClient, d *schema.ResourceData, jobId string) (string, error) {
	err := cloudimages.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), jobId)
	if err != nil {
		return "", err
	}

	imageId, err := cloudimages.GetJobEntity(client, jobId, "image_id")
	if err != nil {
		return "", err
	}

	v, ok := imageId.(string)
	if !ok {
		return "", errors.New("an unexpected conversion error occurred with image_id")
	}

	return v, nil
}

// This method is being used by other resources and will be deleted in the future.
func GetImageList(client *golangsdk.ServiceClient, imageId string) ([]cloudimages.Image, error) {
	// If the `enterprise_project_id` is not filled, the list API will query images under all enterprise projects.
	// So there's no need to fill `enterprise_project_id` here.
	listOpts := &cloudimages.ListOpts{
		ID: imageId,
	}
	allPages, err := cloudimages.List(client, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return nil, fmt.Errorf("unable to extract images: %s", err)
	}

	return allImages, nil
}

// This method is being used by other resources and will be deleted in the future.
func flattenImageTags(d *schema.ResourceData, client *golangsdk.ServiceClient) map[string]string {
	tagList, err := tags.Get(client, d.Id()).Extract()
	if err == nil {
		tagMap := make(map[string]string)
		for _, val := range tagList.Tags {
			tagMap[val.Key] = val.Value
		}

		return tagMap
	}
	log.Printf("[WARN] failed fetch image tags: %s", err)

	return nil
}

// This method is being used by other resources and will be deleted in the future.
func convertTagMapToTags(tagMap map[string]interface{}) []tags.Tag {
	var tagList []tags.Tag

	for k, v := range tagMap {
		tag := tags.Tag{
			Key:   k,
			Value: v.(string),
		}
		tagList = append(tagList, tag)
	}

	return tagList
}

// This method is being used by other resources and will be deleted in the future.
func updateImage(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		imageId = d.Id()
		region  = cfg.GetRegion(d)
	)

	if d.HasChange("name") {
		updateOpts := make(cloudimages.UpdateOpts, 0)
		nameOpt := cloudimages.UpdateImageProperty{
			Op:    cloudimages.ReplaceOp,
			Name:  "name",
			Value: d.Get("name").(string),
		}
		updateOpts = append(updateOpts, nameOpt)
		_, err := cloudimages.Update(client, imageId, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("failed update name: %s", err)
		}
	}

	if d.HasChange("description") {
		updateOpts := make(cloudimages.UpdateOpts, 0)
		descriptionOpt := cloudimages.UpdateImageProperty{
			Op:    cloudimages.ReplaceOp,
			Name:  "__description",
			Value: d.Get("description").(string),
		}
		updateOpts = append(updateOpts, descriptionOpt)
		_, err := cloudimages.Update(client, imageId, updateOpts).Extract()
		if err != nil {
			err = dealUpdateDescriptionErr(d, client, err)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("max_ram") {
		updateOpts := make(cloudimages.UpdateOpts, 0)
		maxRAMOpt := cloudimages.UpdateImageProperty{
			Op:    cloudimages.ReplaceOp,
			Name:  "max_ram",
			Value: strconv.Itoa(d.Get("max_ram").(int)),
		}
		updateOpts = append(updateOpts, maxRAMOpt)
		_, err := cloudimages.Update(client, imageId, updateOpts).Extract()
		if err != nil {
			err = dealUpdateMaxRAMErr(d, client, err)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("min_ram") {
		updateOpts := make(cloudimages.UpdateOpts, 0)
		minRAMOpt := cloudimages.UpdateImageProperty{
			Op:    cloudimages.ReplaceOp,
			Name:  "min_ram",
			Value: d.Get("min_ram").(int),
		}
		updateOpts = append(updateOpts, minRAMOpt)
		_, err := cloudimages.Update(client, imageId, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("failed update min_ram: %s", err)
		}
	}

	if d.HasChange("tags") {
		err := updateImageTags(client, d)
		if err != nil {
			return err
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
			return err
		}
	}

	return nil
}

// This method is being used by other resources and will be deleted in the future.
func updateImageTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		oRaw, nRaw = d.GetChange("tags")
		oMap       = oRaw.(map[string]interface{})
		nMap       = nRaw.(map[string]interface{})
		imageId    = d.Id()
	)

	// Remove old tags
	if len(oMap) > 0 {
		deleteOpts := tags.BatchOpts{Action: tags.ActionDelete, Tags: convertTagMapToTags(oMap)}
		deleteTags := tags.BatchAction(client, imageId, deleteOpts)
		if deleteTags.Err != nil {
			return fmt.Errorf("faild delete old tags: %s", deleteTags.Err)
		}
	}

	// Create new tags
	if len(nMap) > 0 {
		createOpts := tags.BatchOpts{Action: tags.ActionCreate, Tags: convertTagMapToTags(nMap)}
		createTags := tags.BatchAction(client, imageId, createOpts)
		if createTags.Err != nil {
			return fmt.Errorf("faild create new tags: %s", createTags.Err)
		}
	}

	return nil
}

// This method is being used by other resources and will be deleted in the future.
func dealUpdateDescriptionErr(d *schema.ResourceData, client *golangsdk.ServiceClient, err error) error {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return jsonErr
		}
		errorCode, errorCodeErr := jmespath.Search("error.code", apiError)
		if errorCodeErr != nil {
			return fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}
		if errorCode != "IMG.0035" {
			return err
		}
		updateOpts := make(cloudimages.UpdateOpts, 0)
		description := cloudimages.UpdateImageProperty{
			Op:    cloudimages.AddOp,
			Name:  "__description",
			Value: d.Get("description").(string),
		}
		updateOpts = append(updateOpts, description)

		_, err = cloudimages.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("failed update description: %s", err)
		}
		return nil
	}

	return fmt.Errorf("failed update description field: %s", err)
}

// This method is being used by other resources and will be deleted in the future.
func dealUpdateMaxRAMErr(d *schema.ResourceData, client *golangsdk.ServiceClient, err error) error {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return jsonErr
		}

		errorCode, errorCodeErr := jmespath.Search("error.code", apiError)
		if errorCodeErr != nil {
			return fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if errorCode != "IMG.0035" {
			return err
		}

		updateOpts := make(cloudimages.UpdateOpts, 0)
		description := cloudimages.UpdateImageProperty{
			Op:    cloudimages.AddOp,
			Name:  "max_ram",
			Value: d.Get("max_ram").(int),
		}
		updateOpts = append(updateOpts, description)

		_, err = cloudimages.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("failed update max_ram: %s", err)
		}

		return nil
	}

	return fmt.Errorf("failed update max_ram field: %s", err)
}

// This method is being used by other resources and will be deleted in the future.
func resourceImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		imageId = d.Id()
	)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	// Before deleting, call the query API first, if the query result is empty, then process `CheckDeleted` logic.
	imageList, err := GetImageList(client, imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(imageList) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS image")
	}

	if err = images.Delete(client, imageId).Err; err != nil {
		return diag.Errorf("error deleting IMS image: %s", err)
	}

	// Because the delete API always return `204` status code,
	// so we need to call the list query API to check if the image has been successfully deleted.
	err = waitForDeleteImageCompleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for IMS image deleted: %s", err)
	}

	return nil
}

// This method is being used by other resources and will be deleted in the future.
func waitForDeleteImageCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			listResp, err := GetImageList(client, d.Id())
			if err != nil {
				return nil, "ERROR", nil
			}

			if len(listResp) < 1 {
				return "success", "COMPLETED", nil
			}

			return listResp, "PENDING", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
