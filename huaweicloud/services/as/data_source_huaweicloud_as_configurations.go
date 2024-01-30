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
										Elem:        ConfigurationDiskSchema(),
										Description: "The disk group information of the AS configuration.",
									},
									"personality": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        ConfigurationPersonalitySchema(),
										Description: "The customize personality of the AS configuration.",
									},
									"public_ip": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        ConfigurationPublicIpSchema(),
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
					},
				},
			},
		},
	}
}

func ConfigurationDiskSchema() *schema.Resource {
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
		},
	}
}

func ConfigurationPersonalitySchema() *schema.Resource {
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

func ConfigurationPublicIpSchema() *schema.Resource {
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceASConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	opts := configurations.ListOpts{
		Name:    d.Get("name").(string),
		ImageID: d.Get("image_id").(string),
	}
	page, err := configurations.List(asClient, opts).AllPages()
	if err != nil {
		return diag.Errorf("error getting AS Configuration list: %s", err)
	}

	configurationList, err := page.(configurations.ConfigurationPage).Extract()
	if err != nil {
		return diag.Errorf("error extract to AS Configuration list: %s", err)
	}

	ids := make([]string, 0, len(configurationList))
	elements := make([]map[string]interface{}, 0, len(configurationList))
	for _, configuration := range configurationList {
		configurationMap := map[string]interface{}{
			"scaling_configuration_name": configuration.Name,
			"instance_config":            flattenInstanceConfig(configuration.InstanceConfig),
			"status":                     normalizeConfigurationStatus(configuration.ScalingGroupID),
		}
		ids = append(ids, configuration.ID)
		elements = append(elements, configurationMap)
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("configurations", elements),
		d.Set("region", region),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting AS Configuration fields: %s", mErr)
	}
	return nil
}
