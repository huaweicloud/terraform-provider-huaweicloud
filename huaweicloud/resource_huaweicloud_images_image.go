package huaweicloud

import (
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	imageservice_v2 "github.com/chnsz/golangsdk/openstack/imageservice/v2/images"
	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/chnsz/golangsdk/openstack/ims/v2/tags"
)

func ResourceImsImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceImsImageCreate,
		Read:   resourceImsImageRead,
		Update: resourceImsImageUpdate,
		Delete: resourceImsImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
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
			"tags": tagsSchema(),

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
				ExactlyOneOf: []string{"image_url"},
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
				ValidateFunc: validation.StringInSlice([]string{
					"ECS", "FusionCompute", "BMS", "Ironic",
				}, true),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			// following are additional attributus
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
	var tags []cloudimages.ImageTag

	rawTags := d.Get("tags").(map[string]interface{})
	for key, val := range rawTags {
		tagRequest := cloudimages.ImageTag{
			Key:   key,
			Value: val.(string),
		}
		tags = append(tags, tagRequest)
	}
	return tags
}

func resourceImsImageCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	imsClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud image client: %s", err)
	}

	var v *cloudimages.JobResponse
	imageTags := resourceContainerImageTags(d)
	if val, ok := d.GetOk("instance_id"); ok {
		createOpts := &cloudimages.CreateByServerOpts{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			MaxRam:              d.Get("max_ram").(int),
			MinRam:              d.Get("min_ram").(int),
			InstanceId:          val.(string),
			ImageTags:           imageTags,
			EnterpriseProjectID: GetEnterpriseProjectID(d, config),
		}
		logp.Printf("[DEBUG] Create Options: %#v", createOpts)
		v, err = cloudimages.CreateImageByServer(imsClient, createOpts).ExtractJobResponse()
	} else {
		createOpts := &cloudimages.CreateByOBSOpts{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			ImageUrl:            d.Get("image_url").(string),
			MinDisk:             d.Get("min_disk").(int),
			MaxRam:              d.Get("max_ram").(int),
			MinRam:              d.Get("min_ram").(int),
			OsVersion:           d.Get("os_version").(string),
			IsConfig:            d.Get("is_config").(bool),
			CmkId:               d.Get("cmk_id").(string),
			Type:                d.Get("type").(string),
			ImageTags:           imageTags,
			EnterpriseProjectID: GetEnterpriseProjectID(d, config),
		}
		logp.Printf("[DEBUG] Create Options: %#v", createOpts)
		v, err = cloudimages.CreateImageByOBS(imsClient, createOpts).ExtractJobResponse()
	}

	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IMS: %s", err)
	}
	logp.Printf("[INFO] IMS Job ID: %s", v.JobID)

	// Wait for the ims to become available.
	logp.Printf("[DEBUG] Waiting for IMS to become available")
	err = cloudimages.WaitForJobSuccess(imsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return err
	}

	entity, err := cloudimages.GetJobEntity(imsClient, v.JobID, "image_id")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		logp.Printf("[INFO] IMS ID: %s", id)
		// Store the ID now
		d.SetId(id)
		return resourceImsImageRead(d, meta)
	}
	return fmtp.Errorf("Unexpected conversion error in resourceImsImageCreate.")
}

func getCloudimage(client *golangsdk.ServiceClient, id string) (*cloudimages.Image, error) {
	listOpts := &cloudimages.ListOpts{
		ID:    id,
		Limit: 1,
	}
	allPages, err := cloudimages.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmtp.Errorf("Unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return nil, fmtp.Errorf("Unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return nil, fmtp.Errorf("Unable to find images %s: Maybe not existed", id)
	}

	img := allImages[0]
	if img.ID != id {
		return nil, fmtp.Errorf("Unexpected images ID")
	}
	logp.Printf("[DEBUG] Retrieved Image %s: %#v", id, img)
	return &img, nil
}

func getInstanceID(data string) string {
	results := strings.Split(data, ",")
	if len(results) == 2 && results[0] == "instance" {
		return results[1]
	}

	return ""
}

func resourceImsImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	imsClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud image client: %s", err)
	}

	img, err := getCloudimage(imsClient, d.Id())
	if err != nil {
		return fmtp.Errorf("Image %s not found: %s", d.Id(), err)
	}
	logp.Printf("[DEBUG] Retrieved Image %s: %#v", d.Id(), img)

	d.Set("name", img.Name)
	d.Set("description", img.Description)
	d.Set("visibility", img.Visibility)
	d.Set("disk_format", img.DiskFormat)
	d.Set("image_size", img.ImageSize)
	d.Set("enterprise_project_id", img.EnterpriseProjectID)
	d.Set("checksum", img.Checksum)
	d.Set("status", img.Status)

	if img.OsVersion != "" {
		d.Set("os_version", img.OsVersion)
	}
	if img.DataOrigin != "" {
		d.Set("instance_id", getInstanceID(img.DataOrigin))
		d.Set("data_origin", img.DataOrigin)
	}

	// Set image tags
	if Taglist, err := tags.Get(imsClient, d.Id()).Extract(); err == nil {
		tagmap := make(map[string]string)
		for _, val := range Taglist.Tags {
			tagmap[val.Key] = val.Value
		}
		if err := d.Set("tags", tagmap); err != nil {
			return fmtp.Errorf("[DEBUG] Error saving tags for HuaweiCloud image (%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] fetching tags of image failed: %s", err)
	}

	return nil
}

func setTagForImage(d *schema.ResourceData, meta interface{}, imageID string, tagmap map[string]interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud image client: %s", err)
	}

	taglist := []tags.Tag{}
	for k, v := range tagmap {
		tag := tags.Tag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	createOpts := tags.BatchOpts{Action: tags.ActionCreate, Tags: taglist}
	createTags := tags.BatchAction(client, imageID, createOpts)
	if createTags.Err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud image tags: %s", createTags.Err)
	}

	return nil
}

func resourceImsImageUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	imsClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud image client: %s", err)
	}

	if d.HasChange("name") {
		updateOpts := make(imageservice_v2.UpdateOpts, 0)
		v := imageservice_v2.ReplaceImageName{NewName: d.Get("name").(string)}
		updateOpts = append(updateOpts, v)

		logp.Printf("[DEBUG] Update Options: %#v", updateOpts)
		_, err = imageservice_v2.Update(imsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating image: %s", err)
		}
	}

	if d.HasChange("tags") {
		oldTags, err := tags.Get(imsClient, d.Id()).Extract()
		if err != nil {
			return fmtp.Errorf("Error fetching HuaweiCloud image tags: %s", err)
		}
		if len(oldTags.Tags) > 0 {
			deleteopts := tags.BatchOpts{Action: tags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := tags.BatchAction(imsClient, d.Id(), deleteopts)
			if deleteTags.Err != nil {
				return fmtp.Errorf("Error deleting HuaweiCloud image tags: %s", deleteTags.Err)
			}
		}

		if hasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				logp.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForImage(d, meta, d.Id(), tagmap)
				if err != nil {
					return fmtp.Errorf("Error updating HuaweiCloud tags of image:%s", err)
				}
			}
		}
	}

	return resourceImsImageRead(d, meta)
}

func resourceImsImageDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud image client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting Image %s", d.Id())
	if err := imageservice_v2.Delete(imageClient, d.Id()).Err; err != nil {
		return fmtp.Errorf("Error deleting Image: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForImageDelete(imageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting Huaweicloud Image: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForImageDelete(imageClient *golangsdk.ServiceClient, imageID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := imageservice_v2.Get(imageClient, imageID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[INFO] Successfully deleted Huaweicloud image %s", imageID)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
