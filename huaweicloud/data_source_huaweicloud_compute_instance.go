package huaweicloud

import (
	"context"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeInstanceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fixed_ip_v4": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
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
			"security_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"system_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fixed_ip_v6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"volume_attached": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pci_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"scheduler_hints": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildListOptsWithoutStatus(d *schema.ResourceData, conf *config.Config) *cloudservers.ListOpts {
	result := cloudservers.ListOpts{
		Limit:               100,
		EnterpriseProjectID: conf.DataGetEnterpriseProjectID(d),
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor_id").(string),
		IP:                  d.Get("fixed_ip_v4").(string),
	}

	return &result
}

func parseEcsInstanceSecurityGroups(groups []cloudservers.SecurityGroups) []string {
	result := make([]string, len(groups))

	for i, sg := range groups {
		result[i] = sg.Name
	}

	return result
}

func setEcsInstanceSchedulerHints(d *schema.ResourceData, hints cloudservers.OsSchedulerHints) error {
	if len(hints.Group) > 0 {
		return d.Set("scheduler_hints", parseEcsInstanceSchedulerHintInfo(hints))
	}
	return nil
}

func setEcsInstanceTags(d *schema.ResourceData, tags []string) error {
	if len(tags) > 0 {
		return d.Set("tags", parseEcsInstanceTagInfo(tags))
	}
	return nil
}

func setEcsInstancePublicIp(d *schema.ResourceData, client *golangsdk.ServiceClient,
	addresses map[string][]cloudservers.Address) error {
	// Set the instance network and address information

	networks, eip := flattenComputeNetworks(d, client, addresses)
	mErr := multierror.Append(nil,
		d.Set("network", networks),
		d.Set("public_ip", eip),
	)
	return mErr.ErrorOrNil()
}

func setEcsInstanceVolumeAttached(d *schema.ResourceData, ecsClient, evsClient *golangsdk.ServiceClient,
	attached []cloudservers.VolumeAttached) error {
	// Set volume attached
	if len(attached) > 0 {
		bds := make([]map[string]interface{}, len(attached))
		for i, b := range attached {
			// retrieve volume `size` and `type`
			volumeInfo, err := cloudvolumes.Get(evsClient, b.ID).Extract()
			if err != nil {
				return err
			}
			logp.Printf("[DEBUG] Retrieved volume %s: %#v", b.ID, volumeInfo)

			// retrieve volume `pci_address`
			va, err := block_devices.Get(ecsClient, d.Id(), b.ID).Extract()
			if err != nil {
				return err
			}
			logp.Printf("[DEBUG] Retrieved block device %s: %#v", b.ID, va)

			bds[i] = map[string]interface{}{
				"volume_id":   b.ID,
				"size":        volumeInfo.Size,
				"type":        volumeInfo.VolumeType,
				"boot_index":  va.BootIndex,
				"pci_address": va.PciAddress,
			}

			if va.BootIndex == 0 {
				d.Set("system_disk_id", b.ID)
			}
		}
		return d.Set("volume_attached", bds)
	}
	return nil
}

func setEcsInstanceParams(d *schema.ResourceData, config *config.Config, ecsClient *golangsdk.ServiceClient,
	server cloudservers.CloudServer) diag.Diagnostics {
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v2 client: %s", err)
	}
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud EVS client: %s", err)
	}
	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("availability_zone", server.AvailabilityZone),
		d.Set("name", server.Name),
		d.Set("status", server.Status),
		d.Set("flavor_id", server.Flavor.ID),
		d.Set("flavor_name", server.Flavor.Name),
		d.Set("image_id", server.Image.ID),
		d.Set("image_name", server.Metadata.ImageName),
		d.Set("key_pair", server.KeyName),
		d.Set("user_data", server.UserData),
		d.Set("enterprise_project_id", server.EnterpriseProjectID),
		d.Set("security_group_ids", parseEcsInstanceSecurityGroupIds(server.SecurityGroups)),
		d.Set("security_groups", parseEcsInstanceSecurityGroups(server.SecurityGroups)),
		setEcsInstanceSchedulerHints(d, server.OsSchedulerHints),
		setEcsInstancePublicIp(d, networkingClient, server.Addresses),
		setEcsInstanceVolumeAttached(d, ecsClient, blockStorageClient, server.VolumeAttached),
		setEcsInstanceTags(d, server.Tags),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func dataSourceComputeInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud ECS v1 client: %s", err)
	}

	opt := buildListOptsWithoutStatus(d, config)
	allServers, err := queryEcsInstances(ecsClient, opt)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve cloud servers: %s", err)
	}

	if len(allServers) < 1 {
		return fmtp.DiagErrorf("Your query returned no results, please change your search criteria and try again.")
	}
	if len(allServers) > 1 {
		return fmtp.DiagErrorf("Your query returned more than one result, please try a more specific search criteria.")
	}

	server := allServers[0]
	logp.Printf("[DEBUG] fetching the ecs instance: %#v", server)

	d.SetId(server.ID)

	// Set instance parameters

	return setEcsInstanceParams(d, config, ecsClient, server)
}

// flattenComputeNetworks collects instance network information
func flattenComputeNetworks(d *schema.ResourceData, client *golangsdk.ServiceClient,
	addressResp map[string][]cloudservers.Address) ([]map[string]interface{}, string) {
	publicIP := ""
	networks := []map[string]interface{}{}

	for _, addresses := range addressResp {
		for _, addr := range addresses {
			if addr.Type == "floating" {
				publicIP = addr.Addr
				continue
			}

			// get networkID
			var networkID string
			p, err := ports.Get(client, addr.PortID).Extract()
			if err != nil {
				networkID = ""
				logp.Printf("[DEBUG] failed to fetch port %s", addr.PortID)
			} else {
				networkID = p.NetworkID
			}

			v := map[string]interface{}{
				"uuid": networkID,
				"port": addr.PortID,
				"mac":  addr.MacAddr,
			}
			if addr.Version == "6" {
				v["fixed_ip_v6"] = addr.Addr
			} else {
				v["fixed_ip_v4"] = addr.Addr
			}

			networks = append(networks, v)
		}
	}

	logp.Printf("[DEBUG] flatten Instance Networks: %#v", networks)
	return networks, publicIP
}
