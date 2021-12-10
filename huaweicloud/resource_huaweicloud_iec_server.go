package huaweicloud

import (
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/iec/v1/cloudvolumes"
	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/servers"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

var iecServerNicsSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

var iecVolumeAttachedSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_id": {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	},
}

func resourceIecServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceIecServerV1Create,
		Read:   resourceIecServerV1Read,
		Update: resourceIecServerV1Update,
		Delete: resourceIecServerV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"system_disk_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"coverage_sites": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"admin_pass", "key_pair"},
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bind_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"coverage_level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SITE",
			},
			"coverage_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "centralize",
				ValidateFunc: validation.StringInSlice([]string{
					"centralize", "discrete",
				}, true),
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc: utils.HashAndHexEncode,
			},

			// computed fields
			"edgecloud_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edgecloud_name": {
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
			"origin_server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildServerSecGroups(d *schema.ResourceData) []common.SecurityGroup {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]common.SecurityGroup, len(rawSecGroups))

	for i, raw := range rawSecGroups {
		secgroups[i] = common.SecurityGroup{
			ID: raw.(string),
		}
	}
	return secgroups
}

func buildNetworkConfig(d *schema.ResourceData) common.NetConfig {
	netOpts := common.NetConfig{}

	rawSubnets := d.Get("subnet_ids").([]interface{})
	subents := make([]common.SubnetID, len(rawSubnets))
	for i, raw := range rawSubnets {
		subents[i] = common.SubnetID{
			ID: raw.(string),
		}
	}
	netOpts.Subnets = subents
	netOpts.VpcID = d.Get("vpc_id").(string)
	netOpts.NicNum = len(rawSubnets)

	return netOpts
}

func buildServerRootVolume(d *schema.ResourceData) common.RootVolume {
	rootVolume := common.RootVolume{
		VolumeType: d.Get("system_disk_type").(string),
		Size:       d.Get("system_disk_size").(int),
	}

	return rootVolume
}

func buildServerDataVolumes(d *schema.ResourceData) []common.DataVolume {
	rawVols := d.Get("data_disks").([]interface{})
	volList := make([]common.DataVolume, len(rawVols))

	for i, v := range rawVols {
		vol := v.(map[string]interface{})
		volList[i] = common.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
		}
	}

	return volList
}

func buildServerCoverage(d *schema.ResourceData) common.Coverage {
	rawSites := d.Get("coverage_sites").([]interface{})
	sitesList := make([]common.CoverageSite, len(rawSites))

	for i, v := range rawSites {
		site := v.(map[string]interface{})
		sitesList[i] = common.CoverageSite{
			Site: site["site_id"].(string),
			Demands: []common.Demand{
				{
					Operator: site["operator"].(string),
					Count:    1,
				},
			},
		}
	}

	var coverageOpts = common.Coverage{
		CoveragePolicy: d.Get("coverage_policy").(string),
		CoverageLevel:  d.Get("coverage_level").(string),
		CoverageSites:  sitesList,
	}
	logp.Printf("[DEBUG] servers coverage options: %+v", coverageOpts)

	return coverageOpts
}

func resourceIecServerV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	resourceOpts := common.ResourceOpts{
		Count:          1,
		Name:           d.Get("name").(string),
		ImageRef:       d.Get("image_id").(string),
		FlavorRef:      d.Get("flavor_id").(string),
		NetConfig:      buildNetworkConfig(d),
		SecurityGroups: buildServerSecGroups(d),
		RootVolume:     buildServerRootVolume(d),
		DataVolumes:    buildServerDataVolumes(d),
	}
	if d.Get("bind_eip").(bool) {
		resourceOpts.BandWidth = &common.BandWidth{
			ShareType: "WHOLE",
		}
	}

	createOpts := servers.CreateOpts{
		ResourceOpts: resourceOpts,
		Coverage:     buildServerCoverage(d),
	}
	logp.Printf("[DEBUG] Create IEC servers options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	if v, ok := d.GetOk("admin_pass"); ok {
		createOpts.AdminPass = v.(string)
	} else {
		createOpts.KeyName = d.Get("key_pair").(string)
	}

	resp, err := servers.CreateServer(iecClient, createOpts)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC server: %s", err)
	}

	jobID := resp.Job.Id
	serverID := resp.ServerIDs.IDs[0]
	logp.Printf("[INFO] job ID: %s, servers ID: %s", jobID, serverID)
	// Store the ID now
	d.SetId(serverID)

	// Wait for the servers to become running
	logp.Printf("[DEBUG] waiting for IEC server (%s) to become running", serverID)

	// Pending state "DELETED" means the instance has not be ready
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETED", "BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    serverStateRefreshFunc(iecClient, serverID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error waiting for IEC server (%s) to become ready: %s", serverID, err)
	}

	// CreateServer will add an prefix "IEC-xxx-" for the instance name, we should update it.
	serverName := d.Get("name").(string)
	updateOpts := servers.UpdateInstance{
		UpdateServer: servers.UpdateOpts{
			Name: &serverName,
		},
	}
	_, err = servers.UpdateServer(iecClient, updateOpts, d.Id()).ExtractUpdateToServer()
	if err != nil {
		logp.Printf("[WARN] Updating name of HuaweiCloud IEC server (%s) failed: %s", serverID, err)
	}

	return resourceIecServerV1Read(d, meta)
}

