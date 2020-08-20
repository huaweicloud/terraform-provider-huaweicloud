package huaweicloud

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/blockstorage/v2/volumes"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/availabilityzones"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/keypairs"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/schedulerhints"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/secgroups"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/startstop"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/flavors"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/images"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
)

func resourceComputeInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceV2Create,
		Read:   resourceComputeInstanceV2Read,
		Update: resourceComputeInstanceV2Update,
		Delete: resourceComputeInstanceV2Delete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeInstanceV2ImportState,
		},

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
				ForceNew: false,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_ID", nil),
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_NAME", nil),
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
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v6": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_network": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"metadata": {
				Type:          schema.TypeMap,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"system_disk_type", "system_disk_size", "data_disks"},
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ForceNew:  false,
			},
			"access_ip_v4": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"access_ip_v6": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"block_device": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"system_disk_type", "system_disk_size", "data_disks"},
				Deprecated:    "use system_disk_type, system_disk_size, data_disks instead",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"destination_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"boot_index": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"guest_format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"system_disk_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"block_device", "metadata"},
				ValidateFunc: validation.StringInSlice([]string{
					"SATA", "SAS", "SSD", "co-p1", "uh-l1",
				}, true),
			},
			"system_disk_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"block_device", "metadata"},
			},
			"data_disks": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"block_device", "metadata"},
				MinItems:      1,
				MaxItems:      23,
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
			"scheduler_hints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"tenancy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"deh_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceComputeSchedulerHintsHash,
			},
			"stop_before_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enterprise_project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"block_device", "metadata"},
			},
			"delete_disks_on_termination": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"block_device", "metadata"},
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, true),
				ConflictsWith: []string{"block_device", "metadata"},
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, true),
				ConflictsWith: []string{"block_device", "metadata"},
			},
			"period": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"block_device", "metadata"},
			},
			"auto_renew": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"block_device", "metadata"},
			},
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: validateECSTagValue,
			},
			"all_metadata": {
				Type:     schema.TypeMap,
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
						"pci_address": {
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
					},
				},
			},
		},
	}
}

