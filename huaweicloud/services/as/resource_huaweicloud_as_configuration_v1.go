package as

import (
	"regexp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/configurations"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceASConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceASConfigurationCreate,
		Read:   resourceASConfigurationRead,
		Update: nil,
		Delete: resourceASConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: resourceASConfigurationValidateName,
				ForceNew:     true,
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
						},
						"flavor": {
							Type:     schema.TypeString,
							Required: true,
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_name": {
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
						"disk": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"SAS", "SSD", "GPSSD", "ESSD", "SATA",
										}, false),
									},
									"disk_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"SYS", "DATA",
										}, false),
									},
									"kms_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"personality": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 5,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"public_ip": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: resourceASConfigurationValidateIpType,
												},
												"bandwidth": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"size": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: resourceASConfigurationValidateEipBandWidthSize,
															},
															"share_type": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: resourceASConfigurationValidateShareType,
															},
															"charging_mode": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: resourceASConfigurationValidateChargeMode,
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
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
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
		IpType:    eipMap["ip_type"].(string),
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
		ID:          configDataMap["instance_id"].(string),
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
		instanceConfigOpts.PubicIp = publicIps
		logp.Printf("[DEBUG] get publicIps: %#v", publicIps)
	}
	logp.Printf("[DEBUG] get instanceConfig: %#v", instanceConfigOpts)
	return instanceConfigOpts, nil
}

func resourceASConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}
	logp.Printf("[DEBUG] asClient: %#v", asClient)
	configDataMap := d.Get("instance_config").([]interface{})[0].(map[string]interface{})
	logp.Printf("[DEBUG] instance_config is: %#v", configDataMap)
	instanceConfig, err1 := getInstanceConfig(configDataMap)
	if err1 != nil {
		return fmtp.Errorf("Error when getting instance_config info: %s", err1)
	}
	createOpts := configurations.CreateOpts{
		Name:           d.Get("scaling_configuration_name").(string),
		InstanceConfig: instanceConfig,
	}

	logp.Printf("[DEBUG] Create AS configuration Options: %#v", createOpts)
	asConfigId, err := configurations.Create(asClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating ASConfiguration: %s", err)
	}
	logp.Printf("[DEBUG] Create AS Configuration Options: %#v", createOpts)
	d.SetId(asConfigId)
	logp.Printf("[DEBUG] Create AS Configuration %q Success!", asConfigId)
	return resourceASConfigurationRead(d, meta)
}

func resourceASConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}

	asConfig, err := configurations.Get(asClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "AS Configuration")
	}

	logp.Printf("[DEBUG] Retrieved ASConfiguration %q: %+v", d.Id(), asConfig)

	return nil
}

func resourceASConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}
	groups, err1 := getASGroupsByConfiguration(asClient, d.Id())
	if err1 != nil {
		return fmtp.Errorf("Error getting AS groups by configuration ID %q: %s", d.Id(), err1)
	}
	if len(groups) > 0 {
		var groupIds []string
		for _, group := range groups {
			groupIds = append(groupIds, group.ID)
		}
		return fmtp.Errorf("Can not delete the configuration %q, it is used by AS groups %s.", d.Id(), groupIds)
	}
	logp.Printf("[DEBUG] Begin to delete AS configuration %q", d.Id())
	if delErr := configurations.Delete(asClient, d.Id()).ExtractErr(); delErr != nil {
		return fmtp.Errorf("Error deleting AS configuration: %s", delErr)
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

var BandWidthChargeMode = [1]string{"traffic"}

func resourceASConfigurationValidateChargeMode(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range BandWidthChargeMode {
		if value == BandWidthChargeMode[i] {
			return
		}
	}
	errors = append(errors, fmtp.Errorf("%q must be one of %v", k, BandWidthChargeMode))
	return
}

var BandWidthShareType = [1]string{"PER"}

func resourceASConfigurationValidateShareType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range BandWidthShareType {
		if value == BandWidthShareType[i] {
			return
		}
	}
	errors = append(errors, fmtp.Errorf("%q must be one of %v", k, BandWidthShareType))
	return
}

func resourceASConfigurationValidateEipBandWidthSize(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if 1 <= value && value <= 300 {
		return
	}
	errors = append(errors, fmtp.Errorf("%q must be [1, 300], but it is %d", k, value))
	return
}

var IpTypes = [1]string{"5_bgp"}

func resourceASConfigurationValidateIpType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range IpTypes {
		if value == IpTypes[i] {
			return
		}
	}
	errors = append(errors, fmtp.Errorf("%q must be one of %v", k, IpTypes))
	return
}

//lintignore:V001
func resourceASConfigurationValidateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 || len(value) < 1 {
		errors = append(errors, fmtp.Errorf("%q must contain more than 1 and less than 64 characters", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z-_]+$`).MatchString(value) {
		errors = append(errors, fmtp.Errorf("only alphanumeric characters, hyphens, and underscores allowed in %q", k))
	}
	return
}
