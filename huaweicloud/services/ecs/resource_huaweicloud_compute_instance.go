package ecs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/secgroups"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/powers"
	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"
	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"
	groups "github.com/chnsz/golangsdk/openstack/networking/v1/security/securitygroups"
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	powerActionMap = map[string]string{
		"ON":     "os-start",
		"OFF":    "os-stop",
		"REBOOT": "reboot",
	}
	SystemDiskType = "GPSSD"
)

// @API ECS POST /v1.1/{project_id}/cloudservers
// @API ECS POST /v1/{project_id}/cloudservers/delete
// @API ECS PUT /v1/{project_id}/cloudservers/{server_id}
// @API ECS POST /v1/{project_id}/cloudservers/action
// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/metadata
// @API ECS DELETE /v1/{project_id}/cloudservers/{server_id}/metadata/{key}
// @API ECS POST /v1.1/{project_id}/cloudservers/{server_id}/resize
// @API ECS PUT /v1/{project_id}/cloudservers/{server_id}/os-reset-password
// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/tags/action
// @API ECS POST /v2.1/{project_id}/servers/{server_id}/action
// @API ECS GET /v1/{project_id}/cloudservers/{server_id}
// @API ECS GET /v1.1/{project_id}/cloudservers/detail
// @API ECS GET /v1/{project_id}/cloudservers/{server_id}/block_device/{volume_id}
// @API ECS GET /v1/{project_id}/jobs/{job_id}
// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/changevpc
// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/migrate
// @API ECS POST /v1/{project_id}/cloudservers/actions/change-charge-mode
// @API IMS GET /v2/cloudimages
// @API EVS POST /v2.1/{project_id}/cloudvolumes/{volume_id}/action
// @API EVS GET /v2/{project_id}/cloudvolumes/{volume_id}
// @API VPC PUT /v1/{project_id}/ports/{port_id}
// @API VPC GET /v1/{project_id}/security-groups
// @API VPC GET /v1/{project_id}/subnets/{subnet_id}
func ResourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeInstanceCreate,
		ReadContext:   resourceComputeInstanceRead,
		UpdateContext: resourceComputeInstanceUpdate,
		DeleteContext: resourceComputeInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceComputeInstanceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("HW_IMAGE_ID", nil),
			},
			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("HW_IMAGE_NAME", nil),
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("HW_FLAVOR_ID", nil),
				Description: "schema: Required",
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("HW_FLAVOR_NAME", nil),
				Description: "schema: Computed",
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"security_groups": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Description:   "schema: Computed",
				ConflictsWith: []string{"security_group_ids"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"network": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 12,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Required",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "schema: Computed",
						},
						"ipv6_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_dest_check": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"fixed_ip_v6": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "schema: Computed",
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
			"system_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"system_disk_kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"system_disk_iops": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"system_disk_throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"system_disk_dss_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
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
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"throughput": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"dss_pool_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"fault_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "schema: Internal",
						},
						"tenancy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"deh_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				// just stash the hash for state & diff comparisons
				StateFunc: utils.HashAndHexEncode,
				// Suppress changes if we get a base64 format or plaint text user_data
				DiffSuppressFunc: utils.SuppressUserData,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"stop_before_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_disks_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_eip_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"include_data_disks_on_update": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_publicips_on_update": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"eip_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_type", "bandwidth"},
			},
			"eip_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith:  []string{"bandwidth"},
			},
			"bandwidth": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"eip_id"},
				RequiredWith:  []string{"eip_type"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PER", "WHOLE",
							}, true),
						},
						"id": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"bandwidth.0.size", "bandwidth.0.charge_mode"},
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							RequiredWith: []string{"bandwidth.0.charge_mode"},
						},
						"charge_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							RequiredWith: []string{"bandwidth.0.size"},
						},
						"extend_param": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Internal",
			},

			// charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid", "spot",
				}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":   common.SchemaAutoPay(nil),

			"spot_maximum_price": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"spot_duration", "spot_duration_count"},
			},
			"spot_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"spot_duration_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"spot_duration"},
			},

			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"agent_list": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"power_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// If you want to support more actions, please update powerActionMap simultaneously.
				ValidateFunc: validation.StringInSlice([]string{
					"ON", "OFF", "REBOOT", "FORCE-OFF", "FORCE-REBOOT",
				}, false),
			},
			"auto_terminate_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enclave_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			// computed attributes
			"volume_attached": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_id": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pci_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{
					Computed: true,
				}),
			},
			"system_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip_v6": {
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
			"expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getSpotDurationCount(d *schema.ResourceData) int {
	var count = 1
	if c, ok := d.GetOk("spot_duration_count"); ok {
		count = c.(int)
	}
	return count
}

func resourceComputeInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute v1 client: %s", err)
	}
	ecsV11Client, err := cfg.ComputeV11Client(region)
	if err != nil {
		return diag.Errorf("error creating compute v1.1 client: %s", err)
	}
	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating image client: %s", err)
	}
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v1 client: %s", err)
	}
	nicClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v2 client: %s", err)
	}

	// Determines the Image ID using the following rules:
	// If an image_id was specified, use it.
	// If an image_name was specified, look up the image ID, report if error.
	imageId, err := getImageIDFromConfig(d, imsClient)
	if err != nil {
		return diag.FromErr(err)
	}

	flavorId, err := getFlavorID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vpcId, err := getVpcID(d, vpcClient)
	if err != nil {
		return diag.FromErr(err)
	}

	secGroupIDs, err := buildInstanceSecGroupIds(d, vpcClient)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := &cloudservers.CreateOpts{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		ImageRef:          imageId,
		FlavorRef:         flavorId,
		KeyName:           d.Get("key_pair").(string),
		VpcId:             vpcId,
		SecurityGroups:    secGroupIDs,
		AvailabilityZone:  d.Get("availability_zone").(string),
		RootVolume:        buildInstanceRootVolume(d),
		DataVolumes:       buildInstanceDataVolumes(d),
		Nics:              buildInstanceNicsRequest(d),
		PublicIp:          buildInstancePublicIPRequest(d),
		UserData:          []byte(d.Get("user_data").(string)),
		AutoTerminateTime: d.Get("auto_terminate_time").(string),
		EnclaveOptions:    buildInstanceEnclaveOptionsPRequest(d),
	}

	if v := d.Get("enable_jumbo_frame").(bool); v {
		createOpts.EnableJumboFrame = &v
	}

	if tags, ok := d.GetOk("tags"); ok {
		createOpts.ServerTags = utils.ExpandResourceTags(tags.(map[string]interface{}))
	}

	var extendParam cloudservers.ServerExtendParam
	chargingMode := d.Get("charging_mode").(string)
	if chargingMode == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}

		extendParam.ChargingMode = chargingMode
		extendParam.PeriodType = d.Get("period_unit").(string)
		extendParam.PeriodNum = d.Get("period").(int)
		extendParam.IsAutoRenew = d.Get("auto_renew").(string)
		extendParam.IsAutoPay = common.GetAutoPay(d)
	} else if chargingMode == "spot" {
		extendParam.MarketType = "spot"
		extendParam.SpotPrice = d.Get("spot_maximum_price").(string)
		if v, ok := d.GetOk("spot_duration"); ok {
			extendParam.InterruptionPolicy = "immediate"
			extendParam.SpotDurationHours = v.(int)
			extendParam.SpotDurationCount = getSpotDurationCount(d)
		}
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		extendParam.EnterpriseProjectId = epsID
	}
	if extendParam != (cloudservers.ServerExtendParam{}) {
		createOpts.ExtendParam = &extendParam
	}

	var metadata cloudservers.MetaData
	metadata.OpSvcUserId = getOpSvcUserID(d, cfg)

	if v, ok := d.GetOk("agency_name"); ok {
		metadata.AgencyName = v.(string)
	}
	if v, ok := d.GetOk("agent_list"); ok {
		metadata.AgentList = v.(string)
	}
	if metadata != (cloudservers.MetaData{}) {
		createOpts.MetaData = &metadata
	}

	schedulerHintsRaw := d.Get("scheduler_hints").([]interface{})
	if len(schedulerHintsRaw) > 0 {
		if m, ok := schedulerHintsRaw[0].(map[string]interface{}); ok {
			schedulerHints := buildInstanceSchedulerHints(m)
			createOpts.SchedulerHints = &schedulerHints
		} else {
			log.Printf("[WARN] can not build scheduler hints: %+v", schedulerHintsRaw[0])
		}
	}

	log.Printf("[DEBUG] ECS create options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.AdminPass = d.Get("admin_pass").(string)

	if d.Get("charging_mode") == "prePaid" {
		// prePaid.
		n, err := cloudservers.CreatePrePaid(ecsV11Client, createOpts).ExtractOrderResponse()
		if err != nil {
			return diag.Errorf("error creating server: %s", err)
		}
		if len(n.ServerIDs) != 0 {
			d.SetId(n.ServerIDs[0])
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, n.OrderID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, n.OrderID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	} else {
		// postPaid.
		n, err := cloudservers.Create(ecsV11Client, createOpts).ExtractJobResponse()
		if err != nil {
			return diag.Errorf("error creating server: %s", err)
		}
		if len(n.ServerIDs) != 0 {
			d.SetId(n.ServerIDs[0])
		}
		if err := cloudservers.WaitForJobSuccess(ecsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
			return diag.FromErr(err)
		}
		serverId, err := cloudservers.GetJobEntity(ecsClient, n.JobID, "server_id")
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(serverId.(string))
	}

	// update the user-defined metadata if necessary
	if v, ok := d.GetOk("metadata"); ok {
		metadataOpts := v.(map[string]interface{})
		log.Printf("[DEBUG] ECS metadata options: %v", metadataOpts)

		_, err := cloudservers.UpdateMetadata(ecsClient, d.Id(), metadataOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating the metadata: %s", err)
		}
	}

	// Create an instance in the shutdown state.
	if action, ok := d.GetOk("power_action"); ok {
		action := action.(string)
		if action == "OFF" || action == "FORCE-OFF" {
			if err = doPowerAction(ecsClient, d, action); err != nil {
				return diag.Errorf("Doing power action (%s) for instance (%s) failed: %s", action, d.Id(), err)
			}
		} else {
			log.Printf("[WARN] the power action (%s) is invalid after instance created", action)
		}
	}

	// get the original value of source_dest_check in script
	originalNetworks := d.Get("network").([]interface{})
	sourceDestChecks := make([]bool, len(originalNetworks))
	var flag bool

	for i, v := range originalNetworks {
		nic := v.(map[string]interface{})
		sourceDestChecks[i] = nic["source_dest_check"].(bool)
		if !flag && !sourceDestChecks[i] {
			flag = true
		}
	}

	if flag {
		// Get the instance network and address information
		server, err := cloudservers.Get(ecsClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error retrieving compute instance: %s", d.Id())
		}
		networks, err := flattenInstanceNetworks(d, meta, server)
		if err != nil {
			return diag.FromErr(err)
		}

		for i, nic := range networks {
			nicPort := nic["port"].(string)
			if nicPort == "" {
				continue
			}

			if !sourceDestChecks[i] {
				if err := disableSourceDestCheck(nicClient, nicPort); err != nil {
					return diag.Errorf("error disabling source dest check on port(%s) of instance(%s): %s", nicPort, d.Id(), err)
				}
			}
		}
	}

	return resourceComputeInstanceRead(ctx, d, meta)
}

func resourceComputeInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	serverID := d.Id()
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}
	blockStorageClient, err := cfg.BlockStorageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating evs client: %s", err)
	}
	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating image client: %s", err)
	}

	server, err := cloudservers.Get(ecsClient, serverID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving compute instance")
	} else if server.Status == "DELETED" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] retrieved compute instance %s: %+v", serverID, server)
	// Set some attributes
	d.Set("region", region)
	d.Set("enterprise_project_id", server.EnterpriseProjectID)
	d.Set("availability_zone", server.AvailabilityZone)
	d.Set("name", server.Name)
	d.Set("description", server.Description)
	d.Set("hostname", server.Hostname)
	d.Set("status", server.Status)
	d.Set("agency_name", server.Metadata.AgencyName)
	d.Set("agent_list", server.Metadata.AgentList)
	d.Set("charging_mode", normalizeChargingMode(server.Metadata.ChargingMode))
	d.Set("created_at", server.Created.Format(time.RFC3339))
	d.Set("updated_at", server.Updated.Format(time.RFC3339))
	d.Set("auto_terminate_time", server.AutoTerminateTime)
	d.Set("public_ip", computePublicIP(server))
	d.Set("enclave_options", flattenEnclaveOptions(server.EnclaveOptions))

	flavorInfo := server.Flavor
	d.Set("flavor_id", flavorInfo.ID)
	d.Set("flavor_name", flavorInfo.Name)

	// Set the instance's image information appropriately
	if err := setImageInformation(d, imsClient, server.Image.ID); err != nil {
		return diag.FromErr(err)
	}

	if server.KeyName != "" {
		d.Set("key_pair", server.KeyName)
	}

	// Get the instance network and address information
	networks, err := flattenInstanceNetworks(d, meta, server)
	if err != nil {
		return diag.FromErr(err)
	}
	// Determine the best IPv4 and IPv6 addresses to access the instance with
	hostv4, hostv6 := getInstanceAccessAddresses(networks)

	// update hostv4/6 by AccessIPv4/v6
	// currently, AccessIPv4/v6 are Reserved in HuaweiCloud
	if server.AccessIPv4 != "" {
		hostv4 = server.AccessIPv4
	}
	if server.AccessIPv6 != "" {
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

	secGrpNames := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpNames = append(secGrpNames, sg.Name)
	}
	d.Set("security_groups", secGrpNames)

	secGrpIDs := make([]string, len(server.SecurityGroups))
	for i, sg := range server.SecurityGroups {
		secGrpIDs[i] = sg.ID
	}
	d.Set("security_group_ids", secGrpIDs)

	// Set volume attached
	if len(server.VolumeAttached) > 0 {
		bds := make([]map[string]interface{}, len(server.VolumeAttached))
		for i, b := range server.VolumeAttached {
			// retrieve volume `size` and `type`
			volumeInfo, err := cloudvolumes.Get(blockStorageClient, b.ID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] retrieved volume %s: %#v", b.ID, volumeInfo)

			// retrieve volume `pci_address`
			va, err := block_devices.Get(ecsClient, serverID, b.ID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] retrieved block device %s: %#v", b.ID, va)

			bds[i] = map[string]interface{}{
				"volume_id":   b.ID,
				"size":        volumeInfo.Size,
				"type":        volumeInfo.VolumeType,
				"boot_index":  va.BootIndex,
				"pci_address": va.PciAddress,
				"kms_key_id":  volumeInfo.Metadata.SystemCmkID,
			}

			if va.BootIndex == 0 {
				d.Set("system_disk_id", b.ID)
				d.Set("system_disk_size", volumeInfo.Size)
				d.Set("system_disk_type", volumeInfo.VolumeType)
				d.Set("system_disk_kms_key_id", volumeInfo.Metadata.SystemCmkID)
				d.Set("system_disk_iops", volumeInfo.IOPS.TotalVal)
				d.Set("system_disk_throughput", volumeInfo.Throughput.TotalVal)
			}
		}
		d.Set("volume_attached", bds)
	}

	// set scheduler_hints
	osHints := server.OsSchedulerHints
	if len(osHints.Group) > 0 || len(osHints.Tenancy) > 0 || len(osHints.DedicatedHostID) > 0 {
		schedulerHint := make(map[string]interface{})
		if len(osHints.Group) > 0 {
			schedulerHint["group"] = osHints.Group[0]
		}
		if len(osHints.Tenancy) > 0 {
			schedulerHint["tenancy"] = osHints.Tenancy[0]
		}
		if len(osHints.DedicatedHostID) > 0 {
			schedulerHint["deh_id"] = osHints.DedicatedHostID[0]
		}
		d.Set("scheduler_hints", []map[string]interface{}{schedulerHint})
	}

	// Set instance tags
	d.Set("tags", flattenTagsToMap(d, server.Tags))

	// Set expired time for prePaid instance
	if normalizeChargingMode(server.Metadata.ChargingMode) == "prePaid" {
		expiredTime, err := getPrePaidExpiredTime(d, cfg, serverID)
		if err != nil {
			log.Printf("error get prePaid expired time: %s", err)
		}

		d.Set("expired_time", expiredTime)
	}

	return nil
}

