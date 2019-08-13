package huaweicloud

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/secgroups"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/tags"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func resourceEcsInstanceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEcsInstanceV1Create,
		Read:   resourceEcsInstanceV1Read,
		Update: resourceEcsInstanceV1Update,
		Delete: resourceEcsInstanceV1Delete,

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
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						hash := sha1.Sum([]byte(v.(string)))
						return hex.EncodeToString(hash[:])
					default:
						return ""
					}
				},
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 12,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"system_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SATA",
				ValidateFunc: validation.StringInSlice([]string{
					"SATA", "SAS", "SSD", "co-p1", "uh-l1",
				}, true),
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"data_disks": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 23,
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
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "postPaid",
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, true),
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "month",
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, true),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: validateECSTagValue,
			},
			"auto_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"delete_disks_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceEcsInstanceV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV11Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute V1.1 client: %s", err)
	}
	computeV1Client, err := config.computeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute V1 client: %s", err)
	}

	createOpts := &cloudservers.CreateOpts{
		Name:             d.Get("name").(string),
		ImageRef:         d.Get("image_id").(string),
		FlavorRef:        d.Get("flavor").(string),
		KeyName:          d.Get("key_name").(string),
		VpcId:            d.Get("vpc_id").(string),
		SecurityGroups:   resourceInstanceSecGroupsV1(d),
		AvailabilityZone: d.Get("availability_zone").(string),
		Nics:             resourceInstanceNicsV1(d),
		RootVolume:       resourceInstanceRootVolumeV1(d),
		DataVolumes:      resourceInstanceDataVolumesV1(d),
		AdminPass:        d.Get("password").(string),
		UserData:         []byte(d.Get("user_data").(string)),
	}

	if d.Get("charging_mode") == "prePaid" {
		extendparam := cloudservers.ServerExtendParam{
			ChargingMode: d.Get("charging_mode").(string),
			PeriodType:   d.Get("period_unit").(string),
			PeriodNum:    d.Get("period").(int),
		}
		createOpts.ExtendParam = &extendparam
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	n, err := cloudservers.Create(computeClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud server: %s", err)
	}

	if err := cloudservers.WaitForJobSuccess(computeV1Client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	entity, err := cloudservers.GetJobEntity(computeV1Client, n.JobID, "server_id")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		d.SetId(id)

		if hasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			log.Printf("[DEBUG] Setting tags: %v", tagmap)
			err = setTagForInstance(d, meta, id, tagmap)
			if err != nil {
				log.Printf("[WARN] Error setting tags of instance:%s, err=%s", id, err)
			}
		}

		if hasFilledOpt(d, "auto_recovery") {
			ar := d.Get("auto_recovery").(bool)
			log.Printf("[DEBUG] Set auto recovery of instance to %t", ar)
			err = setAutoRecoveryForInstance(d, meta, id, ar)
			if err != nil {
				log.Printf("[WARN] Error setting auto recovery of instance:%s, err=%s", id, err)
			}
		}

		return resourceEcsInstanceV1Read(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceEcsInstanceV1Create.")
}

func resourceEcsInstanceV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	server, err := cloudservers.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "server")
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)
	d.Set("image_id", server.Image.ID)
	d.Set("flavor", server.Flavor.ID)
	d.Set("password", d.Get("password"))
	d.Set("key_name", server.KeyName)
	d.Set("vpc_id", server.Metadata.VpcID)
	d.Set("availability_zone", server.AvailabilityZone)

	secGrpNames := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpNames = append(secGrpNames, sg.Name)
	}
	d.Set("security_groups", secGrpNames)

	// Get the instance network and address information
	nics := flattenInstanceNicsV1(d, meta, server.Addresses)
	d.Set("nics", nics)

	// Set instance tags
	if _, ok := d.GetOk("tags"); ok {
		Taglist, err := tags.Get(computeClient, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("Error fetching HuaweiCloud instance tags: %s", err)
		}

		tagmap := make(map[string]string)
		for _, val := range Taglist.Tags {
			tagmap[val.Key] = val.Value
		}
		if err := d.Set("tags", tagmap); err != nil {
			return fmt.Errorf("[DEBUG] Error saving tag to state for HuaweiCloud instance (%s): %s", d.Id(), err)
		}
	}

	ar, err := resourceECSAutoRecoveryV1Read(d, meta, d.Id())
	if err != nil && !isResourceNotFound(err) {
		return fmt.Errorf("Error reading auto recovery of instance:%s, err=%s", d.Id(), err)
	}
	d.Set("auto_recovery", ar)

	return nil
}

