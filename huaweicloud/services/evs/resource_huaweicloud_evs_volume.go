package evs

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/jobs"
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
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenew(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
			"tags":          common.TagsSchema(),
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

func buildBssParamParams(d *schema.ResourceData) *cloudvolumes.BssParam {
	bssParams := &cloudvolumes.BssParam{
		ChargingMode: d.Get("charging_mode").(string),
		PeriodType:   d.Get("period_unit").(string),
		PeriodNum:    d.Get("period").(int),
		IsAutoRenew:  d.Get("auto_renew").(string),
		IsAutoPay:    common.GetAutoPay(d),
	}
	return bssParams
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
	m := map[string]string{
		"create_for_volume_id": "true",
	}
	if v, ok := d.GetOk("kms_id"); ok {
		m["__system__cmkid"] = v.(string)
		m["__system__encrypted"] = "1"
	}
	if d.Get("device_type").(string) == "SCSI" {
		m["hw:passthrough"] = "true"
	}
	volumeOpts.Metadata = m

	result := cloudvolumes.CreateOpts{
		Volume: volumeOpts,
	}
	if v, ok := d.GetOk("charging_mode"); ok && v == "prePaid" {
		result.ChargeInfo = buildBssParamParams(d)
	}
	return result
}

func resourceEvsVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	// The v2 client is used to obtain the volume detail.
	evsV2Client, err := config.BlockStorageV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud block storage v2 client: %s", err)
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
	if len(job.VolumeIDs) < 1 {
		return fmtp.DiagErrorf("The volume ID was not included in the response to the request to create the volume.")
	}
	d.SetId(job.VolumeIDs[0])

	// If charging mode is PrePaid, wait for the order to be completed.
	if job.OrderID != "" {
		err = common.WaitOrderComplete(ctx, d, config, job.OrderID)
		if err != nil {
			return fmtp.DiagErrorf("The order is not completed while creating EVS volume (%s): %#v", d.Id(), err)
		}
	}

	logp.Printf("[DEBUG] Waiting for the EVS volume to become available, the volume ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    CloudVolumeRefreshFunc(evsV2Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error waiting for the creation of EVS volume (%s) to complete: %s", d.Id(), err)
	}

	return resourceEvsVolumeRead(ctx, d, meta)
}

func setEvsVolumeDeviceType(d *schema.ResourceData, resp *cloudvolumes.Volume) error {
	if resp.Metadata.HwPassthrough == "true" {
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
	}
	logp.Printf("[DEBUG] The relevant attach information for EVS volume is: %v", attachments)
	return d.Set("attachment", attachments)
}

func setEvsVolumeChargingInfo(d *schema.ResourceData, resp *cloudvolumes.Volume) error {
	if resp.Metadata.OrderID != "" {
		return d.Set("charging_mode", "prePaid")
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
		setEvsVolumeChargingInfo(d, resp),
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
		desc := d.Get("description").(string)
		updateOpts := cloudvolumes.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: &desc,
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
			return fmtp.DiagErrorf("Error creating HuaweiCloud block storage v2.1 client: %s", err)
		}
		extendOpts := cloudvolumes.ExtendOpts{
			SizeOpts: cloudvolumes.ExtendSizeOpts{
				NewSize: d.Get("size").(int),
			},
		}

		// If charging mode is PrePaid, the order is automatically paid to adjust the volume size.
		if strings.EqualFold(d.Get("charging_mode").(string), "prePaid") {
			extendOpts.ChargeInfo = &cloudvolumes.ExtendChargeOpts{
				IsAutoPay: common.GetAutoPay(d),
			}
		}

		resp, err := cloudvolumes.ExtendSize(evsV21Client, d.Id(), extendOpts).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Error extending EVS volume (%s) size: %s", d.Id(), err)
		}

		if strings.EqualFold(d.Get("charging_mode").(string), "prePaid") {
			err = common.WaitOrderComplete(ctx, d, config, resp.OrderID)
			if err != nil {
				return fmtp.DiagErrorf("The order (%s) is not completed while extending EVS volume (%s) size: %#v",
					resp.OrderID, d.Id(), err)
			}
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    CloudVolumeRefreshFunc(evsV2Client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmtp.DiagErrorf("Error waiting for EVS volume (%s) to become ready: %s", d.Id(), err)
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
		computeClient, err := config.ComputeV1Client(region)
		if err != nil {
			return fmtp.DiagErrorf("Error creating HuaweiCloud ECS v2 client: %s", err)
		}
		for _, attachment := range v.Attachments {
			logp.Printf("[DEBUG] The attachment is: %v", attachment)
			opts := block_devices.DetachOpts{
				ServerId: attachment.ServerID,
			}
			job, err := block_devices.Detach(computeClient, attachment.AttachmentID, opts)
			if err != nil {
				return diag.FromErr(err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{"RUNNING"},
				Target:     []string{"SUCCESS", "NOTFOUND"},
				Refresh:    AttachmentJobRefreshFunc(computeClient, job.ID),
				Timeout:    10 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			if _, err = stateConf.WaitForState(); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		err = common.UnsubscribePrePaidResource(d, config, []string{d.Id()})
		if err != nil {
			return fmtp.DiagErrorf("error unsubscribing HuaweiCloud EVS volume : %s", err)
		}
	} else {
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
	}

	// Wait for the volume to delete before moving on.
	logp.Printf("[DEBUG] Waiting for the EVS volume (%s) to delete", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "downloading", "available"},
		Target:     []string{"deleted"},
		Refresh:    CloudVolumeRefreshFunc(evsV2Client, d.Id()),
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

func CloudVolumeRefreshFunc(c *golangsdk.ServiceClient, volumeId string) resource.StateRefreshFunc {
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
