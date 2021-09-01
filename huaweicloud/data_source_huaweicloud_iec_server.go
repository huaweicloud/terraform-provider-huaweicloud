package huaweicloud

import (
	"fmt"

	"github.com/chnsz/golangsdk/openstack/iec/v1/servers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func dataSourceIECServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIECServerRead,

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

			"nics":            iecServerNicsSchema,
			"volume_attached": iecVolumeAttachedSchema,
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

func dataSourceIECServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	listOpts := &servers.ListOpts{
		Name:        d.Get("name").(string),
		Status:      d.Get("status").(string),
		EdgeCloudID: d.Get("edgecloud_id").(string),
	}

	logp.Printf("[DEBUG] searching the IEC server by filter: %#v", listOpts)
	allServers, err := servers.List(iecClient, listOpts).Extract()
	if err != nil {
		return err
	}

	total := len(allServers.Servers)
	if total < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if total > 1 {
		return fmtp.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	server := allServers.Servers[0]
	logp.Printf("[DEBUG] fetching the IEC server: %#v", server)

	d.SetId(server.ID)
	d.Set("name", server.Name)
	d.Set("status", server.Status)

	flavorInfo := server.Flavor
	d.Set("flavor_id", flavorInfo.ID)
	d.Set("flavor_name", flavorInfo.Name)
	d.Set("image_id", server.Image.ID)
	d.Set("image_name", server.Metadata.ImageName)

	if server.KeyName != "" {
		d.Set("key_pair", server.KeyName)
	}
	if server.UserData != "" {
		d.Set("user_data", server.UserData)
	}

	// set networking fields
	d.Set("vpc_id", server.Metadata.VpcID)
	secGrpIDs := make([]string, len(server.SecurityGroups))
	for i, sg := range server.SecurityGroups {
		secGrpIDs[i] = sg.ID
	}
	d.Set("security_groups", secGrpIDs)

	allNics, eip := expandIecServerNics(&server)
	d.Set("nics", allNics)
	d.Set("public_ip", eip)

	// set volume fields
	allVolumes, sysDiskID := expandIecServerVolumeAttached(iecClient, &server)
	d.Set("volume_attached", allVolumes)
	d.Set("system_disk_id", sysDiskID)

	// set IEC fields
	location := server.Location
	siteInfo := fmt.Sprintf("%s/%s/%s/%s", location.Country, location.Area, location.Province, location.City)
	siteItem := map[string]interface{}{
		"site_id":   location.ID,
		"site_info": siteInfo,
		"operator":  server.Operator.Name,
	}
	d.Set("coverage_sites", []map[string]interface{}{siteItem})
	d.Set("edgecloud_id", server.EdgeCloudID)
	d.Set("edgecloud_name", server.EdgeCloudName)

	return nil
}
