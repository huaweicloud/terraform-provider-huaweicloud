package kafka

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v2/{project_id}/instances
// @API Kafka GET /v2/available-zones
func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the instances are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the instance to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the instance to be queried.`,
			},
			"fuzzy_match": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to perform fuzzy matching query by the name of the instance.`,
				RequiredWith: []string{
					"name",
				},
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the instance belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the instance to be queried.`,
			},
			"include_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the query results contain instances that failed to create.`,
			},
			// Attributes
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance type.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance description.`,
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of AZ names.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID to which the instance belongs.`,
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product ID used by the instance.`,
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The kafka engine version.`,
						},
						"storage_spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage I/O specification.`,
						},
						"storage_space": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The message storage capacity, in GB unit.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC ID to which the instance belongs.`,
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subnet ID to which the instance belongs.`,
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The security group ID associated with the instance.`,
						},
						"manager_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The username for logging in to the Kafka Manager.`,
						},
						"access_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The access username.`,
						},
						"maintain_begin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time at which a maintenance time window starts, the format is HH:mm.`,
						},
						"maintain_end": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time at which a maintenance time window ends, the format is HH:mm.`,
						},
						"enable_public_ip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether public access to the instance is enabled.`,
						},
						"public_ip_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The IDs of the elastic IP address (EIP).`,
						},
						"public_conn_addresses": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance public access address.`,
						},
						"retention_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The action to be taken when the memory usage reaches the disk capacity threshold.`,
						},
						"dumping": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to dumping is enabled.`,
						},
						"enable_auto_topic": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether automatic topic creation is enabled.`,
						},
						"partition_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of topics in the DMS kafka instance.`,
						},
						"ssl_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the Kafka SASL_SSL is enabled.`,
						},
						"security_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The protocol to use after SASL is enabled.`,
						},
						"enabled_mechanisms": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The authentication mechanisms to use after SASL is enabled.`,
						},
						"used_storage_space": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The used message storage space, in GB unit.`,
						},
						"connect_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address for instance connection.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The port number of the instance.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance status.`,
						},
						"resource_spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource specifications identifier.`,
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user ID who created the instance.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The username who created the instance.`,
						},
						"management_connect_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The connection address of the Kafka manager of an instance.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs to associate with the instance.`,
						},
						"cross_vpc_accesses": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the list of cross-VPC access information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the listener IP address.`,
									},
									"advertised_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the advertised IP address.`,
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the port number.`,
									},
									"port_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the port ID.`,
									},
									// Typo, it is only kept in the code, will not be shown in the docs.
									"lisenter_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Deprecated:  "typo in lisenter_ip, please use \"listener_ip\" instead.",
										Description: `Indicates the listener IP address.`,
									},
								},
							},
						},
						// Deprecated attributes.
						// Typo, it is only kept in the code, will not be shown in the docs.
						"manegement_connect_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "typo in manegement_connect_address, please use \"management_connect_address\" instead.",
							Description: `Indicates the management connection address.`,
						},
					},
				},
				Description: `The list of instances that match the filter parameters.`,
			},
		},
	}
}

