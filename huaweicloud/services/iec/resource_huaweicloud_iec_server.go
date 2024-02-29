package iec

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/iec/v1/cloudvolumes"
	ieccommon "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/servers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var serverNicsSchema = &schema.Schema{
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

var volumeAttachedSchema = &schema.Schema{
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

// @API IEC POST /v1/cloudservers/delete
// @API IEC GET /v1/cloudservers/{server_id}
// @API IEC PUT /v1/cloudservers/{server_id}
// @API IEC POST /v1/cloudservers
// @API IEC GET /v1/cloudvolumes/{volume_id}
func ResourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

func buildServerSecGroups(d *schema.ResourceData) []ieccommon.SecurityGroup {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]ieccommon.SecurityGroup, len(rawSecGroups))

	for i, raw := range rawSecGroups {
		secgroups[i] = ieccommon.SecurityGroup{
			ID: raw.(string),
		}
	}
	return secgroups
}

func buildNetworkConfig(d *schema.ResourceData) ieccommon.NetConfig {
	netOpts := ieccommon.NetConfig{}

	rawSubnets := d.Get("subnet_ids").([]interface{})
	subents := make([]ieccommon.SubnetID, len(rawSubnets))
	for i, raw := range rawSubnets {
		subents[i] = ieccommon.SubnetID{
			ID: raw.(string),
		}
	}
	netOpts.Subnets = subents
	netOpts.VpcID = d.Get("vpc_id").(string)
	netOpts.NicNum = len(rawSubnets)

	return netOpts
}

func buildServerRootVolume(d *schema.ResourceData) ieccommon.RootVolume {
	rootVolume := ieccommon.RootVolume{
		VolumeType: d.Get("system_disk_type").(string),
		Size:       d.Get("system_disk_size").(int),
	}

	return rootVolume
}

func buildServerDataVolumes(d *schema.ResourceData) []ieccommon.DataVolume {
	rawVols := d.Get("data_disks").([]interface{})
	volList := make([]ieccommon.DataVolume, len(rawVols))

	for i, v := range rawVols {
		vol := v.(map[string]interface{})
		volList[i] = ieccommon.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
		}
	}

	return volList
}

func buildServerCoverage(d *schema.ResourceData) ieccommon.Coverage {
	rawSites := d.Get("coverage_sites").([]interface{})
	sitesList := make([]ieccommon.CoverageSite, len(rawSites))

	for i, v := range rawSites {
		site := v.(map[string]interface{})
		sitesList[i] = ieccommon.CoverageSite{
			Site: site["site_id"].(string),
			Demands: []ieccommon.Demand{
				{
					Operator: site["operator"].(string),
					Count:    1,
				},
			},
		}
	}

	var coverageOpts = ieccommon.Coverage{
		CoveragePolicy: d.Get("coverage_policy").(string),
		CoverageLevel:  d.Get("coverage_level").(string),
		CoverageSites:  sitesList,
	}
	log.Printf("[DEBUG] servers coverage options: %+v", coverageOpts)

	return coverageOpts
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	resourceOpts := ieccommon.ResourceOpts{
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
		resourceOpts.BandWidth = &ieccommon.BandWidth{
			ShareType: "WHOLE",
		}
	}

	createOpts := servers.CreateOpts{
		ResourceOpts: resourceOpts,
		Coverage:     buildServerCoverage(d),
	}
	log.Printf("[DEBUG] create IEC servers options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	if v, ok := d.GetOk("admin_pass"); ok {
		createOpts.AdminPass = v.(string)
	} else {
		createOpts.KeyName = d.Get("key_pair").(string)
	}

	resp, err := servers.CreateServer(iecClient, createOpts)
	if err != nil {
		return diag.Errorf("error creating IEC server: %s", err)
	}

	jobID := resp.Job.Id
	serverID := resp.ServerIDs.IDs[0]
	log.Printf("[INFO] job ID: %s, servers ID: %s", jobID, serverID)
	// Store the ID now
	d.SetId(serverID)

	// Wait for the servers to become running
	log.Printf("[DEBUG] waiting for IEC server (%s) to become running", serverID)

	// Pending state "DELETED" means the instance has not be ready
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETED", "BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    serverStateRefreshFunc(iecClient, serverID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for IEC server (%s) to become ready: %s", serverID, err)
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
		log.Printf("[WARN] updating name of IEC server (%s) failed: %s", serverID, err)
	}

	return resourceServerRead(ctx, d, meta)
}

func resourceServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	iecServers, err := servers.GetServer(iecClient, d.Id()).ExtractServerDetail()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "iec server")
	}

	edgeServer := iecServers.Server
	log.Printf("[DEBUG] retrieved server %s: %+v", d.Id(), edgeServer)

	allNics, eip := expandServerNics(edgeServer)
	allVolumes, sysDiskID := expandServerVolumeAttached(iecClient, edgeServer)

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
		return diag.Errorf("error setting fields: %s", err)
	}

	if vpcID := edgeServer.Metadata.VpcID; vpcID != "" {
		d.Set("vpc_id", vpcID)
	}
	if imageID := edgeServer.Image.ID; imageID != "" {
		d.Set("image_id", imageID)
	}

	return nil
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
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
			return diag.Errorf("error updating IEC server: %s", err)
		}

		return resourceServerRead(ctx, d, meta)
	}

	return nil
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	log.Printf("[DEBUG] deleting servers %s", d.Id())
	deleteOpts := servers.DeleteOpts{
		Servers: []cloudservers.Server{
			{
				Id: d.Id(),
			},
		},
	}
	err = servers.DeleteServers(iecClient, deleteOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting server: %s", err)
	}

	// Wait for the servers to delete before moving on.
	log.Printf("[DEBUG] waiting for servers (%s) to delete", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTOFF"},
		Target:     []string{"DELETED", "SOFT_DELETED"},
		Refresh:    serverStateRefreshFunc(iecClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for servers (%s) to delete: %s", d.Id(), err)
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
			return s, "ERROR", fmt.Errorf("the edge instance is error")
		}
		return s, s.Server.Status, nil
	}
}

func expandServerNics(edgeServer *servers.Server) ([]map[string]interface{}, string) {
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

func expandServerVolumeAttached(client *golangsdk.ServiceClient, edgeServer *servers.Server) ([]map[string]interface{}, string) {
	var sysDiskID string
	allVolumes := make([]map[string]interface{}, 0, len(edgeServer.VolumeAttached))

	for _, disk := range edgeServer.VolumeAttached {
		if disk.BootIndex == "0" {
			sysDiskID = disk.ID
		}

		volumeInfo, err := cloudvolumes.Get(client, disk.ID).Extract()
		if err != nil {
			log.Printf("[WARN] failed to retrieve volume %s: %s", disk.ID, err)
			continue
		}

		log.Printf("[DEBUG] retrieved volume %s: %#v", disk.ID, volumeInfo)
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
