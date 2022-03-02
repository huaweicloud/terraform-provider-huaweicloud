package huaweicloud

import (
	"context"
	"strings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceComputeInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
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
						"volume_attached": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_sys_volume": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"scheduler_hints": {
							Type:     schema.TypeList,
							Optional: true,
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

func queryEcsInstances(client *golangsdk.ServiceClient, opt *cloudservers.ListOpts) ([]cloudservers.CloudServer, error) {
	pages, err := cloudservers.List(client, opt).AllPages()
	if err != nil {
		return []cloudservers.CloudServer{}, fmtp.Errorf("Error getting cloud servers: %s", err)
	}
	return cloudservers.ExtractServers(pages)
}

func dataSourceComputeInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	ecsClient, err := conf.ComputeV1Client(GetRegion(d, conf))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud ECS v1 client: %s", err)
	}

	opt := buildListOptsWithoutIP(d, conf)

	allServers, err := queryEcsInstances(ecsClient, opt)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve cloud servers: %s", err)
	}

	servers, ids := filterCloudServers(d, allServers)
	// Save the data source ID using a hash code constructed using all instance IDs.
	d.SetId(hashcode.Strings(ids))

	return setComputeInstancesParams(d, conf, servers)
}

func IsSystemVolume(index string) bool {
	return index == "0"
}

func parseEcsInstanceSecurityGroupIds(groups []cloudservers.SecurityGroups) []string {
	result := make([]string, len(groups))

	for i, sg := range groups {
		result[i] = sg.ID
	}

	return result
}

func parseEcsInstanceSchedulerHintInfo(hints cloudservers.OsSchedulerHints) []map[string]interface{} {
	result := make([]map[string]interface{}, len(hints.Group))

	for i, val := range hints.Group {
		result[i] = map[string]interface{}{
			"group": val,
		}
	}

	return result
}

func parseEcsInstanceVolumeAttachedInfo(attachList []cloudservers.VolumeAttached) []map[string]interface{} {
	result := make([]map[string]interface{}, len(attachList))

	for i, volume := range attachList {
		result[i] = map[string]interface{}{
			"volume_id":     volume.ID,
			"is_sys_volume": IsSystemVolume(volume.BootIndex),
		}
	}

	return result
}

func parseEcsInstanceTagInfo(tags []string) map[string]interface{} {
	result := map[string]interface{}{}
	for _, tag := range tags {
		kv := strings.SplitN(tag, "=", 2)
		if len(kv) != 2 {
			logp.Printf("[WARN] Invalid key/value format of tag: %s", tag)
			continue
		}
		result[kv[0]] = kv[1]
	}

	return result
}

func setComputeInstancesParams(d *schema.ResourceData, config *config.Config,
	servers []cloudservers.CloudServer) diag.Diagnostics {
	result := make([]map[string]interface{}, len(servers))

	for i, val := range servers {
		server := map[string]interface{}{
			"id":                    val.ID,
			"user_data":             val.UserData,
			"name":                  val.Name,
			"flavor_name":           val.Flavor.Name,
			"status":                val.Status,
			"enterprise_project_id": val.EnterpriseProjectID,
			"flavor_id":             val.Flavor.ID,
			"image_id":              val.Image.ID,
			"availability_zone":     val.AvailabilityZone,
			"key_pair":              val.KeyName,
		}

		server["security_group_ids"] = parseEcsInstanceSecurityGroupIds(val.SecurityGroups)

		if len(val.OsSchedulerHints.Group) > 0 {
			server["scheduler_hints"] = parseEcsInstanceSchedulerHintInfo(val.OsSchedulerHints)
		}

		if len(val.VolumeAttached) > 0 {
			server["volume_attached"] = parseEcsInstanceVolumeAttachedInfo(val.VolumeAttached)
		}

		if len(val.Tags) > 0 {
			server["tags"] = parseEcsInstanceTagInfo(val.Tags)
		}

		result[i] = server
	}

	if err := d.Set("instances", result); err != nil {
		return fmtp.DiagErrorf("Error setting cloud server list: %s", err)
	}
	return nil
}
