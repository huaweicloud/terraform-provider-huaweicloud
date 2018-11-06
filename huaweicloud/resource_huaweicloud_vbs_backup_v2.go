package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/backups"
)

func resourceVBSBackupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVBSBackupV2Create,
		Read:   resourceVBSBackupV2Read,
		Delete: resourceVBSBackupV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateVBSBackupName,
			},
			"volume_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateVBSBackupDescription,
			},
			"container": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_metadata": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateVBSTagKey,
						},
						"value": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateVBSTagValue,
						},
					},
				},
			},
		},
	}
}

func resourceVBSBackupTagsV2(d *schema.ResourceData) []backups.Tag {
	rawTags := d.Get("tags").(*schema.Set).List()
	tags := make([]backups.Tag, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = backups.Tag{
			Key:   rawMap["key"].(string),
			Value: rawMap["value"].(string),
		}
	}
	return tags
}

func resourceVBSBackupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating huaweicloud vbs client: %s", err)
	}

	createOpts := backups.CreateOpts{
		Name:        d.Get("name").(string),
		VolumeId:    d.Get("volume_id").(string),
		SnapshotId:  d.Get("snapshot_id").(string),
		Description: d.Get("description").(string),
		Tags:        resourceVBSBackupTagsV2(d),
	}

	n, err := backups.Create(vbsClient, createOpts).ExtractJobResponse()

	if err != nil {
		return fmt.Errorf("Error creating huaweicloud VBS Backup: %s", err)
	}

	if err := backups.WaitForJobSuccess(vbsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	entity, err := backups.GetJobEntity(vbsClient, n.JobID, "backup_id")

	if id, ok := entity.(string); ok {
		d.SetId(id)
		return resourceVBSBackupV2Read(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceVBSBackupV2Create.")
}

func resourceVBSBackupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud Vbs client: %s", err)
	}

	n, err := backups.Get(vbsClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving huaweicloud VBS Backup: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	d.Set("status", n.Status)
	d.Set("availability_zone", n.AvailabilityZone)
	d.Set("snapshot_id", n.SnapshotId)
	d.Set("service_metadata", n.ServiceMetadata)
	d.Set("size", n.Size)
	d.Set("container", n.Container)
	d.Set("volume_id", n.VolumeId)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceVBSBackupV2Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud vbs: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForBackupDelete(vbsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting huaweicloud VBS Backup: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForBackupDelete(client *golangsdk.ServiceClient, backupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := backups.Get(client, backupID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return r, "deleted", nil
			}
			return nil, "available", err
		}

		if r.Status != "deleting" {
			err := backups.Delete(client, backupID).ExtractErr()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					log.Printf("[INFO] Successfully deleted huaweicloud VBS backup %s", backupID)
					return r, "deleted", nil
				}
				return r, r.Status, err
			}
		}
		return r, r.Status, nil
	}
}
