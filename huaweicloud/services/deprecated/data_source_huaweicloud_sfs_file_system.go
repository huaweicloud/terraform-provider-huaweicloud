package deprecated

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/sfs/v2/shares"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SFS GET /v2/{project_id}/shares/detail
// @API SFS POST /v2/{project_id}/shares/{id}/action
// @API SFS GET /v2/{project_id}/shares/{id}/export_locations
func DataSourceSFSFileSystemV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSFSFileSystemV2Read,

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

func dataSourceSFSFileSystemV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS Client: %s", err)
	}

	listOpts := shares.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	refinedSfs, err := shares.List(sfsClient, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve shares: %s", err)
	}

	if len(refinedSfs) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedSfs) > 1 {
		return diag.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	share := refinedSfs[0]

	log.Printf("[INFO] Retrieved Shares using given filter %s: %+v", share.ID, share)
	d.SetId(share.ID)

	mErr := multierror.Append(nil,
		d.Set("availability_zone", share.AvailabilityZone),
		d.Set("description", share.Description),
		d.Set("is_public", share.IsPublic),
		d.Set("name", share.Name),
		d.Set("share_proto", share.ShareProto),
		d.Set("size", share.Size),
		d.Set("status", share.Status),
		d.Set("export_location", share.ExportLocation),
		d.Set("metadata", share.Metadata),
		d.Set("region", region),
	)

	n, err := shares.ListAccessRights(sfsClient, share.ID).ExtractAccessRights()
	if err != nil {
		return diag.Errorf("error extracting the AccessRight slice from the response: %s", err)
	}
	shareaccess := n[0]
	mErr = multierror.Append(mErr,
		d.Set("access_type", shareaccess.AccessType),
		d.Set("access_to", shareaccess.AccessTo),
		d.Set("access_level", shareaccess.AccessLevel),
		d.Set("state", shareaccess.State),
		d.Set("share_access_id", shareaccess.ID),
	)

	mount, err := shares.GetExportLocations(sfsClient, share.ID).ExtractExportLocations()
	if err != nil {
		return diag.Errorf("error getting the Mount/Export Locations of the SFS specified: %S", err)
	}
	mountTarget := mount[0]

	mErr = multierror.Append(mErr,
		d.Set("mount_id", mountTarget.ID),
		d.Set("preferred", mountTarget.Preferred),
		d.Set("share_instance_id", mountTarget.ShareInstanceID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
