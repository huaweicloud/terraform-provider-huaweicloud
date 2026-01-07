package bms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bms/v1/baremetalservers"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API BMS GET /v1/{project_id}/baremetalservers/detail
// @API BMS GET /v1/{project_id}/baremetalservers/{server_id}/tags
func DataSourceBmsInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBmsInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Elem:     serversSchema(),
				Computed: true,
			},
		},
	}
}

func serversSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsComputedSchema(),
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volumes_attached": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_on_termination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_index": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vm_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"launched_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceBmsInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bmsClient, err := cfg.BmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	listOpts := baremetalservers.ListOpts{
		FlavorId:            d.Get("flavor_id").(string),
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
		Tags:                d.Get("tags").(string),
	}

	allInstances, err := baremetalservers.List(bmsClient, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve BMS instances: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("servers", flattenListBMSInstances(d, meta, bmsClient, allInstances)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListBMSInstances(d *schema.ResourceData, meta interface{}, bmsClient *golangsdk.ServiceClient,
	allInstances []baremetalservers.CloudServer) []map[string]interface{} {
	if allInstances == nil {
		return nil
	}

	servers := make([]map[string]interface{}, len(allInstances))
	for i, v := range allInstances {
		servers[i] = map[string]interface{}{
			"id":                    v.ID,
			"name":                  v.Name,
			"status":                v.Status,
			"nics":                  flattenBmsInstanceNicsV1(d, meta, v.Addresses),
			"flavor_id":             v.Flavor.ID,
			"flavor_name":           v.Flavor.Name,
			"vcpus":                 v.Flavor.Vcpus,
			"memory":                v.Flavor.RAM,
			"disk":                  v.Flavor.Disk,
			"created_at":            v.Created,
			"updated_at":            v.Updated,
			"launched_at":           v.Launched,
			"vpc_id":                v.Metadata.VpcID,
			"agency_name":           v.Metadata.AgencyName,
			"locked":                v.Locked,
			"tags":                  faltternBmsTags(v.ID, bmsClient),
			"user_id":               v.UserID,
			"key_pair":              v.KeyName,
			"description":           v.Description,
			"volumes_attached":      falttenVolumesAtached(v.VolumeAttached),
			"vm_state":              v.VMState,
			"disk_config":           v.DiskConfig,
			"availability_zone":     v.AvailabilityZone,
			"root_device_name":      v.RootDeviceName,
			"enterprise_project_id": v.EnterpriseProjectID,
			"user_data":             v.UserData,
			"security_groups":       flattenSecurityGroupIDs(v.SecurityGroups),
			"image_id":              v.Image.ID,
			"image_name":            v.Metadata.ImageName,
			"image_type":            v.Metadata.Imagetype,
		}
	}
	return servers
}

func flattenBmsInstanceNicsV1(d *schema.ResourceData, meta interface{},
	addresses map[string][]baremetalservers.Address) []map[string]interface{} {
	cfg := meta.(*config.Config)
	networkingClient, err := cfg.NetworkingV2Client(cfg.GetRegion(d))
	if err != nil {
		log.Printf("Error creating networking client: %s", err)
	}

	var network string
	var nics []map[string]interface{}
	// Loop through all networks and addresses.
	for _, addrs := range addresses {
		for _, addr := range addrs {
			// Skip if not fixed ip
			if addr.Type != "fixed" {
				continue
			}

			p, err := ports.Get(networkingClient, addr.PortID).Extract()
			if err != nil {
				network = ""
				log.Printf("[DEBUG] flattenInstanceNicsV1: failed to fetch port %s", addr.PortID)
			} else {
				network = p.NetworkID
			}

			v := map[string]interface{}{
				"subnet_id":   network,
				"ip_address":  addr.Addr,
				"mac_address": addr.MacAddr,
				"port_id":     addr.PortID,
			}
			nics = append(nics, v)
		}
	}

	log.Printf("[DEBUG] flattenInstanceNicsV1: %#v", nics)
	return nics
}

func faltternBmsTags(instanceId string, bmsClient *golangsdk.ServiceClient) map[string]string {
	var tagsmap map[string]string
	resourceTags, err := tags.Get(bmsClient, "baremetalservers", instanceId).Extract()
	if err != nil {
		return nil
	}
	tagsmap = utils.TagsToMap(resourceTags.Tags)
	return tagsmap
}

func flattenSecurityGroupIDs(securityGroups []baremetalservers.SecurityGroups) []string {
	var secGrpIDs []string
	for _, sg := range securityGroups {
		secGrpIDs = append(secGrpIDs, sg.ID)
	}

	return secGrpIDs
}

func falttenVolumesAtached(volumesAttached []baremetalservers.VolumeAttached) []map[string]interface{} {
	var volumes []map[string]interface{}
	for _, va := range volumesAttached {
		v := map[string]interface{}{
			"id":                    va.ID,
			"delete_on_termination": va.DeleteOnTermination,
			"boot_index":            va.BootIndex,
			"device":                va.Device,
		}
		volumes = append(volumes, v)
	}
	return volumes
}
