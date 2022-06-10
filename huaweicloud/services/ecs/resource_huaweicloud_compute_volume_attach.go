package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/jobs"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceComputeVolumeAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeVolumeAttachCreate,
		ReadContext:   resourceComputeVolumeAttachRead,
		DeleteContext: resourceComputeVolumeAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceComputeVolumeAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	computeClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating compute v1 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	volumeId := d.Get("volume_id").(string)
	config.MutexKV.Lock(volumeId)
	defer config.MutexKV.Unlock(volumeId)

	var device string
	if v, ok := d.GetOk("device"); ok {
		device = v.(string)
	}

	attachOpts := block_devices.AttachOpts{
		Device:   device,
		VolumeId: volumeId,
		ServerId: instanceId,
	}

	log.Printf("[DEBUG] Creating volume attachment: %#v", attachOpts)
	job, err := block_devices.Attach(computeClient, attachOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] The response of volume attachment request is: %#v", job)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"INIT", "RUNNING"},
		Target:       []string{"SUCCESS"},
		Refresh:      AttachmentJobRefreshFunc(computeClient, job.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error attaching volume: %s", err)
	}

	id := fmt.Sprintf("%s/%s", instanceId, volumeId)
	d.SetId(id)

	return resourceComputeVolumeAttachRead(ctx, d, meta)
}

func resourceComputeVolumeAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	computeClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating compute V1 client: %s", err)
	}

	instanceId, volumeId, err := ParseComputeVolumeAttachmentId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := block_devices.Get(computeClient, instanceId, volumeId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, parseRequestError(err), "error fetching compute_volume_attach")
	}

	log.Printf("[DEBUG] Retrieved volume attachment: %#v", attachment)

	mErr := multierror.Append(nil,
		d.Set("instance_id", attachment.ServerId),
		d.Set("volume_id", attachment.VolumeId),
		d.Set("device", attachment.Device),
		d.Set("region", region),
		d.Set("pci_address", attachment.PciAddress),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceComputeVolumeAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	computeClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating compute V1 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	volumeId := d.Get("volume_id").(string)
	config.MutexKV.Lock(volumeId)
	defer config.MutexKV.Unlock(volumeId)

	opts := block_devices.DetachOpts{
		ServerId: instanceId,
	}
	job, err := block_devices.Detach(computeClient, volumeId, opts)
	if err != nil {
		return diag.FromErr(err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCESS", "NOTFOUND"},
		Refresh:      AttachmentJobRefreshFunc(computeClient, job.ID),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error detaching volume: %s", err)
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

func parseRequestError(respErr error) error {
	var apiErr block_devices.ErrorResponse
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok && errCode.Body != nil {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr == nil && apiErr.Error.Code == "Ecs.1000" && strings.Contains(apiErr.Error.Message, "itemNotFound") {
			return golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte("the volume has been deleted"),
				},
			}
		}
	}
	return respErr
}

func ParseComputeVolumeAttachmentId(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("unable to determine volume attachment ID")
	}

	instanceId := idParts[0]
	volumeId := idParts[1]

	return instanceId, volumeId, nil
}
