package huaweicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/jobs"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceComputeVolumeAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeVolumeAttachCreate,
		Read:   resourceComputeVolumeAttachRead,
		Delete: resourceComputeVolumeAttachDelete,
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
				DiffSuppressFunc: utils.SuppressDiffAll,
			},

			"pci_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeVolumeAttachCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute v1 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	volumeId := d.Get("volume_id").(string)

	var device string
	if v, ok := d.GetOk("device"); ok {
		device = v.(string)
	}

	attachOpts := block_devices.AttachOpts{
		Device:   device,
		VolumeId: volumeId,
		ServerId: instanceId,
	}

	logp.Printf("[DEBUG] Creating volume attachment: %#v", attachOpts)
	job, err := block_devices.Attach(computeClient, attachOpts)
	if err != nil {
		return err
	}

	logp.Printf("[DEBUG] The response of volume attachment request is: %#v", job)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INIT", "RUNNING"},
		Target:     []string{"SUCCESS"},
		Refresh:    AttachmentJobRefreshFunc(computeClient, job.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 15 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmtp.Errorf("Error attaching HuaweiCloud volume: %s", err)
	}

	// The old logic use the instance ID and attachment ID as the resource ID, and the attachment ID equals
	// the volume ID.
	id := fmt.Sprintf("%s/%s", instanceId, volumeId)
	d.SetId(id)

	return resourceComputeVolumeAttachRead(d, meta)
}

func resourceComputeVolumeAttachRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instanceId, volumeId, err := parseComputeVolumeAttachmentId(d.Id())
	if err != nil {
		return err
	}

	attachment, err := block_devices.Get(computeClient, instanceId, volumeId).Extract()
	if err != nil {
		return CheckDeleted(d, err, "compute_volume_attach")
	}

	logp.Printf("[DEBUG] Retrieved volume attachment: %#v", attachment)

	mErr := multierror.Append(nil,
		d.Set("instance_id", attachment.ServerId),
		d.Set("volume_id", attachment.VolumeId),
		d.Set("device", attachment.Device),
		d.Set("region", region),
		d.Set("pci_address", attachment.PciAddress),
	)

	return mErr.ErrorOrNil()
}

func resourceComputeVolumeAttachDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instanceId, volumeId, err := parseComputeVolumeAttachmentId(d.Id())
	if err != nil {
		return err
	}

	opts := block_devices.DetachOpts{
		ServerId: instanceId,
	}
	job, err := block_devices.Detach(computeClient, volumeId, opts)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"RUNNING"},
		Target:     []string{"SUCCESS", "NOTFOUND"},
		Refresh:    AttachmentJobRefreshFunc(computeClient, job.ID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      15 * time.Second,
		MinTimeout: 15 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmtp.Errorf("Error detaching HuaweiCloud volume: %s", err)
	}

	return nil
}

func AttachmentJobRefreshFunc(c *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := jobs.Get(c, jobId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resp, "NOTFOUND", nil
			}
			return resp, "ERROR", err
		}

		return resp, resp.Status, nil
	}
}

func parseComputeVolumeAttachmentId(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmtp.Errorf("Unable to determine volume attachment ID")
	}

	instanceId := idParts[0]
	volumeId := idParts[1]

	return instanceId, volumeId, nil
}
