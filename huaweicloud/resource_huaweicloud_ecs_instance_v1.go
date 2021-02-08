package huaweicloud

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/secgroups"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func resourceEcsInstanceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEcsInstanceV1Create,
		Read:   resourceEcsInstanceV1Read,
		Update: resourceEcsInstanceV1Update,
		Delete: resourceEcsInstanceV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_compute_instance resource instead",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
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
						"port_id": {
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
				Optional: true,
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
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, true),
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, true),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: validateECSTagValue,
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"auto_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"delete_disks_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"op_svc_userid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEcsInstanceV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV11Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute V1.1 client: %s", err)
	}
	computeV1Client, err := config.ComputeV1Client(GetRegion(d, config))
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

	var extendParam cloudservers.ServerExtendParam
	if d.Get("charging_mode") == "prePaid" {
		extendParam.ChargingMode = d.Get("charging_mode").(string)
		extendParam.PeriodType = d.Get("period_unit").(string)
		extendParam.PeriodNum = d.Get("period").(int)
		extendParam.IsAutoPay = "true"
		extendParam.IsAutoRenew = d.Get("auto_renew").(string)
	}
	epsID := GetEnterpriseProjectID(d, config)
	if epsID != "" {
		extendParam.EnterpriseProjectId = epsID
	}
	if extendParam != (cloudservers.ServerExtendParam{}) {
		createOpts.ExtendParam = &extendParam
	}

	var metadata cloudservers.MetaData
	if hasFilledOpt(d, "op_svc_userid") {
		metadata.OpSvcUserId = d.Get("op_svc_userid").(string)
	}
	if metadata != (cloudservers.MetaData{}) {
		createOpts.MetaData = &metadata
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	var instance_id string
	if d.Get("charging_mode") == "prePaid" {
		bssV1Client, err := config.BssV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud bss V1 client: %s", err)
		}
		n, err := cloudservers.CreatePrePaid(computeClient, createOpts).ExtractOrderResponse()
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud server: %s", err)
		}

		if err := cloudservers.WaitForOrderSuccess(bssV1Client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.OrderID); err != nil {
			return err
		}

		resource, err := cloudservers.GetOrderResource(bssV1Client, n.OrderID)
		if err != nil {
			return err
		}
		instance_id = resource.(string)
	} else {
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
		instance_id = entity.(string)
	}

	if instance_id != "" {
		d.SetId(instance_id)

		if hasFilledOpt(d, "tags") {
			tagRaw := d.Get("tags").(map[string]interface{})
			taglist := expandResourceTags(tagRaw)
			tagErr := tags.Create(computeV1Client, "cloudservers", instance_id, taglist).ExtractErr()
			if tagErr != nil {
				log.Printf("[WARN] Error setting tags of instance:%s, err=%s", instance_id, err)
			}
		}

		if hasFilledOpt(d, "auto_recovery") {
			ar := d.Get("auto_recovery").(bool)
			log.Printf("[DEBUG] Set auto recovery of instance to %t", ar)
			err = setAutoRecoveryForInstance(d, meta, instance_id, ar)
			if err != nil {
				log.Printf("[WARN] Error setting auto recovery of instance:%s, err=%s", instance_id, err)
			}
		}

		return resourceEcsInstanceV1Read(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceEcsInstanceV1Create.")
}

func resourceEcsInstanceV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	server, err := cloudservers.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "server")
	}
	if server.Status == "DELETED" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)
	d.Set("image_id", server.Image.ID)
	d.Set("flavor", server.Flavor.ID)
	d.Set("key_name", server.KeyName)
	d.Set("vpc_id", server.Metadata.VpcID)
	d.Set("availability_zone", server.AvailabilityZone)

	// Get the instance network and address information
	nics := flattenInstanceNicsV1(d, meta, server.Addresses)
	d.Set("nics", nics)

	// Set instance tags
	if resourceTags, err := tags.Get(computeClient, "cloudservers", d.Id()).Extract(); err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmt.Errorf("Error saving tags to state for ECS instance (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] Error fetching tags of ECS instance (%s): %s", d.Id(), err)
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
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
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
		ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute v1 client: %s", err)
		}

		tagErr := UpdateResourceTags(ecsClient, d, "cloudservers", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of instance:%s, err:%s", d.Id(), err)
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
	computeV1Client, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}
	computeV2Client, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute V2 client: %s", err)
	}

	if d.Get("charging_mode") == "prePaid" {
		bssV1Client, err := config.BssV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud bss V1 client: %s", err)
		}

		resourceIds := []string{d.Id()}
		deleteOrderOpts := cloudservers.DeleteOrderOpts{
			ResourceIds: resourceIds,
			UnSubType:   1,
		}
		n, err := cloudservers.DeleteOrder(bssV1Client, deleteOrderOpts).ExtractDeleteOrderResponse()
		if err != nil {
			return fmt.Errorf("Error deleting HuaweiCloud server: %s", err)
		}

		if err := cloudservers.WaitForOrderDeleteSuccess(bssV1Client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.OrderIDs[0]); err != nil {
			return err
		}
	} else {
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
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTOFF"},
		Target:     []string{"DELETED", "SOFT_DELETED"},
		Refresh:    ServerV2StateRefreshFunc(computeV2Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to delete: %s",
			d.Id(), err)
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
	disk_type := d.Get("system_disk_type").(string)
	if disk_type == "" {
		disk_type = "GPSSD"
	}
	volRequest := cloudservers.RootVolume{
		VolumeType: disk_type,
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
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
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
				"port_id":     addr.PortID,
			}
			nics = append(nics, v)
		}
	}

	log.Printf("[DEBUG] flattenInstanceNicsV1: %#v", nics)
	return nics
}
