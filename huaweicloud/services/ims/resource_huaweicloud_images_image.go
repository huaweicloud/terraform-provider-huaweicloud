package ims

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/backups"
	"github.com/chnsz/golangsdk/openstack/imageservice/v2/images"
	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/chnsz/golangsdk/openstack/ims/v2/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceImsImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsImageCreate,
		ReadContext:   resourceImsImageRead,
		UpdateContext: resourceImsImageUpdate,
		DeleteContext: resourceImsImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImsImageImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"max_ram": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			// instance_id is required for creating an image from an ECS
			"instance_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"image_url", "backup_id"},
			},
			"vault_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"instance_id"},
			},
			// backup_id is required for creating an image from backup of ECS
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// image_url and min_disk are required for creating an image from an OBS
			"image_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"min_disk"},
			},
			"min_disk": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id"},
			},
			// following are valid for creating an image from an OBS
			"os_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			// following are additional attributes
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceContainerImageTags(d *schema.ResourceData) []cloudimages.ImageTag {
	var tagList []cloudimages.ImageTag

	rawTags := d.Get("tags").(map[string]interface{})
	for key, val := range rawTags {
		tagRequest := cloudimages.ImageTag{
			Key:   key,
			Value: val.(string),
		}
		tagList = append(tagList, tagRequest)
	}
	return tagList
}

func resourceImsImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	var v *cloudimages.JobResponse
	imageTags := resourceContainerImageTags(d)
	if val, ok := d.GetOk("instance_id"); ok {
		v, err = createByInstanceId(d, cfg, imsClient, val.(string), imageTags)
	} else if val, ok = d.GetOk("image_url"); ok {
		v, err = createByImageUrl(d, cfg, imsClient, val.(string), imageTags)
	} else {
		v, err = createByBackupId(d, cfg, imageTags)
	}

	if err != nil {
		return diag.Errorf("error creating IMS image: %s", err)
	}
	log.Printf("[INFO] IMS Job ID: %s", v.JobID)

	// Wait for the image to become available.
	log.Printf("[DEBUG] Waiting for IMS image to become available")
	err = cloudimages.WaitForJobSuccess(imsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	entity, err := cloudimages.GetJobEntity(imsClient, v.JobID, "image_id")
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := entity.(string); ok {
		log.Printf("[INFO] IMS ID: %s", id)
		// Store the ID now
		d.SetId(id)
		return resourceImsImageRead(ctx, d, meta)
	}
	return diag.Errorf("unexpected conversion error in resourceImsImageCreate.")
}

func createByInstanceId(d *schema.ResourceData, cfg *config.Config, client *golangsdk.ServiceClient,
	instanceId string, imageTags []cloudimages.ImageTag) (*cloudimages.JobResponse, error) {
	region := cfg.GetRegion(d)
	// if vault_id is not empty, then a whole image wil be created
	if vaultId, ok := d.GetOk("vault_id"); ok {
		imsClient, err := cfg.ImageV1Client(region)
		if err != nil {
			return nil, fmt.Errorf("error creating IMS client: %s", err)
		}
		createOpts := &cloudimages.CreateWholeImageOpts{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			MaxRam:              d.Get("max_ram").(int),
			MinRam:              d.Get("min_ram").(int),
			InstanceId:          instanceId,
			ImageTags:           imageTags,
			EnterpriseProjectID: common.GetEnterpriseProjectID(d, cfg),
			VaultId:             vaultId.(string),
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		return cloudimages.CreateWholeImageByServer(imsClient, createOpts).ExtractJobResponse()
	}
	createOpts := &cloudimages.CreateByServerOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		MaxRam:              d.Get("max_ram").(int),
		MinRam:              d.Get("min_ram").(int),
		InstanceId:          instanceId,
		ImageTags:           imageTags,
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, cfg),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	return cloudimages.CreateImageByServer(client, createOpts).ExtractJobResponse()
}

func createByImageUrl(d *schema.ResourceData, cfg *config.Config, client *golangsdk.ServiceClient,
	imageUrl string, imageTags []cloudimages.ImageTag) (*cloudimages.JobResponse, error) {
	createOpts := &cloudimages.CreateByOBSOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		ImageUrl:            imageUrl,
		MinDisk:             d.Get("min_disk").(int),
		MaxRam:              d.Get("max_ram").(int),
		MinRam:              d.Get("min_ram").(int),
		OsVersion:           d.Get("os_version").(string),
		IsConfig:            d.Get("is_config").(bool),
		CmkId:               d.Get("cmk_id").(string),
		Type:                d.Get("type").(string),
		ImageTags:           imageTags,
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, cfg),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	return cloudimages.CreateImageByOBS(client, createOpts).ExtractJobResponse()
}

func createByBackupId(d *schema.ResourceData, cfg *config.Config,
	imageTags []cloudimages.ImageTag) (*cloudimages.JobResponse, error) {
	imsClient, err := cfg.ImageV1Client(cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating IMS client: %s", err)
	}
	createOpts := &cloudimages.CreateWholeImageOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		MaxRam:              d.Get("max_ram").(int),
		MinRam:              d.Get("min_ram").(int),
		BackupId:            d.Get("backup_id").(string),
		ImageTags:           imageTags,
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, cfg),
		WholeImageType:      "CBR",
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	return cloudimages.CreateWholeImageByServer(imsClient, createOpts).ExtractJobResponse()
}

func GetCloudImage(client *golangsdk.ServiceClient, id string) (*cloudimages.Image, error) {
	listOpts := &cloudimages.ListOpts{
		ID:    id,
		Limit: 1,
	}
	allPages, err := cloudimages.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	img := allImages[0]
	if img.ID != id {
		return nil, fmt.Errorf("unexpected images ID")
	}
	log.Printf("[DEBUG] Retrieved Image %s: %#v", id, img)
	return &img, nil
}