func resourceComputeInstanceV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	computeV1Client, err := config.computeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	// Determines the Image ID using the following rules:
	// If a bootable block_device was specified, ignore the image altogether.
	// If an image_id was specified, use it.
	// If an image_name was specified, look up the image ID, report if error.
	imageId, err := getImageIDFromConfig(computeClient, d)
	if err != nil {
		return err
	}

	flavorId, err := getFlavorID(computeClient, d)
	if err != nil {
		return err
	}

	// determine if block_device configuration is correct
	// this includes valid combinations and required attributes
	if err := checkBlockDeviceConfig(d); err != nil {
		return err
	}

	// Try to call API of Huawei ECS instead of OpenStack
	if !hasFilledOpt(d, "block_device") && !hasFilledOpt(d, "metadata") {
		client, err := config.computeV11Client(GetRegion(d, config))
		vpcClient, err := config.networkingV1Client(GetRegion(d, config))
		sgClient, err := config.networkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud Client: %s", err)
		}

		vpcId, err := getVpcID(vpcClient, d)
		if err != nil {
			return err
		}

		secGroups, err := resourceInstanceSecGroupIdsV1(sgClient, d)
		if err != nil {
			return err
		}

		createOpts := &cloudservers.CreateOpts{
			Name:             d.Get("name").(string),
			ImageRef:         imageId,
			FlavorRef:        flavorId,
			KeyName:          d.Get("key_pair").(string),
			VpcId:            vpcId,
			SecurityGroups:   secGroups,
			AvailabilityZone: d.Get("availability_zone").(string),
			Nics:             resourceInstanceNicsV2(d),
			RootVolume:       resourceInstanceRootVolumeV1(d),
			DataVolumes:      resourceInstanceDataVolumesV1(d),
			AdminPass:        d.Get("admin_pass").(string),
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
		if hasFilledOpt(d, "enterprise_project_id") {
			extendParam.EnterpriseProjectId = d.Get("enterprise_project_id").(string)
		}
		if extendParam != (cloudservers.ServerExtendParam{}) {
			createOpts.ExtendParam = &extendParam
		}

		schedulerHintsRaw := d.Get("scheduler_hints").(*schema.Set).List()
		if len(schedulerHintsRaw) > 0 {
			log.Printf("[DEBUG] schedulerhints: %+v", schedulerHintsRaw)
			schedulerHints := resourceInstanceSchedulerHintsV1(d, schedulerHintsRaw[0].(map[string]interface{}))
			createOpts.SchedulerHints = &schedulerHints
		}

		log.Printf("[DEBUG] Create Options: %#v", createOpts)

		var server_id string
		if d.Get("charging_mode") == "prePaid" {
			// prePaid.
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
			server_id = resource.(string)
		} else {
			// postPaid.
			n, err := cloudservers.Create(client, createOpts).ExtractJobResponse()
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
			server_id = entity.(string)
		}

		// Store the ID now
		d.SetId(server_id)

	} else {
		// OpenStack API implementation. Clean up this after removing block_device.

		// Build a list of networks with the information given upon creation.
		// Error out if an invalid network configuration was used.
		allInstanceNetworks, err := getAllInstanceNetworks(d, meta)
		if err != nil {
			return err
		}

		// Build a []servers.Network to pass into the create options.
		networks := expandInstanceNetworks(allInstanceNetworks)

		var createOpts servers.CreateOptsBuilder
		createOpts = &servers.CreateOpts{
			Name:             d.Get("name").(string),
			ImageRef:         imageId,
			FlavorRef:        flavorId,
			SecurityGroups:   resourceInstanceSecGroupsV2(d),
			AvailabilityZone: d.Get("availability_zone").(string),
			Networks:         networks,
			Metadata:         resourceInstanceMetadataV2(d),
			AdminPass:        d.Get("admin_pass").(string),
			UserData:         []byte(d.Get("user_data").(string)),
		}

		if keyName, ok := d.Get("key_pair").(string); ok && keyName != "" {
			createOpts = &keypairs.CreateOptsExt{
				CreateOptsBuilder: createOpts,
				KeyName:           keyName,
			}
		}

		if vL, ok := d.GetOk("block_device"); ok {
			blockDevices, err := resourceInstanceBlockDevicesV2(d, vL.([]interface{}))
			if err != nil {
				return err
			}

			createOpts = &bootfromvolume.CreateOptsExt{
				CreateOptsBuilder: createOpts,
				BlockDevice:       blockDevices,
			}
		}

		schedulerHintsRaw := d.Get("scheduler_hints").(*schema.Set).List()
		if len(schedulerHintsRaw) > 0 {
			log.Printf("[DEBUG] schedulerhints: %+v", schedulerHintsRaw)
			schedulerHints := resourceInstanceSchedulerHintsV2(d, schedulerHintsRaw[0].(map[string]interface{}))
			createOpts = &schedulerhints.CreateOptsExt{
				CreateOptsBuilder: createOpts,
				SchedulerHints:    schedulerHints,
			}
		}

		log.Printf("[DEBUG] Create Options: %#v", createOpts)

		// If a block_device is used, use the bootfromvolume.Create function as it allows an empty ImageRef.
		// Otherwise, use the normal servers.Create function.
		var server *servers.Server
		if _, ok := d.GetOk("block_device"); ok {
			server, err = bootfromvolume.Create(computeClient, createOpts).Extract()
		} else {
			server, err = servers.Create(computeClient, createOpts).Extract()
		}

		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud server: %s", err)
		}

		log.Printf("[INFO] Instance ID: %s", server.ID)

		// Store the ID now
		d.SetId(server.ID)

		// Wait for the instance to become running so we can get some attributes
		// that aren't available until later.
		log.Printf(
			"[DEBUG] Waiting for instance (%s) to become running",
			server.ID)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"BUILD"},
			Target:     []string{"ACTIVE"},
			Refresh:    ServerV2StateRefreshFunc(computeClient, server.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				server.ID, err)
		}
	}

	// Set tags
	if hasFilledOpt(d, "tags") {
		tagRaw := d.Get("tags").(map[string]interface{})
		taglist := expandResourceTags(tagRaw)
		tagErr := tags.Create(computeV1Client, "cloudservers", d.Id(), taglist).ExtractErr()
		if tagErr != nil {
			log.Printf("[WARN] Error setting tags of instance:%s, err=%s", d.Id(), err)
		}
	}

	return resourceComputeInstanceV2Read(d, meta)
}

func resourceComputeInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	computeV1Client, err := config.computeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	server, err := servers.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "server")
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)

	// Get the instance network and address information
	networks, err := flattenInstanceNetworks(d, meta)
	if err != nil {
		return err
	}

	// Determine the best IPv4 and IPv6 addresses to access the instance with
	hostv4, hostv6 := getInstanceAccessAddresses(d, networks)

	// AccessIPv4/v6 isn't standard in HuaweiCloud, but there have been reports
	// of them being used in some environments.
	if server.AccessIPv4 != "" && hostv4 == "" {
		hostv4 = server.AccessIPv4
	}

	if server.AccessIPv6 != "" && hostv6 == "" {
		hostv6 = server.AccessIPv6
	}

	d.Set("network", networks)
	d.Set("access_ip_v4", hostv4)
	d.Set("access_ip_v6", hostv6)

	// Determine the best IP address to use for SSH connectivity.
	// Prefer IPv4 over IPv6.
	var preferredSSHAddress string
	if hostv4 != "" {
		preferredSSHAddress = hostv4
	} else if hostv6 != "" {
		preferredSSHAddress = hostv6
	}

	if preferredSSHAddress != "" {
		// Initialize the connection info
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": preferredSSHAddress,
		})
	}

	d.Set("all_metadata", server.Metadata)

	flavorId, ok := server.Flavor["id"].(string)
	if !ok {
		return fmt.Errorf("Error setting HuaweiCloud server's flavor: %v", server.Flavor)
	}
	d.Set("flavor_id", flavorId)

	flavor, err := flavors.Get(computeClient, flavorId).Extract()
	if err != nil {
		return err
	}
	d.Set("flavor_name", flavor.Name)

	// Set volume attached
	bds := []map[string]interface{}{}
	if len(server.VolumesAttached) > 0 {
		for _, b := range server.VolumesAttached {
			va, err := block_devices.Get(computeV1Client, d.Id(), b["id"]).Extract()
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Retrieved volume attachment %s: %#v", d.Id(), va)
			v := map[string]interface{}{
				"pci_address": va.PciAddress,
				"volume_id":   b["id"],
				"boot_index":  va.BootIndex,
				"size":        va.Size,
			}
			bds = append(bds, v)
		}
		d.Set("volume_attached", bds)
	}

	// Set the instance's image information appropriately
	if err := setImageInformation(computeClient, server, d); err != nil {
		return err
	}

	// Build a custom struct for the availability zone extension
	var serverWithAZ struct {
		servers.Server
		availabilityzones.ServerAvailabilityZoneExt
	}

	// Do another Get so the above work is not disturbed.
	err = servers.Get(computeClient, d.Id()).ExtractInto(&serverWithAZ)
	if err != nil {
		return CheckDeleted(d, err, "server")
	}

	// Set the availability zone
	d.Set("availability_zone", serverWithAZ.AvailabilityZone)

	// Set the region
	d.Set("region", GetRegion(d, config))

	// Set instance tags
	resourceTags, err := tags.Get(computeV1Client, "cloudservers", d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching HuaweiCloud instance tags: %s", err)
	}

	tagmap := tagsToMap(resourceTags.Tags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("[DEBUG] Error saving tag to state for HuaweiCloud instance (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceComputeInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
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

	if d.HasChange("metadata") {
		oldMetadata, newMetadata := d.GetChange("metadata")
		var metadataToDelete []string

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey := range oldMetadata.(map[string]interface{}) {
			var found bool
			for newKey := range newMetadata.(map[string]interface{}) {
				if oldKey == newKey {
					found = true
				}
			}

			if !found {
				metadataToDelete = append(metadataToDelete, oldKey)
			}
		}

		for _, key := range metadataToDelete {
			err := servers.DeleteMetadatum(computeClient, d.Id(), key).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error deleting metadata (%s) from server (%s): %s", key, d.Id(), err)
			}
		}

		// Update existing metadata and add any new metadata.
		metadataOpts := make(servers.MetadataOpts)
		for k, v := range newMetadata.(map[string]interface{}) {
			metadataOpts[k] = v.(string)
		}

		_, err := servers.UpdateMetadata(computeClient, d.Id(), metadataOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud server (%s) metadata: %s", d.Id(), err)
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

	if d.HasChange("admin_pass") {
		if newPwd, ok := d.Get("admin_pass").(string); ok {
			err := servers.ChangeAdminPassword(computeClient, d.Id(), newPwd).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error changing admin password of HuaweiCloud server (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("flavor_id") || d.HasChange("flavor_name") {
		var newFlavorId string
		var err error
		if d.HasChange("flavor_id") {
			newFlavorId = d.Get("flavor_id").(string)
		} else {
			newFlavorName := d.Get("flavor_name").(string)
			newFlavorId, err = flavors.IDFromName(computeClient, newFlavorName)
			if err != nil {
				return err
			}
		}

		resizeOpts := &servers.ResizeOpts{
			FlavorRef: newFlavorId,
		}
		log.Printf("[DEBUG] Resize configuration: %#v", resizeOpts)
		err = servers.Resize(computeClient, d.Id(), resizeOpts).ExtractErr()
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
		ecsClient, err := config.computeV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute v1 client: %s", err)
		}

		tagErr := UpdateResourceTags(ecsClient, d, "cloudservers")
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of instance:%s, err:%s", d.Id(), err)
		}
	}

	return resourceComputeInstanceV2Read(d, meta)
}

func resourceComputeInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeV1Client, err := config.computeV1Client(GetRegion(d, config))
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	if d.Get("stop_before_destroy").(bool) {
		err = startstop.Stop(computeClient, d.Id()).ExtractErr()
		if err != nil {
			log.Printf("[WARN] Error stopping HuaweiCloud instance: %s", err)
		} else {
			stopStateConf := &resource.StateChangeConf{
				Pending:    []string{"ACTIVE"},
				Target:     []string{"SHUTOFF"},
				Refresh:    ServerV2StateRefreshFunc(computeClient, d.Id()),
				Timeout:    3 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			log.Printf("[DEBUG] Waiting for instance (%s) to stop", d.Id())
			_, err = stopStateConf.WaitForState()
			if err != nil {
				log.Printf("[WARN] Error waiting for instance (%s) to stop: %s, proceeding to delete", d.Id(), err)
			}
		}
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

	d.SetId("")
	return nil
}

func resourceComputeInstanceV2ImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	var serverWithAttachments struct {
		VolumesAttached []map[string]interface{} `json:"os-extended-volumes:volumes_attached"`
	}

	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	results := make([]*schema.ResourceData, 1)
	err = resourceComputeInstanceV2Read(d, meta)
	if err != nil {
		return nil, fmt.Errorf("Error reading huaweicloud_compute_instance_v2 %s: %s", d.Id(), err)
	}

	raw := servers.Get(computeClient, d.Id())
	if raw.Err != nil {
		return nil, CheckDeleted(d, raw.Err, "huaweicloud_compute_instance_v2")
	}

	err = raw.ExtractInto(&serverWithAttachments)

	if err != nil {
		return nil, fmt.Errorf("Error reading attached volumes: %s", err)
	}

	log.Printf("[DEBUG] Retrieved huaweicloud_compute_instance_v2 %s volume attachments: %#v",
		d.Id(), serverWithAttachments)

	bds := []map[string]interface{}{}
	if len(serverWithAttachments.VolumesAttached) > 0 {
		blockStorageClient, err := config.blockStorageV2Client(GetRegion(d, config))
		if err != nil {
			return nil, fmt.Errorf("Error creating HuaweiCloud volume client: %s", err)
		}

		var volMetaData = struct {
			VolumeImageMetadata map[string]interface{} `json:"volume_image_metadata"`
			Id                  string                 `json:"id"`
			Size                int                    `json:"size"`
			Bootable            string                 `json:"bootable"`
		}{}
		for i, b := range serverWithAttachments.VolumesAttached {
			rawVolume := volumes.Get(blockStorageClient, b["id"].(string))
			err = rawVolume.ExtractInto(&volMetaData)
			if err != nil {
				return nil, fmt.Errorf("Error reading metadata from volume %s: %s", b["id"], err)
			}

			log.Printf("[DEBUG] retrieved volume%+v", volMetaData)
			v := map[string]interface{}{
				"delete_on_termination": true,
				"uuid":                  volMetaData.VolumeImageMetadata["image_id"],
				"boot_index":            i,
				"destination_type":      "volume",
				"source_type":           "image",
				"volume_size":           volMetaData.Size,
				"disk_bus":              "",
				"volume_type":           "",
				"device_type":           "",
			}

			if volMetaData.Bootable == "true" {
				bds = append(bds, v)
			}
		}

		d.Set("block_device", bds)
	}
	metadata, err := servers.Metadata(computeClient, d.Id()).Extract()
	d.Set("metadata", metadata)
	results[0] = d

	return results, nil
}

// ServerV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an HuaweiCloud instance.
func ServerV2StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := servers.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		// get fault message when status is ERROR
		if s.Status == "ERROR" {
			fault := fmt.Errorf("[error code: %d, message: %s]", s.Fault.Code, s.Fault.Message)
			return s, "ERROR", fault
		}
		return s, s.Status, nil
	}
}

func resourceInstanceSecGroupIdsV1(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]cloudservers.SecurityGroup, error) {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]cloudservers.SecurityGroup, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secName := raw.(string)
		secId, err := groups.IDFromName(client, secName)
		if err != nil {
			return secgroups, err
		}
		secgroups[i] = cloudservers.SecurityGroup{
			ID: secId,
		}
	}
	return secgroups, nil
}

