package ecs

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ECS GET /v1/{project_id}/cloudservers/detail
// @API ECS GET /v1/{project_id}/cloudservers/{server_id}/block_device
// @API EVS GET /v2/{project_id}/cloudvolumes/{volume_id}
// @API VPC GET /v2.0/ports/{id}
// @API ECS GET /v1.1/{project_id}/cloudservers/detail
func DataSourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeInstanceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// attributes
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
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_ids": {
				Type:     schema.TypeList,
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
			"network":         computedSchemaNetworks(),
			"volume_attached": computedSchemaVolumeAttached(),
			"scheduler_hints": computedSchemaSchedulerHints(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func computedSchemaNetworks() *schema.Schema {
	computedSchema := schema.Schema{
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
	}

	return &computedSchema
}

func computedSchemaVolumeAttached() *schema.Schema {
	computedSchema := schema.Schema{
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
				"is_sys_volume": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
	}

	return &computedSchema
}

func computedSchemaSchedulerHints() *schema.Schema {
	computedSchema := schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"group": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}

	return &computedSchema
}

func buildListOptsWithoutStatus(d *schema.ResourceData, cfg *config.Config) *cloudservers.ListOpts {
	result := cloudservers.ListOpts{
		Limit:               100,
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor_id").(string),
		IPEqual:             d.Get("fixed_ip_v4").(string),
		Tags:                buildQueryInstancesTagsOpts(d),
	}

	return &result
}

func queryEcsInstances(client *golangsdk.ServiceClient, opt *cloudservers.ListOpts) ([]cloudservers.CloudServer, error) {
	pages, err := cloudservers.List(client, opt).AllPages()
	if err != nil {
		return nil, err
	}
	return cloudservers.ExtractServers(pages)
}

func dataSourceComputeInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	ecsClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	opt := buildListOptsWithoutStatus(d, conf)
	allServers, err := queryEcsInstances(ecsClient, opt)
	if err != nil {
		return diag.Errorf("unable to retrieve ECS instances: %s", err)
	}

	filter := map[string]interface{}{
		"ID": d.Get("instance_id"),
	}

	filterServers, err := utils.FilterSliceWithField(allServers, filter)
	if err != nil {
		return diag.Errorf("filter ECS instances failed: %s", err)
	}

	if len(filterServers) < 1 {
		return diag.Errorf("your query returned no results, please change your search criteria and try again.")
	}
	if len(filterServers) > 1 {
		return diag.Errorf("your query returned more than one result, please try a more specific search criteria.")
	}

	server := filterServers[0].(cloudservers.CloudServer)
	log.Printf("[DEBUG] fetching the ECS instance: %#v", server)

	d.SetId(server.ID)
	return setEcsInstanceParams(d, conf, ecsClient, server)
}