func getInstanceID(data string) string {
	results := strings.Split(data, ",")
	if len(results) == 2 && results[0] == "instance" {
		return results[1]
	}

	return ""
}

func resourceImsImageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	img, err := GetCloudImage(imsClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving image")
	}
	log.Printf("[DEBUG] Retrieved Image %s: %#v", d.Id(), img)

	mErr := multierror.Append(
		d.Set("name", img.Name),
		d.Set("description", img.Description),
		d.Set("visibility", img.Visibility),
		d.Set("disk_format", img.DiskFormat),
		d.Set("image_size", img.ImageSize),
		d.Set("enterprise_project_id", img.EnterpriseProjectID),
		d.Set("checksum", img.Checksum),
		d.Set("status", img.Status),
	)

	if img.OsVersion != "" {
		mErr = multierror.Append(mErr, d.Set("os_version", img.OsVersion))
	}
	if img.WholeImage == "true" {
		// the server will create a CBR backup first when create a whole image by an ECS instance,
		// we can only get the backup_id from the image, and the value of param data_origin only
		// contains backup_id, so we should get the instance_id by backup_id if needed
		if _, ok := d.GetOk("backup_id"); ok {
			mErr = multierror.Append(
				mErr,
				d.Set("backup_id", img.BackupID),
			)
		} else {
			cbrClient, err := cfg.CbrV3Client(region)
			if err != nil {
				return diag.Errorf("error creating CBR v3 client: %s", err)
			}
			backup, err := backups.Get(cbrClient, img.BackupID)
			if err != nil {
				return diag.Errorf("error querying backup detail: %s", err)
			}
			mErr = multierror.Append(
				mErr,
				d.Set("instance_id", backup.ResourceId),
			)
		}
		mErr = multierror.Append(
			mErr,
			d.Set("data_origin", img.DataOrigin),
		)
	}
	if img.DataOrigin != "" && img.WholeImage != "true" {
		mErr = multierror.Append(
			mErr,
			d.Set("instance_id", getInstanceID(img.DataOrigin)),
			d.Set("data_origin", img.DataOrigin),
		)
	}

	// Set image tags
	if tagList, err := tags.Get(imsClient, d.Id()).Extract(); err == nil {
		tagMap := make(map[string]string)
		for _, val := range tagList.Tags {
			tagMap[val.Key] = val.Value
		}
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] fetching tags of image failed: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func setTagForImage(d *schema.ResourceData, meta interface{}, imageID string, tagMap map[string]interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating IMS client: %s", err)
	}

	var tagList []tags.Tag
	for k, v := range tagMap {
		tag := tags.Tag{
			Key:   k,
			Value: v.(string),
		}
		tagList = append(tagList, tag)
	}

	createOpts := tags.BatchOpts{Action: tags.ActionCreate, Tags: tagList}
	createTags := tags.BatchAction(client, imageID, createOpts)
	if createTags.Err != nil {
		return fmt.Errorf("error creating image tags: %s", createTags.Err)
	}

	return nil
}

func resourceImsImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	if d.HasChange("name") {
		updateOpts := make(images.UpdateOpts, 0)
		v := images.ReplaceImageName{NewName: d.Get("name").(string)}
		updateOpts = append(updateOpts, v)

		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
		_, err = images.Update(imsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating image: %s", err)
		}
	}

	if d.HasChange("tags") {
		oldTags, err := tags.Get(imsClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error fetching image tags: %s", err)
		}
		if len(oldTags.Tags) > 0 {
			deleteOpts := tags.BatchOpts{Action: tags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := tags.BatchAction(imsClient, d.Id(), deleteOpts)
			if deleteTags.Err != nil {
				return diag.Errorf("error deleting image tags: %s", deleteTags.Err)
			}
		}

		if common.HasFilledOpt(d, "tags") {
			tagMap := d.Get("tags").(map[string]interface{})
			if len(tagMap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagMap)
				err = setTagForImage(d, meta, d.Id(), tagMap)
				if err != nil {
					return diag.Errorf("error updating tags of image:%s", err)
				}
			}
		}
	}

	return resourceImsImageRead(ctx, d, meta)
}

func resourceImsImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	imageClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	log.Printf("[DEBUG] Deleting Image %s", d.Id())
	if err = images.Delete(imageClient, d.Id()).Err; err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting Image")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForImageDelete(imageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting image: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForImageDelete(imageClient *golangsdk.ServiceClient, imageID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := images.Get(imageClient, imageID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted image %s", imageID)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}

func resourceImsImageImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating IMS client: %s", err)
	}

	img, err := GetCloudImage(imsClient, d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	var mErr *multierror.Error

	if img.WholeImage == "true" {
		cbrClient, err := cfg.CbrV3Client(region)
		if err != nil {
			return []*schema.ResourceData{d}, fmt.Errorf("error creating CBR v3 client: %s", err)
		}
		backup, err := backups.Get(cbrClient, img.BackupID)
		if err != nil {
			return []*schema.ResourceData{d}, fmt.Errorf("error querying backup detail: %s", err)
		}
		mErr = multierror.Append(
			mErr,
			d.Set("instance_id", backup.ResourceId),
		)
	}

	if img.DataOrigin != "" && img.WholeImage != "true" {
		mErr = multierror.Append(
			mErr,
			d.Set("instance_id", getInstanceID(img.DataOrigin)),
		)
	}

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