func resourceInstanceSecGroupsV2(d *schema.ResourceData) []string {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]string, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secgroups[i] = raw.(string)
	}
	return secgroups
}

func resourceInstanceMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceInstanceNicsV2(d *schema.ResourceData) []cloudservers.Nic {
	var nicRequests []cloudservers.Nic

	networks := d.Get("network").([]interface{})
	for _, v := range networks {
		network := v.(map[string]interface{})
		nicRequest := cloudservers.Nic{
			SubnetId:  network["uuid"].(string),
			IpAddress: network["fixed_ip_v4"].(string),
		}

		nicRequests = append(nicRequests, nicRequest)
	}
	return nicRequests
}

func resourceInstanceBlockDevicesV2(d *schema.ResourceData, bds []interface{}) ([]bootfromvolume.BlockDevice, error) {
	blockDeviceOpts := make([]bootfromvolume.BlockDevice, len(bds))
	for i, bd := range bds {
		bdM := bd.(map[string]interface{})
		blockDeviceOpts[i] = bootfromvolume.BlockDevice{
			UUID:                bdM["uuid"].(string),
			VolumeSize:          bdM["volume_size"].(int),
			BootIndex:           bdM["boot_index"].(int),
			DeleteOnTermination: bdM["delete_on_termination"].(bool),
			GuestFormat:         bdM["guest_format"].(string),
		}

		sourceType := bdM["source_type"].(string)
		switch sourceType {
		case "blank":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceBlank
		case "image":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceImage
		case "snapshot":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceSnapshot
		case "volume":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceVolume
		default:
			return blockDeviceOpts, fmt.Errorf("unknown block device source type %s", sourceType)
		}

		destinationType := bdM["destination_type"].(string)
		switch destinationType {
		case "local":
			blockDeviceOpts[i].DestinationType = bootfromvolume.DestinationLocal
		case "volume":
			blockDeviceOpts[i].DestinationType = bootfromvolume.DestinationVolume
		default:
			return blockDeviceOpts, fmt.Errorf("unknown block device destination type %s", destinationType)
		}
	}

	log.Printf("[DEBUG] Block Device Options: %+v", blockDeviceOpts)
	return blockDeviceOpts, nil
}

