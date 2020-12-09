package huaweicloud

import (
	"fmt"
	"log"

	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSFSFileSystemV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSFSFileSystemV2Read,

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
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"share_proto": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"export_locations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"access_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_access_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSFSFileSystemV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud SFS Client: %s", err)
	}

	listOpts := shares.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	refinedSfs, err := shares.List(sfsClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve shares: %s", err)
	}

	if len(refinedSfs) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedSfs) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	share := refinedSfs[0]

	log.Printf("[INFO] Retrieved Shares using given filter %s: %+v", share.ID, share)
	d.SetId(share.ID)

	d.Set("availability_zone", share.AvailabilityZone)
	d.Set("description", share.Description)
	d.Set("host", share.Host)
	d.Set("id", share.ID)
	d.Set("is_public", share.IsPublic)
	d.Set("name", share.Name)
	d.Set("project_id", share.ProjectID)
	d.Set("share_proto", share.ShareProto)
	d.Set("share_type", share.ShareType)
	d.Set("size", share.Size)
	d.Set("status", share.Status)
	d.Set("volume_type", share.VolumeType)
	d.Set("export_location", share.ExportLocation)
	d.Set("export_locations", share.ExportLocations)
	d.Set("metadata", share.Metadata)
	d.Set("region", GetRegion(d, config))

	n, err := shares.ListAccessRights(sfsClient, share.ID).ExtractAccessRights()
	shareaccess := n[0]
	d.Set("access_type", shareaccess.AccessType)
	d.Set("access_to", shareaccess.AccessTo)
	d.Set("access_level", shareaccess.AccessLevel)
	d.Set("state", shareaccess.State)
	d.Set("share_access_id", shareaccess.ID)

	mount, err := shares.GetExportLocations(sfsClient, share.ID).ExtractExportLocations()
	MountTarget := mount[0]

	d.Set("mount_id", MountTarget.ID)
	d.Set("preferred", MountTarget.Preferred)
	d.Set("share_instance_id", MountTarget.ShareInstanceID)

	return nil
}
