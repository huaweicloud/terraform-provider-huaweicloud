package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceRdsInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fixed_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ha_replication_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"param_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"user_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"volume": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_encryption_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"backup_strategy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"keep_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"private_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"public_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud RDS client: %s", err)
	}

	listOpts := instances.ListOpts{
		Name:          d.Get("name").(string),
		Type:          d.Get("type").(string),
		DataStoreType: d.Get("datastore_type").(string),
		VpcId:         d.Get("vpc_id").(string),
		SubnetId:      d.Get("subnet_id").(string),
	}

	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Unable to list instances: %s", err)
	}

	allInstances, err := instances.ExtractRdsInstances(pages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve instances: %s", err)
	}

	filter := map[string]interface{}{
		"EnterpriseProjectId": d.Get("enterprise_project_id").(string),
	}
	filteredInstances, err := utils.FilterSliceWithField(allInstances.Instances, filter)

	if err != nil {
		return fmtp.DiagErrorf("filter RDS flavors failed: %s", err)
	}

	var instancesToSet []map[string]interface{}
	var instancesIds []string

	for _, item := range filteredInstances {
		instanceInAll := item.(instances.RdsInstanceResponse)
		instanceToSet := map[string]interface{}{
			"id":                    instanceInAll.Id,
			"region":                instanceInAll.Region,
			"name":                  instanceInAll.Name,
			"status":                instanceInAll.Status,
			"created":               instanceInAll.Created,
			"ha_replication_mode":   instanceInAll.Ha.ReplicationMode,
			"vpc_id":                instanceInAll.VpcId,
			"subnet_id":             instanceInAll.SubnetId,
			"security_group_id":     instanceInAll.SecurityGroupId,
			"flavor":                instanceInAll.FlavorRef,
			"time_zone":             instanceInAll.TimeZone,
			"enterprise_project_id": instanceInAll.EnterpriseProjectId,
			"tags":                  utils.TagsToMap(instanceInAll.Tags),
		}

		instanceID := instanceInAll.Id
		instancesIds = append(instancesIds, instanceID)

		// publicIps
		publicIps := make([]interface{}, len(instanceInAll.PublicIps))
		for i, v := range instanceInAll.PublicIps {
			publicIps[i] = v
		}
		instanceToSet["public_ips"] = publicIps

		// privateIps
		privateIps := make([]string, len(instanceInAll.PrivateIps))
		for i, v := range instanceInAll.PrivateIps {
			privateIps[i] = v
		}
		instanceToSet["private_ips"] = privateIps

		if len(privateIps) > 0 {
			instanceToSet["fixed_ip"] = privateIps[0]
		}

		// volume
		volume := make([]map[string]interface{}, 1)
		volume[0] = map[string]interface{}{
			"type":               instanceInAll.Volume.Type,
			"size":               instanceInAll.Volume.Size,
			"disk_encryption_id": instanceInAll.DiskEncryptionId,
		}
		instanceToSet["volume"] = volume

		// db
		database := make([]map[string]interface{}, 1)
		database[0] = map[string]interface{}{
			"type":      instanceInAll.DataStore.Type,
			"version":   instanceInAll.DataStore.Version,
			"port":      instanceInAll.Port,
			"user_name": instanceInAll.DbUserName,
		}
		instanceToSet["db"] = database

		// backup
		backup := make([]map[string]interface{}, 1)
		backup[0] = map[string]interface{}{
			"start_time": instanceInAll.BackupStrategy.StartTime,
			"keep_days":  instanceInAll.BackupStrategy.KeepDays,
		}
		instanceToSet["backup_strategy"] = backup

		// nodes
		nodes := make([]map[string]interface{}, len(instanceInAll.Nodes))
		for i, v := range instanceInAll.Nodes {
			nodes[i] = map[string]interface{}{
				"id":                v.Id,
				"name":              v.Name,
				"role":              v.Role,
				"status":            v.Status,
				"availability_zone": v.AvailabilityZone,
			}
		}
		instanceToSet["nodes"] = nodes

		// az
		az1 := instanceInAll.Nodes[0].AvailabilityZone
		if strings.HasSuffix(instanceInAll.FlavorRef, ".ha") {
			if len(instanceInAll.Nodes) < 2 {
				return fmtp.DiagErrorf("[DEBUG] Error saving availability zone to RDS instance (%s): "+
					"HA mode must have two availability zone", instanceID)
			}
			az2 := instanceInAll.Nodes[1].AvailabilityZone
			if instanceInAll.Nodes[1].Role == "master" {
				instanceToSet["availability_zone"] = []string{az2, az1}
			} else {
				instanceToSet["availability_zone"] = []string{az1, az2}
			}
		} else {
			instanceToSet["availability_zone"] = []string{az1}
		}

		instancesToSet = append(instancesToSet, instanceToSet)
	}

	d.SetId(hashcode.Strings(instancesIds))
	d.Set("instances", instancesToSet)

	return nil
}
