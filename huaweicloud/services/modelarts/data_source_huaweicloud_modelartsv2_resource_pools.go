package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ServiceStage GET /v2/{project_id}/pools
func DataSourceV2ResourcePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2ResourcePoolsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the resource pools are located.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID to which the resource pool belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the resource pools to be queried.`,
			},
			"resource_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata": {
							Type:        schema.TypeList,
							Elem:        dataV2ResourcePoolMetadataSchema(),
							Computed:    true,
							Description: `The metadata configuration of the resource pool.`,
						},
						"spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resources": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        dataV2ResourcePoolSpecResourceSchema(),
										Description: `The list of resource specifications in the resource pool.`,
									},
									"scope": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The list of job types supported by the resource pool.`,
									},
									"network": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The name of the network.`,
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The ID of the VPC.`,
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The ID of the subnet.`,
												},
											},
										},
										Description: `The network of the resource pool.`,
									},
									"user_login": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{Schema: map[string]*schema.Schema{
											"key_pair_name": {
												Type:        schema.TypeString,
												Computed:    true,
												Sensitive:   true,
												Description: `The name of the key pair.`,
											},
										}},
										Description: `The user login information of the privileged pool.`,
									},
									"clusters": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{Schema: map[string]*schema.Schema{
											"provider_id": {
												Type:        schema.TypeString,
												Computed:    true,
												Description: `The provider ID of the cluster.`,
											},
											"name": {
												Type:        schema.TypeString,
												Computed:    true,
												Description: `The name of the cluster.`,
											},
										}},
										Description: `The cluster information of the privileged pool.`,
									},
								},
							},
							Description: `The specification of the resource pool.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the resource pool.`,
						},
						// Internal attributes.
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The name of the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"scope": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: utils.SchemaDesc(
								`The list of job types supported by the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataV2ResourcePoolResourceSchema(),
							Description: utils.SchemaDesc(
								`The list of resource specifications in the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The ModelArts network ID of the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The prefix of the user-defined node name of the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The ID of the VPC to which the resource pool belongs.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The network ID of the subnet to which the resource pool belongs.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataV2ResourcePoolClustersSchema(),
							Description: utils.SchemaDesc(
								`The list of the CCE clusters.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"user_login": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataV2ResourcePoolUserLoginSchema(),
							Description: utils.SchemaDesc(
								`The user login info of the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The workspace ID of the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The description of the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"charging_mode": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The charging mode of the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"resource_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(`The resource ID of the resource pool.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
					Description: "All resource pools that match the filter parameters.",
				},
				Description: `All resource pools that matched filter parameters.`,
			},
		},
	}
}

func dataV2ResourcePoolSpecResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor of the resource pool.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The count of the resource pool.`,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The max number of resources of the corresponding flavors.`,
			},
			"node_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of resource pool nodes.`,
			},
			"taints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `The key of the taint.`,
					},
					"effect": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `The effect of the taint.`,
					},
					"value": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `The value of the taint.`,
					},
				}},
				Description: `The taint list of the resource pool.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The key/value pairs labels of resource pool.`,
			},
			"tags": common.TagsComputedSchema(
				`The key/value pairs to associate with the resource pool nodes.`,
			),
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the VPC.`,
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the subnet.`,
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The ID list of the security group.`,
						},
					},
				},
				Description: `The network of the privileged pool.`,
			},
			"extend_params": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extend params of the resource pool, in JSON format.`,
			},
			"creating_step": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The creation step of the resource pool nodes.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the resource pool nodes.`,
						},
					},
				},
				Description: `The creation step configuration of the resource pool nodes.`,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the root volume.`,
						},
						"size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The size of the root volume.`,
						},
					},
				},
				Description: `The root volume of the resource pool nodes.`,
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the data volume.`,
						},
						"size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The size of the data volume.`,
						},
						"extend_params": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The extend parameters of the data volume, in JSON format.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The count of the current data volume configuration.`,
						},
					},
				},
				Description: `The data volumes of the resource pool nodes.`,
			},
			"volume_group_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the volume group.`,
						},
						"docker_thin_pool": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The percentage of container volumes to data volumes on resource pool nodes.`,
						},
						"lvm_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lv_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The LVM write mode.`,
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The volume mount path.`,
									},
								},
							},
							Description: `The configuration of the LVM management.`,
						},
						"types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The storage types of the volume group.`,
						},
					},
				},
				Description: `The extend configurations of the volume groups.`,
			},
			"os": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OS name of the image.`,
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image ID.`,
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image type.`,
						},
					},
				},
				Description: `The image information for the specified OS.`,
			},
			"azs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"az": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The AZ name`,
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of nodes in the AZ.`,
						},
					},
				},
				Description: `The AZ list of the resource pool nodes.`,
			},
		},
	}
}

func dataV2ResourcePoolMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the resource pool.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The annotations of the resource pool.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the resource pool.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the resource pool, in RFC3339 format.`,
			},
		},
	}
}