func resourceInstanceSchedulerHintsV1(d *schema.ResourceData, schedulerHintsRaw map[string]interface{}) cloudservers.SchedulerHints {
	schedulerHints := cloudservers.SchedulerHints{
		Group:           schedulerHintsRaw["group"].(string),
		Tenancy:         schedulerHintsRaw["tenancy"].(string),
		DedicatedHostID: schedulerHintsRaw["deh_id"].(string),
	}

	return schedulerHints
}

func resourceInstanceSchedulerHintsV2(d *schema.ResourceData, schedulerHintsRaw map[string]interface{}) schedulerhints.SchedulerHints {
	schedulerHints := schedulerhints.SchedulerHints{
		Group:           schedulerHintsRaw["group"].(string),
		Tenancy:         schedulerHintsRaw["tenancy"].(string),
		DedicatedHostID: schedulerHintsRaw["deh_id"].(string),
	}

	return schedulerHints
}

func getImageIDFromConfig(computeClient *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			return "", nil
		}
	}

	if imageId := d.Get("image_id").(string); imageId != "" {
		return imageId, nil
	} else {
		// try the OS_IMAGE_ID environment variable
		if v := os.Getenv("OS_IMAGE_ID"); v != "" {
			return v, nil
		}
	}

	imageName := d.Get("image_name").(string)
	if imageName == "" {
		// try the OS_IMAGE_NAME environment variable
		if v := os.Getenv("OS_IMAGE_NAME"); v != "" {
			imageName = v
		}
	}

	if imageName != "" {
		imageId, err := images.IDFromName(computeClient, imageName)
		if err != nil {
			return "", err
		}
		return imageId, nil
	}

	return "", fmt.Errorf("Neither a boot device, image ID, or image name were able to be determined.")
}

