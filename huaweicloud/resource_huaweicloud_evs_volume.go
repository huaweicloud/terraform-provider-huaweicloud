package huaweicloud

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk/openstack/blockstorage/extensions/volumeactions"
	volumes_v2 "github.com/huaweicloud/golangsdk/openstack/blockstorage/v2/volumes"
	"github.com/huaweicloud/golangsdk/openstack/evs/v3/volumes"
)

func ResourceEvsStorageVolumeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceEvsVolumeV3Create,
		Read:   resourceEvsVolumeV3Read,
		Update: resourceEvsVolumeV3Update,
		Delete: resourceBlockStorageVolumeV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SATA", "SAS", "SSD", "co-p1", "uh-l1",
				}, true),
			},
			"device_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "VBD",
				ValidateFunc: validation.StringInSlice([]string{"VBD", "SCSI"}, true),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attachment": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceVolumeAttachmentHash,
			},
			"multiattach": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"kms_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"wwn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cascade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceVolumeAttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}

func resourceEvsVolumeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud EVS storage client: %s", err)
	}
	blockStorageClient.Endpoint = blockStorageClient.ResourceBase
	if !hasFilledOpt(d, "backup_id") && !hasFilledOpt(d, "size") {
		return fmt.Errorf("Missing required argument: 'size' is required, but no definition was found.")
	}
	tags := resourceContainerTags(d)
	createOpts := &volumes.CreateOpts{
		BackupID:         d.Get("backup_id").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		Description:      d.Get("description").(string),
		Size:             d.Get("size").(int),
		Name:             d.Get("name").(string),
		SnapshotID:       d.Get("snapshot_id").(string),
		ImageRef:         d.Get("image_id").(string),
		VolumeType:       d.Get("volume_type").(string),
		Multiattach:      d.Get("multiattach").(bool),
		Tags:             tags,
	}
	m := make(map[string]string)
	if v, ok := d.GetOk("kms_id"); ok {
		m["__system__cmkid"] = v.(string)
		m["__system__encrypted"] = "1"
	}
	if d.Get("device_type").(string) == "SCSI" {
		m["hw:passthrough"] = "true"
	}
	if len(m) > 0 {
		createOpts.Metadata = m
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := volumes.Create(blockStorageClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud EVS volume: %s", err)
	}
	log.Printf("[INFO] Volume Job ID: %s", v.JobID)

	// Wait for the volume to become available.
	log.Printf("[DEBUG] Waiting for volume to become available")
	err = volumes.WaitForJobSuccess(blockStorageClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return err
	}

	entity, err := volumes.GetJobEntity(blockStorageClient, v.JobID, "volume_id")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		log.Printf("[INFO] Volume ID: %s", id)
		// Store the ID now
		d.SetId(id)
		return resourceEvsVolumeV3Read(d, meta)
	}
	return fmt.Errorf("Unexpected conversion error in resourceEvsVolumeV3Create.")
}

func resourceEvsVolumeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud EVS storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "volume")
	}

	log.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), v)

	d.Set("size", v.Size)
	d.Set("description", v.Description)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("name", v.Name)
	d.Set("snapshot_id", v.SnapshotID)
	d.Set("source_vol_id", v.SourceVolID)
	d.Set("volume_type", v.VolumeType)
	d.Set("wwn", v.WWN)

	// set tags
	tags := make(map[string]string)
	for key, val := range v.Tags {
		tags[key] = val
	}
	if err := d.Set("tags", tags); err != nil {
		return fmt.Errorf("[DEBUG] Error saving tags to state for HuaweiCloud evs storage (%s): %s", d.Id(), err)
	}

	// set attachments
	attachments := make([]map[string]interface{}, len(v.Attachments))
	for i, attachment := range v.Attachments {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
		log.Printf("[DEBUG] attachment: %v", attachment)
	}
	if err := d.Set("attachment", attachments); err != nil {
		return fmt.Errorf("[DEBUG] Error saving attachment to state for HuaweiCloud evs storage (%s): %s", d.Id(), err)
	}

	return nil
}

// using OpenStack Cinder API v2 to update volume resource
func resourceEvsVolumeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud block storage client: %s", err)
	}

	updateOpts := volumes_v2.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if d.HasChange("tags") {
		_, err = resourceEVSTagV2Create(d, meta, "volumes", d.Id(), resourceContainerTags(d))
	}

	if d.HasChange("size") {
		extendOpts := volumeactions.ExtendSizeOpts{
			NewSize: d.Get("size").(int),
		}

		err = volumeactions.ExtendSize(blockStorageClient, d.Id(), extendOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error extending huaweicloud_evs_volume %s size: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err := stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for huaweicloud_evs_volume %s to become ready: %s", d.Id(), err)
		}
	}

	_, err = volumes_v2.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud volume: %s", err)
	}

	return resourceEvsVolumeV3Read(d, meta)
}
