package deprecated

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/secgroups"
	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceEcsInstanceV1() *schema.Resource {
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

		CustomizeDiff: config.MergeDefaultTags(),

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
				StateFunc: utils.HashAndHexEncode,
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
			"tags": common.TagsSchema(),
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
	cfg := meta.(*config.Config)
	computeClient, err := cfg.ComputeV11Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute V1.1 client: %s", err)
	}
	computeV1Client, err := cfg.ComputeV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute V1 client: %s", err)
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
	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		extendParam.EnterpriseProjectId = epsID
	}
	if extendParam != (cloudservers.ServerExtendParam{}) {
		createOpts.ExtendParam = &extendParam
	}

	var metadata cloudservers.MetaData
	if common.HasFilledOpt(d, "op_svc_userid") {
		metadata.OpSvcUserId = d.Get("op_svc_userid").(string)
	}
	if metadata != (cloudservers.MetaData{}) {
		createOpts.MetaData = &metadata
	}

	if tags, ok := d.GetOk("tags"); ok {
		createOpts.ServerTags = utils.ExpandResourceTags(tags.(map[string]interface{}))
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.AdminPass = d.Get("password").(string)

	var instance_id string
	if d.Get("charging_mode") == "prePaid" {
		bssV1Client, err := cfg.BssV1Client(cfg.GetRegion(d))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud bss V1 client: %s", err)
		}
		n, err := cloudservers.CreatePrePaid(computeClient, createOpts).ExtractOrderResponse()
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud server: %s", err)
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
			return fmtp.Errorf("Error creating HuaweiCloud server: %s", err)
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

		if common.HasFilledOpt(d, "auto_recovery") {
			ar := d.Get("auto_recovery").(bool)
			logp.Printf("[DEBUG] Set auto recovery of instance to %t", ar)
			err = setAutoRecoveryForInstance(d, meta, instance_id, ar)
			if err != nil {
				logp.Printf("[WARN] Error setting auto recovery of instance:%s, err=%s", instance_id, err)
			}
		}

		return resourceEcsInstanceV1Read(d, meta)
	}

	return fmtp.Errorf("Unexpected conversion error in resourceEcsInstanceV1Create.")
}

func resourceEcsInstanceV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	server, err := cloudservers.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "server")
	}
	if server.Status == "DELETED" {
		d.SetId("")
		return nil
	}

	logp.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)
	d.Set("image_id", server.Image.ID)
	d.Set("flavor", server.Flavor.ID)
	d.Set("key_name", server.KeyName)
	d.Set("vpc_id", server.Metadata.VpcID)
	d.Set("availability_zone", server.AvailabilityZone)

	// set secgroups
	secGrpIDs := make([]string, len(server.SecurityGroups))
	for i, sg := range server.SecurityGroups {
		secGrpIDs[i] = sg.ID
	}
	d.Set("security_groups", secGrpIDs)

	// Get the instance network and address information
	nics := flattenInstanceNicsV1(d, meta, server.Addresses)
	d.Set("nics", nics)

	// Set instance tags
	d.Set("tags", flattenTagsToMap(server.Tags))

	ar, err := resourceECSAutoRecoveryV1Read(d, meta, d.Id())
	if err != nil && !utils.IsResourceNotFound(err) {
		return fmtp.Errorf("Error reading auto recovery of instance:%s, err=%s", d.Id(), err)
	}
	d.Set("auto_recovery", ar)

	return nil
}

