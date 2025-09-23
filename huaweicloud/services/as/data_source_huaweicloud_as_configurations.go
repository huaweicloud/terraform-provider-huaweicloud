package as

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/configurations"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_configuration
func DataSourceASConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceASConfigurationRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the AS configurations are located.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The AS configuration name used to query configuration list.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The AS image id used to query configuration list.",
			},
			"configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The information about AS instance configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the AS configuration.",
						},
						"scaling_configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the AS configuration.",
						},
						"instance_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ECS instance ID of the AS configuration.",
									},
									"flavor": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ECS flavor name of the AS configuration.",
									},
									"image": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ECS image ID of the AS configuration.",
									},
									"key_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the SSH key pair used to log in to the instance.",
									},
									"key_fingerprint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The fingerprint of the SSH key pair used to log in to the instance.",
									},
									"tenancy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates creating ECS instance on DEH.",
									},
									"dedicated_host_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the DEH.",
									},
									"security_group_ids": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "An array of one or more security group IDs.",
									},
									"charging_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The billing mode for ECS.",
									},
									"flavor_priority_policy": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "The priority policy used when there are multiple flavors and " +
											"instances to be created using an AS configuration.",
									},
									"ecs_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ECS group ID of the AS configuration.",
									},
									"disk": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        configurationDataSourceDiskSchema(),
										Description: "The disk group information of the AS configuration.",
									},
									"personality": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        configurationDataSourcePersonalitySchema(),
										Description: "The customize personality of the AS configuration.",
									},
									"public_ip": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        configurationDataSourcePublicIpSchema(),
										Description: "The EIP of the ECS instance.",
									},
									"user_data": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user data to provide when launching the instance.",
									},
									"metadata": {
										Type:        schema.TypeMap,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The key/value pairs to make available from within the instance.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the AS configuration.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the AS configuration.",
						},
					},
				},
			},
		},
	}
}

func configurationDataSourceDiskSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The disk size. The unit is GB.",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The volume type.",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The disk type.",
			},
			"kms_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The encryption KMS ID.",
			},
			"dedicated_storage_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DSS device for the disk.",
			},
			"data_disk_image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the data disk image for creating a data disk.",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The disk backup snapshot ID.",
			},
			"iops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The IOPS of an EVS disk.",
			},
			"throughput": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The throughput of an EVS disk.",
			},
		},
	}
}

func configurationDataSourcePersonalitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path of the injected file.",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content of the injected file.",
			},
		},
	}
}

func configurationDataSourcePublicIpSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"eip": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of EIP configuration that will be automatically assigned to the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The EIP type.",
						},
						"bandwidth": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of bandwidth information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The bandwidth (Mbit/s).",
									},
									"share_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The bandwidth sharing type.",
									},
									"charging_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The bandwidth billing mode.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the bandwidth.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildDataSourceConfigurationOpts(d *schema.ResourceData) configurations.ListOpts {
	return configurations.ListOpts{
		Name:    d.Get("name").(string),
		ImageID: d.Get("image_id").(string),
	}
}

func dataSourceASConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf   = meta.(*config.Config)
		region = conf.GetRegion(d)
		opts   = buildDataSourceConfigurationOpts(d)
	)

	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	page, err := configurations.List(asClient, opts).AllPages()
	if err != nil {
		return diag.Errorf("error getting AS configuration list: %s", err)
	}

	configurationList, err := page.(configurations.ConfigurationPage).Extract()
	if err != nil {
		return diag.Errorf("error extract to AS configuration list: %s", err)
	}

	ids, elements := flattenDataSourceConfigurations(configurationList)
	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("configurations", elements),
		d.Set("region", region),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting AS configuration fields: %s", mErr)
	}
	return nil
}

func flattenDataSourceConfigurations(configurationList []configurations.Configuration) ([]string, []map[string]interface{}) {
	ids := make([]string, 0, len(configurationList))
	elements := make([]map[string]interface{}, 0, len(configurationList))
	for _, configuration := range configurationList {
		configurationMap := map[string]interface{}{
			"scaling_configuration_id":   configuration.ID,
			"scaling_configuration_name": configuration.Name,
			"instance_config":            flattenInstanceConfigs(configuration.InstanceConfig),
			"status":                     normalizeDataSourceConfigurationStatus(configuration.ScalingGroupID),
			"create_time":                configuration.CreateTime,
		}
		ids = append(ids, configuration.ID)
		elements = append(elements, configurationMap)
	}

	return ids, elements
}

func flattenInstanceConfigs(instanceConfig configurations.InstanceConfig) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"charging_mode":          normalizeAsConfigurationChargingMode(instanceConfig.MarketType),
			"instance_id":            instanceConfig.InstanceID,
			"flavor":                 instanceConfig.FlavorRef,
			"image":                  instanceConfig.ImageRef,
			"key_name":               instanceConfig.SSHKey,
			"flavor_priority_policy": instanceConfig.FlavorPriorityPolicy,
			"ecs_group_id":           instanceConfig.ServerGroupID,
			"user_data":              instanceConfig.UserData,
			"metadata":               instanceConfig.Metadata,
			"disk":                   flattenAsInstanceDisks(instanceConfig.Disk),
			"public_ip":              flattenAsInstancePublicIP(instanceConfig.PublicIp.Eip),
			"security_group_ids":     flattenAsSecurityGroupIDs(instanceConfig.SecurityGroups),
			"personality":            flattenAsInstancePersonality(instanceConfig.Personality),
			"key_fingerprint":        instanceConfig.KeyFingerprint,
			"tenancy":                instanceConfig.Tenancy,
			"dedicated_host_id":      instanceConfig.DedicatedHostID,
		},
	}
}

func normalizeDataSourceConfigurationStatus(groupIDs string) string {
	if groupIDs != "" {
		return "Bound"
	}
	return "Unbound"
}

func flattenAsInstanceDisks(disks []configurations.Disk) []map[string]interface{} {
	if len(disks) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(disks))
	for i, item := range disks {
		res[i] = map[string]interface{}{
			"volume_type":          item.VolumeType,
			"size":                 item.Size,
			"disk_type":            item.DiskType,
			"dedicated_storage_id": item.DedicatedStorageID,
			"data_disk_image_id":   item.DataDiskImageID,
			"snapshot_id":          item.SnapshotID,
			"iops":                 item.Iops,
			"throughput":           item.Throughput,
		}

		if kms, ok := item.Metadata["__system__cmkid"]; ok {
			res[i]["kms_id"] = kms
		}
	}
	return res
}

func flattenAsInstancePublicIP(eipObject configurations.Eip) []map[string]interface{} {
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

func flattenAsInstancePersonality(personalities []configurations.Personality) []map[string]interface{} {
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

func flattenAsSecurityGroupIDs(sgs []configurations.SecurityGroup) []string {
	if len(sgs) == 0 {
		return nil
	}

	res := make([]string, len(sgs))
	for i, item := range sgs {
		res[i] = item.ID
	}
	return res
}

func normalizeAsConfigurationChargingMode(marketType string) string {
	if marketType == "" {
		return "postPaid"
	}
	return "spot"
}