func getPrePaidExpiredTime(d *schema.ResourceData, cfg *config.Config, instanceID string) (string, error) {
	product := "ecsv11"
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return "", fmt.Errorf("error creating EVS client: %s", err)
	}

	getPrePaidExpiredTimeHttpUrl := "v1.1/{project_id}/cloudservers/detail?id={instance_id}&expect-fields=market_info"
	getPrePaidExpiredTimePath := client.Endpoint + getPrePaidExpiredTimeHttpUrl
	getPrePaidExpiredTimePath = strings.ReplaceAll(getPrePaidExpiredTimePath, "{project_id}", client.ProjectID)
	getPrePaidExpiredTimePath = strings.ReplaceAll(getPrePaidExpiredTimePath, "{instance_id}", instanceID)

	getPrePaidExpiredTimeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getPrePaidExpiredTimeResp, err := client.Request("GET", getPrePaidExpiredTimePath, &getPrePaidExpiredTimeOpt)
	if err != nil {
		return "", err
	}

	getPrePaidExpiredTimeRespBody, err := utils.FlattenResponse(getPrePaidExpiredTimeResp)
	if err != nil {
		return "", err
	}

	expiredTime := utils.PathSearch("servers[0].market_info.prepaid_info.expired_time", getPrePaidExpiredTimeRespBody, "")

	return expiredTime.(string), nil
}

func normalizeChargingMode(mode string) string {
	var ret string
	switch mode {
	case "1":
		ret = "prePaid"
	case "2":
		ret = "spot"
	default:
		ret = "postPaid"
	}

	return ret
}

func flattenEnclaveOptions(enclaveOptions *cloudservers.EnclaveOptions) []map[string]interface{} {
	if enclaveOptions == nil {
		return nil
	}

	res := []map[string]interface{}{
		{
			"enabled": enclaveOptions.Enabled,
		},
	}

	return res
}

func resourceComputeInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	computeClient, err := cfg.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V2 client: %s", err)
	}
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}
	ecsV11Client, err := cfg.ComputeV11Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1.1 client: %s", err)
	}
	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	serverID := d.Id()
	if d.HasChange("charging_mode") {
		if d.Get("charging_mode").(string) != "prePaid" {
			return diag.Errorf("error updating the charging mode of the ECS instance (%s): %s", d.Id(),
				"only support change to pre-paid")
		}
		if err = updateInstanceChargingMode(ctx, d, ecsClient, bssClient); err != nil {
			return diag.FromErr(err)
		}
	} else if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), serverID); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", serverID, err)
		}
	}

	if d.HasChanges("name", "description", "user_data") {
		var updateOpts cloudservers.UpdateOpts
		updateOpts.Name = d.Get("name").(string)
		updateOpts.UserData = []byte(d.Get("user_data").(string))
		description := d.Get("description").(string)
		updateOpts.Description = &description

		err := cloudservers.Update(ecsClient, serverID, updateOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating server: %s", err)
		}
	}

	if d.HasChanges("agency_name", "agent_list") {
		metadataOpts := make(map[string]interface{})
		if d.HasChange("agency_name") {
			metadataOpts["agency_name"] = d.Get("agency_name").(string)
		}
		if d.HasChange("agent_list") {
			metadataOpts["__support_agent_list"] = d.Get("agent_list").(string)
		}
		_, err = cloudservers.UpdateMetadata(ecsClient, serverID, metadataOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating the metadata(agency_name, agent_list): %s", err)
		}
	}

	if d.HasChanges("metadata") {
		if err := updateInstanceMetaData(d, ecsClient, serverID); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("security_group_ids", "security_groups") {
		var oldSGRaw interface{}
		var newSGRaw interface{}
		if d.HasChange("security_group_ids") {
			oldSGRaw, newSGRaw = d.GetChange("security_group_ids")
		} else {
			oldSGRaw, newSGRaw = d.GetChange("security_groups")
		}
		oldSGSet := oldSGRaw.(*schema.Set)
		newSGSet := newSGRaw.(*schema.Set)
		secgroupsToAdd := newSGSet.Difference(oldSGSet)
		secgroupsToRemove := oldSGSet.Difference(newSGSet)
		log.Printf("[DEBUG] security groups to add: %v", secgroupsToAdd)
		log.Printf("[DEBUG] security groups to remove: %v", secgroupsToRemove)

		for _, g := range secgroupsToRemove.List() {
			err := secgroups.RemoveServer(computeClient, serverID, g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					continue
				}
				return diag.Errorf("error removing security group (%s) from server (%s): %s", g, serverID, err)
			}
			log.Printf("[DEBUG] removed security group (%s) from instance (%s)", g, serverID)
		}

		for _, g := range secgroupsToAdd.List() {
			err := secgroups.AddServer(computeClient, serverID, g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return diag.Errorf("error adding security group (%s) to server (%s): %s", g, serverID, err)
			}
			log.Printf("[DEBUG] added security group (%s) to instance (%s)", g, serverID)
		}
	}

	if d.HasChange("admin_pass") {
		if newPwd, ok := d.Get("admin_pass").(string); ok {
			err := cloudservers.ChangeAdminPassword(ecsClient, serverID, newPwd).ExtractErr()
			if err != nil {
				return diag.Errorf("error changing admin password of server (%s): %s", serverID, err)
			}
		}
	}

	if d.HasChanges("flavor_id", "flavor_name") {
		newFlavorId, err := getFlavorID(d)
		if err != nil {
			return diag.FromErr(err)
		}

		extendParam := &cloudservers.ResizeExtendParam{
			AutoPay: common.GetAutoPay(d),
		}
		resizeOpts := &cloudservers.ResizeOpts{
			FlavorRef:   newFlavorId,
			Mode:        "withStopServer",
			ExtendParam: extendParam,
		}
		log.Printf("[DEBUG] resize configuration: %#v", resizeOpts)
		job, err := cloudservers.Resize(ecsV11Client, resizeOpts, serverID).ExtractJobResponse()
		if err != nil {
			return diag.Errorf("error resizing server: %s", err)
		}

		if err := cloudservers.WaitForJobSuccess(ecsClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), job.JobID); err != nil {
			return diag.Errorf("error waiting for instance (%s) to be resized: %s", serverID, err)
		}
	}

	if d.HasChange("network") {
		var err error
		nicClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating networking client: %s", err)
		}

		if err := updateSourceDestCheck(d, nicClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(ecsClient, d, "cloudservers", serverID)
		if tagErr != nil {
			return diag.Errorf("error updating tags of instance:%s, err:%s", serverID, tagErr)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "ecs",
			RegionId:     region,
			ProjectId:    ecsClient.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("system_disk_size") {
		extendOpts := cloudvolumes.ExtendOpts{
			SizeOpts: cloudvolumes.ExtendSizeOpts{
				NewSize: d.Get("system_disk_size").(int),
			},
		}

		if strings.EqualFold(d.Get("charging_mode").(string), "prePaid") {
			extendOpts.ChargeInfo = &cloudvolumes.ExtendChargeOpts{
				IsAutoPay: common.GetAutoPay(d),
			}
		}

		evsV2Client, err := cfg.BlockStorageV2Client(region)
		if err != nil {
			return diag.Errorf("error creating evs V2 client: %s", err)
		}
		evsV21Client, err := cfg.BlockStorageV21Client(region)
		if err != nil {
			return diag.Errorf("error creating evs V2.1 client: %s", err)
		}

		systemDiskID := d.Get("system_disk_id").(string)

		resp, err := cloudvolumes.ExtendSize(evsV21Client, systemDiskID, extendOpts).Extract()
		if err != nil {
			return diag.Errorf("error extending EVS volume (%s) size: %s", systemDiskID, err)
		}

		if strings.EqualFold(d.Get("charging_mode").(string), "prePaid") {
			err = common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.Errorf("The order (%s) is not completed while extending system disk (%s) size: %v",
					resp.OrderID, serverID, err)
			}
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    cloudVolumeRefreshFunc(evsV2Client, systemDiskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for huaweicloud_compute_instance system disk %s to become ready: %s", systemDiskID, err)
		}
	}

	// update the key_pair before power action
	if d.HasChange("key_pair") {
		kmsClient, err := cfg.KmsV3Client(region)
		if err != nil {
			return diag.Errorf("error creating KMS v3 client: %s", err)
		}

		o, n := d.GetChange("key_pair")
		keyPairOpts := &common.KeypairAuthOpts{
			InstanceID:       serverID,
			InUsedKeyPair:    o.(string),
			NewKeyPair:       n.(string),
			InUsedPrivateKey: d.Get("private_key").(string),
			Password:         d.Get("admin_pass").(string),
			Timeout:          d.Timeout(schema.TimeoutUpdate),
		}
		if err := common.UpdateEcsInstanceKeyPair(ctx, ecsClient, kmsClient, keyPairOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	// The instance power status update needs to be done at the end
	if d.HasChange("power_action") {
		action := d.Get("power_action").(string)
		if err = doPowerAction(ecsClient, d, action); err != nil {
			return diag.Errorf("Doing power action (%s) for instance (%s) failed: %s", action, serverID, err)
		}
	}

	if d.HasChange("auto_terminate_time") {
		terminateTime := d.Get("auto_terminate_time").(string)
		err := cloudservers.UpdateAutoTerminateTime(ecsClient, serverID, terminateTime).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating auto-terminate-time of server (%s): %s", serverID, err)
		}
	}

	if d.HasChanges("network.0.uuid", "network.0.fixed_ip_v4") {
		vpcClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating networking client: %s", err)
		}
		err = updateInstanceNetwork(ctx, d, ecsClient, vpcClient, serverID)
		if err != nil {
			return diag.Errorf("error updating network of server (%s): %s", serverID, err)
		}
	}

	if d.HasChanges("scheduler_hints.0.deh_id") {
		err = updateInstanceDehId(ctx, d, ecsClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceComputeInstanceRead(ctx, d, meta)
}

func cloudVolumeRefreshFunc(c *golangsdk.ServiceClient, volumeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := cloudvolumes.Get(c, volumeId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return response, "deleted", nil
			}
			return response, "ERROR", err
		}
		if response != nil {
			return response, response.Status, nil
		}
		return response, "ERROR", nil
	}
}

func updateInstanceMetaData(d *schema.ResourceData, client *golangsdk.ServiceClient, serverID string) error {
	oldRaw, newRaw := d.GetChange("metadata")
	oldMetadata := oldRaw.(map[string]interface{})
	newMetadata := newRaw.(map[string]interface{})

	// Determine if any metadata keys will be removed from the configuration.
	// Then request those keys to be deleted.
	var metadataToDelete []string
	for oldKey := range oldMetadata {
		var found bool
		for newKey := range newMetadata {
			if oldKey == newKey {
				found = true
				break
			}
		}

		if !found {
			metadataToDelete = append(metadataToDelete, oldKey)
		}
	}

	for _, key := range metadataToDelete {
		err := cloudservers.DeleteMetadatItem(client, serverID, key).ExtractErr()
		if err != nil {
			return fmt.Errorf("error deleting metadata (%s) from server: %s", key, err)
		}
	}

	// Update existing metadata and add any new metadata.
	if len(newMetadata) > 0 {
		_, err := cloudservers.UpdateMetadata(client, serverID, newMetadata).Extract()
		if err != nil {
			return fmt.Errorf("error updating the metadata: %s", err)
		}
	}

	return nil
}

func updateInstanceNetwork(ctx context.Context, d *schema.ResourceData, client, vpcClient *golangsdk.ServiceClient, serverID string) error {
	updateNetworkHttpUrl := "v1/{project_id}/cloudservers/{server_id}/changevpc"
	updateNetworkPath := client.Endpoint + updateNetworkHttpUrl
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{project_id}", client.ProjectID)
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{server_id}", serverID)

	vpcID, err := getVpcID(d, vpcClient)
	if err != nil {
		return err
	}

	updateNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateInstanceNetworkOpts(d, vpcID)),
	}
	updateNetworkResp, err := client.Request("POST", updateNetworkPath, &updateNetworkOpt)
	if err != nil {
		return fmt.Errorf("error udpating ECS network: %s", err)
	}
	updateNetworkRespBody, err := utils.FlattenResponse(updateNetworkResp)
	if err != nil {
		return err
	}
	jobID := utils.PathSearch("job_id", updateNetworkRespBody, "").(string)
	if jobID == "" {
		return errors.New("unable to find the job ID from the API response")
	}

	// Wait for job status become `SUCCESS`.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      getJobRefreshFunc(client, jobID),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ECS network updated: %s", err)
	}

	return nil
}

