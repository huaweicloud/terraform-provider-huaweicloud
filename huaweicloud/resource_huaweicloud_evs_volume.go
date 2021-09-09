package huaweicloud

import (
	"bytes"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/blockstorage/extensions/volumeactions"
	volumes_v2 "github.com/chnsz/golangsdk/openstack/blockstorage/v2/volumes"
	"github.com/chnsz/golangsdk/openstack/evs/v3/volumes"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
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
			"tags": tagsSchema(),
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

func resourceEvsVolumeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud EVS storage client: %s", err)
	}
	// compatiable with job
	blockStorageClient.Endpoint = blockStorageClient.ResourceBase

	tags := resourceContainerTags(d)
	epsID := GetEnterpriseProjectID(d, config)
	createOpts := &volumes.CreateOpts{
		BackupID:            d.Get("backup_id").(string),
		AvailabilityZone:    d.Get("availability_zone").(string),
		Description:         d.Get("description").(string),
		Size:                d.Get("size").(int),
		Name:                d.Get("name").(string),
		SnapshotID:          d.Get("snapshot_id").(string),
		ImageRef:            d.Get("image_id").(string),
		VolumeType:          d.Get("volume_type").(string),
		Multiattach:         d.Get("multiattach").(bool),
		EnterpriseProjectID: epsID,
		Tags:                tags,
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

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := volumes.Create(blockStorageClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud EVS volume: %s", err)
	}
	logp.Printf("[INFO] Volume Job ID: %s", v.JobID)

	// Wait for the volume to become available.
	logp.Printf("[DEBUG] Waiting for volume to become available")
	err = volumes.WaitForJobSuccess(blockStorageClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return err
	}

	entity, err := volumes.GetJobEntity(blockStorageClient, v.JobID, "volume_id")
	if err != nil {
		return err
	}
	if _, ok := entity.(string); !ok {
		return fmtp.Errorf("error retrieving volume ID from job entity")
	}

	// Store the ID now
	d.SetId(entity.(string))
	return resourceEvsVolumeV3Read(d, meta)
}

func setEvsVolumeDeviceType(d *schema.ResourceData, resp *volumes.Volume) error {
	if value, ok := resp.Metadata["hw:passthrough"]; ok && value == "true" {
		return d.Set("device_type", "SCSI")
	}
	return d.Set("device_type", "VBD")
}

func setEvsVolumeImageId(d *schema.ResourceData, resp *volumes.Volume) error {
	if value, ok := resp.VolumeImageMetadata["image_id"]; ok {
		return d.Set("image_id", value)
	}
	return nil
}

func setEvsVolumeAttachment(d *schema.ResourceData, resp *volumes.Volume) error {
	// set attachments
	attachments := make([]map[string]interface{}, len(resp.Attachments))
	for i, attachment := range resp.Attachments {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
		logp.Printf("[DEBUG] attachment: %v", attachment)
	}
	if err := d.Set("attachment", attachments); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving attachment to state for HuaweiCloud evs storage (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceEvsVolumeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud EVS storage client: %s", err)
	}

	resp, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "volume")
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
		d.Set("region", GetRegion(d, config)),
		d.Set("wwn", resp.WWN),
		d.Set("multiattach", resp.Multiattach),
		d.Set("tags", resp.Tags),
		setEvsVolumeDeviceType(d, resp),
		setEvsVolumeImageId(d, resp),
		setEvsVolumeAttachment(d, resp),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting vault fields: %s", err)
	}

	return nil
}

// using OpenStack Cinder API v2 to update volume resource
func resourceEvsVolumeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud block storage client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updateOpts := volumes_v2.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
		_, err = volumes_v2.Update(blockStorageClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud volume: %s", err)
		}
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
			return fmtp.Errorf("Error extending huaweicloud_evs_volume %s size: %s", d.Id(), err)
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
			return fmtp.Errorf(
				"Error waiting for huaweicloud_evs_volume %s to become ready: %s", d.Id(), err)
		}
	}

	return resourceEvsVolumeV3Read(d, meta)
}