func resourceEcsInstanceV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeV2Client, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute v2.1 client: %s", err)
	}

	computeV1Client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute v1 client: %s", err)
	}

	var updateOpts servers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if updateOpts != (servers.UpdateOpts{}) {
		_, err := servers.Update(computeV2Client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud server: %s", err)
		}
	}

	if d.HasChange("security_groups") {
		oldSGRaw, newSGRaw := d.GetChange("security_groups")
		oldSGSet := oldSGRaw.(*schema.Set)
		newSGSet := newSGRaw.(*schema.Set)
		secgroupsToAdd := newSGSet.Difference(oldSGSet)
		secgroupsToRemove := oldSGSet.Difference(newSGSet)

		logp.Printf("[DEBUG] Security groups to add: %v", secgroupsToAdd)

		logp.Printf("[DEBUG] Security groups to remove: %v", secgroupsToRemove)

		for _, g := range secgroupsToRemove.List() {
			err := secgroups.RemoveServer(computeV2Client, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					continue
				}

				return fmtp.Errorf("Error removing security group (%s) from HuaweiCloud server (%s): %s", g, d.Id(), err)
			} else {
				logp.Printf("[DEBUG] Removed security group (%s) from instance (%s)", g, d.Id())
			}
		}

		for _, g := range secgroupsToAdd.List() {
			err := secgroups.AddServer(computeV2Client, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return fmtp.Errorf("Error adding security group (%s) to HuaweiCloud server (%s): %s", g, d.Id(), err)
			}
			logp.Printf("[DEBUG] Added security group (%s) to instance (%s)", g, d.Id())
		}
	}

	if d.HasChange("flavor") {
		newFlavorId := d.Get("flavor").(string)

		resizeOpts := &servers.ResizeOpts{
			FlavorRef: newFlavorId,
		}
		logp.Printf("[DEBUG] Resize configuration: %#v", resizeOpts)
		err := servers.Resize(computeV2Client, d.Id(), resizeOpts).ExtractErr()
		if err != nil {
			return fmtp.Errorf("Error resizing HuaweiCloud server: %s", err)
		}

		// Wait for the instance to finish resizing.
		logp.Printf("[DEBUG] Waiting for instance (%s) to finish resizing", d.Id())

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"RESIZE"},
			Target:     []string{"VERIFY_RESIZE"},
			Refresh:    ServerV1StateRefreshFunc(computeV1Client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmtp.Errorf("Error waiting for instance (%s) to resize: %s", d.Id(), err)
		}

		// Confirm resize.
		logp.Printf("[DEBUG] Confirming resize")
		err = servers.ConfirmResize(computeV2Client, d.Id()).ExtractErr()
		if err != nil {
			return fmtp.Errorf("Error confirming resize of HuaweiCloud server: %s", err)
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"VERIFY_RESIZE"},
			Target:     []string{"ACTIVE"},
			Refresh:    ServerV1StateRefreshFunc(computeV1Client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmtp.Errorf("Error waiting for instance (%s) to confirm resize: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		ecsClient, err := config.ComputeV1Client(config.GetRegion(d))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud compute v1 client: %s", err)
		}

		tagErr := utils.UpdateResourceTags(ecsClient, d, "cloudservers", d.Id())
		if tagErr != nil {
			return fmtp.Errorf("Error updating tags of instance:%s, err:%s", d.Id(), err)
		}
	}

	if d.HasChange("auto_recovery") {
		ar := d.Get("auto_recovery").(bool)
		logp.Printf("[DEBUG] Update auto recovery of instance to %t", ar)
		err = setAutoRecoveryForInstance(d, meta, d.Id(), ar)
		if err != nil {
			return fmtp.Errorf("Error updating auto recovery of instance:%s, err:%s", d.Id(), err)
		}
	}

	return resourceEcsInstanceV1Read(d, meta)
}

func resourceEcsInstanceV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeV1Client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	if d.Get("charging_mode") == "prePaid" {
		bssV1Client, err := config.BssV1Client(config.GetRegion(d))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud bss V1 client: %s", err)
		}

		resourceIds := []string{d.Id()}
		deleteOrderOpts := cloudservers.DeleteOrderOpts{
			ResourceIds: resourceIds,
			UnSubType:   1,
		}
		n, err := cloudservers.DeleteOrder(bssV1Client, deleteOrderOpts).ExtractDeleteOrderResponse()
		if err != nil {
			return fmtp.Errorf("Error deleting HuaweiCloud server: %s", err)
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
			return fmtp.Errorf("Error deleting HuaweiCloud server: %s", err)
		}

		if err := cloudservers.WaitForJobSuccess(computeV1Client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
			return err
		}
	}

	// Wait for the instance to delete before moving on.
	logp.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTOFF"},
		Target:     []string{"DELETED", "SOFT_DELETED"},
		Refresh:    ServerV1StateRefreshFunc(computeV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
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

	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		logp.Printf("Error creating HuaweiCloud networking client: %s", err)
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
				logp.Printf("[DEBUG] flattenInstanceNicsV1: failed to fetch port %s", addr.PortID)
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

	logp.Printf("[DEBUG] flattenInstanceNicsV1: %#v", nics)
	return nics
}

// ServerV1StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch an HuaweiCloud instance.
func ServerV1StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := cloudservers.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		// get fault message when status is ERROR
		if s.Status == "ERROR" {
			fault := fmtp.Errorf("[error code: %d, message: %s]", s.Fault.Code, s.Fault.Message)
			return s, "ERROR", fault
		}
		return s, s.Status, nil
	}
}

func flattenTagsToMap(tags []string) map[string]string {
	result := make(map[string]string)

	for _, tagStr := range tags {
		tag := strings.SplitN(tagStr, "=", 2)
		if len(tag) == 1 {
			result[tag[0]] = ""
		} else if len(tag) == 2 {
			result[tag[0]] = tag[1]
		}
	}

	return result
}
