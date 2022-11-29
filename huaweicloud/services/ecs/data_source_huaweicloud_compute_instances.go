package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

func DataSourceComputeInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
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
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func buildListOptsWithoutIP(d *schema.ResourceData, conf *config.Config) *cloudservers.ListOpts {
	result := cloudservers.ListOpts{
		Limit:               100,
		EnterpriseProjectID: conf.DataGetEnterpriseProjectID(d),
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor_id").(string),
		Status:              d.Get("status").(string),
	}

	return &result
}

func filterCloudServers(d *schema.ResourceData, servers []cloudservers.CloudServer) ([]cloudservers.CloudServer,
	[]string) {
	result := make([]cloudservers.CloudServer, 0, len(servers))
	ids := make([]string, 0, len(servers))

	for _, server := range servers {
		if serverId, ok := d.GetOk("instance_id"); ok && serverId != server.ID {
			continue
		}
		if flavorName, ok := d.GetOk("flavor_name"); ok && flavorName != server.Flavor.Name {
			continue
		}
		if iamgeId, ok := d.GetOk("image_id"); ok && iamgeId != server.Image.ID {
			continue
		}
		if az, ok := d.GetOk("availability_zone"); ok && az != server.AvailabilityZone {
			continue
		}
		if keypair, ok := d.GetOk("key_pair"); ok && keypair != server.KeyName {
			continue
		}
		result = append(result, server)
		ids = append(ids, server.ID)
	}

	return result, ids
}

func dataSourceComputeInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	ecsClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	opts := buildListOptsWithoutIP(d, conf)
	allServers, err := queryEcsInstances(ecsClient, opts)
	if err != nil {
		return diag.Errorf("unable to retrieve cloud servers: %s", err)
	}

	servers, ids := filterCloudServers(d, allServers)
	// Save the data source ID using a hash code constructed using all instance IDs.
	d.SetId(hashcode.Strings(ids))

	return setComputeInstancesParams(d, conf, servers)
}

func setComputeInstancesParams(d *schema.ResourceData, conf *config.Config, servers []cloudservers.CloudServer) diag.Diagnostics {
	region := conf.GetRegion(d)
	ecsClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	evsClient, err := conf.BlockStorageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	networkingClient, err := conf.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v2 client: %s", err)
	}

	result := make([]map[string]interface{}, len(servers))
	for i, item := range servers {
		server := map[string]interface{}{
			"id":                    item.ID,
			"name":                  item.Name,
			"image_id":              item.Image.ID,
			"image_name":            item.Metadata.ImageName,
			"flavor_id":             item.Flavor.ID,
			"flavor_name":           item.Flavor.Name,
			"status":                item.Status,
			"availability_zone":     item.AvailabilityZone,
			"enterprise_project_id": item.EnterpriseProjectID,
			"user_data":             item.UserData,
			"key_pair":              item.KeyName,
			"tags":                  flattenEcsInstanceTags(item.Tags),
			"security_group_ids":    flattenEcsInstanceSecurityGroupIds(item.SecurityGroups),
			"scheduler_hints":       flattenEcsInstanceSchedulerHints(item.OsSchedulerHints),
		}

		if len(item.VolumeAttached) > 0 {
			vols, sysDiskID := flattenEcsInstanceVolumeAttached(ecsClient, evsClient, item.ID)
			server["volume_attached"] = vols
			server["system_disk_id"] = sysDiskID

		}

		if len(item.Addresses) > 0 {
			networks, eip := flattenEcsInstanceNetworks(networkingClient, item.Addresses)
			server["network"] = networks
			server["public_ip"] = eip
		}

		result[i] = server
	}

	if err := d.Set("instances", result); err != nil {
		return diag.Errorf("error setting cloud server list: %s", err)
	}
	return nil
}
