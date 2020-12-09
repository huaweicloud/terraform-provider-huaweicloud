package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/backup"
)

func dataSourceCSBSBackupV1() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceCSBSBackupV1Read,
		DeprecationMessage: "It has been deprecated.",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_record_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_trigger": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"average_speed": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vm_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_backups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_saving_ratio": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"average_speed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incremental": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"vm_metadata": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCSBSBackupV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	backupClient, err := config.CsbsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating csbs client: %s", err)
	}

	listOpts := backup.ListOpts{
		ID:           d.Get("id").(string),
		Name:         d.Get("backup_name").(string),
		Status:       d.Get("status").(string),
		ResourceName: d.Get("resource_name").(string),
		CheckpointId: d.Get("backup_record_id").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceId:   d.Get("resource_id").(string),
		PolicyId:     d.Get("policy_id").(string),
		VmIp:         d.Get("vm_ip").(string),
	}

	refinedbackups, err := backup.List(backupClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve backup: %s", err)
	}

	if len(refinedbackups) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedbackups) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	backupObject := refinedbackups[0]
	log.Printf("[INFO] Retrieved backup %s using given filter", backupObject.Id)

	d.SetId(backupObject.Id)

	d.Set("backup_record_id", backupObject.CheckpointId)
	d.Set("backup_name", backupObject.Name)
	d.Set("resource_id", backupObject.ResourceId)
	d.Set("status", backupObject.Status)
	d.Set("description", backupObject.Description)
	d.Set("resource_type", backupObject.ResourceType)
	d.Set("auto_trigger", backupObject.ExtendInfo.AutoTrigger)
	d.Set("average_speed", backupObject.ExtendInfo.AverageSpeed)
	d.Set("resource_name", backupObject.ExtendInfo.ResourceName)
	d.Set("size", backupObject.ExtendInfo.Size)
	d.Set("volume_backups", flattenCSBSVolumeBackups(&backupObject))
	d.Set("vm_metadata", flattenCSBSVMMetadata(&backupObject))

	d.Set("region", GetRegion(d, config))

	return nil
}
