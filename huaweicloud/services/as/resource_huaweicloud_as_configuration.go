package as

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/configurations"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

var (
	ValidDiskTypes     = []string{"SYS", "DATA"}
	ValidVolumeTypes   = []string{"SAS", "SSD", "GPSSD", "ESSD", "SATA"}
	ValidEipTypes      = []string{"5_bgp", "5_sbgp"}
	ValidShareTypes    = []string{"PER", "WHOLE"}
	ValidChargingModes = []string{"traffic", "bandwidth"}
)

func ResourceASConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASConfigurationCreate,
		ReadContext:   resourceASConfigurationRead,
		DeleteContext: resourceASConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_configuration_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa50-9a-zA-Z-_]+$"),
						"only letters, digits, underscores (_), and hyphens (-) are allowed"),
				),
			},
			"instance_config": {
				Required: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							AtLeastOneOf: []string{
								"instance_config.0.instance_id", "instance_config.0.flavor",
								"instance_config.0.image", "instance_config.0.disk",
							},
						},
						"flavor": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							RequiredWith: []string{"instance_config.0.image", "instance_config.0.disk"},
						},
						"image": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							RequiredWith: []string{"instance_config.0.flavor", "instance_config.0.disk"},
						},
						"key_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							RequiredWith: []string{"instance_config.0.flavor", "instance_config.0.image"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"volume_type": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringInSlice(ValidVolumeTypes, false),
									},
									"disk_type": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringInSlice(ValidDiskTypes, false),
									},
									"kms_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"personality": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							MaxItems: 5,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"content": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"public_ip": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_type": {
													Type:         schema.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringInSlice(ValidEipTypes, false),
												},
												"bandwidth": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													ForceNew: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"size": {
																Type:         schema.TypeInt,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.IntBetween(1, 2000),
															},
															"share_type": {
																Type:         schema.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.StringInSlice(ValidShareTypes, false),
															},
															"charging_mode": {
																Type:         schema.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.StringInSlice(ValidChargingModes, false),
															},
														},
													},
												},
											},
										},
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
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func getDisk(diskMeta []interface{}) ([]configurations.DiskOpts, error) {
	var diskOptsList []configurations.DiskOpts

	for _, v := range diskMeta {
		disk := v.(map[string]interface{})
		size := disk["size"].(int)
		volumeType := disk["volume_type"].(string)
		diskType := disk["disk_type"].(string)
		if diskType == "SYS" {
			if size < 40 || size > 32768 {
				return diskOptsList, fmtp.Errorf("For system disk size should be [40, 32768]")
			}
		}
		if diskType == "DATA" {
			if size < 10 || size > 32768 {
				return diskOptsList, fmtp.Errorf("For data disk size should be [10, 32768]")
			}
		}
		diskOpts := configurations.DiskOpts{
			Size:       size,
			VolumeType: volumeType,
			DiskType:   diskType,
		}
		kmsId := disk["kms_id"].(string)
		if kmsId != "" {
			m := make(map[string]string)
			m["__system__cmkid"] = kmsId
			m["__system__encrypted"] = "1"
			diskOpts.Metadata = m
		}
		diskOptsList = append(diskOptsList, diskOpts)
	}

	return diskOptsList, nil
}

func getPersonality(personalityMeta []interface{}) []configurations.PersonalityOpts {
	var personalityOptsList []configurations.PersonalityOpts

	for _, v := range personalityMeta {
		personality := v.(map[string]interface{})
		personalityOpts := configurations.PersonalityOpts{
			Path:    personality["path"].(string),
			Content: personality["content"].(string),
		}
		personalityOptsList = append(personalityOptsList, personalityOpts)
	}

	return personalityOptsList
}

func getPublicIps(publicIpMeta map[string]interface{}) configurations.PublicIpOpts {
	eipMap := publicIpMeta["eip"].([]interface{})[0].(map[string]interface{})
	bandWidthMap := eipMap["bandwidth"].([]interface{})[0].(map[string]interface{})
	bandWidthOpts := configurations.BandwidthOpts{
		Size:         bandWidthMap["size"].(int),
		ShareType:    bandWidthMap["share_type"].(string),
		ChargingMode: bandWidthMap["charging_mode"].(string),
	}

	eipOpts := configurations.EipOpts{
		Type:      eipMap["ip_type"].(string),
		Bandwidth: bandWidthOpts,
	}

	publicIpOpts := configurations.PublicIpOpts{
		Eip: eipOpts,
	}

	return publicIpOpts
}

