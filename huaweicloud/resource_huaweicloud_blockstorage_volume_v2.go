package huaweicloud

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/blockstorage/extensions/volumeactions"
	"github.com/huaweicloud/golangsdk/openstack/blockstorage/v2/volumes"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/volumeattach"
)

func resourceBlockStorageVolumeV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockStorageVolumeV2Create,
		Read:   resourceBlockStorageVolumeV2Read,
		Update: resourceBlockStorageVolumeV2Update,
		Delete: resourceBlockStorageVolumeV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_evs_volume instead",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_vol_id": {
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
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"consistency_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_replica": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Set: resourceVolumeV2AttachmentHash,
			},
			"cascade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceBlockStorageVolumeV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud block storage client: %s", err)
	}

	createOpts := &volumes.CreateOpts{
		AvailabilityZone:   d.Get("availability_zone").(string),
		ConsistencyGroupID: d.Get("consistency_group_id").(string),
		Description:        d.Get("description").(string),
		ImageID:            d.Get("image_id").(string),
		Metadata:           resourceVolumeMetadataV2(d),
		Name:               d.Get("name").(string),
		Size:               d.Get("size").(int),
		SnapshotID:         d.Get("snapshot_id").(string),
		SourceReplica:      d.Get("source_replica").(string),
		SourceVolID:        d.Get("source_vol_id").(string),
		VolumeType:         d.Get("volume_type").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := volumes.Create(blockStorageClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud volume: %s", err)
	}
	log.Printf("[INFO] Volume ID: %s", v.ID)

	// Wait for the volume to become available.
	log.Printf(
		"[DEBUG] Waiting for volume (%s) to become available",
		v.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"downloading", "creating"},
		Target:     []string{"available"},
		Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, v.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to become ready: %s",
			v.ID, err)
	}

	// Store the ID now
	d.SetId(v.ID)

	return resourceBlockStorageVolumeV2Read(d, meta)
}

func resourceBlockStorageVolumeV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud block storage client: %s", err)
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
	d.Set("region", GetRegion(d, config))

	// NOTE: This tries to remove system metadata on huawei cloud :(
	md := make(map[string]string)
	var sys_keys = [3]string{"billing", "resourceSpecCode", "resourceType"}

OUTER:
	for key, val := range v.Metadata {
		for i := range sys_keys {
			if key == sys_keys[i] {
				continue OUTER
			}
		}
		md[key] = val
	}
	d.Set("metadata", md)

	attachments := make([]map[string]interface{}, len(v.Attachments))
	for i, attachment := range v.Attachments {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
		log.Printf("[DEBUG] attachment: %v", attachment)
	}
	d.Set("attachment", attachments)

	return nil
}

func resourceBlockStorageVolumeV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud block storage client: %s", err)
	}

	updateOpts := volumes.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if d.HasChange("metadata") {
		updateOpts.Metadata = resourceVolumeMetadataV2(d)
	}

	if d.HasChange("size") {
		extendOpts := volumeactions.ExtendSizeOpts{
			NewSize: d.Get("size").(int),
		}

		err = volumeactions.ExtendSize(blockStorageClient, d.Id(), extendOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error extending huaweicloud_blockstorage_volume_v2 %s size: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err := stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for huaweicloud_blockstorage_volume_v2 %s to become ready: %s", d.Id(), err)
		}
	}

	_, err = volumes.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud volume: %s", err)
	}

	return resourceBlockStorageVolumeV2Read(d, meta)
}

func resourceBlockStorageVolumeV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud block storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "volume")
	}

	// make sure this volume is detached from all instances before deleting
	if len(v.Attachments) > 0 {
		log.Printf("[DEBUG] detaching volumes")
		if computeClient, err := config.ComputeV2Client(GetRegion(d, config)); err != nil {
			return err
		} else {
			for _, volumeAttachment := range v.Attachments {
				log.Printf("[DEBUG] Attachment: %v", volumeAttachment)
				if err := volumeattach.Delete(computeClient, volumeAttachment.ServerID, volumeAttachment.ID).ExtractErr(); err != nil {
					return err
				}
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"in-use", "attaching", "detaching"},
				Target:     []string{"available"},
				Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
				Timeout:    10 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForState()
			if err != nil {
				return fmt.Errorf(
					"Error waiting for volume (%s) to become available: %s",
					d.Id(), err)
			}
		}
	}

	// The snapshots associated with the disk are deleted together with the EVS disk if cascade value is true
	deleteOpts := volumes.DeleteOpts{
		Cascade: d.Get("cascade").(bool),
	}
	// It's possible that this volume was used as a boot device and is currently
	// in a "deleting" state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.Status != "deleting" {
		if err := volumes.Delete(blockStorageClient, d.Id(), deleteOpts).ExtractErr(); err != nil {
			return CheckDeleted(d, err, "volume")
		}
	}

	// Wait for the volume to delete before moving on.
	log.Printf("[DEBUG] Waiting for volume (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "downloading", "available"},
		Target:     []string{"deleted"},
		Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceVolumeMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

// VolumeV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an HuaweiCloud volume.
func VolumeV2StateRefreshFunc(client *golangsdk.ServiceClient, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := volumes.Get(client, volumeID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "deleted", nil
			}
			return nil, "", err
		}

		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("There was an error creating the volume. " +
				"Please check with your cloud admin or check the Block Storage " +
				"API logs to see why this error occurred.")
		}

		return v, v.Status, nil
	}
}

func resourceVolumeV2AttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}