func dataV2ResourcePoolResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The resource flavor ID.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: utils.SchemaDesc(
					`The number of resources of the corresponding flavors.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"node_pool": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The name of resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"max_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: utils.SchemaDesc(
					`The max number of resources of the corresponding flavors.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The ID of the VPC to which the the resource pool nodes belong.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The network ID of a subnet to which the the resource pool nodes belong.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc(
					`The security group IDs to which the the resource pool nodes belong.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"azs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataV2ResourcePoolResourceAzsSchema(),
				Description: utils.SchemaDesc(
					`The availability zones for the resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"taints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataV2ResourcePoolResourceTaintsSchema(),
				Description: utils.SchemaDesc(
					`The taints added to resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc(
					`The labels of resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"tags": common.TagsComputedSchema(
				utils.SchemaDesc(
					`The key/value pairs to associate with the resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			),
			"extend_params": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The extend parameters of the resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"root_volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataV2ResourcePoolResourceRootVolumeSchema(),
				Description: utils.SchemaDesc(
					`The root volume of the resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataV2ResourcePoolResourceDataVolumesSchema(),
				Description: utils.SchemaDesc(
					`The list of data volumes of the resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"volume_group_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataV2ResourcePoolResourceVolumeGroupConfigsSchema(),
				Description: utils.SchemaDesc(
					`The extend configurations of the volume groups.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"creating_step": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"step": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: utils.SchemaDesc(
								`The creation step of the resource pool nodes.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The type of the resource pool nodes.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The creation step configuration of the resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func dataV2ResourcePoolResourceAzsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"az": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The AZ name.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: utils.SchemaDesc(
					`The number of nodes for the corresponding AZ.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func dataV2ResourcePoolResourceTaintsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The key of the taint.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The value of the taint.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"effect": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The effect of the taint.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func dataV2ResourcePoolResourceRootVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The type of the root volume.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The size of the root volume.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func dataV2ResourcePoolResourceDataVolumesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The type of the data volume.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The size of the data volume.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"extend_params": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The extend parameters of the data volume.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: utils.SchemaDesc(
					`The count of the current data volume configuration.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func dataV2ResourcePoolResourceVolumeGroupConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_group": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The name of the volume group.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"docker_thin_pool": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: utils.SchemaDesc(
					`The percentage of container volumes to data volumes on resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"lvm_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lv_type": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The LVM write mode.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The volume mount path.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The configuration of the LVM management.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc(
					`The list of storage types of the volume group.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func dataV2ResourcePoolClustersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The ID of the CCE cluster that resource pool used.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The name of the CCE cluster that resource pool used.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func dataV2ResourcePoolUserLoginSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The key pair name of the login user.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildV2ResourcePoolsQueryParams(d *schema.ResourceData) string {
	result := ""

	if workspaceId, ok := d.GetOk("workspace_id"); ok {
		result += fmt.Sprintf("&workspaceId=%v", workspaceId)
	}

	if status, ok := d.GetOk("status"); ok {
		result += fmt.Sprintf("&status=%v", status)
	}

	if result == "" {
		return result
	}
	return "?" + result[1:]
}

func listV2ResourcePools(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/pools"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildV2ResourcePoolsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenV2ResourcePoolMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":        utils.PathSearch("name", metadata, nil),
			"annotations": utils.PathSearch("annotations", metadata, make(map[string]interface{})),
			"labels":      utils.PathSearch("labels", metadata, make(map[string]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("creationTimestamp",
				metadata, "").(string))/1000, false),
		},
	}
}

func flattenV2ResourcePoolResourceAzs(azList []interface{}) []map[string]interface{} {
	if len(azList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(azList))
	for _, az := range azList {
		result = append(result, map[string]interface{}{
			"az":    utils.PathSearch("az", az, nil),
			"count": utils.PathSearch("count", az, nil),
		})
	}

	return result
}

func flattenV2ResourcePoolResourceTaints(taints []interface{}) []map[string]interface{} {
	if len(taints) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(taints))
	for _, taint := range taints {
		result = append(result, map[string]interface{}{
			"key":    utils.PathSearch("key", taint, nil),
			"value":  utils.PathSearch("value", taint, nil),
			"effect": utils.PathSearch("effect", taint, nil),
		})
	}

	return result
}

func flattenV2ResourcePoolResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"flavor_id":          utils.PathSearch("flavor", resource, nil),
			"count":              utils.PathSearch("count", resource, nil),
			"max_count":          utils.PathSearch("maxCount", resource, nil),
			"node_pool":          utils.PathSearch("nodePool", resource, nil),
			"vpc_id":             utils.PathSearch("network.vpc", resource, nil),
			"subnet_id":          utils.PathSearch("network.subnet", resource, nil),
			"security_group_ids": utils.PathSearch("network.securityGroups", resource, nil),
			"azs":                flattenV2ResourcePoolResourceAzs(utils.PathSearch("azs", resource, make([]interface{}, 0)).([]interface{})),
			"taints":             flattenV2ResourcePoolResourceTaints(utils.PathSearch("taints", resource, make([]interface{}, 0)).([]interface{})),
			"labels":             utils.PathSearch("labels", resource, nil),
			"tags":               flattenResourcePoolResourcesTags(resource),
			"extend_params":      utils.JsonToString(utils.PathSearch("extendParams", resource, nil)),
			"root_volume":        flattenResourcePoolResourcesRootVolume(utils.PathSearch("rootVolume", resource, nil)),
			"data_volumes": flattenResourcePoolResourcesDataVolumes(utils.PathSearch("dataVolumes",
				resource, make([]interface{}, 0)).([]interface{})),
			"volume_group_configs": flattenResourcePoolResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
				resource, make([]interface{}, 0)).([]interface{})),
			"creating_step": flattenResourcePoolResourcesCreatingStep(utils.PathSearch("creatingStep", resource, nil)),
		})
	}

	return result
}

func flattenV2ResourcePoolClusters(clusters []interface{}) []map[string]interface{} {
	if len(clusters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		result = append(result, map[string]interface{}{
			"provider_id": utils.PathSearch("providerId", cluster, nil),
			"name":        utils.PathSearch("name", cluster, nil),
		})
	}

	return result
}

func flattenV2ResourcePoolUserLogin(userLogin interface{}) []map[string]interface{} {
	if userLogin == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"key_pair_name": utils.PathSearch("keyPairName", userLogin, nil),
		},
	}
}

func flattenV2ResourcePoolChargingMode(chargingMode string) interface{} {
	if chargingMode == "1" {
		return "prePaid"
	}

	return nil
}

func flattenV2DataResourcePoolResourcesOs(osInfo interface{}) []map[string]interface{} {
	if osInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":       utils.PathSearch("name", osInfo, nil),
			"image_id":   utils.PathSearch("imageId", osInfo, nil),
			"image_type": utils.PathSearch("iamgeType", osInfo, nil),
		},
	}
}

func flattenV2DataResourcePoolVolumeGroupConfigsLvmConfig(lvmConfig interface{}) []map[string]interface{} {
	if lvmConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"lv_type": utils.PathSearch("lvType", lvmConfig, nil),
			"path":    utils.PathSearch("path", lvmConfig, nil),
		},
	}
}

func flattenResourcePoolSpecResourcesVolumeGroupConfigs(volumeGroupConfigs []interface{}) []map[string]interface{} {
	if len(volumeGroupConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(volumeGroupConfigs))
	for _, volumeGroupConfig := range volumeGroupConfigs {
		result = append(result, map[string]interface{}{
			"volume_group":     utils.PathSearch("volumeGroup", volumeGroupConfig, nil),
			"docker_thin_pool": utils.PathSearch("dockerThinPool", volumeGroupConfig, nil),
			"lvm_config": flattenV2DataResourcePoolVolumeGroupConfigsLvmConfig(utils.PathSearch("lvmConfig",
				volumeGroupConfig, nil)),
			"types": utils.PathSearch("types", volumeGroupConfig, make([]interface{}, 0)),
		})
	}
	return result
}

func flattenV2DataResourcePoolResourcesCreatingStep(creatingStep interface{}) []map[string]interface{} {
	if creatingStep == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"step": utils.PathSearch("step", creatingStep, nil),
			"type": utils.PathSearch("type", creatingStep, nil),
		},
	}
}

func flattenV2DataResourcePoolResourcesRootVolume(rootVolume interface{}) []map[string]interface{} {
	if rootVolume == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"volume_type": utils.PathSearch("volumeType", rootVolume, nil),
			"size":        utils.PathSearch("size", rootVolume, nil),
		},
	}
}

func flattenV2DataResourcePoolResourcesDataVolumes(dataVolumes []interface{}) []map[string]interface{} {
	if len(dataVolumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataVolumes))
	for _, dataVolume := range dataVolumes {
		result = append(result, map[string]interface{}{
			"volume_type":   utils.PathSearch("volumeType", dataVolume, nil),
			"size":          utils.PathSearch("size", dataVolume, nil),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", dataVolume, nil)),
			"count":         utils.PathSearch("count", dataVolume, nil),
		})
	}

	return result
}

func flattenV2DataResourcePoolSpecResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"flavor":        utils.PathSearch("flavor", resource, nil),
			"count":         utils.PathSearch("count", resource, nil),
			"max_count":     utils.PathSearch("maxCount", resource, nil),
			"node_pool":     utils.PathSearch("nodePool", resource, nil),
			"taints":        flattenV2ResourcePoolResourceTaints(utils.PathSearch("taints", resource, make([]interface{}, 0)).([]interface{})),
			"labels":        utils.PathSearch("labels", resource, nil),
			"tags":          utils.FlattenTagsToMap(utils.PathSearch("tags", resource, nil)),
			"network":       flattenV2DataSourceResourcePoolNetwork(utils.PathSearch("network", resource, nil)),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", resource, nil)),
			"creating_step": flattenV2DataResourcePoolResourcesCreatingStep(utils.PathSearch("creatingStep", resource, nil)),
			"root_volume":   flattenV2DataResourcePoolResourcesRootVolume(utils.PathSearch("rootVolume", resource, nil)),
			"data_volumes": flattenV2DataResourcePoolResourcesDataVolumes(utils.PathSearch("dataVolumes", resource,
				make([]interface{}, 0)).([]interface{})),
			"volume_group_configs": flattenResourcePoolSpecResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
				resource, make([]interface{}, 0)).([]interface{})),
			"os":  flattenV2DataResourcePoolResourcesOs(utils.PathSearch("os", resource, nil)),
			"azs": flattenV2ResourcePoolResourceAzs(utils.PathSearch("azs", resource, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenV2DataSourceResourcePoolNetwork(network interface{}) []map[string]interface{} {
	if network == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"vpc":             utils.PathSearch("vpc", network, nil),
			"subnet":          utils.PathSearch("subnet", network, nil),
			"security_groups": utils.PathSearch("securityGroups", network, nil),
		},
	}
}

func flattenV2DataSourceResourcePoolSpecNetwork(network interface{}) []map[string]interface{} {
	if network == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":      utils.PathSearch("name", network, nil),
			"vpc_id":    utils.PathSearch("vpcId", network, nil),
			"subnet_id": utils.PathSearch("subnetId", network, nil),
		},
	}
}

func flattenV2ResourcePoolSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"resources":  flattenV2DataResourcePoolSpecResources(utils.PathSearch("resources", spec, make([]interface{}, 0)).([]interface{})),
			"scope":      utils.PathSearch("scope", spec, nil),
			"network":    flattenV2DataSourceResourcePoolSpecNetwork(utils.PathSearch("network", spec, nil)),
			"user_login": flattenV2ResourcePoolUserLogin(utils.PathSearch("userLogin", spec, nil)),
			"clusters":   flattenV2ResourcePoolClusters(utils.PathSearch("clusters", spec, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenV2ResourcePools(resourcePools []interface{}) []map[string]interface{} {
	if len(resourcePools) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resourcePools))
	for _, resourcePool := range resourcePools {
		result = append(result, map[string]interface{}{
			"metadata": flattenV2ResourcePoolMetadata(utils.PathSearch("metadata", resourcePool, nil)),
			"spec":     flattenV2ResourcePoolSpec(utils.PathSearch("spec", resourcePool, nil)),
			// Internal attributes.
			"name":  utils.PathSearch(`metadata.labels."os.modelarts/name"`, resourcePool, nil),
			"scope": utils.PathSearch("spec.scope", resourcePool, nil),
			"resources": flattenV2ResourcePoolResources(utils.PathSearch("spec.resources",
				resourcePool, make([]interface{}, 0)).([]interface{})),
			"network_id": utils.PathSearch("spec.network.name", resourcePool, nil),
			"prefix":     utils.PathSearch(`metadata.labels."os.modelarts/node.prefix"`, resourcePool, nil),
			"vpc_id":     utils.PathSearch("spec.network.vpcId", resourcePool, nil),
			"subnet_id":  utils.PathSearch("spec.network.subnetId", resourcePool, nil),
			"clusters": flattenV2ResourcePoolClusters(utils.PathSearch("spec.clusters",
				resourcePool, make([]interface{}, 0)).([]interface{})),
			"user_login":   flattenV2ResourcePoolUserLogin(utils.PathSearch("spec.userLogin", resourcePool, nil)),
			"workspace_id": utils.PathSearch(`metadata.labels."os.modelarts/workspace.id"`, resourcePool, nil),
			"description":  utils.PathSearch(`metadata.annotations."os.modelarts/description"`, resourcePool, nil),
			"charging_mode": flattenV2ResourcePoolChargingMode(utils.PathSearch(`metadata.annotations."os.modelarts/billing.mode"`,
				resourcePool, "").(string)),
			"status":           utils.PathSearch("status.phase", resourcePool, nil),
			"resource_pool_id": utils.PathSearch(`metadata.labels."os.modelarts/resource.id"`, resourcePool, nil),
		})
	}

	return result
}

func dataSourceV2ResourcePoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resourcePools, err := listV2ResourcePools(client, d)
	if err != nil {
		return diag.Errorf("error getting resource pools: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resource_pools", flattenV2ResourcePools(resourcePools)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
