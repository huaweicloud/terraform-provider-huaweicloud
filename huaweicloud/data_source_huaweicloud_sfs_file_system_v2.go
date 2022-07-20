package huaweicloud

import (
	"github.com/chnsz/golangsdk/openstack/sfs/v2/shares"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSFSFileSystemV2() *schema.Resource {
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
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"share_proto": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_level": {
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
			"share_access_id": {
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
	config := meta.(*config.Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud SFS Client: %s", err)
	}

	listOpts := shares.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	refinedSfs, err := shares.List(sfsClient, listOpts)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve shares: %s", err)
	}

	if len(refinedSfs) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedSfs) > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	share := refinedSfs[0]

	logp.Printf("[INFO] Retrieved Shares using given filter %s: %+v", share.ID, share)
	d.SetId(share.ID)

	d.Set("availability_zone", share.AvailabilityZone)
	d.Set("description", share.Description)
	d.Set("is_public", share.IsPublic)
	d.Set("name", share.Name)
	d.Set("share_proto", share.ShareProto)
	d.Set("size", share.Size)
	d.Set("status", share.Status)
	d.Set("export_location", share.ExportLocation)
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