func flattenInstances(conf *config.Config, region string, instances []interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(instances))

	for i, val := range instances {
		partitionNum, err := strconv.Atoi(utils.PathSearch("partition_num", val, "0").(string))
		if err != nil {
			return nil, err
		}

		// convert the ids of the availability zone into codes
		availableZoneIDs := utils.ExpandToStringList(utils.PathSearch("available_zones", val, make([]interface{}, 0)).([]interface{}))
		availableZoneCodes, err := GetAvailableZoneCodeByID(conf, region, availableZoneIDs)
		if err != nil {
			return nil, err
		}

		instance := map[string]interface{}{
			"id":                         utils.PathSearch("instance_id", val, nil),
			"type":                       utils.PathSearch("type", val, nil),
			"name":                       utils.PathSearch("name", val, nil),
			"description":                utils.PathSearch("description", val, nil),
			"availability_zones":         availableZoneCodes,
			"enterprise_project_id":      utils.PathSearch("enterprise_project_id", val, nil),
			"product_id":                 utils.PathSearch("product_id", val, nil),
			"engine_version":             utils.PathSearch("engine_version", val, nil),
			"storage_spec_code":          utils.PathSearch("storage_spec_code", val, nil),
			"storage_space":              utils.PathSearch("total_storage_space", val, nil),
			"vpc_id":                     utils.PathSearch("vpc_id", val, nil),
			"network_id":                 utils.PathSearch("subnet_id", val, nil),
			"security_group_id":          utils.PathSearch("security_group_id", val, nil),
			"manager_user":               utils.PathSearch("kafka_manager_user", val, nil),
			"access_user":                utils.PathSearch("access_user", val, nil),
			"maintain_begin":             utils.PathSearch("maintain_begin", val, nil),
			"maintain_end":               utils.PathSearch("maintain_end", val, nil),
			"retention_policy":           utils.PathSearch("retention_policy", val, nil),
			"dumping":                    utils.PathSearch("connector_enable", val, nil),
			"enable_auto_topic":          utils.PathSearch("enable_auto_topic", val, nil),
			"partition_num":              partitionNum,
			"ssl_enable":                 utils.PathSearch("ssl_enable", val, nil),
			"security_protocol":          utils.PathSearch("kafka_security_protocol", val, nil),
			"enabled_mechanisms":         utils.PathSearch("sasl_enabled_mechanisms", val, nil),
			"used_storage_space":         utils.PathSearch("used_storage_space", val, nil),
			"connect_address":            utils.PathSearch("connect_address", val, nil),
			"port":                       utils.PathSearch("port", val, nil),
			"status":                     utils.PathSearch("status", val, nil),
			"resource_spec_code":         utils.PathSearch("resource_spec_code", val, nil),
			"user_id":                    utils.PathSearch("user_id", val, nil),
			"user_name":                  utils.PathSearch("user_name", val, nil),
			"management_connect_address": utils.PathSearch("management_connect_address", val, nil),
			"manegement_connect_address": utils.PathSearch("management_connect_address", val, nil),
			"tags":                       utils.FlattenTagsToMap(utils.PathSearch("tags", val, nil)),
		}

		if enablePublicIP := utils.PathSearch("enable_publicip", val, false).(bool); enablePublicIP {
			publicConnectionAddress := strings.TrimSpace(utils.PathSearch("public_connect_address", val, "").(string))
			addrList := strings.Split(publicConnectionAddress, ",")
			log.Printf("[DEBUG] The address list is: %v", addrList)

			publicIPs := make([]string, len(addrList))
			re := regexp.MustCompile(`(.*):\d+`)
			for i, val := range addrList {
				resp := re.FindStringSubmatch(val)
				if len(resp) < 2 {
					return nil, fmt.Errorf("wrong public IP format, want '{public IP}:{port}', but '%v'", val)
				}
				publicIPs[i] = resp[1]
			}
			instance["public_ip_ids"] = publicIPs
			instance["enable_public_ip"] = enablePublicIP
			instance["public_conn_addresses"] = publicConnectionAddress
		}

		crossVpcAccess, err := FlattenCrossVpcInfo(utils.PathSearch("cross_vpc_info", val, "").(string))
		if err != nil {
			return nil, fmt.Errorf("error retrieving details of the cross-VPC information: %v", err)
		}
		instance["cross_vpc_accesses"] = crossVpcAccess

		result[i] = instance
	}

	return result, nil
}

func listInstances(client *golangsdk.ServiceClient, d *schema.ResourceData, epsId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/instances"
		// Default limit is 10, the maximum value is 50.
		limit  = 50
		offset = 0
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listPath += buildInstancesQueryParams(d, epsId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		instances := utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, instances...)
		// The `offset` cannot be greater than or equal to the total number of instances. Otherwise, the response is as follows:
		// {"error_code": "DMS.40050010","error_msg": "Offset parameter is invalid."}
		offset += len(instances)
		if offset >= int(utils.PathSearch("instance_num", respBody, float64(0)).(float64)) {
			break
		}
	}

	return result, nil
}

func buildInstancesQueryParams(d *schema.ResourceData, epsId string) string {
	res := "&engine=kafka"
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}

	rawConfig := d.GetRawConfig()
	if v := utils.GetNestedObjectFromRawConfig(rawConfig, "fuzzy_match"); v != nil {
		res = fmt.Sprintf("%s&exact_match_name=%v", res, !v.(bool))
	} else {
		res = fmt.Sprintf("%s&exact_match_name=true", res)
	}

	if v := utils.GetNestedObjectFromRawConfig(rawConfig, "include_failure"); v != nil {
		res = fmt.Sprintf("%s&include_failure=%v", res, v)
	}
	return res
}

func dataSourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	client, err := conf.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instances, err := listInstances(client, d, conf.GetEnterpriseProjectID(d))
	if err != nil {
		return diag.Errorf("error querying kafka instance list: %s", err)
	}

	result, err := flattenInstances(conf, region, instances)
	if err != nil {
		return diag.Errorf("error flattening kafka instance list: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
