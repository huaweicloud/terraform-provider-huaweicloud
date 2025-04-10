package as

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/configurations"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_configuration/{id}
// @API AS DELETE /autoscaling-api/v1/{project_id}/scaling_configuration/{id}
// @API AS POST /autoscaling-api/v1/{project_id}/scaling_configuration
// @API AS GET /autoscaling-api/v1/{project_id}/scaling_group
func ResourceASConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASConfigurationCreate,
		ReadContext:   resourceASConfigurationRead,
		DeleteContext: resourceASConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "schema: Required",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"charging_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "postPaid",
						},
						"flavor_priority_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"ecs_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"disk_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"kms_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"throughput": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"dedicated_storage_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"data_disk_image_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"snapshot_id": {
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
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
												"bandwidth": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													ForceNew: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"size": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
																ForceNew: true,
															},
															"share_type": {
																Type:     schema.TypeString,
																Required: true,
																ForceNew: true,
															},
															"charging_mode": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
																ForceNew: true,
															},
															"id": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
																ForceNew: true,
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
						"tenancy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"dedicated_host_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							// just stash the hash for state & diff comparisons
							StateFunc:        utils.HashAndHexEncode,
							DiffSuppressFunc: utils.SuppressUserData,
						},
						"admin_pass": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"key_fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDiskOpts(diskMeta []interface{}) []configurations.DiskOpts {
	var diskOptsList []configurations.DiskOpts

	for _, v := range diskMeta {
		disk := v.(map[string]interface{})
		diskOpts := configurations.DiskOpts{
			Size:               disk["size"].(int),
			VolumeType:         disk["volume_type"].(string),
			DiskType:           disk["disk_type"].(string),
			Iops:               disk["iops"].(int),
			Throughput:         disk["throughput"].(int),
			DedicatedStorageID: disk["dedicated_storage_id"].(string),
			DataDiskImageID:    disk["data_disk_image_id"].(string),
			SnapshotId:         disk["snapshot_id"].(string),
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

	return diskOptsList
}

func buildPersonalityOpts(personalityMeta []interface{}) []configurations.PersonalityOpts {
	if len(personalityMeta) == 0 {
		return nil
	}

	personalityOptsList := make([]configurations.PersonalityOpts, len(personalityMeta))
	for i, v := range personalityMeta {
		personality := v.(map[string]interface{})
		personalityOptsList[i] = configurations.PersonalityOpts{
			Path:    personality["path"].(string),
			Content: personality["content"].(string),
		}
	}

	return personalityOptsList
}

func buildPublicIpOpts(publicIpMeta map[string]interface{}) configurations.PublicIpOpts {
	eipMap := publicIpMeta["eip"].([]interface{})[0].(map[string]interface{})
	bandWidthMap := eipMap["bandwidth"].([]interface{})[0].(map[string]interface{})
	bandWidthOpts := configurations.BandwidthOpts{
		Size:         bandWidthMap["size"].(int),
		ShareType:    bandWidthMap["share_type"].(string),
		ChargingMode: bandWidthMap["charging_mode"].(string),
		ID:           bandWidthMap["id"].(string),
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

func buildSecurityGroupIDsOpts(secGroups []interface{}) []configurations.SecurityGroupOpts {
	if len(secGroups) == 0 {
		return nil
	}

	res := make([]configurations.SecurityGroupOpts, len(secGroups))
	for i, v := range secGroups {
		res[i] = configurations.SecurityGroupOpts{
			ID: v.(string),
		}
	}
	return res
}

func buildMetadataOpts(configDataMap map[string]interface{}) map[string]interface{} {
	metadataMap := configDataMap["metadata"].(map[string]interface{})
	adminPass := configDataMap["admin_pass"].(string)
	if metadataMap == nil && adminPass == "" {
		return nil
	}

	resultMap := make(map[string]interface{})
	if adminPass != "" {
		resultMap["admin_pass"] = adminPass
	}

	for k, v := range metadataMap {
		resultMap[k] = v
	}
	return resultMap
}

func buildInstanceConfig(configDataMap map[string]interface{}) configurations.InstanceConfigOpts {
	instanceConfigOpts := configurations.InstanceConfigOpts{
		InstanceID:           configDataMap["instance_id"].(string),
		FlavorRef:            configDataMap["flavor"].(string),
		ImageRef:             configDataMap["image"].(string),
		SSHKey:               configDataMap["key_name"].(string),
		FlavorPriorityPolicy: configDataMap["flavor_priority_policy"].(string),
		ServerGroupID:        configDataMap["ecs_group_id"].(string),
		Tenancy:              configDataMap["tenancy"].(string),
		DedicatedHostID:      configDataMap["dedicated_host_id"].(string),
		UserData:             []byte(configDataMap["user_data"].(string)),
		Metadata:             buildMetadataOpts(configDataMap),
		SecurityGroups:       buildSecurityGroupIDsOpts(configDataMap["security_group_ids"].([]interface{})),
		Personality:          buildPersonalityOpts(configDataMap["personality"].([]interface{})),
		Disk:                 buildDiskOpts(configDataMap["disk"].([]interface{})),
	}

	if mode, ok := configDataMap["charging_mode"]; ok && mode.(string) == "spot" {
		instanceConfigOpts.MarketType = "spot"
	}

	pubicIpData := configDataMap["public_ip"].([]interface{})
	if len(pubicIpData) == 1 {
		publicIpMap := pubicIpData[0].(map[string]interface{})
		publicIps := buildPublicIpOpts(publicIpMap)
		instanceConfigOpts.PubicIP = &publicIps
	}

	return instanceConfigOpts
}

func resourceASConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	configDataMap := d.Get("instance_config").([]interface{})[0].(map[string]interface{})
	instanceConfig := buildInstanceConfig(configDataMap)

	createOpts := configurations.CreateOpts{
		Name:           d.Get("scaling_configuration_name").(string),
		InstanceConfig: instanceConfig,
	}

	asConfigId, err := configurations.Create(asClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating AS configuration: %s", err)
	}

	d.SetId(asConfigId)

	// instance_id is missing from Get response, we should set it to the state before READ method
	if v, ok := d.GetOk("instance_config.0.instance_id"); ok {
		instanceConfig := []map[string]interface{}{
			{"instance_id": v},
		}
		d.Set("instance_config", instanceConfig)
	}
	return resourceASConfigurationRead(ctx, d, meta)
}

func resourceASConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	configId := d.Id()
	asConfig, err := configurations.Get(asClient, configId).Extract()
	if err != nil {
		// When the resource does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving AS configuration")
	}

	// update InstanceID if necessary
	if v, ok := d.GetOk("instance_config.0.instance_id"); ok {
		if asConfig.InstanceConfig.InstanceID == "" {
			asConfig.InstanceConfig.InstanceID = v.(string)
		}
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("scaling_configuration_name", asConfig.Name),
		d.Set("instance_config", flattenInstanceConfig(asConfig.InstanceConfig, d)),
		d.Set("status", normalizeConfigurationStatus(asConfig.ScalingGroupID)),
		d.Set("create_time", asConfig.CreateTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceASConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	asClient, err := conf.AutoscalingV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	configId := d.Id()
	asGroups, err := getASGroupsByConfiguration(asClient, configId)
	if err != nil {
		return diag.Errorf("error getting AS groups by configuration ID %s: %s", configId, err)
	}

	if len(asGroups) > 0 {
		var groupIds []string
		for _, group := range asGroups {
			groupIds = append(groupIds, group.ID)
		}
		return diag.Errorf("can not delete the configuration %s, it is used by AS groups %v", configId, groupIds)
	}

	if delErr := configurations.Delete(asClient, configId).ExtractErr(); delErr != nil {
		// When the resource does not exist, the response HTTP status code of the delete API is 404.
		return common.CheckDeletedDiag(d, delErr, "error deleting AS configuration")
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
		return gs, fmt.Errorf("error getting AS groups by configuration %s: %s", configurationID, err)
	}

	gs, err = page.(groups.GroupPage).Extract()
	return gs, err
}

// In order to avoid triggering force_new changes, write `admin_pass` and `metadata` in local files.
func flattenInstanceConfig(instanceConfig configurations.InstanceConfig, d *schema.ResourceData) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"charging_mode":          normalizeConfigurationChargingMode(instanceConfig.MarketType),
			"instance_id":            instanceConfig.InstanceID,
			"flavor":                 instanceConfig.FlavorRef,
			"image":                  instanceConfig.ImageRef,
			"key_name":               instanceConfig.SSHKey,
			"flavor_priority_policy": instanceConfig.FlavorPriorityPolicy,
			"ecs_group_id":           instanceConfig.ServerGroupID,
			"tenancy":                instanceConfig.Tenancy,
			"dedicated_host_id":      instanceConfig.DedicatedHostID,
			"user_data":              d.Get("instance_config.0.user_data"),
			"admin_pass":             d.Get("instance_config.0.admin_pass"),
			"metadata":               d.Get("instance_config.0.metadata"),
			"disk":                   flattenInstanceDisks(instanceConfig.Disk),
			"public_ip":              flattenInstancePublicIP(instanceConfig.PublicIp.Eip),
			"security_group_ids":     flattenSecurityGroupIDs(instanceConfig.SecurityGroups),
			"personality":            flattenInstancePersonality(instanceConfig.Personality),
			"key_fingerprint":        instanceConfig.KeyFingerprint,
		},
	}
}

