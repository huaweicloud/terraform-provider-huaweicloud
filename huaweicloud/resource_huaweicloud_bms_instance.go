package huaweicloud

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/bms/v1/baremetalservers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceBmsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBmsInstanceCreate,
		Read:   resourceBmsInstanceRead,
		Update: resourceBmsInstanceUpdate,
		Delete: resourceBmsInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
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
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"admin_pass": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
				ExactlyOneOf: []string{
					"admin_pass", "key_pair",
				},
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"admin_pass", "key_pair",
				},
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"iptype", "eip_charge_mode", "bandwidth_charge_mode", "bandwidth_size", "sharetype",
				},
			},
			"iptype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"5_bgp", "5_sbgp",
				}, true),
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"eip_charge_mode", "sharetype", "bandwidth_size",
				},
			},
			"eip_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, true),
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "sharetype", "bandwidth_size",
				},
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PER", "WHOLE",
				}, true),
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "eip_charge_mode", "bandwidth_size",
				},
			},
			"bandwidth_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "eip_charge_mode", "sharetype",
				},
			},
			"bandwidth_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"traffic", "bandwidth",
				}, true),
				ConflictsWith: []string{"eip_id"},
			},
			"system_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SAS", "SSD", "GPSSD", "ESSD",
				}, true),
				RequiredWith: []string{
					"system_disk_size",
				},
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"system_disk_type",
				},
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 59,
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
			"charging_mode": schemeChargingMode([]string{}),
			"period_unit":   schemaPeriodUnit([]string{}),
			"period":        schemaPeriod([]string{}),
			"auto_renew":    schemaAutoRenew([]string{}),

			"tags": tagsForceNewSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBmsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	bmsClient, err := config.BmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud bms client: %s", err)
	}

	createOpts := &baremetalservers.CreateOpts{
		Name:      d.Get("name").(string),
		ImageRef:  d.Get("image_id").(string),
		FlavorRef: d.Get("flavor_id").(string),
		MetaData: baremetalservers.MetaData{
			OpSvcUserId: d.Get("user_id").(string),
			AgencyName:  d.Get("agency_name").(string),
		},
		UserData:         []byte(d.Get("user_data").(string)),
		AdminPass:        d.Get("admin_pass").(string),
		KeyName:          d.Get("key_pair").(string),
		VpcId:            d.Get("vpc_id").(string),
		SecurityGroups:   resourceBmsInstanceSecGroupsV1(d),
		AvailabilityZone: d.Get("availability_zone").(string),
		Nics:             resourceBmsInstanceNicsV1(d),
		DataVolumes:      resourceBmsInstanceDataVolumesV1(d),
		ExtendParam: baremetalservers.ServerExtendParam{
			ChargingMode:        d.Get("charging_mode").(string),
			PeriodType:          d.Get("period_unit").(string),
			PeriodNum:           d.Get("period").(int),
			IsAutoPay:           "true",
			IsAutoRenew:         d.Get("auto_renew").(string),
			EnterpriseProjectId: GetEnterpriseProjectID(d, config),
		},
	}

	if hasFilledOpt(d, "eip_id") || hasFilledOpt(d, "iptype") {
		var publicIp baremetalservers.PublicIp
		if hasFilledOpt(d, "eip_id") {
			publicIp.Id = d.Get("eip_id").(string)
		} else {
			publicIp.Eip = &baremetalservers.Eip{
				IpType: d.Get("iptype").(string),
				BandWidth: baremetalservers.BandWidth{
					ShareType:  d.Get("sharetype").(string),
					Size:       d.Get("bandwidth_size").(int),
					ChargeMode: d.Get("bandwidth_charge_mode").(string),
				},
				ExtendParam: baremetalservers.EipExtendParam{
					ChargingMode: d.Get("eip_charge_mode").(string),
				},
			}
		}
		createOpts.PublicIp = &publicIp
	}

	if v, ok := d.GetOk("system_disk_type"); ok {
		volRequest := baremetalservers.RootVolume{
			VolumeType: v.(string),
			Size:       d.Get("system_disk_size").(int),
		}
		createOpts.RootVolume = &volRequest
	}

	if hasFilledOpt(d, "tags") {
		tagRaw := d.Get("tags").(map[string]interface{})
		taglist := utils.ExpandResourceTags(tagRaw)
		createOpts.ServerTags = taglist
	}

	n, err := baremetalservers.CreatePrePaid(bmsClient, createOpts).ExtractOrderResponse()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud BMS server: %s", err)
	}
	job_id := n.JobID

	if err := baremetalservers.WaitForJobSuccess(bmsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), job_id); err != nil {
		return err
	}

	entity, err := baremetalservers.GetJobEntity(bmsClient, job_id, "server_id")
	if err != nil {
		return err
	}
	server_id := entity.(string)

	if server_id != "" {
		d.SetId(server_id)
		return resourceBmsInstanceRead(d, meta)
	}

	return fmtp.Errorf("Unexpected conversion error in resourceBmsInstanceCreate")
}

func resourceBmsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	bmsClient, err := config.BmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	server, err := baremetalservers.Get(bmsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "server")
	}
	if server.Status == "DELETED" {
		d.SetId("")
		return nil
	}

	logp.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("region", GetRegion(d, config))
	d.Set("name", server.Name)
	d.Set("image_id", server.Image.ID)
	d.Set("flavor_id", server.Flavor.ID)
	d.Set("host_id", server.HostID)

	// Set fixed and floating ip
	if eip := bmsPublicIP(server); eip != "" {
		d.Set("public_ip", eip)
	}
	nics := flattenBmsInstanceNicsV1(d, meta, server.Addresses)
	d.Set("nics", nics)

	d.Set("key_pair", server.KeyName)
	// Set security groups
	secGrpIds := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpIds = append(secGrpIds, sg.ID)
	}
	d.Set("security_groups", secGrpIds)
	d.Set("status", server.Status)
	d.Set("user_id", server.Metadata.OpSvcUserId)
	d.Set("image_name", server.Metadata.ImageName)
	d.Set("vpc_id", server.Metadata.VpcID)
	d.Set("availability_zone", server.AvailabilityZone)
	d.Set("description", server.Description)
	d.Set("user_data", server.UserData)
	d.Set("tags", server.Tags)
	d.Set("enterprise_project_id", server.EnterpriseProjectID)
	// Set disk ids
	diskIds := []string{}
	for _, disk := range server.VolumeAttached {
		diskIds = append(diskIds, disk.ID)
	}
	d.Set("disk_ids", diskIds)
	return nil
}

func resourceBmsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	bmsClient, err := config.BmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	if d.HasChange("name") {
		var updateOpts baremetalservers.UpdateOpts
		updateOpts.Name = d.Get("name").(string)

		_, err = baremetalservers.Update(bmsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud bms server: %s", err)
		}
	}

	return resourceBmsInstanceRead(d, meta)
}

func resourceBmsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	bmsClient, err := config.BmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}
	serverID := d.Id()
	publicIP := d.Get("public_ip").(string)
	diskIds := d.Get("disk_ids").([]interface{})

	resourceIDs := make([]string, 0, 2+len(diskIds))

	resourceIDs = append(resourceIDs, serverID)

	if len(diskIds) > 0 {
		for _, diskId := range diskIds {
			resourceIDs = append(resourceIDs, diskId.(string))
		}
	}

	// unsubscribe the eip if necessary
	if _, ok := d.GetOk("iptype"); ok && publicIP != "" && d.Get("eip_charge_mode").(string) == "prePaid" {
		eipClient, err := config.NetworkingV1Client(GetRegion(d, config))
		if err != nil {
			return fmtp.Errorf("Error creating networking client: %s", err)
		}

		if eipID, err := getEipIDbyAddress(eipClient, publicIP); err == nil {
			resourceIDs = append(resourceIDs, eipID)
		} else {
			return fmtp.Errorf("Error fetching EIP ID of BMS server (%s): %s", d.Id(), err)
		}
	}

	if err := UnsubscribePrePaidResource(d, config, resourceIDs); err != nil {
		return fmtp.Errorf("Error unsubscribing HuaweiCloud BMS server: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting", "ACTIVE", "SHUTOFF"},
		Target:       []string{"DELETED"},
		Refresh:      waitForBmsInstanceDelete(bmsClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud BMS instance: %s", err)
	}

	d.SetId("")
	return nil
}
func resourceBmsInstanceNicsV1(d *schema.ResourceData) []baremetalservers.Nic {
	var nicRequests []baremetalservers.Nic

	nics := d.Get("nics").([]interface{})
	for i := range nics {
		nic := nics[i].(map[string]interface{})
		nicRequest := baremetalservers.Nic{
			SubnetId:  nic["subnet_id"].(string),
			IpAddress: nic["ip_address"].(string),
		}

		nicRequests = append(nicRequests, nicRequest)
	}
	return nicRequests
}

func resourceBmsInstanceDataVolumesV1(d *schema.ResourceData) []baremetalservers.DataVolume {
	var volRequests []baremetalservers.DataVolume

	vols := d.Get("data_disks").([]interface{})
	for i := range vols {
		vol := vols[i].(map[string]interface{})
		volRequest := baremetalservers.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
		}
		volRequests = append(volRequests, volRequest)
	}
	return volRequests
}

func resourceBmsInstanceSecGroupsV1(d *schema.ResourceData) []baremetalservers.SecurityGroup {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]baremetalservers.SecurityGroup, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secgroups[i] = baremetalservers.SecurityGroup{
			ID: raw.(string),
		}
	}
	return secgroups
}

func flattenBmsInstanceNicsV1(d *schema.ResourceData, meta interface{},
	addresses map[string][]baremetalservers.Address) []map[string]interface{} {

	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
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
				"subnet_id":   network,
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

func bmsPublicIP(server *baremetalservers.CloudServer) string {
	var publicIP string

	for _, addresses := range server.Addresses {
		for _, addr := range addresses {
			if addr.Type == "floating" {
				publicIP = addr.Addr
				break
			}
		}
	}

	return publicIP
}

func waitForBmsInstanceDelete(bmsClient *golangsdk.ServiceClient, ServerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud BMS instance %s", ServerId)

		r, err := baremetalservers.Get(bmsClient, ServerId).Extract()

		if err != nil {
			return r, "Deleting", err
		}

		return r, r.Status, nil
	}
}