func resourceEcsInstanceV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	var updateOpts servers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if updateOpts != (servers.UpdateOpts{}) {
		_, err := servers.Update(computeClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud server: %s", err)
		}
	}

	if d.HasChange("security_groups") {
		oldSGRaw, newSGRaw := d.GetChange("security_groups")
		oldSGSet := oldSGRaw.(*schema.Set)
		newSGSet := newSGRaw.(*schema.Set)
		secgroupsToAdd := newSGSet.Difference(oldSGSet)
		secgroupsToRemove := oldSGSet.Difference(newSGSet)

		log.Printf("[DEBUG] Security groups to add: %v", secgroupsToAdd)

		log.Printf("[DEBUG] Security groups to remove: %v", secgroupsToRemove)

		for _, g := range secgroupsToRemove.List() {
			err := secgroups.RemoveServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					continue
				}

				return fmt.Errorf("Error removing security group (%s) from HuaweiCloud server (%s): %s", g, d.Id(), err)
			} else {
				log.Printf("[DEBUG] Removed security group (%s) from instance (%s)", g, d.Id())
			}
		}

		for _, g := range secgroupsToAdd.List() {
			err := secgroups.AddServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return fmt.Errorf("Error adding security group (%s) to HuaweiCloud server (%s): %s", g, d.Id(), err)
			}
			log.Printf("[DEBUG] Added security group (%s) to instance (%s)", g, d.Id())
		}
	}

	if d.HasChange("flavor") {
		newFlavorId := d.Get("flavor").(string)

		resizeOpts := &servers.ResizeOpts{
			FlavorRef: newFlavorId,
		}
		log.Printf("[DEBUG] Resize configuration: %#v", resizeOpts)
		err := servers.Resize(computeClient, d.Id(), resizeOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error resizing HuaweiCloud server: %s", err)
		}

		// Wait for the instance to finish resizing.
		log.Printf("[DEBUG] Waiting for instance (%s) to finish resizing", d.Id())

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"RESIZE"},
			Target:     []string{"VERIFY_RESIZE"},
			Refresh:    ServerV2StateRefreshFunc(computeClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for instance (%s) to resize: %s", d.Id(), err)
		}

		// Confirm resize.
		log.Printf("[DEBUG] Confirming resize")
		err = servers.ConfirmResize(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error confirming resize of HuaweiCloud server: %s", err)
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"VERIFY_RESIZE"},
			Target:     []string{"ACTIVE"},
			Refresh:    ServerV2StateRefreshFunc(computeClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for instance (%s) to confirm resize: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		computeClient, err := config.computeV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute v1 client: %s", err)
		}
		oldTags, err := tags.Get(computeClient, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("Error fetching HuaweiCloud instance tags: %s", err)
		}
		if len(oldTags.Tags) > 0 {
			deleteopts := tags.BatchOpts{Action: tags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := tags.BatchAction(computeClient, d.Id(), deleteopts)
			if deleteTags.Err != nil {
				return fmt.Errorf("Error updating HuaweiCloud instance tags: %s", deleteTags.Err)
			}
		}

		if hasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForInstance(d, meta, d.Id(), tagmap)
				if err != nil {
					return fmt.Errorf("Error updating tags of instance:%s, err:%s", d.Id(), err)
				}
			}
		}
	}

	if d.HasChange("auto_recovery") {
		ar := d.Get("auto_recovery").(bool)
		log.Printf("[DEBUG] Update auto recovery of instance to %t", ar)
		err = setAutoRecoveryForInstance(d, meta, d.Id(), ar)
		if err != nil {
			return fmt.Errorf("Error updating auto recovery of instance:%s, err:%s", d.Id(), err)
		}
	}

	return resourceEcsInstanceV1Read(d, meta)
}

func resourceEcsInstanceV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeV1Client, err := config.computeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	var serverRequests []cloudservers.Server
	server := cloudservers.Server{
		Id: d.Id(),
	}
	serverRequests = append(serverRequests, server)

	deleteOpts := cloudservers.DeleteOpts{
		Servers:      serverRequests,
		DeleteVolume: d.Get("delete_disks_on_termination").(bool),
	}

	n, err := cloudservers.Delete(computeV1Client, deleteOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud server: %s", err)
	}

	if err := cloudservers.WaitForJobSuccess(computeV1Client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceInstanceNicsV1(d *schema.ResourceData) []cloudservers.Nic {
	var nicRequests []cloudservers.Nic

	nics := d.Get("nics").([]interface{})
	for i := range nics {
		nic := nics[i].(map[string]interface{})
		nicRequest := cloudservers.Nic{
			SubnetId:  nic["network_id"].(string),
			IpAddress: nic["ip_address"].(string),
		}

		nicRequests = append(nicRequests, nicRequest)
	}
	return nicRequests
}

func resourceInstanceRootVolumeV1(d *schema.ResourceData) cloudservers.RootVolume {
	volRequest := cloudservers.RootVolume{
		VolumeType: d.Get("system_disk_type").(string),
		Size:       d.Get("system_disk_size").(int),
	}
	return volRequest
}

func resourceInstanceDataVolumesV1(d *schema.ResourceData) []cloudservers.DataVolume {
	var volRequests []cloudservers.DataVolume

	vols := d.Get("data_disks").([]interface{})
	for i := range vols {
		vol := vols[i].(map[string]interface{})
		volRequest := cloudservers.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
		}
		if vol["snapshot_id"] != "" {
			extendparam := cloudservers.VolumeExtendParam{
				SnapshotId: vol["snapshot_id"].(string),
			}
			volRequest.Extendparam = &extendparam
		}

		volRequests = append(volRequests, volRequest)
	}
	return volRequests
}

func resourceInstanceSecGroupsV1(d *schema.ResourceData) []cloudservers.SecurityGroup {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]cloudservers.SecurityGroup, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secgroups[i] = cloudservers.SecurityGroup{
			ID: raw.(string),
		}
	}
	return secgroups
}

func flattenInstanceNicsV1(
	d *schema.ResourceData, meta interface{}, addresses map[string][]cloudservers.Address) []map[string]interface{} {

	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		log.Printf("Error creating HuaweiCloud networking client: %s", err)
	}

	var network string
	nics := []map[string]interface{}{}
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
				"network_id":  network,
				"ip_address":  addr.Addr,
				"mac_address": addr.MacAddr,
			}
			nics = append(nics, v)
		}
	}

	log.Printf("[DEBUG] flattenInstanceNicsV1: %#v", nics)
	return nics
}
