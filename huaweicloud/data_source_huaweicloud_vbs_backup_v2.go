package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/backups"
)

func dataSourceVBSBackupV2() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceVBSBackupV2Read,
		DeprecationMessage: "It has been deprecated.",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"container": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVBSBackupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.VbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud vbs client: %s", err)
	}

	listOpts := backups.ListOpts{
		Id:         d.Get("id").(string),
		Name:       d.Get("name").(string),
		Status:     d.Get("status").(string),
		VolumeId:   d.Get("volume_id").(string),
		SnapshotId: d.Get("snapshot_id").(string),
	}

	refinedBackups, err := backups.List(vbsClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve backups: %s", err)
	}

	if len(refinedBackups) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedBackups) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Backup := refinedBackups[0]

	log.Printf("[INFO] Retrieved Backup using given filter %s: %+v", Backup.Id, Backup)
	d.SetId(Backup.Id)

	d.Set("name", Backup.Name)
	d.Set("description", Backup.Description)
	d.Set("status", Backup.Status)
	d.Set("availability_zone", Backup.AvailabilityZone)
	d.Set("snapshot_id", Backup.SnapshotId)
	d.Set("service_metadata", Backup.ServiceMetadata)
	d.Set("size", Backup.Size)
	d.Set("container", Backup.Container)
	d.Set("volume_id", Backup.VolumeId)
	d.Set("region", GetRegion(d, config))

	return nil
}