func resourceIecServerV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	servers, err := servers.GetServer(iecClient, d.Id()).ExtractServerDetail()
	if err != nil {
		return CheckDeleted(d, err, "iec server")
	}

	edgeServer := servers.Server
	logp.Printf("[DEBUG] Retrieved server %s: %+v", d.Id(), edgeServer)

	allNics, eip := expandIecServerNics(edgeServer)
	allVolumes, sysDiskID := expandIecServerVolumeAttached(iecClient, edgeServer)

	mErr := multierror.Append(
		d.Set("name", edgeServer.Name),
		d.Set("status", edgeServer.Status),
		d.Set("edgecloud_id", edgeServer.EdgeCloudID),
		d.Set("edgecloud_name", edgeServer.EdgeCloudName),
		d.Set("origin_server_id", edgeServer.ServerID),
		d.Set("flavor_id", edgeServer.Flavor.ID),
		d.Set("flavor_name", edgeServer.Flavor.Name),
		d.Set("image_name", edgeServer.Metadata.ImageName),
		d.Set("nics", allNics),
		d.Set("public_ip", eip),
		d.Set("volume_attached", allVolumes),
		d.Set("system_disk_id", sysDiskID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting fields: %s", err)
	}

	if vpcID := edgeServer.Metadata.VpcID; vpcID != "" {
		d.Set("vpc_id", vpcID)
	}
	if imageID := edgeServer.Image.ID; imageID != "" {
		d.Set("image_id", imageID)
	}

	return nil
}

func resourceIecServerV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	if d.HasChange("name") {
		serverName := d.Get("name").(string)
		updateOpts := servers.UpdateInstance{
			UpdateServer: servers.UpdateOpts{
				Name: &serverName,
			},
		}

		_, err := servers.UpdateServer(iecClient, updateOpts, d.Id()).ExtractUpdateToServer()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud IEC server: %s", err)
		}

		return resourceIecServerV1Read(d, meta)
	}

	return nil
}

func resourceIecServerV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting HuaweiCloud servers %s", d.Id())
	deleteOpts := servers.DeleteOpts{
		Servers: []cloudservers.Server{
			{
				Id: d.Id(),
			},
		},
	}
	err = servers.DeleteServers(iecClient, deleteOpts).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud server: %s", err)
	}

	// Wait for the servers to delete before moving on.
	logp.Printf("[DEBUG] Waiting for servers (%s) to delete", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTOFF"},
		Target:     []string{"DELETED", "SOFT_DELETED"},
		Refresh:    serverStateRefreshFunc(iecClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error waiting for servers (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

// serverStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an HuaweiCloud IEC servers.
func serverStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := servers.GetServer(client, id).ExtractServerDetail()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		// get fault message when status is ERROR
		if s.Server.Status == "ERROR" {
			return s, "ERROR", fmtp.Errorf("the edge instance is error")
		}
		return s, s.Server.Status, nil
	}
}

func expandIecServerNics(edgeServer *servers.Server) ([]map[string]interface{}, string) {
	var publicIP string
	allNics := make([]map[string]interface{}, 0)

	for _, val := range edgeServer.Addresses {
		for _, nicRaw := range val {
			if nicRaw.Type == "floating" {
				publicIP = nicRaw.Addr
				continue
			}

			nicItem := map[string]interface{}{
				"port":    nicRaw.PortID,
				"mac":     nicRaw.MacAddr,
				"address": nicRaw.Addr,
			}
			allNics = append(allNics, nicItem)
		}
	}
	return allNics, publicIP
}

func expandIecServerVolumeAttached(client *golangsdk.ServiceClient, edgeServer *servers.Server) ([]map[string]interface{}, string) {
	var sysDiskID string
	allVolumes := make([]map[string]interface{}, 0, len(edgeServer.VolumeAttached))

	for _, disk := range edgeServer.VolumeAttached {
		if disk.BootIndex == "0" {
			sysDiskID = disk.ID
		}

		volumeInfo, err := cloudvolumes.Get(client, disk.ID).Extract()
		if err != nil {
			logp.Printf("[WARN] failed to retrieve volume %s: %s", disk.ID, err)
			continue
		}

		logp.Printf("[DEBUG] Retrieved volume %s: %#v", disk.ID, volumeInfo)
		volumeItem := map[string]interface{}{
			"volume_id":  disk.ID,
			"boot_index": disk.BootIndex,
			"device":     disk.Device,
			"size":       volumeInfo.Volume.Size,
			"type":       volumeInfo.Volume.VolumeType,
		}
		allVolumes = append(allVolumes, volumeItem)
	}

	return allVolumes, sysDiskID
}