func getInstanceConfig(configDataMap map[string]interface{}) (configurations.InstanceConfigOpts, error) {
	disksData := configDataMap["disk"].([]interface{})
	disks, err := getDisk(disksData)
	if err != nil {
		return configurations.InstanceConfigOpts{}, fmtp.Errorf("Error happened when validating disk size: %s", err)
	}
	logp.Printf("[DEBUG] get disks: %#v", disks)

	personalityData := configDataMap["personality"].([]interface{})
	personalities := getPersonality(personalityData)
	logp.Printf("[DEBUG] get personality: %#v", personalities)

	instanceConfigOpts := configurations.InstanceConfigOpts{
		InstanceID:  configDataMap["instance_id"].(string),
		FlavorRef:   configDataMap["flavor"].(string),
		ImageRef:    configDataMap["image"].(string),
		SSHKey:      configDataMap["key_name"].(string),
		UserData:    []byte(configDataMap["user_data"].(string)),
		Disk:        disks,
		Personality: personalities,
		Metadata:    configDataMap["metadata"].(map[string]interface{}),
	}
	logp.Printf("[DEBUG] instanceConfigOpts: %#v", instanceConfigOpts)
	pubicIpData := configDataMap["public_ip"].([]interface{})
	logp.Printf("[DEBUG] pubicIpData: %#v", pubicIpData)
	// user specify public_ip
	if len(pubicIpData) == 1 {
		publicIpMap := pubicIpData[0].(map[string]interface{})
		publicIps := getPublicIps(publicIpMap)
		instanceConfigOpts.PubicIP = &publicIps
		logp.Printf("[DEBUG] get publicIps: %#v", publicIps)
	}
	logp.Printf("[DEBUG] get instanceConfig: %#v", instanceConfigOpts)
	return instanceConfigOpts, nil
}

func resourceASConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating autoscaling client: %s", err)
	}

	configDataMap := d.Get("instance_config").([]interface{})[0].(map[string]interface{})
	instanceConfig, err := getInstanceConfig(configDataMap)
	if err != nil {
		return diag.Errorf("Error when getting instance_config info: %s", err)
	}
	createOpts := configurations.CreateOpts{
		Name:           d.Get("scaling_configuration_name").(string),
		InstanceConfig: instanceConfig,
	}

	logp.Printf("[DEBUG] Create AS configuration Options: %#v", createOpts)
	asConfigId, err := configurations.Create(asClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating ASConfiguration: %s", err)
	}

	d.SetId(asConfigId)
	return resourceASConfigurationRead(ctx, d, meta)
}

func resourceASConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating autoscaling client: %s", err)
	}

	asConfig, err := configurations.Get(asClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AS Configuration")
	}

	logp.Printf("[DEBUG] Retrieved ASConfiguration %q: %+v", d.Id(), asConfig)

	return nil
}

func resourceASConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating autoscaling client: %s", err)
	}

	groups, err := getASGroupsByConfiguration(asClient, d.Id())
	if err != nil {
		return diag.Errorf("Error getting AS groups by configuration ID %q: %s", d.Id(), err)
	}
	if len(groups) > 0 {
		var groupIds []string
		for _, group := range groups {
			groupIds = append(groupIds, group.ID)
		}
		return diag.Errorf("Can not delete the configuration %q, it is used by AS groups %s.", d.Id(), groupIds)
	}

	logp.Printf("[DEBUG] Begin to delete AS configuration %q", d.Id())
	if delErr := configurations.Delete(asClient, d.Id()).ExtractErr(); delErr != nil {
		return diag.Errorf("Error deleting AS configuration: %s", delErr)
	}

	return nil
}

func getASGroupsByConfiguration(asClient *golangsdk.ServiceClient, configurationID string) ([]groups.Group, error) {
	var gs []groups.Group
	listOpts := groups.ListOpts{
		ConfigurationID:     configurationID,
		EnterpriseProjectID: "all_granted_eps",
	}
	page, err := groups.List(asClient, listOpts).AllPages()
	if err != nil {
		return gs, fmtp.Errorf("Error getting ASGroups by configuration %q: %s", configurationID, err)
	}
	gs, err = page.(groups.GroupPage).Extract()
	return gs, err
}
