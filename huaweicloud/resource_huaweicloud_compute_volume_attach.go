package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/volumeattach"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/block_devices"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ResourceComputeVolumeAttachV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeVolumeAttachV2Create,
		Read:   resourceComputeVolumeAttachV2Read,
		Delete: resourceComputeVolumeAttachV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"device": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: suppressDiffAll,
			},

			"pci_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeVolumeAttachV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	volumeId := d.Get("volume_id").(string)

	var device string
	if v, ok := d.GetOk("device"); ok {
		device = v.(string)
	}

	attachOpts := volumeattach.CreateOpts{
		Device:   device,
		VolumeID: volumeId,
	}

	log.Printf("[DEBUG] Creating volume attachment: %#v", attachOpts)

	attachment, err := volumeattach.Create(computeClient, instanceId, attachOpts).Extract()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    resourceComputeVolumeAttachV2AttachFunc(computeClient, instanceId, attachment.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 15 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error attaching HuaweiCloud volume: %s", err)
	}

	log.Printf("[DEBUG] Created volume attachment: %#v", attachment)

	// Use the instance ID and attachment ID as the resource ID.
	// This is because an attachment cannot be retrieved just by its ID alone.
	id := fmt.Sprintf("%s/%s", instanceId, attachment.ID)

	d.SetId(id)

	return resourceComputeVolumeAttachV2Read(d, meta)
}

func resourceComputeVolumeAttachV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instanceId, attachmentId, err := parseComputeVolumeAttachmentId(d.Id())
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	attachment, err := volumeattach.Get(computeClient, instanceId, attachmentId).Extract()
	if err != nil {
		return CheckDeleted(d, err, "compute_volume_attach")
	}

	log.Printf("[DEBUG] Retrieved volume attachment: %#v", attachment)

	d.Set("instance_id", attachment.ServerID)
	d.Set("volume_id", attachment.VolumeID)
	d.Set("device", attachment.Device)
	d.Set("region", GetRegion(d, config))

	computeV1Client, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		log.Printf("[WARN] Error creating HuaweiCloud compute V1 client: %s", err)
	} else {
		bd, err := block_devices.Get(computeV1Client, attachment.ServerID, attachment.VolumeID).Extract()
		if err != nil {
			log.Printf("[WARN] Error fetching HuaweiCloud block device: %s", err)
		}
		log.Printf("[DEBUG] Retrieved volume attachment extra info: %#v", bd)
		d.Set("pci_address", bd.PciAddress)
	}

	return nil
}

func resourceComputeVolumeAttachV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instanceId, attachmentId, err := parseComputeVolumeAttachmentId(d.Id())
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    resourceComputeVolumeAttachV2DetachFunc(computeClient, instanceId, attachmentId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      15 * time.Second,
		MinTimeout: 15 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error detaching HuaweiCloud volume: %s", err)
	}

	return nil
}

func resourceComputeVolumeAttachV2AttachFunc(
	computeClient *golangsdk.ServiceClient, instanceId, attachmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := volumeattach.Get(computeClient, instanceId, attachmentId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		return va, "ATTACHED", nil
	}
}

func resourceComputeVolumeAttachV2DetachFunc(
	computeClient *golangsdk.ServiceClient, instanceId, attachmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to detach HuaweiCloud volume %s from instance %s",
			attachmentId, instanceId)

		va, err := volumeattach.Get(computeClient, instanceId, attachmentId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = volumeattach.Delete(computeClient, instanceId, attachmentId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud Volume Attachment (%s) is still active.", attachmentId)
		return nil, "", nil
	}
}

func parseComputeVolumeAttachmentId(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine volume attachment ID")
	}

	instanceId := idParts[0]
	attachmentId := idParts[1]

	return instanceId, attachmentId, nil
}