func buildUpdateInstanceNetworkOpts(d *schema.ResourceData, vpcID string) map[string]interface{} {
	var ipAddress interface{}
	networks := d.Get("network")

	networkID := utils.PathSearch("[0].uuid", networks, nil)
	if d.HasChange("network.0.fixed_ip_v4") {
		ipAddress = utils.PathSearch("[0].fixed_ip_v4", networks, nil)
	}

	bodyParam := map[string]interface{}{
		"vpc_id": vpcID,
		"nic": map[string]interface{}{
			"subnet_id":       networkID,
			"ip_address":      utils.ValueIgnoreEmpty(ipAddress),
			"security_groups": buildUpdateInstanceNetworkSecgroupOpts(d),
		},
	}

	return bodyParam
}

func buildUpdateInstanceNetworkSecgroupOpts(d *schema.ResourceData) []map[string]interface{} {
	secgroupIDs := d.Get("security_group_ids").(*schema.Set).List()
	bodyParams := make([]map[string]interface{}, len(secgroupIDs))

	for i, v := range secgroupIDs {
		bodyParams[i] = map[string]interface{}{
			"id": v,
		}
	}

	return bodyParams
}

func updateInstanceChargingMode(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	httpUrl := "v1/{project_id}/cloudservers/actions/change-charge-mode"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateInstanceChargingMode(d),
	}
	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error udpating ECS charging mode: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}
	orderId := utils.PathSearch("order_id", updateRespBody, "").(string)
	if orderId == "" {
		return errors.New("order_id is not found in the API response")
	}

	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateInstanceChargingMode(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"server_ids":  []string{d.Id()},
		"charge_mode": "prePaid",
		"prepaid_options": map[string]interface{}{
			"include_data_disks": d.Get("include_data_disks_on_update"),
			"include_publicips":  d.Get("include_publicips_on_update"),
			"period_type":        d.Get("period_unit"),
			"period_num":         d.Get("period"),
			"auto_pay":           true,
			"auto_renew":         d.Get("auto_renew"),
		},
	}

	return bodyParam
}

