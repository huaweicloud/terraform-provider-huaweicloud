package evs

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/volumeattach"
	"github.com/chnsz/golangsdk/openstack/evs/v1/jobs"
	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceEvsVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsVolumeCreate,
		ReadContext:   resourceEvsVolumeRead,
		UpdateContext: resourceEvsVolumeUpdate,
		DeleteContext: resourceEvsVolumeDelete,

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
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GPSSD", "SSD", "ESSD", "SAS",
				}, true),
			},
			"device_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "VBD",
				ValidateFunc: validation.StringInSlice([]string{"VBD", "SCSI"}, true),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"size", "backup_id"},
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kms_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"multiattach": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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

func buildEvsVolumeCreateOpts(d *schema.ResourceData, config *config.Config) cloudvolumes.CreateOpts {
	volumeOpts := cloudvolumes.VolumeOpts{
		AvailabilityZone:    d.Get("availability_zone").(string),
		VolumeType:          d.Get("volume_type").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Size:                d.Get("size").(int),
		BackupID:            d.Get("backup_id").(string),
		SnapshotID:          d.Get("snapshot_id").(string),
		ImageID:             d.Get("image_id").(string),
		Multiattach:         d.Get("multiattach").(bool),
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, config),
		Tags:                resourceContainerTags(d),
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
		volumeOpts.Metadata = m
	}

	return cloudvolumes.CreateOpts{
		Volume: volumeOpts,
	}
}

func resourceEvsVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	// The v1 client is used to get the volume jobs.
	evsV1Client, err := config.BlockStorageV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud block storage v1 client: %s", err)
	}
	// The v2.1 client is used to create the volume.
	evsV21Client, err := config.BlockStorageV21Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud block storage v2.1 client: %s", err)
	}

	opt := buildEvsVolumeCreateOpts(d, config)
	logp.Printf("[DEBUG] Create Options: %#v", opt)
	job, err := cloudvolumes.Create(evsV21Client, opt).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud EVS volume: %s", err)
	}
	logp.Printf("[INFO] Volume Job ID: %s", job.JobID)

	logp.Printf("[DEBUG] Waiting for the EVS volume to become available, the job ID is %s.", job.JobID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"SUCCESS"},
		Refresh:    jobRefreshFunc(evsV1Client, job.JobID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error waiting for the creation of EVS volume (%s) to complete: %s", d.Id(), err)
	}

	d.SetId(resp.(*jobs.Job).Entities.VolumeID)
	return resourceEvsVolumeRead(ctx, d, meta)
}

func setEvsVolumeDeviceType(d *schema.ResourceData, resp *cloudvolumes.Volume) error {
	if value, ok := resp.ImageMetadata["hw:passthrough"]; ok && value == "true" {
		return d.Set("device_type", "SCSI")
	}
	return d.Set("device_type", "VBD")
}

func setEvsVolumeImageId(d *schema.ResourceData, resp *cloudvolumes.Volume) error {
	if value, ok := resp.ImageMetadata["image_id"]; ok {
		return d.Set("image_id", value)
	}
	return nil
}

func setEvsVolumeAttachment(d *schema.ResourceData, resp *cloudvolumes.Volume) error {
	attachments := make([]map[string]interface{}, len(resp.Attachments))
	for i, attachment := range resp.Attachments {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.AttachmentID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
		logp.Printf("[DEBUG] attachment: %v", attachment)
	}
	if err := d.Set("attachment", attachments); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving attachment to state for HuaweiCloud EVS volume (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceEvsVolumeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	evsV2Client, err := config.BlockStorageV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud block storage v2 client: %s", err)
	}

	resp, err := cloudvolumes.Get(evsV2Client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "EVS volume")
	}

	logp.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), resp)
	mErr := multierror.Append(
		d.Set("name", resp.Name),
		d.Set("size", resp.Size),
		d.Set("description", resp.Description),
		d.Set("availability_zone", resp.AvailabilityZone),
		d.Set("snapshot_id", resp.SnapshotID),
		d.Set("volume_type", resp.VolumeType),
		d.Set("enterprise_project_id", resp.EnterpriseProjectID),
		d.Set("region", config.GetRegion(d)),
		d.Set("wwn", resp.WWN),
		d.Set("multiattach", resp.Multiattach),
		d.Set("tags", resp.Tags),
		setEvsVolumeDeviceType(d, resp),
		setEvsVolumeImageId(d, resp),
		setEvsVolumeAttachment(d, resp),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting volume fields: %s", err)
	}

	return nil
}

func resourceEvsVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	evsV2Client, err := config.BlockStorageV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud block storage v2 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updateOpts := cloudvolumes.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: golangsdk.MaybeString(d.Get("description").(string)),
		}
		_, err = cloudvolumes.Update(evsV2Client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloud volume: %s", err)
		}
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(evsV2Client, d, "cloudvolumes", d.Id())
		if tagErr != nil {
			return fmtp.DiagErrorf("Error updating tags of HuaweiCloud volume:%s, err:%s", d.Id(), tagErr)
		}
	}

	if d.HasChange("size") {
		evsV21Client, err := config.BlockStorageV21Client(config.GetRegion(d))
		if err != nil {
			return fmtp.DiagErrorf("Error creating HuaweiCloud block storage client: %s", err)
		}
		extendOpts := cloudvolumes.ExtendOpts{
			SizeOpts: cloudvolumes.ExtendSizeOpts{
				NewSize: d.Get("size").(int),
			},
		}

		_, err = cloudvolumes.ExtendSize(evsV21Client, d.Id(), extendOpts).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Error extending EVS volume (%s) size: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    cloudVolumeRefreshFunc(evsV2Client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmtp.DiagErrorf("Error waiting for huaweicloud_evs_volume %s to become ready: %s", d.Id(), err)
		}
	}

	return resourceEvsVolumeRead(ctx, d, meta)
}

func resourceContainerTags(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("tags").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceEvsVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	evsV2Client, err := config.BlockStorageV2Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud block storage v2 client: %s", err)
	}

	v, err := cloudvolumes.Get(evsV2Client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "EVS volume")
	}

	// Make sure this volume is detached from all instances before deleting.
	if len(v.Attachments) > 0 {
		logp.Printf("[DEBUG] Start to detaching volumes.")
		computeClient, err := config.ComputeV2Client(region)
		if err != nil {
			return fmtp.DiagErrorf("Error creating HuaweiCloud ECS v2 client: %s", err)
		}
		for _, attachment := range v.Attachments {
			logp.Printf("[DEBUG] The attachment is: %v", attachment)
			err = volumeattach.Delete(computeClient, attachment.ServerID, attachment.AttachmentID).ExtractErr()
			if err != nil {
				return fmtp.DiagErrorf("An error occurred while detaching the volume from the object instance: %s", err)
			}
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"in-use", "attaching", "detaching"},
			Target:     []string{"available"},
			Refresh:    cloudVolumeRefreshFunc(evsV2Client, d.Id()),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmtp.DiagErrorf("Error waiting for volume (%s) to become available: %s", d.Id(), err)
		}
	}

	// The snapshots associated with the disk are deleted together with the EVS disk if cascade value is true
	deleteOpts := cloudvolumes.DeleteOpts{
		Cascade: d.Get("cascade").(bool),
	}
	// It's possible that this volume was used as a boot device and is currently
	// in a "deleting" state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.Status != "deleting" {
		if err := cloudvolumes.Delete(evsV2Client, d.Id(), deleteOpts).ExtractErr(); err != nil {
			return common.CheckDeletedDiag(d, err, "EVS volume")
		}
	}

	// Wait for the volume to delete before moving on.
	logp.Printf("[DEBUG] Waiting for the EVS volume (%s) to delete", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "downloading", "available"},
		Target:     []string{"deleted"},
		Refresh:    cloudVolumeRefreshFunc(evsV2Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.DiagErrorf("Error waiting for the EVS volume (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func jobRefreshFunc(c *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := jobs.GetJobDetails(c, jobId).ExtractJob()
		if err != nil {
			return response, "ERROR", err
		}
		if response != nil {
			return response, response.Status, nil
		}
		return response, "ERROR", nil
	}
}

func cloudVolumeRefreshFunc(c *golangsdk.ServiceClient, volumeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := cloudvolumes.Get(c, volumeId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return response, "deleted", nil
			}
			return response, "ERROR", err
		}
		if response != nil {
			return response, response.Status, nil
		}
		return response, "ERROR", nil
	}
}
