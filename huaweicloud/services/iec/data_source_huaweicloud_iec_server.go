package iec

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/servers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/cloudservers
func DataSourceServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"edgecloud_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"edgecloud_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"coverage_sites": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"nics":            serverNicsSchema,
			"volume_attached": volumeAttachedSchema,
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	listOpts := &servers.ListOpts{
		Name:        d.Get("name").(string),
		Status:      d.Get("status").(string),
		EdgeCloudID: d.Get("edgecloud_id").(string),
	}

	log.Printf("[DEBUG] searching the IEC server by filter: %#v", listOpts)
	allServers, err := servers.List(iecClient, listOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to extract IEC server: %s", err)
	}

	total := len(allServers.Servers)
	if total < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if total > 1 {
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	server := allServers.Servers[0]
	log.Printf("[DEBUG] fetching the IEC server: %#v", server)

	d.SetId(server.ID)

	allNics, eip := expandServerNics(&server)
	// set volume fields
	allVolumes, sysDiskID := expandServerVolumeAttached(iecClient, &server)

	mErr := multierror.Append(nil,
		d.Set("name", server.Name),
		d.Set("status", server.Status),
		d.Set("flavor_id", server.Flavor.ID),
		d.Set("flavor_name", server.Flavor.Name),
		d.Set("image_id", server.Image.ID),
		d.Set("image_name", server.Metadata.ImageName),
		d.Set("vpc_id", server.Metadata.VpcID),
		d.Set("nics", allNics),
		d.Set("public_ip", eip),
		d.Set("volume_attached", allVolumes),
		d.Set("system_disk_id", sysDiskID),
		d.Set("edgecloud_id", server.EdgeCloudID),
		d.Set("edgecloud_name", server.EdgeCloudName),
	)

	if server.KeyName != "" {
		mErr = multierror.Append(mErr, d.Set("key_pair", server.KeyName))
	}
	if server.UserData != "" {
		mErr = multierror.Append(mErr, d.Set("user_data", server.UserData))
	}

	// set networking fields
	secGrpIDs := make([]string, len(server.SecurityGroups))
	for i, sg := range server.SecurityGroups {
		secGrpIDs[i] = sg.ID
	}
	mErr = multierror.Append(mErr, d.Set("security_groups", secGrpIDs))

	// set IEC fields
	location := server.Location
	siteInfo := fmt.Sprintf("%s/%s/%s/%s", location.Country, location.Area, location.Province, location.City)
	siteItem := map[string]interface{}{
		"site_id":   location.ID,
		"site_info": siteInfo,
		"operator":  server.Operator.Name,
	}
	mErr = multierror.Append(mErr, d.Set("coverage_sites", []map[string]interface{}{siteItem}))

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving IEC server: %s", err)
	}

	return nil
}