func updateInstanceDehId(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	httpUrl := "v1/{project_id}/cloudservers/{server_id}/migrate"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateInstanceDehIdOpts(d),
	}
	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error udpating ECS dedicated host ID: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}
	jobID := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobID == "" {
		return errors.New("unable to find the job ID from the API response")
	}

	// Wait for job status become `SUCCESS`.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      getJobRefreshFunc(client, jobID),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ECS dedicated host ID updated: %s", err)
	}

	return nil
}

func buildUpdateInstanceDehIdOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"migrate": nil,
	}
	schedulerHintsRaw := d.Get("scheduler_hints").([]interface{})
	if len(schedulerHintsRaw) > 0 {
		if m, ok := schedulerHintsRaw[0].(map[string]interface{}); ok {
			dedicatedHostId := m["deh_id"]
			if dedicatedHostId != nil && len(dedicatedHostId.(string)) > 0 {
				bodyParam["migrate"] = map[string]interface{}{
					"dedicated_host_id": dedicatedHostId,
				}
			}
		}
	}

	return bodyParam
}

func resourceComputeInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	if d.Get("stop_before_destroy").(bool) {
		if err = doPowerAction(ecsClient, d, "FORCE-OFF"); err != nil {
			log.Printf("[WARN] error stopping instance: %s", err)
		} else {
			log.Printf("[DEBUG] waiting for instance (%s) to stop", d.Id())
			pending := []string{"ACTIVE"}
			target := []string{"SHUTOFF"}
			timeout := d.Timeout(schema.TimeoutDelete)
			if err := waitForServerTargetState(ctx, ecsClient, d.Id(), pending, target, timeout); err != nil {
				return diag.Errorf("State waiting timeout: %s", err)
			}
		}
	}

	if d.Get("charging_mode") == "prePaid" {
		resources, err := calcUnsubscribeResources(d, cfg)
		if err != nil {
			return diag.Errorf("error unsubscribe ECS server: %s", err)
		}

		log.Printf("[DEBUG] %v will be unsubscribed", resources)
		if err := common.UnsubscribePrePaidResource(d, cfg, resources); err != nil {
			return diag.Errorf("error unsubscribe ECS server: %s", err)
		}
	} else {
		serverRequests := []cloudservers.Server{
			{Id: d.Id()},
		}

		deleteOpts := cloudservers.DeleteOpts{
			Servers:        serverRequests,
			DeleteVolume:   d.Get("delete_disks_on_termination").(bool),
			DeletePublicIP: d.Get("delete_eip_on_termination").(bool),
		}

		n, err := cloudservers.Delete(ecsClient, deleteOpts).ExtractJobResponse()
		if err != nil {
			return diag.Errorf("error deleting server: %s", err)
		}

		if err := cloudservers.WaitForJobSuccess(ecsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
			return diag.FromErr(err)
		}
	}

	// Instance may still exist after Order/Job succeed.
	pending := []string{"ACTIVE", "SHUTOFF"}
	target := []string{"DELETED", "SOFT_DELETED"}
	deleteTimeout := d.Timeout(schema.TimeoutDelete)
	if err := waitForServerTargetState(ctx, ecsClient, d.Id(), pending, target, deleteTimeout); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceComputeInstanceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating compute client: %s", err)
	}

	server, err := cloudservers.Get(ecsClient, d.Id()).Extract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "compute instance")
	}

	allInstanceNics, err := getInstanceAddresses(d, meta, server)
	if err != nil {
		return nil, fmt.Errorf("error fetching networks of compute instance %s: %s", d.Id(), err)
	}

	networks := []map[string]interface{}{}
	for _, nic := range allInstanceNics {
		v := map[string]interface{}{
			"uuid":              nic.NetworkID,
			"port":              nic.PortID,
			"fixed_ip_v4":       nic.FixedIPv4,
			"fixed_ip_v6":       nic.FixedIPv6,
			"ipv6_enable":       nic.FixedIPv6 != "",
			"source_dest_check": nic.SourceDestCheck,
			"mac":               nic.MAC,
		}
		networks = append(networks, v)
	}

	log.Printf("[DEBUG] flatten Instance Networks: %#v", networks)
	d.Set("network", networks)

	return []*schema.ResourceData{d}, nil
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
			fault := fmt.Errorf("error code: %d, message: %s", s.Fault.Code, s.Fault.Message)
			return s, "ERROR", fault
		}
		return s, s.Status, nil
	}
}

func buildInstanceSecGroupIds(d *schema.ResourceData, client *golangsdk.ServiceClient) ([]cloudservers.SecurityGroup, error) {
	if v, ok := d.GetOk("security_group_ids"); ok {
		rawSecGroups := v.(*schema.Set).List()
		secGroups := make([]cloudservers.SecurityGroup, len(rawSecGroups))
		for i, raw := range rawSecGroups {
			secGroups[i] = cloudservers.SecurityGroup{
				ID: raw.(string),
			}
		}
		return secGroups, nil
	}

	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secGroups := make([]cloudservers.SecurityGroup, 0, len(rawSecGroups))

	opt := groups.ListOpts{
		EnterpriseProjectId: "all_granted_eps",
	}
	pages, err := groups.List(client, opt).AllPages()
	if err != nil {
		return nil, err
	}
	resp, err := groups.ExtractSecurityGroups(pages)
	if err != nil {
		return nil, err
	}

	for _, raw := range rawSecGroups {
		secName := raw.(string)
		for _, secGroup := range resp {
			if secName == secGroup.Name {
				secGroups = append(secGroups, cloudservers.SecurityGroup{
					ID: secGroup.ID,
				})
				break
			}
		}
	}
	if len(secGroups) != len(rawSecGroups) {
		return nil, fmt.Errorf("the list contains invalid security groups (num: %d), please check your entry",
			len(rawSecGroups)-len(secGroups))
	}

	return secGroups, nil
}

func getOpSvcUserID(d *schema.ResourceData, conf *config.Config) string {
	if v, ok := d.GetOk("user_id"); ok {
		return v.(string)
	}
	return conf.UserID
}

func buildInstanceNicsRequest(d *schema.ResourceData) []cloudservers.Nic {
	var nicRequests []cloudservers.Nic

	networks := d.Get("network").([]interface{})
	for _, v := range networks {
		network := v.(map[string]interface{})
		nicRequest := cloudservers.Nic{
			SubnetId:   network["uuid"].(string),
			IpAddress:  network["fixed_ip_v4"].(string),
			Ipv6Enable: network["ipv6_enable"].(bool),
		}

		nicRequests = append(nicRequests, nicRequest)
	}
	return nicRequests
}

