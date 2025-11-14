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
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage type of the instance.`,
						},
						"storage_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage resource ID of the instance.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC ID to which the instance belongs.`,
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC name to which the instance belongs.`,
						},
						"vpc_client_plain": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the intra-VPC plaintext access is enabled.`,
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
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The security group name associated with the instance.`,
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
						"connector_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dump task.`,
						},
						"connector_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of dump nodes.`,
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
						"specification": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The specification of the instance.`,
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
						"broker_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The broker numbers of the instance.`,
						},
						"ces_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CES version corresponding to the instance.`,
						},
						"charging_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The charging mode of the instance.`,
						},
						"enable_log_collection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether log collection is enabled.`,
						},
						"extend_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of times the instance has been extended.`,
						},
						"ipv6_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether IPv6 is enabled.`,
						},
						"ipv6_connect_addresses": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The IPv6 connect addresses of the instance.`,
						},
						"is_logical_volume": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `whether the expansion is new instance.`,
						},
						"message_query_inst_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether message query is enabled.`,
						},
						"new_auth_cert": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the new auth cert is enabled.`,
						},
						"new_spec_billing_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `whether the new billing specification is enabled.`,
						},
						"node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of nodes of the instance.`,
						},
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The order ID of the instance.`,
						},
						"pod_connect_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The connection address on the tenant side.`,
						},
						"port_protocol": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private_plain_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether private plaintext access is enabled.`,
									},
									"private_sasl_ssl_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether private SASL SSL access is enabled.`,
									},
									"private_sasl_plaintext_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether private SASL plaintext access is enabled.`,
									},
									"public_plain_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether public plaintext access is enabled.`,
									},
									"public_sasl_ssl_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether public SASL SSL access is enabled.`,
									},
									"public_sasl_plaintext_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether public SASL plaintext access is enabled.`,
									},
									"private_plain_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The address of the private plaintext access.`,
									},
									"private_plain_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain name of the private plaintext access.`,
									},
									"private_sasl_ssl_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The address of the private SASL SSL access.`,
									},
									"private_sasl_ssl_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain name of the private SASL SSL access.`,
									},
									"private_sasl_plaintext_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The address of the private SASL plaintext access.`,
									},
									"private_sasl_plaintext_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain name of the private SASL plaintext access.`,
									},
									"public_plain_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The address of the public plaintext access.`,
									},
									"public_plain_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain name of the public plaintext access.`,
									},
									"public_sasl_ssl_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The address of the public SASL SSL access.`,
									},
									"public_sasl_ssl_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain name of the public SASL SSL access.`,
									},
									"public_sasl_plaintext_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The address of the public SASL plaintext access.`,
									},
									"public_sasl_plaintext_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain name of the public SASL plaintext access.`,
									},
								},
							},
							Description: `The port protocol information of the Kafka instance.`,
						},
						"support_features": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The support features of the instance.`,
						},
						"ssl_two_way_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the two-way authentication is enabled.`,
						},
						"public_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The public bandwidth of the instance.`,
						},
						"public_boundwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The public boundwidth of the instance.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the instance, in RFC3339 format.`,
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

		createdAt, err := strconv.Atoi(utils.PathSearch("created_at", val, "").(string))
		if err != nil {
			log.Printf("[ERROR] error converting created_at to int: %v", err)
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
			"storage_type":               utils.PathSearch("storage_type", val, nil),
			"storage_resource_id":        utils.PathSearch("storage_resource_id", val, nil),
			"vpc_id":                     utils.PathSearch("vpc_id", val, nil),
			"vpc_name":                   utils.PathSearch("vpc_name", val, nil),
			"vpc_client_plain":           utils.PathSearch("vpc_client_plain", val, nil),
			"network_id":                 utils.PathSearch("subnet_id", val, nil),
			"security_group_id":          utils.PathSearch("security_group_id", val, nil),
			"security_group_name":        utils.PathSearch("security_group_name", val, nil),
			"manager_user":               utils.PathSearch("kafka_manager_user", val, nil),
			"access_user":                utils.PathSearch("access_user", val, nil),
			"maintain_begin":             utils.PathSearch("maintain_begin", val, nil),
			"maintain_end":               utils.PathSearch("maintain_end", val, nil),
			"retention_policy":           utils.PathSearch("retention_policy", val, nil),
			"dumping":                    utils.PathSearch("connector_enable", val, nil),
			"connector_id":               utils.PathSearch("connector_id", val, nil),
			"connector_node_num":         utils.PathSearch("connector_node_num", val, nil),
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
			"specification":              utils.PathSearch("specification", val, nil),
			"user_id":                    utils.PathSearch("user_id", val, nil),
			"user_name":                  utils.PathSearch("user_name", val, nil),
			"management_connect_address": utils.PathSearch("management_connect_address", val, nil),
			"tags":                       utils.FlattenTagsToMap(utils.PathSearch("tags", val, nil)),
			"broker_num":                 utils.PathSearch("broker_num", val, nil),
			"ces_version":                utils.PathSearch("ces_version", val, nil),
			"charging_mode":              convertChargingMode(int(utils.PathSearch("charging_mode", val, float64(0)).(float64))),
			"enable_log_collection":      utils.PathSearch("enable_log_collection", val, nil),
			"extend_times":               utils.PathSearch("extend_times", val, nil),
			"ipv6_enable":                utils.PathSearch("ipv6_enable", val, nil),
			"ipv6_connect_addresses":     utils.PathSearch("ipv6_connect_addresses", val, nil),
			"is_logical_volume":          utils.PathSearch("is_logical_volume", val, nil),
			"message_query_inst_enable":  utils.PathSearch("message_query_inst_enable", val, nil),
			"new_auth_cert":              utils.PathSearch("new_auth_cert", val, nil),
			"new_spec_billing_enable":    utils.PathSearch("new_spec_billing_enable", val, nil),
			"node_num":                   utils.PathSearch("node_num", val, nil),
			"order_id":                   utils.PathSearch("order_id", val, nil),
			"pod_connect_address":        utils.PathSearch("pod_connect_address", val, nil),
			"port_protocol":              flattenInstancesPortProtocol(utils.PathSearch("port_protocols", val, nil)),
			"support_features":           utils.PathSearch("support_features", val, nil),
			"ssl_two_way_enable":         utils.PathSearch("ssl_two_way_enable", val, nil),
			"public_bandwidth":           utils.PathSearch("public_bandwidth", val, nil),
			"public_boundwidth":          utils.PathSearch("public_boundwidth", val, nil),
			"created_at":                 utils.FormatTimeStampRFC3339(int64(createdAt)/1000, false),
			// Deprecated attributes.
			"manegement_connect_address": utils.PathSearch("management_connect_address", val, nil),
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

func convertChargingMode(chargingMode int) interface{} {
	switch chargingMode {
	case 0:
		return "prePaid"
	case 1:
		return "postPaid"
	}
	return chargingMode
}

func flattenInstancesPortProtocol(portProtocol interface{}) interface{} {
	return []map[string]interface{}{
		{
			"private_plain_enable":               utils.PathSearch("private_plain_enable", portProtocol, nil),
			"private_plain_address":              utils.PathSearch("private_plain_address", portProtocol, nil),
			"private_plain_domain_name":          utils.PathSearch("private_plain_domain_name", portProtocol, nil),
			"private_sasl_ssl_enable":            utils.PathSearch("private_sasl_ssl_enable", portProtocol, nil),
			"private_sasl_ssl_address":           utils.PathSearch("private_sasl_ssl_address", portProtocol, nil),
			"private_sasl_ssl_domain_name":       utils.PathSearch("private_sasl_ssl_domain_name", portProtocol, nil),
			"private_sasl_plaintext_enable":      utils.PathSearch("private_sasl_plaintext_enable", portProtocol, nil),
			"private_sasl_plaintext_address":     utils.PathSearch("private_sasl_plaintext_address", portProtocol, nil),
			"private_sasl_plaintext_domain_name": utils.PathSearch("private_sasl_plaintext_domain_name", portProtocol, nil),
			"public_plain_enable":                utils.PathSearch("public_plain_enable", portProtocol, nil),
			"public_plain_address":               utils.PathSearch("public_plain_address", portProtocol, nil),
			"public_plain_domain_name":           utils.PathSearch("public_plain_domain_name", portProtocol, nil),
			"public_sasl_ssl_enable":             utils.PathSearch("public_sasl_ssl_enable", portProtocol, nil),
			"public_sasl_ssl_address":            utils.PathSearch("public_sasl_ssl_address", portProtocol, nil),
			"public_sasl_ssl_domain_name":        utils.PathSearch("public_sasl_ssl_domain_name", portProtocol, nil),
			"public_sasl_plaintext_enable":       utils.PathSearch("public_sasl_plaintext_enable", portProtocol, nil),
			"public_sasl_plaintext_address":      utils.PathSearch("public_sasl_plaintext_address", portProtocol, nil),
			"public_sasl_plaintext_domain_name":  utils.PathSearch("public_sasl_plaintext_domain_name", portProtocol, nil),
		},
	}
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
