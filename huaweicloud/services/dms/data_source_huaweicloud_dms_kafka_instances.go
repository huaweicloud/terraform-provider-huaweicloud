package dms

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDmsKafkaInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsKafkaInstances,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fuzzy_match": {
				Type:     schema.TypeBool,
				Optional: true,
				RequiredWith: []string{
					"name",
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_failure": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// Attributes
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_space": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manager_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_begin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_end": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_public_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"public_ip_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"public_conn_addresses": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dumping": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_auto_topic": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"partition_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ssl_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"used_storage_space": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connect_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manegement_connect_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cross_vpc_accesses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lisenter_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"advertised_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"port_id": {
										Type:     schema.TypeString,
										Computed: true,
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

func isPubilcIPEnabled(val instances.Instance) bool {
	return val.EnablePublicIP
}

func flattenKafkaInstanceList(client *golangsdk.ServiceClient, conf *config.Config, region string,
	list []instances.Instance) ([]map[string]interface{}, []string, error) {
	ids := make([]string, len(list))
	result := make([]map[string]interface{}, len(list))

	for i, val := range list {
		partitionNum, err := strconv.ParseInt(val.PartitionNum, 10, 64)
		if err != nil {
			return nil, nil, err
		}

		// convert the ids of the availability zone into codes
		availableZoneIDs := val.AvailableZones
		availableZoneCodes, err := getAvailableZoneCodeByID(conf, region, availableZoneIDs)
		if err != nil {
			return nil, nil, err
		}

		instance := map[string]interface{}{
			"id":                         val.InstanceID,
			"type":                       val.Type,
			"name":                       val.Name,
			"description":                val.Description,
			"availability_zones":         availableZoneCodes,
			"enterprise_project_id":      val.EnterpriseProjectID,
			"product_id":                 val.ProductID,
			"engine_version":             val.EngineVersion,
			"storage_spec_code":          val.StorageSpecCode,
			"storage_space":              val.TotalStorageSpace,
			"vpc_id":                     val.VPCID,
			"network_id":                 val.SubnetID,
			"security_group_id":          val.SecurityGroupID,
			"manager_user":               val.KafkaManagerUser,
			"access_user":                val.AccessUser,
			"maintain_begin":             val.MaintainBegin,
			"maintain_end":               val.MaintainEnd,
			"retention_policy":           val.RetentionPolicy,
			"dumping":                    val.ConnectorEnalbe,
			"enable_auto_topic":          val.EnableAutoTopic,
			"partition_num":              partitionNum,
			"ssl_enable":                 val.SslEnable,
			"used_storage_space":         val.UsedStorageSpace,
			"connect_address":            val.ConnectAddress,
			"port":                       val.Port,
			"status":                     val.Status,
			"resource_spec_code":         val.ResourceSpecCode,
			"user_id":                    val.UserID,
			"user_name":                  val.UserName,
			"manegement_connect_address": val.ManagementConnectAddress,
			"tags":                       utils.TagsToMap(val.Tags),
		}
		if isPubilcIPEnabled(val) {
			addrList := strings.Split(strings.TrimSpace(val.PublicConnectionAddress), ",")
			log.Printf("[DEBUG] The address list is: %v", addrList)

			publicIPs := make([]string, len(addrList))
			re := regexp.MustCompile(`(.*):\d+`)
			for i, val := range addrList {
				resp := re.FindStringSubmatch(val)
				if len(resp) < 2 {
					return nil, nil, fmt.Errorf("wrong public IP format, want '{public IP}:{port}', but '%v'", val)
				}
				publicIPs[i] = resp[1]
			}
			instance["public_ip_ids"] = publicIPs
			instance["enable_public_ip"] = val.EnablePublicIP
			instance["public_conn_addresses"] = strings.TrimSpace(val.PublicConnectionAddress)
		}

		crossVpcAccess, err := flattenConnectPorts(val.CrossVpcInfo)
		if err != nil {
			return nil, nil, fmt.Errorf("error retriving details of the cross-VPC information: %v", err)
		}
		instance["cross_vpc_accesses"] = crossVpcAccess

		result[i] = instance
		ids[i] = val.InstanceID
	}

	return result, ids, nil
}

func dataSourceDmsKafkaInstances(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	client, err := conf.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DMS instance client: %s", err)
	}

	isFuzzyMatch := d.Get("fuzzy_match").(bool)
	opt := instances.ListOpts{
		InstanceId:          d.Get("instance_id").(string),
		Name:                d.Get("name").(string),
		ExactMatchName:      strconv.FormatBool(!isFuzzyMatch),
		Engine:              "kafka",
		Status:              d.Get("status").(string),
		IncludeFailure:      strconv.FormatBool(d.Get("include_failure").(bool)),
		EnterpriseProjectID: conf.GetEnterpriseProjectID(d),
	}
	pages, err := instances.List(client, opt).AllPages()
	if err != nil {
		return diag.Errorf("error querying DMS kafka instance list：%v", err)
	}
	list, err := instances.ExtractInstances(pages)
	if err != nil {
		return diag.Errorf("error parsing DMS kafka instance list：%v", err)
	}

	log.Printf("[DEBUG] The result of the DMS kafka instance list query is: %v", list)
	result, ids, err := flattenKafkaInstanceList(client, conf, region, list.Instances)
	if err != nil {
		return diag.Errorf("error flattening DMS kafka instance list：%v", err)
	}
	d.SetId(hashcode.Strings(ids))

	return diag.FromErr(d.Set("instances", result))
}