func buildInstancePublicIPRequest(d *schema.ResourceData) *cloudservers.PublicIp {
	if v, ok := d.GetOk("eip_id"); ok {
		return &cloudservers.PublicIp{
			Id: v.(string),
		}
	}

	bandWidthRaw := d.Get("bandwidth").([]interface{})
	if len(bandWidthRaw) != 1 {
		return nil
	}

	bandWidth := bandWidthRaw[0].(map[string]interface{})
	bwOpts := cloudservers.BandWidth{
		ShareType:  bandWidth["share_type"].(string),
		Id:         bandWidth["id"].(string),
		ChargeMode: bandWidth["charge_mode"].(string),
		Size:       bandWidth["size"].(int),
	}

	var eipExtendOpts *cloudservers.EipExtendParam
	extendParam := bandWidth["extend_param"].(map[string]interface{})
	if v, ok := extendParam["charging_mode"]; ok {
		eipExtendOpts = &cloudservers.EipExtendParam{
			ChargingMode: v.(string),
		}
	}

	return &cloudservers.PublicIp{
		Eip: &cloudservers.Eip{
			IpType:      d.Get("eip_type").(string),
			BandWidth:   &bwOpts,
			ExtendParam: eipExtendOpts,
		},
		DeleteOnTermination: d.Get("delete_eip_on_termination").(bool),
	}
}

func buildInstanceSchedulerHints(schedulerHintsRaw map[string]interface{}) cloudservers.SchedulerHints {
	schedulerHints := cloudservers.SchedulerHints{
		Group:           schedulerHintsRaw["group"].(string),
		FaultDomain:     schedulerHintsRaw["fault_domain"].(string),
		Tenancy:         schedulerHintsRaw["tenancy"].(string),
		DedicatedHostID: schedulerHintsRaw["deh_id"].(string),
	}

	return schedulerHints
}

func buildInstanceEnclaveOptionsPRequest(d *schema.ResourceData) *cloudservers.EnclaveOptions {
	v, ok := d.GetOk("enclave_options")
	if !ok {
		return nil
	}

	enclaveOptionsRaw := v.([]interface{})[0]

	res := cloudservers.EnclaveOptions{
		Enabled: utils.PathSearch("enabled", enclaveOptionsRaw, false).(bool),
	}

	return &res
}

func getImage(client *golangsdk.ServiceClient, id, name string) (*cloudimages.Image, error) {
	listOpts := &cloudimages.ListOpts{
		ID:                  id,
		Name:                name,
		Limit:               1,
		EnterpriseProjectID: "all_granted_eps",
	}
	allPages, err := cloudimages.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return nil, fmt.Errorf("unable to find images %s: Maybe not existed", id)
	}

	img := allImages[0]
	if id != "" && img.ID != id {
		return nil, fmt.Errorf("unexpected images ID")
	}
	if name != "" && img.Name != name {
		return nil, fmt.Errorf("unexpected images Name")
	}
	log.Printf("[DEBUG] retrieved image %s: %#v", id, img)
	return &img, nil
}

func getImageIDFromConfig(d *schema.ResourceData, imsClient *golangsdk.ServiceClient) (string, error) {
	if imageID := d.Get("image_id").(string); imageID != "" {
		return imageID, nil
	}

	if imageName := d.Get("image_name").(string); imageName != "" {
		img, err := getImage(imsClient, "", imageName)
		if err != nil {
			return "", err
		}
		return img.ID, nil
	}

	return "", fmt.Errorf("neither a image ID or image name were able to be determined")
}

func setImageInformation(d *schema.ResourceData, imsClient *golangsdk.ServiceClient, imageID string) error {
	d.Set("image_id", imageID)
	image, err := getImage(imsClient, imageID, "")
	if err != nil {
		// If the image name can't be found, set the value to "Image not found".
		// The most likely scenario is that the image no longer exists in the Image Service
		// but the instance still has a record from when it existed.
		d.Set("image_name", "Image not found")
		return nil
	}
	d.Set("image_name", image.Name)

	return nil
}

// computePublicIP get the first floating address
func computePublicIP(server *cloudservers.CloudServer) string {
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

func getFlavorID(d *schema.ResourceData) (string, error) {
	var flavorID string

	// both flavor_id and flavor_name are the same value
	if v1, ok := d.GetOk("flavor_id"); ok {
		flavorID = v1.(string)
	} else if v2, ok := d.GetOk("flavor_name"); ok {
		flavorID = v2.(string)
	}

	if flavorID == "" {
		return "", fmt.Errorf("missing required argument: the `flavor_id` must be specified")
	}
	return flavorID, nil
}

func getVpcID(d *schema.ResourceData, client *golangsdk.ServiceClient) (string, error) {
	var networkID string

	networks := d.Get("network").([]interface{})
	if len(networks) > 0 {
		// all networks belongs to one VPC
		network := networks[0].(map[string]interface{})
		networkID = network["uuid"].(string)
	}

	if networkID == "" {
		return "", fmt.Errorf("network ID should not be empty")
	}

	subnet, err := subnets.Get(client, networkID).Extract()
	if err != nil {
		return "", fmt.Errorf("error retrieving subnets: %s", err)
	}

	return subnet.VPC_ID, nil
}

func waitForServerTargetState(ctx context.Context, client *golangsdk.ServiceClient, instanceID string, pending, target []string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      pending,
		Target:       target,
		Refresh:      ServerV1StateRefreshFunc(client, instanceID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to become target state (%v): %s", instanceID, target, err)
	}
	return nil
}

// doPowerAction is a method for instance power doing shutdown, startup and reboot actions.
func doPowerAction(client *golangsdk.ServiceClient, d *schema.ResourceData, action string) error {
	var jobResp *cloudservers.JobResponse
	powerOpts := powers.PowerOpts{
		Servers: []powers.ServerInfo{
			{ID: d.Id()},
		},
	}
	// In the reboot structure, Type is a required option.
	// Since the type of power off and reboot is 'SOFT' by default, setting this value has solved the power structural
	// compatibility problem between optional and required.
	if action != "ON" {
		powerOpts.Type = "SOFT"
	}
	if strings.HasPrefix(action, "FORCE-") {
		powerOpts.Type = "HARD"
		action = strings.TrimPrefix(action, "FORCE-")
	}
	op, ok := powerActionMap[action]
	if !ok {
		return fmt.Errorf("the powerMap does not contain option (%s)", action)
	}
	jobResp, err := powers.PowerAction(client, powerOpts, op).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("doing power action (%s) for instance (%s) failed: %s", action, d.Id(), err)
	}

	// The time of the power on/off and reboot is usually between 15 and 35 seconds.
	timeout := 3 * time.Minute
	if err := cloudservers.WaitForJobSuccess(client, int(timeout/time.Second), jobResp.JobID); err != nil {
		return fmt.Errorf("waiting power action (%s) for instance (%s) failed: %s", action, d.Id(), err)
	}
	return nil
}