func setImageInformation(computeClient *golangsdk.ServiceClient, server *servers.Server, d *schema.ResourceData) error {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			d.Set("image_id", "Attempt to boot from volume - no image supplied")
			return nil
		}
	}

	imageId := server.Image["id"].(string)
	if imageId != "" {
		d.Set("image_id", imageId)
		if image, err := images.Get(computeClient, imageId).Extract(); err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				// If the image name can't be found, set the value to "Image not found".
				// The most likely scenario is that the image no longer exists in the Image Service
				// but the instance still has a record from when it existed.
				d.Set("image_name", "Image not found")
				return nil
			}
			return err
		} else {
			d.Set("image_name", image.Name)
		}
	}

	return nil
}

func getFlavorID(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	flavorId := d.Get("flavor_id").(string)

	if flavorId != "" {
		return flavorId, nil
	}

	flavorName := d.Get("flavor_name").(string)
	return flavors.IDFromName(client, flavorName)
}

func getVpcID(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	var networkID string

	networks := d.Get("network").([]interface{})
	for _, v := range networks {
		network := v.(map[string]interface{})
		networkID = network["uuid"].(string)
	}

	if networkID == "" {
		return "", fmt.Errorf("Network ID should not be empty.")
	}

	subnet, err := subnets.Get(client, networkID).Extract()
	if err != nil {
		return "", fmt.Errorf("Error retrieving Huaweicloud Subnets: %s", err)
	}

	return subnet.VPC_ID, nil
}

func resourceComputeSchedulerHintsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["group"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["group"].(string)))
	}

	if m["tenancy"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["tenancy"].(string)))
	}

	if m["deh_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["deh_id"].(string)))
	}

	return hashcode.String(buf.String())
}

func checkBlockDeviceConfig(d *schema.ResourceData) error {
	if vL, ok := d.GetOk("block_device"); ok {
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})

			if vM["source_type"] != "blank" && vM["uuid"] == "" {
				return fmt.Errorf("You must specify a uuid for %s block device types", vM["source_type"])
			}

			if vM["source_type"] == "image" && vM["destination_type"] == "volume" {
				if vM["volume_size"] == 0 {
					return fmt.Errorf("You must specify a volume_size when creating a volume from an image")
				}
			}

			if vM["source_type"] == "blank" && vM["destination_type"] == "local" {
				if vM["volume_size"] == 0 {
					return fmt.Errorf("You must specify a volume_size when creating a blank block device")
				}
			}
		}
	}

	return nil
}