func flattenInstanceDisks(disks []configurations.Disk) []map[string]interface{} {
	if len(disks) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(disks))
	for i, item := range disks {
		res[i] = map[string]interface{}{
			"volume_type":          item.VolumeType,
			"size":                 item.Size,
			"disk_type":            item.DiskType,
			"iops":                 item.Iops,
			"throughput":           item.Throughput,
			"dedicated_storage_id": item.DedicatedStorageID,
			"data_disk_image_id":   item.DataDiskImageID,
			"snapshot_id":          item.SnapshotID,
		}

		if kms, ok := item.Metadata["__system__cmkid"]; ok {
			res[i]["kms_id"] = kms
		}
	}
	return res
}

func flattenInstancePublicIP(eipObject configurations.Eip) []map[string]interface{} {
	if eipObject.Type == "" {
		return nil
	}

	bwInfo := []map[string]interface{}{
		{
			"share_type":    eipObject.Bandwidth.ShareType,
			"size":          eipObject.Bandwidth.Size,
			"charging_mode": eipObject.Bandwidth.ChargingMode,
			"id":            eipObject.Bandwidth.ID,
		},
	}

	eipInfo := []map[string]interface{}{
		{
			"ip_type":   eipObject.Type,
			"bandwidth": bwInfo,
		},
	}

	return []map[string]interface{}{
		{"eip": eipInfo},
	}
}

func flattenInstancePersonality(personalities []configurations.Personality) []map[string]interface{} {
	if len(personalities) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(personalities))
	for i, item := range personalities {
		res[i] = map[string]interface{}{
			"path":    item.Path,
			"content": item.Content,
		}
	}
	return res
}

func flattenSecurityGroupIDs(sgs []configurations.SecurityGroup) []string {
	if len(sgs) == 0 {
		return nil
	}

	res := make([]string, len(sgs))
	for i, item := range sgs {
		res[i] = item.ID
	}
	return res
}

func normalizeConfigurationStatus(groupIDs string) string {
	if groupIDs != "" {
		return "Bound"
	}
	return "Unbound"
}

func normalizeConfigurationChargingMode(marketType string) string {
	if marketType == "" {
		return "postPaid"
	}
	return "spot"
}