func disableSourceDestCheck(networkClient *golangsdk.ServiceClient, portID string) error {
	// Update the allowed-address-pairs of the port to 1.1.1.1/0
	// to disable the source/destination check
	portpairs := ports.UpdateOpts{
		AllowedAddressPairs: []ports.AddressPair{
			{
				IpAddress: "1.1.1.1/0",
			},
		},
	}

	_, err := ports.Update(networkClient, portID, portpairs)
	return err
}

func enableSourceDestCheck(networkClient *golangsdk.ServiceClient, portID string) error {
	// cancle all allowed-address-pairs to enable the source/destination check
	portpairs := ports.UpdateOpts{
		AllowedAddressPairs: make([]ports.AddressPair, 0),
	}

	_, err := ports.Update(networkClient, portID, portpairs)
	return err
}

func updateSourceDestCheck(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var err error

	networks := d.Get("network").([]interface{})
	for i, v := range networks {
		nic := v.(map[string]interface{})
		nicPort := nic["port"].(string)
		if nicPort == "" {
			continue
		}

		if d.HasChange(fmt.Sprintf("network.%d.source_dest_check", i)) {
			sourceDestCheck := nic["source_dest_check"].(bool)
			if !sourceDestCheck {
				err = disableSourceDestCheck(client, nicPort)
			} else {
				err = enableSourceDestCheck(client, nicPort)
			}

			if err != nil {
				return fmt.Errorf("error updating source_dest_check on port(%s) of instance(%s) failed: %s", nicPort, d.Id(), err)
			}
		}
	}

	return nil
}

func calcUnsubscribeResources(d *schema.ResourceData, cfg *config.Config) ([]string, error) {
	var mainResources = []string{d.Id()}

	if shouldUnsubscribeEIP(d) {
		region := cfg.GetRegion(d)
		eipClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return nil, fmt.Errorf("error creating networking client: %s", err)
		}

		eipAddr := d.Get("public_ip").(string)
		eipID, err := common.GetEipIDbyAddress(eipClient, eipAddr, "all_granted_eps")
		if err != nil {
			return nil, fmt.Errorf("error fetching EIP ID of ECS (%s): %s", d.Id(), err)
		}

		mainResources = append(mainResources, eipID)
	}

	return mainResources, nil
}

func shouldUnsubscribeEIP(d *schema.ResourceData) bool {
	deleteEIP := d.Get("delete_eip_on_termination").(bool)
	eipAddr := d.Get("public_ip").(string)
	eipType := d.Get("eip_type").(string)
	_, sharebw := d.GetOk("bandwidth.0.id")

	return deleteEIP && eipAddr != "" && eipType != "" && !sharebw
}

func flattenTagsToMap(d *schema.ResourceData, tags []string) map[string]string {
	result := make(map[string]string)

	tagsRaw := d.Get("tags").(map[string]interface{})

	if len(tagsRaw) != 0 {
		for _, tagStr := range tags {
			tag := strings.Split(tagStr, "=")
			if _, ok := tagsRaw[tag[0]]; ok {
				if len(tag) == 1 {
					result[tag[0]] = ""
				} else if len(tag) == 2 {
					result[tag[0]] = tag[1]
				}
			}
		}
	} else {
		for _, tagStr := range tags {
			tag := strings.Split(tagStr, "=")
			if len(tag) == 1 {
				result[tag[0]] = ""
			} else if len(tag) == 2 {
				result[tag[0]] = tag[1]
			}
		}
	}

	return result
}

func buildInstanceRootVolume(d *schema.ResourceData) cloudservers.RootVolume {
	diskType := d.Get("system_disk_type").(string)
	if diskType == "" {
		diskType = SystemDiskType
	}
	volRequest := cloudservers.RootVolume{
		VolumeType: diskType,
		Size:       d.Get("system_disk_size").(int),
		IOPS:       d.Get("system_disk_iops").(int),
		Throughput: d.Get("system_disk_throughput").(int),
	}

	if v, ok := d.GetOk("system_disk_kms_key_id"); ok {
		matadata := cloudservers.VolumeMetadata{
			SystemEncrypted: "1",
			SystemCmkid:     v.(string),
		}
		volRequest.Metadata = &matadata
	}

	if v, ok := d.GetOk("system_disk_dss_pool_id"); ok {
		volRequest.ClusterType = "DSS"
		volRequest.ClusterId = v.(string)
	}

	return volRequest
}

func buildInstanceDataVolumes(d *schema.ResourceData) []cloudservers.DataVolume {
	var volRequests []cloudservers.DataVolume

	vols := d.Get("data_disks").([]interface{})
	for i := range vols {
		vol := vols[i].(map[string]interface{})
		volRequest := cloudservers.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
			IOPS:       vol["iops"].(int),
			Throughput: vol["throughput"].(int),
		}

		if vol["snapshot_id"] != "" {
			extendparam := cloudservers.VolumeExtendParam{
				SnapshotId: vol["snapshot_id"].(string),
			}
			volRequest.Extendparam = &extendparam
		}

		if vol["kms_key_id"] != "" {
			matadata := cloudservers.VolumeMetadata{
				SystemEncrypted: "1",
				SystemCmkid:     vol["kms_key_id"].(string),
			}
			volRequest.Metadata = &matadata
		}

		if vol["dss_pool_id"] != "" {
			volRequest.ClusterType = "DSS"
			volRequest.ClusterId = vol["dss_pool_id"].(string)
		}

		volRequests = append(volRequests, volRequest)
	}
	return volRequests
}