func setEcsInstanceParams(d *schema.ResourceData, conf *config.Config, ecsClient *golangsdk.ServiceClient,
	server cloudservers.CloudServer) diag.Diagnostics {
	region := conf.GetRegion(d)
	networkingClient, err := conf.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v2 client: %s", err)
	}
	blockStorageClient, err := conf.BlockStorageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
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
		d.Set("tags", flattenEcsInstanceTags(server.Tags)),
		d.Set("security_group_ids", flattenEcsInstanceSecurityGroupIds(server.SecurityGroups)),
		d.Set("security_groups", flattenEcsInstanceSecurityGroups(server.SecurityGroups)),
		d.Set("scheduler_hints", flattenEcsInstanceSchedulerHints(server.OsSchedulerHints)),
		d.Set("charging_mode", normalizeChargingMode(server.Metadata.ChargingMode)),

		setEcsInstanceNetworks(d, networkingClient, server.Addresses),
		setEcsInstanceVolumeAttached(d, ecsClient, blockStorageClient, server.VolumeAttached),
	)

	// Set expired time for prePaid instance
	if normalizeChargingMode(server.Metadata.ChargingMode) == "prePaid" {
		expiredTime, err := getPrePaidExpiredTime(d, conf, server.ID)
		if err != nil {
			log.Printf("error get prePaid expired time: %s", err)
		}

		mErr = multierror.Append(mErr, d.Set("expired_time", expiredTime))
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func setEcsInstanceNetworks(d *schema.ResourceData, client *golangsdk.ServiceClient,
	addresses map[string][]cloudservers.Address) error {
	if len(addresses) == 0 {
		return nil
	}

	networks, eip := flattenEcsInstanceNetworks(client, addresses)
	mErr := multierror.Append(nil,
		d.Set("network", networks),
		d.Set("public_ip", eip),
	)
	return mErr.ErrorOrNil()
}

func setEcsInstanceVolumeAttached(d *schema.ResourceData, ecsClient, evsClient *golangsdk.ServiceClient,
	attached []cloudservers.VolumeAttached) error {
	if len(attached) == 0 {
		return nil
	}

	vols, sysDiskID := flattenEcsInstanceVolumeAttached(ecsClient, evsClient, d.Id())
	mErr := multierror.Append(nil,
		d.Set("system_disk_id", sysDiskID),
		d.Set("volume_attached", vols),
	)
	return mErr.ErrorOrNil()
}

func flattenEcsInstanceSecurityGroups(groups []cloudservers.SecurityGroups) []string {
	if len(groups) == 0 {
		return nil
	}

	result := make([]string, len(groups))
	for i, sg := range groups {
		result[i] = sg.Name
	}
	return result
}

func flattenEcsInstanceSecurityGroupIds(groups []cloudservers.SecurityGroups) []string {
	if len(groups) == 0 {
		return nil
	}

	result := make([]string, len(groups))
	for i, sg := range groups {
		result[i] = sg.ID
	}
	return result
}

func flattenEcsInstanceSchedulerHints(hints cloudservers.OsSchedulerHints) []map[string]interface{} {
	if len(hints.Group) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(hints.Group))
	for i, val := range hints.Group {
		result[i] = map[string]interface{}{
			"group": val,
		}
	}
	return result
}

func flattenEcsInstanceTags(tags []string) map[string]interface{} {
	result := map[string]interface{}{}

	for _, tag := range tags {
		kv := strings.SplitN(tag, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else {
			result[kv[0]] = ""
		}
	}

	return result
}

// flattenEcsInstanceNetworks collects instance network information
func flattenEcsInstanceNetworks(client *golangsdk.ServiceClient,
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
				log.Printf("[WARN] failed to fetch port %s: %s", addr.PortID, err)
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

	log.Printf("[DEBUG] flatten Instance Networks: %#v", networks)
	return networks, publicIP
}

func flattenEcsInstanceVolumeAttached(ecsClient, evsClient *golangsdk.ServiceClient,
	instanceID string) ([]map[string]interface{}, string) {
	devices, err := block_devices.List(ecsClient, instanceID)
	if err != nil {
		log.Printf("[WARN] failed to retrieve volumes in %s: %s", instanceID, err)
		return nil, ""
	}

	var systemDiskID string
	allVolumes := make([]map[string]interface{}, len(devices))
	for i, vol := range devices {
		allVolumes[i] = map[string]interface{}{
			"volume_id":   vol.Id,
			"size":        vol.Size,
			"boot_index":  vol.BootIndex,
			"pci_address": vol.PciAddress,
		}

		if vol.BootIndex == 0 {
			allVolumes[i]["is_sys_volume"] = true
			systemDiskID = vol.Id
		}

		// retrieve volume type
		volumeInfo, err := cloudvolumes.Get(evsClient, vol.Id).Extract()
		if err != nil {
			log.Printf("[WARN] failed to retrieve volume %s: %s", vol.Id, err)
		} else {
			allVolumes[i]["type"] = volumeInfo.VolumeType
		}
	}
	return allVolumes, systemDiskID
}
